package main

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/neutrino"
	"github.com/stretchr/testify/require"
)

func Test_ProxyFunctions_Testnet3(t *testing.T) {
	t.Skip("manual test that requires local bitcoin proxy and may take time to sync; skipped in automated runs")
	//Do not run as part of CI.

	//This test requires 2 additional pieces of setup, so will not work out of the box
	/*
		1. The 30 second sleep below gives the Neutrino nodes time to sync, an error will be returned unless the node
		has sufficient time to download the block headers, once the database has been built the sleep can be removed if further
		testing is run.
		2. A Bitcoin Proxy needs to be running locally "http://127.0.0.1:1234", "user", "secret" so the alternative mechanism
		of obtaining a block header can be tested.
	*/

	ns := neutrino.NeutrinoServer{}
	ns.Initialize("testnet", "http://127.0.0.1:1234", "user", "secret", "./neutrino", 12345, "/neutrino_")

	time.Sleep(30 * time.Second)

	//Get via the Neutrino Node
	header1, hash1, _, err := ns.GetBlockHeaderByHeight("testnet3", 1000)
	require.Nil(t, err, "error getting block header")

	//Get via the Proxy
	header2, hash2, _, err := ns.ProxyGetBlockHeaderByHeight("testnet3", 1000)
	require.Nil(t, err, "error getting block header")

	//ignore height it takes a while for it to build the neutrino filters

	require.Equal(t, header1.Nonce, header2.Nonce, "block header mismatch")
	require.Equal(t, header1.Version, header2.Version, "block header mismatch")
	require.Equal(t, header1.Timestamp, header2.Timestamp, "block header mismatch")
	require.Equal(t, header1.Bits, header2.Bits, "block header mismatch")
	require.Equal(t, header1.MerkleRoot, header2.MerkleRoot, "block header mismatch")
	require.Equal(t, header1.PrevBlock, header2.PrevBlock, "block header mismatch")

	require.Equal(t, hash1, hash2, "block hash mismatch")
}

func Test_ProxyFunctions_Testnet4(t *testing.T) {
	t.Skip("manual test that requires local bitcoin proxy and may take time to sync; skipped in automated runs")
	//Do not run as part of CI.

	//This test requires 2 additional pieces of setup, so will not work out of the box
	/*
		1. The 30 second sleep below gives the Neutrino nodes time to sync, an error will be returned unless the node
		has sufficient time to download the block headers, once the database has been built the sleep can be removed if further
		testing is run.
		2. A Bitcoin Proxy needs to be running locally "http://127.0.0.1:1234", "user", "secret" so the alternative mechanism
		of obtaining a block header can be tested.
	*/

	ns := neutrino.NeutrinoServer{}
	ns.Initialize("testnet", "http://127.0.0.1:1234", "user", "secret", "./neutrino", 12345, "/neutrino_")
	time.Sleep(30 * time.Second)

	//Get via the Neutrino Node
	_, _, _, err := ns.GetBlockHeaderByHeight("testnet4", 1000)
	require.Error(t, err, "Testnet4 should return an error")

	//Get via the Proxy
	_, hash2, _, err := ns.ProxyGetBlockHeaderByHeight("testnet4", 1000)
	require.Nil(t, err, "error getting block header")

	hex := hex.EncodeToString(hash2[:])
	hex = ReverseHex(hex)
	require.Equal(t, hex, "00000000b747d47c3b38161693ad05e26924b3775a8be669751f969da836311e", "hash is invalid")

}
func ReverseHex(hexStr string) string {
	n := len(hexStr)
	if n%2 != 0 {
		// Ensure the hex string length is even
		return "invalid hex string"
	}
	result := make([]byte, n)
	for i := 0; i < n; i += 2 {
		// Copy two characters (one byte) at a time from the end to the beginning
		result[n-i-2], result[n-i-1] = hexStr[i], hexStr[i+1]
	}
	return string(result)
}
