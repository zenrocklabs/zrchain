from cosmos.app.v1alpha1 import module_pb2 as _module_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Module(_message.Message):
    __slots__ = ("genesis_params",)
    GENESIS_PARAMS_FIELD_NUMBER: _ClassVar[int]
    genesis_params: GeneratorParams
    def __init__(self, genesis_params: _Optional[_Union[GeneratorParams, _Mapping]] = ...) -> None: ...

class GeneratorParams(_message.Message):
    __slots__ = ("seed", "bucket_count", "key_mean", "key_std_dev", "value_mean", "value_std_dev", "genesis_count", "insert_weight", "update_weight", "get_weight", "delete_weight")
    SEED_FIELD_NUMBER: _ClassVar[int]
    BUCKET_COUNT_FIELD_NUMBER: _ClassVar[int]
    KEY_MEAN_FIELD_NUMBER: _ClassVar[int]
    KEY_STD_DEV_FIELD_NUMBER: _ClassVar[int]
    VALUE_MEAN_FIELD_NUMBER: _ClassVar[int]
    VALUE_STD_DEV_FIELD_NUMBER: _ClassVar[int]
    GENESIS_COUNT_FIELD_NUMBER: _ClassVar[int]
    INSERT_WEIGHT_FIELD_NUMBER: _ClassVar[int]
    UPDATE_WEIGHT_FIELD_NUMBER: _ClassVar[int]
    GET_WEIGHT_FIELD_NUMBER: _ClassVar[int]
    DELETE_WEIGHT_FIELD_NUMBER: _ClassVar[int]
    seed: int
    bucket_count: int
    key_mean: int
    key_std_dev: int
    value_mean: int
    value_std_dev: int
    genesis_count: int
    insert_weight: float
    update_weight: float
    get_weight: float
    delete_weight: float
    def __init__(self, seed: _Optional[int] = ..., bucket_count: _Optional[int] = ..., key_mean: _Optional[int] = ..., key_std_dev: _Optional[int] = ..., value_mean: _Optional[int] = ..., value_std_dev: _Optional[int] = ..., genesis_count: _Optional[int] = ..., insert_weight: _Optional[float] = ..., update_weight: _Optional[float] = ..., get_weight: _Optional[float] = ..., delete_weight: _Optional[float] = ...) -> None: ...
