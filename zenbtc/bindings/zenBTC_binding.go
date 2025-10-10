// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// ZenBTCMetaData contains all meta data concerning the ZenBTC contract.
var ZenBTCMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"basisPointsRateUpdated\",\"type\":\"uint16\"}],\"name\":\"BasisPointsRateUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newFeeAddress\",\"type\":\"address\"}],\"name\":\"FeeAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newRedeemFee\",\"type\":\"uint256\"}],\"name\":\"RedeemFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"destAddr\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"TokenRedemption\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"}],\"name\":\"TokensMintedWithFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"zenBTControllerUpdated\",\"type\":\"address\"}],\"name\":\"ZenBTControllerUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"estimateFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBasisPointsRate\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"}],\"name\":\"initializeV1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destAddr\",\"type\":\"bytes\"}],\"name\":\"unwrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"basisPointsRate_\",\"type\":\"uint16\"}],\"name\":\"updateBasisPointsRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeAddress_\",\"type\":\"address\"}],\"name\":\"updateFeeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"}],\"name\":\"wrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100d0565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161561006e5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100cd5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b6119f3806100df6000396000f3fe608060405234801561001057600080fd5b50600436106101c45760003560e01c80634e7ceacb116100f9578063bbcaac3811610097578063dd62ed3e11610071578063dd62ed3e1461042e578063f1d96dd314610441578063fd154cc714610454578063fded25291461046757600080fd5b8063bbcaac38146103f3578063d539139314610406578063d547741f1461041b57600080fd5b806395d89b41116100d357806395d89b41146103bd578063a217fddf146103c5578063a9059cbb146103cd578063b413148e146103e057600080fd5b80634e7ceacb1461033a57806370a082311461037457806391d14854146103aa57600080fd5b80632a0276f811610166578063313ce56711610140578063313ce567146102d657806336568abe1461030157806340c10f191461031457806342966c681461032757600080fd5b80632a0276f8146102875780632a4b2fd9146102ae5780632f2ff15d146102c357600080fd5b8063127e8e4d116101a2578063127e8e4d1461021957806318160ddd1461023a57806323b872dd14610261578063248a9ca31461027457600080fd5b806301ffc9a7146101c957806306fdde03146101f1578063095ea7b314610206575b600080fd5b6101dc6101d736600461139e565b61048a565b60405190151581526020015b60405180910390f35b6101f96104c1565b6040516101e8919061140e565b6101dc61021436600461143d565b610584565b61022c610227366004611467565b61059c565b6040519081526020016101e8565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace025461022c565b6101dc61026f366004611480565b6105a7565b61022c610282366004611467565b6105cb565b61022c7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf2881565b6102c16102bc366004611567565b6105ed565b005b6102c16102d13660046115e4565b6106ff565b60008051602061193e8339815191525462010000900460ff1660405160ff90911681526020016101e8565b6102c161030f3660046115e4565b610721565b6102c161032236600461143d565b610759565b6102c1610335366004611467565b61077b565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace06546040516001600160a01b0390911681526020016101e8565b61022c610382366004611610565b6001600160a01b0316600090815260008051602061195e833981519152602052604090205490565b6101dc6103b83660046115e4565b6107a1565b6101f96107d9565b61022c600081565b6101dc6103db36600461143d565b610818565b6102c16103ee36600461162b565b610826565b6102c1610401366004611610565b61090c565b61022c60008051602061199e83398151915281565b6102c16104293660046115e4565b61093f565b61022c61043c366004611685565b61095b565b6102c161044f3660046116c6565b6109a5565b6102c1610462366004611709565b610b20565b60008051602061193e8339815191525460405161ffff90911681526020016101e8565b60006001600160e01b03198216637965db0b60e01b14806104bb57506301ffc9a760e01b6001600160e01b03198316145b92915050565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace04805460609160008051602061195e833981519152916105009061172d565b80601f016020809104026020016040519081016040528092919081815260200182805461052c9061172d565b80156105795780601f1061054e57610100808354040283529160200191610579565b820191906000526020600020905b81548152906001019060200180831161055c57829003601f168201915b505050505091505090565b600033610592818585610b53565b5060019392505050565b60006104bb82610b60565b6000336105b5858285610c3d565b6105c0858585610c9d565b506001949350505050565b600090815260008051602061197e833981519152602052604090206001015490565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff1615906001600160401b03166000811580156106325750825b90506000826001600160401b0316600114801561064e5750303b155b90508115801561065c575080155b1561067a5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156106a457845460ff60401b1916600160401b1785555b6106af888888610d27565b83156106f557845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050505050565b610708826105cb565b61071181610d4a565b61071b8383610d57565b50505050565b6001600160a01b038116331461074a5760405163334bd91960e11b815260040160405180910390fd5b6107548282610e03565b505050565b60008051602061199e83398151915261077181610d4a565b6107548383610e7f565b60008051602061199e83398151915261079381610d4a565b61079d3383610eb5565b5050565b600091825260008051602061197e833981519152602090815260408084206001600160a01b0393909316845291905290205460ff1690565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace05805460609160008051602061195e833981519152916105009061172d565b600033610592818585610c9d565b60008051602061195e8339815191526001600160401b038311156108835760405162461bcd60e51b815260206004820152600f60248201526e56616c756520746f6f206c6172676560881b60448201526064015b60405180910390fd5b600061088e84610b60565b905061089a3385610eb5565b60068201546108b2906001600160a01b031682610e7f565b60006108be828661177d565b9050336001600160a01b03167f4c971c8b2abb197a17896b2fc57f597830db78d3556831e3faa337a596150f228286856040516108fd93929190611790565b60405180910390a25050505050565b7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf2861093681610d4a565b61079d82610eeb565b610948826105cb565b61095181610d4a565b61071b8383610e03565b6001600160a01b0391821660009081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace016020908152604080832093909416825291909152205490565b60008051602061199e8339815191526109bd81610d4a565b6001600160401b038381161115610a085760405162461bcd60e51b815260206004820152600f60248201526e56616c756520746f6f206c6172676560881b604482015260640161087a565b826001600160401b0316826001600160401b031610610a785760405162461bcd60e51b815260206004820152602660248201527f5a656e4254433a2046656520657863656564732074686520616d6f756e7420746044820152651bc81b5a5b9d60d21b606482015260840161087a565b60008051602061195e8339815191526000610a9384866117c2565b6006830154909150610ab7906001600160a01b03166001600160401b038616610e7f565b610aca86826001600160401b0316610e7f565b604080516001600160401b038084168252861660208201526001600160a01b038816917f890f2fd0578a1ba6f8735fc60c94d1ec453ad71501d82ba4d5a9041f9ac14f2e910160405180910390a2505050505050565b7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf28610b4a81610d4a565b61079d82610f66565b6107548383836001611039565b60006127108211610ba55760405162461bcd60e51b815260206004820152600f60248201526e15985b1d59481d1bdbc81cdb585b1b608a1b604482015260640161087a565b60008051602061193e8339815191525460008051602061195e83398151915290600090610bd790859061ffff16611121565b9050838110610c365760405162461bcd60e51b815260206004820152602560248201527f5a656e4254433a2072656465656d2066656520657863656564732074686520616044820152641b5bdd5b9d60da1b606482015260840161087a565b9392505050565b6000610c49848461095b565b9050600019811461071b5781811015610c8e57604051637dc7a0d960e11b81526001600160a01b0384166004820152602481018290526044810183905260640161087a565b61071b84848484036000611039565b6001600160a01b038316610cc757604051634b637e8f60e11b81526000600482015260240161087a565b6001600160a01b038216610cf15760405163ec442f0560e01b81526000600482015260240161087a565b306001600160a01b03831603610d1c5760405163ec442f0560e01b815230600482015260240161087a565b61075483838361113e565b610d2f61127c565b610d376112c7565b610d3f6112cf565b6107548383836112df565b610d548133611352565b50565b600060008051602061197e833981519152610d7284846107a1565b610df2576000848152602082815260408083206001600160a01b03871684529091529020805460ff19166001179055610da83390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019150506104bb565b60009150506104bb565b5092915050565b600060008051602061197e833981519152610e1e84846107a1565b15610df2576000848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a460019150506104bb565b6001600160a01b038216610ea95760405163ec442f0560e01b81526000600482015260240161087a565b61079d6000838361113e565b6001600160a01b038216610edf57604051634b637e8f60e11b81526000600482015260240161087a565b61079d8260008361113e565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0680546001600160a01b0319166001600160a01b03831690811790915560405160008051602061195e83398151915291907f446e39bcf1b47cfadfaa23442cb4b34682cfe6bd9220da084894e3b1f834e4f390600090a25050565b61271061ffff82161115610fd05760405162461bcd60e51b815260206004820152602b60248201527f626173697320706f696e747320726174652063616e6e6f74206578636565642060448201526a31303030302075696e747360a81b606482015260840161087a565b60008051602061193e833981519152805461ffff831661ffff199091168117909155604080519182525160008051602061195e833981519152917f79d2a06c43b232cf9d1835b5e915efe74621561aaf75fecec91fd75a940f9d70919081900360200190a15050565b60008051602061195e8339815191526001600160a01b0385166110725760405163e602df0560e01b81526000600482015260240161087a565b6001600160a01b03841661109c57604051634a1406b160e11b81526000600482015260240161087a565b6001600160a01b0380861660009081526001830160209081526040808320938816835292905220839055811561111a57836001600160a01b0316856001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258560405161111191815260200190565b60405180910390a35b5050505050565b600061271061113461ffff8416856117e2565b610c3691906117f9565b60008051602061195e8339815191526001600160a01b03841661117a578181600201600082825461116f919061181b565b909155506111ec9050565b6001600160a01b038416600090815260208290526040902054828110156111cd5760405163391434e360e21b81526001600160a01b0386166004820152602481018290526044810184905260640161087a565b6001600160a01b03851660009081526020839052604090209083900390555b6001600160a01b03831661120a576002810180548390039055611229565b6001600160a01b03831660009081526020829052604090208054830190555b826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161126e91815260200190565b60405180910390a350505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054600160401b900460ff166112c557604051631afcd79f60e31b815260040160405180910390fd5b565b6112c561127c565b6112d761127c565b6112c561138b565b6112e761127c565b60008051602061195e8339815191527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace04611321858261187e565b5060058101611330848261187e565b50600301805460ff909216620100000262ff0000199092169190911790555050565b61135c82826107a1565b61079d5760405163e2517d3f60e01b81526001600160a01b03821660048201526024810183905260440161087a565b61139361127c565b610d54600033610d57565b6000602082840312156113b057600080fd5b81356001600160e01b031981168114610c3657600080fd5b6000815180845260005b818110156113ee576020818501810151868301820152016113d2565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000610c3660208301846113c8565b80356001600160a01b038116811461143857600080fd5b919050565b6000806040838503121561145057600080fd5b61145983611421565b946020939093013593505050565b60006020828403121561147957600080fd5b5035919050565b60008060006060848603121561149557600080fd5b61149e84611421565b92506114ac60208501611421565b9150604084013590509250925092565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b03808411156114ec576114ec6114bc565b604051601f8501601f19908116603f01168101908282118183101715611514576115146114bc565b8160405280935085815286868601111561152d57600080fd5b858560208301376000602087830101525050509392505050565b600082601f83011261155857600080fd5b610c36838335602085016114d2565b60008060006060848603121561157c57600080fd5b83356001600160401b038082111561159357600080fd5b61159f87838801611547565b945060208601359150808211156115b557600080fd5b506115c286828701611547565b925050604084013560ff811681146115d957600080fd5b809150509250925092565b600080604083850312156115f757600080fd5b8235915061160760208401611421565b90509250929050565b60006020828403121561162257600080fd5b610c3682611421565b6000806040838503121561163e57600080fd5b8235915060208301356001600160401b0381111561165b57600080fd5b8301601f8101851361166c57600080fd5b61167b858235602084016114d2565b9150509250929050565b6000806040838503121561169857600080fd5b6116a183611421565b915061160760208401611421565b80356001600160401b038116811461143857600080fd5b6000806000606084860312156116db57600080fd5b6116e484611421565b92506116f2602085016116af565b9150611700604085016116af565b90509250925092565b60006020828403121561171b57600080fd5b813561ffff81168114610c3657600080fd5b600181811c9082168061174157607f821691505b60208210810361176157634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b818103818111156104bb576104bb611767565b6001600160401b03841681526060602082015260006117b260608301856113c8565b9050826040830152949350505050565b6001600160401b03828116828216039080821115610dfc57610dfc611767565b80820281158282048414176104bb576104bb611767565b60008261181657634e487b7160e01b600052601260045260246000fd5b500490565b808201808211156104bb576104bb611767565b601f821115610754576000816000526020600020601f850160051c810160208610156118575750805b601f850160051c820191505b8181101561187657828155600101611863565b505050505050565b81516001600160401b03811115611897576118976114bc565b6118ab816118a5845461172d565b8461182e565b602080601f8311600181146118e057600084156118c85750858301515b600019600386901b1c1916600185901b178555611876565b600085815260208120601f198616915b8281101561190f578886015182559484019460019091019084016118f0565b508582101561192d5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fe52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0352c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0002dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800b458d13b0f4ce9e6aa65d297a27b10f75fdc6d0957bb29e1f2a30c8766b35415a2646970667358221220446a1716cd625489fe5bb84ebc9ef5a483e658b647524734d19ceedfd6d67e2c64736f6c63430008180033",
}

// ZenBTCABI is the input ABI used to generate the binding from.
// Deprecated: Use ZenBTCMetaData.ABI instead.
var ZenBTCABI = ZenBTCMetaData.ABI

// ZenBTCBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZenBTCMetaData.Bin instead.
var ZenBTCBin = ZenBTCMetaData.Bin

// DeployZenBTC deploys a new Ethereum contract, binding an instance of ZenBTC to it.
func DeployZenBTC(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ZenBTC, error) {
	parsed, err := ZenBTCMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZenBTCBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ZenBTC{ZenBTCCaller: ZenBTCCaller{contract: contract}, ZenBTCTransactor: ZenBTCTransactor{contract: contract}, ZenBTCFilterer: ZenBTCFilterer{contract: contract}}, nil
}

// ZenBTC is an auto generated Go binding around an Ethereum contract.
type ZenBTC struct {
	ZenBTCCaller     // Read-only binding to the contract
	ZenBTCTransactor // Write-only binding to the contract
	ZenBTCFilterer   // Log filterer for contract events
}

// ZenBTCCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZenBTCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBTCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZenBTCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBTCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZenBTCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBTCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZenBTCSession struct {
	Contract     *ZenBTC           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZenBTCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZenBTCCallerSession struct {
	Contract *ZenBTCCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ZenBTCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZenBTCTransactorSession struct {
	Contract     *ZenBTCTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZenBTCRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZenBTCRaw struct {
	Contract *ZenBTC // Generic contract binding to access the raw methods on
}

// ZenBTCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZenBTCCallerRaw struct {
	Contract *ZenBTCCaller // Generic read-only contract binding to access the raw methods on
}

// ZenBTCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZenBTCTransactorRaw struct {
	Contract *ZenBTCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZenBTC creates a new instance of ZenBTC, bound to a specific deployed contract.
func NewZenBTC(address common.Address, backend bind.ContractBackend) (*ZenBTC, error) {
	contract, err := bindZenBTC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZenBTC{ZenBTCCaller: ZenBTCCaller{contract: contract}, ZenBTCTransactor: ZenBTCTransactor{contract: contract}, ZenBTCFilterer: ZenBTCFilterer{contract: contract}}, nil
}

// NewZenBTCCaller creates a new read-only instance of ZenBTC, bound to a specific deployed contract.
func NewZenBTCCaller(address common.Address, caller bind.ContractCaller) (*ZenBTCCaller, error) {
	contract, err := bindZenBTC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZenBTCCaller{contract: contract}, nil
}

// NewZenBTCTransactor creates a new write-only instance of ZenBTC, bound to a specific deployed contract.
func NewZenBTCTransactor(address common.Address, transactor bind.ContractTransactor) (*ZenBTCTransactor, error) {
	contract, err := bindZenBTC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZenBTCTransactor{contract: contract}, nil
}

// NewZenBTCFilterer creates a new log filterer instance of ZenBTC, bound to a specific deployed contract.
func NewZenBTCFilterer(address common.Address, filterer bind.ContractFilterer) (*ZenBTCFilterer, error) {
	contract, err := bindZenBTC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZenBTCFilterer{contract: contract}, nil
}

// bindZenBTC binds a generic wrapper to an already deployed contract.
func bindZenBTC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZenBTCMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenBTC *ZenBTCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenBTC.Contract.ZenBTCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenBTC *ZenBTCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenBTC.Contract.ZenBTCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenBTC *ZenBTCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenBTC.Contract.ZenBTCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenBTC *ZenBTCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenBTC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenBTC *ZenBTCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenBTC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenBTC *ZenBTCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenBTC.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ZenBTC.Contract.DEFAULTADMINROLE(&_ZenBTC.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ZenBTC.Contract.DEFAULTADMINROLE(&_ZenBTC.CallOpts)
}

// FEEROLE is a free data retrieval call binding the contract method 0x2a0276f8.
//
// Solidity: function FEE_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCCaller) FEEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "FEE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FEEROLE is a free data retrieval call binding the contract method 0x2a0276f8.
//
// Solidity: function FEE_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCSession) FEEROLE() ([32]byte, error) {
	return _ZenBTC.Contract.FEEROLE(&_ZenBTC.CallOpts)
}

// FEEROLE is a free data retrieval call binding the contract method 0x2a0276f8.
//
// Solidity: function FEE_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCCallerSession) FEEROLE() ([32]byte, error) {
	return _ZenBTC.Contract.FEEROLE(&_ZenBTC.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCSession) MINTERROLE() ([32]byte, error) {
	return _ZenBTC.Contract.MINTERROLE(&_ZenBTC.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_ZenBTC *ZenBTCCallerSession) MINTERROLE() ([32]byte, error) {
	return _ZenBTC.Contract.MINTERROLE(&_ZenBTC.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBTC *ZenBTCCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBTC *ZenBTCSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZenBTC.Contract.Allowance(&_ZenBTC.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBTC *ZenBTCCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZenBTC.Contract.Allowance(&_ZenBTC.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBTC *ZenBTCCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBTC *ZenBTCSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZenBTC.Contract.BalanceOf(&_ZenBTC.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBTC *ZenBTCCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZenBTC.Contract.BalanceOf(&_ZenBTC.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBTC *ZenBTCCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBTC *ZenBTCSession) Decimals() (uint8, error) {
	return _ZenBTC.Contract.Decimals(&_ZenBTC.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBTC *ZenBTCCallerSession) Decimals() (uint8, error) {
	return _ZenBTC.Contract.Decimals(&_ZenBTC.CallOpts)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(uint256 value) view returns(uint256)
func (_ZenBTC *ZenBTCCaller) EstimateFee(opts *bind.CallOpts, value *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "estimateFee", value)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(uint256 value) view returns(uint256)
func (_ZenBTC *ZenBTCSession) EstimateFee(value *big.Int) (*big.Int, error) {
	return _ZenBTC.Contract.EstimateFee(&_ZenBTC.CallOpts, value)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(uint256 value) view returns(uint256)
func (_ZenBTC *ZenBTCCallerSession) EstimateFee(value *big.Int) (*big.Int, error) {
	return _ZenBTC.Contract.EstimateFee(&_ZenBTC.CallOpts, value)
}

// GetBasisPointsRate is a free data retrieval call binding the contract method 0xfded2529.
//
// Solidity: function getBasisPointsRate() view returns(uint16)
func (_ZenBTC *ZenBTCCaller) GetBasisPointsRate(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "getBasisPointsRate")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetBasisPointsRate is a free data retrieval call binding the contract method 0xfded2529.
//
// Solidity: function getBasisPointsRate() view returns(uint16)
func (_ZenBTC *ZenBTCSession) GetBasisPointsRate() (uint16, error) {
	return _ZenBTC.Contract.GetBasisPointsRate(&_ZenBTC.CallOpts)
}

// GetBasisPointsRate is a free data retrieval call binding the contract method 0xfded2529.
//
// Solidity: function getBasisPointsRate() view returns(uint16)
func (_ZenBTC *ZenBTCCallerSession) GetBasisPointsRate() (uint16, error) {
	return _ZenBTC.Contract.GetBasisPointsRate(&_ZenBTC.CallOpts)
}

// GetFeeAddress is a free data retrieval call binding the contract method 0x4e7ceacb.
//
// Solidity: function getFeeAddress() view returns(address)
func (_ZenBTC *ZenBTCCaller) GetFeeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "getFeeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFeeAddress is a free data retrieval call binding the contract method 0x4e7ceacb.
//
// Solidity: function getFeeAddress() view returns(address)
func (_ZenBTC *ZenBTCSession) GetFeeAddress() (common.Address, error) {
	return _ZenBTC.Contract.GetFeeAddress(&_ZenBTC.CallOpts)
}

// GetFeeAddress is a free data retrieval call binding the contract method 0x4e7ceacb.
//
// Solidity: function getFeeAddress() view returns(address)
func (_ZenBTC *ZenBTCCallerSession) GetFeeAddress() (common.Address, error) {
	return _ZenBTC.Contract.GetFeeAddress(&_ZenBTC.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ZenBTC *ZenBTCCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ZenBTC *ZenBTCSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ZenBTC.Contract.GetRoleAdmin(&_ZenBTC.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ZenBTC *ZenBTCCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ZenBTC.Contract.GetRoleAdmin(&_ZenBTC.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ZenBTC *ZenBTCCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ZenBTC *ZenBTCSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ZenBTC.Contract.HasRole(&_ZenBTC.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ZenBTC *ZenBTCCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ZenBTC.Contract.HasRole(&_ZenBTC.CallOpts, role, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBTC *ZenBTCCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBTC *ZenBTCSession) Name() (string, error) {
	return _ZenBTC.Contract.Name(&_ZenBTC.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBTC *ZenBTCCallerSession) Name() (string, error) {
	return _ZenBTC.Contract.Name(&_ZenBTC.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ZenBTC *ZenBTCCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ZenBTC *ZenBTCSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ZenBTC.Contract.SupportsInterface(&_ZenBTC.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ZenBTC *ZenBTCCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ZenBTC.Contract.SupportsInterface(&_ZenBTC.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBTC *ZenBTCCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBTC *ZenBTCSession) Symbol() (string, error) {
	return _ZenBTC.Contract.Symbol(&_ZenBTC.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBTC *ZenBTCCallerSession) Symbol() (string, error) {
	return _ZenBTC.Contract.Symbol(&_ZenBTC.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBTC *ZenBTCCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBTC *ZenBTCSession) TotalSupply() (*big.Int, error) {
	return _ZenBTC.Contract.TotalSupply(&_ZenBTC.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBTC *ZenBTCCallerSession) TotalSupply() (*big.Int, error) {
	return _ZenBTC.Contract.TotalSupply(&_ZenBTC.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.Approve(&_ZenBTC.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.Approve(&_ZenBTC.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_ZenBTC *ZenBTCTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_ZenBTC *ZenBTCSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.Burn(&_ZenBTC.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_ZenBTC *ZenBTCTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.Burn(&_ZenBTC.TransactOpts, value)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ZenBTC *ZenBTCTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ZenBTC *ZenBTCSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ZenBTC.Contract.GrantRole(&_ZenBTC.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ZenBTC *ZenBTCTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ZenBTC.Contract.GrantRole(&_ZenBTC.TransactOpts, role, account)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0x2a4b2fd9.
//
// Solidity: function initializeV1(string name_, string symbol_, uint8 decimals_) returns()
func (_ZenBTC *ZenBTCTransactor) InitializeV1(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "initializeV1", name_, symbol_, decimals_)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0x2a4b2fd9.
//
// Solidity: function initializeV1(string name_, string symbol_, uint8 decimals_) returns()
func (_ZenBTC *ZenBTCSession) InitializeV1(name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _ZenBTC.Contract.InitializeV1(&_ZenBTC.TransactOpts, name_, symbol_, decimals_)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0x2a4b2fd9.
//
// Solidity: function initializeV1(string name_, string symbol_, uint8 decimals_) returns()
func (_ZenBTC *ZenBTCTransactorSession) InitializeV1(name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _ZenBTC.Contract.InitializeV1(&_ZenBTC.TransactOpts, name_, symbol_, decimals_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 value) returns()
func (_ZenBTC *ZenBTCTransactor) Mint(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "mint", account, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 value) returns()
func (_ZenBTC *ZenBTCSession) Mint(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.Mint(&_ZenBTC.TransactOpts, account, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 value) returns()
func (_ZenBTC *ZenBTCTransactorSession) Mint(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.Mint(&_ZenBTC.TransactOpts, account, value)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ZenBTC *ZenBTCTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ZenBTC *ZenBTCSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ZenBTC.Contract.RenounceRole(&_ZenBTC.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ZenBTC *ZenBTCTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ZenBTC.Contract.RenounceRole(&_ZenBTC.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ZenBTC *ZenBTCTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ZenBTC *ZenBTCSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ZenBTC.Contract.RevokeRole(&_ZenBTC.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ZenBTC *ZenBTCTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ZenBTC.Contract.RevokeRole(&_ZenBTC.TransactOpts, role, account)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.Transfer(&_ZenBTC.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.Transfer(&_ZenBTC.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.TransferFrom(&_ZenBTC.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBTC *ZenBTCTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBTC.Contract.TransferFrom(&_ZenBTC.TransactOpts, from, to, value)
}

// Unwrap is a paid mutator transaction binding the contract method 0xb413148e.
//
// Solidity: function unwrap(uint256 value, bytes destAddr) returns()
func (_ZenBTC *ZenBTCTransactor) Unwrap(opts *bind.TransactOpts, value *big.Int, destAddr []byte) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "unwrap", value, destAddr)
}

// Unwrap is a paid mutator transaction binding the contract method 0xb413148e.
//
// Solidity: function unwrap(uint256 value, bytes destAddr) returns()
func (_ZenBTC *ZenBTCSession) Unwrap(value *big.Int, destAddr []byte) (*types.Transaction, error) {
	return _ZenBTC.Contract.Unwrap(&_ZenBTC.TransactOpts, value, destAddr)
}

// Unwrap is a paid mutator transaction binding the contract method 0xb413148e.
//
// Solidity: function unwrap(uint256 value, bytes destAddr) returns()
func (_ZenBTC *ZenBTCTransactorSession) Unwrap(value *big.Int, destAddr []byte) (*types.Transaction, error) {
	return _ZenBTC.Contract.Unwrap(&_ZenBTC.TransactOpts, value, destAddr)
}

// UpdateBasisPointsRate is a paid mutator transaction binding the contract method 0xfd154cc7.
//
// Solidity: function updateBasisPointsRate(uint16 basisPointsRate_) returns()
func (_ZenBTC *ZenBTCTransactor) UpdateBasisPointsRate(opts *bind.TransactOpts, basisPointsRate_ uint16) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "updateBasisPointsRate", basisPointsRate_)
}

// UpdateBasisPointsRate is a paid mutator transaction binding the contract method 0xfd154cc7.
//
// Solidity: function updateBasisPointsRate(uint16 basisPointsRate_) returns()
func (_ZenBTC *ZenBTCSession) UpdateBasisPointsRate(basisPointsRate_ uint16) (*types.Transaction, error) {
	return _ZenBTC.Contract.UpdateBasisPointsRate(&_ZenBTC.TransactOpts, basisPointsRate_)
}

// UpdateBasisPointsRate is a paid mutator transaction binding the contract method 0xfd154cc7.
//
// Solidity: function updateBasisPointsRate(uint16 basisPointsRate_) returns()
func (_ZenBTC *ZenBTCTransactorSession) UpdateBasisPointsRate(basisPointsRate_ uint16) (*types.Transaction, error) {
	return _ZenBTC.Contract.UpdateBasisPointsRate(&_ZenBTC.TransactOpts, basisPointsRate_)
}

// UpdateFeeAddress is a paid mutator transaction binding the contract method 0xbbcaac38.
//
// Solidity: function updateFeeAddress(address feeAddress_) returns()
func (_ZenBTC *ZenBTCTransactor) UpdateFeeAddress(opts *bind.TransactOpts, feeAddress_ common.Address) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "updateFeeAddress", feeAddress_)
}

// UpdateFeeAddress is a paid mutator transaction binding the contract method 0xbbcaac38.
//
// Solidity: function updateFeeAddress(address feeAddress_) returns()
func (_ZenBTC *ZenBTCSession) UpdateFeeAddress(feeAddress_ common.Address) (*types.Transaction, error) {
	return _ZenBTC.Contract.UpdateFeeAddress(&_ZenBTC.TransactOpts, feeAddress_)
}

// UpdateFeeAddress is a paid mutator transaction binding the contract method 0xbbcaac38.
//
// Solidity: function updateFeeAddress(address feeAddress_) returns()
func (_ZenBTC *ZenBTCTransactorSession) UpdateFeeAddress(feeAddress_ common.Address) (*types.Transaction, error) {
	return _ZenBTC.Contract.UpdateFeeAddress(&_ZenBTC.TransactOpts, feeAddress_)
}

// Wrap is a paid mutator transaction binding the contract method 0xf1d96dd3.
//
// Solidity: function wrap(address account, uint64 value, uint64 fee) returns()
func (_ZenBTC *ZenBTCTransactor) Wrap(opts *bind.TransactOpts, account common.Address, value uint64, fee uint64) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "wrap", account, value, fee)
}

// Wrap is a paid mutator transaction binding the contract method 0xf1d96dd3.
//
// Solidity: function wrap(address account, uint64 value, uint64 fee) returns()
func (_ZenBTC *ZenBTCSession) Wrap(account common.Address, value uint64, fee uint64) (*types.Transaction, error) {
	return _ZenBTC.Contract.Wrap(&_ZenBTC.TransactOpts, account, value, fee)
}

// Wrap is a paid mutator transaction binding the contract method 0xf1d96dd3.
//
// Solidity: function wrap(address account, uint64 value, uint64 fee) returns()
func (_ZenBTC *ZenBTCTransactorSession) Wrap(account common.Address, value uint64, fee uint64) (*types.Transaction, error) {
	return _ZenBTC.Contract.Wrap(&_ZenBTC.TransactOpts, account, value, fee)
}

// ZenBTCApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ZenBTC contract.
type ZenBTCApprovalIterator struct {
	Event *ZenBTCApproval // Event containing the contract specifics and raw log

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
func (it *ZenBTCApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCApproval)
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
		it.Event = new(ZenBTCApproval)
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
func (it *ZenBTCApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCApproval represents a Approval event raised by the ZenBTC contract.
type ZenBTCApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenBTC *ZenBTCFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ZenBTCApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ZenBTCApprovalIterator{contract: _ZenBTC.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenBTC *ZenBTCFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZenBTCApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCApproval)
				if err := _ZenBTC.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseApproval(log types.Log) (*ZenBTCApproval, error) {
	event := new(ZenBTCApproval)
	if err := _ZenBTC.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCBasisPointsRateUpdatedIterator is returned from FilterBasisPointsRateUpdated and is used to iterate over the raw logs and unpacked data for BasisPointsRateUpdated events raised by the ZenBTC contract.
type ZenBTCBasisPointsRateUpdatedIterator struct {
	Event *ZenBTCBasisPointsRateUpdated // Event containing the contract specifics and raw log

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
func (it *ZenBTCBasisPointsRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCBasisPointsRateUpdated)
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
		it.Event = new(ZenBTCBasisPointsRateUpdated)
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
func (it *ZenBTCBasisPointsRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCBasisPointsRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCBasisPointsRateUpdated represents a BasisPointsRateUpdated event raised by the ZenBTC contract.
type ZenBTCBasisPointsRateUpdated struct {
	BasisPointsRateUpdated uint16
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterBasisPointsRateUpdated is a free log retrieval operation binding the contract event 0x79d2a06c43b232cf9d1835b5e915efe74621561aaf75fecec91fd75a940f9d70.
//
// Solidity: event BasisPointsRateUpdated(uint16 basisPointsRateUpdated)
func (_ZenBTC *ZenBTCFilterer) FilterBasisPointsRateUpdated(opts *bind.FilterOpts) (*ZenBTCBasisPointsRateUpdatedIterator, error) {

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "BasisPointsRateUpdated")
	if err != nil {
		return nil, err
	}
	return &ZenBTCBasisPointsRateUpdatedIterator{contract: _ZenBTC.contract, event: "BasisPointsRateUpdated", logs: logs, sub: sub}, nil
}

// WatchBasisPointsRateUpdated is a free log subscription operation binding the contract event 0x79d2a06c43b232cf9d1835b5e915efe74621561aaf75fecec91fd75a940f9d70.
//
// Solidity: event BasisPointsRateUpdated(uint16 basisPointsRateUpdated)
func (_ZenBTC *ZenBTCFilterer) WatchBasisPointsRateUpdated(opts *bind.WatchOpts, sink chan<- *ZenBTCBasisPointsRateUpdated) (event.Subscription, error) {

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "BasisPointsRateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCBasisPointsRateUpdated)
				if err := _ZenBTC.contract.UnpackLog(event, "BasisPointsRateUpdated", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseBasisPointsRateUpdated(log types.Log) (*ZenBTCBasisPointsRateUpdated, error) {
	event := new(ZenBTCBasisPointsRateUpdated)
	if err := _ZenBTC.contract.UnpackLog(event, "BasisPointsRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCFeeAddressUpdatedIterator is returned from FilterFeeAddressUpdated and is used to iterate over the raw logs and unpacked data for FeeAddressUpdated events raised by the ZenBTC contract.
type ZenBTCFeeAddressUpdatedIterator struct {
	Event *ZenBTCFeeAddressUpdated // Event containing the contract specifics and raw log

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
func (it *ZenBTCFeeAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCFeeAddressUpdated)
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
		it.Event = new(ZenBTCFeeAddressUpdated)
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
func (it *ZenBTCFeeAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCFeeAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCFeeAddressUpdated represents a FeeAddressUpdated event raised by the ZenBTC contract.
type ZenBTCFeeAddressUpdated struct {
	NewFeeAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterFeeAddressUpdated is a free log retrieval operation binding the contract event 0x446e39bcf1b47cfadfaa23442cb4b34682cfe6bd9220da084894e3b1f834e4f3.
//
// Solidity: event FeeAddressUpdated(address indexed newFeeAddress)
func (_ZenBTC *ZenBTCFilterer) FilterFeeAddressUpdated(opts *bind.FilterOpts, newFeeAddress []common.Address) (*ZenBTCFeeAddressUpdatedIterator, error) {

	var newFeeAddressRule []interface{}
	for _, newFeeAddressItem := range newFeeAddress {
		newFeeAddressRule = append(newFeeAddressRule, newFeeAddressItem)
	}

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "FeeAddressUpdated", newFeeAddressRule)
	if err != nil {
		return nil, err
	}
	return &ZenBTCFeeAddressUpdatedIterator{contract: _ZenBTC.contract, event: "FeeAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchFeeAddressUpdated is a free log subscription operation binding the contract event 0x446e39bcf1b47cfadfaa23442cb4b34682cfe6bd9220da084894e3b1f834e4f3.
//
// Solidity: event FeeAddressUpdated(address indexed newFeeAddress)
func (_ZenBTC *ZenBTCFilterer) WatchFeeAddressUpdated(opts *bind.WatchOpts, sink chan<- *ZenBTCFeeAddressUpdated, newFeeAddress []common.Address) (event.Subscription, error) {

	var newFeeAddressRule []interface{}
	for _, newFeeAddressItem := range newFeeAddress {
		newFeeAddressRule = append(newFeeAddressRule, newFeeAddressItem)
	}

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "FeeAddressUpdated", newFeeAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCFeeAddressUpdated)
				if err := _ZenBTC.contract.UnpackLog(event, "FeeAddressUpdated", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseFeeAddressUpdated(log types.Log) (*ZenBTCFeeAddressUpdated, error) {
	event := new(ZenBTCFeeAddressUpdated)
	if err := _ZenBTC.contract.UnpackLog(event, "FeeAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ZenBTC contract.
type ZenBTCInitializedIterator struct {
	Event *ZenBTCInitialized // Event containing the contract specifics and raw log

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
func (it *ZenBTCInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCInitialized)
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
		it.Event = new(ZenBTCInitialized)
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
func (it *ZenBTCInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCInitialized represents a Initialized event raised by the ZenBTC contract.
type ZenBTCInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ZenBTC *ZenBTCFilterer) FilterInitialized(opts *bind.FilterOpts) (*ZenBTCInitializedIterator, error) {

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ZenBTCInitializedIterator{contract: _ZenBTC.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ZenBTC *ZenBTCFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ZenBTCInitialized) (event.Subscription, error) {

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCInitialized)
				if err := _ZenBTC.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseInitialized(log types.Log) (*ZenBTCInitialized, error) {
	event := new(ZenBTCInitialized)
	if err := _ZenBTC.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCRedeemFeeUpdatedIterator is returned from FilterRedeemFeeUpdated and is used to iterate over the raw logs and unpacked data for RedeemFeeUpdated events raised by the ZenBTC contract.
type ZenBTCRedeemFeeUpdatedIterator struct {
	Event *ZenBTCRedeemFeeUpdated // Event containing the contract specifics and raw log

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
func (it *ZenBTCRedeemFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCRedeemFeeUpdated)
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
		it.Event = new(ZenBTCRedeemFeeUpdated)
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
func (it *ZenBTCRedeemFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCRedeemFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCRedeemFeeUpdated represents a RedeemFeeUpdated event raised by the ZenBTC contract.
type ZenBTCRedeemFeeUpdated struct {
	NewRedeemFee *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRedeemFeeUpdated is a free log retrieval operation binding the contract event 0xd6c7508d6658ccee36b7b7d7fd72e5cbaeefb40c64eff24e9ae7470e846304ee.
//
// Solidity: event RedeemFeeUpdated(uint256 newRedeemFee)
func (_ZenBTC *ZenBTCFilterer) FilterRedeemFeeUpdated(opts *bind.FilterOpts) (*ZenBTCRedeemFeeUpdatedIterator, error) {

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "RedeemFeeUpdated")
	if err != nil {
		return nil, err
	}
	return &ZenBTCRedeemFeeUpdatedIterator{contract: _ZenBTC.contract, event: "RedeemFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchRedeemFeeUpdated is a free log subscription operation binding the contract event 0xd6c7508d6658ccee36b7b7d7fd72e5cbaeefb40c64eff24e9ae7470e846304ee.
//
// Solidity: event RedeemFeeUpdated(uint256 newRedeemFee)
func (_ZenBTC *ZenBTCFilterer) WatchRedeemFeeUpdated(opts *bind.WatchOpts, sink chan<- *ZenBTCRedeemFeeUpdated) (event.Subscription, error) {

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "RedeemFeeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCRedeemFeeUpdated)
				if err := _ZenBTC.contract.UnpackLog(event, "RedeemFeeUpdated", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseRedeemFeeUpdated(log types.Log) (*ZenBTCRedeemFeeUpdated, error) {
	event := new(ZenBTCRedeemFeeUpdated)
	if err := _ZenBTC.contract.UnpackLog(event, "RedeemFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ZenBTC contract.
type ZenBTCRoleAdminChangedIterator struct {
	Event *ZenBTCRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ZenBTCRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCRoleAdminChanged)
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
		it.Event = new(ZenBTCRoleAdminChanged)
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
func (it *ZenBTCRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCRoleAdminChanged represents a RoleAdminChanged event raised by the ZenBTC contract.
type ZenBTCRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ZenBTC *ZenBTCFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ZenBTCRoleAdminChangedIterator, error) {

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

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ZenBTCRoleAdminChangedIterator{contract: _ZenBTC.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ZenBTC *ZenBTCFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ZenBTCRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCRoleAdminChanged)
				if err := _ZenBTC.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseRoleAdminChanged(log types.Log) (*ZenBTCRoleAdminChanged, error) {
	event := new(ZenBTCRoleAdminChanged)
	if err := _ZenBTC.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ZenBTC contract.
type ZenBTCRoleGrantedIterator struct {
	Event *ZenBTCRoleGranted // Event containing the contract specifics and raw log

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
func (it *ZenBTCRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCRoleGranted)
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
		it.Event = new(ZenBTCRoleGranted)
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
func (it *ZenBTCRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCRoleGranted represents a RoleGranted event raised by the ZenBTC contract.
type ZenBTCRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ZenBTC *ZenBTCFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ZenBTCRoleGrantedIterator, error) {

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

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ZenBTCRoleGrantedIterator{contract: _ZenBTC.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ZenBTC *ZenBTCFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ZenBTCRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCRoleGranted)
				if err := _ZenBTC.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseRoleGranted(log types.Log) (*ZenBTCRoleGranted, error) {
	event := new(ZenBTCRoleGranted)
	if err := _ZenBTC.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ZenBTC contract.
type ZenBTCRoleRevokedIterator struct {
	Event *ZenBTCRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ZenBTCRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCRoleRevoked)
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
		it.Event = new(ZenBTCRoleRevoked)
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
func (it *ZenBTCRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCRoleRevoked represents a RoleRevoked event raised by the ZenBTC contract.
type ZenBTCRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ZenBTC *ZenBTCFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ZenBTCRoleRevokedIterator, error) {

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

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ZenBTCRoleRevokedIterator{contract: _ZenBTC.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ZenBTC *ZenBTCFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ZenBTCRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCRoleRevoked)
				if err := _ZenBTC.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseRoleRevoked(log types.Log) (*ZenBTCRoleRevoked, error) {
	event := new(ZenBTCRoleRevoked)
	if err := _ZenBTC.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCTokenRedemptionIterator is returned from FilterTokenRedemption and is used to iterate over the raw logs and unpacked data for TokenRedemption events raised by the ZenBTC contract.
type ZenBTCTokenRedemptionIterator struct {
	Event *ZenBTCTokenRedemption // Event containing the contract specifics and raw log

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
func (it *ZenBTCTokenRedemptionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCTokenRedemption)
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
		it.Event = new(ZenBTCTokenRedemption)
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
func (it *ZenBTCTokenRedemptionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCTokenRedemptionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCTokenRedemption represents a TokenRedemption event raised by the ZenBTC contract.
type ZenBTCTokenRedemption struct {
	Redeemer common.Address
	Value    uint64
	DestAddr []byte
	Fee      *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokenRedemption is a free log retrieval operation binding the contract event 0x4c971c8b2abb197a17896b2fc57f597830db78d3556831e3faa337a596150f22.
//
// Solidity: event TokenRedemption(address indexed redeemer, uint64 value, bytes destAddr, uint256 fee)
func (_ZenBTC *ZenBTCFilterer) FilterTokenRedemption(opts *bind.FilterOpts, redeemer []common.Address) (*ZenBTCTokenRedemptionIterator, error) {

	var redeemerRule []interface{}
	for _, redeemerItem := range redeemer {
		redeemerRule = append(redeemerRule, redeemerItem)
	}

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "TokenRedemption", redeemerRule)
	if err != nil {
		return nil, err
	}
	return &ZenBTCTokenRedemptionIterator{contract: _ZenBTC.contract, event: "TokenRedemption", logs: logs, sub: sub}, nil
}

// WatchTokenRedemption is a free log subscription operation binding the contract event 0x4c971c8b2abb197a17896b2fc57f597830db78d3556831e3faa337a596150f22.
//
// Solidity: event TokenRedemption(address indexed redeemer, uint64 value, bytes destAddr, uint256 fee)
func (_ZenBTC *ZenBTCFilterer) WatchTokenRedemption(opts *bind.WatchOpts, sink chan<- *ZenBTCTokenRedemption, redeemer []common.Address) (event.Subscription, error) {

	var redeemerRule []interface{}
	for _, redeemerItem := range redeemer {
		redeemerRule = append(redeemerRule, redeemerItem)
	}

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "TokenRedemption", redeemerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCTokenRedemption)
				if err := _ZenBTC.contract.UnpackLog(event, "TokenRedemption", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseTokenRedemption(log types.Log) (*ZenBTCTokenRedemption, error) {
	event := new(ZenBTCTokenRedemption)
	if err := _ZenBTC.contract.UnpackLog(event, "TokenRedemption", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCTokensMintedWithFeeIterator is returned from FilterTokensMintedWithFee and is used to iterate over the raw logs and unpacked data for TokensMintedWithFee events raised by the ZenBTC contract.
type ZenBTCTokensMintedWithFeeIterator struct {
	Event *ZenBTCTokensMintedWithFee // Event containing the contract specifics and raw log

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
func (it *ZenBTCTokensMintedWithFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCTokensMintedWithFee)
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
		it.Event = new(ZenBTCTokensMintedWithFee)
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
func (it *ZenBTCTokensMintedWithFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCTokensMintedWithFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCTokensMintedWithFee represents a TokensMintedWithFee event raised by the ZenBTC contract.
type ZenBTCTokensMintedWithFee struct {
	Recipient common.Address
	Value     uint64
	Fee       uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokensMintedWithFee is a free log retrieval operation binding the contract event 0x890f2fd0578a1ba6f8735fc60c94d1ec453ad71501d82ba4d5a9041f9ac14f2e.
//
// Solidity: event TokensMintedWithFee(address indexed recipient, uint64 value, uint64 fee)
func (_ZenBTC *ZenBTCFilterer) FilterTokensMintedWithFee(opts *bind.FilterOpts, recipient []common.Address) (*ZenBTCTokensMintedWithFeeIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "TokensMintedWithFee", recipientRule)
	if err != nil {
		return nil, err
	}
	return &ZenBTCTokensMintedWithFeeIterator{contract: _ZenBTC.contract, event: "TokensMintedWithFee", logs: logs, sub: sub}, nil
}

// WatchTokensMintedWithFee is a free log subscription operation binding the contract event 0x890f2fd0578a1ba6f8735fc60c94d1ec453ad71501d82ba4d5a9041f9ac14f2e.
//
// Solidity: event TokensMintedWithFee(address indexed recipient, uint64 value, uint64 fee)
func (_ZenBTC *ZenBTCFilterer) WatchTokensMintedWithFee(opts *bind.WatchOpts, sink chan<- *ZenBTCTokensMintedWithFee, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "TokensMintedWithFee", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCTokensMintedWithFee)
				if err := _ZenBTC.contract.UnpackLog(event, "TokensMintedWithFee", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseTokensMintedWithFee(log types.Log) (*ZenBTCTokensMintedWithFee, error) {
	event := new(ZenBTCTokensMintedWithFee)
	if err := _ZenBTC.contract.UnpackLog(event, "TokensMintedWithFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ZenBTC contract.
type ZenBTCTransferIterator struct {
	Event *ZenBTCTransfer // Event containing the contract specifics and raw log

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
func (it *ZenBTCTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCTransfer)
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
		it.Event = new(ZenBTCTransfer)
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
func (it *ZenBTCTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCTransfer represents a Transfer event raised by the ZenBTC contract.
type ZenBTCTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenBTC *ZenBTCFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ZenBTCTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ZenBTCTransferIterator{contract: _ZenBTC.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenBTC *ZenBTCFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZenBTCTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCTransfer)
				if err := _ZenBTC.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseTransfer(log types.Log) (*ZenBTCTransfer, error) {
	event := new(ZenBTCTransfer)
	if err := _ZenBTC.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBTCZenBTControllerUpdatedIterator is returned from FilterZenBTControllerUpdated and is used to iterate over the raw logs and unpacked data for ZenBTControllerUpdated events raised by the ZenBTC contract.
type ZenBTCZenBTControllerUpdatedIterator struct {
	Event *ZenBTCZenBTControllerUpdated // Event containing the contract specifics and raw log

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
func (it *ZenBTCZenBTControllerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBTCZenBTControllerUpdated)
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
		it.Event = new(ZenBTCZenBTControllerUpdated)
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
func (it *ZenBTCZenBTControllerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBTCZenBTControllerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBTCZenBTControllerUpdated represents a ZenBTControllerUpdated event raised by the ZenBTC contract.
type ZenBTCZenBTControllerUpdated struct {
	ZenBTControllerUpdated common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterZenBTControllerUpdated is a free log retrieval operation binding the contract event 0xc4c1343b8690c9be210951438cee3b760df1a37031ac0799c086b90d6af53950.
//
// Solidity: event ZenBTControllerUpdated(address zenBTControllerUpdated)
func (_ZenBTC *ZenBTCFilterer) FilterZenBTControllerUpdated(opts *bind.FilterOpts) (*ZenBTCZenBTControllerUpdatedIterator, error) {

	logs, sub, err := _ZenBTC.contract.FilterLogs(opts, "ZenBTControllerUpdated")
	if err != nil {
		return nil, err
	}
	return &ZenBTCZenBTControllerUpdatedIterator{contract: _ZenBTC.contract, event: "ZenBTControllerUpdated", logs: logs, sub: sub}, nil
}

// WatchZenBTControllerUpdated is a free log subscription operation binding the contract event 0xc4c1343b8690c9be210951438cee3b760df1a37031ac0799c086b90d6af53950.
//
// Solidity: event ZenBTControllerUpdated(address zenBTControllerUpdated)
func (_ZenBTC *ZenBTCFilterer) WatchZenBTControllerUpdated(opts *bind.WatchOpts, sink chan<- *ZenBTCZenBTControllerUpdated) (event.Subscription, error) {

	logs, sub, err := _ZenBTC.contract.WatchLogs(opts, "ZenBTControllerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBTCZenBTControllerUpdated)
				if err := _ZenBTC.contract.UnpackLog(event, "ZenBTControllerUpdated", log); err != nil {
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
func (_ZenBTC *ZenBTCFilterer) ParseZenBTControllerUpdated(log types.Log) (*ZenBTCZenBTControllerUpdated, error) {
	event := new(ZenBTCZenBTControllerUpdated)
	if err := _ZenBTC.contract.UnpackLog(event, "ZenBTControllerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}