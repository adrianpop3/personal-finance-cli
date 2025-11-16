package parser

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"personal-finance-cli/db"
)

// ParsedTransaction is a simple representation used by the parser
type ParsedTransaction struct {
	Amount      float64
	Description string
	Date        time.Time
	Category    string
}

// DetectAndParse reads the file contents and dispatches to CSV or OFX parser.
func DetectAndParse(r io.Reader, filename string) ([]ParsedTransaction, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".csv":
		return parseCSV(r)
	case ".ofx", ".qfx":
		return parseOFX(r)
	default:
		// try to sniff: read first few bytes
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

// parseCSV expects a header with: date,amount,description,category (category optional)
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

	// detect header indices (flexible)
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

	// If header row didn't contain recognized headings, assume fixed order:
	if idxDate == -1 && len(headers) >= 1 {
		idxDate = 0
	}
	if idxAmount == -1 && len(headers) >= 2 {
		idxAmount = 1
	}
	if idxDesc == -1 && len(headers) >= 3 {
		idxDesc = 2
	}
	// idxCat can remain -1 if not present

	var parsed []ParsedTransaction
	for i, row := range records[1:] {
		// skip empty lines
		if len(row) == 0 {
			continue
		}
		// guard indexes
		get := func(idx int) string {
			if idx >= 0 && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}
		dateStr := get(idxDate)
		amtStr := get(idxAmount)
		desc := get(idxDesc)
		cat := get(idxCat)

		// parse amount
		amtStr = strings.ReplaceAll(amtStr, ",", "") // remove thousands separators
		amount, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			// skip invalid rows
			continue
		}

		// parse date (try several formats)
		var txDate time.Time
		parseDate := func(s string) (time.Time, error) {
			s = strings.TrimSpace(s)
			if s == "" {
				return time.Time{}, errors.New("empty date")
			}
			formats := []string{"2006-01-02", "20060102", "02/01/2006", "1/2/2006"}
			for _, f := range formats {
				if t, e := time.Parse(f, s); e == nil {
					return t, nil
				}
			}
			return time.Time{}, fmt.Errorf("unrecognized date: %s", s)
		}
		if d, err := parseDate(dateStr); err == nil {
			txDate = d
		} else {
			// fallback to today
			txDate = time.Now()
		}

		pt := ParsedTransaction{
			Amount:      amount,
			Description: desc,
			Date:        txDate,
			Category:    cat,
		}
		// auto-categorize if missing
		if strings.TrimSpace(pt.Category) == "" {
			pt.Category = inferCategory(pt.Description)
		}
		parsed = append(parsed, pt)

		// safety: avoid extremely large parse loops (defensive)
		_ = i
	}

	return parsed, nil
}

// Very small OFX parser that extracts STMTTRN entries (DTPOSTED, TRNAMT, NAME or MEMO)
func parseOFX(r io.Reader) ([]ParsedTransaction, error) {
	scanner := bufio.NewScanner(r)
	var parsed []ParsedTransaction
	var inTxn bool
	var dateStr, amtStr, name, memo string

	reset := func() {
		inTxn = false
		dateStr, amtStr, name, memo = "", "", "", ""
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineLower := strings.ToLower(line)
		if strings.HasPrefix(lineLower, "<stmttrn") {
			inTxn = true
			continue
		}
		if strings.HasPrefix(lineLower, "</stmttrn") {
			// commit transaction
			// parse date, amount
			var txDate time.Time
			if dateStr != "" {
				// OFX DTPOSTED often like YYYYMMDD or YYYYMMDDHHMMSS
				ds := dateStr
				if len(ds) >= 8 {
					ds = ds[:8]
				}
				if d, err := time.Parse("20060102", ds); err == nil {
					txDate = d
				} else {
					txDate = time.Now()
				}
			} else {
				txDate = time.Now()
			}
			amtStr = strings.ReplaceAll(amtStr, ",", "")
			amount, err := strconv.ParseFloat(amtStr, 64)
			if err != nil {
				// skip invalid
				reset()
				continue
			}
			desc := name
			if desc == "" {
				desc = memo
			}
			pt := ParsedTransaction{
				Amount:      amount,
				Description: desc,
				Date:        txDate,
				Category:    "",
			}
			if strings.TrimSpace(pt.Category) == "" {
				pt.Category = inferCategory(pt.Description)
			}
			parsed = append(parsed, pt)
			reset()
			continue
		}
		if !inTxn {
			continue
		}
		// extract tags simply
		if strings.HasPrefix(lineLower, "<dtposted>") {
			dateStr = strings.TrimSpace(line[len("<dtposted>"):])
		} else if strings.HasPrefix(lineLower, "<trnamt>") {
			amtStr = strings.TrimSpace(line[len("<trnamt>"):])
		} else if strings.HasPrefix(lineLower, "<name>") {
			name = strings.TrimSpace(line[len("<name>"):])
		} else if strings.HasPrefix(lineLower, "<memo>") {
			memo = strings.TrimSpace(line[len("<memo>"):])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return parsed, nil
}

// InsertParsedTransactions converts to db.Transaction and inserts into DB
func InsertParsedTransactions(parsed []ParsedTransaction) error {
	for _, p := range parsed {
		tx := db.Transaction{
			Amount:      p.Amount,
			Description: p.Description,
			Category:    p.Category,
			Date:        p.Date,
		}
		if err := db.InsertTransaction(tx); err != nil {
			return err
		}
	}
	return nil
}

// ------------------ Auto-categorization ------------------

// rule holds a compiled regexp and the category to use when it matches.
type rule struct {
	re       *regexp.Regexp
	category string
}

// defaultRules is a small set of heuristics for inferring categories.
var defaultRules = []rule{
	{regexp.MustCompile(`\b(supermarket|grocery|groceries|aldi|lidl|tesco|spar)\b`), "Food"},
	{regexp.MustCompile(`\b(coffee|cafe|starbucks|espresso)\b`), "Coffee"},
	{regexp.MustCompile(`\b(salary|payroll|pay)\b`), "Income"},
	{regexp.MustCompile(`\b(electricity|water bill|gas bill|utility|utilities)\b`), "Utilities"},
	{regexp.MustCompile(`\b(rent|landlord)\b`), "Rent"},
	{regexp.MustCompile(`\b(uber|taxi|lyft|cab|transport)\b`), "Transport"},
	{regexp.MustCompile(`\b(restaurant|dinner|lunch|breakfast|bar)\b`), "Dining"},
	{regexp.MustCompile(`\b(insurance)\b`), "Insurance"},
}

// inferCategory checks description against rules and returns the first match or "Uncategorized".
func inferCategory(description string) string {
	s := strings.ToLower(description)
	for _, r := range defaultRules {
		if r.re.MatchString(s) {
			return r.category
		}
	}
	return "Uncategorized"
}

// Helper to allow parsing a file path (used by TUI)
func ParseFileByPath(path string) ([]ParsedTransaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DetectAndParse(f, path)
}
