# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: services/verifiable-credentials/templates/v1/templates.proto
# plugin: python-betterproto
from dataclasses import dataclass
from typing import Dict, List

import betterproto
from betterproto.grpc.grpclib_server import ServiceBase
import grpclib


@dataclass(eq=False, repr=False)
class CreateCredentialTemplateRequest(betterproto.Message):
    name: str = betterproto.string_field(1)
    schema: "___common_v1__.JsonPayload" = betterproto.message_field(2)
    base_uri: str = betterproto.string_field(3)


@dataclass(eq=False, repr=False)
class CreateCredentialTemplateResponse(betterproto.Message):
    id: str = betterproto.string_field(1)
    uri: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class GetCredentialTemplateRequest(betterproto.Message):
    id: str = betterproto.string_field(1)


@dataclass(eq=False, repr=False)
class GetCredentialTemplateResponse(betterproto.Message):
    template: "CredentialTemplate" = betterproto.message_field(1)


@dataclass(eq=False, repr=False)
class SearchCredentialTemplatesRequest(betterproto.Message):
    query: str = betterproto.string_field(1)
    continuation_token: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class SearchCredentialTemplatesResponse(betterproto.Message):
    templates: List["CredentialTemplate"] = betterproto.message_field(1)
    has_more: bool = betterproto.bool_field(2)
    count: int = betterproto.int32_field(3)
    continuation_token: str = betterproto.string_field(4)


@dataclass(eq=False, repr=False)
class UpdateCredentialTemplateRequest(betterproto.Message):
    id: str = betterproto.string_field(1)
    name: str = betterproto.string_field(2)
    schema: "___common_v1__.JsonPayload" = betterproto.message_field(3)


@dataclass(eq=False, repr=False)
class UpdateCredentialTemplateResponse(betterproto.Message):
    template: "CredentialTemplate" = betterproto.message_field(1)


@dataclass(eq=False, repr=False)
class DeleteCredentialTemplateRequest(betterproto.Message):
    id: str = betterproto.string_field(1)


@dataclass(eq=False, repr=False)
class DeleteCredentialTemplateResponse(betterproto.Message):
    pass


@dataclass(eq=False, repr=False)
class CredentialTemplate(betterproto.Message):
    id: str = betterproto.string_field(1)
    name: str = betterproto.string_field(2)
    version: str = betterproto.string_field(3)
    schema: "___common_v1__.JsonPayload" = betterproto.message_field(4)
    uri: str = betterproto.string_field(5)


class CredentialTemplatesStub(betterproto.ServiceStub):
    async def create(
        self,
        *,
        name: str = "",
        schema: "___common_v1__.JsonPayload" = None,
        base_uri: str = "",
    ) -> "CreateCredentialTemplateResponse":

        request = CreateCredentialTemplateRequest()
        request.name = name
        if schema is not None:
            request.schema = schema
        request.base_uri = base_uri

        return await self._unary_unary(
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Create",
            request,
            CreateCredentialTemplateResponse,
        )

    async def get(self, *, id: str = "") -> "GetCredentialTemplateResponse":

        request = GetCredentialTemplateRequest()
        request.id = id

        return await self._unary_unary(
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Get",
            request,
            GetCredentialTemplateResponse,
        )

    async def search(
        self, *, query: str = "", continuation_token: str = ""
    ) -> "SearchCredentialTemplatesResponse":

        request = SearchCredentialTemplatesRequest()
        request.query = query
        request.continuation_token = continuation_token

        return await self._unary_unary(
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Search",
            request,
            SearchCredentialTemplatesResponse,
        )

    async def update(
        self,
        *,
        id: str = "",
        name: str = "",
        schema: "___common_v1__.JsonPayload" = None,
    ) -> "UpdateCredentialTemplateResponse":

        request = UpdateCredentialTemplateRequest()
        request.id = id
        request.name = name
        if schema is not None:
            request.schema = schema

        return await self._unary_unary(
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Update",
            request,
            UpdateCredentialTemplateResponse,
        )

    async def delete(self, *, id: str = "") -> "DeleteCredentialTemplateResponse":

        request = DeleteCredentialTemplateRequest()
        request.id = id

        return await self._unary_unary(
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Delete",
            request,
            DeleteCredentialTemplateResponse,
        )


class CredentialTemplatesBase(ServiceBase):
    async def create(
        self, name: str, schema: "___common_v1__.JsonPayload", base_uri: str
    ) -> "CreateCredentialTemplateResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def get(self, id: str) -> "GetCredentialTemplateResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def search(
        self, query: str, continuation_token: str
    ) -> "SearchCredentialTemplatesResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def update(
        self, id: str, name: str, schema: "___common_v1__.JsonPayload"
    ) -> "UpdateCredentialTemplateResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def delete(self, id: str) -> "DeleteCredentialTemplateResponse":
        raise grpclib.GRPCError(grpclib.const.Status.UNIMPLEMENTED)

    async def __rpc_create(self, stream: grpclib.server.Stream) -> None:
        request = await stream.recv_message()

        request_kwargs = {
            "name": request.name,
            "schema": request.schema,
            "base_uri": request.base_uri,
        }

        response = await self.create(**request_kwargs)
        await stream.send_message(response)

    async def __rpc_get(self, stream: grpclib.server.Stream) -> None:
        request = await stream.recv_message()

        request_kwargs = {
            "id": request.id,
        }

        response = await self.get(**request_kwargs)
        await stream.send_message(response)

    async def __rpc_search(self, stream: grpclib.server.Stream) -> None:
        request = await stream.recv_message()

        request_kwargs = {
            "query": request.query,
            "continuation_token": request.continuation_token,
        }

        response = await self.search(**request_kwargs)
        await stream.send_message(response)

    async def __rpc_update(self, stream: grpclib.server.Stream) -> None:
        request = await stream.recv_message()

        request_kwargs = {
            "id": request.id,
            "name": request.name,
            "schema": request.schema,
        }

        response = await self.update(**request_kwargs)
        await stream.send_message(response)

    async def __rpc_delete(self, stream: grpclib.server.Stream) -> None:
        request = await stream.recv_message()

        request_kwargs = {
            "id": request.id,
        }

        response = await self.delete(**request_kwargs)
        await stream.send_message(response)

    def __mapping__(self) -> Dict[str, grpclib.const.Handler]:
        return {
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Create": grpclib.const.Handler(
                self.__rpc_create,
                grpclib.const.Cardinality.UNARY_UNARY,
                CreateCredentialTemplateRequest,
                CreateCredentialTemplateResponse,
            ),
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Get": grpclib.const.Handler(
                self.__rpc_get,
                grpclib.const.Cardinality.UNARY_UNARY,
                GetCredentialTemplateRequest,
                GetCredentialTemplateResponse,
            ),
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Search": grpclib.const.Handler(
                self.__rpc_search,
                grpclib.const.Cardinality.UNARY_UNARY,
                SearchCredentialTemplatesRequest,
                SearchCredentialTemplatesResponse,
            ),
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Update": grpclib.const.Handler(
                self.__rpc_update,
                grpclib.const.Cardinality.UNARY_UNARY,
                UpdateCredentialTemplateRequest,
                UpdateCredentialTemplateResponse,
            ),
            "/services.verifiablecredentials.templates.v1.CredentialTemplates/Delete": grpclib.const.Handler(
                self.__rpc_delete,
                grpclib.const.Cardinality.UNARY_UNARY,
                DeleteCredentialTemplateRequest,
                DeleteCredentialTemplateResponse,
            ),
        }


from ....common import v1 as ___common_v1__
