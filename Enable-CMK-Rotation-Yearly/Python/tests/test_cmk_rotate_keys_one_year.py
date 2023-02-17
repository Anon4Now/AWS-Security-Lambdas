"""Module containing tests for cmk_rotate_keys_one_year."""

# Standard Library imports
import unittest
from unittest.mock import patch

# Third-party imports
import boto3
from moto import mock_kms, mock_sts

# Local App imports
from cmk_rotate_keys_one_year import KMSKey, lambda_handler


class TestCmkRotation(unittest.TestCase):
    """Class to test the rotation of CMK keys."""

    mock_sts = mock_sts()
    mock_kms = mock_kms()

    def setUp(self) -> None:
        # start the mocks
        self.mock_sts.start()
        self.mock_kms.start()

        # create the mocked clients
        self.sts_client = boto3.client('sts', 'us-east-1')
        self.kms_client = boto3.client('kms', 'us-east-1')

        # get the mocked account id
        self.account_id = self.sts_client.get_caller_identity()['Account']

        # create two keys in the mocked account
        self.key1 = self.kms_client.create_key()
        self.key2 = self.kms_client.create_key()

        # instantiate the class with the mocked info
        self.kms_obj = KMSKey(self.kms_client, self.account_id)

    def test_describe_key(self):
        """Test the parsing for active keys in the account that are CMK and not pending deletion."""
        self.kms_obj._describe_key()
        # test that the two kms keys appear in the list
        self.assertEqual(len(self.kms_obj._active_customer_key_list), 2)

    def test_check_key_rotation(self):
        """Test the parsing of keys that are active but with no rotation enabled."""
        self.kms_obj._describe_key()
        self.kms_obj._check_key_rotation_status()

        # test that the two kms keys appear in the list
        self.assertEqual(len(self.kms_obj._non_rotated_key_list), 2)

    def test_dispatch(self):
        """Test the dispatch method for when non-rotated keys are found."""

        # test that keys exist without rotation enabled
        with self.assertLogs() as captured:
            self.kms_obj.dispatch()
            self.assertIn("Successfully added key rotation", captured.output[0])

        # test that keys do not exist with rotation disabled/or just don't exist in account
        with self.assertLogs() as captured:
            self.kms_obj_alt = KMSKey(self.kms_client, '111111111111')
            self.kms_obj_alt.dispatch()
            self.assertIn("No KMS key(s) found with rotation disabled, no action taken.", captured.output[0])

    @patch('cmk_rotate_keys_one_year.get_account_id')
    @patch('cmk_rotate_keys_one_year.create_boto3')
    def test_lambda_handler(self, mock_boto3, mock_account_id):
        """Test the lambda handler code."""
        mock_boto3.return_value = self.kms_client
        mock_account_id.return_value = self.account_id

        lambda_handler('', '')

    def tearDown(self) -> None:
        # stop the mocks
        self.mock_sts.stop()
        self.mock_kms.stop()
