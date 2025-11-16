package parser

import (
	"bufio"
	"errors"
	"io"
	"path/filepath"
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
		br := bufio.NewReader(r)
		peek, _ := br.Peek(512)
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

func parseCSV(r io.Reader) ([]ParsedTransaction, error) {
	// TODO: implement real CSV parsing. For now return empty but non-error slice.
	return []ParsedTransaction{}, nil
}

func parseOFX(r io.Reader) ([]ParsedTransaction, error) {
	// TODO: implement OFX parsing.
	return []ParsedTransaction{}, nil
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
