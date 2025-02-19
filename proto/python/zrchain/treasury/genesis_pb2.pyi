from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.treasury import params_pb2 as _params_pb2
from zrchain.treasury import key_pb2 as _key_pb2
from zrchain.treasury import mpcsign_pb2 as _mpcsign_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("params", "port_id", "keys", "key_requests", "sign_requests", "sign_tx_requests", "ica_tx_requests", "no_fee_msgs")
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    KEYS_FIELD_NUMBER: _ClassVar[int]
    KEY_REQUESTS_FIELD_NUMBER: _ClassVar[int]
    SIGN_REQUESTS_FIELD_NUMBER: _ClassVar[int]
    SIGN_TX_REQUESTS_FIELD_NUMBER: _ClassVar[int]
    ICA_TX_REQUESTS_FIELD_NUMBER: _ClassVar[int]
    NO_FEE_MSGS_FIELD_NUMBER: _ClassVar[int]
    params: _params_pb2.Params
    port_id: str
    keys: _containers.RepeatedCompositeFieldContainer[_key_pb2.Key]
    key_requests: _containers.RepeatedCompositeFieldContainer[_key_pb2.KeyRequest]
    sign_requests: _containers.RepeatedCompositeFieldContainer[_mpcsign_pb2.SignRequest]
    sign_tx_requests: _containers.RepeatedCompositeFieldContainer[_mpcsign_pb2.SignTransactionRequest]
    ica_tx_requests: _containers.RepeatedCompositeFieldContainer[_mpcsign_pb2.ICATransactionRequest]
    no_fee_msgs: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ..., port_id: _Optional[str] = ..., keys: _Optional[_Iterable[_Union[_key_pb2.Key, _Mapping]]] = ..., key_requests: _Optional[_Iterable[_Union[_key_pb2.KeyRequest, _Mapping]]] = ..., sign_requests: _Optional[_Iterable[_Union[_mpcsign_pb2.SignRequest, _Mapping]]] = ..., sign_tx_requests: _Optional[_Iterable[_Union[_mpcsign_pb2.SignTransactionRequest, _Mapping]]] = ..., ica_tx_requests: _Optional[_Iterable[_Union[_mpcsign_pb2.ICATransactionRequest, _Mapping]]] = ..., no_fee_msgs: _Optional[_Iterable[str]] = ...) -> None: ...
