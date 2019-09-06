package awsssmenv

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type (
	// SSM defines the type implementing the SSM call.
	SSM interface {
		GetParameterWithContext(aws.Context, *ssm.GetParameterInput, ...request.Option) (*ssm.GetParameterOutput, error)
	}

	field struct {
		Index     int
		Name      string
		Encrypted bool
	}
)

const (
	tagEncrypted = "encrypted"
	tagName      = "ssm"
)

var (
	// ErrInvalidConfiguration indicates that a configuration is of the wrong type.
	ErrInvalidConfiguration = errors.New("configuration must be a struct pointer")

	// ErrInvalidFieldType indicates that a field is of the wrong type.
	ErrInvalidFieldType = errors.New("field must be a string")

	// ErrInvalidFieldAccess indicates that a field is of the wrong access.
	ErrInvalidFieldAccess = errors.New("field must be exported")
)

func loadFields(conf interface{}) ([]field, error) {
	rv := reflect.ValueOf(conf)
	if rv.Kind() != reflect.Ptr {
		return nil, ErrInvalidConfiguration
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, ErrInvalidConfiguration
	}

	var fields []field
	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		rf := rv.Field(i)

		var (
			f     field
			ok    bool
			value string
		)

		if value, ok = t.Field(i).Tag.Lookup(tagEncrypted); ok {
			f.Encrypted = true
		}

		if value, ok = t.Field(i).Tag.Lookup(tagName); ok {
			f.Index = i
			f.Name = value
		}

		if ok {
			if rf.Kind() != reflect.String {
				return nil, ErrInvalidFieldType
			}
			if !rf.CanSet() {
				return nil, ErrInvalidFieldAccess
			}

			fields = append(fields, f)
		}
	}

	return fields, nil
}

// Get uses environment variables to set field values in the conf struct, when
// necessary it requests remote values using AWS SSM.
func Get(ctx context.Context, conf interface{}, svc SSM) error {
	fields, err := loadFields(conf)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(conf).Elem()

	for _, f := range fields {
		rv.Field(f.Index).SetString(os.Getenv(f.Name))

		if evssm := os.Getenv(fmt.Sprintf("%s_SSM", f.Name)); evssm != "" {
			p := ssm.GetParameterInput{Name: &evssm, WithDecryption: &f.Encrypted}
			v, err := svc.GetParameterWithContext(ctx, &p)
			if err != nil {
				return err
			}
			rv.Field(f.Index).SetString(*v.Parameter.Value)
		}
	}

	return nil
}
