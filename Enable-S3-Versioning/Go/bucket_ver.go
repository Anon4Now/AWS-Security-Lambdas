package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)


type Bucket struct {
	Client S3ActionsApi
	BucketList []string
}

func (b *Bucket) bucketList() {	
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

func (b *Bucket) updateBucketVersion(bucket string, message string) {
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
}

func (b *Bucket) getBucketVersion(bucket string) (*s3.GetBucketVersioningOutput, error) {

	params := &s3.GetBucketVersioningInput {
		Bucket: &bucket,
	}
	return b.Client.GetBucketVersioning(context.TODO(), params)
}

func (b *Bucket) checkBucketVersion() map[string]string{
	m := make(map[string]string)
	for i, bucket := range b.BucketList {
		resp, err := b.getBucketVersion(bucket)

		if err != nil {
			log.Println(err)
		}

			if resp.Status == "Suspended" {
				m["suspended" + string(i)] = bucket
			} else if resp.Status == "Enabled"{
				m["enabled" + string(i)] = bucket
			} else {
				m["disabled" + string(i)] = bucket
			}
	}
	return m
}

func (b *Bucket) removeItemFromDict(mapKey string, bucketMap map[string]string) bool {
	if _, ok := bucketMap[mapKey]; ok {
		delete(bucketMap, mapKey)
		return true
	}
	return false
}

func (b *Bucket) Dispatch() {
	b.bucketList()

	if len(b.BucketList) == 0 {
		log.Println("No buckets found in account.")
	}

	bucketMap := b.checkBucketVersion()

	for status, bucket := range bucketMap {
		copyMap := bucketMap

		if strings.Contains(status, "suspended") || strings.Contains(status, "disabled") {
			if b.removeItemFromDict(status, copyMap) {
				b.updateBucketVersion(bucket, status)
				log.Printf("Enabled versioning on bucket %v\n", bucket)

				bucketMap = copyMap
			}

		} else {
			log.Printf("Bucket %v already had versioning\n", bucket)
		}
	}

}