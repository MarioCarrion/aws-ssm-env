package main

import (
	"context"
	"fmt"

	awsssmenv "github.com/MarioCarrion/aws-ssm-env"
)

func main() {
	v := struct {
		Username string `ssm:"USER"`
	}{}

	if err := awsssmenv.Get(context.Background(), &v, nil); err != nil {
		fmt.Println("Error", err)
	}

	fmt.Printf("%+v\n", v)
}
