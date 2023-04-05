/* Main package that will run the code to block pub access */

package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/s3control/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// interface that implments all of the AWS API calls needed
// provides the ability for mocks during testing
//go:generate moq -out s3control_moq_test.go . S3ControlActionsAPI
type S3ControlActionsAPI interface {
	PutPublicAccessBlock(ctx context.Context, params *s3control.PutPublicAccessBlockInput, optFns ...func(*s3control.Options)) (*s3control.PutPublicAccessBlockOutput, error)
	GetPublicAccessBlock(ctx context.Context, params *s3control.GetPublicAccessBlockInput, optFns ...func(*s3control.Options)) (*s3control.GetPublicAccessBlockOutput, error)
}

func parseResults(status string) bool {
	return strings.Contains(status, "false")
}

func getAccountId(client *sts.Client) string {
	resp, err := client.GetCallerIdentity(context.TODO(), nil)

	if err != nil {
		panic("Cannot obtain AWS account ID" + err.Error())
	}
	return *resp.Account
}

func putPublicAccessBlock(client S3ControlActionsAPI, accountID string) bool {
	params := &s3control.PutPublicAccessBlockInput {
		AccountId: aws.String(accountID),
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration {
			BlockPublicAcls: true,
			BlockPublicPolicy: true,
			IgnorePublicAcls: true,
			RestrictPublicBuckets: true,
		},
	}

	_, err := client.PutPublicAccessBlock(context.TODO(), params)

	if err != nil {
		panic("AWS API Error" + err.Error())
	}

	return true

}

func getPublicAccessBlock(client S3ControlActionsAPI, accountID string) bool {
	params := &s3control.GetPublicAccessBlockInput {
		AccountId: aws.String(accountID),
	}

	resp, err := client.GetPublicAccessBlock(context.TODO(), params)
	if err != nil {
		log.Println(err)
		return false
	}

	b, err := json.Marshal(resp)
	log.Println(string(b))

	if parseResults(string(b)) {
		return false
	} else {
		return true
	}
}

func HandleRequest(ctx context.Context) {
	// load the SDK client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	s3ControlClient := s3control.NewFromConfig(cfg)
	stsClient := sts.NewFromConfig(cfg)

	acctId := getAccountId(stsClient)

	resp := getPublicAccessBlock(s3ControlClient, acctId)

	if !resp {
		log.Println("[!] Account access is not set to block, attempting to set config.")
		putPublicAccessBlock(s3ControlClient, acctId)
		log.Println("[+] Successfully set the AWS account access to public block.")
	} else {
		log.Println("[+] Account access to S3 already set to block.")
	}
}

func main() {
	lambda.Start(HandleRequest)
}