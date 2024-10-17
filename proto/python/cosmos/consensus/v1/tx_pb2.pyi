from amino import amino_pb2 as _amino_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from tendermint.types import params_pb2 as _params_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "block", "evidence", "validator", "abci")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    BLOCK_FIELD_NUMBER: _ClassVar[int]
    EVIDENCE_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    ABCI_FIELD_NUMBER: _ClassVar[int]
    authority: str
    block: _params_pb2.BlockParams
    evidence: _params_pb2.EvidenceParams
    validator: _params_pb2.ValidatorParams
    abci: _params_pb2.ABCIParams
    def __init__(self, authority: _Optional[str] = ..., block: _Optional[_Union[_params_pb2.BlockParams, _Mapping]] = ..., evidence: _Optional[_Union[_params_pb2.EvidenceParams, _Mapping]] = ..., validator: _Optional[_Union[_params_pb2.ValidatorParams, _Mapping]] = ..., abci: _Optional[_Union[_params_pb2.ABCIParams, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
