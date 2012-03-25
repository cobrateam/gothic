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

	qm := make([]string, len(fieldNames)) // supply the question marks for the sql stmt
	for i := 0; i < len(qm); i++ {
		qm[i] = "?"
	}

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", t.Name(), strings.Join(fieldNames, ", "), strings.Join(qm, ", "))
	return strings.ToLower(sql)
}

func fieldNames(t reflect.Type) []string {
	fieldNames := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldNames[i] = t.Field(i).Name
	}

	return fieldNames
}
