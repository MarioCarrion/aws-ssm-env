package awsssmenv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadFields(t *testing.T) {
	type expected struct {
		fields []field
		err    error
	}

	tests := []struct {
		name     string
		input    interface{}
		expected expected
	}{
		{
			"OK",
			&struct {
				Encrypted    string `ssm:"ENCRYPTED" encrypted:"-"`
				NotEncrypted string `ssm:"NOT_ENCRYPTED"`
			}{},
			expected{
				fields: []field{
					{
						Index:     0,
						Name:      "ENCRYPTED",
						Encrypted: true,
					},
					{
						Index: 1,
						Name:  "NOT_ENCRYPTED",
					},
				},
			},
		},
		{
			"OK: ignoring untagged",
			&struct {
				One bool
				Two string `ssm:"ONE"`
			}{},
			expected{
				fields: []field{
					{
						Index: 1,
						Name:  "ONE",
					},
				},
			},
		},
		{
			"Error: not a pointer (ErrInvalidConfiguration)",
			struct{}{},
			expected{err: ErrInvalidConfiguration},
		},
		{
			"Error: not a struct (ErrInvalidConfiguration)",
			new(int64),
			expected{err: ErrInvalidConfiguration},
		},
		{
			"Error: tagged is not string (ErrInvalidFieldType)",
			&struct {
				NotString bool `ssm:"ERROR"`
			}{},
			expected{err: ErrInvalidFieldType},
		},
		{
			"Error: tagged is not exported (ErrInvalidFieldAccess)",
			&struct {
				unexported string `ssm:"ERROR"`
			}{},
			expected{err: ErrInvalidFieldAccess},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			gotfields, goterr := loadFields(test.input)

			if !cmp.Equal(test.expected.fields, gotfields) {
				t.Errorf("expected values do not match\n%s", cmp.Diff(test.expected.fields, gotfields))
			}

			if test.expected.err != goterr {
				t.Errorf("expected %s, got %s", test.expected.err, goterr)
			}
		})
	}
}
