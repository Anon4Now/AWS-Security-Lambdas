"""Module that contains code that will enable versioning on S3 buckets."""

# Standard Library imports
from typing import List, Any, Dict, Optional
import re
import copy

# Third-party imports
import boto3

# Local App imports
from utils import (
    logger,
    create_boto3,
    error_handler,
    get_account_id
)


class BucketVersion:
    """Action-oriented class for checking and modifying S3 bucket versioning."""

    def __init__(self, s3_boto_client: boto3.client, current_account: str) -> None:
        """
        Init method with two params.

        :param s3_boto_client: (required) An instantiated S3 boto3 client
        :param current_account: (required) A string containing the current account ID
        :return None
        """
        self._boto_client = s3_boto_client
        self._account = current_account
        self._bucket_status_dict = {}

    @property
    @error_handler
    def _bucket_list(self) -> List[str]:
        """
        A class property that gets a list of buckets for versioning check.

        :return: A list of bucket names as strings
        """
        response = self._boto_client.list_buckets()
        return [bucket['Name'] for bucket in response['Buckets'] if 'Name' in bucket]

    @error_handler
    def _update_bucket_version(self, bucket: str, message: str) -> Dict[str, Any]:
        """
        Update the bucket based on the name passed from _check_bucket_version method.

        :param bucket: (required) A string containing an S3 bucket name to have versioning enabled
        :param message: (required) A string containing either 'suspended' or 'disabled'
        :return: The default AWS ResponseMetadata
        """
        logger.info("[!] S3 bucket '%s' in account '%s' has a current status of '%s', attempting to update.",
                    bucket, self._account, message)
        return self._boto_client.put_bucket_versioning(
            Bucket=bucket,
            VersioningConfiguration={
                'MFADelete': 'Disabled',
                'Status': 'Enabled'
            },
        )

    @error_handler
    def _get_bucket_version(self, bucket: str) -> dict:
        """
        Method to return the current versioning status of a bucket.

        :param bucket: (required) A string containing the name of an S3 bucket
        :return: Dictionary containing the status of versioning for the bucket
        """
        return self._boto_client.get_bucket_versioning(Bucket=bucket)

    @error_handler
    def _check_bucket_version(self) -> None:
        """
        Method to check a bucket for versioning and add as k/v pair to dict.

        Use the 'counter' to concatenate a string-num to the status, to
        avoid conflicts within the dict of the same key.
        (e.g. 'disabled0', 'disabled1', 'enabled0')

        Will update the instance attr 'self._bucket_status_dict' with the data.

        :return: None
        """

        for index, bucket in enumerate(self._bucket_list):
            response = self._get_bucket_version(bucket)
            if 'Status' in response and response['Status'] == 'Suspended':
                self._bucket_status_dict[f'suspended{str(index)}'] = bucket
            elif 'Status' not in response:
                self._bucket_status_dict[f'disabled{str(index)}'] = bucket

    @staticmethod
    @error_handler
    def _remove_item_from_dict(dict_key: str, dictionary: dict) -> bool:
        """
        Use param key to remove item from param dict.

        :param dict_key: (required) String value that relates to a dict key
        :param dictionary: (required) Dict containing str key and str vals (i.e. {'key': 'val'}
        :return: A bool that determines whether k/v was removed from the dict
        """
        if dict_key in dictionary:
            del dictionary[dict_key]
            return True
        return False

    @error_handler
    def dispatch(self) -> Optional[Dict[str, int]]:
        """
        Main method that calls the other methods to check versioning on a bucket.

        After it evaluates a bucket status and logs out the status, it will remove
        that entry from the dict to avoid duplicates results. This method will
        return either a Dict containing {'error': -1} to signify
        that the bucket list is empty or None which signifies
        successful execution.

        :return: Dictionary containing {'error': -1} or None
        """
        if not self._bucket_list:
            return {'error': -1}
        self._check_bucket_version()

        for status, bucket in self._bucket_status_dict.items():
            # create shallow copy of the dict to avoid iterate dict alteration error
            copy_dict = copy.copy(self._bucket_status_dict)
            # strip off the number from end the status (i.e. disabled0 -> disabled)
            stripped_status = re.search(r'(\D*)', status)[0]
            if stripped_status in ['suspended', 'disabled'] and self._remove_item_from_dict(status, copy_dict):
                self._update_bucket_version(bucket, stripped_status)
                logger.info("[+] Successfully enabled versioning on bucket: %s in account %s", bucket, self._account)
                # reassign the status_dict var to the shallow copy that has been altered
                self._bucket_status_dict = copy_dict


@error_handler
def lambda_handler(event, context) -> None:
    """
    Pass the resource 's3' to the class for instantiation and call the run method.

    :param event: (required) Default param for JSON events sent from custom resource trigger(s)
    :param context: (optional) Lambda execution context
    :return: None
    """
    s3_boto_client = create_boto3('s3', 'boto_client')
    current_account = get_account_id(create_boto3('sts', 'boto_client'))

    logger.info("[!] Checking S3 buckets for versioning")
    s3_object = BucketVersion(s3_boto_client, current_account)
    s3_object.dispatch()
    logger.info("[+] Completed bucket version checks")
