package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)


type Bucket struct {
	/*
	Struct that contains the attrs and methods to set S3 versioning on buckets.
	*/
	Client S3ActionsApi
	BucketList []string
}

func (b *Bucket) bucketList() {
	/*
	Private method that checks what buckets are available to the role in the AWS account.

	:return: nil
	*/	
	resp, err := b.Client.ListBuckets(context.TODO(), nil)

	if err != nil {
		log.Println(err)
	}

	for _, bucket := range resp.Buckets {
		if len(*bucket.Name) != 0 {
			b.BucketList = append(b.BucketList, *bucket.Name)
		}
	}

}

func (b *Bucket) updateBucketVersion(bucket string) bool {
	/*
	Private method to set versioning on a specific S3 bucket within the AWS account.

	Will return 'true' if the method runs successfully.

	:param bucket: (required) A string containing the S3 bucket name
	:return: A boolean result when run successfully
	*/
	params := &s3.PutBucketVersioningInput {
		Bucket: &bucket,
		VersioningConfiguration: &types.VersioningConfiguration {
			Status: "Enabled",
		},
	}

	_, err := b.Client.PutBucketVersioning(context.TODO(), params)

	if err != nil {
		log.Println(err)
	}

	return true
}

func (b *Bucket) getBucketVersion(bucket string) (*s3.GetBucketVersioningOutput, error) {
	/*
	Private method that will get the current status of versioning on an S3 bucket.

	:param bucket: (required) A string containing the name of the S3 bucket
	:return: Will return either a struct output or an error from AWS
	*/

	params := &s3.GetBucketVersioningInput {
		Bucket: &bucket,
	}
	return b.Client.GetBucketVersioning(context.TODO(), params)
}

func (b *Bucket) checkBucketVersion() map[string]string{
	/*
	Private method that will make a map key with the bucket version status.

	The map key will be concatenated with the index value of the struct slice.

	:return: A map with a string key and string value (i.e., {"enabled0": "bucketName", "suspended1": "bucketnane1"})
	*/
	m := make(map[string]string)
	for i, bucket := range b.BucketList {
		resp, err := b.getBucketVersion(bucket)

		if err != nil {
			log.Println(err)
		}
			// need to perform a conversion on the index number to avoid a testing error
			if resp.Status == "Suspended" {
				m["suspended" + strconv.FormatInt(int64(i), 10)] = bucket
			} else if resp.Status == "Enabled"{
				m["enabled" + strconv.FormatInt(int64(i), 10)] = bucket
			} else {
				m["disabled" + strconv.FormatInt(int64(i), 10)] = bucket
			}
	}
	return m
}

func (b *Bucket) removeItemFromDict(mapKey string, bucketMap map[string]string) bool {
	/*
	Private method that will remove items from the shallow copy of the map.

	:param mapKey: (required) The string key that will point to the val in the map to remove
	:param bucketMap: (required) A map with a string key and a string value for bucket name and status
	:return: A boolean result that depends on whether the map deletion was successful
	*/
	if _, ok := bucketMap[mapKey]; ok {
		delete(bucketMap, mapKey)
		return true
	}
	return false
}

func (b *Bucket) Dispatch() {
	/*
	Public method that will call all the private methods in the correct order.

	Will check that key contains either "suspended" or "disabled" to perform the bucket versioning.

	:return: nil
	*/
	b.bucketList()

	if len(b.BucketList) == 0 {
		log.Println("No buckets found in account.")
	}

	bucketMap := b.checkBucketVersion()

	for status, bucket := range bucketMap {
		copyMap := bucketMap

		if strings.Contains(status, "suspended") || strings.Contains(status, "disabled") {
			if b.removeItemFromDict(status, copyMap) {
				b.updateBucketVersion(bucket)
				log.Printf("Enabled versioning on bucket %v\n", bucket)

				bucketMap = copyMap
			}

		} else {
			log.Printf("Bucket %v already had versioning\n", bucket)
		}
	}

}