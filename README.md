# structs

`structs` provides utility features related of the struct.

## Features

- propagation struct fields

### propagation struct fields

In the same struct, struct fields can be copied from struct variable A to B.

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

