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
	fieldNames := []string{}
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}

	sql := fmt.Sprintf("select %s from %s", strings.Join(fieldNames, ", "), t.Name())
	return strings.ToLower(sql)
}
