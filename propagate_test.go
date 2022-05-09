package structs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_Propagate_same_struct(t *testing.T) {
	type A struct {
		ID    int
		Name  string
		Age   int
		Email string
		IQ    int
	}

	src := &A{
		ID:    1,
		Name:  "foo",
		Age:   10,
		Email: "foo@example.com",
		IQ:    150,
	}
	dst := &A{
		ID:    2,
		Name:  "bar",
		Age:   1,
		Email: "foo@example.com",
		IQ:    10,
	}

	want := &A{
		ID:    2,                 // X: ignore fields
		Name:  "bar",             // X: ignore fields
		Age:   10,                // O: copy
		Email: "foo@example.com", // X: unmatch value
		IQ:    10,                // X: unmatch value
	}

	got, err := structs.PropagateValues(
		src,
		dst,
		structs.PropagateWithIgnoreFields("ID", "Name"),
		structs.PropagateWithValue(10),
	)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func Test_Propagate_different_struct(t *testing.T) {
	type A struct {
		ID   int
		Name string
		Age  int
		Age2 int
	}

	type B struct {
		ID    int
		Name  string
		Age   int
		Email string
		IQ    int
	}

	src := &A{
		ID:   1,
		Name: "foo",
		Age:  10,
		Age2: 10,
	}
	dst := &B{
		ID:    2,
		Name:  "bar",
		Age:   1,
		Email: "foo@example.com",
		IQ:    10,
	}

	want := &B{
		ID:    2,                 // X: ignore fields
		Name:  "bar",             // X: ignore fields
		Age:   10,                // O: copy
		Email: "foo@example.com", // X: not exist A
		IQ:    10,                // X: not exist A
	}

	got, err := structs.PropagateValues(
		src,
		dst,
		structs.PropagateWithIgnoreFields("ID", "Name"),
		structs.PropagateWithValue(10),
	)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
