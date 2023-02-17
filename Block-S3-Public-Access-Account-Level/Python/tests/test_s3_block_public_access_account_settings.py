"""Module containing tests for the s3-block-public-access-settings-script."""

# Standard Library imports
import unittest
from unittest.mock import patch

# Third-party imports
import boto3
from botocore.exceptions import ClientError
from moto import mock_s3control, mock_sts

# Local App imports
from s3_block_public_access_account_settings import (
    _get_public_access_block_settings,
    _parse_settings_data,
    lambda_handler
)

from . import MockedSuccessClient, MockedFailClient


class TestS3BlockPublicAccess(unittest.TestCase):
    """Test the lambda function that will block s3 public access (acct level)."""

    mock_s3control = mock_s3control()
    mock_sts = mock_sts()

    def setUp(self) -> None:
        """Set up the moto mocks for test usage."""
        self.mock_s3control.start()
        self.mock_sts.start()

        self.sts = boto3.client('sts', 'us-east-2')
        self.s3control = boto3.client('s3control', 'us-east-2')

    def test_get_public_access_block_settings(self):
        """Test getting the public access block settings for the account."""

        # get a faked account number
        account_id = self.sts.get_caller_identity()['Account']

        # check that is the 'NoSuchPublicAccessBlockConfiguration' ClientError appears the below dict is returned
        self.assertEqual(_get_public_access_block_settings(self.s3control, account_id), {'error', -1})

        # check that is ANY ClientError other than 'NoSuchPublicAccessBlockConfiguration' it is raised
        with self.assertRaises(ClientError):
            _get_public_access_block_settings(self.s3control, '111111111111')

    def test_parse_settings_data(self):
        """Test that the function will parse a passed dict successfully."""

        # this data represents an account setting that is blocking public access
        input_data_all_true = {
            'PublicAccessBlockConfiguration': {
                'BlockPublicAcls': True,
                'IgnorePublicAcls': True,
                'BlockPublicPolicy': True,
                'RestrictPublicBuckets': True
            }
        }

        input_data_mixed = {
            'PublicAccessBlockConfiguration': {
                'BlockPublicAcls': True,
                'IgnorePublicAcls': True,
                'BlockPublicPolicy': True,
                'RestrictPublicBuckets': False
            }
        }

        # this data represents an account that has no public access settings configured
        input_data_error = {'error': -1}

        # this should result in a dict containing an empty list (i.e. {'result':[]})
        # which means that there were no False values found in the 'PublicAccessBlockConfiguration' settings
        self.assertFalse(_parse_settings_data(input_data_all_true)['result'])

        # this should result in a dict containing a list that contains some content (i.e. {'result':[False]})
        # which means that there were AT LEAST ONE False value(s) found in the 'PublicAccessBlockConfiguration' settings
        self.assertTrue(_parse_settings_data(input_data_mixed)['result'])

        # this should result in a dict containing a negative integer (i.e. {'error':-1})
        # which means that the 'PublicAccessBlockConfiguration' settings don't currently exist
        self.assertEqual(_parse_settings_data(input_data_error), {'error': -1})

    @patch('s3_block_public_access_account_settings.get_account_id')
    @patch('s3_block_public_access_account_settings.create_boto3')
    def test_lambda_handler(self, mock_create_boto3, mock_get_account_id):
        """This tests the lambda code by using a mocked boto3 client."""

        mock_create_boto3.return_value = MockedFailClient()
        mock_get_account_id.return_value = '11111111111111'

        with self.assertLogs() as captured:
            lambda_handler('', '')
            self.assertIn("Account does not have public access block setting enabled", captured.output[0])
            self.assertIn("Account enabling ended with a status of 'success'", captured.output[1])

    def tearDown(self) -> None:
        """Stop the moto mocks"""
        self.mock_sts.stop()
        self.mock_s3control.stop()
