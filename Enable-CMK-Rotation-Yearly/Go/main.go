package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

// interface that implments all of the AWS API calls needed
// provides the ability for mocks during testing
//go:generate moq -out kms_moq_test.go . KMSActionsAPI
type KMSActionsAPI interface {
	EnableKeyRotation(ctx context.Context, params *kms.EnableKeyRotationInput, optFns ...func(*kms.Options)) (*kms.EnableKeyRotationOutput, error)
	GetKeyRotationStatus(ctx context.Context, params *kms.GetKeyRotationStatusInput, optFns ...func(*kms.Options)) (*kms.GetKeyRotationStatusOutput, error)
	DescribeKey(ctx context.Context, params *kms.DescribeKeyInput, optFns ...func(*kms.Options)) (*kms.DescribeKeyOutput, error)
	ListKeys(ctx context.Context, params *kms.ListKeysInput, optFns ...func(*kms.Options)) (*kms.ListKeysOutput, error)

}


func setKeyRotation(client KMSActionsAPI, nonRotatedKeys []kms.DescribeKeyOutput) bool {
	/*
	Function that sets the found CMK's to rotate yearly via API call.

	:param client: An instantiated struct that contains methods matching the KMSActionsAPI interface
	:param nonRotatedKeys: A slice of key data for non-rotated CMK's
	:return: A bool result if all actions in function successfully run
	*/

	
	for _, el := range nonRotatedKeys {
		params := &kms.EnableKeyRotationInput {
			KeyId: aws.String(*el.KeyMetadata.KeyId),
		}

		_, err := client.EnableKeyRotation(context.TODO(), params)

		if err != nil {
			log.Println(err)
		}

		log.Printf("Key: %v has been set to rotate successfully.\n", el.KeyMetadata.KeyId)
	}
	return true
}

func getRotationStatus(client KMSActionsAPI, custKeys []kms.DescribeKeyOutput) []kms.DescribeKeyOutput {
	/*
	Function that finds the current rotation status of the CMK's.

	:param client: An instantiated struct that contains methods matching the KMSActionsAPI interface
	:param custKeys: A slice of key data for KMS keys in account/region that are customer managed.
	:return: A slice containing the key data for the non-rotated keys.
	*/

	var nonRotatedKeys []kms.DescribeKeyOutput

	for _, el := range custKeys {
		params := &kms.GetKeyRotationStatusInput {
			KeyId: aws.String(*el.KeyMetadata.KeyId),
		}

		resp, err := client.GetKeyRotationStatus(context.TODO(), params)

		if err != nil {
			log.Println(err)
		}

		if resp.KeyRotationEnabled == false {
			nonRotatedKeys = append(nonRotatedKeys, el)
		}
	}
	return nonRotatedKeys
}

func getCustKeys(client KMSActionsAPI, keys []types.KeyListEntry) []kms.DescribeKeyOutput {
	/*
	Function that finds CMK's in the current AWS account/region.

	:param client: An instantiated struct that contains methods matching the KMSActionsAPI interface
	:param keys: A slice containing all key data for the current AWS account/region.
	:return: A slice containing the key data for all customer managed keys.
	*/

	var custKeys []kms.DescribeKeyOutput

	for _, el := range keys {
		params := &kms.DescribeKeyInput {
			KeyId: aws.String(*el.KeyId),
		}

		resp, err := client.DescribeKey(context.TODO(), params)
		if err != nil {
			log.Println(err)
		}

		if resp.KeyMetadata.KeyManager == "CUSTOMER" && resp.KeyMetadata.KeyState != "PendingDeletion" {
			custKeys = append(custKeys, *resp)
		}
	}
	return custKeys
}


func listKeys(client KMSActionsAPI) []types.KeyListEntry {
	/*
	Function that obtains key data for all keys in the current AWS account/region.

	:param client: An instantiated struct that contains methods matching the KMSActionsAPI interface
	:return: A slice containing the key data for all keys in AWS account/region.
	*/

	resp, err := client.ListKeys(context.TODO(), nil)

	if err != nil {
		log.Println(err)
	}

	return resp.Keys
}


func HandleRequest(ctx context.Context) {
	/*
	Main handler for the Lambda that dispatches calls to functions.

	:param ctx: The default Lambda context during execution.
	:return: None
	*/

	// load the KMS client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
        panic("unable to load SDK config, " + err.Error())
	}
	client := kms.NewFromConfig(cfg)

	// get an array of KMS keys
	listOfKeys := listKeys(client)

	if len(listOfKeys) == 0 {
		log.Println("[!] No keys found in account.")
		os.Exit(0)
	}

	// get an array of CMK's
	custKeys := getCustKeys(client, listOfKeys)

	if len(custKeys) == 0 {
		log.Println("[!] No customer managed keys found in account.")
		os.Exit(0)
	}

	// get the rotation status of the CMK's
	statusOfKeys := getRotationStatus(client, custKeys)

	if len(statusOfKeys) == 0 {
		log.Println("[+] All keys set to rotate in account, no action taken.")
		os.Exit(0)
	}

	// set the CMK's to rotate
	log.Printf("[!] Attempting to set keys %v to rotate.\n", statusOfKeys)
	result := setKeyRotation(client, statusOfKeys)

	if result {
		log.Println("[+] All keys set to rotate successfully.")
	}

}


func main() {
	lambda.Start(HandleRequest)
}