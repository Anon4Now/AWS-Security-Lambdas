"""Main script that checks/updates the S3 public access settings."""

# Standard Library imports
from typing import Any, Dict

# Third-party imports
from botocore.exceptions import ClientError
import boto3

# Local app imports
from utils import (
    logger,
    create_boto3,
    error_handler,
    get_account_id
)


@error_handler
def _update_s3_account_access_settings(s3_client: boto3.client, account_id: str) -> Dict[str, Any]:
    """
    Perform the update on the S3 account to block public access to buckets.

    :param s3_client: (required) An instantiated boto3 client for S3
    :param account_id:  (required) A string containing the current account id
    :return: Standard AWS API response
    """
    return s3_client.put_public_access_block(
        PublicAccessBlockConfiguration={
            'BlockPublicAcls': True,
            'IgnorePublicAcls': True,
            'BlockPublicPolicy': True,
            'RestrictPublicBuckets': True,
        },
        AccountId=account_id,
    )


# DO NOT DECORATE WITH ERROR_HANDLER
def _get_public_access_block_settings(s3_client: boto3.client, account_id: str) -> Dict[str, Any]:
    """
    Function that obtains the current account settings data.

    The error_handler CANNOT be used here, as the AWS API will return a ClientError
    if the public_access_block has never been set. This error is used
    by the _parse_settings_data func to determine next steps.

    :param s3_client: (required) An instantiated boto3 S3 client
    :param account_id: (required) A string containing the current account ID
    :return: A dict containing string keys and potentially complex data structures as values
    :raise A NoSuchPublicAccessBlockConfiguration is the account has never had the S3 public
    block set before
    """
    try:
        return s3_client.get_public_access_block(AccountId=account_id)
    except ClientError as e:
        if 'NoSuchPublicAccessBlockConfiguration' in str(e):
            return {'error': -1}
        raise


@error_handler
def _parse_settings_data(data: Dict[str, Any]) -> Dict[str, Any]:
    """
    Function that parses the response data from AWS and returns a dictionary.

    The possible returned data includes the three options below:
    - {'result': []} which means that the account is successfully blocks public access
    - {'result': [False, etc...]} which means that a 'False' value was found settings (i.e. account block isn't complete)
    - {'error': -1} which means that the account has no public access settings (i.e. account block needs to be enabled)

    :param data: (required) A dict containing string keys, and potentially a complex data structure
    :return: Will return a dictionary that will be used to decide next steps
    """
    if 'error' in data:
        return data
    return dict(result=list(filter(lambda element: element == 0, list(data['PublicAccessBlockConfiguration'].values()))))


@error_handler
def lambda_handler(event, context) -> None:
    """
    Create the Lambda handler and call appropriate functions.

    :param event: (optional) Data that can be used as a trigger
    :param context: (optional) Execution env created for Lambda
    :return: None
    """
    sts_client = create_boto3('sts', boto_client=True)
    s3_client = create_boto3('s3control', boto_client=True)
    account_id = get_account_id(sts_client)

    parsed_data = _parse_settings_data(_get_public_access_block_settings(s3_client, account_id))
    if 'error' in parsed_data or parsed_data.get('result'):
        logger.info("[!] Account does not have public access block setting enabled, attempting to enable for %s",
                    account_id)
        response = 'success' if _update_s3_account_access_settings(s3_client, account_id) else 'fail'
        logger.info("[!] Account enabling ended with a status of '%s'", response)
        return
    logger.info("[+] S3 public access block is already set, no action taken")
