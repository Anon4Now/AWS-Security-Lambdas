package main

import (
	"testing"
	"context"
	"encoding/json"
	"io/ioutil"

	"gotest.tools/assert"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

var custKeysMock []kms.DescribeKeyOutput
var keyId string = "1234abcd-12ab-34cd-56ef-1234567890ab"


func TestListKeysAndGetCustKeys(t *testing.T) {

	/*
	This test function will look at both the listKeys and getCustKeys
	functions in the main module. These fuctions call the AWS SDK methods
	'ListKeys' and 'DescribeKey' respectively. Using the moq data these client calls will
	be mocked and example json data will be returned which will then
	be used by the induvidual functions in the main module.
	*/

		// make and configure a mocked KMSActionsAPI
		mockedKMSActionsAPIForList := &KMSActionsAPIMock{
			ListKeysFunc: func(ctx context.Context, params *kms.ListKeysInput, optFns ...func(*kms.Options)) (*kms.ListKeysOutput, error) {
				
				var kmsOutput kms.ListKeysOutput
				// Read json file contaning example data for KMS keys in an account
				data, _ := ioutil.ReadFile("test-data/kms-key-list.json")
				
				json.Unmarshal(data, &kmsOutput);
				return &kmsOutput,nil;
			},
		}
	
	
		// use mockedKMSActionsAPI in code that requires KMSActionsAPI
		// and then make assertions.
	
		listOfKeys := listKeys(mockedKMSActionsAPIForList)

		assert.Equal(t,keyId, *listOfKeys[0].KeyId)

		// make and configure a mocked KMSActionsAPI

		mockedKMSActionsAPIDescribe := &KMSActionsAPIMock{
			DescribeKeyFunc: func(ctx context.Context, params *kms.DescribeKeyInput, optFns ...func(*kms.Options)) (*kms.DescribeKeyOutput, error) {
				
				var kmsOutput kms.DescribeKeyOutput
				// Read json file containing example key detail for a single key
				data, _ := ioutil.ReadFile("test-data/key-details.json")
				
				// this unmarshal will return the data in the correct struct format
				json.Unmarshal(data, &kmsOutput);
				return &kmsOutput,nil;
			},
		}
		customerManagedKeys := getCustKeys(mockedKMSActionsAPIDescribe, listOfKeys)

		// this append process is so that the global var can be used in the following tests
		custKeysMock = append(custKeysMock, customerManagedKeys[0])

		assert.Equal(t, keyId, *customerManagedKeys[0].KeyMetadata.KeyId)

}


func TestGetKeyRotationStatus(t *testing.T) {
	/*
	This test function will look at both the getRotationStatus
	function in the main module. This function will call the AWS SDK method
	'GetKeyRotationStatus'. Using the moq data this client call will
	be mocked and example json data will be returned which will then
	be used by the function in the main module.
	*/

		// make and configure a mocked KMSActionsAPI
		mockedKMSActionsAPI := &KMSActionsAPIMock{
			GetKeyRotationStatusFunc: func(ctx context.Context, params *kms.GetKeyRotationStatusInput, optFns ...func(*kms.Options)) (*kms.GetKeyRotationStatusOutput, error) {
				
				var kmsOutput kms.GetKeyRotationStatusOutput
				// Read json file containing example data for a non-rotated key status
				data, _ := ioutil.ReadFile("test-data/get-rotation-status.json")
				
				// this unmarshal will return the data in the correct struct format
				json.Unmarshal(data, &kmsOutput);
				return &kmsOutput,nil;
			},
		}

		nonRotatedKeyList := getRotationStatus(mockedKMSActionsAPI, custKeysMock)

		assert.Equal(t, keyId, *nonRotatedKeyList[0].KeyMetadata.KeyId)


}

func TestSetKeyRotation(t *testing.T) {
	/*
	This test function will look at both the setKeyRotation
	function in the main module. This function will call the AWS SDK method
	'EnableKeyRotation'. Using the moq data this client call will
	be mocked and example json data will be returned which will then
	be used by the function in the main module.
	*/

	// mocked KMS client call to enablekeyrotation method
	mockedKMSActionsAPI := &KMSActionsAPIMock{
		EnableKeyRotationFunc: func(ctx context.Context, params *kms.EnableKeyRotationInput, optFns ...func(*kms.Options)) (*kms.EnableKeyRotationOutput, error) {
			
			// return a blank payload as this method does not
			// return anything from AWS
			var kmsOutput kms.EnableKeyRotationOutput
			return &kmsOutput,nil;
		},
	}

	result := setKeyRotation(mockedKMSActionsAPI, custKeysMock)

	assert.Equal(t, true, result)


}