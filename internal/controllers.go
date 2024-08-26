package internal

import (
	"bytes"
	"crud/internal/repo"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// * PUT /$table/$id - обновляет запись, данные приходят в теле запроса (POST-параметры)
func handlePut(w http.ResponseWriter, r *http.Request, dbWrapper *repo.DBWrapper, tableName string, id int) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		log.Printf("handlePost [%s, %s]: %v", err, tableName)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var b interface{}
	if err := json.Unmarshal(buf.Bytes(), &b); err != nil {
		log.Printf("handlePost [%s, %s]: %v", err, tableName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqRow := b.(map[string]interface{})
	newRow, err := validateUpdatingRow(reqRow, tableName, dbWrapper)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFieldInvalidType := ErrFieldInvalidType{}
		if errors2.As(err, &errFieldInvalidType) {
			resp, _ := json.Marshal(RespError{Error: err.Error()})
			w.Write(resp)
			return
		}

		log.Println(err)
		return
	}

	idResp, err := dbWrapper.UpdateRow(tableName, id, newRow)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(RespOk{RespInner{Updated: &idResp}})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

// * POST /$table - создаёт новую запись, данный по записи в теле запроса (POST-параметры)
func handlePost(w http.ResponseWriter, r *http.Request, dbWrapper *repo.DBWrapper, tableName string) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		log.Printf("handlePost [%s, %s]: %v", err, tableName)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var body interface{}
	if err := json.Unmarshal(buf.Bytes(), &body); err != nil {
		log.Printf("handlePost [%s, %s]: %v", err, tableName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqRow := body.(map[string]interface{})
	newRow, err := validateNewRow(reqRow, tableName, dbWrapper)
	if err != nil {
		fmt.Println(newRow)
		w.WriteHeader(http.StatusBadRequest)

		errFieldInvalidType := ErrFieldInvalidType{}
		if errors2.As(err, &errFieldInvalidType) {
			fmt.Println(err)
			resp, _ := json.Marshal(RespError{Error: err.Error()})
			w.Write(resp)
			return
		}

		log.Println(err)
		return
	}

	id, err := dbWrapper.AddRow(tableName, newRow)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prKey := dbWrapper.DbInfo[tableName].PrKey
	resp, err := json.Marshal(RespOk{map[string]int{prKey: id}})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)

}

// * DELETE /$table/$id - удаляет запись
func handleDelete(w http.ResponseWriter, r *http.Request, dbWrapper *repo.DBWrapper, tableName string, id int) {
	num, err := dbWrapper.DeleteRow(tableName, id)
	if err != nil {
		log.Printf("handleDelete [%s, %s]: %v", err, tableName, id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(RespOk{RespInner{Deleted: &num}})
	if err != nil {
		log.Printf("handleDelete [%s, %s] : %v", tableName, id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)

}

// * GET /$table/$id - возвращает информацию о самой записи или 404
func handleGetById(w http.ResponseWriter, r *http.Request, dbWrapper *repo.DBWrapper, tableName string, id int) {
	inf, err := dbWrapper.GetRowById(id, tableName)
	if len(inf) == 0 {
		w.WriteHeader(http.StatusNotFound)
		resp, _ := json.Marshal(RespError{ErrRecordNotFound.Error()})
		w.Write(resp)
		return
	}
	resp, err := json.Marshal(RespOk{RespInner{Record: inf}})
	if err != nil {
		log.Printf("handleGetById %q : %v", inf, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

// * GET / - возвращает список все таблиц (которые мы можем использовать в дальнейших запросах)
func handleGetTablesName(w http.ResponseWriter, dbInfo repo.DBInfo) {
	names := make([]string, len(dbInfo))
	i := 0
	for k := range dbInfo {
		names[i] = k
		i++
	}

	resp, err := json.Marshal(RespOk{RespInner{Tables: names}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Write(resp)
}

// * GET /$table?limit=5&offset=7 - возвращает список из 5 записей (limit) начиная с 7-й (offset) из таблицы $table. limit
// по-умолчанию 5, offset 0
func handleGetList(w http.ResponseWriter, r *http.Request, dbWrapper *repo.DBWrapper, tableName string) {
	params := r.URL.Query()
	limit := getIntParam("limit", limitDefVal, params)
	offset := getIntParam("offset", offsetDefVal, params)

	list, err := dbWrapper.GetRows(tableName, offset, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	res, err := json.Marshal(RespOk{RespInner{Records: list}})
	if err != nil {
		fmt.Printf("handleGetList %q : %v", list, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
