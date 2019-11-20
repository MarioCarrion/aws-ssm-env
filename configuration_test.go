package awsssmenv_test

import (
	"context"
	"errors"
	"os"
	"testing"

	awsssm "github.com/aws/aws-sdk-go/service/ssm"

	awsssmenv "github.com/MarioCarrion/aws-ssm-env"
	"github.com/MarioCarrion/aws-ssm-env/awsssmenvtesting"
	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	type (
		conf struct {
			WithSSM         string `ssm:"WITH"`
			WithoutSSM      string
			WitSSMEncrypted string `ssm:"WITHENC,encrypted"`
		}

		expected struct {
			value interface{}
			err   bool
		}
	)

	tests := []struct {
		name     string
		input    interface{}
		expected expected
		setup    func(*awsssmenvtesting.FakeSSM)
		teardown func()
	}{
		{
			"Error: loading struct fields",
			new(int),
			expected{err: true},
			func(*awsssmenvtesting.FakeSSM) {},
			func() {},
		},
		{
			"OK: env vars exist but no postfixed one",
			&conf{
				WithSSM: "otherValue",
			},
			expected{
				value: &conf{
					WithSSM: "yesSSM",
				},
			},
			func(*awsssmenvtesting.FakeSSM) { os.Setenv("WITH", "yesSSM") },
			func() { os.Unsetenv("WITH") },
		},
		{
			"OK: non decorated values not modified",
			&conf{
				WithoutSSM: "notModified",
			},
			expected{
				value: &conf{
					WithoutSSM: "notModified",
				},
			},
			func(*awsssmenvtesting.FakeSSM) {},
			func() {},
		},
		{
			"OK: env vars exist and also postfixed one",
			&conf{
				WithSSM: "otherValue",
			},
			expected{
				value: &conf{
					WithSSM: "something remote",
				},
			},
			func(f *awsssmenvtesting.FakeSSM) {
				os.Setenv("WITH", "yesSSM")
				os.Setenv("WITH_SSM", "/remote/value")
				f.GetParameterWithContextReturns(&awsssm.GetParameterOutput{
					Parameter: &awsssm.Parameter{
						Value: func() *string {
							value := "something remote"
							return &value
						}(),
					},
				}, nil)
			},
			func() {
				os.Unsetenv("WITH")
				os.Unsetenv("WITH_SSM")
			},
		},
		{
			"Error: env vars exist and also postfixed one AND but remote fails",
			&conf{
				WithSSM: "otherValueFailed",
			},
			expected{err: true},
			func(f *awsssmenvtesting.FakeSSM) {
				os.Setenv("WITH", "yesSSM_failed")
				os.Setenv("WITH_SSM", "/remote/value/failed")
				f.GetParameterWithContextReturns(&awsssm.GetParameterOutput{}, errors.New("remote failed"))
			},
			func() {
				os.Unsetenv("WITH")
				os.Unsetenv("WITH_SSM")
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			defer test.teardown()

			ssm := awsssmenvtesting.FakeSSM{}
			test.setup(&ssm)

			goterr := awsssmenv.Get(context.Background(), test.input, &ssm)
			if (goterr != nil) != test.expected.err {
				t.Fatalf("expected err, got nothing")
			}

			if test.expected.err {
				return
			}

			if !cmp.Equal(test.expected.value, test.input) {
				t.Errorf("expected values do not match\n%s", cmp.Diff(test.expected.value, test.input))
			}
		})
	}
}
