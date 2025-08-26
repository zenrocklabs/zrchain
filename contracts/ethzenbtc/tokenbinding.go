// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zenbtc

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IZenBTCSidecarHelperZenBTCBurn is an auto generated low-level Go binding around an user-defined struct.
type IZenBTCSidecarHelperZenBTCBurn struct {
	NetValue    *big.Int
	BlockNumber uint64
	Sender      common.Address
	Receiver    []byte
}

// ZenbtcMetaData contains all meta data concerning the Zenbtc contract.
var ZenbtcMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"basisPointsRateUpdated\",\"type\":\"uint16\"}],\"name\":\"BasisPointsRateUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newFeeAddress\",\"type\":\"address\"}],\"name\":\"FeeAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newRedeemFee\",\"type\":\"uint256\"}],\"name\":\"RedeemFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldHelper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newHelper\",\"type\":\"address\"}],\"name\":\"SidecarHelperUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"destAddr\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"TokenRedemption\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"}],\"name\":\"TokensMintedWithFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"zenBTControllerUpdated\",\"type\":\"address\"}],\"name\":\"ZenBTControllerUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentVersion\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"estimateFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllBurns\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"netValue\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"}],\"internalType\":\"structIZenBTCSidecarHelper.ZenBTCBurn[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBasisPointsRate\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"netValue\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"}],\"internalType\":\"structIZenBTCSidecarHelper.ZenBTCBurn\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBurnCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"from\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"to\",\"type\":\"uint256\"}],\"name\":\"getBurnsFromTo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"netValue\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"}],\"internalType\":\"structIZenBTCSidecarHelper.ZenBTCBurn[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSidecarHelper\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"}],\"name\":\"initializeV1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"}],\"name\":\"initializeV2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sidecarHelper_\",\"type\":\"address\"}],\"name\":\"initializeV3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destAddr\",\"type\":\"bytes\"}],\"name\":\"unwrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"basisPointsRate_\",\"type\":\"uint16\"}],\"name\":\"updateBasisPointsRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeAddress_\",\"type\":\"address\"}],\"name\":\"updateFeeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newHelper\",\"type\":\"address\"}],\"name\":\"updateSidecarHelper\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"}],\"name\":\"wrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b50620000216200002760201b60201c565b62000191565b5f620000386200012b60201b60201c565b9050805f0160089054906101000a900460ff161562000083576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8016815f015f9054906101000a900467ffffffffffffffff1667ffffffffffffffff1614620001285767ffffffffffffffff815f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d267ffffffffffffffff6040516200011f919062000176565b60405180910390a15b50565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b5f67ffffffffffffffff82169050919050565b620001708162000152565b82525050565b5f6020820190506200018b5f83018462000165565b92915050565b6142a4806200019f5f395ff3fe608060405234801561000f575f80fd5b5060043610610230575f3560e01c80635acc22c71161012e578063bbcaac38116100b6578063e7cf548c1161007a578063e7cf548c146106c2578063f1d96dd3146106e0578063faa85988146106fc578063fd154cc714610718578063fded25291461073457610230565b8063bbcaac3814610620578063d53913931461063c578063d547741f1461065a578063dd62ed3e14610676578063deb42f13146106a657610230565b806395d89b41116100fd57806395d89b411461057a5780639d888e8614610598578063a217fddf146105b6578063a9059cbb146105d4578063b413148e1461060457610230565b80635acc22c7146104ce57806370a08231146104fe57806379cc67901461052e57806391d148541461054a57610230565b80632a4b2fd9116101bc578063313ce56711610180578063313ce5671461043e57806336568abe1461045c57806340c10f191461047857806342966c68146104945780634e7ceacb146104b057610230565b80632a4b2fd91461039c5780632b9b5728146103b85780632eb3f49b146103d65780632f2ff15d146104065780633101cfcb1461042257610230565b8063127e8e4d11610203578063127e8e4d146102d057806318160ddd1461030057806323b872dd1461031e578063248a9ca31461034e5780632a0276f81461037e57610230565b806301ffc9a71461023457806303b41f741461026457806306fdde0314610282578063095ea7b3146102a0575b5f80fd5b61024e60048036038101906102499190612ba3565b610752565b60405161025b9190612be8565b60405180910390f35b61026c6107cb565b6040516102799190612e1f565b60405180910390f35b61028a610900565b6040516102979190612e91565b60405180910390f35b6102ba60048036038101906102b59190612f05565b61099e565b6040516102c79190612be8565b60405180910390f35b6102ea60048036038101906102e59190612f43565b6109c0565b6040516102f79190612f7d565b60405180910390f35b6103086109d1565b6040516103159190612f7d565b60405180910390f35b61033860048036038101906103339190612f96565b6109e8565b6040516103459190612be8565b60405180910390f35b61036860048036038101906103639190613019565b610a16565b6040516103759190613053565b60405180910390f35b610386610a40565b6040516103939190613053565b60405180910390f35b6103b660048036038101906103b191906131ce565b610a66565b005b6103c0610beb565b6040516103cd9190613265565b60405180910390f35b6103f060048036038101906103eb9190612f43565b610c20565b6040516103fd91906132de565b60405180910390f35b610420600480360381019061041b91906132fe565b610d68565b005b61043c6004803603810190610437919061333c565b610d8a565b005b610446610f6b565b6040516104539190613376565b60405180910390f35b610476600480360381019061047191906132fe565b610f8f565b005b610492600480360381019061048d9190612f05565b61100a565b005b6104ae60048036038101906104a99190612f43565b611045565b005b6104b8611086565b6040516104c59190613265565b60405180910390f35b6104e860048036038101906104e3919061338f565b6110bc565b6040516104f59190612e1f565b60405180910390f35b6105186004803603810190610513919061333c565b611201565b6040516105259190612f7d565b60405180910390f35b61054860048036038101906105439190612f05565b611254565b005b610564600480360381019061055f91906132fe565b6112a1565b6040516105719190612be8565b60405180910390f35b610582611312565b60405161058f9190612e91565b60405180910390f35b6105a06113b0565b6040516105ad91906133dc565b60405180910390f35b6105be6113be565b6040516105cb9190613053565b60405180910390f35b6105ee60048036038101906105e99190612f05565b6113c4565b6040516105fb9190612be8565b60405180910390f35b61061e60048036038101906106199190613493565b6113e6565b005b61063a6004803603810190610635919061333c565b61151a565b005b610644611553565b6040516106519190613053565b60405180910390f35b610674600480360381019061066f91906132fe565b611579565b005b610690600480360381019061068b91906134ed565b61159b565b60405161069d9190612f7d565b60405180910390f35b6106c060048036038101906106bb91906131ce565b61162b565b005b6106ca61175e565b6040516106d79190612f7d565b60405180910390f35b6106fa60048036038101906106f59190613555565b61188f565b005b6107166004803603810190610711919061333c565b611a25565b005b610732600480360381019061072d91906135dc565b611b71565b005b61073c611baa565b6040516107499190613616565b60405180910390f35b5f7f7965db0b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806107c457506107c382611bce565b5b9050919050565b60605f6107d6611c37565b90505f73ffffffffffffffffffffffffffffffffffffffff16815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610868576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161085f90613679565b60405180910390fd5b805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166303b41f746040518163ffffffff1660e01b81526004015f60405180830381865afa1580156108d2573d5f803e3d5ffd5b505050506040513d5f823e3d601f19601f820116820180604052508101906108fa91906138bc565b91505090565b60605f61090b611c9d565b905080600401805461091c90613930565b80601f016020809104026020016040519081016040528092919081815260200182805461094890613930565b80156109935780601f1061096a57610100808354040283529160200191610993565b820191905f5260205f20905b81548152906001019060200180831161097657829003601f168201915b505050505091505090565b5f806109a8611cc4565b90506109b5818585611ccb565b600191505092915050565b5f6109ca82611cdd565b9050919050565b5f806109db611c9d565b9050806002015491505090565b5f806109f2611cc4565b90506109ff858285611d9b565b610a0a858585611e2d565b60019150509392505050565b5f80610a20611f1d565b9050805f015f8481526020019081526020015f2060010154915050919050565b7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf285f1b81565b5f610a6f611f44565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff16148015610ab75750825b90505f60018367ffffffffffffffff16148015610aea57505f3073ffffffffffffffffffffffffffffffffffffffff163b145b905081158015610af8575080155b15610b2f576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508315610b7c576001855f0160086101000a81548160ff0219169083151502179055505b610b87888888611f6b565b8315610be1575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d26001604051610bd891906139a2565b60405180910390a15b5050505050505050565b5f80610bf5611c37565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b610c28612af8565b5f610c31611c37565b90505f73ffffffffffffffffffffffffffffffffffffffff16815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610cc3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cba90613679565b60405180910390fd5b805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632eb3f49b846040518263ffffffff1660e01b8152600401610d1e9190612f7d565b5f60405180830381865afa158015610d38573d5f803e3d5ffd5b505050506040513d5f823e3d601f19601f82011682018060405250810190610d6091906139bb565b915050919050565b610d7182610a16565b610d7a81611f93565b610d848383611fa7565b50505050565b60035f610d95611f44565b9050805f0160089054906101000a900460ff1680610ddd57508167ffffffffffffffff16815f015f9054906101000a900467ffffffffffffffff1667ffffffffffffffff1610155b15610e14576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81815f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506001815f0160086101000a81548160ff0219169083151502179055505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610ec7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ebe90613a72565b60405180910390fd5b5f610ed0611c37565b905083815f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505f815f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d282604051610f5e91906133dc565b60405180910390a1505050565b5f80610f75611c9d565b90508060030160029054906101000a900460ff1691505090565b610f97611cc4565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610ffb576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611005828261209f565b505050565b7fb458d13b0f4ce9e6aa65d297a27b10f75fdc6d0957bb29e1f2a30c8766b354155f1b61103681611f93565b6110408383612197565b505050565b7fb458d13b0f4ce9e6aa65d297a27b10f75fdc6d0957bb29e1f2a30c8766b354155f1b61107181611f93565b61108261107c611cc4565b83612216565b5050565b5f80611090611c9d565b9050806006015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b60605f6110c7611c37565b90505f73ffffffffffffffffffffffffffffffffffffffff16815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603611159576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161115090613679565b60405180910390fd5b805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635acc22c785856040518363ffffffff1660e01b81526004016111b6929190613a90565b5f60405180830381865afa1580156111d0573d5f803e3d5ffd5b505050506040513d5f823e3d601f19601f820116820180604052508101906111f891906138bc565b91505092915050565b5f8061120b611c9d565b9050805f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054915050919050565b7fb458d13b0f4ce9e6aa65d297a27b10f75fdc6d0957bb29e1f2a30c8766b354155f1b61128081611f93565b6112928361128c611cc4565b84611d9b565b61129c8383612216565b505050565b5f806112ab611f1d565b9050805f015f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1691505092915050565b60605f61131d611c9d565b905080600501805461132e90613930565b80601f016020809104026020016040519081016040528092919081815260200182805461135a90613930565b80156113a55780601f1061137c576101008083540402835291602001916113a5565b820191905f5260205f20905b81548152906001019060200180831161138857829003601f168201915b505050505091505090565b5f6113b9612295565b905090565b5f801b81565b5f806113ce611cc4565b90506113db818585611e2d565b600191505092915050565b6113f082826122b9565b5f6113f9611c37565b90505f73ffffffffffffffffffffffffffffffffffffffff16815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614611515575f61145a84611cdd565b90505f81856114699190613ae4565b9050825f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166398e886a1826114b2611cc4565b876040518463ffffffff1660e01b81526004016114d193929190613b5f565b6020604051808303815f875af11580156114ed573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906115119190613b9b565b5050505b505050565b7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf285f1b61154681611f93565b61154f826123ca565b5050565b7fb458d13b0f4ce9e6aa65d297a27b10f75fdc6d0957bb29e1f2a30c8766b354155f1b81565b61158282610a16565b61158b81611f93565b611595838361209f565b50505050565b5f806115a5611c9d565b9050806001015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205491505092915050565b60025f611636611f44565b9050805f0160089054906101000a900460ff168061167e57508167ffffffffffffffff16815f015f9054906101000a900467ffffffffffffffff1667ffffffffffffffff1610155b156116b5576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81815f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506001815f0160086101000a81548160ff02191690831515021790555061170585858561245e565b5f815f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d28260405161174f91906133dc565b60405180910390a15050505050565b5f80611768611c37565b90505f73ffffffffffffffffffffffffffffffffffffffff16815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036117fa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016117f190613679565b60405180910390fd5b805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e7cf548c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611865573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906118899190613b9b565b91505090565b7fb458d13b0f4ce9e6aa65d297a27b10f75fdc6d0957bb29e1f2a30c8766b354155f1b6118bb81611f93565b67ffffffffffffffff80168367ffffffffffffffff161115611912576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161190990613c10565b60405180910390fd5b8267ffffffffffffffff168267ffffffffffffffff1610611968576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161195f90613c9e565b60405180910390fd5b5f611971611c9d565b90505f83856119809190613cbc565b90506119b9826006015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff168567ffffffffffffffff16612197565b6119cd868267ffffffffffffffff16612197565b8573ffffffffffffffffffffffffffffffffffffffff167f890f2fd0578a1ba6f8735fc60c94d1ec453ad71501d82ba4d5a9041f9ac14f2e8286604051611a15929190613cf7565b60405180910390a2505050505050565b5f801b611a3181611f93565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611a9f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a9690613a72565b60405180910390fd5b5f611aa8611c37565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905083825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167fae0961ff5b393c8e05811c7cccbbbbb4453f9336fcd0197d1d87ee17072bdbcf60405160405180910390a350505050565b7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf285f1b611b9d81611f93565b611ba6826124b8565b5050565b5f80611bb4611c9d565b9050806003015f9054906101000a900461ffff1691505090565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b5f8060ff5f1b1960017fdb8ca13531a22b74b6340685cd30f2b2fde811711771bb2bdc9faaa2932826d55f1c611c6d9190613ae4565b604051602001611c7d9190612f7d565b604051602081830303815290604052805190602001201690508091505090565b5f7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00905090565b5f33905090565b611cd88383836001612569565b505050565b5f61271061ffff168211611d26576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611d1d90613d68565b60405180910390fd5b5f611d2f611c9d565b90505f611d4d84836003015f9054906101000a900461ffff16612746565b9050838110611d91576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611d8890613df6565b60405180910390fd5b8092505050919050565b5f611da6848461159b565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114611e275781811015611e18578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401611e0f93929190613e14565b60405180910390fd5b611e2684848484035f612569565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603611e9d575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401611e949190613265565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611f0d575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401611f049190613265565b60405180910390fd5b611f18838383612770565b505050565b5f7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800905090565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b611f73612a0f565b611f7b612a4f565b611f83612a59565b611f8e83838361245e565b505050565b611fa481611f9f611cc4565b612a6b565b50565b5f80611fb1611f1d565b9050611fbd84846112a1565b612094576001815f015f8681526020019081526020015f205f015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550612030611cc4565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a46001915050612099565b5f9150505b92915050565b5f806120a9611f1d565b90506120b584846112a1565b1561218c575f815f015f8681526020019081526020015f205f015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550612128611cc4565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a46001915050612191565b5f9150505b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603612207575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016121fe9190613265565b60405180910390fd5b6122125f8383612770565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603612286575f6040517f96c6fd1e00000000000000000000000000000000000000000000000000000000815260040161227d9190613265565b60405180910390fd5b612291825f83612770565b5050565b5f61229e611f44565b5f015f9054906101000a900467ffffffffffffffff16905090565b5f6122c2611c9d565b905067ffffffffffffffff8016831115612311576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161230890613c10565b60405180910390fd5b5f61231b84611cdd565b905061232e612328611cc4565b85612216565b61235b826006015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682612197565b5f81856123689190613ae4565b9050612372611cc4565b73ffffffffffffffffffffffffffffffffffffffff167f4c971c8b2abb197a17896b2fc57f597830db78d3556831e3faa337a596150f228286856040516123bb93929190613e49565b60405180910390a25050505050565b5f6123d3611c9d565b905081816006015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff167f446e39bcf1b47cfadfaa23442cb4b34682cfe6bd9220da084894e3b1f834e4f360405160405180910390a25050565b612466612a0f565b5f61246f611c9d565b9050838160040190816124829190614019565b50828160050190816124949190614019565b50818160030160026101000a81548160ff021916908360ff16021790555050505050565b61271061ffff168161ffff161115612505576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016124fc90614158565b60405180910390fd5b5f61250e611c9d565b905081816003015f6101000a81548161ffff021916908361ffff1602179055507f79d2a06c43b232cf9d1835b5e915efe74621561aaf75fecec91fd75a940f9d708260405161255d9190613616565b60405180910390a15050565b5f612572611c9d565b90505f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16036125e4575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016125db9190613265565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603612654575f6040517f94280d6200000000000000000000000000000000000000000000000000000000815260040161264b9190613265565b60405180910390fd5b82816001015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550811561273f578373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925856040516127369190612f7d565b60405180910390a35b5050505050565b5f61271061ffff168261ffff168461275e9190614176565b61276891906141e4565b905092915050565b3073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036127e057306040517fec442f050000000000000000000000000000000000000000000000000000000081526004016127d79190613265565b60405180910390fd5b5f6127e9611c9d565b90505f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff160361283d5781816002015f8282546128319190614214565b9250508190555061290f565b5f815f015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050828110156128c8578481846040517fe450d38c0000000000000000000000000000000000000000000000000000000081526004016128bf93929190613e14565b60405180910390fd5b828103825f015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036129585781816002015f82825403925050819055506129a4565b81815f015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051612a019190612f7d565b60405180910390a350505050565b612a17612abc565b612a4d576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b612a57612a0f565b565b612a61612a0f565b612a69612ada565b565b612a7582826112a1565b612ab85780826040517fe2517d3f000000000000000000000000000000000000000000000000000000008152600401612aaf929190614247565b60405180910390fd5b5050565b5f612ac5611f44565b5f0160089054906101000a900460ff16905090565b612ae2612a0f565b612af55f801b612af0611cc4565b611fa7565b50565b60405180608001604052805f81526020015f67ffffffffffffffff1681526020015f73ffffffffffffffffffffffffffffffffffffffff168152602001606081525090565b5f604051905090565b5f80fd5b5f80fd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b612b8281612b4e565b8114612b8c575f80fd5b50565b5f81359050612b9d81612b79565b92915050565b5f60208284031215612bb857612bb7612b46565b5b5f612bc584828501612b8f565b91505092915050565b5f8115159050919050565b612be281612bce565b82525050565b5f602082019050612bfb5f830184612bd9565b92915050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b5f819050919050565b612c3c81612c2a565b82525050565b5f67ffffffffffffffff82169050919050565b612c5e81612c42565b82525050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f612c8d82612c64565b9050919050565b612c9d81612c83565b82525050565b5f81519050919050565b5f82825260208201905092915050565b5f5b83811015612cda578082015181840152602081019050612cbf565b5f8484015250505050565b5f601f19601f8301169050919050565b5f612cff82612ca3565b612d098185612cad565b9350612d19818560208601612cbd565b612d2281612ce5565b840191505092915050565b5f608083015f830151612d425f860182612c33565b506020830151612d556020860182612c55565b506040830151612d686040860182612c94565b5060608301518482036060860152612d808282612cf5565b9150508091505092915050565b5f612d988383612d2d565b905092915050565b5f602082019050919050565b5f612db682612c01565b612dc08185612c0b565b935083602082028501612dd285612c1b565b805f5b85811015612e0d5784840389528151612dee8582612d8d565b9450612df983612da0565b925060208a01995050600181019050612dd5565b50829750879550505050505092915050565b5f6020820190508181035f830152612e378184612dac565b905092915050565b5f81519050919050565b5f82825260208201905092915050565b5f612e6382612e3f565b612e6d8185612e49565b9350612e7d818560208601612cbd565b612e8681612ce5565b840191505092915050565b5f6020820190508181035f830152612ea98184612e59565b905092915050565b612eba81612c83565b8114612ec4575f80fd5b50565b5f81359050612ed581612eb1565b92915050565b612ee481612c2a565b8114612eee575f80fd5b50565b5f81359050612eff81612edb565b92915050565b5f8060408385031215612f1b57612f1a612b46565b5b5f612f2885828601612ec7565b9250506020612f3985828601612ef1565b9150509250929050565b5f60208284031215612f5857612f57612b46565b5b5f612f6584828501612ef1565b91505092915050565b612f7781612c2a565b82525050565b5f602082019050612f905f830184612f6e565b92915050565b5f805f60608486031215612fad57612fac612b46565b5b5f612fba86828701612ec7565b9350506020612fcb86828701612ec7565b9250506040612fdc86828701612ef1565b9150509250925092565b5f819050919050565b612ff881612fe6565b8114613002575f80fd5b50565b5f8135905061301381612fef565b92915050565b5f6020828403121561302e5761302d612b46565b5b5f61303b84828501613005565b91505092915050565b61304d81612fe6565b82525050565b5f6020820190506130665f830184613044565b92915050565b5f80fd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6130aa82612ce5565b810181811067ffffffffffffffff821117156130c9576130c8613074565b5b80604052505050565b5f6130db612b3d565b90506130e782826130a1565b919050565b5f67ffffffffffffffff82111561310657613105613074565b5b61310f82612ce5565b9050602081019050919050565b828183375f83830152505050565b5f61313c613137846130ec565b6130d2565b90508281526020810184848401111561315857613157613070565b5b61316384828561311c565b509392505050565b5f82601f83011261317f5761317e61306c565b5b813561318f84826020860161312a565b91505092915050565b5f60ff82169050919050565b6131ad81613198565b81146131b7575f80fd5b50565b5f813590506131c8816131a4565b92915050565b5f805f606084860312156131e5576131e4612b46565b5b5f84013567ffffffffffffffff81111561320257613201612b4a565b5b61320e8682870161316b565b935050602084013567ffffffffffffffff81111561322f5761322e612b4a565b5b61323b8682870161316b565b925050604061324c868287016131ba565b9150509250925092565b61325f81612c83565b82525050565b5f6020820190506132785f830184613256565b92915050565b5f608083015f8301516132935f860182612c33565b5060208301516132a66020860182612c55565b5060408301516132b96040860182612c94565b50606083015184820360608601526132d18282612cf5565b9150508091505092915050565b5f6020820190508181035f8301526132f6818461327e565b905092915050565b5f806040838503121561331457613313612b46565b5b5f61332185828601613005565b925050602061333285828601612ec7565b9150509250929050565b5f6020828403121561335157613350612b46565b5b5f61335e84828501612ec7565b91505092915050565b61337081613198565b82525050565b5f6020820190506133895f830184613367565b92915050565b5f80604083850312156133a5576133a4612b46565b5b5f6133b285828601612ef1565b92505060206133c385828601612ef1565b9150509250929050565b6133d681612c42565b82525050565b5f6020820190506133ef5f8301846133cd565b92915050565b5f67ffffffffffffffff82111561340f5761340e613074565b5b61341882612ce5565b9050602081019050919050565b5f613437613432846133f5565b6130d2565b90508281526020810184848401111561345357613452613070565b5b61345e84828561311c565b509392505050565b5f82601f83011261347a5761347961306c565b5b813561348a848260208601613425565b91505092915050565b5f80604083850312156134a9576134a8612b46565b5b5f6134b685828601612ef1565b925050602083013567ffffffffffffffff8111156134d7576134d6612b4a565b5b6134e385828601613466565b9150509250929050565b5f806040838503121561350357613502612b46565b5b5f61351085828601612ec7565b925050602061352185828601612ec7565b9150509250929050565b61353481612c42565b811461353e575f80fd5b50565b5f8135905061354f8161352b565b92915050565b5f805f6060848603121561356c5761356b612b46565b5b5f61357986828701612ec7565b935050602061358a86828701613541565b925050604061359b86828701613541565b9150509250925092565b5f61ffff82169050919050565b6135bb816135a5565b81146135c5575f80fd5b50565b5f813590506135d6816135b2565b92915050565b5f602082840312156135f1576135f0612b46565b5b5f6135fe848285016135c8565b91505092915050565b613610816135a5565b82525050565b5f6020820190506136295f830184613607565b92915050565b7f536964656361722068656c706572206e6f7420736574000000000000000000005f82015250565b5f613663601683612e49565b915061366e8261362f565b602082019050919050565b5f6020820190508181035f83015261369081613657565b9050919050565b5f67ffffffffffffffff8211156136b1576136b0613074565b5b602082029050602081019050919050565b5f80fd5b5f80fd5b5f80fd5b5f815190506136dc81612edb565b92915050565b5f815190506136f08161352b565b92915050565b5f8151905061370481612eb1565b92915050565b5f61371c613717846133f5565b6130d2565b90508281526020810184848401111561373857613737613070565b5b613743848285612cbd565b509392505050565b5f82601f83011261375f5761375e61306c565b5b815161376f84826020860161370a565b91505092915050565b5f6080828403121561378d5761378c6136c6565b5b61379760806130d2565b90505f6137a6848285016136ce565b5f8301525060206137b9848285016136e2565b60208301525060406137cd848285016136f6565b604083015250606082015167ffffffffffffffff8111156137f1576137f06136ca565b5b6137fd8482850161374b565b60608301525092915050565b5f61381b61381684613697565b6130d2565b9050808382526020820190506020840283018581111561383e5761383d6136c2565b5b835b8181101561388557805167ffffffffffffffff8111156138635761386261306c565b5b8086016138708982613778565b85526020850194505050602081019050613840565b5050509392505050565b5f82601f8301126138a3576138a261306c565b5b81516138b3848260208601613809565b91505092915050565b5f602082840312156138d1576138d0612b46565b5b5f82015167ffffffffffffffff8111156138ee576138ed612b4a565b5b6138fa8482850161388f565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061394757607f821691505b60208210810361395a57613959613903565b5b50919050565b5f819050919050565b5f819050919050565b5f61398c61398761398284613960565b613969565b612c42565b9050919050565b61399c81613972565b82525050565b5f6020820190506139b55f830184613993565b92915050565b5f602082840312156139d0576139cf612b46565b5b5f82015167ffffffffffffffff8111156139ed576139ec612b4a565b5b6139f984828501613778565b91505092915050565b7f536964656361722068656c7065722063616e6e6f74206265207a65726f2061645f8201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b5f613a5c602583612e49565b9150613a6782613a02565b604082019050919050565b5f6020820190508181035f830152613a8981613a50565b9050919050565b5f604082019050613aa35f830185612f6e565b613ab06020830184612f6e565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f613aee82612c2a565b9150613af983612c2a565b9250828203905081811115613b1157613b10613ab7565b5b92915050565b5f82825260208201905092915050565b5f613b3182612ca3565b613b3b8185613b17565b9350613b4b818560208601612cbd565b613b5481612ce5565b840191505092915050565b5f606082019050613b725f830186612f6e565b613b7f6020830185613256565b8181036040830152613b918184613b27565b9050949350505050565b5f60208284031215613bb057613baf612b46565b5b5f613bbd848285016136ce565b91505092915050565b7f56616c756520746f6f206c6172676500000000000000000000000000000000005f82015250565b5f613bfa600f83612e49565b9150613c0582613bc6565b602082019050919050565b5f6020820190508181035f830152613c2781613bee565b9050919050565b7f5a656e4254433a2046656520657863656564732074686520616d6f756e7420745f8201527f6f206d696e740000000000000000000000000000000000000000000000000000602082015250565b5f613c88602683612e49565b9150613c9382613c2e565b604082019050919050565b5f6020820190508181035f830152613cb581613c7c565b9050919050565b5f613cc682612c42565b9150613cd183612c42565b9250828203905067ffffffffffffffff811115613cf157613cf0613ab7565b5b92915050565b5f604082019050613d0a5f8301856133cd565b613d1760208301846133cd565b9392505050565b7f56616c756520746f6f20736d616c6c00000000000000000000000000000000005f82015250565b5f613d52600f83612e49565b9150613d5d82613d1e565b602082019050919050565b5f6020820190508181035f830152613d7f81613d46565b9050919050565b7f5a656e4254433a2072656465656d2066656520657863656564732074686520615f8201527f6d6f756e74000000000000000000000000000000000000000000000000000000602082015250565b5f613de0602583612e49565b9150613deb82613d86565b604082019050919050565b5f6020820190508181035f830152613e0d81613dd4565b9050919050565b5f606082019050613e275f830186613256565b613e346020830185612f6e565b613e416040830184612f6e565b949350505050565b5f606082019050613e5c5f8301866133cd565b8181036020830152613e6e8185613b27565b9050613e7d6040830184612f6e565b949350505050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302613ee17fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82613ea6565b613eeb8683613ea6565b95508019841693508086168417925050509392505050565b5f613f1d613f18613f1384612c2a565b613969565b612c2a565b9050919050565b5f819050919050565b613f3683613f03565b613f4a613f4282613f24565b848454613eb2565b825550505050565b5f90565b613f5e613f52565b613f69818484613f2d565b505050565b5b81811015613f8c57613f815f82613f56565b600181019050613f6f565b5050565b601f821115613fd157613fa281613e85565b613fab84613e97565b81016020851015613fba578190505b613fce613fc685613e97565b830182613f6e565b50505b505050565b5f82821c905092915050565b5f613ff15f1984600802613fd6565b1980831691505092915050565b5f6140098383613fe2565b9150826002028217905092915050565b61402282612e3f565b67ffffffffffffffff81111561403b5761403a613074565b5b6140458254613930565b614050828285613f90565b5f60209050601f831160018114614081575f841561406f578287015190505b6140798582613ffe565b8655506140e0565b601f19841661408f86613e85565b5f5b828110156140b657848901518255600182019150602085019450602081019050614091565b868310156140d357848901516140cf601f891682613fe2565b8355505b6001600288020188555050505b505050505050565b7f626173697320706f696e747320726174652063616e6e6f7420657863656564205f8201527f31303030302075696e7473000000000000000000000000000000000000000000602082015250565b5f614142602b83612e49565b915061414d826140e8565b604082019050919050565b5f6020820190508181035f83015261416f81614136565b9050919050565b5f61418082612c2a565b915061418b83612c2a565b925082820261419981612c2a565b915082820484148315176141b0576141af613ab7565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f6141ee82612c2a565b91506141f983612c2a565b925082614209576142086141b7565b5b828204905092915050565b5f61421e82612c2a565b915061422983612c2a565b925082820190508082111561424157614240613ab7565b5b92915050565b5f60408201905061425a5f830185613256565b6142676020830184613044565b939250505056fea26469706673582212208593e8106611419e71aed7148704da8aa81f23fd31a87f84f959ce5169008a7b64736f6c63430008180033",
}

// ZenbtcABI is the input ABI used to generate the binding from.
// Deprecated: Use ZenbtcMetaData.ABI instead.
var ZenbtcABI = ZenbtcMetaData.ABI

// ZenbtcBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZenbtcMetaData.Bin instead.
var ZenbtcBin = ZenbtcMetaData.Bin

// DeployZenbtc deploys a new Ethereum contract, binding an instance of Zenbtc to it.
func DeployZenbtc(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Zenbtc, error) {
	parsed, err := ZenbtcMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZenbtcBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Zenbtc{ZenbtcCaller: ZenbtcCaller{contract: contract}, ZenbtcTransactor: ZenbtcTransactor{contract: contract}, ZenbtcFilterer: ZenbtcFilterer{contract: contract}}, nil
}

// Zenbtc is an auto generated Go binding around an Ethereum contract.
type Zenbtc struct {
	ZenbtcCaller     // Read-only binding to the contract
	ZenbtcTransactor // Write-only binding to the contract
	ZenbtcFilterer   // Log filterer for contract events
}

// ZenbtcCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZenbtcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenbtcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZenbtcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenbtcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZenbtcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenbtcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZenbtcSession struct {
	Contract     *Zenbtc           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZenbtcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZenbtcCallerSession struct {
	Contract *ZenbtcCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ZenbtcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZenbtcTransactorSession struct {
	Contract     *ZenbtcTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZenbtcRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZenbtcRaw struct {
	Contract *Zenbtc // Generic contract binding to access the raw methods on
}

// ZenbtcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZenbtcCallerRaw struct {
	Contract *ZenbtcCaller // Generic read-only contract binding to access the raw methods on
}

// ZenbtcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZenbtcTransactorRaw struct {
	Contract *ZenbtcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZenbtc creates a new instance of Zenbtc, bound to a specific deployed contract.
func NewZenbtc(address common.Address, backend bind.ContractBackend) (*Zenbtc, error) {
	contract, err := bindZenbtc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Zenbtc{ZenbtcCaller: ZenbtcCaller{contract: contract}, ZenbtcTransactor: ZenbtcTransactor{contract: contract}, ZenbtcFilterer: ZenbtcFilterer{contract: contract}}, nil
}

// NewZenbtcCaller creates a new read-only instance of Zenbtc, bound to a specific deployed contract.
func NewZenbtcCaller(address common.Address, caller bind.ContractCaller) (*ZenbtcCaller, error) {
	contract, err := bindZenbtc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZenbtcCaller{contract: contract}, nil
}

// NewZenbtcTransactor creates a new write-only instance of Zenbtc, bound to a specific deployed contract.
func NewZenbtcTransactor(address common.Address, transactor bind.ContractTransactor) (*ZenbtcTransactor, error) {
	contract, err := bindZenbtc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZenbtcTransactor{contract: contract}, nil
}

// NewZenbtcFilterer creates a new log filterer instance of Zenbtc, bound to a specific deployed contract.
func NewZenbtcFilterer(address common.Address, filterer bind.ContractFilterer) (*ZenbtcFilterer, error) {
	contract, err := bindZenbtc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZenbtcFilterer{contract: contract}, nil
}

// bindZenbtc binds a generic wrapper to an already deployed contract.
func bindZenbtc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZenbtcMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Zenbtc *ZenbtcRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Zenbtc.Contract.ZenbtcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Zenbtc *ZenbtcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Zenbtc.Contract.ZenbtcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Zenbtc *ZenbtcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Zenbtc.Contract.ZenbtcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Zenbtc *ZenbtcCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Zenbtc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Zenbtc *ZenbtcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Zenbtc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Zenbtc *ZenbtcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Zenbtc.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Zenbtc.Contract.DEFAULTADMINROLE(&_Zenbtc.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Zenbtc.Contract.DEFAULTADMINROLE(&_Zenbtc.CallOpts)
}

// FEEROLE is a free data retrieval call binding the contract method 0x2a0276f8.
//
// Solidity: function FEE_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcCaller) FEEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "FEE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FEEROLE is a free data retrieval call binding the contract method 0x2a0276f8.
//
// Solidity: function FEE_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcSession) FEEROLE() ([32]byte, error) {
	return _Zenbtc.Contract.FEEROLE(&_Zenbtc.CallOpts)
}

// FEEROLE is a free data retrieval call binding the contract method 0x2a0276f8.
//
// Solidity: function FEE_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcCallerSession) FEEROLE() ([32]byte, error) {
	return _Zenbtc.Contract.FEEROLE(&_Zenbtc.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcSession) MINTERROLE() ([32]byte, error) {
	return _Zenbtc.Contract.MINTERROLE(&_Zenbtc.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_Zenbtc *ZenbtcCallerSession) MINTERROLE() ([32]byte, error) {
	return _Zenbtc.Contract.MINTERROLE(&_Zenbtc.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Zenbtc *ZenbtcCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Zenbtc *ZenbtcSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Zenbtc.Contract.Allowance(&_Zenbtc.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Zenbtc *ZenbtcCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Zenbtc.Contract.Allowance(&_Zenbtc.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Zenbtc *ZenbtcCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Zenbtc *ZenbtcSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Zenbtc.Contract.BalanceOf(&_Zenbtc.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Zenbtc *ZenbtcCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Zenbtc.Contract.BalanceOf(&_Zenbtc.CallOpts, account)
}

// CurrentVersion is a free data retrieval call binding the contract method 0x9d888e86.
//
// Solidity: function currentVersion() view returns(uint64)
func (_Zenbtc *ZenbtcCaller) CurrentVersion(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "currentVersion")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CurrentVersion is a free data retrieval call binding the contract method 0x9d888e86.
//
// Solidity: function currentVersion() view returns(uint64)
func (_Zenbtc *ZenbtcSession) CurrentVersion() (uint64, error) {
	return _Zenbtc.Contract.CurrentVersion(&_Zenbtc.CallOpts)
}

// CurrentVersion is a free data retrieval call binding the contract method 0x9d888e86.
//
// Solidity: function currentVersion() view returns(uint64)
func (_Zenbtc *ZenbtcCallerSession) CurrentVersion() (uint64, error) {
	return _Zenbtc.Contract.CurrentVersion(&_Zenbtc.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Zenbtc *ZenbtcCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Zenbtc *ZenbtcSession) Decimals() (uint8, error) {
	return _Zenbtc.Contract.Decimals(&_Zenbtc.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Zenbtc *ZenbtcCallerSession) Decimals() (uint8, error) {
	return _Zenbtc.Contract.Decimals(&_Zenbtc.CallOpts)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(uint256 value) view returns(uint256)
func (_Zenbtc *ZenbtcCaller) EstimateFee(opts *bind.CallOpts, value *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "estimateFee", value)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(uint256 value) view returns(uint256)
func (_Zenbtc *ZenbtcSession) EstimateFee(value *big.Int) (*big.Int, error) {
	return _Zenbtc.Contract.EstimateFee(&_Zenbtc.CallOpts, value)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(uint256 value) view returns(uint256)
func (_Zenbtc *ZenbtcCallerSession) EstimateFee(value *big.Int) (*big.Int, error) {
	return _Zenbtc.Contract.EstimateFee(&_Zenbtc.CallOpts, value)
}

// GetAllBurns is a free data retrieval call binding the contract method 0x03b41f74.
//
// Solidity: function getAllBurns() view returns((uint256,uint64,address,bytes)[])
func (_Zenbtc *ZenbtcCaller) GetAllBurns(opts *bind.CallOpts) ([]IZenBTCSidecarHelperZenBTCBurn, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "getAllBurns")

	if err != nil {
		return *new([]IZenBTCSidecarHelperZenBTCBurn), err
	}

	out0 := *abi.ConvertType(out[0], new([]IZenBTCSidecarHelperZenBTCBurn)).(*[]IZenBTCSidecarHelperZenBTCBurn)

	return out0, err

}

// GetAllBurns is a free data retrieval call binding the contract method 0x03b41f74.
//
// Solidity: function getAllBurns() view returns((uint256,uint64,address,bytes)[])
func (_Zenbtc *ZenbtcSession) GetAllBurns() ([]IZenBTCSidecarHelperZenBTCBurn, error) {
	return _Zenbtc.Contract.GetAllBurns(&_Zenbtc.CallOpts)
}

// GetAllBurns is a free data retrieval call binding the contract method 0x03b41f74.
//
// Solidity: function getAllBurns() view returns((uint256,uint64,address,bytes)[])
func (_Zenbtc *ZenbtcCallerSession) GetAllBurns() ([]IZenBTCSidecarHelperZenBTCBurn, error) {
	return _Zenbtc.Contract.GetAllBurns(&_Zenbtc.CallOpts)
}

// GetBasisPointsRate is a free data retrieval call binding the contract method 0xfded2529.
//
// Solidity: function getBasisPointsRate() view returns(uint16)
func (_Zenbtc *ZenbtcCaller) GetBasisPointsRate(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "getBasisPointsRate")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetBasisPointsRate is a free data retrieval call binding the contract method 0xfded2529.
//
// Solidity: function getBasisPointsRate() view returns(uint16)
func (_Zenbtc *ZenbtcSession) GetBasisPointsRate() (uint16, error) {
	return _Zenbtc.Contract.GetBasisPointsRate(&_Zenbtc.CallOpts)
}

// GetBasisPointsRate is a free data retrieval call binding the contract method 0xfded2529.
//
// Solidity: function getBasisPointsRate() view returns(uint16)
func (_Zenbtc *ZenbtcCallerSession) GetBasisPointsRate() (uint16, error) {
	return _Zenbtc.Contract.GetBasisPointsRate(&_Zenbtc.CallOpts)
}

// GetBurn is a free data retrieval call binding the contract method 0x2eb3f49b.
//
// Solidity: function getBurn(uint256 index) view returns((uint256,uint64,address,bytes))
func (_Zenbtc *ZenbtcCaller) GetBurn(opts *bind.CallOpts, index *big.Int) (IZenBTCSidecarHelperZenBTCBurn, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "getBurn", index)

	if err != nil {
		return *new(IZenBTCSidecarHelperZenBTCBurn), err
	}

	out0 := *abi.ConvertType(out[0], new(IZenBTCSidecarHelperZenBTCBurn)).(*IZenBTCSidecarHelperZenBTCBurn)

	return out0, err

}

// GetBurn is a free data retrieval call binding the contract method 0x2eb3f49b.
//
// Solidity: function getBurn(uint256 index) view returns((uint256,uint64,address,bytes))
func (_Zenbtc *ZenbtcSession) GetBurn(index *big.Int) (IZenBTCSidecarHelperZenBTCBurn, error) {
	return _Zenbtc.Contract.GetBurn(&_Zenbtc.CallOpts, index)
}

// GetBurn is a free data retrieval call binding the contract method 0x2eb3f49b.
//
// Solidity: function getBurn(uint256 index) view returns((uint256,uint64,address,bytes))
func (_Zenbtc *ZenbtcCallerSession) GetBurn(index *big.Int) (IZenBTCSidecarHelperZenBTCBurn, error) {
	return _Zenbtc.Contract.GetBurn(&_Zenbtc.CallOpts, index)
}

// GetBurnCount is a free data retrieval call binding the contract method 0xe7cf548c.
//
// Solidity: function getBurnCount() view returns(uint256)
func (_Zenbtc *ZenbtcCaller) GetBurnCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "getBurnCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBurnCount is a free data retrieval call binding the contract method 0xe7cf548c.
//
// Solidity: function getBurnCount() view returns(uint256)
func (_Zenbtc *ZenbtcSession) GetBurnCount() (*big.Int, error) {
	return _Zenbtc.Contract.GetBurnCount(&_Zenbtc.CallOpts)
}

// GetBurnCount is a free data retrieval call binding the contract method 0xe7cf548c.
//
// Solidity: function getBurnCount() view returns(uint256)
func (_Zenbtc *ZenbtcCallerSession) GetBurnCount() (*big.Int, error) {
	return _Zenbtc.Contract.GetBurnCount(&_Zenbtc.CallOpts)
}

// GetBurnsFromTo is a free data retrieval call binding the contract method 0x5acc22c7.
//
// Solidity: function getBurnsFromTo(uint256 from, uint256 to) view returns((uint256,uint64,address,bytes)[])
func (_Zenbtc *ZenbtcCaller) GetBurnsFromTo(opts *bind.CallOpts, from *big.Int, to *big.Int) ([]IZenBTCSidecarHelperZenBTCBurn, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "getBurnsFromTo", from, to)

	if err != nil {
		return *new([]IZenBTCSidecarHelperZenBTCBurn), err
	}

	out0 := *abi.ConvertType(out[0], new([]IZenBTCSidecarHelperZenBTCBurn)).(*[]IZenBTCSidecarHelperZenBTCBurn)

	return out0, err

}

// GetBurnsFromTo is a free data retrieval call binding the contract method 0x5acc22c7.
//
// Solidity: function getBurnsFromTo(uint256 from, uint256 to) view returns((uint256,uint64,address,bytes)[])
func (_Zenbtc *ZenbtcSession) GetBurnsFromTo(from *big.Int, to *big.Int) ([]IZenBTCSidecarHelperZenBTCBurn, error) {
	return _Zenbtc.Contract.GetBurnsFromTo(&_Zenbtc.CallOpts, from, to)
}

// GetBurnsFromTo is a free data retrieval call binding the contract method 0x5acc22c7.
//
// Solidity: function getBurnsFromTo(uint256 from, uint256 to) view returns((uint256,uint64,address,bytes)[])
func (_Zenbtc *ZenbtcCallerSession) GetBurnsFromTo(from *big.Int, to *big.Int) ([]IZenBTCSidecarHelperZenBTCBurn, error) {
	return _Zenbtc.Contract.GetBurnsFromTo(&_Zenbtc.CallOpts, from, to)
}

// GetFeeAddress is a free data retrieval call binding the contract method 0x4e7ceacb.
//
// Solidity: function getFeeAddress() view returns(address)
func (_Zenbtc *ZenbtcCaller) GetFeeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "getFeeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFeeAddress is a free data retrieval call binding the contract method 0x4e7ceacb.
//
// Solidity: function getFeeAddress() view returns(address)
func (_Zenbtc *ZenbtcSession) GetFeeAddress() (common.Address, error) {
	return _Zenbtc.Contract.GetFeeAddress(&_Zenbtc.CallOpts)
}

// GetFeeAddress is a free data retrieval call binding the contract method 0x4e7ceacb.
//
// Solidity: function getFeeAddress() view returns(address)
func (_Zenbtc *ZenbtcCallerSession) GetFeeAddress() (common.Address, error) {
	return _Zenbtc.Contract.GetFeeAddress(&_Zenbtc.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Zenbtc *ZenbtcCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Zenbtc *ZenbtcSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Zenbtc.Contract.GetRoleAdmin(&_Zenbtc.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Zenbtc *ZenbtcCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Zenbtc.Contract.GetRoleAdmin(&_Zenbtc.CallOpts, role)
}

// GetSidecarHelper is a free data retrieval call binding the contract method 0x2b9b5728.
//
// Solidity: function getSidecarHelper() view returns(address)
func (_Zenbtc *ZenbtcCaller) GetSidecarHelper(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "getSidecarHelper")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSidecarHelper is a free data retrieval call binding the contract method 0x2b9b5728.
//
// Solidity: function getSidecarHelper() view returns(address)
func (_Zenbtc *ZenbtcSession) GetSidecarHelper() (common.Address, error) {
	return _Zenbtc.Contract.GetSidecarHelper(&_Zenbtc.CallOpts)
}

// GetSidecarHelper is a free data retrieval call binding the contract method 0x2b9b5728.
//
// Solidity: function getSidecarHelper() view returns(address)
func (_Zenbtc *ZenbtcCallerSession) GetSidecarHelper() (common.Address, error) {
	return _Zenbtc.Contract.GetSidecarHelper(&_Zenbtc.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Zenbtc *ZenbtcCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Zenbtc *ZenbtcSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Zenbtc.Contract.HasRole(&_Zenbtc.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Zenbtc *ZenbtcCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Zenbtc.Contract.HasRole(&_Zenbtc.CallOpts, role, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Zenbtc *ZenbtcCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Zenbtc *ZenbtcSession) Name() (string, error) {
	return _Zenbtc.Contract.Name(&_Zenbtc.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Zenbtc *ZenbtcCallerSession) Name() (string, error) {
	return _Zenbtc.Contract.Name(&_Zenbtc.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Zenbtc *ZenbtcCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Zenbtc *ZenbtcSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Zenbtc.Contract.SupportsInterface(&_Zenbtc.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Zenbtc *ZenbtcCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Zenbtc.Contract.SupportsInterface(&_Zenbtc.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Zenbtc *ZenbtcCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Zenbtc *ZenbtcSession) Symbol() (string, error) {
	return _Zenbtc.Contract.Symbol(&_Zenbtc.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Zenbtc *ZenbtcCallerSession) Symbol() (string, error) {
	return _Zenbtc.Contract.Symbol(&_Zenbtc.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Zenbtc *ZenbtcCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Zenbtc.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Zenbtc *ZenbtcSession) TotalSupply() (*big.Int, error) {
	return _Zenbtc.Contract.TotalSupply(&_Zenbtc.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Zenbtc *ZenbtcCallerSession) TotalSupply() (*big.Int, error) {
	return _Zenbtc.Contract.TotalSupply(&_Zenbtc.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.Approve(&_Zenbtc.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.Approve(&_Zenbtc.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_Zenbtc *ZenbtcTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_Zenbtc *ZenbtcSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.Burn(&_Zenbtc.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_Zenbtc *ZenbtcTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.Burn(&_Zenbtc.TransactOpts, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_Zenbtc *ZenbtcTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "burnFrom", account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_Zenbtc *ZenbtcSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.BurnFrom(&_Zenbtc.TransactOpts, account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_Zenbtc *ZenbtcTransactorSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.BurnFrom(&_Zenbtc.TransactOpts, account, value)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Zenbtc *ZenbtcTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Zenbtc *ZenbtcSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.GrantRole(&_Zenbtc.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Zenbtc *ZenbtcTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.GrantRole(&_Zenbtc.TransactOpts, role, account)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0x2a4b2fd9.
//
// Solidity: function initializeV1(string name_, string symbol_, uint8 decimals_) returns()
func (_Zenbtc *ZenbtcTransactor) InitializeV1(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "initializeV1", name_, symbol_, decimals_)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0x2a4b2fd9.
//
// Solidity: function initializeV1(string name_, string symbol_, uint8 decimals_) returns()
func (_Zenbtc *ZenbtcSession) InitializeV1(name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _Zenbtc.Contract.InitializeV1(&_Zenbtc.TransactOpts, name_, symbol_, decimals_)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0x2a4b2fd9.
//
// Solidity: function initializeV1(string name_, string symbol_, uint8 decimals_) returns()
func (_Zenbtc *ZenbtcTransactorSession) InitializeV1(name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _Zenbtc.Contract.InitializeV1(&_Zenbtc.TransactOpts, name_, symbol_, decimals_)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0xdeb42f13.
//
// Solidity: function initializeV2(string name_, string symbol_, uint8 decimals_) returns()
func (_Zenbtc *ZenbtcTransactor) InitializeV2(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "initializeV2", name_, symbol_, decimals_)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0xdeb42f13.
//
// Solidity: function initializeV2(string name_, string symbol_, uint8 decimals_) returns()
func (_Zenbtc *ZenbtcSession) InitializeV2(name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _Zenbtc.Contract.InitializeV2(&_Zenbtc.TransactOpts, name_, symbol_, decimals_)
}

// InitializeV2 is a paid mutator transaction binding the contract method 0xdeb42f13.
//
// Solidity: function initializeV2(string name_, string symbol_, uint8 decimals_) returns()
func (_Zenbtc *ZenbtcTransactorSession) InitializeV2(name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _Zenbtc.Contract.InitializeV2(&_Zenbtc.TransactOpts, name_, symbol_, decimals_)
}

// InitializeV3 is a paid mutator transaction binding the contract method 0x3101cfcb.
//
// Solidity: function initializeV3(address sidecarHelper_) returns()
func (_Zenbtc *ZenbtcTransactor) InitializeV3(opts *bind.TransactOpts, sidecarHelper_ common.Address) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "initializeV3", sidecarHelper_)
}

// InitializeV3 is a paid mutator transaction binding the contract method 0x3101cfcb.
//
// Solidity: function initializeV3(address sidecarHelper_) returns()
func (_Zenbtc *ZenbtcSession) InitializeV3(sidecarHelper_ common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.InitializeV3(&_Zenbtc.TransactOpts, sidecarHelper_)
}

// InitializeV3 is a paid mutator transaction binding the contract method 0x3101cfcb.
//
// Solidity: function initializeV3(address sidecarHelper_) returns()
func (_Zenbtc *ZenbtcTransactorSession) InitializeV3(sidecarHelper_ common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.InitializeV3(&_Zenbtc.TransactOpts, sidecarHelper_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 value) returns()
func (_Zenbtc *ZenbtcTransactor) Mint(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "mint", account, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 value) returns()
func (_Zenbtc *ZenbtcSession) Mint(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.Mint(&_Zenbtc.TransactOpts, account, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 value) returns()
func (_Zenbtc *ZenbtcTransactorSession) Mint(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.Mint(&_Zenbtc.TransactOpts, account, value)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Zenbtc *ZenbtcTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Zenbtc *ZenbtcSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.RenounceRole(&_Zenbtc.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Zenbtc *ZenbtcTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.RenounceRole(&_Zenbtc.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Zenbtc *ZenbtcTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Zenbtc *ZenbtcSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.RevokeRole(&_Zenbtc.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Zenbtc *ZenbtcTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.RevokeRole(&_Zenbtc.TransactOpts, role, account)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.Transfer(&_Zenbtc.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.Transfer(&_Zenbtc.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.TransferFrom(&_Zenbtc.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Zenbtc *ZenbtcTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Zenbtc.Contract.TransferFrom(&_Zenbtc.TransactOpts, from, to, value)
}

// Unwrap is a paid mutator transaction binding the contract method 0xb413148e.
//
// Solidity: function unwrap(uint256 value, bytes destAddr) returns()
func (_Zenbtc *ZenbtcTransactor) Unwrap(opts *bind.TransactOpts, value *big.Int, destAddr []byte) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "unwrap", value, destAddr)
}

// Unwrap is a paid mutator transaction binding the contract method 0xb413148e.
//
// Solidity: function unwrap(uint256 value, bytes destAddr) returns()
func (_Zenbtc *ZenbtcSession) Unwrap(value *big.Int, destAddr []byte) (*types.Transaction, error) {
	return _Zenbtc.Contract.Unwrap(&_Zenbtc.TransactOpts, value, destAddr)
}

// Unwrap is a paid mutator transaction binding the contract method 0xb413148e.
//
// Solidity: function unwrap(uint256 value, bytes destAddr) returns()
func (_Zenbtc *ZenbtcTransactorSession) Unwrap(value *big.Int, destAddr []byte) (*types.Transaction, error) {
	return _Zenbtc.Contract.Unwrap(&_Zenbtc.TransactOpts, value, destAddr)
}

// UpdateBasisPointsRate is a paid mutator transaction binding the contract method 0xfd154cc7.
//
// Solidity: function updateBasisPointsRate(uint16 basisPointsRate_) returns()
func (_Zenbtc *ZenbtcTransactor) UpdateBasisPointsRate(opts *bind.TransactOpts, basisPointsRate_ uint16) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "updateBasisPointsRate", basisPointsRate_)
}

// UpdateBasisPointsRate is a paid mutator transaction binding the contract method 0xfd154cc7.
//
// Solidity: function updateBasisPointsRate(uint16 basisPointsRate_) returns()
func (_Zenbtc *ZenbtcSession) UpdateBasisPointsRate(basisPointsRate_ uint16) (*types.Transaction, error) {
	return _Zenbtc.Contract.UpdateBasisPointsRate(&_Zenbtc.TransactOpts, basisPointsRate_)
}

// UpdateBasisPointsRate is a paid mutator transaction binding the contract method 0xfd154cc7.
//
// Solidity: function updateBasisPointsRate(uint16 basisPointsRate_) returns()
func (_Zenbtc *ZenbtcTransactorSession) UpdateBasisPointsRate(basisPointsRate_ uint16) (*types.Transaction, error) {
	return _Zenbtc.Contract.UpdateBasisPointsRate(&_Zenbtc.TransactOpts, basisPointsRate_)
}

// UpdateFeeAddress is a paid mutator transaction binding the contract method 0xbbcaac38.
//
// Solidity: function updateFeeAddress(address feeAddress_) returns()
func (_Zenbtc *ZenbtcTransactor) UpdateFeeAddress(opts *bind.TransactOpts, feeAddress_ common.Address) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "updateFeeAddress", feeAddress_)
}

// UpdateFeeAddress is a paid mutator transaction binding the contract method 0xbbcaac38.
//
// Solidity: function updateFeeAddress(address feeAddress_) returns()
func (_Zenbtc *ZenbtcSession) UpdateFeeAddress(feeAddress_ common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.UpdateFeeAddress(&_Zenbtc.TransactOpts, feeAddress_)
}

// UpdateFeeAddress is a paid mutator transaction binding the contract method 0xbbcaac38.
//
// Solidity: function updateFeeAddress(address feeAddress_) returns()
func (_Zenbtc *ZenbtcTransactorSession) UpdateFeeAddress(feeAddress_ common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.UpdateFeeAddress(&_Zenbtc.TransactOpts, feeAddress_)
}

// UpdateSidecarHelper is a paid mutator transaction binding the contract method 0xfaa85988.
//
// Solidity: function updateSidecarHelper(address newHelper) returns()
func (_Zenbtc *ZenbtcTransactor) UpdateSidecarHelper(opts *bind.TransactOpts, newHelper common.Address) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "updateSidecarHelper", newHelper)
}

// UpdateSidecarHelper is a paid mutator transaction binding the contract method 0xfaa85988.
//
// Solidity: function updateSidecarHelper(address newHelper) returns()
func (_Zenbtc *ZenbtcSession) UpdateSidecarHelper(newHelper common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.UpdateSidecarHelper(&_Zenbtc.TransactOpts, newHelper)
}

// UpdateSidecarHelper is a paid mutator transaction binding the contract method 0xfaa85988.
//
// Solidity: function updateSidecarHelper(address newHelper) returns()
func (_Zenbtc *ZenbtcTransactorSession) UpdateSidecarHelper(newHelper common.Address) (*types.Transaction, error) {
	return _Zenbtc.Contract.UpdateSidecarHelper(&_Zenbtc.TransactOpts, newHelper)
}

// Wrap is a paid mutator transaction binding the contract method 0xf1d96dd3.
//
// Solidity: function wrap(address account, uint64 value, uint64 fee) returns()
func (_Zenbtc *ZenbtcTransactor) Wrap(opts *bind.TransactOpts, account common.Address, value uint64, fee uint64) (*types.Transaction, error) {
	return _Zenbtc.contract.Transact(opts, "wrap", account, value, fee)
}

// Wrap is a paid mutator transaction binding the contract method 0xf1d96dd3.
//
// Solidity: function wrap(address account, uint64 value, uint64 fee) returns()
func (_Zenbtc *ZenbtcSession) Wrap(account common.Address, value uint64, fee uint64) (*types.Transaction, error) {
	return _Zenbtc.Contract.Wrap(&_Zenbtc.TransactOpts, account, value, fee)
}

// Wrap is a paid mutator transaction binding the contract method 0xf1d96dd3.
//
// Solidity: function wrap(address account, uint64 value, uint64 fee) returns()
func (_Zenbtc *ZenbtcTransactorSession) Wrap(account common.Address, value uint64, fee uint64) (*types.Transaction, error) {
	return _Zenbtc.Contract.Wrap(&_Zenbtc.TransactOpts, account, value, fee)
}

// ZenbtcApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Zenbtc contract.
type ZenbtcApprovalIterator struct {
	Event *ZenbtcApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcApproval represents a Approval event raised by the Zenbtc contract.
type ZenbtcApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Zenbtc *ZenbtcFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ZenbtcApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcApprovalIterator{contract: _Zenbtc.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Zenbtc *ZenbtcFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZenbtcApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcApproval)
				if err := _Zenbtc.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Zenbtc *ZenbtcFilterer) ParseApproval(log types.Log) (*ZenbtcApproval, error) {
	event := new(ZenbtcApproval)
	if err := _Zenbtc.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcBasisPointsRateUpdatedIterator is returned from FilterBasisPointsRateUpdated and is used to iterate over the raw logs and unpacked data for BasisPointsRateUpdated events raised by the Zenbtc contract.
type ZenbtcBasisPointsRateUpdatedIterator struct {
	Event *ZenbtcBasisPointsRateUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcBasisPointsRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcBasisPointsRateUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcBasisPointsRateUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcBasisPointsRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcBasisPointsRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcBasisPointsRateUpdated represents a BasisPointsRateUpdated event raised by the Zenbtc contract.
type ZenbtcBasisPointsRateUpdated struct {
	BasisPointsRateUpdated uint16
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterBasisPointsRateUpdated is a free log retrieval operation binding the contract event 0x79d2a06c43b232cf9d1835b5e915efe74621561aaf75fecec91fd75a940f9d70.
//
// Solidity: event BasisPointsRateUpdated(uint16 basisPointsRateUpdated)
func (_Zenbtc *ZenbtcFilterer) FilterBasisPointsRateUpdated(opts *bind.FilterOpts) (*ZenbtcBasisPointsRateUpdatedIterator, error) {

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "BasisPointsRateUpdated")
	if err != nil {
		return nil, err
	}
	return &ZenbtcBasisPointsRateUpdatedIterator{contract: _Zenbtc.contract, event: "BasisPointsRateUpdated", logs: logs, sub: sub}, nil
}

// WatchBasisPointsRateUpdated is a free log subscription operation binding the contract event 0x79d2a06c43b232cf9d1835b5e915efe74621561aaf75fecec91fd75a940f9d70.
//
// Solidity: event BasisPointsRateUpdated(uint16 basisPointsRateUpdated)
func (_Zenbtc *ZenbtcFilterer) WatchBasisPointsRateUpdated(opts *bind.WatchOpts, sink chan<- *ZenbtcBasisPointsRateUpdated) (event.Subscription, error) {

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "BasisPointsRateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcBasisPointsRateUpdated)
				if err := _Zenbtc.contract.UnpackLog(event, "BasisPointsRateUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBasisPointsRateUpdated is a log parse operation binding the contract event 0x79d2a06c43b232cf9d1835b5e915efe74621561aaf75fecec91fd75a940f9d70.
//
// Solidity: event BasisPointsRateUpdated(uint16 basisPointsRateUpdated)
func (_Zenbtc *ZenbtcFilterer) ParseBasisPointsRateUpdated(log types.Log) (*ZenbtcBasisPointsRateUpdated, error) {
	event := new(ZenbtcBasisPointsRateUpdated)
	if err := _Zenbtc.contract.UnpackLog(event, "BasisPointsRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcFeeAddressUpdatedIterator is returned from FilterFeeAddressUpdated and is used to iterate over the raw logs and unpacked data for FeeAddressUpdated events raised by the Zenbtc contract.
type ZenbtcFeeAddressUpdatedIterator struct {
	Event *ZenbtcFeeAddressUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcFeeAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcFeeAddressUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcFeeAddressUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcFeeAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcFeeAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcFeeAddressUpdated represents a FeeAddressUpdated event raised by the Zenbtc contract.
type ZenbtcFeeAddressUpdated struct {
	NewFeeAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterFeeAddressUpdated is a free log retrieval operation binding the contract event 0x446e39bcf1b47cfadfaa23442cb4b34682cfe6bd9220da084894e3b1f834e4f3.
//
// Solidity: event FeeAddressUpdated(address indexed newFeeAddress)
func (_Zenbtc *ZenbtcFilterer) FilterFeeAddressUpdated(opts *bind.FilterOpts, newFeeAddress []common.Address) (*ZenbtcFeeAddressUpdatedIterator, error) {

	var newFeeAddressRule []interface{}
	for _, newFeeAddressItem := range newFeeAddress {
		newFeeAddressRule = append(newFeeAddressRule, newFeeAddressItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "FeeAddressUpdated", newFeeAddressRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcFeeAddressUpdatedIterator{contract: _Zenbtc.contract, event: "FeeAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchFeeAddressUpdated is a free log subscription operation binding the contract event 0x446e39bcf1b47cfadfaa23442cb4b34682cfe6bd9220da084894e3b1f834e4f3.
//
// Solidity: event FeeAddressUpdated(address indexed newFeeAddress)
func (_Zenbtc *ZenbtcFilterer) WatchFeeAddressUpdated(opts *bind.WatchOpts, sink chan<- *ZenbtcFeeAddressUpdated, newFeeAddress []common.Address) (event.Subscription, error) {

	var newFeeAddressRule []interface{}
	for _, newFeeAddressItem := range newFeeAddress {
		newFeeAddressRule = append(newFeeAddressRule, newFeeAddressItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "FeeAddressUpdated", newFeeAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcFeeAddressUpdated)
				if err := _Zenbtc.contract.UnpackLog(event, "FeeAddressUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeeAddressUpdated is a log parse operation binding the contract event 0x446e39bcf1b47cfadfaa23442cb4b34682cfe6bd9220da084894e3b1f834e4f3.
//
// Solidity: event FeeAddressUpdated(address indexed newFeeAddress)
func (_Zenbtc *ZenbtcFilterer) ParseFeeAddressUpdated(log types.Log) (*ZenbtcFeeAddressUpdated, error) {
	event := new(ZenbtcFeeAddressUpdated)
	if err := _Zenbtc.contract.UnpackLog(event, "FeeAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Zenbtc contract.
type ZenbtcInitializedIterator struct {
	Event *ZenbtcInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcInitialized represents a Initialized event raised by the Zenbtc contract.
type ZenbtcInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Zenbtc *ZenbtcFilterer) FilterInitialized(opts *bind.FilterOpts) (*ZenbtcInitializedIterator, error) {

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ZenbtcInitializedIterator{contract: _Zenbtc.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Zenbtc *ZenbtcFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ZenbtcInitialized) (event.Subscription, error) {

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcInitialized)
				if err := _Zenbtc.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Zenbtc *ZenbtcFilterer) ParseInitialized(log types.Log) (*ZenbtcInitialized, error) {
	event := new(ZenbtcInitialized)
	if err := _Zenbtc.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcRedeemFeeUpdatedIterator is returned from FilterRedeemFeeUpdated and is used to iterate over the raw logs and unpacked data for RedeemFeeUpdated events raised by the Zenbtc contract.
type ZenbtcRedeemFeeUpdatedIterator struct {
	Event *ZenbtcRedeemFeeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcRedeemFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcRedeemFeeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcRedeemFeeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcRedeemFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcRedeemFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcRedeemFeeUpdated represents a RedeemFeeUpdated event raised by the Zenbtc contract.
type ZenbtcRedeemFeeUpdated struct {
	NewRedeemFee *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRedeemFeeUpdated is a free log retrieval operation binding the contract event 0xd6c7508d6658ccee36b7b7d7fd72e5cbaeefb40c64eff24e9ae7470e846304ee.
//
// Solidity: event RedeemFeeUpdated(uint256 newRedeemFee)
func (_Zenbtc *ZenbtcFilterer) FilterRedeemFeeUpdated(opts *bind.FilterOpts) (*ZenbtcRedeemFeeUpdatedIterator, error) {

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "RedeemFeeUpdated")
	if err != nil {
		return nil, err
	}
	return &ZenbtcRedeemFeeUpdatedIterator{contract: _Zenbtc.contract, event: "RedeemFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchRedeemFeeUpdated is a free log subscription operation binding the contract event 0xd6c7508d6658ccee36b7b7d7fd72e5cbaeefb40c64eff24e9ae7470e846304ee.
//
// Solidity: event RedeemFeeUpdated(uint256 newRedeemFee)
func (_Zenbtc *ZenbtcFilterer) WatchRedeemFeeUpdated(opts *bind.WatchOpts, sink chan<- *ZenbtcRedeemFeeUpdated) (event.Subscription, error) {

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "RedeemFeeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcRedeemFeeUpdated)
				if err := _Zenbtc.contract.UnpackLog(event, "RedeemFeeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeemFeeUpdated is a log parse operation binding the contract event 0xd6c7508d6658ccee36b7b7d7fd72e5cbaeefb40c64eff24e9ae7470e846304ee.
//
// Solidity: event RedeemFeeUpdated(uint256 newRedeemFee)
func (_Zenbtc *ZenbtcFilterer) ParseRedeemFeeUpdated(log types.Log) (*ZenbtcRedeemFeeUpdated, error) {
	event := new(ZenbtcRedeemFeeUpdated)
	if err := _Zenbtc.contract.UnpackLog(event, "RedeemFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Zenbtc contract.
type ZenbtcRoleAdminChangedIterator struct {
	Event *ZenbtcRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcRoleAdminChanged represents a RoleAdminChanged event raised by the Zenbtc contract.
type ZenbtcRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Zenbtc *ZenbtcFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ZenbtcRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcRoleAdminChangedIterator{contract: _Zenbtc.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Zenbtc *ZenbtcFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ZenbtcRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcRoleAdminChanged)
				if err := _Zenbtc.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Zenbtc *ZenbtcFilterer) ParseRoleAdminChanged(log types.Log) (*ZenbtcRoleAdminChanged, error) {
	event := new(ZenbtcRoleAdminChanged)
	if err := _Zenbtc.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Zenbtc contract.
type ZenbtcRoleGrantedIterator struct {
	Event *ZenbtcRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcRoleGranted represents a RoleGranted event raised by the Zenbtc contract.
type ZenbtcRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zenbtc *ZenbtcFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ZenbtcRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcRoleGrantedIterator{contract: _Zenbtc.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zenbtc *ZenbtcFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ZenbtcRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcRoleGranted)
				if err := _Zenbtc.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zenbtc *ZenbtcFilterer) ParseRoleGranted(log types.Log) (*ZenbtcRoleGranted, error) {
	event := new(ZenbtcRoleGranted)
	if err := _Zenbtc.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Zenbtc contract.
type ZenbtcRoleRevokedIterator struct {
	Event *ZenbtcRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcRoleRevoked represents a RoleRevoked event raised by the Zenbtc contract.
type ZenbtcRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zenbtc *ZenbtcFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ZenbtcRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcRoleRevokedIterator{contract: _Zenbtc.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zenbtc *ZenbtcFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ZenbtcRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcRoleRevoked)
				if err := _Zenbtc.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Zenbtc *ZenbtcFilterer) ParseRoleRevoked(log types.Log) (*ZenbtcRoleRevoked, error) {
	event := new(ZenbtcRoleRevoked)
	if err := _Zenbtc.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcSidecarHelperUpdatedIterator is returned from FilterSidecarHelperUpdated and is used to iterate over the raw logs and unpacked data for SidecarHelperUpdated events raised by the Zenbtc contract.
type ZenbtcSidecarHelperUpdatedIterator struct {
	Event *ZenbtcSidecarHelperUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcSidecarHelperUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcSidecarHelperUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcSidecarHelperUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcSidecarHelperUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcSidecarHelperUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcSidecarHelperUpdated represents a SidecarHelperUpdated event raised by the Zenbtc contract.
type ZenbtcSidecarHelperUpdated struct {
	OldHelper common.Address
	NewHelper common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSidecarHelperUpdated is a free log retrieval operation binding the contract event 0xae0961ff5b393c8e05811c7cccbbbbb4453f9336fcd0197d1d87ee17072bdbcf.
//
// Solidity: event SidecarHelperUpdated(address indexed oldHelper, address indexed newHelper)
func (_Zenbtc *ZenbtcFilterer) FilterSidecarHelperUpdated(opts *bind.FilterOpts, oldHelper []common.Address, newHelper []common.Address) (*ZenbtcSidecarHelperUpdatedIterator, error) {

	var oldHelperRule []interface{}
	for _, oldHelperItem := range oldHelper {
		oldHelperRule = append(oldHelperRule, oldHelperItem)
	}
	var newHelperRule []interface{}
	for _, newHelperItem := range newHelper {
		newHelperRule = append(newHelperRule, newHelperItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "SidecarHelperUpdated", oldHelperRule, newHelperRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcSidecarHelperUpdatedIterator{contract: _Zenbtc.contract, event: "SidecarHelperUpdated", logs: logs, sub: sub}, nil
}

// WatchSidecarHelperUpdated is a free log subscription operation binding the contract event 0xae0961ff5b393c8e05811c7cccbbbbb4453f9336fcd0197d1d87ee17072bdbcf.
//
// Solidity: event SidecarHelperUpdated(address indexed oldHelper, address indexed newHelper)
func (_Zenbtc *ZenbtcFilterer) WatchSidecarHelperUpdated(opts *bind.WatchOpts, sink chan<- *ZenbtcSidecarHelperUpdated, oldHelper []common.Address, newHelper []common.Address) (event.Subscription, error) {

	var oldHelperRule []interface{}
	for _, oldHelperItem := range oldHelper {
		oldHelperRule = append(oldHelperRule, oldHelperItem)
	}
	var newHelperRule []interface{}
	for _, newHelperItem := range newHelper {
		newHelperRule = append(newHelperRule, newHelperItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "SidecarHelperUpdated", oldHelperRule, newHelperRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcSidecarHelperUpdated)
				if err := _Zenbtc.contract.UnpackLog(event, "SidecarHelperUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSidecarHelperUpdated is a log parse operation binding the contract event 0xae0961ff5b393c8e05811c7cccbbbbb4453f9336fcd0197d1d87ee17072bdbcf.
//
// Solidity: event SidecarHelperUpdated(address indexed oldHelper, address indexed newHelper)
func (_Zenbtc *ZenbtcFilterer) ParseSidecarHelperUpdated(log types.Log) (*ZenbtcSidecarHelperUpdated, error) {
	event := new(ZenbtcSidecarHelperUpdated)
	if err := _Zenbtc.contract.UnpackLog(event, "SidecarHelperUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcTokenRedemptionIterator is returned from FilterTokenRedemption and is used to iterate over the raw logs and unpacked data for TokenRedemption events raised by the Zenbtc contract.
type ZenbtcTokenRedemptionIterator struct {
	Event *ZenbtcTokenRedemption // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcTokenRedemptionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcTokenRedemption)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcTokenRedemption)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcTokenRedemptionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcTokenRedemptionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcTokenRedemption represents a TokenRedemption event raised by the Zenbtc contract.
type ZenbtcTokenRedemption struct {
	Redeemer common.Address
	Value    uint64
	DestAddr []byte
	Fee      *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokenRedemption is a free log retrieval operation binding the contract event 0x4c971c8b2abb197a17896b2fc57f597830db78d3556831e3faa337a596150f22.
//
// Solidity: event TokenRedemption(address indexed redeemer, uint64 value, bytes destAddr, uint256 fee)
func (_Zenbtc *ZenbtcFilterer) FilterTokenRedemption(opts *bind.FilterOpts, redeemer []common.Address) (*ZenbtcTokenRedemptionIterator, error) {

	var redeemerRule []interface{}
	for _, redeemerItem := range redeemer {
		redeemerRule = append(redeemerRule, redeemerItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "TokenRedemption", redeemerRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcTokenRedemptionIterator{contract: _Zenbtc.contract, event: "TokenRedemption", logs: logs, sub: sub}, nil
}

// WatchTokenRedemption is a free log subscription operation binding the contract event 0x4c971c8b2abb197a17896b2fc57f597830db78d3556831e3faa337a596150f22.
//
// Solidity: event TokenRedemption(address indexed redeemer, uint64 value, bytes destAddr, uint256 fee)
func (_Zenbtc *ZenbtcFilterer) WatchTokenRedemption(opts *bind.WatchOpts, sink chan<- *ZenbtcTokenRedemption, redeemer []common.Address) (event.Subscription, error) {

	var redeemerRule []interface{}
	for _, redeemerItem := range redeemer {
		redeemerRule = append(redeemerRule, redeemerItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "TokenRedemption", redeemerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcTokenRedemption)
				if err := _Zenbtc.contract.UnpackLog(event, "TokenRedemption", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTokenRedemption is a log parse operation binding the contract event 0x4c971c8b2abb197a17896b2fc57f597830db78d3556831e3faa337a596150f22.
//
// Solidity: event TokenRedemption(address indexed redeemer, uint64 value, bytes destAddr, uint256 fee)
func (_Zenbtc *ZenbtcFilterer) ParseTokenRedemption(log types.Log) (*ZenbtcTokenRedemption, error) {
	event := new(ZenbtcTokenRedemption)
	if err := _Zenbtc.contract.UnpackLog(event, "TokenRedemption", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcTokensMintedWithFeeIterator is returned from FilterTokensMintedWithFee and is used to iterate over the raw logs and unpacked data for TokensMintedWithFee events raised by the Zenbtc contract.
type ZenbtcTokensMintedWithFeeIterator struct {
	Event *ZenbtcTokensMintedWithFee // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcTokensMintedWithFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcTokensMintedWithFee)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcTokensMintedWithFee)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcTokensMintedWithFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcTokensMintedWithFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcTokensMintedWithFee represents a TokensMintedWithFee event raised by the Zenbtc contract.
type ZenbtcTokensMintedWithFee struct {
	Recipient common.Address
	Value     uint64
	Fee       uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokensMintedWithFee is a free log retrieval operation binding the contract event 0x890f2fd0578a1ba6f8735fc60c94d1ec453ad71501d82ba4d5a9041f9ac14f2e.
//
// Solidity: event TokensMintedWithFee(address indexed recipient, uint64 value, uint64 fee)
func (_Zenbtc *ZenbtcFilterer) FilterTokensMintedWithFee(opts *bind.FilterOpts, recipient []common.Address) (*ZenbtcTokensMintedWithFeeIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "TokensMintedWithFee", recipientRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcTokensMintedWithFeeIterator{contract: _Zenbtc.contract, event: "TokensMintedWithFee", logs: logs, sub: sub}, nil
}

// WatchTokensMintedWithFee is a free log subscription operation binding the contract event 0x890f2fd0578a1ba6f8735fc60c94d1ec453ad71501d82ba4d5a9041f9ac14f2e.
//
// Solidity: event TokensMintedWithFee(address indexed recipient, uint64 value, uint64 fee)
func (_Zenbtc *ZenbtcFilterer) WatchTokensMintedWithFee(opts *bind.WatchOpts, sink chan<- *ZenbtcTokensMintedWithFee, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "TokensMintedWithFee", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcTokensMintedWithFee)
				if err := _Zenbtc.contract.UnpackLog(event, "TokensMintedWithFee", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTokensMintedWithFee is a log parse operation binding the contract event 0x890f2fd0578a1ba6f8735fc60c94d1ec453ad71501d82ba4d5a9041f9ac14f2e.
//
// Solidity: event TokensMintedWithFee(address indexed recipient, uint64 value, uint64 fee)
func (_Zenbtc *ZenbtcFilterer) ParseTokensMintedWithFee(log types.Log) (*ZenbtcTokensMintedWithFee, error) {
	event := new(ZenbtcTokensMintedWithFee)
	if err := _Zenbtc.contract.UnpackLog(event, "TokensMintedWithFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Zenbtc contract.
type ZenbtcTransferIterator struct {
	Event *ZenbtcTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcTransfer represents a Transfer event raised by the Zenbtc contract.
type ZenbtcTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Zenbtc *ZenbtcFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ZenbtcTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ZenbtcTransferIterator{contract: _Zenbtc.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Zenbtc *ZenbtcFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZenbtcTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcTransfer)
				if err := _Zenbtc.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Zenbtc *ZenbtcFilterer) ParseTransfer(log types.Log) (*ZenbtcTransfer, error) {
	event := new(ZenbtcTransfer)
	if err := _Zenbtc.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenbtcZenBTControllerUpdatedIterator is returned from FilterZenBTControllerUpdated and is used to iterate over the raw logs and unpacked data for ZenBTControllerUpdated events raised by the Zenbtc contract.
type ZenbtcZenBTControllerUpdatedIterator struct {
	Event *ZenbtcZenBTControllerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ZenbtcZenBTControllerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenbtcZenBTControllerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ZenbtcZenBTControllerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ZenbtcZenBTControllerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenbtcZenBTControllerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenbtcZenBTControllerUpdated represents a ZenBTControllerUpdated event raised by the Zenbtc contract.
type ZenbtcZenBTControllerUpdated struct {
	ZenBTControllerUpdated common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterZenBTControllerUpdated is a free log retrieval operation binding the contract event 0xc4c1343b8690c9be210951438cee3b760df1a37031ac0799c086b90d6af53950.
//
// Solidity: event ZenBTControllerUpdated(address zenBTControllerUpdated)
func (_Zenbtc *ZenbtcFilterer) FilterZenBTControllerUpdated(opts *bind.FilterOpts) (*ZenbtcZenBTControllerUpdatedIterator, error) {

	logs, sub, err := _Zenbtc.contract.FilterLogs(opts, "ZenBTControllerUpdated")
	if err != nil {
		return nil, err
	}
	return &ZenbtcZenBTControllerUpdatedIterator{contract: _Zenbtc.contract, event: "ZenBTControllerUpdated", logs: logs, sub: sub}, nil
}

// WatchZenBTControllerUpdated is a free log subscription operation binding the contract event 0xc4c1343b8690c9be210951438cee3b760df1a37031ac0799c086b90d6af53950.
//
// Solidity: event ZenBTControllerUpdated(address zenBTControllerUpdated)
func (_Zenbtc *ZenbtcFilterer) WatchZenBTControllerUpdated(opts *bind.WatchOpts, sink chan<- *ZenbtcZenBTControllerUpdated) (event.Subscription, error) {

	logs, sub, err := _Zenbtc.contract.WatchLogs(opts, "ZenBTControllerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenbtcZenBTControllerUpdated)
				if err := _Zenbtc.contract.UnpackLog(event, "ZenBTControllerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseZenBTControllerUpdated is a log parse operation binding the contract event 0xc4c1343b8690c9be210951438cee3b760df1a37031ac0799c086b90d6af53950.
//
// Solidity: event ZenBTControllerUpdated(address zenBTControllerUpdated)
func (_Zenbtc *ZenbtcFilterer) ParseZenBTControllerUpdated(log types.Log) (*ZenbtcZenBTControllerUpdated, error) {
	event := new(ZenbtcZenBTControllerUpdated)
	if err := _Zenbtc.contract.UnpackLog(event, "ZenBTControllerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
