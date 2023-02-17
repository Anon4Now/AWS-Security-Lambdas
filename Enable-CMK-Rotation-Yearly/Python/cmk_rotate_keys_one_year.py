"""Module containing code to set yearly rotation on KMS keys."""

# Standard Library imports
from typing import List

# Third-party imports
import boto3

# Local App imports
from utils import (
    logger,
    error_handler,
    create_boto3,
    get_account_id
)


class KMSKey:
    """Action-oriented class to eval and remediate KMS keys."""

    def __init__(self, kms_boto_client: boto3.client, current_account: str) -> None:
        self._boto_client = kms_boto_client
        self._account = current_account

        self._active_customer_key_list = []
        self._non_rotated_key_list = []

    @property
    @error_handler
    def _key_list(self) -> List[str]:
        """
        A property that obtains a full list of the keys that exist in the account.

        :return: A list containing strings of the KMS key ID's
        """
        return [key['KeyId'] for key in self._boto_client.list_keys()['Keys']]

    @error_handler
    def _describe_key(self) -> None:
        """
        Obtain a parsed list of keys that match conditions.

        The keys that will be added to the instance attr must be
        customer created and not pending deletion.

        :return: None
        """
        for key in self._key_list:
            key_data = self._boto_client.describe_key(KeyId=key)['KeyMetadata']
            if key_data['KeyManager'] == 'CUSTOMER' and key_data['KeyState'] != 'PendingDeletion':
                self._active_customer_key_list.append(key)

    @error_handler
    def _check_key_rotation_status(self) -> None:
        """
        Obtain a list of keys that do not have rotation enabled.

        The keys that will be added to the instance attr
        must not have rotation enabled.

        :return: None
        """
        for key in self._active_customer_key_list:
            key_data = self._boto_client.get_key_rotation_status(KeyId=key)['KeyRotationEnabled']
            if not key_data:
                self._non_rotated_key_list.append(key)

    @error_handler
    def _enable_key_rotation(self) -> List[dict]:
        """
        Enable rotation on keys that do not have it.

        :return: A list containing standard AWS API ResponseMetadata
        """
        return [self._boto_client.enable_key_rotation(KeyId=key) for key in self._non_rotated_key_list]

    @error_handler
    def dispatch(self) -> None:
        """
        Dispatch the method calls in the correct order.

        :return: None
        """
        self._describe_key()
        self._check_key_rotation_status()
        if self._non_rotated_key_list:
            self._enable_key_rotation()
            logger.info("[+] Successfully added key rotation to '%s'", self._non_rotated_key_list)
            return
        logger.info("[+] No KMS key(s) found with rotation disabled, no action taken")


@error_handler
def lambda_handler(event, context) -> None:
    """
    Create the Lambda handler and call appropriate resources.

    :param event: (optional) Used to trigger events
    :param context: (optional) Execution env created by AWS
    :return: None
    """
    sts_client = create_boto3('sts', 'boto_client')
    kms_client = create_boto3('kms', 'boto_client')
    account_id = get_account_id(sts_client)

    kms_object = KMSKey(kms_client, account_id)
    kms_object.dispatch()
