package internal

import (
	"crud/internal/repo"
)

func validateUpdatingRow(requestRow map[string]interface{}, tableName string, dbWrapper *repo.DBWrapper) (map[string]interface{}, error) {
	prKey := dbWrapper.DbInfo[tableName].PrKey

	_, contains := requestRow[prKey]
	if len(requestRow) == 1 && contains {
		return requestRow, ErrFieldInvalidType{Item: prKey}
	}

	newRow := make(map[string]interface{}, len(dbWrapper.DbInfo[tableName].ColInfo))
	for colName, val := range requestRow {
		colType, contains := dbWrapper.DbInfo[tableName].ColInfo[colName]
		if colName == prKey || !contains {
			continue
		}

		switch val.(type) {
		case float64:
			if colType.ScanType().Name() == "int32" {
				newRow[colName] = val.(int)
			} else {
				return requestRow, ErrFieldInvalidType{Item: colName}
			}
		case string:
			if colType.ScanType().Name() != "RawBytes" {
				return requestRow, ErrFieldInvalidType{Item: colName}
			}
			newRow[colName] = val
		case nil:
			isNullable, _ := colType.Nullable()
			if isNullable {
				newRow[colName] = val
			} else {
				return requestRow, ErrFieldInvalidType{Item: colName}
			}
		default:
			return requestRow, ErrFieldInvalidType{Item: colName}
		}

	}

	return newRow, nil
}

func validateNewRow(requestRow map[string]interface{}, tableName string, dbWrapper *repo.DBWrapper) (map[string]interface{}, error) {
	prKey := dbWrapper.DbInfo[tableName].PrKey
	newRow := make(map[string]interface{}, len(dbWrapper.DbInfo[tableName].ColInfo))
	for colName, colType := range dbWrapper.DbInfo[tableName].ColInfo {
		if colName == prKey {
			continue
		}

		val, ok := requestRow[colName]

		isNullable, _ := colType.Nullable()
		if !ok { // check nullable field
			if isNullable {
				newRow[colName] = nil
			} else {
				return requestRow, ErrFieldInvalidType{Item: colName}
			}
		} else {
			switch val.(type) {
			case float64:
				if colType.ScanType().Name() == "int32" { // TODO
					newRow[colName] = val.(int)
				} else {
					return requestRow, ErrFieldInvalidType{Item: colName}
				}
			case string:
				if colType.ScanType().Name() != "RawBytes" {
					return requestRow, ErrFieldInvalidType{Item: colName}
				}
				newRow[colName] = val
			case nil:
				if isNullable {
					newRow[colName] = val
				} else {
					return requestRow, ErrFieldInvalidType{Item: colName}
				}
			default:
				return requestRow, ErrFieldInvalidType{Item: colName}
			}

		}
	}

	return newRow, nil
}
