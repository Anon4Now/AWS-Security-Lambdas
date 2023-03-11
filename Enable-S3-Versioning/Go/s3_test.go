package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gotest.tools/assert"
)


func TestListBuckets(t *testing.T) {
	/*
	This test is used to test the functionality of the 'bucketList' method.

	This will assert that the struct attribute slice successfully is filled
	with S3 bucket names.
	*/

	mockedS3ActionsApi := &S3ActionsApiMock{
		ListBucketsFunc: func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
			
			var s3Output s3.ListBucketsOutput

			//read json from file containing bucket list data
			data, _ := ioutil.ReadFile("test_data/bucket-list-data.json")

			// this unmarchal will retun the data to correct format for test
			json.Unmarshal(data, &s3Output);
			return &s3Output,nil;

		},
	}
	b := Bucket{Client: mockedS3ActionsApi}
	b.bucketList()
	assert.Equal(t, "bucket1", b.BucketList[0])
}

func TestPutBucketVersioning(t *testing.T) {
	/*
	This test is used to test the functionality of the 'updateBucketVersion' method.

	This will assert that the method returns a boolean 'true' if successfully run.
	*/
	
	mockedS3ActionsApi := &S3ActionsApiMock{
		PutBucketVersioningFunc: func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
			
			var s3Output s3.PutBucketVersioningOutput
			// return a blank payload as this method does not
			// return anything from AWS
			return &s3Output,nil;

		},
	}
	b := Bucket{Client: mockedS3ActionsApi}
	result := b.updateBucketVersion("bucket1")
	assert.Equal(t, true, result)

}

func TestGetBucketVersioning(t *testing.T) {
	/*
	This test is used to test the functionality of the 'getBucketVersion' method.

	This will assert that the method returns a string containing the status of the 
	S3 buckets versioning.
	*/
	mockedS3ActionsApi := &S3ActionsApiMock{
		GetBucketVersioningFunc: func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
			
			var s3Output s3.GetBucketVersioningOutput
			// read a json file and return the content for mock
			data, _ := ioutil.ReadFile("test_data/get-bucket-versioning-data.json")

			// this unmarshal will convert the json data into a usable mock format
			json.Unmarshal(data, &s3Output);
			return &s3Output,nil;

		},
	}

	b := Bucket{Client: mockedS3ActionsApi}
	resp, _ := b.getBucketVersion("bucket1")
	assert.Equal(t, "Enabled", string(resp.Status))
}

func TestCheckBucketVersion(t *testing.T) {
	/*
	This test is used to test the functionality of the 'checkBucketVersion' method.

	This will assert that the method returns a map containing a string key and key.
	The key will will be the versioning status (i.e., "enabled", "suspended", "disabled")
	plus the index of the bucket name in the struct slice.
	*/
	mockedS3ActionsApi := &S3ActionsApiMock{
		GetBucketVersioningFunc: func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
			
			var s3Output s3.GetBucketVersioningOutput
			// read a json file and return the content for mock
			data, _ := ioutil.ReadFile("test_data/get-bucket-versioning-data.json")

			// this unmarshal will convert the json data into a usable mock format
			json.Unmarshal(data, &s3Output);
			return &s3Output,nil;

		},
	}

	b := Bucket{Client: mockedS3ActionsApi}
	b.BucketList = append(b.BucketList, "bucket1")
	b.BucketList = append(b.BucketList, "bucket2")
	result := b.checkBucketVersion()
	assert.Equal(t, "bucket2", result["enabled1"])
}

func TestRemoveItemFromDict(t *testing.T){
	/*
	This test is used to test the functionality of the 'removeItemFromDict' method.

	This will assert that the method returns a boolean value depending on whether 
	the item is removed from the map.
	*/
	mockedS3ActionsApi := &S3ActionsApiMock{
		GetBucketVersioningFunc: func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
			
			var s3Output s3.GetBucketVersioningOutput
			// read a json file and return the content for mock
			data, _ := ioutil.ReadFile("test_data/get-bucket-versioning-data.json")

			// this unmarshal will convert the json data into a usable mock format
			json.Unmarshal(data, &s3Output);
			return &s3Output,nil;

		},
	}

	b := Bucket{Client: mockedS3ActionsApi}
	b.BucketList = append(b.BucketList, "bucket1")
	b.BucketList = append(b.BucketList, "bucket2")
	
	bucketVersion := b.checkBucketVersion()
	result1 := b.removeItemFromDict("enabled0", bucketVersion)
	result2 := b.removeItemFromDict("enabled", bucketVersion)
	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)

}

func TestDispatch(t *testing.T) {
	// THIS IS A FUNCTIONAL TEST NOT A UNIT TEST FOR DISPATCH METHOD
	mockedS3ActionsApi := &S3ActionsApiMock{
		ListBucketsFunc: func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
			
			var s3Output s3.ListBucketsOutput

			//read json from file containing bucket list data
			data, _ := ioutil.ReadFile("test_data/bucket-list-data.json")

			// this unmarchal will retun the data to correct format for test
			json.Unmarshal(data, &s3Output);
			return &s3Output,nil;

		},
		PutBucketVersioningFunc: func(ctx context.Context, params *s3.PutBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.PutBucketVersioningOutput, error) {
			
			var s3Output s3.PutBucketVersioningOutput
			// return a blank payload as this method does not
			// return anything from AWS
			return &s3Output,nil;

		},
		GetBucketVersioningFunc: func(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
			
			var s3Output s3.GetBucketVersioningOutput
			// read a json file and return the content for mock
			data, _ := ioutil.ReadFile("test_data/get-bucket-versioning-data-disabled.json")

			// this unmarshal will convert the json data into a usable mock format
			json.Unmarshal(data, &s3Output);
			return &s3Output,nil;

		},
	}
	b := Bucket{Client: mockedS3ActionsApi}
	b.BucketList = append(b.BucketList, "bucket2")
	b.BucketList = append(b.BucketList, "bucket3")
	b.Dispatch()

}