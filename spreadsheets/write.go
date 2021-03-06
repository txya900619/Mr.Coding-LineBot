package spreadsheets

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

//AppendRow can send answers to spreadsheets
func (ss *Spreadsheets) AppendRow(row []string) error {
	interfaceRow := make([]interface{}, len(row))
	for i, v := range row {
		interfaceRow[i] = v
	}
	_, err := ss.Values.Append(ss.SpreadsheetsID, "表單回應 1", &sheets.ValueRange{Values: [][]interface{}{interfaceRow}}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("spreadsheet appendrow file, err: %v", err)
	}

	return nil
}
