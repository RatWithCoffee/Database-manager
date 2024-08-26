package internal

import (
	"crud/internal/repo"
	"net/url"
	"strconv"
	"strings"
)

func getIntParam(paramName string, defVal int, params url.Values) int {
	if len(params[paramName]) == 0 {
		return defVal
	}

	paramVal := params[paramName][0]
	res, err := strconv.Atoi(paramVal)
	if err != nil || res < 0 {
		return defVal
	}
	return res
}

func trimPath(url *url.URL) []string {
	path := strings.Split(url.Path, "/")
	if path[0] == "" {
		path = path[1:]
	}
	if path[len(path)-1] == "" {
		path = path[:len(path)-1]
	}
	return path
}

func isTableNameValid(tableName string, dbInfo repo.DBInfo) bool {
	_, ok := dbInfo[tableName]
	if !ok {
		return false
	}
	return true
}
