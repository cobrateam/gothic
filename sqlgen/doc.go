// Copyright 2012 Cobrateam members. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

	Package sqlgen is used for SQL generation. It provides functions that
	generates insert's, update's, delete's and select's from struct
	instances. It does not generate ready-to-run SQLs, but prepared
	statements.

	You can declare a struct as following:

		type Person struct {
			Name string
			Age  int
		}

	And then use the sqlgen to generate your SQL instructions. For
	instance, a select statement without any filters could be generated
	this way:

		var p Person
		sql, err := sqlgen.Select(p)

	This would generate:

		select name, age from person

*/
package sqlgen
