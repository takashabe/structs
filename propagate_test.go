package structs_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/takashabe/structs"
	"gopkg.in/guregu/null.v4"
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
		structs.WithIgnoreFields("ID", "Name"),
		structs.WithValue(1),
	)
	fmt.Printf("%+v", got)

	// Output:
	// &{ID:2 Name:bar Age:1}
}

func Test_Propagate_same_struct(t *testing.T) {
	type A struct {
		ID          int
		Name        string
		Age         int
		Email       string
		IQ          int
		_unexported int
	}

	src := &A{
		ID:          1,
		Name:        "foo",
		Age:         10,
		Email:       "foo@example.com",
		IQ:          150,
		_unexported: 10,
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
		structs.WithIgnoreFields("ID", "Name"),
		structs.WithValue(10),
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
		structs.WithIgnoreFields("ID", "Name"),
		structs.WithValue(10),
	)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func Test_Propagate_withDefaultVaue(t *testing.T) {
	type B struct {
		Position string
		Grade    int
	}
	type A struct {
		ID          int
		Name        string
		Age         null.Int
		IsSmoker    null.Bool
		Address     string
		DateOfBirth time.Time
		Height      float64
		Weight      null.Float
		B           *B
	}

	src := &A{
		ID:          1,
		Name:        "foo",
		Age:         null.IntFrom(10),
		IsSmoker:    null.BoolFrom(false),
		Address:     "Chiyoda, Tokyo",
		DateOfBirth: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		Height:      180.01,
		Weight:      null.FloatFrom(60.01),
		B:           &B{Grade: 10},
	}
	dst := &A{
		ID:   2,
		Name: "bar",
	}

	want := &A{
		ID:          2,                                             // X: already set
		Name:        "bar",                                         // X: already set
		Age:         null.IntFrom(10),                              // O: copy
		IsSmoker:    null.BoolFrom(false),                          // O: copy
		Address:     "Chiyoda, Tokyo",                              // O: copy
		DateOfBirth: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), // O: copy
		Height:      180.01,                                        // O: copy
		Weight:      null.FloatFrom(60.01),                         // O: copy
		B:           &B{Grade: 10},                                 // O: copy
	}

	got, err := structs.PropagateValues(
		src,
		dst,
		structs.WithDefaultValue(),
	)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func Test_Propagate_different_pointerField(t *testing.T) {
	type B struct {
		Position string
		Grade    int
	}
	type C struct {
		Position string
		Grade    int64
		Salary   int64
	}
	type A struct {
		ID   int
		Name string
		B    *B
		B2   *B
	}
	type A2 struct {
		ID   int
		Name string
		B    *C // different type
		B2   *B
	}

	tests := []struct {
		name string
		opts []structs.Option
		src  *A
		dst  *A2
		want *A2
	}{
		{
			name: "none option",
			opts: []structs.Option{},
			src: &A{
				ID:   1,
				Name: "foo",
				B:    &B{Position: "Sales", Grade: 10},
				B2:   &B{Position: "Consultant", Grade: 9},
			},
			dst: &A2{
				ID:   2,
				Name: "bar",
				B:    &C{Position: "Programmer", Grade: 12, Salary: 30000000},
				B2:   &B{Position: "Manager", Grade: 8},
			},
			want: &A2{
				ID:   1,                                    // O: copy
				Name: "foo",                                // O: copy
				B:    &C{Position: "Sales", Grade: 10},     // O: copy
				B2:   &B{Position: "Consultant", Grade: 9}, // O: copy
			},
		},
		{
			name: "src nil",
			opts: []structs.Option{},
			src: &A{
				ID:   1,
				Name: "foo",
				B:    nil,
				B2:   nil,
			},
			dst: &A2{
				ID:   2,
				Name: "bar",
				B:    nil,
				B2:   &B{Position: "Manager", Grade: 8},
			},
			want: &A2{
				ID:   1,     // O: copy
				Name: "foo", // O: copy
				B:    nil,   // O: copy
				B2:   nil,   // O: copy
			},
		},
		{
			name: "WithDefaultValue",
			opts: []structs.Option{
				structs.WithDefaultValue(),
			},
			src: &A{
				ID:   1,
				Name: "foo",
				B:    &B{Position: "Sales", Grade: 10},
				B2:   &B{Position: "Consultant", Grade: 9},
			},
			dst: &A2{
				ID:   2,
				Name: "bar",
				B:    &C{Position: "Programmer", Grade: 12, Salary: 30000000},
			},
			want: &A2{
				ID:   2,                                                       // X: already set
				Name: "bar",                                                   // X: already set
				B:    &C{Position: "Programmer", Grade: 12, Salary: 30000000}, // X: already set
				B2:   &B{Position: "Consultant", Grade: 9},                    // O: copy
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := structs.PropagateValues(tt.src, tt.dst, tt.opts...)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Propagate_nested_struct(t *testing.T) {
	type Item struct {
		Price int
	}
	type Hobby struct {
		Name string
		Item *Item
	}
	type A struct {
		Age   int
		Hobby *Hobby
	}

	type Item2 struct {
		Price int
	}
	type Hobby2 struct {
		Name string
		Item *Item2
	}
	type A2 struct {
		Age   int
		Hobby *Hobby2
	}

	src := &A{
		Age: 10,
		Hobby: &Hobby{
			Name: "Go",
			Item: &Item{
				Price: 100,
			},
		},
	}
	dst := &A2{
		Age: 20,
	}
	want := &A2{
		Age: 10,
		Hobby: &Hobby2{
			Name: "Go",
			Item: &Item2{
				Price: 100,
			},
		},
	}

	got, err := structs.PropagateValues(src, dst)
	assert.NoError(t, err)
	diff := cmp.Diff(want, got)
	assert.Empty(t, diff)
}

func Test_Propagate_slice(t *testing.T) {
	type Hobby struct {
		Name string
	}
	type A struct {
		Age     int
		Hobbies []*Hobby
	}

	src := &A{
		Age: 10,
		Hobbies: []*Hobby{
			{Name: "Go"},
			{Name: "Rust"},
		},
	}
	dst := &A{
		Age: 20,
	}
	want := &A{
		Age: 10,
		Hobbies: []*Hobby{
			{Name: "Go"},
			{Name: "Rust"},
		},
	}

	got, err := structs.PropagateValues(src, dst)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
