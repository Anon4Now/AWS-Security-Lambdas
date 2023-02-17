package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var client *s3.Client

func init() {
	// initalizer func that will create
	// the usable s3control client

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	client = s3.NewFromConfig(cfg)
}

func HandleRequest(ctx context.Context) {
	b := Bucket{Client: client,}
	b.Dispatch()
}

func main() {
	lambda.Start(HandleRequest)
}