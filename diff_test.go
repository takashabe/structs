package structs_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/takashabe/structs"
	"gopkg.in/guregu/null.v4"
)

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
		structs.WithIgnoreFields("ID", "Gender"),
	)
	fmt.Printf("%+v", got)

	// Output:
	// [name]
}

func Test_DiffFields_withOptions(t *testing.T) {
	type B struct {
		Position string
		Grade    int
	}
	type A struct {
		ID          int        `json:"id"`
		Name        string     `json:"name"`
		Age         null.Int   `json:"age"`
		IsSmoker    null.Bool  `json:"is_smoker"`
		Address     string     `json:"address"`
		DateOfBirth time.Time  `json:"date_of_birth"`
		Height      float64    `json:"height"`
		Weight      null.Float `json:"weight"`
		B           *B         `json:"b"`
	}

	tests := []struct {
		name string
		src  *A
		dst  *A
		opts []structs.Option
		want []string
	}{
		{
			name: "WithDefaultValue",
			src: &A{
				ID:          1,
				Name:        "foo",
				Age:         null.IntFrom(10),
				IsSmoker:    null.BoolFrom(false),
				Address:     "Chiyoda, Tokyo",
				DateOfBirth: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
				Height:      180.01,
				Weight:      null.FloatFrom(60.01),
				B:           &B{Grade: 10},
			},
			dst: &A{
				ID:   2,
				Name: "bar",
			},
			opts: []structs.Option{structs.WithDefaultValue()},
			want: []string{"age", "is_smoker", "address", "date_of_birth", "height", "weight", "b"},
		},
		{
			name: "WithDefaultValue, WithIgnoreSourceDefaultValue",
			src: &A{
				ID:   1,
				Name: "foo",
			},
			dst: &A{
				ID: 2,
			},
			opts: []structs.Option{structs.WithDefaultValue(), structs.WithIgnoreSourceDefaultValue()},
			want: []string{"name"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := structs.DiffFields(
				tt.src,
				tt.dst,
				tt.opts...,
			)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
