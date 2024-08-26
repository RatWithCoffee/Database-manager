package repo

import (
	"database/sql"
	"fmt"
	"strings"
)

func (dbWrapper *DBWrapper) UpdateRow(tableName string, idReq int, reqRow map[string]interface{}) (int, error) {
	var queryVals strings.Builder
	newVals := make([]interface{}, 0)
	i := 0
	idLastStr := dbWrapper.DbInfo[tableName].PrKey + "=?"
	for name, val := range reqRow {
		if name == dbWrapper.DbInfo[tableName].PrKey {
			continue
		}
		if _, ok := dbWrapper.DbInfo[tableName].ColInfo[name]; ok {
			newVals = append(newVals, val)
			queryVals.WriteString(name)
			queryVals.WriteString("=?,")
			i++
		}
	}

	if len(newVals) == 0 {
		return 0, fmt.Errorf("updateRow [%s, %v] : %v", tableName, reqRow, "wrong body")
	}

	newVals = append(newVals, idReq)
	queryValsStr := queryVals.String()[:queryVals.Len()-1]
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, queryValsStr, idLastStr)
	insertRes, err := dbWrapper.Db.Exec(query, newVals...)
	if err != nil {
		return 0, fmt.Errorf("updateRow [%s, %v] : %v", tableName, reqRow, err)
	}

	numAffected, err := insertRes.RowsAffected()
	if numAffected == 0 {
		return 0, fmt.Errorf("updateRow [%s, %v] : %v", tableName, reqRow, "wrong id")

	}
	if err != nil {
		return 0, fmt.Errorf("updateRow [%s, %v] : %v", tableName, reqRow, err)
	}
	return int(numAffected), nil
}

func (dbWrapper *DBWrapper) AddRow(tableName string, reqRow map[string]interface{}) (int, error) {
	var p strings.Builder
	var n strings.Builder
	insertValues := make([]interface{}, len(dbWrapper.DbInfo[tableName].ColInfo)-1)
	i := 0
	for name, _ := range dbWrapper.DbInfo[tableName].ColInfo {
		if name == dbWrapper.DbInfo[tableName].PrKey {
			continue
		}
		insertValues[i] = reqRow[name]
		n.WriteString(name)
		n.WriteString(",")
		p.WriteString("?,")
		i++
	}

	placeholders := p.String()[:p.Len()-1]
	colNamesStr := n.String()[:n.Len()-1]

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, colNamesStr, placeholders)
	insertRes, err := dbWrapper.Db.Exec(query, insertValues...)
	if err != nil {
		return 0, fmt.Errorf("addRow [%s, %v] : %v", tableName, reqRow, err)
	}

	res, err := insertRes.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addRow [%s, %v] : %v", tableName, reqRow, err)
	}
	return int(res), nil
}

func (dbWrapper *DBWrapper) DeleteRow(tableName string, id int) (int, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, dbWrapper.DbInfo[tableName].PrKey)

	res, err := dbWrapper.Db.Exec(query, id)

	if err != nil {
		return 0, fmt.Errorf("deleteRow [%s]: %v", id, err)
	}

	num, _ := res.RowsAffected()
	return int(num), err
}

func (dbWrapper *DBWrapper) GetRowById(id int, tableName string) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, dbWrapper.DbInfo[tableName].PrKey)
	rows, err := dbWrapper.Db.Query(query, id)
	if err != nil {
		return res, fmt.Errorf("getRowById : %v", err)
	}
	defer rows.Close()

	colNames, _ := rows.Columns()
	colTypes, _ := rows.ColumnTypes()
	n := len(colNames)
	values := make([]sql.RawBytes, n)
	pointers := make([]interface{}, n)
	for i := 0; i < n; i++ {
		pointers[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(pointers...); err != nil {
			return res, fmt.Errorf("getRowById %d: %v", id, err)
		}

		res = getRow(values, colTypes, colNames)
	}

	return res, nil
}

// * GET /$table?limit=5&offset=7 - возвращает список из 5 записей (limit) начиная с 7-й (offset) из таблицы $table. limit
// по-умолчанию 5, offset 0
func (dbWrapper *DBWrapper) GetRows(tableName string, offset int, limit int) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)

	query := "SELECT * FROM " + tableName + " LIMIT ?, ?"
	rows, err := dbWrapper.Db.Query(query, offset, limit)
	if err != nil {
		return res, fmt.Errorf("getRows [%s, %d, %d]: %v", tableName, offset, limit, err)
	}
	defer rows.Close()

	colNames, _ := rows.Columns()
	colTypes, _ := rows.ColumnTypes()
	n := len(colNames)
	values := make([]sql.RawBytes, n)
	pointers := make([]interface{}, n)
	for i := 0; i < n; i++ {
		pointers[i] = &values[i]
	}
	for rows.Next() {
		if err := rows.Scan(pointers...); err != nil {
			return res, fmt.Errorf("getRows [%s, %d, %d]: %v", tableName, offset, limit, err)
		}

		rowMap := getRow(values, colTypes, colNames)
		res = append(res, rowMap)
	}

	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("getRows [%s, %d, %d]: %v", tableName, offset, limit, err)
	}

	return res, nil
}
