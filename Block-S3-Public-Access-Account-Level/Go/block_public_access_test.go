// Module containing unit tests for the main.go module

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"gotest.tools/assert"
)

func TestParseResults(t *testing.T) {
	/* 
	This test is used to check whether the 'parseResults' function
	returns the correct boolean based on the arg passed.
	*/

	respSuc := parseResults("false")
	assert.Equal(t, true, respSuc)

	respFal := parseResults("true")
	assert.Equal(t, false, respFal)
}

func TestPutPublicAccessBlock(t *testing.T) {
	/*
	This tests the functionality of the 'putPublicAccessBlock' function.

	It asserts that a boolean true will be returned if the function successfully calls AWS API.
	*/
	mockedS3ControlActionsAPI := &S3ControlActionsAPIMock{
		PutPublicAccessBlockFunc: func(ctx context.Context, params *s3control.PutPublicAccessBlockInput, optFns ...func(*s3control.Options)) (*s3control.PutPublicAccessBlockOutput, error) {
			var s3ControlOutput s3control.PutPublicAccessBlockOutput
			// return a blank payload as this method does not
			// return anything from AWS
			return &s3ControlOutput,nil;

		},
	}
	resp := putPublicAccessBlock(mockedS3ControlActionsAPI)
	assert.Equal(t, true, resp)
}


func TestGetPublicAccessBlockOpen(t *testing.T) {
	/*
	This tests the functionality of the 'getPublicAccessBlock' function.

	It asserts that a boolean false will be returned if the AWS account being checked is NOT blocked
	from public access.
	*/
	
	mockedS3ControlActionsAPI := &S3ControlActionsAPIMock{
		GetPublicAccessBlockFunc: func(ctx context.Context, params *s3control.GetPublicAccessBlockInput, optFns ...func(*s3control.Options)) (*s3control.GetPublicAccessBlockOutput, error) {
			var s3ControlOutput s3control.GetPublicAccessBlockOutput
			// read a json file and return the content for mock
			data, _ := ioutil.ReadFile("test_data/open_account.json")

			// this unmarshal will convert the json data into a usable mock format
			json.Unmarshal(data, &s3ControlOutput);
			return &s3ControlOutput,nil;
		},
	}

	resp := getPublicAccessBlock(mockedS3ControlActionsAPI)
	assert.Equal(t, false, resp)
}

func TestGetPublicAccessBlockClosed(t *testing.T) {
	/*
	This tests the functionality of the 'getPublicAccessBlock' function.

	It asserts that a boolean true will be returned if the AWS account being checked is blocked 
	from public access.
	*/
	
	mockedS3ControlActionsAPI := &S3ControlActionsAPIMock{
		GetPublicAccessBlockFunc: func(ctx context.Context, params *s3control.GetPublicAccessBlockInput, optFns ...func(*s3control.Options)) (*s3control.GetPublicAccessBlockOutput, error) {
			var s3ControlOutput s3control.GetPublicAccessBlockOutput
			// read a json file and return the content for mock
			data, _ := ioutil.ReadFile("test_data/closed_account.json")

			// this unmarshal will convert the json data into a usable mock format
			json.Unmarshal(data, &s3ControlOutput);
			return &s3ControlOutput,nil;
		},
	}

	resp := getPublicAccessBlock(mockedS3ControlActionsAPI)
	assert.Equal(t, true, resp)
}