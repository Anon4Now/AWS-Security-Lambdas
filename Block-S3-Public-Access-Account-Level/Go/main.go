/* Main package that will run the code to block pub access */

package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/s3control/types"
)

// Set the package level SDK client
var client *s3control.Client

func init() {
	// initalizer func that will create
	// the usable s3control client

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	client = s3control.NewFromConfig(cfg)
}

func putPublicAccessBlock() {
	params := &s3control.PutPublicAccessBlockInput {
		AccountId: aws.String("ACCOUNT_ID"),
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration {
			BlockPublicAcls: true,
			BlockPublicPolicy: true,
			IgnorePublicAcls: true,
			RestrictPublicBuckets: true,
		},
	}

	resp, err := client.PutPublicAccessBlock(context.TODO(), params)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	printOut(string(b))

}

func getPublicAccessBlock() {
	params := &s3control.GetPublicAccessBlock {
		AccountId: aws.String("ACCOUNT_ID"),
	}

	resp, err := client.GetPublicAccessBlock(context.TODO(), params)
	if err != nil {
		putPublicAccessBlock()
	}

	b, err := json.Marshal(resp)

	if parseResults(string(b)) {
		putPublicAccessBlock()
	} else {
		printOut("Account is blocked from public access.")
	}
}

func HandleRequest(ctx context.Context) {
	getPublicAccessBlock()
}

func main() {
	lambda.Start(HandleRequest)
}