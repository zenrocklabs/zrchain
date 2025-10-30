package bitcoin

import (
	"fmt"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/require"
)

func Test_VerifyBTCUnlockTransaction(t *testing.T) {
	t.Skip("Skipping on CI - requires Bitcoin block data validation")

	// Test takes data from Proxy + merkleRoot from the Neutrino Node
	// VerifyBTCLockTransaction recalculates the Transaction ID
	// Using that calculated TXid hashs against the proof and checks the result is the same as the Merkle Root
	// If valid VerifyBTCLockTransaction also returns a set of outputs
	// Any output which is a valid ZenBTC address can mint ZenBTC
	// Store the TXID

	// 3 Fields from the Proxy

	index := 240
	rawTX := "01000000000102471371d5514faf456d8e6108c8f048c2d20d1e4a60bbd78cc3e98ff6daa8b4ee0000000000ffffffff6b789d277a9540deef30d4e75578b9890c2a40c792b26561f540efa4ed85cd8f0000000000ffffffff012e0b0600000000001976a9141b3f406c1e5eafccdfa0a3a3ff48aa6b020afb4e88ac0141735bc5871028c17f242010d86b4ccd462c0058462940a7b8aa5f07a9019a9c3565aa2578fb4772f29027ac7128690b78716ee702bff36aa2bebc1ff4388bb0da010141254cc51e219e22ac1e22b15a32a64daf56e8f87216df4f653441a0ccba1b559324750f9711cc0d2d1e08084c7a2db3ed8313622f4f884d72a54e2e56f26266270100000000"
	proof := []string{
		"9184bd1b362f72344f02e7def078fd187933ab03f843776641e94b11dc1d10f9",
		"7fdd1ce23535050d3ee83b88a2f5090ff3590220a9c649dfbba148f37671a1fa",
		"a01b1960df18f207976c7fecab97d916937f1f541a09167557be39bfdfa2a827",
		"85edccf2f8aa0878153472fa1b5d4c536cd61c0247bd39de68a1f26c61c50d80",
		"99dc8b031d2c606d267f22702a527736fa51a8618346b2934b192ea9e31a016a",
		"594e2c7446aa225cb987dfe86bf230b46f5339dae6835a7bdec18293def91ff8",
		"3528fd1f6f35771c70f7f91ae532a0827184377497b5d604f311769233f1712e",
		"40c94dc15e7eed3dff4e0ac2bb84e1756292d7af53856a57ceb3a1d417a76d15",
		"1f578553fa6159e50221b44d81373da02e9c604f75395412f33486bb90017781",
		"48df9cf4044d402c53bae68b536c8f4c9e3355040e3907fb2f93cd6d04ba2830",
		"dd1936b9bf265c51bfb41547ccd4739392058d3731495c8374fa61ae7a4321ba",
		"12160c339ad1aae42884c1bf1b55ba2750446792ff54d5c77f64592db269c8bd",
	}

	//1 Field from the Neutrino Node
	blockHeader := &api.BTCBlockHeader{
		Version:    543301632,
		PrevBlock:  "0000000000000000000063c2c472446997e1dd3abcb1605e0058024d15279635",
		MerkleRoot: "4f22b75ba55d0d7180b42af565fe4e11a113161970528f4fdaaae8b4858e6c61",
		TimeStamp:  1708539953,
		Bits:       386101681,
		Nonce:      2275713329,
		BlockHash:  "0000000000000000000207bd7147ace5cda4cead0cf46babdfb8e6addfb0f8e2",
	}

	outputs, txid, err := VerifyBTCLockTransaction(rawTX, "mainnet", index, proof, blockHeader, []string{})
	require.NoError(t, err)
	require.True(t, txid == "0a9787e73cee590730cca5b787f8f1fd9c2205a7fa2c8f165b89fc43b9e22cfa", "invalid calculated txid")
	require.True(t, len(outputs) == 1, "should be two outputs") //in this eg. 1 is real, 1 has 0 amount, no address
	fmt.Println(outputs)
}

func Test_VerifyBTCLockTransaction(t *testing.T) {
	t.Skip("Skipping on CI - requires Bitcoin block data validation")

	// Test takes data from Proxy + merkleRoot from the Neutrino Node
	// VerifyBTCLockTransaction recalculates the Transaction ID
	// Using that calculated TXid hashs against the proof and checks the result is the same as the Merkle Root
	// If valid VerifyBTCLockTransaction also returns a set of outputs
	// Any output which is a valid ZenBTC address can mint ZenBTC
	// Store the TXID

	// 3 Fields from the Proxy
	index := 3462
	rawTX := "0200000001252ef4787fe208a370fecf707674605ca4c803036e636ebdcbc0a287ca12f74c000000006b483045022100caa45533d411ce7140d6e9e9a0a1d00947304d09e546d1c725a2ef5679f11ba202201fc19ec89b4b2def1119fd0d95853d9f2e2f118867cf8dbbb60e5c8a793bd1d20121021848a8e26bf13437b711cf2504083650722a98a15d00283575260737cb4c050dffffffff0275da4c04000000001976a91424299bb977941a79afba01728c34d018ae02c7c888ac0000000000000000536a4c5048454d49010064f9080015ce6c8c4ea762b9a5fccc0caaf96953afa9ab4d3bb24ae2665ab175216fff8f3525c402168d60e83afc7f1ba011da7f7857612c4f35790dba5b93840c1db1f8945ae5ec9b6569492c00"
	proof := []string{
		"45f55616ce67e7eb8c92fb0b25dd342f61e1d8906f300e372e64260e464b7c03",
		"0a184d91e652c35d987072d5f1ce33c77a411b3f470cf185097d035a275f0e1c",
		"ba6526507eb6601eb5ab1113aca7429d888e8fc1e77562766eb2d554924ad489",
		"c1b30eaea97e20d61c2311260811080f9f35daf57f134bffb167eef0bfbd8337",
		"e679739ab8a223122552cf32a94bcfaa43492258d15dd5f08bca00d472d64d89",
		"1f15e67691be2b5e6de6384913e2a98bac62690778d7d43e26a5ef8cf543ce80",
		"a7ffdc429d24c3e26a4d54e94e3af88db6507b6d615ef4dd720e91982e4ce02c",
		"c20c42d2d3dfacc90e35d38cd328e73e1e1c76596d4b11af0bc9823480d0e6ed",
		"ed246e78bd59bf29d9d3b912e872fbf87f8b9f1757b9917d99644fe7dcc717c0",
		"3488bb935d647d178d7bdfa442a5cd02923693fadca295c0bfcc33b184ee07d2",
		"554f43d6696d786092c4017d60b3a3ff1b7ae92911e9d702a2af66c52bc519f3",
		"c6b371722f132ee15c2dc09d4c14b40149e0507a73d89f31b5abf6bdaa702694",
	}

	//1 Field from the Neutrino Node
	blockHeader := &api.BTCBlockHeader{
		Version:    536870912,
		PrevBlock:  "000000000000000f10b5de36d015586d3bf3f63a0faa418b73cb91aaff5de064",
		MerkleRoot: "503d9d9abd29b58233568efded9af1516f395fea7ca3eb1aaafdc8b63406d214",
		TimeStamp:  1725270640,
		Bits:       486604799,
		Nonce:      3645002625,
		BlockHash:  "000000003c057188e715529359ae9a45a9f7b8462f0c9b8b9f2f24ac196eb6fb",
	}

	outputs, txid, err := VerifyBTCLockTransaction(rawTX, "testnet", index, proof, blockHeader, []string{})
	require.NoError(t, err)
	require.True(t, txid == "f0ace3a53686463bda143371613ac204780f0fddb6f19fe1cc88c2f0f2588541", "invalid calculated txid")
	require.True(t, len(outputs) == 2, "should be two outputs") //in this eg. 1 is real, 1 has 0 amount, no address
	fmt.Println(outputs)
}

func Test_CheckBlockHeader(t *testing.T) {
	t.Skip("Skipping on CI - requires Bitcoin block data validation")

	//Given this Blockheader
	// Check the blockhash is correctly derived from the other data
	// Check the difficult for the blockhash (leading zero's) are within permitted values (bits fields)
	blockHeader := &api.BTCBlockHeader{
		Version:    536870912,
		PrevBlock:  "000000000000000f10b5de36d015586d3bf3f63a0faa418b73cb91aaff5de064",
		MerkleRoot: "503d9d9abd29b58233568efded9af1516f395fea7ca3eb1aaafdc8b63406d214",
		TimeStamp:  1725270640,
		Bits:       486604799,
		Nonce:      3645002625,
		BlockHash:  "000000003c057188e715529359ae9a45a9f7b8462f0c9b8b9f2f24ac196eb6fb",
	}

	err := CheckBlockHeader(blockHeader)
	require.NoError(t, err)

}

func Test_VerboseVerifyBTCLockTransaction(t *testing.T) {

	//chainName := "testnet"
	//blockhash := "000000003c057188e715529359ae9a45a9f7b8462f0c9b8b9f2f24ac196eb6fb"
	txid := "f0ace3a53686463bda143371613ac204780f0fddb6f19fe1cc88c2f0f2588541"
	index := 3462
	//blockHeight := 2902385

	rawTX := "0200000001252ef4787fe208a370fecf707674605ca4c803036e636ebdcbc0a287ca12f74c000000006b483045022100caa45533d411ce7140d6e9e9a0a1d00947304d09e546d1c725a2ef5679f11ba202201fc19ec89b4b2def1119fd0d95853d9f2e2f118867cf8dbbb60e5c8a793bd1d20121021848a8e26bf13437b711cf2504083650722a98a15d00283575260737cb4c050dffffffff0275da4c04000000001976a91424299bb977941a79afba01728c34d018ae02c7c888ac0000000000000000536a4c5048454d49010064f9080015ce6c8c4ea762b9a5fccc0caaf96953afa9ab4d3bb24ae2665ab175216fff8f3525c402168d60e83afc7f1ba011da7f7857612c4f35790dba5b93840c1db1f8945ae5ec9b6569492c00"

	proof := []string{
		"45f55616ce67e7eb8c92fb0b25dd342f61e1d8906f300e372e64260e464b7c03",
		"0a184d91e652c35d987072d5f1ce33c77a411b3f470cf185097d035a275f0e1c",
		"ba6526507eb6601eb5ab1113aca7429d888e8fc1e77562766eb2d554924ad489",
		"c1b30eaea97e20d61c2311260811080f9f35daf57f134bffb167eef0bfbd8337",
		"e679739ab8a223122552cf32a94bcfaa43492258d15dd5f08bca00d472d64d89",
		"1f15e67691be2b5e6de6384913e2a98bac62690778d7d43e26a5ef8cf543ce80",
		"a7ffdc429d24c3e26a4d54e94e3af88db6507b6d615ef4dd720e91982e4ce02c",
		"c20c42d2d3dfacc90e35d38cd328e73e1e1c76596d4b11af0bc9823480d0e6ed",
		"ed246e78bd59bf29d9d3b912e872fbf87f8b9f1757b9917d99644fe7dcc717c0",
		"3488bb935d647d178d7bdfa442a5cd02923693fadca295c0bfcc33b184ee07d2",
		"554f43d6696d786092c4017d60b3a3ff1b7ae92911e9d702a2af66c52bc519f3",
		"c6b371722f132ee15c2dc09d4c14b40149e0507a73d89f31b5abf6bdaa702694",
	}

	//This data comes from the Neutrino Node
	blockHeader := &api.BTCBlockHeader{
		Version:    536870912,
		PrevBlock:  "000000000000000f10b5de36d015586d3bf3f63a0faa418b73cb91aaff5de064",
		MerkleRoot: "503d9d9abd29b58233568efded9af1516f395fea7ca3eb1aaafdc8b63406d214",
		TimeStamp:  1725270640,
		Bits:       486604799,
		Nonce:      3645002625,
		BlockHash:  "000000003c057188e715529359ae9a45a9f7b8462f0c9b8b9f2f24ac196eb6fb",
	}
	_ = blockHeader

	//Check The Transaction ID is correctly derived
	calculatedTxID, err := CalculateTXID(rawTX, "testnet")
	require.NoError(t, err)
	calculatedIDString := ReverseHex(calculatedTxID.String())
	require.Equal(t, txid, calculatedIDString, "txid should be equal")

	//Check the Merkle Proof, derive the merkle root using the proof and compare to the actual MerkleRoot
	targetHash, err := chainhash.NewHashFromStr(calculatedIDString)
	require.NoError(t, err)
	i := index
	for _, sibling := range proof {
		siblingHash, err := chainhash.NewHashFromStr(sibling)
		require.NoError(t, err)
		if i%2 == 0 {
			targetHash = MergeHashes(targetHash, siblingHash)
		} else {
			targetHash = MergeHashes(siblingHash, targetHash)
		}
		i /= 2
	}
	merkleRootBytes, err := chainhash.NewHashFromStr(blockHeader.MerkleRoot)
	require.NoError(t, err)
	require.True(t, targetHash.IsEqual(merkleRootBytes), "merkle root should be equal")

	//TODO Check the Block Hash by rehashing the block

}
