package spreadsheets

import (
	"google.golang.org/api/sheets/v4"
	"strconv"
	"time"
)

func (ss *Spreadsheets) SaveValueToSpecificCell(value, ranges string) error {
	_, err := ss.Values.Update(ss.SpreadsheetsID, ranges, &sheets.ValueRange{Values: [][]interface{}{{value}}}).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}
	return nil
}

// Complete question
func (ss *Spreadsheets) DeleteUserID(rowID int) error {
	rowIdStr := strconv.Itoa(rowID)
	ranges := "I" + rowIdStr + ":" + "I" + rowIdStr
	err := ss.SaveValueToSpecificCell("Complete", ranges)
	return err
}

func (ss *Spreadsheets) InsertTimestampAndUserID(userID string, rowID int) error {
	rowIdStr := strconv.Itoa(rowID)
	timeStampRange := "A" + rowIdStr + ":" + "A" + rowIdStr
	userIdRange := "I" + rowIdStr + ":" + "I" + rowIdStr

	err := ss.SaveValueToSpecificCell(time.Now().String(), timeStampRange)
	if err != nil {
		return err
	}

	err = ss.SaveValueToSpecificCell(userID, userIdRange)
	if err != nil {
		return err
	}

	return nil
}
