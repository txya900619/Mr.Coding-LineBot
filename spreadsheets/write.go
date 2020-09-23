package spreadsheets

import (
	"google.golang.org/api/sheets/v4"
)

func (ss *Spreadsheets) AppendRow(row []interface{}) error {

	_, err := ss.Values.Append(ss.SpreadsheetsID, "表單回應 1", &sheets.ValueRange{Values: [][]interface{}{row}}).ValueInputOption("INPUT_VALUE_OPTION_UNSPECIFIED").Do()
	if err != nil {
		return err
	}

	return nil
}
