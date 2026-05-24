package transaction

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/commonerr"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var csvRequiredHeaders = []string{"transaction_id", "amount", "type", "transaction_time"}

func (u *transactionUC) UploadFromCSV(ctx context.Context, bankCode string, file []byte) (*model.UploadMstTransactionCSVResponse, error) {
	normalizedBankCode, err := normalizeBankCode(bankCode)
	if err != nil {
		return nil, err
	}

	rows, err := parseTransactionCSV(file)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, commonerr.SetNewBadRequest("empty_csv", "csv file contains no transaction rows")
	}

	transactions := make([]model.MstTransaction, 0, len(rows))
	for i, row := range rows {
		tx, err := u.buildTransactionFromCSVRow(normalizedBankCode, i+2, row)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *tx)
	}

	if err := u.repo.InsertBatch(ctx, transactions); err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"UploadFromCSV.InsertBatch")
	}

	return &model.UploadMstTransactionCSVResponse{Inserted: len(transactions)}, nil
}

func (u *transactionUC) buildTransactionFromCSVRow(bankCode string, lineNumber int, row csvTransactionRow) (*model.MstTransaction, error) {
	if strings.TrimSpace(row.TransactionID) == "" {
		return nil, commonerr.SetNewBadRequest("invalid_csv", fmt.Sprintf("row %d: transaction_id is required", lineNumber))
	}

	amount, err := decimal.NewFromString(strings.TrimSpace(row.Amount))
	if err != nil || !amount.GreaterThan(decimal.Zero) {
		return nil, commonerr.SetNewBadRequest("invalid_csv", fmt.Sprintf("row %d: amount must be a number greater than zero", lineNumber))
	}

	txType := model.TransactionType(strings.ToUpper(strings.TrimSpace(row.Type)))
	if !txType.IsValid() {
		return nil, commonerr.SetNewBadRequest("invalid_csv", fmt.Sprintf("row %d: type must be DEBIT or CREDIT", lineNumber))
	}

	transactionTime, err := strconv.ParseInt(strings.TrimSpace(row.TransactionTime), 10, 64)
	if err != nil || transactionTime <= 0 {
		return nil, commonerr.SetNewBadRequest("invalid_csv", fmt.Sprintf("row %d: transaction_time must be a positive unix timestamp", lineNumber))
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"UploadFromCSV.NewV7")
	}

	return &model.MstTransaction{
		ID:                   id.String(),
		BankCode:             bankCode,
		TransactionReference: strings.TrimSpace(row.TransactionID),
		Amount:               amount,
		Type:                 txType,
		TransactionTime:      transactionTime,
	}, nil
}

type csvTransactionRow struct {
	TransactionID   string
	Amount          string
	Type            string
	TransactionTime string
}

func parseTransactionCSV(file []byte) ([]csvTransactionRow, error) {
	reader := csv.NewReader(bytes.NewReader(file))
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, commonerr.SetNewBadRequest("invalid_csv", "csv file is empty")
		}
		return nil, commonerr.SetNewBadRequest("invalid_csv", "failed to read csv header")
	}

	colIndex, err := mapCSVHeaderIndexes(header)
	if err != nil {
		return nil, err
	}

	var rows []csvTransactionRow
	lineNumber := 1
	for {
		lineNumber++
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, commonerr.SetNewBadRequest("invalid_csv", fmt.Sprintf("row %d: failed to parse csv", lineNumber))
		}

		if isEmptyCSVRecord(record) {
			continue
		}

		if maxColumnIndex(colIndex) >= len(record) {
			return nil, commonerr.SetNewBadRequest("invalid_csv", fmt.Sprintf("row %d: insufficient columns", lineNumber))
		}

		rows = append(rows, csvTransactionRow{
			TransactionID:   record[colIndex["transaction_id"]],
			Amount:          record[colIndex["amount"]],
			Type:            record[colIndex["type"]],
			TransactionTime: record[colIndex["transaction_time"]],
		})
	}

	return rows, nil
}

func mapCSVHeaderIndexes(header []string) (map[string]int, error) {
	indexes := make(map[string]int, len(csvRequiredHeaders))
	for i, col := range header {
		indexes[strings.ToLower(strings.TrimSpace(col))] = i
	}

	for _, required := range csvRequiredHeaders {
		if _, ok := indexes[required]; !ok {
			return nil, commonerr.SetNewBadRequest("invalid_csv", fmt.Sprintf("missing required column: %s", required))
		}
	}

	return indexes, nil
}

func isEmptyCSVRecord(record []string) bool {
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
