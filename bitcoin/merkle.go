package bitcoin

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func VerifyBTCLockTransaction(rawTX string, chainName string, index int, proof []string, blockHeader *api.BTCBlockHeader, ignoreAddresses []string) ([]TXOutputs, string, error) {
	//1st Check the blockheader is valid
	err := CheckBlockHeader(blockHeader)
	if err != nil {
		return nil, "", fmt.Errorf("Fail to Check BlockHeader: %w", err)
	}

	merkleRootBytes, err := hex.DecodeString(blockHeader.MerkleRoot)
	merkleRootBytes = ReverseBytes(merkleRootBytes)
	calculatedTxID, err := CalculateTXID(rawTX, chainName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to calculate txid: %w", err)
	}
	calculatedIDString := ReverseHex(calculatedTxID.String())

	targetHash, err := chainhash.NewHashFromStr(calculatedIDString)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse tx hash: %w", err)
	}
	i := index
	for _, sibling := range proof {
		siblingHash, _ := chainhash.NewHashFromStr(sibling)
		if i%2 == 0 {
			targetHash = MergeHashes(targetHash, siblingHash)
		} else {
			targetHash = MergeHashes(siblingHash, targetHash)
		}
		i /= 2
	}
	merkleRootHash, err := chainhash.NewHash(merkleRootBytes)
	if err != nil {
		return nil, "", err
	}

	//invalid merkle verification
	if !targetHash.IsEqual(merkleRootHash) {
		return nil, "", fmt.Errorf("invalid merkle verification")
	}

	//Verifies, get the outputs
	outputs, err := DecodeOutputs(rawTX, chainName)
	if err != nil {
		return nil, "", err
	}

	//Remove ignoreAddresses from outputs
	cleanedOutputs := filterTXOutputs(outputs, ignoreAddresses)

	return cleanedOutputs, calculatedIDString, nil
}

func filterTXOutputs(outputs []TXOutputs, ignoreAddresses []string) []TXOutputs {
	ignoreMap := make(map[string]struct{}, len(ignoreAddresses))

	// Populate the ignoreMap with addresses to be ignored
	for _, addr := range ignoreAddresses {
		ignoreMap[addr] = struct{}{}
	}

	var filteredOutputs []TXOutputs
	for _, output := range outputs {
		if _, found := ignoreMap[output.Address]; !found {
			filteredOutputs = append(filteredOutputs, output)
		}
	}

	return filteredOutputs
}

func CheckBlockHeader(b *api.BTCBlockHeader) error {
	ok, err := deriveBlockHash(b)
	if err != nil {
		return fmt.Errorf("fail to derive blockhash %w", err)
	}
	if ok != true {
		return fmt.Errorf("invalid blockhash")
	}

	ok, err = blockHashCompliesWithDifficulty(b)
	if err != nil {
		return fmt.Errorf("fail to calculate difficulty compliance %w", err)
	}
	if ok != true {
		return fmt.Errorf("blockhash does not comply with difficulty")
	}
	return nil
}

func deriveBlockHash(b *api.BTCBlockHeader) (bool, error) {
	// Check if this is a Zcash header (has NonceHex populated)
	isZcash := b.NonceHex != ""

	var buf []byte
	if isZcash {
		// Zcash block headers are 140 bytes
		buf = make([]byte, 140)
	} else {
		// Bitcoin block headers are 80 bytes
		buf = make([]byte, 80)
	}

	version := b.Version
	prevBlock, err := hex.DecodeString(ReverseHex(b.PrevBlock))
	if err != nil {
		return false, err
	}
	merkleRoot, err := hex.DecodeString(ReverseHex(b.MerkleRoot))
	if err != nil {
		return false, err
	}
	bitHex := fmt.Sprintf("%08x", b.Bits)
	bits, err := hex.DecodeString(ReverseHex(bitHex))
	if err != nil {
		return false, err
	}

	binary.LittleEndian.PutUint32(buf[0:4], uint32(version))
	copy(buf[4:36], prevBlock)
	copy(buf[36:68], merkleRoot)

	if isZcash {
		// Zcash has a 32-byte block commitments field (hashReserved) after merkleRoot
		// BlockCommitments from RPC is in big-endian (display order), reverse to little-endian for header
		blockCommitments, err := hex.DecodeString(ReverseHex(b.BlockCommitments))
		if err != nil {
			return false, fmt.Errorf("failed to decode Zcash BlockCommitments: %w", err)
		}
		if len(blockCommitments) != 32 {
			return false, fmt.Errorf("invalid Zcash BlockCommitments length: expected 32 bytes, got %d", len(blockCommitments))
		}
		copy(buf[68:100], blockCommitments)

		binary.LittleEndian.PutUint32(buf[100:104], uint32(b.TimeStamp))
		copy(buf[104:108], bits)

		// Zcash uses a 256-bit (32-byte) nonce
		// The nonce from RPC needs to be reversed for the header
		nonceBytes, err := hex.DecodeString(ReverseHex(b.NonceHex))
		if err != nil {
			return false, fmt.Errorf("failed to decode Zcash nonce: %w", err)
		}
		if len(nonceBytes) != 32 {
			return false, fmt.Errorf("invalid Zcash nonce length: expected 32 bytes, got %d", len(nonceBytes))
		}
		copy(buf[108:140], nonceBytes)
	} else {
		// Bitcoin format
		binary.LittleEndian.PutUint32(buf[68:72], uint32(b.TimeStamp))
		copy(buf[72:76], bits)
		binary.LittleEndian.PutUint32(buf[76:80], uint32(b.Nonce))
	}

	//hash := doubleSha256(serializedHeader)
	// For Zcash, include the solution in the hash calculation
	var hashInput []byte
	if isZcash && b.Solution != "" {
		// Decode the solution and append it to the header
		solutionBytes, err := hex.DecodeString(b.Solution)
		if err != nil {
			return false, fmt.Errorf("failed to decode Zcash solution: %w", err)
		}
		hashInput = append(buf, solutionBytes...)
	} else {
		hashInput = buf
	}

	first := sha256.Sum256(hashInput)
	second := sha256.Sum256(first[:])

	// For both Bitcoin and Zcash, use the hash directly
	// chainhash.NewHash expects bytes in internal format (little-endian for block hashes)
	// SHA256 outputs big-endian, so we don't reverse
	blockHash, err := chainhash.NewHash(second[:])
	if err != nil {
		return false, err
	}
	chainBlockHash, err := chainhash.NewHashFromStr(b.BlockHash)
	if !blockHash.IsEqual(chainBlockHash) {
		return false, err
	}
	return true, nil
}

func blockHashCompliesWithDifficulty(b *api.BTCBlockHeader) (bool, error) {
	hashInt, validHash := new(big.Int).SetString(b.BlockHash, 16)
	if !validHash {
		return false, fmt.Errorf("error convert Convert Hash to big.Int")
	}
	bitHex := fmt.Sprintf("%08x", b.Bits)
	bitsHex, err := hex.DecodeString(bitHex)
	if err != nil {
		return false, fmt.Errorf("error convert Convert Hash to big.Int")
	}
	bts := binary.BigEndian.Uint32(bitsHex)
	target := calcTarget(bts)
	if hashInt.Cmp(target) > 0 {
		return false, fmt.Errorf("Error: Block hash does not meet the difficulty target")
	}
	return true, nil
}

func calcTarget(bits uint32) (target *big.Int) {
	// Get the bytes from bits
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)

	exponent := int(bytes[0])
	var coefficient uint32
	coefficient = uint32(bytes[1])<<16 + uint32(bytes[2])<<8 + uint32(bytes[3])

	target = new(big.Int).Lsh(big.NewInt(int64(coefficient)), uint(8*(exponent-3)))
	return target
}
