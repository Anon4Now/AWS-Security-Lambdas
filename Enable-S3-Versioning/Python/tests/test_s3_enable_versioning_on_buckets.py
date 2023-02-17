"""Module containing test for s3_enable_versioning_on_buckets script."""

# Standard Library imports
import unittest
from unittest.mock import patch

# Third-party imports
import boto3
from moto import mock_sts, mock_s3

# Local App imports
from s3_enable_versioning_on_buckets import (
    BucketVersion,
    lambda_handler
)


class TestS3EnableVersioning(unittest.TestCase):
    """Tests for whether script enables versioning on S3 buckets."""
    mock_sts = mock_sts()
    mock_s3 = mock_s3()

    def setUp(self) -> None:
        """Sets up the moto mocks."""

        self.mock_sts.start()
        self.mock_s3.start()

        self.s3_client = boto3.client('s3', 'us-east-1')
        self.sts_client = boto3.client('sts', 'us-east-1')

        self.s3_client.create_bucket(Bucket='bucket1', CreateBucketConfiguration={'LocationConstraint': 'us-east-1'})
        self.s3_client.create_bucket(Bucket='bucket2', CreateBucketConfiguration={'LocationConstraint': 'us-east-1'})
        self.s3_client.create_bucket(Bucket='bucket3', CreateBucketConfiguration={'LocationConstraint': 'us-east-1'})

    def test_get_bucket_list(self):
        """Tests to see if a list containing buckets is returned."""
        self.s3_obj = BucketVersion(self.s3_client, '111111111111')
        self.assertEqual(self.s3_obj._bucket_list, ['bucket1', 'bucket2', 'bucket3'])

    def test_update_bucket_version(self):
        """Tests to see wnether versioning is enabled on bucket."""

        self.s3_obj = BucketVersion(self.s3_client, '1111111111111')

        # this returns data based on a suspended bucket
        response = self.s3_obj._update_bucket_version('bucket1', 'suspended')
        self.assertEqual(response['ResponseMetadata']['HTTPStatusCode'], 200)

        # this retusn data based on a disabled bucket
        response = self.s3_obj._update_bucket_version('bucket1', 'disabled')
        self.assertEqual(response['ResponseMetadata']['HTTPStatusCode'], 200)

    def test_check_bucket_version(self):
        """Tests the logic of evaluating buckets for versioning status."""

        # this will test that all buckets are in a disabled state
        self.s3_obj = BucketVersion(self.s3_client, '1111111111111')
        self.s3_obj._check_bucket_version()
        self.assertEqual(self.s3_obj._bucket_status_dict,
                         {'disabled0': 'bucket1', 'disabled1': 'bucket2', 'disabled3': 'bucket3'})

        # this will test if some bucket(s) are in a suspended state
        self.s3_obj = BucketVersion(self.s3_client, '111111111111')
        self.s3_client.put_bucket_versioning(
            Bucket='bucket1',
            VersioningConfiguration={
                'MFADelete': 'Disabled',
                'Status': 'Suspended'
            }
        )

        self.s3_obj._check_bucket_version()
        self.assertEqual(self.s3_obj._bucket_status_dict,
                         {'suspended0': 'bucket1', 'disabled1': 'bucket2', 'disabled2': 'bucket3'})

        # this will test if all the bucket(s) are in an enabled state
        self.s3_obj = BucketVersion(self.s3_client, '111111111111')
        for bucket in self.s3_obj._bucket_list:
            self.s3_client.put_bucket_versioning(
                Bucket=bucket,
                VersioningConfiguration={
                    'MFADelete': 'Disabled',
                    'Status': 'Enabled'
                }
            )
        self.s3_obj._check_bucket_version()
        self.assertEqual(self.s3_obj._bucket_status_dict, {})

    def test_remove_item_from_dict(self):
        """Tests to see if an item is removed from a dictionary."""
        self.s3_obj = BucketVersion(self.s3_client, '111111111111')
        test_dict = {'item1': '1', 'item2': '2', 'item3': '3'}
        self.assertTrue(self.s3_obj._remove_item_from_dict('item1', test_dict))
        self.assertFalse(self.s3_obj._remove_item_from_dict('key_doesnt_exist', test_dict))

    def test_dispatch_calls_success(self):
        """Tests the dispatcher function to ensure it handles action events."""
        self.s3_obj = BucketVersion(self.s3_client, '111111111111')
        with self.assertLogs() as captured:
            self.s3_obj.dispatch()
            self.assertIn("Successfully enabled versioning on bucket", captured.output[1])

    def test_dispatch_calls_error(self):
        """Tests the dispatcher function to ensure it handles error event."""
        self.s3_obj = BucketVersion(self.s3_client, '111111111111')
        for bucket in self.s3_obj._bucket_list:
            self.s3_client.delete_bucket(Bucket=bucket)
        self.assertEqual(self.s3_obj.dispatch(), {'error': -1})

    @patch('s3_enable_versioning_on_buckets.get_acount_id')
    @patch('s3_enable_versioning_on_buckets.create_boto3')
    def test_lambda_handler(self, mock_create_boto3, mock_get_account_id):
        """Tests the lambda handler code."""
        mock_create_boto3.return_value = self.s3_client
        mock_get_account_id.return_value = '111111111111'

        with self.assertLogs() as captured:
            lambda_handler('', '')
            self.assertIn('Completed bucket version checks', captured.output[-1])

    def tearDown(self) -> None:
        """Stops the moto mocks."""
        self.mock_sts.stop()
        self.mock_s3.stop()
