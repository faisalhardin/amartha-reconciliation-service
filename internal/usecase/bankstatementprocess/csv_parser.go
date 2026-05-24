package bankstatementprocess

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/shopspring/decimal"
)

var bankStatementCSVRequiredHeaders = []string{"unique_identifier", "amount", "date"}

type bankStatementCSVRow struct {
	UniqueIdentifier string
	Amount           string
	Date             string
}

func parseBankStatementCSV(file []byte) ([]bankStatementCSVRow, error) {
	reader := csv.NewReader(bytes.NewReader(file))
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("csv file is empty")
		}
		return nil, fmt.Errorf("failed to read csv header")
	}

	colIndex, err := mapBankStatementCSVHeaderIndexes(header)
	if err != nil {
		return nil, err
	}

	var rows []bankStatementCSVRow
	lineNumber := 1
	for {
		lineNumber++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("row %d: failed to parse csv", lineNumber)
		}

		if isEmptyBankStatementCSVRecord(record) {
			continue
		}

		if maxColumnIndex(colIndex) >= len(record) {
			return nil, fmt.Errorf("row %d: insufficient columns", lineNumber)
		}

		rows = append(rows, bankStatementCSVRow{
			UniqueIdentifier: record[colIndex["unique_identifier"]],
			Amount:           record[colIndex["amount"]],
			Date:             record[colIndex["date"]],
		})
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("csv file contains no statement rows")
	}

	return rows, nil
}

func buildBankStatementFromCSVRow(file *model.MstBankStatementFile, lineNumber int, row bankStatementCSVRow) (*model.MstBankStatement, error) {
	if strings.TrimSpace(row.UniqueIdentifier) == "" {
		return nil, fmt.Errorf("row %d: unique_identifier is required", lineNumber)
	}

	amount, err := decimal.NewFromString(strings.TrimSpace(row.Amount))
	if err != nil {
		return nil, fmt.Errorf("row %d: amount must be a valid decimal", lineNumber)
	}

	date, err := strconv.ParseInt(strings.TrimSpace(row.Date), 10, 64)
	if err != nil || date <= 0 {
		return nil, fmt.Errorf("row %d: date must be a positive unix timestamp", lineNumber)
	}

	if date < file.StartDate || date > file.EndDate {
		return nil, fmt.Errorf("row %d: date must be within file start_date and end_date", lineNumber)
	}

	return &model.MstBankStatement{
		BankCode:            file.BankCode,
		BankStatementFileID: file.ID,
		UniqueIdentifier:    strings.TrimSpace(row.UniqueIdentifier),
		Amount:              amount,
		Date:                date,
	}, nil
}

func mapBankStatementCSVHeaderIndexes(header []string) (map[string]int, error) {
	indexes := make(map[string]int, len(bankStatementCSVRequiredHeaders))
	for i, col := range header {
		indexes[strings.ToLower(strings.TrimSpace(col))] = i
	}

	for _, required := range bankStatementCSVRequiredHeaders {
		if _, ok := indexes[required]; !ok {
			return nil, fmt.Errorf("missing required column: %s", required)
		}
	}

	return indexes, nil
}

func isEmptyBankStatementCSVRecord(record []string) bool {
	for _, value := range record {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

func maxColumnIndex(indexes map[string]int) int {
	max := 0
	for _, idx := range indexes {
		if idx > max {
			max = idx
		}
	}
	return max
}
