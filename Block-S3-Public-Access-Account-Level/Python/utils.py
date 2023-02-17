"""Utility funcs to be imported when needed."""

# Standard Library imports
import logging
import os

# Third-party imports
import boto3
from botocore.exceptions import ClientError

#################
# ENV Vars
#################
LOG_LEVEL = os.environ.get("log_level")


#################
# Create logger
#################
def create_handler():
    log = logging.getLogger()
    log.setLevel(LOG_LEVEL)

    logging.getLogger('boto').setLevel(logging.CRITICAL)
    logging.getLogger('botocore').setLevel(logging.CRITICAL)
    return log


logger = create_handler()


def error_handler(func):
    """Take a func and pass to the inner_func"""

    def inner_func(*args, **kwargs):
        """
        Take the params passed to the decorator and passed func.

        :param args: Any number of args
        :param kwargs: Any number of keyword args
        :return: The results of the func being executed
        """
        try:
            result = func(*args, **kwargs)
            return result
        except ClientError as err:
            logger.error("[-] ClientError: error=%s, func=%s", err, func.__name__)
        except KeyError as err:
            logger.error("[-] KeyError: error=%s, func=%s", err, func.__name__)
        except Exception as err:
            logger.error("[-] General Exception: error=%s, func=%s", err, func.__name__)

    return inner_func


@error_handler
def get_account_id(sts_client: boto3.client) -> str:
    """
    Func that gets the current AWS account id.

    :param sts_client: (required) An instantiated boto3 STS client
    :return: A string containing the AWS account ID
    """
    return sts_client.get_caller_identity()['Account']


@error_handler
def create_boto3(
        service: str,
        boto_client: str = False,
        boto_resource: str = False,
        region: str = None,
        access_key: str = None,
        secret_key: str = None,
        session_token: str = None
) -> boto3.client:
    """
    Create a boto3 client or resource based on AWS service passed (e.g. 'sts', 's3)

    :param service: (required) AWS service passed as string
    :param boto_client: (*required) If a boto client is needed change to True
    :param boto_resource: (*required) If a boto resource is needed change to True
    :param region: (optional) AWS region passed as string
    :param access_key: (optional) AWS STS Access Key string obtained for cross-account access
    :param secret_key: (optional) AWS STS Secret Key string obtained for cross-account access
    :param session_token: (optional) AWS STS Session token string obtained for cross-account access
    :return: An instantiated boto3 client or resource
    :raise AWS API boto3 ClientErrors
    """
    if boto_client:
        return boto3.client(service, region_name=region, aws_access_key_id=access_key, aws_secret_access_key=secret_key,
                            aws_session_token=session_token)
    elif boto_resource:
        return boto3.resource(service, region_name=region, aws_access_key_id=access_key,
                              aws_secret_access_key=secret_key,
                              aws_session_token=session_token)
