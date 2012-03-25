// Copyright 2012 Cobrateam members. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlgen

import (
	"reflect"
	"testing"
)

type Person struct {
	Id   int
	Name string
	Age  int
}

func TestFieldNames(t *testing.T) {
	var p Person
	expected := []string{"Id", "Name", "Age"}
	got := fieldNames(reflect.TypeOf(p))
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %q. Got %q.", expected, got)
	}
}

func TestFieldValues(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	expected := []string{"id=1", `name="Umi"`, "age=6"}
	got := fieldValues(reflect.ValueOf(p), []string{"id", "name", "age"})

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %q. Got %q.", expected, got)
	}
}

func TestGenerateSelectFromStruct(t *testing.T) {
	var p Person
	expected := "select id, name, age from person"
	got := Select(p)
	if expected != got {
		t.Errorf(`SELECT generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestGenerateInsertFromStruct(t *testing.T) {
	var p Person
	expected := "insert into person (id, name, age) values (?, ?, ?)"
	got := Insert(p)
	if expected != got {
		t.Errorf(`INSERT generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestGenerateUpdateFromStruct(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	expected := `update person name="Umi", age=6 where id=1`
	got := Update(p, []string{"name", "age"}, []string{"id"})

	if expected != got {
		t.Errorf(`UPDATE generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}
