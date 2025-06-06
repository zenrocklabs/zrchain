# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: tendermint/types/types.proto
# Protobuf Python Version: 6.31.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    6,
    31,
    1,
    '',
    'tendermint/types/types.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from tendermint.crypto import proof_pb2 as tendermint_dot_crypto_dot_proof__pb2
from tendermint.version import types_pb2 as tendermint_dot_version_dot_types__pb2
from tendermint.types import validator_pb2 as tendermint_dot_types_dot_validator__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1ctendermint/types/types.proto\x12\x10tendermint.types\x1a\x14gogoproto/gogo.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1dtendermint/crypto/proof.proto\x1a\x1etendermint/version/types.proto\x1a tendermint/types/validator.proto\"9\n\rPartSetHeader\x12\x14\n\x05total\x18\x01 \x01(\rR\x05total\x12\x12\n\x04hash\x18\x02 \x01(\x0cR\x04hash\"h\n\x04Part\x12\x14\n\x05index\x18\x01 \x01(\rR\x05index\x12\x14\n\x05\x62ytes\x18\x02 \x01(\x0cR\x05\x62ytes\x12\x34\n\x05proof\x18\x03 \x01(\x0b\x32\x18.tendermint.crypto.ProofB\x04\xc8\xde\x1f\x00R\x05proof\"l\n\x07\x42lockID\x12\x12\n\x04hash\x18\x01 \x01(\x0cR\x04hash\x12M\n\x0fpart_set_header\x18\x02 \x01(\x0b\x32\x1f.tendermint.types.PartSetHeaderB\x04\xc8\xde\x1f\x00R\rpartSetHeader\"\xe6\x04\n\x06Header\x12=\n\x07version\x18\x01 \x01(\x0b\x32\x1d.tendermint.version.ConsensusB\x04\xc8\xde\x1f\x00R\x07version\x12&\n\x08\x63hain_id\x18\x02 \x01(\tB\x0b\xe2\xde\x1f\x07\x43hainIDR\x07\x63hainId\x12\x16\n\x06height\x18\x03 \x01(\x03R\x06height\x12\x38\n\x04time\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x04time\x12\x43\n\rlast_block_id\x18\x05 \x01(\x0b\x32\x19.tendermint.types.BlockIDB\x04\xc8\xde\x1f\x00R\x0blastBlockId\x12(\n\x10last_commit_hash\x18\x06 \x01(\x0cR\x0elastCommitHash\x12\x1b\n\tdata_hash\x18\x07 \x01(\x0cR\x08\x64\x61taHash\x12\'\n\x0fvalidators_hash\x18\x08 \x01(\x0cR\x0evalidatorsHash\x12\x30\n\x14next_validators_hash\x18\t \x01(\x0cR\x12nextValidatorsHash\x12%\n\x0e\x63onsensus_hash\x18\n \x01(\x0cR\rconsensusHash\x12\x19\n\x08\x61pp_hash\x18\x0b \x01(\x0cR\x07\x61ppHash\x12*\n\x11last_results_hash\x18\x0c \x01(\x0cR\x0flastResultsHash\x12#\n\revidence_hash\x18\r \x01(\x0cR\x0c\x65videnceHash\x12)\n\x10proposer_address\x18\x0e \x01(\x0cR\x0fproposerAddress\"\x18\n\x04\x44\x61ta\x12\x10\n\x03txs\x18\x01 \x03(\x0cR\x03txs\"\xb7\x03\n\x04Vote\x12\x33\n\x04type\x18\x01 \x01(\x0e\x32\x1f.tendermint.types.SignedMsgTypeR\x04type\x12\x16\n\x06height\x18\x02 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x03 \x01(\x05R\x05round\x12\x45\n\x08\x62lock_id\x18\x04 \x01(\x0b\x32\x19.tendermint.types.BlockIDB\x0f\xc8\xde\x1f\x00\xe2\xde\x1f\x07\x42lockIDR\x07\x62lockId\x12\x42\n\ttimestamp\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\ttimestamp\x12+\n\x11validator_address\x18\x06 \x01(\x0cR\x10validatorAddress\x12\'\n\x0fvalidator_index\x18\x07 \x01(\x05R\x0evalidatorIndex\x12\x1c\n\tsignature\x18\x08 \x01(\x0cR\tsignature\x12\x1c\n\textension\x18\t \x01(\x0cR\textension\x12/\n\x13\x65xtension_signature\x18\n \x01(\x0cR\x12\x65xtensionSignature\"\xc0\x01\n\x06\x43ommit\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12\x45\n\x08\x62lock_id\x18\x03 \x01(\x0b\x32\x19.tendermint.types.BlockIDB\x0f\xc8\xde\x1f\x00\xe2\xde\x1f\x07\x42lockIDR\x07\x62lockId\x12\x41\n\nsignatures\x18\x04 \x03(\x0b\x32\x1b.tendermint.types.CommitSigB\x04\xc8\xde\x1f\x00R\nsignatures\"\xdd\x01\n\tCommitSig\x12\x41\n\rblock_id_flag\x18\x01 \x01(\x0e\x32\x1d.tendermint.types.BlockIDFlagR\x0b\x62lockIdFlag\x12+\n\x11validator_address\x18\x02 \x01(\x0cR\x10validatorAddress\x12\x42\n\ttimestamp\x18\x03 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\ttimestamp\x12\x1c\n\tsignature\x18\x04 \x01(\x0cR\tsignature\"\xe1\x01\n\x0e\x45xtendedCommit\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12\x45\n\x08\x62lock_id\x18\x03 \x01(\x0b\x32\x19.tendermint.types.BlockIDB\x0f\xc8\xde\x1f\x00\xe2\xde\x1f\x07\x42lockIDR\x07\x62lockId\x12Z\n\x13\x65xtended_signatures\x18\x04 \x03(\x0b\x32#.tendermint.types.ExtendedCommitSigB\x04\xc8\xde\x1f\x00R\x12\x65xtendedSignatures\"\xb4\x02\n\x11\x45xtendedCommitSig\x12\x41\n\rblock_id_flag\x18\x01 \x01(\x0e\x32\x1d.tendermint.types.BlockIDFlagR\x0b\x62lockIdFlag\x12+\n\x11validator_address\x18\x02 \x01(\x0cR\x10validatorAddress\x12\x42\n\ttimestamp\x18\x03 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\ttimestamp\x12\x1c\n\tsignature\x18\x04 \x01(\x0cR\tsignature\x12\x1c\n\textension\x18\x05 \x01(\x0cR\textension\x12/\n\x13\x65xtension_signature\x18\x06 \x01(\x0cR\x12\x65xtensionSignature\"\xb3\x02\n\x08Proposal\x12\x33\n\x04type\x18\x01 \x01(\x0e\x32\x1f.tendermint.types.SignedMsgTypeR\x04type\x12\x16\n\x06height\x18\x02 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x03 \x01(\x05R\x05round\x12\x1b\n\tpol_round\x18\x04 \x01(\x05R\x08polRound\x12\x45\n\x08\x62lock_id\x18\x05 \x01(\x0b\x32\x19.tendermint.types.BlockIDB\x0f\xc8\xde\x1f\x00\xe2\xde\x1f\x07\x42lockIDR\x07\x62lockId\x12\x42\n\ttimestamp\x18\x06 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\ttimestamp\x12\x1c\n\tsignature\x18\x07 \x01(\x0cR\tsignature\"r\n\x0cSignedHeader\x12\x30\n\x06header\x18\x01 \x01(\x0b\x32\x18.tendermint.types.HeaderR\x06header\x12\x30\n\x06\x63ommit\x18\x02 \x01(\x0b\x32\x18.tendermint.types.CommitR\x06\x63ommit\"\x96\x01\n\nLightBlock\x12\x43\n\rsigned_header\x18\x01 \x01(\x0b\x32\x1e.tendermint.types.SignedHeaderR\x0csignedHeader\x12\x43\n\rvalidator_set\x18\x02 \x01(\x0b\x32\x1e.tendermint.types.ValidatorSetR\x0cvalidatorSet\"\xc2\x01\n\tBlockMeta\x12\x45\n\x08\x62lock_id\x18\x01 \x01(\x0b\x32\x19.tendermint.types.BlockIDB\x0f\xc8\xde\x1f\x00\xe2\xde\x1f\x07\x42lockIDR\x07\x62lockId\x12\x1d\n\nblock_size\x18\x02 \x01(\x03R\tblockSize\x12\x36\n\x06header\x18\x03 \x01(\x0b\x32\x18.tendermint.types.HeaderB\x04\xc8\xde\x1f\x00R\x06header\x12\x17\n\x07num_txs\x18\x04 \x01(\x03R\x06numTxs\"j\n\x07TxProof\x12\x1b\n\troot_hash\x18\x01 \x01(\x0cR\x08rootHash\x12\x12\n\x04\x64\x61ta\x18\x02 \x01(\x0cR\x04\x64\x61ta\x12.\n\x05proof\x18\x03 \x01(\x0b\x32\x18.tendermint.crypto.ProofR\x05proof*\xd7\x01\n\rSignedMsgType\x12,\n\x17SIGNED_MSG_TYPE_UNKNOWN\x10\x00\x1a\x0f\x8a\x9d \x0bUnknownType\x12,\n\x17SIGNED_MSG_TYPE_PREVOTE\x10\x01\x1a\x0f\x8a\x9d \x0bPrevoteType\x12\x30\n\x19SIGNED_MSG_TYPE_PRECOMMIT\x10\x02\x1a\x11\x8a\x9d \rPrecommitType\x12.\n\x18SIGNED_MSG_TYPE_PROPOSAL\x10 \x1a\x10\x8a\x9d \x0cProposalType\x1a\x08\x88\xa3\x1e\x00\xa8\xa4\x1e\x01\x42\x35Z3github.com/cometbft/cometbft/proto/tendermint/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'tendermint.types.types_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z3github.com/cometbft/cometbft/proto/tendermint/types'
  _globals['_SIGNEDMSGTYPE']._loaded_options = None
  _globals['_SIGNEDMSGTYPE']._serialized_options = b'\210\243\036\000\250\244\036\001'
  _globals['_SIGNEDMSGTYPE'].values_by_name["SIGNED_MSG_TYPE_UNKNOWN"]._loaded_options = None
  _globals['_SIGNEDMSGTYPE'].values_by_name["SIGNED_MSG_TYPE_UNKNOWN"]._serialized_options = b'\212\235 \013UnknownType'
  _globals['_SIGNEDMSGTYPE'].values_by_name["SIGNED_MSG_TYPE_PREVOTE"]._loaded_options = None
  _globals['_SIGNEDMSGTYPE'].values_by_name["SIGNED_MSG_TYPE_PREVOTE"]._serialized_options = b'\212\235 \013PrevoteType'
  _globals['_SIGNEDMSGTYPE'].values_by_name["SIGNED_MSG_TYPE_PRECOMMIT"]._loaded_options = None
  _globals['_SIGNEDMSGTYPE'].values_by_name["SIGNED_MSG_TYPE_PRECOMMIT"]._serialized_options = b'\212\235 \rPrecommitType'
  _globals['_SIGNEDMSGTYPE'].values_by_name["SIGNED_MSG_TYPE_PROPOSAL"]._loaded_options = None
  _globals['_SIGNEDMSGTYPE'].values_by_name["SIGNED_MSG_TYPE_PROPOSAL"]._serialized_options = b'\212\235 \014ProposalType'
  _globals['_PART'].fields_by_name['proof']._loaded_options = None
  _globals['_PART'].fields_by_name['proof']._serialized_options = b'\310\336\037\000'
  _globals['_BLOCKID'].fields_by_name['part_set_header']._loaded_options = None
  _globals['_BLOCKID'].fields_by_name['part_set_header']._serialized_options = b'\310\336\037\000'
  _globals['_HEADER'].fields_by_name['version']._loaded_options = None
  _globals['_HEADER'].fields_by_name['version']._serialized_options = b'\310\336\037\000'
  _globals['_HEADER'].fields_by_name['chain_id']._loaded_options = None
  _globals['_HEADER'].fields_by_name['chain_id']._serialized_options = b'\342\336\037\007ChainID'
  _globals['_HEADER'].fields_by_name['time']._loaded_options = None
  _globals['_HEADER'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_HEADER'].fields_by_name['last_block_id']._loaded_options = None
  _globals['_HEADER'].fields_by_name['last_block_id']._serialized_options = b'\310\336\037\000'
  _globals['_VOTE'].fields_by_name['block_id']._loaded_options = None
  _globals['_VOTE'].fields_by_name['block_id']._serialized_options = b'\310\336\037\000\342\336\037\007BlockID'
  _globals['_VOTE'].fields_by_name['timestamp']._loaded_options = None
  _globals['_VOTE'].fields_by_name['timestamp']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_COMMIT'].fields_by_name['block_id']._loaded_options = None
  _globals['_COMMIT'].fields_by_name['block_id']._serialized_options = b'\310\336\037\000\342\336\037\007BlockID'
  _globals['_COMMIT'].fields_by_name['signatures']._loaded_options = None
  _globals['_COMMIT'].fields_by_name['signatures']._serialized_options = b'\310\336\037\000'
  _globals['_COMMITSIG'].fields_by_name['timestamp']._loaded_options = None
  _globals['_COMMITSIG'].fields_by_name['timestamp']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_EXTENDEDCOMMIT'].fields_by_name['block_id']._loaded_options = None
  _globals['_EXTENDEDCOMMIT'].fields_by_name['block_id']._serialized_options = b'\310\336\037\000\342\336\037\007BlockID'
  _globals['_EXTENDEDCOMMIT'].fields_by_name['extended_signatures']._loaded_options = None
  _globals['_EXTENDEDCOMMIT'].fields_by_name['extended_signatures']._serialized_options = b'\310\336\037\000'
  _globals['_EXTENDEDCOMMITSIG'].fields_by_name['timestamp']._loaded_options = None
  _globals['_EXTENDEDCOMMITSIG'].fields_by_name['timestamp']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_PROPOSAL'].fields_by_name['block_id']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['block_id']._serialized_options = b'\310\336\037\000\342\336\037\007BlockID'
  _globals['_PROPOSAL'].fields_by_name['timestamp']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['timestamp']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_BLOCKMETA'].fields_by_name['block_id']._loaded_options = None
  _globals['_BLOCKMETA'].fields_by_name['block_id']._serialized_options = b'\310\336\037\000\342\336\037\007BlockID'
  _globals['_BLOCKMETA'].fields_by_name['header']._loaded_options = None
  _globals['_BLOCKMETA'].fields_by_name['header']._serialized_options = b'\310\336\037\000'
  _globals['_SIGNEDMSGTYPE']._serialized_start=3405
  _globals['_SIGNEDMSGTYPE']._serialized_end=3620
  _globals['_PARTSETHEADER']._serialized_start=202
  _globals['_PARTSETHEADER']._serialized_end=259
  _globals['_PART']._serialized_start=261
  _globals['_PART']._serialized_end=365
  _globals['_BLOCKID']._serialized_start=367
  _globals['_BLOCKID']._serialized_end=475
  _globals['_HEADER']._serialized_start=478
  _globals['_HEADER']._serialized_end=1092
  _globals['_DATA']._serialized_start=1094
  _globals['_DATA']._serialized_end=1118
  _globals['_VOTE']._serialized_start=1121
  _globals['_VOTE']._serialized_end=1560
  _globals['_COMMIT']._serialized_start=1563
  _globals['_COMMIT']._serialized_end=1755
  _globals['_COMMITSIG']._serialized_start=1758
  _globals['_COMMITSIG']._serialized_end=1979
  _globals['_EXTENDEDCOMMIT']._serialized_start=1982
  _globals['_EXTENDEDCOMMIT']._serialized_end=2207
  _globals['_EXTENDEDCOMMITSIG']._serialized_start=2210
  _globals['_EXTENDEDCOMMITSIG']._serialized_end=2518
  _globals['_PROPOSAL']._serialized_start=2521
  _globals['_PROPOSAL']._serialized_end=2828
  _globals['_SIGNEDHEADER']._serialized_start=2830
  _globals['_SIGNEDHEADER']._serialized_end=2944
  _globals['_LIGHTBLOCK']._serialized_start=2947
  _globals['_LIGHTBLOCK']._serialized_end=3097
  _globals['_BLOCKMETA']._serialized_start=3100
  _globals['_BLOCKMETA']._serialized_end=3294
  _globals['_TXPROOF']._serialized_start=3296
  _globals['_TXPROOF']._serialized_end=3402
# @@protoc_insertion_point(module_scope)
