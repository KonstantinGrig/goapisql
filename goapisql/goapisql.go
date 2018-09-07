package goapisql

import (
	"database/sql"
	"encoding/json"
)

//GetQueryResult retrieves JSON result from DB or error
func GetQueryResult(db *sql.DB, sql string) (string, error) {
	var resString string
	tx, errBeginTx := db.Begin()
	defer tx.Commit()
	if errBeginTx != nil {
		return resString, errBeginTx
	}

	res := []map[string]interface{}{}

	rows, err := tx.Query(sql)
	if err != nil {
		return resString, err
	}

	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return "", err
		}
		var scanValues, rowMap = getEmptyRow(columns)
		rows.Scan(scanValues...)
		res = append(res, rowMap)
	}
	rows.Close()
	resByte, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	resString = string(resByte)
	return resString, nil
}

func getEmptyRow(sliceFieldNames []string) ([]interface{}, map[string]interface{}) {
	var sliceFieldValues []interface{}
	var resultMap = make(map[string]interface{})
	for _, fieldName := range sliceFieldNames {
		var fieldValue interface{}
		sliceFieldValues = append(sliceFieldValues, &fieldValue)
		resultMap[fieldName] = &fieldValue
	}
	return sliceFieldValues, resultMap
}
