// Copyright 2012 Cobrateam members. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlgen

import (
	"fmt"
	"reflect"
	"strings"
)

func Select(obj interface{}) string {
	t := reflect.TypeOf(obj)
	fieldNames := fieldNames(t)

	sql := fmt.Sprintf("select %s from %s", strings.Join(fieldNames, ", "), t.Name())
	return strings.ToLower(sql)
}

func Insert(obj interface{}) string {
	t := reflect.TypeOf(obj)
	fieldNames := fieldNames(t)

	qm := make([]string, len(fieldNames))
	for i := 0; i < len(qm); i++ {
		qm[i] = "?"
	}

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", t.Name(), strings.Join(fieldNames, ", "), strings.Join(qm, ", "))
	return strings.ToLower(sql)
}

func Delete(obj interface{}, filters []string) string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	filter_array := fieldValues(v, filters)
	filter_sql := strings.Join(filter_array, " and ")

	sql := fmt.Sprintf("delete from %s where ", t.Name())
	sql = strings.ToLower(sql)
	sql = sql + filter_sql
	return sql
}

// obj is the struct to be updated in the database
// uFields are the fields that are gonna be update
// fFields are the fields that are gonna be used as filter
// to the where clause
func Update(obj interface{}, uFields, fFields []string) string {
	t := reflect.TypeOf(obj)
	s := reflect.ValueOf(obj)

	fieldsAndValues := fieldValues(s, uFields)
	filters := fieldValues(s, fFields)

	sql := fmt.Sprintf("update %s %s where %s", strings.ToLower(t.Name()), strings.Join(fieldsAndValues, ", "), strings.Join(filters, ", "))

	return sql
}

// Receives a reflect.Value and the fields you want form the struct
// returns the respective values from the fields passed in the form of
// field=value, if value is a string, add " around it
func fieldValues(s reflect.Value, fields []string) []string {
	fieldValues := make([]string, len(fields))

	for i, v := range fields {
		f := s.FieldByName(strings.Title(v))

		var stmt string
		if f.Type().Kind() == reflect.String {
			stmt = fmt.Sprintf(`%s="%v"`, strings.ToLower(v), f.Interface())
		} else {
			stmt = fmt.Sprintf("%s=%v", strings.ToLower(v), f.Interface())
		}
		fieldValues[i] = stmt
	}

	return fieldValues
}

func fieldNames(t reflect.Type) []string {
	fieldNames := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldNames[i] = t.Field(i).Name
	}

	return fieldNames
}
