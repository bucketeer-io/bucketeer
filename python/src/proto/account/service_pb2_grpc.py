# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from proto.account import service_pb2 as proto_dot_account_dot_service__pb2


class AccountServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetMe = channel.unary_unary(
                '/bucketeer.account.AccountService/GetMe',
                request_serializer=proto_dot_account_dot_service__pb2.GetMeRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.GetMeResponse.FromString,
                )
        self.GetMeByEmail = channel.unary_unary(
                '/bucketeer.account.AccountService/GetMeByEmail',
                request_serializer=proto_dot_account_dot_service__pb2.GetMeByEmailRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.GetMeResponse.FromString,
                )
        self.CreateAdminAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/CreateAdminAccount',
                request_serializer=proto_dot_account_dot_service__pb2.CreateAdminAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.CreateAdminAccountResponse.FromString,
                )
        self.EnableAdminAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/EnableAdminAccount',
                request_serializer=proto_dot_account_dot_service__pb2.EnableAdminAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.EnableAdminAccountResponse.FromString,
                )
        self.DisableAdminAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/DisableAdminAccount',
                request_serializer=proto_dot_account_dot_service__pb2.DisableAdminAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.DisableAdminAccountResponse.FromString,
                )
        self.GetAdminAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/GetAdminAccount',
                request_serializer=proto_dot_account_dot_service__pb2.GetAdminAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.GetAdminAccountResponse.FromString,
                )
        self.ListAdminAccounts = channel.unary_unary(
                '/bucketeer.account.AccountService/ListAdminAccounts',
                request_serializer=proto_dot_account_dot_service__pb2.ListAdminAccountsRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.ListAdminAccountsResponse.FromString,
                )
        self.ConvertAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/ConvertAccount',
                request_serializer=proto_dot_account_dot_service__pb2.ConvertAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.ConvertAccountResponse.FromString,
                )
        self.CreateAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/CreateAccount',
                request_serializer=proto_dot_account_dot_service__pb2.CreateAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.CreateAccountResponse.FromString,
                )
        self.EnableAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/EnableAccount',
                request_serializer=proto_dot_account_dot_service__pb2.EnableAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.EnableAccountResponse.FromString,
                )
        self.DisableAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/DisableAccount',
                request_serializer=proto_dot_account_dot_service__pb2.DisableAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.DisableAccountResponse.FromString,
                )
        self.ChangeAccountRole = channel.unary_unary(
                '/bucketeer.account.AccountService/ChangeAccountRole',
                request_serializer=proto_dot_account_dot_service__pb2.ChangeAccountRoleRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.ChangeAccountRoleResponse.FromString,
                )
        self.GetAccount = channel.unary_unary(
                '/bucketeer.account.AccountService/GetAccount',
                request_serializer=proto_dot_account_dot_service__pb2.GetAccountRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.GetAccountResponse.FromString,
                )
        self.ListAccounts = channel.unary_unary(
                '/bucketeer.account.AccountService/ListAccounts',
                request_serializer=proto_dot_account_dot_service__pb2.ListAccountsRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.ListAccountsResponse.FromString,
                )
        self.CreateAPIKey = channel.unary_unary(
                '/bucketeer.account.AccountService/CreateAPIKey',
                request_serializer=proto_dot_account_dot_service__pb2.CreateAPIKeyRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.CreateAPIKeyResponse.FromString,
                )
        self.ChangeAPIKeyName = channel.unary_unary(
                '/bucketeer.account.AccountService/ChangeAPIKeyName',
                request_serializer=proto_dot_account_dot_service__pb2.ChangeAPIKeyNameRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.ChangeAPIKeyNameResponse.FromString,
                )
        self.EnableAPIKey = channel.unary_unary(
                '/bucketeer.account.AccountService/EnableAPIKey',
                request_serializer=proto_dot_account_dot_service__pb2.EnableAPIKeyRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.EnableAPIKeyResponse.FromString,
                )
        self.DisableAPIKey = channel.unary_unary(
                '/bucketeer.account.AccountService/DisableAPIKey',
                request_serializer=proto_dot_account_dot_service__pb2.DisableAPIKeyRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.DisableAPIKeyResponse.FromString,
                )
        self.GetAPIKey = channel.unary_unary(
                '/bucketeer.account.AccountService/GetAPIKey',
                request_serializer=proto_dot_account_dot_service__pb2.GetAPIKeyRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.GetAPIKeyResponse.FromString,
                )
        self.ListAPIKeys = channel.unary_unary(
                '/bucketeer.account.AccountService/ListAPIKeys',
                request_serializer=proto_dot_account_dot_service__pb2.ListAPIKeysRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.ListAPIKeysResponse.FromString,
                )
        self.GetAPIKeyBySearchingAllEnvironments = channel.unary_unary(
                '/bucketeer.account.AccountService/GetAPIKeyBySearchingAllEnvironments',
                request_serializer=proto_dot_account_dot_service__pb2.GetAPIKeyBySearchingAllEnvironmentsRequest.SerializeToString,
                response_deserializer=proto_dot_account_dot_service__pb2.GetAPIKeyBySearchingAllEnvironmentsResponse.FromString,
                )


class AccountServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetMe(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetMeByEmail(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def CreateAdminAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def EnableAdminAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def DisableAdminAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetAdminAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListAdminAccounts(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ConvertAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def CreateAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def EnableAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def DisableAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ChangeAccountRole(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetAccount(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListAccounts(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def CreateAPIKey(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ChangeAPIKeyName(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def EnableAPIKey(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def DisableAPIKey(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetAPIKey(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListAPIKeys(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetAPIKeyBySearchingAllEnvironments(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AccountServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetMe': grpc.unary_unary_rpc_method_handler(
                    servicer.GetMe,
                    request_deserializer=proto_dot_account_dot_service__pb2.GetMeRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.GetMeResponse.SerializeToString,
            ),
            'GetMeByEmail': grpc.unary_unary_rpc_method_handler(
                    servicer.GetMeByEmail,
                    request_deserializer=proto_dot_account_dot_service__pb2.GetMeByEmailRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.GetMeResponse.SerializeToString,
            ),
            'CreateAdminAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateAdminAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.CreateAdminAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.CreateAdminAccountResponse.SerializeToString,
            ),
            'EnableAdminAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.EnableAdminAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.EnableAdminAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.EnableAdminAccountResponse.SerializeToString,
            ),
            'DisableAdminAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.DisableAdminAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.DisableAdminAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.DisableAdminAccountResponse.SerializeToString,
            ),
            'GetAdminAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.GetAdminAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.GetAdminAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.GetAdminAccountResponse.SerializeToString,
            ),
            'ListAdminAccounts': grpc.unary_unary_rpc_method_handler(
                    servicer.ListAdminAccounts,
                    request_deserializer=proto_dot_account_dot_service__pb2.ListAdminAccountsRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.ListAdminAccountsResponse.SerializeToString,
            ),
            'ConvertAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.ConvertAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.ConvertAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.ConvertAccountResponse.SerializeToString,
            ),
            'CreateAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.CreateAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.CreateAccountResponse.SerializeToString,
            ),
            'EnableAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.EnableAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.EnableAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.EnableAccountResponse.SerializeToString,
            ),
            'DisableAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.DisableAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.DisableAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.DisableAccountResponse.SerializeToString,
            ),
            'ChangeAccountRole': grpc.unary_unary_rpc_method_handler(
                    servicer.ChangeAccountRole,
                    request_deserializer=proto_dot_account_dot_service__pb2.ChangeAccountRoleRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.ChangeAccountRoleResponse.SerializeToString,
            ),
            'GetAccount': grpc.unary_unary_rpc_method_handler(
                    servicer.GetAccount,
                    request_deserializer=proto_dot_account_dot_service__pb2.GetAccountRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.GetAccountResponse.SerializeToString,
            ),
            'ListAccounts': grpc.unary_unary_rpc_method_handler(
                    servicer.ListAccounts,
                    request_deserializer=proto_dot_account_dot_service__pb2.ListAccountsRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.ListAccountsResponse.SerializeToString,
            ),
            'CreateAPIKey': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateAPIKey,
                    request_deserializer=proto_dot_account_dot_service__pb2.CreateAPIKeyRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.CreateAPIKeyResponse.SerializeToString,
            ),
            'ChangeAPIKeyName': grpc.unary_unary_rpc_method_handler(
                    servicer.ChangeAPIKeyName,
                    request_deserializer=proto_dot_account_dot_service__pb2.ChangeAPIKeyNameRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.ChangeAPIKeyNameResponse.SerializeToString,
            ),
            'EnableAPIKey': grpc.unary_unary_rpc_method_handler(
                    servicer.EnableAPIKey,
                    request_deserializer=proto_dot_account_dot_service__pb2.EnableAPIKeyRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.EnableAPIKeyResponse.SerializeToString,
            ),
            'DisableAPIKey': grpc.unary_unary_rpc_method_handler(
                    servicer.DisableAPIKey,
                    request_deserializer=proto_dot_account_dot_service__pb2.DisableAPIKeyRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.DisableAPIKeyResponse.SerializeToString,
            ),
            'GetAPIKey': grpc.unary_unary_rpc_method_handler(
                    servicer.GetAPIKey,
                    request_deserializer=proto_dot_account_dot_service__pb2.GetAPIKeyRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.GetAPIKeyResponse.SerializeToString,
            ),
            'ListAPIKeys': grpc.unary_unary_rpc_method_handler(
                    servicer.ListAPIKeys,
                    request_deserializer=proto_dot_account_dot_service__pb2.ListAPIKeysRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.ListAPIKeysResponse.SerializeToString,
            ),
            'GetAPIKeyBySearchingAllEnvironments': grpc.unary_unary_rpc_method_handler(
                    servicer.GetAPIKeyBySearchingAllEnvironments,
                    request_deserializer=proto_dot_account_dot_service__pb2.GetAPIKeyBySearchingAllEnvironmentsRequest.FromString,
                    response_serializer=proto_dot_account_dot_service__pb2.GetAPIKeyBySearchingAllEnvironmentsResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'bucketeer.account.AccountService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class AccountService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetMe(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/GetMe',
            proto_dot_account_dot_service__pb2.GetMeRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.GetMeResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetMeByEmail(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/GetMeByEmail',
            proto_dot_account_dot_service__pb2.GetMeByEmailRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.GetMeResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def CreateAdminAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/CreateAdminAccount',
            proto_dot_account_dot_service__pb2.CreateAdminAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.CreateAdminAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def EnableAdminAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/EnableAdminAccount',
            proto_dot_account_dot_service__pb2.EnableAdminAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.EnableAdminAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def DisableAdminAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/DisableAdminAccount',
            proto_dot_account_dot_service__pb2.DisableAdminAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.DisableAdminAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetAdminAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/GetAdminAccount',
            proto_dot_account_dot_service__pb2.GetAdminAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.GetAdminAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListAdminAccounts(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/ListAdminAccounts',
            proto_dot_account_dot_service__pb2.ListAdminAccountsRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.ListAdminAccountsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ConvertAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/ConvertAccount',
            proto_dot_account_dot_service__pb2.ConvertAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.ConvertAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def CreateAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/CreateAccount',
            proto_dot_account_dot_service__pb2.CreateAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.CreateAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def EnableAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/EnableAccount',
            proto_dot_account_dot_service__pb2.EnableAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.EnableAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def DisableAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/DisableAccount',
            proto_dot_account_dot_service__pb2.DisableAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.DisableAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ChangeAccountRole(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/ChangeAccountRole',
            proto_dot_account_dot_service__pb2.ChangeAccountRoleRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.ChangeAccountRoleResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetAccount(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/GetAccount',
            proto_dot_account_dot_service__pb2.GetAccountRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.GetAccountResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListAccounts(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/ListAccounts',
            proto_dot_account_dot_service__pb2.ListAccountsRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.ListAccountsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def CreateAPIKey(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/CreateAPIKey',
            proto_dot_account_dot_service__pb2.CreateAPIKeyRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.CreateAPIKeyResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ChangeAPIKeyName(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/ChangeAPIKeyName',
            proto_dot_account_dot_service__pb2.ChangeAPIKeyNameRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.ChangeAPIKeyNameResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def EnableAPIKey(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/EnableAPIKey',
            proto_dot_account_dot_service__pb2.EnableAPIKeyRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.EnableAPIKeyResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def DisableAPIKey(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/DisableAPIKey',
            proto_dot_account_dot_service__pb2.DisableAPIKeyRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.DisableAPIKeyResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetAPIKey(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/GetAPIKey',
            proto_dot_account_dot_service__pb2.GetAPIKeyRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.GetAPIKeyResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListAPIKeys(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/ListAPIKeys',
            proto_dot_account_dot_service__pb2.ListAPIKeysRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.ListAPIKeysResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetAPIKeyBySearchingAllEnvironments(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/bucketeer.account.AccountService/GetAPIKeyBySearchingAllEnvironments',
            proto_dot_account_dot_service__pb2.GetAPIKeyBySearchingAllEnvironmentsRequest.SerializeToString,
            proto_dot_account_dot_service__pb2.GetAPIKeyBySearchingAllEnvironmentsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
