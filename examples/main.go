// Run it as:
// AWS_REGION=us-east1 go run examples/main.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"

	awsssmenv "github.com/MarioCarrion/aws-ssm-env"
)

func main() {
	// XXX Error validation omitted
	session, _ := awssession.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	ssm := awsssm.New(session)

	v := struct {
		Username string `ssm:"USER"`
	}{}

	// If USER_SSM is defined on the system we will contact AWS SSM to get
	// the remote value defined in this variable.
	// For example if USER_SSM is "/remote/user" then AWS SSM will be queried
	// using "/remote/user" and the result will be stored in "v.Username"

	if err := awsssmenv.Get(context.Background(), &v, ssm); err != nil {
		fmt.Println("Error", err)
	}

	fmt.Printf("%+v\n", v)
}
