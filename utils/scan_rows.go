package utils

import (
	"database/sql"
	"errors"
	"reflect"
)

var errNotStruct = errors.New("value passed to scanrow(s) is not of type struct")

func ScanRows[T any](rows *sql.Rows) ([]T, error) {
	var results []T

	typ := reflect.TypeOf(new(T)).Elem()
	if typ.Kind() != reflect.Struct {
		return nil, errNotStruct
	}
	for rows.Next() {
		instanceValue := reflect.New(typ).Elem()
		var scanArgs []any
		for i := range instanceValue.NumField() {
			fieldValue := instanceValue.Field(i)
			if instanceValue.Field(i).Type().Kind() == reflect.Struct {
				recursiveScan(&scanArgs, instanceValue)
			} else {
				scanArgs = append(scanArgs, fieldValue.Addr().Interface())
			}
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		results = append(results, instanceValue.Interface().(T))
	}
	return []T{}, nil
}

func recursiveScan(scanArgs *[]any, scannedStruct reflect.Value) {
	for i := range scannedStruct.NumField() {
		fieldValue := scannedStruct.Field(i)
		if scannedStruct.Field(i).Type().Kind() == reflect.Struct {
			recursiveScan(scanArgs, scannedStruct)
		} else {
			*scanArgs = append(*scanArgs, fieldValue.Addr().Interface())
		}
	}
}
