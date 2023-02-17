package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

var client *kms.client
var custKeys []kms.DescribeKeyOutput
var nonRotatedKeys []kms.DescribeKeyOutput

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	client = kms.NewFromConfig(cfg)
}

func setKeyRotation() {
	for _, el := range nonRotatedKeys {
		params := &kms.EnableKeyRotationInput {
			KeyId: aws.String(*el.KeyMetadata.KeyId),
		}

		_, err := client.EnableKeyRotation(context.TODO(), params)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Key: %v has been set to rotate successfully.\n", el.KeyMetadata.KeyId)
	}
}

func getRotationStatus() {
	for _, el := range custKeys {
		params := &kms.GetKeyRotationStatusInput {
			KeyId: aws.String(*el.KeyMetadata.KeyId),
		}

		resp, err := client.GetKeyRotationStatus(context.TODO(), params)

		if err != nil {
			fmt.Println(err)
		}

		if resp.KeyRotationEnabled == false {
			nonRotatedKeys = append(nonRotatedKeys, el)
		}
	}
}

func getCustKeys(keys []types.KeyListEntry) {
	for _, el := range keys {
		params := &kms.DescribeKeyInput {
			KeyId: aws.String(*el.KeyId),
		}

		resp, err := client.DescribeKey(context.TODO(), params)
		if err != nil {
			fmt.Println(err)
		}

		if resp.KeyMetadata.KeyManager == "CUSTOMER" && resp.KeyMetadata.KeyState != "PendingDeletion" {
			custKeys = append(custKeys, *resp)
		}
	}
}


func listKeys() []types.KeyListEntry {
	resp, err := client.ListKeys(context.TODO(), nil)

	if err != nil {
		fmt.Println(err)
	}

	return resp.Keys
}