package repo

import (
	"database/sql"
	"fmt"
)

func GetTableNames(db *sql.DB) ([]string, error) {
	res := make([]string, 1)

	rows, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA='crud'")
	if err != nil {
		return res, fmt.Errorf("getTableNames : %v", err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			if err != nil {
				return res, fmt.Errorf("getTableNames : %v", err)
			}
		}
		tables = append(tables, name)
	}
	return tables, nil
}

func GetDBInfo(db *sql.DB) (DBInfo, error) {
	tableNames, err := GetTableNames(db)
	dbInfo := make(map[string]TableInfo, len(tableNames))

	if err != nil {
		return nil, fmt.Errorf("getDBInfo : %v", err)
	}

	for _, name := range tableNames {
		tableInfo, err := GetTableInfo(db, name, tableNames)
		if err != nil {
			return dbInfo, nil
		}
		dbInfo[name] = tableInfo
	}

	return dbInfo, nil
}

func GetTableInfo(db *sql.DB, tableName string, tableNames []string) (TableInfo, error) {
	var tableInfo TableInfo
	ok := contains(tableNames, tableName)
	if !ok {
		return tableInfo, fmt.Errorf("getTableInfo [%s, %v]: %s", tableName, tableNames, "invalid table name")
	}

	query := "SELECT * FROM " + tableName
	rows, err := db.Query(query)
	if err != nil {
		return tableInfo, fmt.Errorf("getTableInfo [%s, %v]: %s", tableName, tableNames, err)
	}
	defer rows.Close()

	colNames, _ := rows.Columns()
	colTypes, _ := rows.ColumnTypes()

	colNamesMap := make(map[string]struct{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		colNamesMap[colNames[i]] = struct{}{}
	}

	colInfo := make(map[string]*sql.ColumnType, len(colNames))
	for i := 0; i < len(colNames); i++ {
		colInfo[colNames[i]] = colTypes[i]
	}
	prKey := colNames[0]
	//return TableInfo{colNames, colTypes, colNamesMap}, nil
	return TableInfo{colInfo, prKey}, nil
}
