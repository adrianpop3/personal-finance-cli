package parser

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"personal-finance-cli/db"
)

type ParsedTransaction struct {
	Amount      float64
	Description string
	Date        time.Time
	Category    string
	ExternalID  string // FITID for OFX (empty for CSV)
}

func DetectAndParse(r io.Reader, filename string) ([]ParsedTransaction, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".csv":
		return parseCSV(r)
	case ".ofx", ".qfx":
		return parseOFX(r)
	default:
		br := bufio.NewReader(r)
		peek, _ := br.Peek(2048)
		s := strings.ToLower(string(peek))
		if strings.Contains(s, "<ofx") {
			return parseOFX(br)
		}
		if strings.Contains(s, ",") {
			return parseCSV(br)
		}
		return nil, errors.New("unsupported file format")
	}
}

//
// ------------------ CSV PARSER ------------------
//

func parseCSV(r io.Reader) ([]ParsedTransaction, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true

	records, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	headers := records[0]
	idxDate, idxAmount, idxDesc, idxCat := -1, -1, -1, -1

	for i, h := range headers {
		h = strings.ToLower(strings.TrimSpace(h))
		switch h {
		case "date", "dt":
			idxDate = i
		case "amount", "amt", "value":
			idxAmount = i
		case "description", "desc", "name", "memo":
			idxDesc = i
		case "category", "cat":
			idxCat = i
		}
	}

	if idxDate == -1 && len(headers) >= 1 {
		idxDate = 0
	}
	if idxAmount == -1 && len(headers) >= 2 {
		idxAmount = 1
	}
	if idxDesc == -1 && len(headers) >= 3 {
		idxDesc = 2
	}

	var parsed []ParsedTransaction

	for _, row := range records[1:] {
		if len(row) == 0 {
			continue
		}

		get := func(idx int) string {
			if idx >= 0 && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		dateStr := get(idxDate)
		amtStr := strings.ReplaceAll(get(idxAmount), ",", "")
		desc := get(idxDesc)
		cat := get(idxCat)

		amount, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			continue
		}

		txDate := parseFlexibleDate(dateStr)

		pt := ParsedTransaction{
			Amount:      amount,
			Description: desc,
			Date:        txDate,
			Category:    cat,
			ExternalID:  "",
		}

		if strings.TrimSpace(pt.Category) == "" {
			pt.Category = InferCategory(pt.Description)
		}

		parsed = append(parsed, pt)
	}

	return parsed, nil
}

func parseOFX(r io.Reader) ([]ParsedTransaction, error) {
	scanner := bufio.NewScanner(r)

	var parsed []ParsedTransaction
	var inTxn bool

	var dateStr, amtStr, name, memo, fitid string

	reset := func() {
		inTxn = false
		dateStr, amtStr, name, memo, fitid = "", "", "", "", ""
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineLower := strings.ToLower(line)

		if strings.Contains(lineLower, "<stmttrn") {
			inTxn = true
			continue
		}

		if strings.Contains(lineLower, "</stmttrn") {
			amount, err := strconv.ParseFloat(strings.ReplaceAll(amtStr, ",", ""), 64)
			if err != nil {
				reset()
				continue
			}

			txDate := parseOFXDate(dateStr)

			desc := name
			if desc == "" {
				desc = memo
			}

			pt := ParsedTransaction{
				Amount:      amount,
				Description: desc,
				Date:        txDate,
				Category:    InferCategory(desc),
				ExternalID:  fitid,
			}

			parsed = append(parsed, pt)
			reset()
			continue
		}

		if !inTxn {
			continue
		}

		if v := extractTagValue(lineLower, line, "dtposted"); v != "" {
			dateStr = v
		}
		if v := extractTagValue(lineLower, line, "trnamt"); v != "" {
			amtStr = v
		}
		if v := extractTagValue(lineLower, line, "name"); v != "" {
			name = v
		}
		if v := extractTagValue(lineLower, line, "memo"); v != "" {
			memo = v
		}
		if v := extractTagValue(lineLower, line, "fitid"); v != "" {
			fitid = v
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parsed, nil
}

func InsertParsedTransactions(parsed []ParsedTransaction) error {
	for _, p := range parsed {

		if p.ExternalID != "" {
			exists, err := db.TransactionExistsByExternalID(p.ExternalID)
			if err != nil {
				return err
			}
			if exists {
				continue
			}
		}

		tx := db.Transaction{
			Amount:      p.Amount,
			Description: p.Description,
			Category:    p.Category,
			Date:        p.Date,
			ExternalID:  p.ExternalID,
		}

		if err := db.InsertTransaction(tx); err != nil {
			return err
		}
	}
	return nil
}

func extractTagValue(lowerLine, originalLine, tag string) string {
	open := "<" + tag + ">"
	idx := strings.Index(lowerLine, open)
	if idx == -1 {
		return ""
	}
	return strings.TrimSpace(originalLine[idx+len(open):])
}

func parseFlexibleDate(s string) time.Time {
	s = strings.TrimSpace(s)
	formats := []string{"2006-01-02", "20060102", "02/01/2006", "1/2/2006"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Now()
}

func parseOFXDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Now()
	}
	if idx := strings.IndexAny(s, ".[-"); idx != -1 {
		s = s[:idx]
	}
	if len(s) >= 8 {
		s = s[:8]
	}
	if t, err := time.Parse("20060102", s); err == nil {
		return t
	}
	return time.Now()
}

type rule struct {
	re       *regexp.Regexp
	category string
}

var defaultRules = []rule{
	{regexp.MustCompile(`\b(supermarket|grocery|groceries|aldi|lidl|tesco|kaufland|spar)\b`), "Food"},
	{regexp.MustCompile(`\b(coffee|cafe|starbucks|espresso)\b`), "Coffee"},
	{regexp.MustCompile(`\b(salary|payroll|pay|donation)\b`), "Income"},
	{regexp.MustCompile(`\b(electricity|water bill|gas bill|utility|utilities)\b`), "Utilities"},
	{regexp.MustCompile(`\b(rent|landlord)\b`), "Rent"},
	{regexp.MustCompile(`\b(uber|taxi|lyft|cab|transport)\b`), "Transport"},
	{regexp.MustCompile(`\b(restaurant|dinner|lunch|breakfast|bar|date)\b`), "Dining"},
	{regexp.MustCompile(`\b(insurance)\b`), "Insurance"},
	{regexp.MustCompile(`\b(gym|fitness|yoga)\b`), "Fitness"},
	{regexp.MustCompile(`\b(movie|netflix|cinema)\b`), "Entertainment"},
}

func InferCategory(description string) string {
	s := strings.ToLower(description)
	for _, r := range defaultRules {
		if r.re.MatchString(s) {
			return r.category
		}
	}
	return "Uncategorized"
}

func ParseFileByPath(path string) ([]ParsedTransaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DetectAndParse(f, path)
}
