package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// url?limit=50&page=3
// database will leave/will not show calculated entries from the beginning, (page -1) * limit(1-1*50=0)
// page 1, 1-1*limit=0 , first 50 entries 0-49
// page 2, 2-1*limit=50 = 50, next 50 entries 50-99
func GetPaginationParams(r *http.Request) (int, int) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	return limit, page
}

func GenerateInsertQuery(tableName string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println("dbTag:", dbTag)
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" { //skip ID field if its auto-incremented
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			placeholders += "?"
		}
	}
	//INSERT INTO XXX(first_name, last_name, email, class,subject) VALUES(?,?,?,?,?)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}

func GetStructValues(model interface{}) []interface{} {
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	values := []interface{}{}
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" { //skip ID field if its auto-incremented
			values = append(values, modelVal.Field(i).Interface())
		}
	}
	return values
}

// /students/?sortby=name:asc&sortby=class:desc
// SELECT id,first_name,last_name,email,class,subject FROM students WHERE 1=1 ORDER BY name ASC, class DESC
func AddSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY "
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], strings.ToUpper(parts[1])
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ", "
			}
			query += " " + field + " " + order
		}
	}
	return query
}

func AddFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
		"age":        "age",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + "=?"
			args = append(args, value)
		}
	}
	return query, args
}

func isValidSortOrder(order string) bool {
	return order == "ASC" || order == "DESC"
}

func isValidSortField(field string) bool {
	validField := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validField[field]
}
