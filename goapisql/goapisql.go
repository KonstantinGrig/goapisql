package goapisql

import (
	"database/sql"
	"encoding/json"
	"github.com/KonstantinGrig/goapisql/config"
)

//GetQueryResult retrieves JSON result from DB or error
func GetQueryResult(role string, sql string) ([]byte, error) {
	db := config.GetDbConnection(role)
	tx, errBeginTx := db.Begin()
	if errBeginTx != nil {
		return nil, errBeginTx
	}

	res := []map[string]interface{}{}

	rows, err := tx.Query(sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		columnTypesName, err := getColumnTypesName(rows)
		if err != nil {
			return nil, err
		}

		var scanValues, rowMap = getEmptyRow(columns, columnTypesName)
		err = rows.Scan(scanValues...)
		if err != nil {
			return nil, err
		}
		res = append(res, rowMap)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	resByte, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return resByte, nil
}

func getColumnTypesName(rows *sql.Rows) ([]string, error) {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var columnTypesName []string
	for i := 0; i < len(columnTypes); i++ {
		columnTypesName = append(columnTypesName, columnTypes[i].DatabaseTypeName())
	}
	return columnTypesName, err
}

func getEmptyRow(sliceFieldNames []string, columnTypesName []string) ([]interface{}, map[string]interface{}) {
	var sliceFieldValues []interface{}
	var resultMap = make(map[string]interface{})
	for i, fieldName := range sliceFieldNames {
		switch columnTypesName[i] {
		case "":
			var fieldValue string
			sliceFieldValues = append(sliceFieldValues, &fieldValue)
			resultMap[fieldName] = &fieldValue
		default:
			var fieldValue interface{}
			sliceFieldValues = append(sliceFieldValues, &fieldValue)
			resultMap[fieldName] = &fieldValue
		}
	}
	return sliceFieldValues, resultMap
}
