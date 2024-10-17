from google.protobuf import any_pb2 as _any_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class AccountQueryRequest(_message.Message):
    __slots__ = ("target", "request")
    TARGET_FIELD_NUMBER: _ClassVar[int]
    REQUEST_FIELD_NUMBER: _ClassVar[int]
    target: str
    request: _any_pb2.Any
    def __init__(self, target: _Optional[str] = ..., request: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class AccountQueryResponse(_message.Message):
    __slots__ = ("response",)
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    response: _any_pb2.Any
    def __init__(self, response: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class SchemaRequest(_message.Message):
    __slots__ = ("account_type",)
    ACCOUNT_TYPE_FIELD_NUMBER: _ClassVar[int]
    account_type: str
    def __init__(self, account_type: _Optional[str] = ...) -> None: ...

class SchemaResponse(_message.Message):
    __slots__ = ("init_schema", "execute_handlers", "query_handlers")
    class Handler(_message.Message):
        __slots__ = ("request", "response")
        REQUEST_FIELD_NUMBER: _ClassVar[int]
        RESPONSE_FIELD_NUMBER: _ClassVar[int]
        request: str
        response: str
        def __init__(self, request: _Optional[str] = ..., response: _Optional[str] = ...) -> None: ...
    INIT_SCHEMA_FIELD_NUMBER: _ClassVar[int]
    EXECUTE_HANDLERS_FIELD_NUMBER: _ClassVar[int]
    QUERY_HANDLERS_FIELD_NUMBER: _ClassVar[int]
    init_schema: SchemaResponse.Handler
    execute_handlers: _containers.RepeatedCompositeFieldContainer[SchemaResponse.Handler]
    query_handlers: _containers.RepeatedCompositeFieldContainer[SchemaResponse.Handler]
    def __init__(self, init_schema: _Optional[_Union[SchemaResponse.Handler, _Mapping]] = ..., execute_handlers: _Optional[_Iterable[_Union[SchemaResponse.Handler, _Mapping]]] = ..., query_handlers: _Optional[_Iterable[_Union[SchemaResponse.Handler, _Mapping]]] = ...) -> None: ...

class AccountTypeRequest(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class AccountTypeResponse(_message.Message):
    __slots__ = ("account_type",)
    ACCOUNT_TYPE_FIELD_NUMBER: _ClassVar[int]
    account_type: str
    def __init__(self, account_type: _Optional[str] = ...) -> None: ...

class AccountNumberRequest(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class AccountNumberResponse(_message.Message):
    __slots__ = ("number",)
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    number: int
    def __init__(self, number: _Optional[int] = ...) -> None: ...
