# structs

`structs` provides utility features related of the struct.

## Features

- propagation struct fields
- output diff struct fields

### propagation struct fields

struct fields in the struct can be copied from variable A to B.

#### Usage

```go
package structs_test

import (
	"fmt"

	"github.com/takashabe/structs"
)

func ExamplePropagateValues() {
	type User struct {
		ID   int
		Name string
		Age  int
	}

	a := &User{ID: 1, Name: "foo", Age: 1}
	b := &User{ID: 2, Name: "bar", Age: 2}
	got, _ := structs.PropagateValues(
		a,
		b,
		structs.PropagateWithIgnoreFields("ID", "Name"),
		structs.PropagateWithValue(1),
	)
	fmt.Printf("%+v", got)

	// Output:
	// &{ID:2 Name:bar Age:1}
}
```

### output diff struct fields

provide output of field names with different values

#### Usage

```go
func ExampleDiffFields() {
	type User struct {
		ID     int         `json:"id"`
		Name   string      `json:"name"`
		Age    int         `json:"age"`
		Gender null.String `json:"gender"`
	}

	a := &User{ID: 1, Name: "foo", Age: 1, Gender: null.StringFrom("Male")}
	b := &User{ID: 2, Name: "bar", Age: 1}
	got, _ := structs.DiffFields(
		a,
		b,
	)
	fmt.Printf("%+v", got)

	// Output:
	// [id name gender]
}
```
