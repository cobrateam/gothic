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

func TestGenerateSelectFromStruct(t *testing.T) {
	var p Person
	expected := "select name, age from person"
	got := Select(p)
	if expected != got {
		t.Errorf(`SELECT generation for %q. Was expecting "%s", got %s.`, reflect.TypeOf(p), expected, got)
	}
}
