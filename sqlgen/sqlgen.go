// Copyright 2012 Cobrateam members. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlgen

import (
	"fmt"
	"reflect"
	"strings"
)

func Select(obj interface{}, args ...[]string) string {
	t := reflect.TypeOf(obj).Elem()
	var sql string

	if len(args) != 0 {
		sql = fmt.Sprintf("select %s from %s", strings.Join(args[0], ", "), t.Name())
	} else {
		fieldNames := fieldNames(t)
		sql = fmt.Sprintf("select %s from %s", strings.Join(fieldNames, ", "), t.Name())
	}

	return strings.ToLower(sql)
}

func Insert(obj interface{}) string {
	t := reflect.TypeOf(obj).Elem()
	fieldNames := fieldNames(t)

	qm := make([]string, len(fieldNames))
	for i := 0; i < len(qm); i++ {
		qm[i] = "?"
	}

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", t.Name(), strings.Join(fieldNames, ", "), strings.Join(qm, ", "))
	return strings.ToLower(sql)
}

func Delete(obj interface{}, filters []string) string {
	t := reflect.TypeOf(obj).Elem()

	filter_array := preparedFields(filters)
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
	t := reflect.TypeOf(obj).Elem()

	fieldsAndValues := preparedFields(uFields)
	filters := preparedFields(fFields)

	sql := fmt.Sprintf("update %s set %s where %s", strings.ToLower(t.Name()), strings.Join(fieldsAndValues, ", "), strings.Join(filters, " and "))

	return sql
}

// preparedFields receives a slice of fields an returns a slice with fields in the
// form of field=? that represents a placeholder to be replace for a value
func preparedFields(fields []string) []string {
	preparedFields := make([]string, len(fields))

	for i, v := range fields {
		preparedFields[i] = fmt.Sprintf(`%s=?`, strings.ToLower(v))
	}

	return preparedFields
}

func fieldNames(t reflect.Type) []string {
	fieldNames := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldNames[i] = t.Field(i).Name
	}

	return fieldNames
}
