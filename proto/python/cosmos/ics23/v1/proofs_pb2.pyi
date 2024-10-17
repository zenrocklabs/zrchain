from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class HashOp(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    NO_HASH: _ClassVar[HashOp]
    SHA256: _ClassVar[HashOp]
    SHA512: _ClassVar[HashOp]
    KECCAK256: _ClassVar[HashOp]
    RIPEMD160: _ClassVar[HashOp]
    BITCOIN: _ClassVar[HashOp]
    SHA512_256: _ClassVar[HashOp]
    BLAKE2B_512: _ClassVar[HashOp]
    BLAKE2S_256: _ClassVar[HashOp]
    BLAKE3: _ClassVar[HashOp]

class LengthOp(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    NO_PREFIX: _ClassVar[LengthOp]
    VAR_PROTO: _ClassVar[LengthOp]
    VAR_RLP: _ClassVar[LengthOp]
    FIXED32_BIG: _ClassVar[LengthOp]
    FIXED32_LITTLE: _ClassVar[LengthOp]
    FIXED64_BIG: _ClassVar[LengthOp]
    FIXED64_LITTLE: _ClassVar[LengthOp]
    REQUIRE_32_BYTES: _ClassVar[LengthOp]
    REQUIRE_64_BYTES: _ClassVar[LengthOp]
NO_HASH: HashOp
SHA256: HashOp
SHA512: HashOp
KECCAK256: HashOp
RIPEMD160: HashOp
BITCOIN: HashOp
SHA512_256: HashOp
BLAKE2B_512: HashOp
BLAKE2S_256: HashOp
BLAKE3: HashOp
NO_PREFIX: LengthOp
VAR_PROTO: LengthOp
VAR_RLP: LengthOp
FIXED32_BIG: LengthOp
FIXED32_LITTLE: LengthOp
FIXED64_BIG: LengthOp
FIXED64_LITTLE: LengthOp
REQUIRE_32_BYTES: LengthOp
REQUIRE_64_BYTES: LengthOp

class ExistenceProof(_message.Message):
    __slots__ = ("key", "value", "leaf", "path")
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    LEAF_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    key: bytes
    value: bytes
    leaf: LeafOp
    path: _containers.RepeatedCompositeFieldContainer[InnerOp]
    def __init__(self, key: _Optional[bytes] = ..., value: _Optional[bytes] = ..., leaf: _Optional[_Union[LeafOp, _Mapping]] = ..., path: _Optional[_Iterable[_Union[InnerOp, _Mapping]]] = ...) -> None: ...

class NonExistenceProof(_message.Message):
    __slots__ = ("key", "left", "right")
    KEY_FIELD_NUMBER: _ClassVar[int]
    LEFT_FIELD_NUMBER: _ClassVar[int]
    RIGHT_FIELD_NUMBER: _ClassVar[int]
    key: bytes
    left: ExistenceProof
    right: ExistenceProof
    def __init__(self, key: _Optional[bytes] = ..., left: _Optional[_Union[ExistenceProof, _Mapping]] = ..., right: _Optional[_Union[ExistenceProof, _Mapping]] = ...) -> None: ...

class CommitmentProof(_message.Message):
    __slots__ = ("exist", "nonexist", "batch", "compressed")
    EXIST_FIELD_NUMBER: _ClassVar[int]
    NONEXIST_FIELD_NUMBER: _ClassVar[int]
    BATCH_FIELD_NUMBER: _ClassVar[int]
    COMPRESSED_FIELD_NUMBER: _ClassVar[int]
    exist: ExistenceProof
    nonexist: NonExistenceProof
    batch: BatchProof
    compressed: CompressedBatchProof
    def __init__(self, exist: _Optional[_Union[ExistenceProof, _Mapping]] = ..., nonexist: _Optional[_Union[NonExistenceProof, _Mapping]] = ..., batch: _Optional[_Union[BatchProof, _Mapping]] = ..., compressed: _Optional[_Union[CompressedBatchProof, _Mapping]] = ...) -> None: ...

class LeafOp(_message.Message):
    __slots__ = ("hash", "prehash_key", "prehash_value", "length", "prefix")
    HASH_FIELD_NUMBER: _ClassVar[int]
    PREHASH_KEY_FIELD_NUMBER: _ClassVar[int]
    PREHASH_VALUE_FIELD_NUMBER: _ClassVar[int]
    LENGTH_FIELD_NUMBER: _ClassVar[int]
    PREFIX_FIELD_NUMBER: _ClassVar[int]
    hash: HashOp
    prehash_key: HashOp
    prehash_value: HashOp
    length: LengthOp
    prefix: bytes
    def __init__(self, hash: _Optional[_Union[HashOp, str]] = ..., prehash_key: _Optional[_Union[HashOp, str]] = ..., prehash_value: _Optional[_Union[HashOp, str]] = ..., length: _Optional[_Union[LengthOp, str]] = ..., prefix: _Optional[bytes] = ...) -> None: ...

class InnerOp(_message.Message):
    __slots__ = ("hash", "prefix", "suffix")
    HASH_FIELD_NUMBER: _ClassVar[int]
    PREFIX_FIELD_NUMBER: _ClassVar[int]
    SUFFIX_FIELD_NUMBER: _ClassVar[int]
    hash: HashOp
    prefix: bytes
    suffix: bytes
    def __init__(self, hash: _Optional[_Union[HashOp, str]] = ..., prefix: _Optional[bytes] = ..., suffix: _Optional[bytes] = ...) -> None: ...

class ProofSpec(_message.Message):
    __slots__ = ("leaf_spec", "inner_spec", "max_depth", "min_depth", "prehash_key_before_comparison")
    LEAF_SPEC_FIELD_NUMBER: _ClassVar[int]
    INNER_SPEC_FIELD_NUMBER: _ClassVar[int]
    MAX_DEPTH_FIELD_NUMBER: _ClassVar[int]
    MIN_DEPTH_FIELD_NUMBER: _ClassVar[int]
    PREHASH_KEY_BEFORE_COMPARISON_FIELD_NUMBER: _ClassVar[int]
    leaf_spec: LeafOp
    inner_spec: InnerSpec
    max_depth: int
    min_depth: int
    prehash_key_before_comparison: bool
    def __init__(self, leaf_spec: _Optional[_Union[LeafOp, _Mapping]] = ..., inner_spec: _Optional[_Union[InnerSpec, _Mapping]] = ..., max_depth: _Optional[int] = ..., min_depth: _Optional[int] = ..., prehash_key_before_comparison: bool = ...) -> None: ...

class InnerSpec(_message.Message):
    __slots__ = ("child_order", "child_size", "min_prefix_length", "max_prefix_length", "empty_child", "hash")
    CHILD_ORDER_FIELD_NUMBER: _ClassVar[int]
    CHILD_SIZE_FIELD_NUMBER: _ClassVar[int]
    MIN_PREFIX_LENGTH_FIELD_NUMBER: _ClassVar[int]
    MAX_PREFIX_LENGTH_FIELD_NUMBER: _ClassVar[int]
    EMPTY_CHILD_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    child_order: _containers.RepeatedScalarFieldContainer[int]
    child_size: int
    min_prefix_length: int
    max_prefix_length: int
    empty_child: bytes
    hash: HashOp
    def __init__(self, child_order: _Optional[_Iterable[int]] = ..., child_size: _Optional[int] = ..., min_prefix_length: _Optional[int] = ..., max_prefix_length: _Optional[int] = ..., empty_child: _Optional[bytes] = ..., hash: _Optional[_Union[HashOp, str]] = ...) -> None: ...

class BatchProof(_message.Message):
    __slots__ = ("entries",)
    ENTRIES_FIELD_NUMBER: _ClassVar[int]
    entries: _containers.RepeatedCompositeFieldContainer[BatchEntry]
    def __init__(self, entries: _Optional[_Iterable[_Union[BatchEntry, _Mapping]]] = ...) -> None: ...

class BatchEntry(_message.Message):
    __slots__ = ("exist", "nonexist")
    EXIST_FIELD_NUMBER: _ClassVar[int]
    NONEXIST_FIELD_NUMBER: _ClassVar[int]
    exist: ExistenceProof
    nonexist: NonExistenceProof
    def __init__(self, exist: _Optional[_Union[ExistenceProof, _Mapping]] = ..., nonexist: _Optional[_Union[NonExistenceProof, _Mapping]] = ...) -> None: ...

class CompressedBatchProof(_message.Message):
    __slots__ = ("entries", "lookup_inners")
    ENTRIES_FIELD_NUMBER: _ClassVar[int]
    LOOKUP_INNERS_FIELD_NUMBER: _ClassVar[int]
    entries: _containers.RepeatedCompositeFieldContainer[CompressedBatchEntry]
    lookup_inners: _containers.RepeatedCompositeFieldContainer[InnerOp]
    def __init__(self, entries: _Optional[_Iterable[_Union[CompressedBatchEntry, _Mapping]]] = ..., lookup_inners: _Optional[_Iterable[_Union[InnerOp, _Mapping]]] = ...) -> None: ...

class CompressedBatchEntry(_message.Message):
    __slots__ = ("exist", "nonexist")
    EXIST_FIELD_NUMBER: _ClassVar[int]
    NONEXIST_FIELD_NUMBER: _ClassVar[int]
    exist: CompressedExistenceProof
    nonexist: CompressedNonExistenceProof
    def __init__(self, exist: _Optional[_Union[CompressedExistenceProof, _Mapping]] = ..., nonexist: _Optional[_Union[CompressedNonExistenceProof, _Mapping]] = ...) -> None: ...

class CompressedExistenceProof(_message.Message):
    __slots__ = ("key", "value", "leaf", "path")
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    LEAF_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    key: bytes
    value: bytes
    leaf: LeafOp
    path: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, key: _Optional[bytes] = ..., value: _Optional[bytes] = ..., leaf: _Optional[_Union[LeafOp, _Mapping]] = ..., path: _Optional[_Iterable[int]] = ...) -> None: ...

class CompressedNonExistenceProof(_message.Message):
    __slots__ = ("key", "left", "right")
    KEY_FIELD_NUMBER: _ClassVar[int]
    LEFT_FIELD_NUMBER: _ClassVar[int]
    RIGHT_FIELD_NUMBER: _ClassVar[int]
    key: bytes
    left: CompressedExistenceProof
    right: CompressedExistenceProof
    def __init__(self, key: _Optional[bytes] = ..., left: _Optional[_Union[CompressedExistenceProof, _Mapping]] = ..., right: _Optional[_Union[CompressedExistenceProof, _Mapping]] = ...) -> None: ...
