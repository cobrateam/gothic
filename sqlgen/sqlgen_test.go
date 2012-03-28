// Copyright 2012 Cobrateam members. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlgen

import (
	"reflect"
	"strings"
	"testing"
)

type Person struct {
	Id   int
	Name string
	Age  int
}

func TestFieldNames(t *testing.T) {
	var p Person
	expected := []string{"id", "name", "age"}
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

func TestSelectAllFieldsFromStruct(t *testing.T) {
	var p Person
	expected := "select id, name, age from person"
	got, _ := Select(&p)
	if expected != got {
		t.Errorf(`SELECT generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestSelectOneFieldFromStruct(t *testing.T) {
	var p Person
	expected := "select name from person"
	got, _ := Select(&p, "name")

	if expected != got {
		t.Errorf(`SELECT generation for %q: Was expecting "%s", got "%s".`, reflect.TypeOf(p), expected, got)
	}
}

func TestSelectMultipleFieldsFromStruct(t *testing.T) {
	var p Person
	expected := "select age, name from person"
	got, _ := Select(&p, "age", "name")

	if expected != got {
		t.Errorf(`SELECT generation for %q: Was expecting "%s", got "%s".`, reflect.TypeOf(p), expected, got)
	}
}

func TestSelectAcceptStructValue(t *testing.T) {
	var p Person
	expected := "select id, name, age from person"
	got, _ := Select(p)
	if expected != got {
		t.Errorf(`SELECT generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestSelectReturnsErrorWhenObjectIsNotAStructNorAPointerToAStruct(t *testing.T) {
	i := 10
	_, err := Select(i)
	if err == nil || !strings.Contains(err.Error(), "provide a struct") {
		t.Errorf("SELECT generation: should not accept non-struct values/pointers")
	}
}

func TestSelectReturnsErrorWhenOneFieldIsNotInTheStruct(t *testing.T) {
	var p Person
	_, err := Select(p, "name", "weight")
	if err == nil || !strings.Contains(err.Error(), `Person does not have a field called "weight"`) {
		t.Errorf("SELECT generation: should return error when selecting fields not present in the struct %q", err)
	}
}

func TestInsertFromStructPointer(t *testing.T) {
	var p Person
	expected := "insert into person (id, name, age) values (?, ?, ?)"
	got, _ := Insert(&p)
	if expected != got {
		t.Errorf(`INSERT generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestInsertFromStructValue(t *testing.T) {
	var p Person
	expected := "insert into person (id, name, age) values (?, ?, ?)"
	got, _ := Insert(p)
	if expected != got {
		t.Errorf(`INSERT generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestInsertShouldReturnErrorWhenGivenObjectIsNotAStruct(t *testing.T) {
	i := 10
	_, err := Insert(i)
	if err == nil || !strings.Contains(err.Error(), "provide a struct") {
		t.Errorf("INSERT generation: should return error when the given type is not a struct")
	}
}

func TestSimpleDeleteFromStruct(t *testing.T) {
	p := Person{1, "Chuck", 32}
	expected := "delete from person where name=?"
	got := Delete(&p, []string{"Name"})
	if expected != got {
		t.Errorf(`DELETE generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestMultipleFilterDeleteFromStruct(t *testing.T) {
	p := Person{1, "Chuck", 32}
	expected := "delete from person where name=? and age=?"
	got := Delete(&p, []string{"Name", "Age"})
	if expected != got {
		t.Errorf(`DELETE generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestUpdateFromStructPointer(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	expected := "update person set name=?, age=? where id=?"
	got, _ := Update(&p, []string{"name", "age"}, []string{"id"})

	if expected != got {
		t.Errorf(`UPDATE generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestUpdateFromStructValue(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	expected := "update person set name=?, age=? where id=?"
	got, _ := Update(p, []string{"name", "age"}, []string{"id"})

	if expected != got {
		t.Errorf(`UPDATE generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestMultipleFilterUpdateFromStruct(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	expected := "update person set age=? where id=? and name=?"
	got, _ := Update(&p, []string{"age"}, []string{"id", "name"})

	if expected != got {
		t.Errorf(`UPDATE generation for %q: Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}

func TestUpdateReturnErrorIfTheGivenObjectIsNotAStruct(t *testing.T) {
	i := 10
	_, err := Update(i, []string{"name", "age"}, []string{"id"})

	if err == nil || !strings.Contains(err.Error(), "provide a struct") {
		t.Errorf("UPDATE generation: should return an error when obj is not a struct")
	}
}

func TestUpdateReturnErrorIfOneOfTheUpdateFieldsIsNotMemberOfObj(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	p.Age++
	_, err := Update(p, []string{"age", "weight"}, []string{"id"})

	if err == nil || !strings.Contains(err.Error(), `Person does not have a field called "weight"`) {
		t.Errorf("UPDATE generation: should return an error when one or more of the update fields is not member of the given struct")
	}
}

func TestUpdateReturnErrorIfOneOfTheFilterFieldsIsNotMemberOfObj(t *testing.T) {
	p := Person{Id: 1, Name: "Umi", Age: 6}
	p.Age++
	_, err := Update(p, []string{"age"}, []string{"weight"})

	if err == nil || !strings.Contains(err.Error(), `Person does not have a field called "weight"`) {
		t.Errorf("UPDATE generation: should return an error when one or more of the filter fields is not member of the given struct, %s")
	}
}

func TestCheckTypeReturnsTheTypeForStructPointer(t *testing.T) {
	var p = new(Person)
	tp, _ := checkType(p)

	if tp != reflect.TypeOf(*p) {
		t.Errorf("Check type should accept struct pointer")
	}
}

func TestCheckTypeReturnsTheTypeForStructValue(t *testing.T) {
	var p Person
	tp, _ := checkType(p)

	if tp != reflect.TypeOf(p) {
		t.Errorf("Check type should accept struct value")
	}
}

func TestCheckTypeReturnsErrorWhenTheTypeIsNotAnStruct(t *testing.T) {
	i := 10
	_, err := checkType(i)

	if err == nil || err.Error() != "Error generating SQL, you must provide a struct value or pointer" {
		t.Errorf("Check type should accept only structs and pointers to structs")
	}
}

func TestCheckFieldPresenceReturnsNilIfAllFieldsAreMemberOfTheStruct(t *testing.T) {
	var p Person
	err := checkPresenceOfFields(reflect.TypeOf(p), []string{"age", "name"})
	if err != nil {
		t.Errorf("Check presence of fields should return nil when all fields are present")
	}
}

func TestCheckFieldPresenceReturnsAnErrorIfAtLeastOneOfTheFieldsIsNotMemberOfTheStruct(t *testing.T) {
	var p Person
	err := checkPresenceOfFields(reflect.TypeOf(p), []string{"age", "school"})
	if err == nil || !strings.Contains(err.Error(), `Person does not have a field called "school"`) {
		t.Errorf("Check presence of field should return an error when one or more of the fields is not present")
	}
}
