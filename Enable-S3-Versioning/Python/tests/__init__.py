import os

os.environ.setdefault("log_level", "INFO")


class MockedSuccessClient:
    """An object that will contain successful boto calls."""

    @property
    def client(self):
        return {}

    @staticmethod
    def get_public_access_block(*args, **kwargs):
        return {
            'PublicAccessBlockConfiguration': {
                'BlockPublicAcls': True,
                'IgnorePublicAcls': True,
                'BlockPublicPolicy': True,
                'RestrictPublicBuckets': True
            }
        }

    @staticmethod
    def put_public_access_block(*args, **kwargs):
        return True


class MockedFailClient:
    """An object that will contain failed boto calls."""

    @property
    def client(self):
        return {}

    @staticmethod
    def get_public_access_block(*args, **kwargs):
        return {
            'PublicAccessBlockConfiguration': {
                'BlockPublicAcls': False,
                'IgnorePublicAcls': False,
                'BlockPublicPolicy': True,
                'RestrictPublicBuckets': True
            }
        }

    @staticmethod
    def put_public_access_block(*args, **kwargs):
        return True
