# Vote Examiner Utility

This utility provides tools to examine validator participation in a Cosmos SDK-based blockchain, specifically focusing on vote extensions and block consensus.

## Build

To build the utility, navigate to the `utils/vote-examiner` directory and run:

```bash
go build
```

This will produce a `vote-examiner` executable in the current directory.

## Modes of Operation

The utility has two main modes: Vote Extension Analysis and Block Consensus Report.

### 1. Vote Extension Analysis Mode (Default)

This mode analyzes the vote extensions submitted by validators in a given block. It shows which validators submitted extensions, the content of these extensions, participation statistics, and highlights differences between submitted extensions.

**Usage:**

```bash
./vote-examiner [flags]
```

**Flags for Vote Extension Analysis:**

*   `-file <filepath>`: (Optional) Use a local file containing the JSON output of a block query (specifically the `data.txs[0]` part, base64 decoded, which holds the `ConsensusData`) instead of querying an RPC node. This is useful for analyzing past blocks or offline analysis.
    *Example format for the file content: `{"ConsensusData":{"votes":[{"validator":{"address":"...","power":...},"vote_extension":"..."},...]}}`*
*   `-network <name>`: (Optional) Specify the network to connect to. Defaults to `mainnet`. 
    *   Supported: `localnet` (or `local`, `localhost`), `devnet` (or `dev`, `amber`), `testnet` (or `test`, `gardia`), `mainnet` (or `main`, `diamond`).
    *   This flag is ignored if `-node` is provided.
*   `-node <rpc_url>`: (Optional) Specify the RPC node URL directly (e.g., `http://localhost:26657`). Overrides the `-network` flag.
*   `-height <block_height>`: (Optional) Specify the block height to query. Defaults to the latest block.
*   `-missing-only`: (Optional) If set, only displays validators who did *not* submit a vote extension. This is useful for quickly identifying non-participating validators.

**Example (Vote Extension Analysis):**

```bash
# Analyze latest block on mainnet for vote extensions
./vote-examiner

# Analyze block 12345 on testnet
./vote-examiner -network testnet -height 12345

# Analyze block 5000 using a specific RPC node and show only missing extensions
./vote-examiner -node https://rpc.example.com:443 -height 5000 -missing-only

# Analyze vote extensions from a local file
./vote-examiner -file /path/to/block_consensus_data.json
```

### 2. Block Consensus Report Mode

This mode provides a detailed report on validator participation for a specific block, based on the commit signatures found in the subsequent block. It categorizes validators from the active set by their recorded voting status:
*   **COMMIT**: Voted for the specific block hash.
*   **NIL**: Voted NIL, indicating abstention or disagreement with the proposed block.
*   **ABSENT**: Recorded by the system with an ABSENT status (vote not received).
*   **NO SIGNATURE FOUND**: No signature was found in the commit data from the next block.

The report also includes voting power percentages for each category.

**Usage:**

```bash
./vote-examiner -consensus-report [flags]
```

**Flags for Block Consensus Report:**

*   `-consensus-report`: **Required** to activate this mode.
*   `-network <name>`: (Optional) Specify the network. Same as in vote extension mode. Defaults to `mainnet`.
*   `-node <rpc_url>`: (Optional) Specify the RPC node URL. Same as in vote extension mode. Overrides `-network`.
*   `-height <block_height>`: (Optional) Specify the block height for the consensus report. Defaults to the latest block.
    *Note: This mode fetches data for the specified height AND the next block (height+1) to get signatures.*

**Example (Block Consensus Report):**

```bash
# Get consensus report for the latest block on mainnet
./vote-examiner -consensus-report

# Get consensus report for block 12345 on testnet
./vote-examiner -consensus-report -network testnet -height 12345

# Get consensus report for block 5000 using a specific RPC node
./vote-examiner -consensus-report -node https://rpc.example.com:443 -height 5000
```

## Validator Monikers

For both modes, the utility attempts to fetch validator monikers using `zenrockd` CLI commands (`zenrockd q consensus comet validator-set` and `zenrockd q validation validators`) via the specified RPC node. Ensure `zenrockd` is in your `PATH` and configured correctly if you want to see monikers displayed alongside validator addresses.
If fetching monikers fails, the program will proceed using addresses only or "Unknown" for monikers. 