package internal

import (
	"crud/internal/repo"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	dbInfo, err := repo.GetDBInfo(db)
	if err != nil {
		return nil, err
	}

	fmt.Println(dbInfo)

	return router(repo.DBWrapper{Db: db, DbInfo: dbInfo}), nil
}

func router(dbWrapper repo.DBWrapper) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := trimPath(r.URL)

		switch len(path) {
		case 0:
			routerUrlRoot(w, r, dbWrapper)
		case 1:
			tableName := path[0]
			ok := isTableNameValid(tableName, dbWrapper.DbInfo)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				errMsg, _ := json.Marshal(RespError{ErrUnknownTable.Error()})
				w.Write(errMsg)
				return
			}
			routerUrlTable(w, r, dbWrapper, tableName)
		case 2:
			tableName := path[0]
			idStr := path[1]

			id, err := strconv.Atoi(idStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			ok := isTableNameValid(tableName, dbWrapper.DbInfo)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			routerUrlTableId(w, r, dbWrapper, tableName, id)
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

	}
	return http.HandlerFunc(fn)
}

func routerUrlRoot(w http.ResponseWriter, r *http.Request, dbWrapper repo.DBWrapper) {
	switch r.Method {
	case http.MethodGet:
		handleGetTablesName(w, dbWrapper.DbInfo)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func routerUrlTable(w http.ResponseWriter, r *http.Request, dbWrapper repo.DBWrapper, tableName string) {
	switch r.Method {
	case http.MethodGet:
		handleGetList(w, r, &dbWrapper, tableName)
	case http.MethodPost:
		handlePost(w, r, &dbWrapper, tableName)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func routerUrlTableId(w http.ResponseWriter, r *http.Request, dbWrapper repo.DBWrapper, tableName string, id int) {
	switch r.Method {
	case http.MethodGet:
		handleGetById(w, r, &dbWrapper, tableName, id)
	case http.MethodDelete:
		handleDelete(w, r, &dbWrapper, tableName, id)
	case http.MethodPut:
		handlePut(w, r, &dbWrapper, tableName, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
