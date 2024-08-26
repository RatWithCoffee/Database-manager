package repo

import "database/sql"

type DBInfo map[string]TableInfo

type TableInfo struct {
	ColInfo map[string]*sql.ColumnType
	PrKey   string
}

type DBWrapper struct {
	Db     *sql.DB
	DbInfo DBInfo
}
