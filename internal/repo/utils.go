package repo

import (
	"database/sql"
	"strconv"
)

func getRow(values []sql.RawBytes, colTypes []*sql.ColumnType, colNames []string) map[string]interface{} {
	res := make(map[string]interface{})
	var vStr string
	for i, v := range values {
		vStr = string(v)
		switch colTypes[i].ScanType().Name() {
		case "int32":
			res[colNames[i]], _ = strconv.Atoi(vStr)
		default:
			if vStr == "" {
				if isNullable, _ := colTypes[i].Nullable(); isNullable {
					res[colNames[i]] = nil
				} else {
					res[colNames[i]] = ""
				}
			} else {
				res[colNames[i]] = vStr
			}
		}

	}
	return res
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
