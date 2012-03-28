// Copyright 2012 Cobrateam members. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlgen

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Select generates a SELECT statement, selecting only "fields" from
// the table given by the name of the struct "obj" lowercased. If the type
// of obj is not a struct, this method returns an empty string and an error.
//
// If the type of obj is a struct, but one of the given fields is not a member
// of the struct (lowercased), it returns an empty string and another error.
//
// Otherwise it returns the SQL instruction and a nil error.
func Select(obj interface{}, fields ...string) (string, error) {
	t, err := checkType(obj)
	if err != nil {
		return "", err
	}
	if len(fields) == 0 {
		fields = fieldNames(t)
	} else {
		err := checkPresenceOfFields(t, fields)
		if err != nil {
			return "", err
		}
	}

	sql := fmt.Sprintf("select %s from %s", strings.Join(fields, ", "), strings.ToLower(t.Name()))
	return sql, nil
}

// Insert generates an INSERT statement, using all values for the members of
// obj. If obj is not a struct nor a pointer to a struct, this method returns
// an empty string and an error. Otherwise, it returns the SQL instruction and
// a nil error.
func Insert(obj interface{}) (string, error) {
	t, err := checkType(obj)
	if err != nil {
		return "", err
	}

	fields := fieldNames(t)
	qm := make([]string, len(fields))
	for i := 0; i < len(qm); i++ {
		qm[i] = "?"
	}
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", strings.ToLower(t.Name()), strings.Join(fields, ", "), strings.Join(qm, ", "))
	return sql, nil
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

// Update generates an UPDATE statement, using the given updateFields and
// filterFields. updateFields contains a list of fields to include in the update
// and filterFields contains a list of field to use in WHERE clause.
//
// The following code:
//
//		p := Person{Id: 10, Name: "Batman", Age: 90}
//		sql, _ := Update(p, []string{"name"}, []string{"id"})
//
// will generate the following SQL:
//
//		update person set name = ? where id = ?
//
// Note that it generates just a string to build a prepared statement.
//
// If the given object is not a struct, or one of the fields in the update
// or filter list is not member of the given struct, this function returns
// an empty string and an error. Otherwise, it returns the SQL for build a
// prepared statement and nil error.
func Update(obj interface{}, updateFields, filterFields []string) (string, error) {
	t, err := checkType(obj)
	if err != nil {
		return "", err
	}
	err = checkPresenceOfFields(t, updateFields)
	if err != nil {
		return "", err
	}
	err = checkPresenceOfFields(t, filterFields)
	if err != nil {
		return "", err
	}

	fieldsAndValues := preparedFields(updateFields)
	filters := preparedFields(filterFields)
	sql := fmt.Sprintf("update %s set %s where %s", strings.ToLower(t.Name()), strings.Join(fieldsAndValues, ", "), strings.Join(filters, " and "))
	return sql, nil
}

func checkPresenceOfFields(t reflect.Type, fields []string) error {
	names := fieldNames(t)
	for _, givenName := range fields {
		var found bool
		for _, structName := range names {
			if givenName == structName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf(`%s does not have a field called "%s"`, t.Name(), givenName)
		}
	}
	return nil
}

func checkType(obj interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		return t, nil
	}
	return nil, errors.New("Error generating SQL, you must provide a struct value or pointer")
}

func fieldNames(t reflect.Type) []string {
	fieldNames := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldNames[i] = strings.ToLower(t.Field(i).Name)
	}

	return fieldNames
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
