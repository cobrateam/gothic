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

func TestPrepraredFields(t *testing.T) {
	expected := []string{"id=?", "name=?", "age=?"}
	got := preparedFields([]string{"id", "name", "age"})

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %q. Got %q.", expected, got)
	}
}

func TestGenerateSelectAllFieldsFromStruct(t *testing.T) {
	var p Person
	expected := "select id, name, age from person"
	got := Select(&p)
	if expected != got {
		t.Errorf(`SELECT generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestSelectOneFieldFromStruct(t *testing.T) {
	var p Person
	expected := "select name from person"
	got := Select(&p, []string{"name"})

	if expected != got {
		t.Errorf(`SELECT generation for %q. Was expecting "%s", got "%s".`, reflect.TypeOf(p), expected, got)
	}
}

func TestGenerateInsertFromStruct(t *testing.T) {
	var p Person
	expected := "insert into person (id, name, age) values (?, ?, ?)"
	got := Insert(&p)
	if expected != got {
		t.Errorf(`INSERT generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestSimpleDeleteFromStruct(t *testing.T) {
	p := Person{1, "Chuck", 32}
	expected := "delete from person where name=?"
	got := Delete(&p, []string{"Name"})
	if expected != got {
		t.Errorf(`DELETE generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestMultipleFilterDeleteFromStruct(t *testing.T) {
	p := Person{1, "Chuck", 32}
	expected := "delete from person where name=? and age=?"
	got := Delete(&p, []string{"Name", "Age"})
	if expected != got {
		t.Errorf(`DELETE generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestGenerateUpdateFromStruct(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	expected := "update person set name=?, age=? where id=?"
	got := Update(&p, []string{"name", "age"}, []string{"id"})

	if expected != got {
		t.Errorf(`UPDATE generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestMultipleFilterUpdateFromStructure(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	expected := "update person set age=? where id=? and name=?"
	got := Update(&p, []string{"age"}, []string{"id", "name"})

	if expected != got {
		t.Errorf(`UPDATE generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}
