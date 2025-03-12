# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: tendermint/abci/types.proto
# Protobuf Python Version: 6.30.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    6,
    30,
    0,
    '',
    'tendermint/abci/types.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from tendermint.crypto import proof_pb2 as tendermint_dot_crypto_dot_proof__pb2
from tendermint.crypto import keys_pb2 as tendermint_dot_crypto_dot_keys__pb2
from tendermint.types import params_pb2 as tendermint_dot_types_dot_params__pb2
from tendermint.types import validator_pb2 as tendermint_dot_types_dot_validator__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1btendermint/abci/types.proto\x12\x0ftendermint.abci\x1a\x1dtendermint/crypto/proof.proto\x1a\x1ctendermint/crypto/keys.proto\x1a\x1dtendermint/types/params.proto\x1a tendermint/types/validator.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x14gogoproto/gogo.proto\"\xbf\t\n\x07Request\x12\x32\n\x04\x65\x63ho\x18\x01 \x01(\x0b\x32\x1c.tendermint.abci.RequestEchoH\x00R\x04\x65\x63ho\x12\x35\n\x05\x66lush\x18\x02 \x01(\x0b\x32\x1d.tendermint.abci.RequestFlushH\x00R\x05\x66lush\x12\x32\n\x04info\x18\x03 \x01(\x0b\x32\x1c.tendermint.abci.RequestInfoH\x00R\x04info\x12\x42\n\ninit_chain\x18\x05 \x01(\x0b\x32!.tendermint.abci.RequestInitChainH\x00R\tinitChain\x12\x35\n\x05query\x18\x06 \x01(\x0b\x32\x1d.tendermint.abci.RequestQueryH\x00R\x05query\x12<\n\x08\x63heck_tx\x18\x08 \x01(\x0b\x32\x1f.tendermint.abci.RequestCheckTxH\x00R\x07\x63heckTx\x12\x38\n\x06\x63ommit\x18\x0b \x01(\x0b\x32\x1e.tendermint.abci.RequestCommitH\x00R\x06\x63ommit\x12N\n\x0elist_snapshots\x18\x0c \x01(\x0b\x32%.tendermint.abci.RequestListSnapshotsH\x00R\rlistSnapshots\x12N\n\x0eoffer_snapshot\x18\r \x01(\x0b\x32%.tendermint.abci.RequestOfferSnapshotH\x00R\rofferSnapshot\x12[\n\x13load_snapshot_chunk\x18\x0e \x01(\x0b\x32).tendermint.abci.RequestLoadSnapshotChunkH\x00R\x11loadSnapshotChunk\x12^\n\x14\x61pply_snapshot_chunk\x18\x0f \x01(\x0b\x32*.tendermint.abci.RequestApplySnapshotChunkH\x00R\x12\x61pplySnapshotChunk\x12T\n\x10prepare_proposal\x18\x10 \x01(\x0b\x32\'.tendermint.abci.RequestPrepareProposalH\x00R\x0fprepareProposal\x12T\n\x10process_proposal\x18\x11 \x01(\x0b\x32\'.tendermint.abci.RequestProcessProposalH\x00R\x0fprocessProposal\x12\x45\n\x0b\x65xtend_vote\x18\x12 \x01(\x0b\x32\".tendermint.abci.RequestExtendVoteH\x00R\nextendVote\x12\x61\n\x15verify_vote_extension\x18\x13 \x01(\x0b\x32+.tendermint.abci.RequestVerifyVoteExtensionH\x00R\x13verifyVoteExtension\x12N\n\x0e\x66inalize_block\x18\x14 \x01(\x0b\x32%.tendermint.abci.RequestFinalizeBlockH\x00R\rfinalizeBlockB\x07\n\x05valueJ\x04\x08\x04\x10\x05J\x04\x08\x07\x10\x08J\x04\x08\t\x10\nJ\x04\x08\n\x10\x0b\"\'\n\x0bRequestEcho\x12\x18\n\x07message\x18\x01 \x01(\tR\x07message\"\x0e\n\x0cRequestFlush\"\x90\x01\n\x0bRequestInfo\x12\x18\n\x07version\x18\x01 \x01(\tR\x07version\x12#\n\rblock_version\x18\x02 \x01(\x04R\x0c\x62lockVersion\x12\x1f\n\x0bp2p_version\x18\x03 \x01(\x04R\np2pVersion\x12!\n\x0c\x61\x62\x63i_version\x18\x04 \x01(\tR\x0b\x61\x62\x63iVersion\"\xcc\x02\n\x10RequestInitChain\x12\x38\n\x04time\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x04time\x12\x19\n\x08\x63hain_id\x18\x02 \x01(\tR\x07\x63hainId\x12L\n\x10\x63onsensus_params\x18\x03 \x01(\x0b\x32!.tendermint.types.ConsensusParamsR\x0f\x63onsensusParams\x12\x46\n\nvalidators\x18\x04 \x03(\x0b\x32 .tendermint.abci.ValidatorUpdateB\x04\xc8\xde\x1f\x00R\nvalidators\x12&\n\x0f\x61pp_state_bytes\x18\x05 \x01(\x0cR\rappStateBytes\x12%\n\x0einitial_height\x18\x06 \x01(\x03R\rinitialHeight\"d\n\x0cRequestQuery\x12\x12\n\x04\x64\x61ta\x18\x01 \x01(\x0cR\x04\x64\x61ta\x12\x12\n\x04path\x18\x02 \x01(\tR\x04path\x12\x16\n\x06height\x18\x03 \x01(\x03R\x06height\x12\x14\n\x05prove\x18\x04 \x01(\x08R\x05prove\"R\n\x0eRequestCheckTx\x12\x0e\n\x02tx\x18\x01 \x01(\x0cR\x02tx\x12\x30\n\x04type\x18\x02 \x01(\x0e\x32\x1c.tendermint.abci.CheckTxTypeR\x04type\"\x0f\n\rRequestCommit\"\x16\n\x14RequestListSnapshots\"h\n\x14RequestOfferSnapshot\x12\x35\n\x08snapshot\x18\x01 \x01(\x0b\x32\x19.tendermint.abci.SnapshotR\x08snapshot\x12\x19\n\x08\x61pp_hash\x18\x02 \x01(\x0cR\x07\x61ppHash\"`\n\x18RequestLoadSnapshotChunk\x12\x16\n\x06height\x18\x01 \x01(\x04R\x06height\x12\x16\n\x06\x66ormat\x18\x02 \x01(\rR\x06\x66ormat\x12\x14\n\x05\x63hunk\x18\x03 \x01(\rR\x05\x63hunk\"_\n\x19RequestApplySnapshotChunk\x12\x14\n\x05index\x18\x01 \x01(\rR\x05index\x12\x14\n\x05\x63hunk\x18\x02 \x01(\x0cR\x05\x63hunk\x12\x16\n\x06sender\x18\x03 \x01(\tR\x06sender\"\x98\x03\n\x16RequestPrepareProposal\x12 \n\x0cmax_tx_bytes\x18\x01 \x01(\x03R\nmaxTxBytes\x12\x10\n\x03txs\x18\x02 \x03(\x0cR\x03txs\x12U\n\x11local_last_commit\x18\x03 \x01(\x0b\x32#.tendermint.abci.ExtendedCommitInfoB\x04\xc8\xde\x1f\x00R\x0flocalLastCommit\x12\x44\n\x0bmisbehavior\x18\x04 \x03(\x0b\x32\x1c.tendermint.abci.MisbehaviorB\x04\xc8\xde\x1f\x00R\x0bmisbehavior\x12\x16\n\x06height\x18\x05 \x01(\x03R\x06height\x12\x38\n\x04time\x18\x06 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x04time\x12\x30\n\x14next_validators_hash\x18\x07 \x01(\x0cR\x12nextValidatorsHash\x12)\n\x10proposer_address\x18\x08 \x01(\x0cR\x0fproposerAddress\"\x88\x03\n\x16RequestProcessProposal\x12\x10\n\x03txs\x18\x01 \x03(\x0cR\x03txs\x12S\n\x14proposed_last_commit\x18\x02 \x01(\x0b\x32\x1b.tendermint.abci.CommitInfoB\x04\xc8\xde\x1f\x00R\x12proposedLastCommit\x12\x44\n\x0bmisbehavior\x18\x03 \x03(\x0b\x32\x1c.tendermint.abci.MisbehaviorB\x04\xc8\xde\x1f\x00R\x0bmisbehavior\x12\x12\n\x04hash\x18\x04 \x01(\x0cR\x04hash\x12\x16\n\x06height\x18\x05 \x01(\x03R\x06height\x12\x38\n\x04time\x18\x06 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x04time\x12\x30\n\x14next_validators_hash\x18\x07 \x01(\x0cR\x12nextValidatorsHash\x12)\n\x10proposer_address\x18\x08 \x01(\x0cR\x0fproposerAddress\"\x83\x03\n\x11RequestExtendVote\x12\x12\n\x04hash\x18\x01 \x01(\x0cR\x04hash\x12\x16\n\x06height\x18\x02 \x01(\x03R\x06height\x12\x38\n\x04time\x18\x03 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x04time\x12\x10\n\x03txs\x18\x04 \x03(\x0cR\x03txs\x12S\n\x14proposed_last_commit\x18\x05 \x01(\x0b\x32\x1b.tendermint.abci.CommitInfoB\x04\xc8\xde\x1f\x00R\x12proposedLastCommit\x12\x44\n\x0bmisbehavior\x18\x06 \x03(\x0b\x32\x1c.tendermint.abci.MisbehaviorB\x04\xc8\xde\x1f\x00R\x0bmisbehavior\x12\x30\n\x14next_validators_hash\x18\x07 \x01(\x0cR\x12nextValidatorsHash\x12)\n\x10proposer_address\x18\x08 \x01(\x0cR\x0fproposerAddress\"\x9c\x01\n\x1aRequestVerifyVoteExtension\x12\x12\n\x04hash\x18\x01 \x01(\x0cR\x04hash\x12+\n\x11validator_address\x18\x02 \x01(\x0cR\x10validatorAddress\x12\x16\n\x06height\x18\x03 \x01(\x03R\x06height\x12%\n\x0evote_extension\x18\x04 \x01(\x0cR\rvoteExtension\"\x84\x03\n\x14RequestFinalizeBlock\x12\x10\n\x03txs\x18\x01 \x03(\x0cR\x03txs\x12Q\n\x13\x64\x65\x63ided_last_commit\x18\x02 \x01(\x0b\x32\x1b.tendermint.abci.CommitInfoB\x04\xc8\xde\x1f\x00R\x11\x64\x65\x63idedLastCommit\x12\x44\n\x0bmisbehavior\x18\x03 \x03(\x0b\x32\x1c.tendermint.abci.MisbehaviorB\x04\xc8\xde\x1f\x00R\x0bmisbehavior\x12\x12\n\x04hash\x18\x04 \x01(\x0cR\x04hash\x12\x16\n\x06height\x18\x05 \x01(\x03R\x06height\x12\x38\n\x04time\x18\x06 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x04time\x12\x30\n\x14next_validators_hash\x18\x07 \x01(\x0cR\x12nextValidatorsHash\x12)\n\x10proposer_address\x18\x08 \x01(\x0cR\x0fproposerAddress\"\x94\n\n\x08Response\x12\x42\n\texception\x18\x01 \x01(\x0b\x32\".tendermint.abci.ResponseExceptionH\x00R\texception\x12\x33\n\x04\x65\x63ho\x18\x02 \x01(\x0b\x32\x1d.tendermint.abci.ResponseEchoH\x00R\x04\x65\x63ho\x12\x36\n\x05\x66lush\x18\x03 \x01(\x0b\x32\x1e.tendermint.abci.ResponseFlushH\x00R\x05\x66lush\x12\x33\n\x04info\x18\x04 \x01(\x0b\x32\x1d.tendermint.abci.ResponseInfoH\x00R\x04info\x12\x43\n\ninit_chain\x18\x06 \x01(\x0b\x32\".tendermint.abci.ResponseInitChainH\x00R\tinitChain\x12\x36\n\x05query\x18\x07 \x01(\x0b\x32\x1e.tendermint.abci.ResponseQueryH\x00R\x05query\x12=\n\x08\x63heck_tx\x18\t \x01(\x0b\x32 .tendermint.abci.ResponseCheckTxH\x00R\x07\x63heckTx\x12\x39\n\x06\x63ommit\x18\x0c \x01(\x0b\x32\x1f.tendermint.abci.ResponseCommitH\x00R\x06\x63ommit\x12O\n\x0elist_snapshots\x18\r \x01(\x0b\x32&.tendermint.abci.ResponseListSnapshotsH\x00R\rlistSnapshots\x12O\n\x0eoffer_snapshot\x18\x0e \x01(\x0b\x32&.tendermint.abci.ResponseOfferSnapshotH\x00R\rofferSnapshot\x12\\\n\x13load_snapshot_chunk\x18\x0f \x01(\x0b\x32*.tendermint.abci.ResponseLoadSnapshotChunkH\x00R\x11loadSnapshotChunk\x12_\n\x14\x61pply_snapshot_chunk\x18\x10 \x01(\x0b\x32+.tendermint.abci.ResponseApplySnapshotChunkH\x00R\x12\x61pplySnapshotChunk\x12U\n\x10prepare_proposal\x18\x11 \x01(\x0b\x32(.tendermint.abci.ResponsePrepareProposalH\x00R\x0fprepareProposal\x12U\n\x10process_proposal\x18\x12 \x01(\x0b\x32(.tendermint.abci.ResponseProcessProposalH\x00R\x0fprocessProposal\x12\x46\n\x0b\x65xtend_vote\x18\x13 \x01(\x0b\x32#.tendermint.abci.ResponseExtendVoteH\x00R\nextendVote\x12\x62\n\x15verify_vote_extension\x18\x14 \x01(\x0b\x32,.tendermint.abci.ResponseVerifyVoteExtensionH\x00R\x13verifyVoteExtension\x12O\n\x0e\x66inalize_block\x18\x15 \x01(\x0b\x32&.tendermint.abci.ResponseFinalizeBlockH\x00R\rfinalizeBlockB\x07\n\x05valueJ\x04\x08\x05\x10\x06J\x04\x08\x08\x10\tJ\x04\x08\n\x10\x0bJ\x04\x08\x0b\x10\x0c\")\n\x11ResponseException\x12\x14\n\x05\x65rror\x18\x01 \x01(\tR\x05\x65rror\"(\n\x0cResponseEcho\x12\x18\n\x07message\x18\x01 \x01(\tR\x07message\"\x0f\n\rResponseFlush\"\xb8\x01\n\x0cResponseInfo\x12\x12\n\x04\x64\x61ta\x18\x01 \x01(\tR\x04\x64\x61ta\x12\x18\n\x07version\x18\x02 \x01(\tR\x07version\x12\x1f\n\x0b\x61pp_version\x18\x03 \x01(\x04R\nappVersion\x12*\n\x11last_block_height\x18\x04 \x01(\x03R\x0flastBlockHeight\x12-\n\x13last_block_app_hash\x18\x05 \x01(\x0cR\x10lastBlockAppHash\"\xc4\x01\n\x11ResponseInitChain\x12L\n\x10\x63onsensus_params\x18\x01 \x01(\x0b\x32!.tendermint.types.ConsensusParamsR\x0f\x63onsensusParams\x12\x46\n\nvalidators\x18\x02 \x03(\x0b\x32 .tendermint.abci.ValidatorUpdateB\x04\xc8\xde\x1f\x00R\nvalidators\x12\x19\n\x08\x61pp_hash\x18\x03 \x01(\x0cR\x07\x61ppHash\"\xf7\x01\n\rResponseQuery\x12\x12\n\x04\x63ode\x18\x01 \x01(\rR\x04\x63ode\x12\x10\n\x03log\x18\x03 \x01(\tR\x03log\x12\x12\n\x04info\x18\x04 \x01(\tR\x04info\x12\x14\n\x05index\x18\x05 \x01(\x03R\x05index\x12\x10\n\x03key\x18\x06 \x01(\x0cR\x03key\x12\x14\n\x05value\x18\x07 \x01(\x0cR\x05value\x12\x38\n\tproof_ops\x18\x08 \x01(\x0b\x32\x1b.tendermint.crypto.ProofOpsR\x08proofOps\x12\x16\n\x06height\x18\t \x01(\x03R\x06height\x12\x1c\n\tcodespace\x18\n \x01(\tR\tcodespace\"\xaa\x02\n\x0fResponseCheckTx\x12\x12\n\x04\x63ode\x18\x01 \x01(\rR\x04\x63ode\x12\x12\n\x04\x64\x61ta\x18\x02 \x01(\x0cR\x04\x64\x61ta\x12\x10\n\x03log\x18\x03 \x01(\tR\x03log\x12\x12\n\x04info\x18\x04 \x01(\tR\x04info\x12\x1e\n\ngas_wanted\x18\x05 \x01(\x03R\ngas_wanted\x12\x1a\n\x08gas_used\x18\x06 \x01(\x03R\x08gas_used\x12H\n\x06\x65vents\x18\x07 \x03(\x0b\x32\x16.tendermint.abci.EventB\x18\xc8\xde\x1f\x00\xea\xde\x1f\x10\x65vents,omitemptyR\x06\x65vents\x12\x1c\n\tcodespace\x18\x08 \x01(\tR\tcodespaceJ\x04\x08\t\x10\x0cR\x06senderR\x08priorityR\rmempool_error\"A\n\x0eResponseCommit\x12#\n\rretain_height\x18\x03 \x01(\x03R\x0cretainHeightJ\x04\x08\x01\x10\x02J\x04\x08\x02\x10\x03\"P\n\x15ResponseListSnapshots\x12\x37\n\tsnapshots\x18\x01 \x03(\x0b\x32\x19.tendermint.abci.SnapshotR\tsnapshots\"\xbe\x01\n\x15ResponseOfferSnapshot\x12\x45\n\x06result\x18\x01 \x01(\x0e\x32-.tendermint.abci.ResponseOfferSnapshot.ResultR\x06result\"^\n\x06Result\x12\x0b\n\x07UNKNOWN\x10\x00\x12\n\n\x06\x41\x43\x43\x45PT\x10\x01\x12\t\n\x05\x41\x42ORT\x10\x02\x12\n\n\x06REJECT\x10\x03\x12\x11\n\rREJECT_FORMAT\x10\x04\x12\x11\n\rREJECT_SENDER\x10\x05\"1\n\x19ResponseLoadSnapshotChunk\x12\x14\n\x05\x63hunk\x18\x01 \x01(\x0cR\x05\x63hunk\"\x98\x02\n\x1aResponseApplySnapshotChunk\x12J\n\x06result\x18\x01 \x01(\x0e\x32\x32.tendermint.abci.ResponseApplySnapshotChunk.ResultR\x06result\x12%\n\x0erefetch_chunks\x18\x02 \x03(\rR\rrefetchChunks\x12%\n\x0ereject_senders\x18\x03 \x03(\tR\rrejectSenders\"`\n\x06Result\x12\x0b\n\x07UNKNOWN\x10\x00\x12\n\n\x06\x41\x43\x43\x45PT\x10\x01\x12\t\n\x05\x41\x42ORT\x10\x02\x12\t\n\x05RETRY\x10\x03\x12\x12\n\x0eRETRY_SNAPSHOT\x10\x04\x12\x13\n\x0fREJECT_SNAPSHOT\x10\x05\"+\n\x17ResponsePrepareProposal\x12\x10\n\x03txs\x18\x01 \x03(\x0cR\x03txs\"\xa1\x01\n\x17ResponseProcessProposal\x12O\n\x06status\x18\x01 \x01(\x0e\x32\x37.tendermint.abci.ResponseProcessProposal.ProposalStatusR\x06status\"5\n\x0eProposalStatus\x12\x0b\n\x07UNKNOWN\x10\x00\x12\n\n\x06\x41\x43\x43\x45PT\x10\x01\x12\n\n\x06REJECT\x10\x02\";\n\x12ResponseExtendVote\x12%\n\x0evote_extension\x18\x01 \x01(\x0cR\rvoteExtension\"\xa5\x01\n\x1bResponseVerifyVoteExtension\x12Q\n\x06status\x18\x01 \x01(\x0e\x32\x39.tendermint.abci.ResponseVerifyVoteExtension.VerifyStatusR\x06status\"3\n\x0cVerifyStatus\x12\x0b\n\x07UNKNOWN\x10\x00\x12\n\n\x06\x41\x43\x43\x45PT\x10\x01\x12\n\n\x06REJECT\x10\x02\"\xea\x02\n\x15ResponseFinalizeBlock\x12H\n\x06\x65vents\x18\x01 \x03(\x0b\x32\x16.tendermint.abci.EventB\x18\xc8\xde\x1f\x00\xea\xde\x1f\x10\x65vents,omitemptyR\x06\x65vents\x12<\n\ntx_results\x18\x02 \x03(\x0b\x32\x1d.tendermint.abci.ExecTxResultR\ttxResults\x12S\n\x11validator_updates\x18\x03 \x03(\x0b\x32 .tendermint.abci.ValidatorUpdateB\x04\xc8\xde\x1f\x00R\x10validatorUpdates\x12Y\n\x17\x63onsensus_param_updates\x18\x04 \x01(\x0b\x32!.tendermint.types.ConsensusParamsR\x15\x63onsensusParamUpdates\x12\x19\n\x08\x61pp_hash\x18\x05 \x01(\x0cR\x07\x61ppHash\"Y\n\nCommitInfo\x12\x14\n\x05round\x18\x01 \x01(\x05R\x05round\x12\x35\n\x05votes\x18\x02 \x03(\x0b\x32\x19.tendermint.abci.VoteInfoB\x04\xc8\xde\x1f\x00R\x05votes\"i\n\x12\x45xtendedCommitInfo\x12\x14\n\x05round\x18\x01 \x01(\x05R\x05round\x12=\n\x05votes\x18\x02 \x03(\x0b\x32!.tendermint.abci.ExtendedVoteInfoB\x04\xc8\xde\x1f\x00R\x05votes\"z\n\x05\x45vent\x12\x12\n\x04type\x18\x01 \x01(\tR\x04type\x12]\n\nattributes\x18\x02 \x03(\x0b\x32\x1f.tendermint.abci.EventAttributeB\x1c\xc8\xde\x1f\x00\xea\xde\x1f\x14\x61ttributes,omitemptyR\nattributes\"N\n\x0e\x45ventAttribute\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value\x12\x14\n\x05index\x18\x03 \x01(\x08R\x05index\"\x80\x02\n\x0c\x45xecTxResult\x12\x12\n\x04\x63ode\x18\x01 \x01(\rR\x04\x63ode\x12\x12\n\x04\x64\x61ta\x18\x02 \x01(\x0cR\x04\x64\x61ta\x12\x10\n\x03log\x18\x03 \x01(\tR\x03log\x12\x12\n\x04info\x18\x04 \x01(\tR\x04info\x12\x1e\n\ngas_wanted\x18\x05 \x01(\x03R\ngas_wanted\x12\x1a\n\x08gas_used\x18\x06 \x01(\x03R\x08gas_used\x12H\n\x06\x65vents\x18\x07 \x03(\x0b\x32\x16.tendermint.abci.EventB\x18\xc8\xde\x1f\x00\xea\xde\x1f\x10\x65vents,omitemptyR\x06\x65vents\x12\x1c\n\tcodespace\x18\x08 \x01(\tR\tcodespace\"\x85\x01\n\x08TxResult\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05index\x18\x02 \x01(\rR\x05index\x12\x0e\n\x02tx\x18\x03 \x01(\x0cR\x02tx\x12;\n\x06result\x18\x04 \x01(\x0b\x32\x1d.tendermint.abci.ExecTxResultB\x04\xc8\xde\x1f\x00R\x06result\";\n\tValidator\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\x0cR\x07\x61\x64\x64ress\x12\x14\n\x05power\x18\x03 \x01(\x03R\x05power\"d\n\x0fValidatorUpdate\x12;\n\x07pub_key\x18\x01 \x01(\x0b\x32\x1c.tendermint.crypto.PublicKeyB\x04\xc8\xde\x1f\x00R\x06pubKey\x12\x14\n\x05power\x18\x02 \x01(\x03R\x05power\"\x93\x01\n\x08VoteInfo\x12>\n\tvalidator\x18\x01 \x01(\x0b\x32\x1a.tendermint.abci.ValidatorB\x04\xc8\xde\x1f\x00R\tvalidator\x12\x41\n\rblock_id_flag\x18\x03 \x01(\x0e\x32\x1d.tendermint.types.BlockIDFlagR\x0b\x62lockIdFlagJ\x04\x08\x02\x10\x03\"\xf3\x01\n\x10\x45xtendedVoteInfo\x12>\n\tvalidator\x18\x01 \x01(\x0b\x32\x1a.tendermint.abci.ValidatorB\x04\xc8\xde\x1f\x00R\tvalidator\x12%\n\x0evote_extension\x18\x03 \x01(\x0cR\rvoteExtension\x12/\n\x13\x65xtension_signature\x18\x04 \x01(\x0cR\x12\x65xtensionSignature\x12\x41\n\rblock_id_flag\x18\x05 \x01(\x0e\x32\x1d.tendermint.types.BlockIDFlagR\x0b\x62lockIdFlagJ\x04\x08\x02\x10\x03\"\x83\x02\n\x0bMisbehavior\x12\x34\n\x04type\x18\x01 \x01(\x0e\x32 .tendermint.abci.MisbehaviorTypeR\x04type\x12>\n\tvalidator\x18\x02 \x01(\x0b\x32\x1a.tendermint.abci.ValidatorB\x04\xc8\xde\x1f\x00R\tvalidator\x12\x16\n\x06height\x18\x03 \x01(\x03R\x06height\x12\x38\n\x04time\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x04time\x12,\n\x12total_voting_power\x18\x05 \x01(\x03R\x10totalVotingPower\"\x82\x01\n\x08Snapshot\x12\x16\n\x06height\x18\x01 \x01(\x04R\x06height\x12\x16\n\x06\x66ormat\x18\x02 \x01(\rR\x06\x66ormat\x12\x16\n\x06\x63hunks\x18\x03 \x01(\rR\x06\x63hunks\x12\x12\n\x04hash\x18\x04 \x01(\x0cR\x04hash\x12\x1a\n\x08metadata\x18\x05 \x01(\x0cR\x08metadata*9\n\x0b\x43heckTxType\x12\x10\n\x03NEW\x10\x00\x1a\x07\x8a\x9d \x03New\x12\x18\n\x07RECHECK\x10\x01\x1a\x0b\x8a\x9d \x07Recheck*K\n\x0fMisbehaviorType\x12\x0b\n\x07UNKNOWN\x10\x00\x12\x12\n\x0e\x44UPLICATE_VOTE\x10\x01\x12\x17\n\x13LIGHT_CLIENT_ATTACK\x10\x02\x32\x9d\x0b\n\x04\x41\x42\x43I\x12\x43\n\x04\x45\x63ho\x12\x1c.tendermint.abci.RequestEcho\x1a\x1d.tendermint.abci.ResponseEcho\x12\x46\n\x05\x46lush\x12\x1d.tendermint.abci.RequestFlush\x1a\x1e.tendermint.abci.ResponseFlush\x12\x43\n\x04Info\x12\x1c.tendermint.abci.RequestInfo\x1a\x1d.tendermint.abci.ResponseInfo\x12L\n\x07\x43heckTx\x12\x1f.tendermint.abci.RequestCheckTx\x1a .tendermint.abci.ResponseCheckTx\x12\x46\n\x05Query\x12\x1d.tendermint.abci.RequestQuery\x1a\x1e.tendermint.abci.ResponseQuery\x12I\n\x06\x43ommit\x12\x1e.tendermint.abci.RequestCommit\x1a\x1f.tendermint.abci.ResponseCommit\x12R\n\tInitChain\x12!.tendermint.abci.RequestInitChain\x1a\".tendermint.abci.ResponseInitChain\x12^\n\rListSnapshots\x12%.tendermint.abci.RequestListSnapshots\x1a&.tendermint.abci.ResponseListSnapshots\x12^\n\rOfferSnapshot\x12%.tendermint.abci.RequestOfferSnapshot\x1a&.tendermint.abci.ResponseOfferSnapshot\x12j\n\x11LoadSnapshotChunk\x12).tendermint.abci.RequestLoadSnapshotChunk\x1a*.tendermint.abci.ResponseLoadSnapshotChunk\x12m\n\x12\x41pplySnapshotChunk\x12*.tendermint.abci.RequestApplySnapshotChunk\x1a+.tendermint.abci.ResponseApplySnapshotChunk\x12\x64\n\x0fPrepareProposal\x12\'.tendermint.abci.RequestPrepareProposal\x1a(.tendermint.abci.ResponsePrepareProposal\x12\x64\n\x0fProcessProposal\x12\'.tendermint.abci.RequestProcessProposal\x1a(.tendermint.abci.ResponseProcessProposal\x12U\n\nExtendVote\x12\".tendermint.abci.RequestExtendVote\x1a#.tendermint.abci.ResponseExtendVote\x12p\n\x13VerifyVoteExtension\x12+.tendermint.abci.RequestVerifyVoteExtension\x1a,.tendermint.abci.ResponseVerifyVoteExtension\x12^\n\rFinalizeBlock\x12%.tendermint.abci.RequestFinalizeBlock\x1a&.tendermint.abci.ResponseFinalizeBlockB)Z\'github.com/cometbft/cometbft/abci/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'tendermint.abci.types_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\'github.com/cometbft/cometbft/abci/types'
  _globals['_CHECKTXTYPE'].values_by_name["NEW"]._loaded_options = None
  _globals['_CHECKTXTYPE'].values_by_name["NEW"]._serialized_options = b'\212\235 \003New'
  _globals['_CHECKTXTYPE'].values_by_name["RECHECK"]._loaded_options = None
  _globals['_CHECKTXTYPE'].values_by_name["RECHECK"]._serialized_options = b'\212\235 \007Recheck'
  _globals['_REQUESTINITCHAIN'].fields_by_name['time']._loaded_options = None
  _globals['_REQUESTINITCHAIN'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_REQUESTINITCHAIN'].fields_by_name['validators']._loaded_options = None
  _globals['_REQUESTINITCHAIN'].fields_by_name['validators']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTPREPAREPROPOSAL'].fields_by_name['local_last_commit']._loaded_options = None
  _globals['_REQUESTPREPAREPROPOSAL'].fields_by_name['local_last_commit']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTPREPAREPROPOSAL'].fields_by_name['misbehavior']._loaded_options = None
  _globals['_REQUESTPREPAREPROPOSAL'].fields_by_name['misbehavior']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTPREPAREPROPOSAL'].fields_by_name['time']._loaded_options = None
  _globals['_REQUESTPREPAREPROPOSAL'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_REQUESTPROCESSPROPOSAL'].fields_by_name['proposed_last_commit']._loaded_options = None
  _globals['_REQUESTPROCESSPROPOSAL'].fields_by_name['proposed_last_commit']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTPROCESSPROPOSAL'].fields_by_name['misbehavior']._loaded_options = None
  _globals['_REQUESTPROCESSPROPOSAL'].fields_by_name['misbehavior']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTPROCESSPROPOSAL'].fields_by_name['time']._loaded_options = None
  _globals['_REQUESTPROCESSPROPOSAL'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_REQUESTEXTENDVOTE'].fields_by_name['time']._loaded_options = None
  _globals['_REQUESTEXTENDVOTE'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_REQUESTEXTENDVOTE'].fields_by_name['proposed_last_commit']._loaded_options = None
  _globals['_REQUESTEXTENDVOTE'].fields_by_name['proposed_last_commit']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTEXTENDVOTE'].fields_by_name['misbehavior']._loaded_options = None
  _globals['_REQUESTEXTENDVOTE'].fields_by_name['misbehavior']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTFINALIZEBLOCK'].fields_by_name['decided_last_commit']._loaded_options = None
  _globals['_REQUESTFINALIZEBLOCK'].fields_by_name['decided_last_commit']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTFINALIZEBLOCK'].fields_by_name['misbehavior']._loaded_options = None
  _globals['_REQUESTFINALIZEBLOCK'].fields_by_name['misbehavior']._serialized_options = b'\310\336\037\000'
  _globals['_REQUESTFINALIZEBLOCK'].fields_by_name['time']._loaded_options = None
  _globals['_REQUESTFINALIZEBLOCK'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_RESPONSEINITCHAIN'].fields_by_name['validators']._loaded_options = None
  _globals['_RESPONSEINITCHAIN'].fields_by_name['validators']._serialized_options = b'\310\336\037\000'
  _globals['_RESPONSECHECKTX'].fields_by_name['events']._loaded_options = None
  _globals['_RESPONSECHECKTX'].fields_by_name['events']._serialized_options = b'\310\336\037\000\352\336\037\020events,omitempty'
  _globals['_RESPONSEFINALIZEBLOCK'].fields_by_name['events']._loaded_options = None
  _globals['_RESPONSEFINALIZEBLOCK'].fields_by_name['events']._serialized_options = b'\310\336\037\000\352\336\037\020events,omitempty'
  _globals['_RESPONSEFINALIZEBLOCK'].fields_by_name['validator_updates']._loaded_options = None
  _globals['_RESPONSEFINALIZEBLOCK'].fields_by_name['validator_updates']._serialized_options = b'\310\336\037\000'
  _globals['_COMMITINFO'].fields_by_name['votes']._loaded_options = None
  _globals['_COMMITINFO'].fields_by_name['votes']._serialized_options = b'\310\336\037\000'
  _globals['_EXTENDEDCOMMITINFO'].fields_by_name['votes']._loaded_options = None
  _globals['_EXTENDEDCOMMITINFO'].fields_by_name['votes']._serialized_options = b'\310\336\037\000'
  _globals['_EVENT'].fields_by_name['attributes']._loaded_options = None
  _globals['_EVENT'].fields_by_name['attributes']._serialized_options = b'\310\336\037\000\352\336\037\024attributes,omitempty'
  _globals['_EXECTXRESULT'].fields_by_name['events']._loaded_options = None
  _globals['_EXECTXRESULT'].fields_by_name['events']._serialized_options = b'\310\336\037\000\352\336\037\020events,omitempty'
  _globals['_TXRESULT'].fields_by_name['result']._loaded_options = None
  _globals['_TXRESULT'].fields_by_name['result']._serialized_options = b'\310\336\037\000'
  _globals['_VALIDATORUPDATE'].fields_by_name['pub_key']._loaded_options = None
  _globals['_VALIDATORUPDATE'].fields_by_name['pub_key']._serialized_options = b'\310\336\037\000'
  _globals['_VOTEINFO'].fields_by_name['validator']._loaded_options = None
  _globals['_VOTEINFO'].fields_by_name['validator']._serialized_options = b'\310\336\037\000'
  _globals['_EXTENDEDVOTEINFO'].fields_by_name['validator']._loaded_options = None
  _globals['_EXTENDEDVOTEINFO'].fields_by_name['validator']._serialized_options = b'\310\336\037\000'
  _globals['_MISBEHAVIOR'].fields_by_name['validator']._loaded_options = None
  _globals['_MISBEHAVIOR'].fields_by_name['validator']._serialized_options = b'\310\336\037\000'
  _globals['_MISBEHAVIOR'].fields_by_name['time']._loaded_options = None
  _globals['_MISBEHAVIOR'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_CHECKTXTYPE']._serialized_start=9832
  _globals['_CHECKTXTYPE']._serialized_end=9889
  _globals['_MISBEHAVIORTYPE']._serialized_start=9891
  _globals['_MISBEHAVIORTYPE']._serialized_end=9966
  _globals['_REQUEST']._serialized_start=230
  _globals['_REQUEST']._serialized_end=1445
  _globals['_REQUESTECHO']._serialized_start=1447
  _globals['_REQUESTECHO']._serialized_end=1486
  _globals['_REQUESTFLUSH']._serialized_start=1488
  _globals['_REQUESTFLUSH']._serialized_end=1502
  _globals['_REQUESTINFO']._serialized_start=1505
  _globals['_REQUESTINFO']._serialized_end=1649
  _globals['_REQUESTINITCHAIN']._serialized_start=1652
  _globals['_REQUESTINITCHAIN']._serialized_end=1984
  _globals['_REQUESTQUERY']._serialized_start=1986
  _globals['_REQUESTQUERY']._serialized_end=2086
  _globals['_REQUESTCHECKTX']._serialized_start=2088
  _globals['_REQUESTCHECKTX']._serialized_end=2170
  _globals['_REQUESTCOMMIT']._serialized_start=2172
  _globals['_REQUESTCOMMIT']._serialized_end=2187
  _globals['_REQUESTLISTSNAPSHOTS']._serialized_start=2189
  _globals['_REQUESTLISTSNAPSHOTS']._serialized_end=2211
  _globals['_REQUESTOFFERSNAPSHOT']._serialized_start=2213
  _globals['_REQUESTOFFERSNAPSHOT']._serialized_end=2317
  _globals['_REQUESTLOADSNAPSHOTCHUNK']._serialized_start=2319
  _globals['_REQUESTLOADSNAPSHOTCHUNK']._serialized_end=2415
  _globals['_REQUESTAPPLYSNAPSHOTCHUNK']._serialized_start=2417
  _globals['_REQUESTAPPLYSNAPSHOTCHUNK']._serialized_end=2512
  _globals['_REQUESTPREPAREPROPOSAL']._serialized_start=2515
  _globals['_REQUESTPREPAREPROPOSAL']._serialized_end=2923
  _globals['_REQUESTPROCESSPROPOSAL']._serialized_start=2926
  _globals['_REQUESTPROCESSPROPOSAL']._serialized_end=3318
  _globals['_REQUESTEXTENDVOTE']._serialized_start=3321
  _globals['_REQUESTEXTENDVOTE']._serialized_end=3708
  _globals['_REQUESTVERIFYVOTEEXTENSION']._serialized_start=3711
  _globals['_REQUESTVERIFYVOTEEXTENSION']._serialized_end=3867
  _globals['_REQUESTFINALIZEBLOCK']._serialized_start=3870
  _globals['_REQUESTFINALIZEBLOCK']._serialized_end=4258
  _globals['_RESPONSE']._serialized_start=4261
  _globals['_RESPONSE']._serialized_end=5561
  _globals['_RESPONSEEXCEPTION']._serialized_start=5563
  _globals['_RESPONSEEXCEPTION']._serialized_end=5604
  _globals['_RESPONSEECHO']._serialized_start=5606
  _globals['_RESPONSEECHO']._serialized_end=5646
  _globals['_RESPONSEFLUSH']._serialized_start=5648
  _globals['_RESPONSEFLUSH']._serialized_end=5663
  _globals['_RESPONSEINFO']._serialized_start=5666
  _globals['_RESPONSEINFO']._serialized_end=5850
  _globals['_RESPONSEINITCHAIN']._serialized_start=5853
  _globals['_RESPONSEINITCHAIN']._serialized_end=6049
  _globals['_RESPONSEQUERY']._serialized_start=6052
  _globals['_RESPONSEQUERY']._serialized_end=6299
  _globals['_RESPONSECHECKTX']._serialized_start=6302
  _globals['_RESPONSECHECKTX']._serialized_end=6600
  _globals['_RESPONSECOMMIT']._serialized_start=6602
  _globals['_RESPONSECOMMIT']._serialized_end=6667
  _globals['_RESPONSELISTSNAPSHOTS']._serialized_start=6669
  _globals['_RESPONSELISTSNAPSHOTS']._serialized_end=6749
  _globals['_RESPONSEOFFERSNAPSHOT']._serialized_start=6752
  _globals['_RESPONSEOFFERSNAPSHOT']._serialized_end=6942
  _globals['_RESPONSEOFFERSNAPSHOT_RESULT']._serialized_start=6848
  _globals['_RESPONSEOFFERSNAPSHOT_RESULT']._serialized_end=6942
  _globals['_RESPONSELOADSNAPSHOTCHUNK']._serialized_start=6944
  _globals['_RESPONSELOADSNAPSHOTCHUNK']._serialized_end=6993
  _globals['_RESPONSEAPPLYSNAPSHOTCHUNK']._serialized_start=6996
  _globals['_RESPONSEAPPLYSNAPSHOTCHUNK']._serialized_end=7276
  _globals['_RESPONSEAPPLYSNAPSHOTCHUNK_RESULT']._serialized_start=7180
  _globals['_RESPONSEAPPLYSNAPSHOTCHUNK_RESULT']._serialized_end=7276
  _globals['_RESPONSEPREPAREPROPOSAL']._serialized_start=7278
  _globals['_RESPONSEPREPAREPROPOSAL']._serialized_end=7321
  _globals['_RESPONSEPROCESSPROPOSAL']._serialized_start=7324
  _globals['_RESPONSEPROCESSPROPOSAL']._serialized_end=7485
  _globals['_RESPONSEPROCESSPROPOSAL_PROPOSALSTATUS']._serialized_start=7432
  _globals['_RESPONSEPROCESSPROPOSAL_PROPOSALSTATUS']._serialized_end=7485
  _globals['_RESPONSEEXTENDVOTE']._serialized_start=7487
  _globals['_RESPONSEEXTENDVOTE']._serialized_end=7546
  _globals['_RESPONSEVERIFYVOTEEXTENSION']._serialized_start=7549
  _globals['_RESPONSEVERIFYVOTEEXTENSION']._serialized_end=7714
  _globals['_RESPONSEVERIFYVOTEEXTENSION_VERIFYSTATUS']._serialized_start=7663
  _globals['_RESPONSEVERIFYVOTEEXTENSION_VERIFYSTATUS']._serialized_end=7714
  _globals['_RESPONSEFINALIZEBLOCK']._serialized_start=7717
  _globals['_RESPONSEFINALIZEBLOCK']._serialized_end=8079
  _globals['_COMMITINFO']._serialized_start=8081
  _globals['_COMMITINFO']._serialized_end=8170
  _globals['_EXTENDEDCOMMITINFO']._serialized_start=8172
  _globals['_EXTENDEDCOMMITINFO']._serialized_end=8277
  _globals['_EVENT']._serialized_start=8279
  _globals['_EVENT']._serialized_end=8401
  _globals['_EVENTATTRIBUTE']._serialized_start=8403
  _globals['_EVENTATTRIBUTE']._serialized_end=8481
  _globals['_EXECTXRESULT']._serialized_start=8484
  _globals['_EXECTXRESULT']._serialized_end=8740
  _globals['_TXRESULT']._serialized_start=8743
  _globals['_TXRESULT']._serialized_end=8876
  _globals['_VALIDATOR']._serialized_start=8878
  _globals['_VALIDATOR']._serialized_end=8937
  _globals['_VALIDATORUPDATE']._serialized_start=8939
  _globals['_VALIDATORUPDATE']._serialized_end=9039
  _globals['_VOTEINFO']._serialized_start=9042
  _globals['_VOTEINFO']._serialized_end=9189
  _globals['_EXTENDEDVOTEINFO']._serialized_start=9192
  _globals['_EXTENDEDVOTEINFO']._serialized_end=9435
  _globals['_MISBEHAVIOR']._serialized_start=9438
  _globals['_MISBEHAVIOR']._serialized_end=9697
  _globals['_SNAPSHOT']._serialized_start=9700
  _globals['_SNAPSHOT']._serialized_end=9830
  _globals['_ABCI']._serialized_start=9969
  _globals['_ABCI']._serialized_end=11406
# @@protoc_insertion_point(module_scope)
