// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newFeeAddress\",\"type\":\"address\"}],\"name\":\"FeeAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newRedeemFee\",\"type\":\"uint256\"}],\"name\":\"RedeemFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"destAddr\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"}],\"name\":\"TokenRedemption\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"}],\"name\":\"TokensMintedWithFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numBlocks\",\"type\":\"uint16\"}],\"name\":\"getRecentRedemptionData\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"ids\",\"type\":\"uint64[]\"},{\"internalType\":\"uint64[]\",\"name\":\"amounts\",\"type\":\"uint64[]\"},{\"internalType\":\"bytes[]\",\"name\":\"destination\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRedeemFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"}],\"name\":\"initializeV1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"destAddr\",\"type\":\"bytes\"}],\"name\":\"unwrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feeAddress_\",\"type\":\"address\"}],\"name\":\"updateFeeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"redeemFee_\",\"type\":\"uint64\"}],\"name\":\"updateRedeemFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"}],\"name\":\"wrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600436106101c45760003560e01c80634e7ceacb116100f9578063c6d98f1a11610097578063d547741f11610071578063d547741f1461042b578063dd62ed3e1461043e578063f1690dbe14610451578063f1d96dd31461047357600080fd5b8063c6d98f1a146103de578063ce29ea6214610403578063d53913931461041657600080fd5b806395d89b41116100d357806395d89b41146103a8578063a217fddf146103b0578063a9059cbb146103b8578063bbcaac38146103cb57600080fd5b80634e7ceacb1461033757806370a082311461035f57806391d148541461039557600080fd5b80632a4b2fd91161016657806336568abe1161014057806336568abe146102eb57806340c10f19146102fe57806342966c68146103115780634677d2791461032457600080fd5b80632a4b2fd9146102975780632f2ff15d146102ac578063313ce567146102bf57600080fd5b806318160ddd116101a257806318160ddd1461021957806323b872dd1461024a578063248a9ca31461025d5780632a0276f81461027057600080fd5b806301ffc9a7146101c957806306fdde03146101f1578063095ea7b314610206575b600080fd5b6101dc6101d7366004611b18565b610486565b60405190151581526020015b60405180910390f35b6101f96104bd565b6040516101e89190611b8f565b6101dc610214366004611bbe565b610580565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace02545b6040519081526020016101e8565b6101dc610258366004611be8565b610598565b61023c61026b366004611c24565b6105bc565b61023c7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf2881565b6102aa6102a5366004611ce8565b6105de565b005b6102aa6102ba366004611d65565b6106f8565b6000805160206121b783398151915254600160e01b900460ff1660405160ff90911681526020016101e8565b6102aa6102f9366004611d65565b61071a565b6102aa61030c366004611bbe565b610752565b6102aa61031f366004611c24565b610774565b6102aa610332366004611da8565b61079a565b6000805160206121b7833981519152546040516001600160a01b0390911681526020016101e8565b61023c61036d366004611e09565b6001600160a01b031660009081526000805160206121d7833981519152602052604090205490565b6101dc6103a3366004611d65565b61092f565b6101f9610967565b61023c600081565b6101dc6103c6366004611bbe565b6109a6565b6102aa6103d9366004611e09565b6109b4565b6000805160206121b783398151915254600160a01b90046001600160401b031661023c565b6102aa610411366004611e24565b6109e7565b61023c60008051602061223783398151915281565b6102aa610439366004611d65565b610a1a565b61023c61044c366004611e3f565b610a36565b61046461045f366004611e69565b610a80565b6040516101e893929190611ed2565b6102aa610481366004611f59565b610eef565b60006001600160e01b03198216637965db0b60e01b14806104b757506301ffc9a760e01b6001600160e01b03198316145b92915050565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0480546060916000805160206121d7833981519152916104fc90611f9c565b80601f016020809104026020016040519081016040528092919081815260200182805461052890611f9c565b80156105755780601f1061054a57610100808354040283529160200191610575565b820191906000526020600020905b81548152906001019060200180831161055857829003601f168201915b505050505091505090565b60003361058e81858561106a565b5060019392505050565b6000336105a6858285611077565b6105b18585856110d7565b506001949350505050565b6000908152600080516020612217833981519152602052604090206001015490565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff1615906001600160401b03166000811580156106235750825b90506000826001600160401b0316600114801561063f5750303b155b90508115801561064d575080155b1561066b5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561069557845460ff60401b1916600160401b1785555b6106a0888888611161565b6106a8611184565b83156106ee57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050505050565b610701826105bc565b61070a81611196565b61071483836111a3565b50505050565b6001600160a01b03811633146107435760405163334bd91960e11b815260040160405180910390fd5b61074d828261124f565b505050565b60008051602061223783398151915261076a81611196565b61074d83836112cb565b60008051602061223783398151915261078c81611196565b6107963383611301565b5050565b6000805160206121d78339815191526001600160401b0383811611156107f95760405162461bcd60e51b815260206004820152600f60248201526e56616c756520746f6f206c6172676560881b60448201526064015b60405180910390fd5b60038101546001600160401b03808516600160a01b909204161061086d5760405162461bcd60e51b815260206004820152602560248201527f5a656e4254433a2072656465656d2066656520657863656564732074686520616044820152641b5bdd5b9d60da1b60648201526084016107f0565b61088033846001600160401b0316611301565b60038101546108a9906001600160a01b03811690600160a01b90046001600160401b03166112cb565b60038101546000906108cb90600160a01b90046001600160401b031685611fec565b90506108d78184611337565b600382015460405133917f4f8266db69b83687df444f19eab599803494f35fe2a1c061ba66d871bda06824916109219185918891600160a01b90046001600160401b03169061200c565b60405180910390a250505050565b6000918252600080516020612217833981519152602090815260408084206001600160a01b0393909316845291905290205460ff1690565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0580546060916000805160206121d7833981519152916104fc90611f9c565b60003361058e8185856110d7565b7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf286109de81611196565b61079682611514565b7f13b7ad447453d194d272cdda9bb09d7d357cda1ab7de80d865b4c1cbefc3cf28610a1181611196565b6107968261157d565b610a23826105bc565b610a2c81611196565b610714838361124f565b6001600160a01b0391821660009081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace016020908152604080832093909416825291909152205490565b60608080606461ffff85161115610acb5760405162461bcd60e51b815260206004820152600f60248201526e52616e676520746f6f206c6172676560881b60448201526064016107f0565b6000805160206121f78339815191526000805b8661ffff168161ffff161015610b7557825460009060649061ffff908116849003606301160690506000846001018261ffff1681548110610b2157610b21612058565b6000918252602090912060088204015460079091166004026101000a900463ffffffff1690508015610b6b5763ffffffff8116600090815260028601602052604090205493909301925b5050600101610ade565b50806001600160401b03811115610b8e57610b8e611c3d565b604051908082528060200260200182016040528015610bb7578160200160208202803683370190505b509450806001600160401b03811115610bd257610bd2611c3d565b604051908082528060200260200182016040528015610bfb578160200160208202803683370190505b509350806001600160401b03811115610c1657610c16611c3d565b604051908082528060200260200182016040528015610c4957816020015b6060815260200190600190039081610c345790505b5092506000805b8761ffff168161ffff161015610ee457835460009060649061ffff908116849003606301160690506000856001018261ffff1681548110610c9357610c93612058565b6000918252602090912060088204015460079091166004026101000a900463ffffffff1690508015610eda5763ffffffff81166000908152600287016020908152604080832080548251818502810185019093528083529192909190849084015b82821015610dde576000848152602090819020604080516060810182526002860290920180546001600160401b038082168552600160401b909104169383019390935260018301805492939291840191610d4d90611f9c565b80601f0160208091040260200160405190810160405280929190818152602001828054610d7990611f9c565b8015610dc65780601f10610d9b57610100808354040283529160200191610dc6565b820191906000526020600020905b815481529060010190602001808311610da957829003601f168201915b50505050508152505081526020019060010190610cf4565b50505050905060005b8151811015610ed757818181518110610e0257610e02612058565b6020026020010151600001518b8781518110610e2057610e20612058565b60200260200101906001600160401b031690816001600160401b031681525050818181518110610e5257610e52612058565b6020026020010151602001518a8781518110610e7057610e70612058565b60200260200101906001600160401b031690816001600160401b031681525050818181518110610ea257610ea2612058565b602002602001015160400151898781518110610ec057610ec0612058565b602090810291909101015260019586019501610de7565b50505b5050600101610c50565b505050509193909250565b600080516020612237833981519152610f0781611196565b6001600160401b038381161115610f525760405162461bcd60e51b815260206004820152600f60248201526e56616c756520746f6f206c6172676560881b60448201526064016107f0565b826001600160401b0316826001600160401b031610610fc25760405162461bcd60e51b815260206004820152602660248201527f5a656e4254433a2046656520657863656564732074686520616d6f756e7420746044820152651bc81b5a5b9d60d21b60648201526084016107f0565b6000805160206121d78339815191526000610fdd8486611fec565b6003830154909150611001906001600160a01b03166001600160401b0386166112cb565b61101486826001600160401b03166112cb565b604080516001600160401b038084168252861660208201526001600160a01b038816917f890f2fd0578a1ba6f8735fc60c94d1ec453ad71501d82ba4d5a9041f9ac14f2e910160405180910390a2505050505050565b61074d83838360016115fd565b60006110838484610a36565b9050600019811461071457818110156110c857604051637dc7a0d960e11b81526001600160a01b038416600482015260248101829052604481018390526064016107f0565b610714848484840360006115fd565b6001600160a01b03831661110157604051634b637e8f60e11b8152600060048201526024016107f0565b6001600160a01b03821661112b5760405163ec442f0560e01b8152600060048201526024016107f0565b306001600160a01b038316036111565760405163ec442f0560e01b81523060048201526024016107f0565b61074d8383836116e5565b611169611823565b61117161186c565b611179611874565b61074d838383611884565b61118c611823565b6111946118f9565b565b6111a08133611944565b50565b60006000805160206122178339815191526111be848461092f565b61123e576000848152602082815260408083206001600160a01b03871684529091529020805460ff191660011790556111f43390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019150506104b7565b60009150506104b7565b5092915050565b600060008051602061221783398151915261126a848461092f565b1561123e576000848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a460019150506104b7565b6001600160a01b0382166112f55760405163ec442f0560e01b8152600060048201526024016107f0565b610796600083836116e5565b6001600160a01b03821661132b57604051634b637e8f60e11b8152600060048201526024016107f0565b610796826000836116e5565b4363ffffffff811660009081527fc87a15b711c3ce3a8db8c4c853ae9c01bdbe93da3d2e98c03c16ae2fd0346f026020526040812080546000805160206121f783398151915293920361145d5782546001840180546113ce9261ffff169081106113a3576113a3612058565b90600052602060002090600891828204019190066004029054906101000a900463ffffffff1661197d565b8254600184018054849261ffff169081106113eb576113eb612058565b90600052602060002090600891828204019190066004026101000a81548163ffffffff021916908363ffffffff160217905550606461ffff168360000160009054906101000a900461ffff1660010161ffff168161144b5761144b612042565b845461ffff191691900661ffff161783555b82546201000090046001600160401b031683600261147a8361206e565b82546001600160401b039182166101009390930a928302928202191691909117909155604080516060810182528654620100009004831681528883166020808301918252928201898152865460018181018955600089815295909520845160029290920201805493518716600160401b026001600160801b0319909416919096161791909117845551909350908201906106ee90826120e4565b6000805160206121b783398151915280546001600160a01b0319166001600160a01b0383169081179091556040516000805160206121d783398151915291907f446e39bcf1b47cfadfaa23442cb4b34682cfe6bd9220da084894e3b1f834e4f390600090a25050565b6000805160206121b783398151915280546001600160401b038316600160a01b810267ffffffffffffffff60a01b199092169190911790915560408051918252516000805160206121d7833981519152917fd6c7508d6658ccee36b7b7d7fd72e5cbaeefb40c64eff24e9ae7470e846304ee919081900360200190a15050565b6000805160206121d78339815191526001600160a01b0385166116365760405163e602df0560e01b8152600060048201526024016107f0565b6001600160a01b03841661166057604051634a1406b160e11b8152600060048201526024016107f0565b6001600160a01b038086166000908152600183016020908152604080832093881683529290522083905581156116de57836001600160a01b0316856001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925856040516116d591815260200190565b60405180910390a35b5050505050565b6000805160206121d78339815191526001600160a01b038416611721578181600201600082825461171691906121a3565b909155506117939050565b6001600160a01b038416600090815260208290526040902054828110156117745760405163391434e360e21b81526001600160a01b038616600482015260248101829052604481018490526064016107f0565b6001600160a01b03851660009081526020839052604090209083900390555b6001600160a01b0383166117b15760028101805483900390556117d0565b6001600160a01b03831660009081526020829052604090208054830190555b826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161181591815260200190565b60405180910390a350505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054600160401b900460ff1661119457604051631afcd79f60e31b815260040160405180910390fd5b611194611823565b61187c611823565b6111946119b8565b61188c611823565b6000805160206121d78339815191527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace046118c685826120e4565b50600581016118d584826120e4565b50600301805460ff909216600160e01b0260ff60e01b199092169190911790555050565b611901611823565b604080516064808252610ca082019092526000805160206121f78339815191529160208201610c80803683375050815161079692600185019250602001906119cb565b61194e828261092f565b6107965760405163e2517d3f60e01b81526001600160a01b0382166004820152602481018390526044016107f0565b6000805160206121f783398151915263ffffffff8216156107965763ffffffff82166000908152600282016020526040812061079691611a7a565b6119c0611823565b6111a06000336111a3565b82805482825590600052602060002090600701600890048101928215611a6a5791602002820160005b83821115611a3857835183826101000a81548163ffffffff021916908363ffffffff16021790555092602001926004016020816003010492830192600103026119f4565b8015611a685782816101000a81549063ffffffff0219169055600401602081600301049283019260010302611a38565b505b50611a76929150611a9b565b5090565b50805460008255600202906000526020600020908101906111a09190611ab0565b5b80821115611a765760008155600101611a9c565b80821115611a765780546001600160801b03191681556000611ad56001830182611ade565b50600201611ab0565b508054611aea90611f9c565b6000825580601f10611afa575050565b601f0160209004906000526020600020908101906111a09190611a9b565b600060208284031215611b2a57600080fd5b81356001600160e01b031981168114611b4257600080fd5b9392505050565b6000815180845260005b81811015611b6f57602081850181015186830182015201611b53565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000611b426020830184611b49565b80356001600160a01b0381168114611bb957600080fd5b919050565b60008060408385031215611bd157600080fd5b611bda83611ba2565b946020939093013593505050565b600080600060608486031215611bfd57600080fd5b611c0684611ba2565b9250611c1460208501611ba2565b9150604084013590509250925092565b600060208284031215611c3657600080fd5b5035919050565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b0380841115611c6d57611c6d611c3d565b604051601f8501601f19908116603f01168101908282118183101715611c9557611c95611c3d565b81604052809350858152868686011115611cae57600080fd5b858560208301376000602087830101525050509392505050565b600082601f830112611cd957600080fd5b611b4283833560208501611c53565b600080600060608486031215611cfd57600080fd5b83356001600160401b0380821115611d1457600080fd5b611d2087838801611cc8565b94506020860135915080821115611d3657600080fd5b50611d4386828701611cc8565b925050604084013560ff81168114611d5a57600080fd5b809150509250925092565b60008060408385031215611d7857600080fd5b82359150611d8860208401611ba2565b90509250929050565b80356001600160401b0381168114611bb957600080fd5b60008060408385031215611dbb57600080fd5b611dc483611d91565b915060208301356001600160401b03811115611ddf57600080fd5b8301601f81018513611df057600080fd5b611dff85823560208401611c53565b9150509250929050565b600060208284031215611e1b57600080fd5b611b4282611ba2565b600060208284031215611e3657600080fd5b611b4282611d91565b60008060408385031215611e5257600080fd5b611e5b83611ba2565b9150611d8860208401611ba2565b600060208284031215611e7b57600080fd5b813561ffff81168114611b4257600080fd5b60008151808452602080850194506020840160005b83811015611ec75781516001600160401b031687529582019590820190600101611ea2565b509495945050505050565b606081526000611ee56060830186611e8d565b602083820381850152611ef88287611e8d565b915083820360408501528185518084528284019150828160051b85010183880160005b83811015611f4957601f19878403018552611f37838351611b49565b94860194925090850190600101611f1b565b50909a9950505050505050505050565b600080600060608486031215611f6e57600080fd5b611f7784611ba2565b9250611f8560208501611d91565b9150611f9360408501611d91565b90509250925092565b600181811c90821680611fb057607f821691505b602082108103611fd057634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b6001600160401b0382811682821603908082111561124857611248611fd6565b60006001600160401b0380861683526060602084015261202f6060840186611b49565b9150808416604084015250949350505050565b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052603260045260246000fd5b60006001600160401b0380831681810361208a5761208a611fd6565b6001019392505050565b601f82111561074d576000816000526020600020601f850160051c810160208610156120bd5750805b601f850160051c820191505b818110156120dc578281556001016120c9565b505050505050565b81516001600160401b038111156120fd576120fd611c3d565b6121118161210b8454611f9c565b84612094565b602080601f831160018114612146576000841561212e5750858301515b600019600386901b1c1916600185901b1785556120dc565b600085815260208120601f198616915b8281101561217557888601518255948401946001909101908401612156565b50858210156121935787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b808201808211156104b7576104b7611fd656fe52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0352c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00c87a15b711c3ce3a8db8c4c853ae9c01bdbe93da3d2e98c03c16ae2fd0346f0002dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800b458d13b0f4ce9e6aa65d297a27b10f75fdc6d0957bb29e1f2a30c8766b35415a26469706673582212202dc2c0d79d5c246c0250481504369ff02aebbd8c0d65993ebc27493c4c63b06464736f6c63430008180033",
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

// GetRecentRedemptionData is a free data retrieval call binding the contract method 0xf1690dbe.
//
// Solidity: function getRecentRedemptionData(uint16 numBlocks) view returns(uint64[] ids, uint64[] amounts, bytes[] destination)
func (_ZenBTC *ZenBTCCaller) GetRecentRedemptionData(opts *bind.CallOpts, numBlocks uint16) (struct {
	Ids         []uint64
	Amounts     []uint64
	Destination [][]byte
}, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "getRecentRedemptionData", numBlocks)

	outstruct := new(struct {
		Ids         []uint64
		Amounts     []uint64
		Destination [][]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Ids = *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)
	outstruct.Amounts = *abi.ConvertType(out[1], new([]uint64)).(*[]uint64)
	outstruct.Destination = *abi.ConvertType(out[2], new([][]byte)).(*[][]byte)

	return *outstruct, err

}

// GetRecentRedemptionData is a free data retrieval call binding the contract method 0xf1690dbe.
//
// Solidity: function getRecentRedemptionData(uint16 numBlocks) view returns(uint64[] ids, uint64[] amounts, bytes[] destination)
func (_ZenBTC *ZenBTCSession) GetRecentRedemptionData(numBlocks uint16) (struct {
	Ids         []uint64
	Amounts     []uint64
	Destination [][]byte
}, error) {
	return _ZenBTC.Contract.GetRecentRedemptionData(&_ZenBTC.CallOpts, numBlocks)
}

// GetRecentRedemptionData is a free data retrieval call binding the contract method 0xf1690dbe.
//
// Solidity: function getRecentRedemptionData(uint16 numBlocks) view returns(uint64[] ids, uint64[] amounts, bytes[] destination)
func (_ZenBTC *ZenBTCCallerSession) GetRecentRedemptionData(numBlocks uint16) (struct {
	Ids         []uint64
	Amounts     []uint64
	Destination [][]byte
}, error) {
	return _ZenBTC.Contract.GetRecentRedemptionData(&_ZenBTC.CallOpts, numBlocks)
}

// GetRedeemFee is a free data retrieval call binding the contract method 0xc6d98f1a.
//
// Solidity: function getRedeemFee() view returns(uint256)
func (_ZenBTC *ZenBTCCaller) GetRedeemFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZenBTC.contract.Call(opts, &out, "getRedeemFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRedeemFee is a free data retrieval call binding the contract method 0xc6d98f1a.
//
// Solidity: function getRedeemFee() view returns(uint256)
func (_ZenBTC *ZenBTCSession) GetRedeemFee() (*big.Int, error) {
	return _ZenBTC.Contract.GetRedeemFee(&_ZenBTC.CallOpts)
}

// GetRedeemFee is a free data retrieval call binding the contract method 0xc6d98f1a.
//
// Solidity: function getRedeemFee() view returns(uint256)
func (_ZenBTC *ZenBTCCallerSession) GetRedeemFee() (*big.Int, error) {
	return _ZenBTC.Contract.GetRedeemFee(&_ZenBTC.CallOpts)
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

// Unwrap is a paid mutator transaction binding the contract method 0x4677d279.
//
// Solidity: function unwrap(uint64 value, bytes destAddr) returns()
func (_ZenBTC *ZenBTCTransactor) Unwrap(opts *bind.TransactOpts, value uint64, destAddr []byte) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "unwrap", value, destAddr)
}

// Unwrap is a paid mutator transaction binding the contract method 0x4677d279.
//
// Solidity: function unwrap(uint64 value, bytes destAddr) returns()
func (_ZenBTC *ZenBTCSession) Unwrap(value uint64, destAddr []byte) (*types.Transaction, error) {
	return _ZenBTC.Contract.Unwrap(&_ZenBTC.TransactOpts, value, destAddr)
}

// Unwrap is a paid mutator transaction binding the contract method 0x4677d279.
//
// Solidity: function unwrap(uint64 value, bytes destAddr) returns()
func (_ZenBTC *ZenBTCTransactorSession) Unwrap(value uint64, destAddr []byte) (*types.Transaction, error) {
	return _ZenBTC.Contract.Unwrap(&_ZenBTC.TransactOpts, value, destAddr)
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

// UpdateRedeemFee is a paid mutator transaction binding the contract method 0xce29ea62.
//
// Solidity: function updateRedeemFee(uint64 redeemFee_) returns()
func (_ZenBTC *ZenBTCTransactor) UpdateRedeemFee(opts *bind.TransactOpts, redeemFee_ uint64) (*types.Transaction, error) {
	return _ZenBTC.contract.Transact(opts, "updateRedeemFee", redeemFee_)
}

// UpdateRedeemFee is a paid mutator transaction binding the contract method 0xce29ea62.
//
// Solidity: function updateRedeemFee(uint64 redeemFee_) returns()
func (_ZenBTC *ZenBTCSession) UpdateRedeemFee(redeemFee_ uint64) (*types.Transaction, error) {
	return _ZenBTC.Contract.UpdateRedeemFee(&_ZenBTC.TransactOpts, redeemFee_)
}

// UpdateRedeemFee is a paid mutator transaction binding the contract method 0xce29ea62.
//
// Solidity: function updateRedeemFee(uint64 redeemFee_) returns()
func (_ZenBTC *ZenBTCTransactorSession) UpdateRedeemFee(redeemFee_ uint64) (*types.Transaction, error) {
	return _ZenBTC.Contract.UpdateRedeemFee(&_ZenBTC.TransactOpts, redeemFee_)
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
	Fee      uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokenRedemption is a free log retrieval operation binding the contract event 0x4f8266db69b83687df444f19eab599803494f35fe2a1c061ba66d871bda06824.
//
// Solidity: event TokenRedemption(address indexed redeemer, uint64 value, bytes destAddr, uint64 fee)
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

// WatchTokenRedemption is a free log subscription operation binding the contract event 0x4f8266db69b83687df444f19eab599803494f35fe2a1c061ba66d871bda06824.
//
// Solidity: event TokenRedemption(address indexed redeemer, uint64 value, bytes destAddr, uint64 fee)
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

// ParseTokenRedemption is a log parse operation binding the contract event 0x4f8266db69b83687df444f19eab599803494f35fe2a1c061ba66d871bda06824.
//
// Solidity: event TokenRedemption(address indexed redeemer, uint64 value, bytes destAddr, uint64 fee)
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