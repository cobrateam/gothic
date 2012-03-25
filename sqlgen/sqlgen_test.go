// Copyright 2012 Cobrateam members. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlgen

import (
	"reflect"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func TestFieldNames(t *testing.T) {
	var p Person
	expected := []string{"Name", "Age"}
	got := fieldNames(reflect.TypeOf(p))
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %q. Got %q.", expected, got)
	}
}

func TestGenerateSelectFromStruct(t *testing.T) {
	var p Person
	expected := "select name, age from person"
	got := Select(p)
	if expected != got {
		t.Errorf(`SELECT generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestGenerateInsertFromStruct(t *testing.T) {
	var p Person
	expected := "insert into person (name, age) values (?, ?)"
	got := Insert(p)
	if expected != got {
		t.Errorf(`INSERT generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestSimpleDeleteFromStruct(t *testing.T) {
	p := Person{"Chuck", 32}
	expected := "delete from person where name = Chuck"
	got := Delete(p, []string{"Name"})
	if expected != got {
		t.Errorf(`DELETE generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}
