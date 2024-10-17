from gogoproto import gogo_pb2 as _gogo_pb2
from tendermint.consensus import types_pb2 as _types_pb2
from tendermint.types import events_pb2 as _events_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgInfo(_message.Message):
    __slots__ = ("msg", "peer_id")
    MSG_FIELD_NUMBER: _ClassVar[int]
    PEER_ID_FIELD_NUMBER: _ClassVar[int]
    msg: _types_pb2.Message
    peer_id: str
    def __init__(self, msg: _Optional[_Union[_types_pb2.Message, _Mapping]] = ..., peer_id: _Optional[str] = ...) -> None: ...

class TimeoutInfo(_message.Message):
    __slots__ = ("duration", "height", "round", "step")
    DURATION_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    STEP_FIELD_NUMBER: _ClassVar[int]
    duration: _duration_pb2.Duration
    height: int
    round: int
    step: int
    def __init__(self, duration: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., height: _Optional[int] = ..., round: _Optional[int] = ..., step: _Optional[int] = ...) -> None: ...

class EndHeight(_message.Message):
    __slots__ = ("height",)
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    height: int
    def __init__(self, height: _Optional[int] = ...) -> None: ...

class WALMessage(_message.Message):
    __slots__ = ("event_data_round_state", "msg_info", "timeout_info", "end_height")
    EVENT_DATA_ROUND_STATE_FIELD_NUMBER: _ClassVar[int]
    MSG_INFO_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_INFO_FIELD_NUMBER: _ClassVar[int]
    END_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    event_data_round_state: _events_pb2.EventDataRoundState
    msg_info: MsgInfo
    timeout_info: TimeoutInfo
    end_height: EndHeight
    def __init__(self, event_data_round_state: _Optional[_Union[_events_pb2.EventDataRoundState, _Mapping]] = ..., msg_info: _Optional[_Union[MsgInfo, _Mapping]] = ..., timeout_info: _Optional[_Union[TimeoutInfo, _Mapping]] = ..., end_height: _Optional[_Union[EndHeight, _Mapping]] = ...) -> None: ...

class TimedWALMessage(_message.Message):
    __slots__ = ("time", "msg")
    TIME_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    time: _timestamp_pb2.Timestamp
    msg: WALMessage
    def __init__(self, time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., msg: _Optional[_Union[WALMessage, _Mapping]] = ...) -> None: ...
