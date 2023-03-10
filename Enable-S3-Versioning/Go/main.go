package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// interface that implements all of the AWS API calls needed
// provides the ability for mocks during testing
//go:generate moq -out s3_moq_test.go . S3ActionsAPI
type S3ActionsApi interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error)
	GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
}


func HandleRequest(ctx context.Context) {

	// load the S3 client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
        panic("unable to load SDK config, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	b := Bucket{Client: client,}
	b.Dispatch()
}

func main() {
	lambda.Start(HandleRequest)
}