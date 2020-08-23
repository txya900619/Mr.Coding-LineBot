package mrcoding

import (
	"Mr.Coding-LineBot/spreadsheets"
	"strconv"
)

func getRange(rowID int, colID spreadsheets.ColumnID) string {
	colIdStr := colID.String()
	rowIDStr := strconv.Itoa(rowID)
	return colIdStr + rowIDStr + ":" + colIdStr + rowIDStr

}
