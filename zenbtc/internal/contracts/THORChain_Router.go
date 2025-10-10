// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// THORChainRouterCoin is an auto generated low-level Go binding around an user-defined struct.
type THORChainRouterCoin struct {
	Asset  common.Address
	Amount *big.Int
}

// ThorchainrouterMetaData contains all meta data concerning the Thorchainrouter contract.
var ThorchainrouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rune\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldVault\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newVault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"TransferAllowance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"TransferOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"finalAsset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"TransferOutAndCall\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldVault\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newVault\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structTHORChain_Router.Coin[]\",\"name\":\"coins\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"VaultTransfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"RUNE\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"name\":\"depositWithExpiry\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"asgard\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structTHORChain_Router.Coin[]\",\"name\":\"coins\",\"type\":\"tuple[]\"},{\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"returnVaultAssets\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newVault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"transferAllowance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"transferOut\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"finalToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"transferOutAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"vaultAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620023fe380380620023fe8339818101604052810190620000379190620000f0565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060016002819055505062000122565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000b8826200008b565b9050919050565b620000ca81620000ab565b8114620000d657600080fd5b50565b600081519050620000ea81620000bf565b92915050565b60006020828403121562000109576200010862000086565b5b60006200011984828501620000d9565b91505092915050565b6122cc80620001326000396000f3fe6080604052600436106100705760003560e01c80634039fd4b1161004e5780634039fd4b146100f757806344bc937b14610113578063574da7171461012f57806393e4eaa91461014b57610070565b806303b6a673146100755780631b738b32146100b25780632923e82e146100db575b600080fd5b34801561008157600080fd5b5061009c600480360381019061009791906114e9565b610176565b6040516100a99190611542565b60405180910390f35b3480156100be57600080fd5b506100d960048036038101906100d491906116cf565b6101fd565b005b6100f560048036038101906100f091906118c1565b610311565b005b610111600480360381019061010c9190611960565b61051e565b005b61012d600480360381019061012891906119f7565b610763565b005b61014960048036038101906101449190611a8e565b6107b8565b005b34801561015757600080fd5b50610160610b08565b60405161016d9190611b20565b60405180910390f35b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6002805403610241576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161023890611b98565b60405180910390fd5b600280819055503073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16036102f457610286848484610b2c565b8373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f05b90458f953d3fcb2d7fb25616a2fddeca749d0c47cc5c9832d0266b5346eea8585856040516102e793929190611c2f565b60405180910390a3610302565b6103018585858585610c57565b5b60016002819055505050505050565b6002805403610355576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161034c90611b98565b60405180910390fd5b600280819055503073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036104625760005b82518110156103f5576103e2848483815181106103b2576103b1611c6d565b5b6020026020010151600001518584815181106103d1576103d0611c6d565b5b602002602001015160200151610b2c565b80806103ed90611ccb565b915050610392565b508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f281daef48d91e5cd3d32db0784f6af69cd8d8d2e8c612a3568dca51ded51e08f8484604051610455929190611e0f565b60405180910390a36104cc565b60005b82518110156104ca576104b7858585848151811061048657610485611c6d565b5b6020026020010151600001518685815181106104a5576104a4611c6d565b5b60200260200101516020015186610c57565b80806104c290611ccb565b915050610465565b505b60008373ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f1935050505090508061050f57600080fd5b50600160028190555050505050565b6002805403610562576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161055990611b98565b60405180910390fd5b60028081905550600034905060008673ffffffffffffffffffffffffffffffffffffffff168287878760405160240161059d93929190611e46565b6040516020818303038152906040527f48c314f4000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516106279190611ec4565b60006040518083038185875af1925050503d8060008114610664576040519150601f19603f3d011682016040523d82523d6000602084013e610669565b606091505b50509050806106fa5760008573ffffffffffffffffffffffffffffffffffffffff166108fc849081150290604051600060405180830381858888f193505050509050806106f8573373ffffffffffffffffffffffffffffffffffffffff166108fc849081150290604051600060405180830381858888f193505050501580156106f6573d6000803e3d6000fd5b505b505b3373ffffffffffffffffffffffffffffffffffffffff167f8e5841bcd195b858d53b38bcf91b38d47f3bc800469b6812d35451ab619c6f6c88848989898960405161074a96959493929190611f3a565b60405180910390a2505060016002819055505050505050565b8042106107a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161079c90611fee565b60405180910390fd5b6107b185858585610e8b565b5050505050565b60028054036107fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107f390611b98565b60405180910390fd5b6002808190555060008073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036108c75734905060008573ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f193505050509050806108c1573373ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f193505050501580156108bf573d6000803e3d6000fd5b505b50610a90565b82600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610953919061200e565b925050819055506000808573ffffffffffffffffffffffffffffffffffffffff168786604051602401610987929190612051565b6040516020818303038152906040527fa9059cbb000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610a119190611ec4565b6000604051808303816000865af19150503d8060008114610a4e576040519150601f19603f3d011682016040523d82523d6000602084013e610a53565b606091505b5091509150818015610a815750600081511480610a80575080806020019051810190610a7f91906120b2565b5b5b610a8a57600080fd5b84925050505b8473ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fa9cd03aa3c1b4515114539cd53d22085129d495cb9e9f9af77864526240f1bf7868486604051610af193929190611c2f565b60405180910390a350600160028190555050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610bb8919061200e565b9250508190555080600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610c4b91906120df565b92505081905550505050565b81600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610ce3919061200e565b9250508190555060008373ffffffffffffffffffffffffffffffffffffffff168684604051602401610d16929190612135565b6040516020818303038152906040527f095ea7b3000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610da09190611ec4565b6000604051808303816000865af19150503d8060008114610ddd576040519150601f19603f3d011682016040523d82523d6000602084013e610de2565b606091505b5050905080610df057600080fd5b8573ffffffffffffffffffffffffffffffffffffffff166344bc937b868686867fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6040518663ffffffff1660e01b8152600401610e5195949392919061215e565b600060405180830381600087803b158015610e6b57600080fd5b505af1158015610e7f573d6000803e3d6000fd5b50505050505050505050565b6002805403610ecf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ec690611b98565b60405180910390fd5b6002808190555060008073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610f575734905060008573ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f19350505050905080610f5157600080fd5b506111c0565b60003414610f9a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f9190612204565b60405180910390fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff160361111f5782905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632ccb1b3030856040518363ffffffff1660e01b815260040161104b929190612135565b6020604051808303816000875af115801561106a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061108e91906120b2565b5060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166342966c68846040518263ffffffff1660e01b81526004016110e89190611542565b600060405180830381600087803b15801561110257600080fd5b505af1158015611116573d6000803e3d6000fd5b505050506111bf565b6111298484611236565b905080600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546111b791906120df565b925050819055505b5b8373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fef519b7eb82aaf6ac376a6df2d793843ebfd593de5f1a0601d3cc6ab49ebb395838560405161121f929190612224565b60405180910390a350600160028190555050505050565b6000808373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016112729190611b20565b602060405180830381865afa15801561128f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112b39190612269565b90506000808573ffffffffffffffffffffffffffffffffffffffff163330876040516024016112e493929190611e46565b6040516020818303038152906040527f23b872dd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161136e9190611ec4565b6000604051808303816000865af19150503d80600081146113ab576040519150601f19603f3d011682016040523d82523d6000602084013e6113b0565b606091505b50915091508180156113de57506000815114806113dd5750808060200190518101906113dc91906120b2565b5b5b6113e757600080fd5b828673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016114219190611b20565b602060405180830381865afa15801561143e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114629190612269565b61146c919061200e565b935050505092915050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006114b68261148b565b9050919050565b6114c6816114ab565b81146114d157600080fd5b50565b6000813590506114e3816114bd565b92915050565b60008060408385031215611500576114ff611481565b5b600061150e858286016114d4565b925050602061151f858286016114d4565b9150509250929050565b6000819050919050565b61153c81611529565b82525050565b60006020820190506115576000830184611533565b92915050565b61156681611529565b811461157157600080fd5b50565b6000813590506115838161155d565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6115dc82611593565b810181811067ffffffffffffffff821117156115fb576115fa6115a4565b5b80604052505050565b600061160e611477565b905061161a82826115d3565b919050565b600067ffffffffffffffff82111561163a576116396115a4565b5b61164382611593565b9050602081019050919050565b82818337600083830152505050565b600061167261166d8461161f565b611604565b90508281526020810184848401111561168e5761168d61158e565b5b611699848285611650565b509392505050565b600082601f8301126116b6576116b5611589565b5b81356116c684826020860161165f565b91505092915050565b600080600080600060a086880312156116eb576116ea611481565b5b60006116f9888289016114d4565b955050602061170a888289016114d4565b945050604061171b888289016114d4565b935050606061172c88828901611574565b925050608086013567ffffffffffffffff81111561174d5761174c611486565b5b611759888289016116a1565b9150509295509295909350565b60006117718261148b565b9050919050565b61178181611766565b811461178c57600080fd5b50565b60008135905061179e81611778565b92915050565b600067ffffffffffffffff8211156117bf576117be6115a4565b5b602082029050602081019050919050565b600080fd5b600080fd5b6000604082840312156117f0576117ef6117d5565b5b6117fa6040611604565b9050600061180a848285016114d4565b600083015250602061181e84828501611574565b60208301525092915050565b600061183d611838846117a4565b611604565b905080838252602082019050604084028301858111156118605761185f6117d0565b5b835b81811015611889578061187588826117da565b845260208401935050604081019050611862565b5050509392505050565b600082601f8301126118a8576118a7611589565b5b81356118b884826020860161182a565b91505092915050565b600080600080608085870312156118db576118da611481565b5b60006118e9878288016114d4565b94505060206118fa8782880161178f565b935050604085013567ffffffffffffffff81111561191b5761191a611486565b5b61192787828801611893565b925050606085013567ffffffffffffffff81111561194857611947611486565b5b611954878288016116a1565b91505092959194509250565b600080600080600060a0868803121561197c5761197b611481565b5b600061198a8882890161178f565b955050602061199b888289016114d4565b94505060406119ac888289016114d4565b93505060606119bd88828901611574565b925050608086013567ffffffffffffffff8111156119de576119dd611486565b5b6119ea888289016116a1565b9150509295509295909350565b600080600080600060a08688031215611a1357611a12611481565b5b6000611a218882890161178f565b9550506020611a32888289016114d4565b9450506040611a4388828901611574565b935050606086013567ffffffffffffffff811115611a6457611a63611486565b5b611a70888289016116a1565b9250506080611a8188828901611574565b9150509295509295909350565b60008060008060808587031215611aa857611aa7611481565b5b6000611ab68782880161178f565b9450506020611ac7878288016114d4565b9350506040611ad887828801611574565b925050606085013567ffffffffffffffff811115611af957611af8611486565b5b611b05878288016116a1565b91505092959194509250565b611b1a816114ab565b82525050565b6000602082019050611b356000830184611b11565b92915050565b600082825260208201905092915050565b7f5265656e7472616e637947756172643a207265656e7472616e742063616c6c00600082015250565b6000611b82601f83611b3b565b9150611b8d82611b4c565b602082019050919050565b60006020820190508181036000830152611bb181611b75565b9050919050565b600081519050919050565b60005b83811015611be1578082015181840152602081019050611bc6565b83811115611bf0576000848401525b50505050565b6000611c0182611bb8565b611c0b8185611b3b565b9350611c1b818560208601611bc3565b611c2481611593565b840191505092915050565b6000606082019050611c446000830186611b11565b611c516020830185611533565b8181036040830152611c638184611bf6565b9050949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611cd682611529565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611d0857611d07611c9c565b5b600182019050919050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b611d48816114ab565b82525050565b611d5781611529565b82525050565b604082016000820151611d736000850182611d3f565b506020820151611d866020850182611d4e565b50505050565b6000611d988383611d5d565b60408301905092915050565b6000602082019050919050565b6000611dbc82611d13565b611dc68185611d1e565b9350611dd183611d2f565b8060005b83811015611e02578151611de98882611d8c565b9750611df483611da4565b925050600181019050611dd5565b5085935050505092915050565b60006040820190508181036000830152611e298185611db1565b90508181036020830152611e3d8184611bf6565b90509392505050565b6000606082019050611e5b6000830186611b11565b611e686020830185611b11565b611e756040830184611533565b949350505050565b600081519050919050565b600081905092915050565b6000611e9e82611e7d565b611ea88185611e88565b9350611eb8818560208601611bc3565b80840191505092915050565b6000611ed08284611e93565b915081905092915050565b6000819050919050565b6000611f00611efb611ef68461148b565b611edb565b61148b565b9050919050565b6000611f1282611ee5565b9050919050565b6000611f2482611f07565b9050919050565b611f3481611f19565b82525050565b600060c082019050611f4f6000830189611f2b565b611f5c6020830188611533565b611f696040830187611b11565b611f766060830186611b11565b611f836080830185611533565b81810360a0830152611f958184611bf6565b9050979650505050505050565b7f54484f52436861696e5f526f757465723a206578706972656400000000000000600082015250565b6000611fd8601983611b3b565b9150611fe382611fa2565b602082019050919050565b6000602082019050818103600083015261200781611fcb565b9050919050565b600061201982611529565b915061202483611529565b92508282101561203757612036611c9c565b5b828203905092915050565b61204b81611766565b82525050565b60006040820190506120666000830185612042565b6120736020830184611533565b9392505050565b60008115159050919050565b61208f8161207a565b811461209a57600080fd5b50565b6000815190506120ac81612086565b92915050565b6000602082840312156120c8576120c7611481565b5b60006120d68482850161209d565b91505092915050565b60006120ea82611529565b91506120f583611529565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561212a57612129611c9c565b5b828201905092915050565b600060408201905061214a6000830185611b11565b6121576020830184611533565b9392505050565b600060a0820190506121736000830188611b11565b6121806020830187611b11565b61218d6040830186611533565b818103606083015261219f8185611bf6565b90506121ae6080830184611533565b9695505050505050565b7f756e657870656374656420657468000000000000000000000000000000000000600082015250565b60006121ee600e83611b3b565b91506121f9826121b8565b602082019050919050565b6000602082019050818103600083015261221d816121e1565b9050919050565b60006040820190506122396000830185611533565b818103602083015261224b8184611bf6565b90509392505050565b6000815190506122638161155d565b92915050565b60006020828403121561227f5761227e611481565b5b600061228d84828501612254565b9150509291505056fea2646970667358221220392d3f7fdcff07be4d290d8e27af80806b03b248abc5b3483eba0e5e2ae2cbaa64736f6c634300080d0033",
}

// ThorchainrouterABI is the input ABI used to generate the binding from.
// Deprecated: Use ThorchainrouterMetaData.ABI instead.
var ThorchainrouterABI = ThorchainrouterMetaData.ABI

// ThorchainrouterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ThorchainrouterMetaData.Bin instead.
var ThorchainrouterBin = ThorchainrouterMetaData.Bin

// DeployThorchainrouter deploys a new Ethereum contract, binding an instance of Thorchainrouter to it.
func DeployThorchainrouter(auth *bind.TransactOpts, backend bind.ContractBackend, rune common.Address) (common.Address, *types.Transaction, *Thorchainrouter, error) {
	parsed, err := ThorchainrouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ThorchainrouterBin), backend, rune)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Thorchainrouter{ThorchainrouterCaller: ThorchainrouterCaller{contract: contract}, ThorchainrouterTransactor: ThorchainrouterTransactor{contract: contract}, ThorchainrouterFilterer: ThorchainrouterFilterer{contract: contract}}, nil
}

// Thorchainrouter is an auto generated Go binding around an Ethereum contract.
type Thorchainrouter struct {
	ThorchainrouterCaller     // Read-only binding to the contract
	ThorchainrouterTransactor // Write-only binding to the contract
	ThorchainrouterFilterer   // Log filterer for contract events
}

// ThorchainrouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type ThorchainrouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ThorchainrouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ThorchainrouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ThorchainrouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ThorchainrouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ThorchainrouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ThorchainrouterSession struct {
	Contract     *Thorchainrouter  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ThorchainrouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ThorchainrouterCallerSession struct {
	Contract *ThorchainrouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ThorchainrouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ThorchainrouterTransactorSession struct {
	Contract     *ThorchainrouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ThorchainrouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type ThorchainrouterRaw struct {
	Contract *Thorchainrouter // Generic contract binding to access the raw methods on
}

// ThorchainrouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ThorchainrouterCallerRaw struct {
	Contract *ThorchainrouterCaller // Generic read-only contract binding to access the raw methods on
}

// ThorchainrouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ThorchainrouterTransactorRaw struct {
	Contract *ThorchainrouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewThorchainrouter creates a new instance of Thorchainrouter, bound to a specific deployed contract.
func NewThorchainrouter(address common.Address, backend bind.ContractBackend) (*Thorchainrouter, error) {
	contract, err := bindThorchainrouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Thorchainrouter{ThorchainrouterCaller: ThorchainrouterCaller{contract: contract}, ThorchainrouterTransactor: ThorchainrouterTransactor{contract: contract}, ThorchainrouterFilterer: ThorchainrouterFilterer{contract: contract}}, nil
}

// NewThorchainrouterCaller creates a new read-only instance of Thorchainrouter, bound to a specific deployed contract.
func NewThorchainrouterCaller(address common.Address, caller bind.ContractCaller) (*ThorchainrouterCaller, error) {
	contract, err := bindThorchainrouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ThorchainrouterCaller{contract: contract}, nil
}

// NewThorchainrouterTransactor creates a new write-only instance of Thorchainrouter, bound to a specific deployed contract.
func NewThorchainrouterTransactor(address common.Address, transactor bind.ContractTransactor) (*ThorchainrouterTransactor, error) {
	contract, err := bindThorchainrouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ThorchainrouterTransactor{contract: contract}, nil
}

// NewThorchainrouterFilterer creates a new log filterer instance of Thorchainrouter, bound to a specific deployed contract.
func NewThorchainrouterFilterer(address common.Address, filterer bind.ContractFilterer) (*ThorchainrouterFilterer, error) {
	contract, err := bindThorchainrouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ThorchainrouterFilterer{contract: contract}, nil
}

// bindThorchainrouter binds a generic wrapper to an already deployed contract.
func bindThorchainrouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ThorchainrouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Thorchainrouter *ThorchainrouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Thorchainrouter.Contract.ThorchainrouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Thorchainrouter *ThorchainrouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.ThorchainrouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Thorchainrouter *ThorchainrouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.ThorchainrouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Thorchainrouter *ThorchainrouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Thorchainrouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Thorchainrouter *ThorchainrouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Thorchainrouter *ThorchainrouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.contract.Transact(opts, method, params...)
}

// RUNE is a free data retrieval call binding the contract method 0x93e4eaa9.
//
// Solidity: function RUNE() view returns(address)
func (_Thorchainrouter *ThorchainrouterCaller) RUNE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Thorchainrouter.contract.Call(opts, &out, "RUNE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RUNE is a free data retrieval call binding the contract method 0x93e4eaa9.
//
// Solidity: function RUNE() view returns(address)
func (_Thorchainrouter *ThorchainrouterSession) RUNE() (common.Address, error) {
	return _Thorchainrouter.Contract.RUNE(&_Thorchainrouter.CallOpts)
}

// RUNE is a free data retrieval call binding the contract method 0x93e4eaa9.
//
// Solidity: function RUNE() view returns(address)
func (_Thorchainrouter *ThorchainrouterCallerSession) RUNE() (common.Address, error) {
	return _Thorchainrouter.Contract.RUNE(&_Thorchainrouter.CallOpts)
}

// VaultAllowance is a free data retrieval call binding the contract method 0x03b6a673.
//
// Solidity: function vaultAllowance(address vault, address token) view returns(uint256 amount)
func (_Thorchainrouter *ThorchainrouterCaller) VaultAllowance(opts *bind.CallOpts, vault common.Address, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Thorchainrouter.contract.Call(opts, &out, "vaultAllowance", vault, token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VaultAllowance is a free data retrieval call binding the contract method 0x03b6a673.
//
// Solidity: function vaultAllowance(address vault, address token) view returns(uint256 amount)
func (_Thorchainrouter *ThorchainrouterSession) VaultAllowance(vault common.Address, token common.Address) (*big.Int, error) {
	return _Thorchainrouter.Contract.VaultAllowance(&_Thorchainrouter.CallOpts, vault, token)
}

// VaultAllowance is a free data retrieval call binding the contract method 0x03b6a673.
//
// Solidity: function vaultAllowance(address vault, address token) view returns(uint256 amount)
func (_Thorchainrouter *ThorchainrouterCallerSession) VaultAllowance(vault common.Address, token common.Address) (*big.Int, error) {
	return _Thorchainrouter.Contract.VaultAllowance(&_Thorchainrouter.CallOpts, vault, token)
}

// DepositWithExpiry is a paid mutator transaction binding the contract method 0x44bc937b.
//
// Solidity: function depositWithExpiry(address vault, address asset, uint256 amount, string memo, uint256 expiration) payable returns()
func (_Thorchainrouter *ThorchainrouterTransactor) DepositWithExpiry(opts *bind.TransactOpts, vault common.Address, asset common.Address, amount *big.Int, memo string, expiration *big.Int) (*types.Transaction, error) {
	return _Thorchainrouter.contract.Transact(opts, "depositWithExpiry", vault, asset, amount, memo, expiration)
}

// DepositWithExpiry is a paid mutator transaction binding the contract method 0x44bc937b.
//
// Solidity: function depositWithExpiry(address vault, address asset, uint256 amount, string memo, uint256 expiration) payable returns()
func (_Thorchainrouter *ThorchainrouterSession) DepositWithExpiry(vault common.Address, asset common.Address, amount *big.Int, memo string, expiration *big.Int) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.DepositWithExpiry(&_Thorchainrouter.TransactOpts, vault, asset, amount, memo, expiration)
}

// DepositWithExpiry is a paid mutator transaction binding the contract method 0x44bc937b.
//
// Solidity: function depositWithExpiry(address vault, address asset, uint256 amount, string memo, uint256 expiration) payable returns()
func (_Thorchainrouter *ThorchainrouterTransactorSession) DepositWithExpiry(vault common.Address, asset common.Address, amount *big.Int, memo string, expiration *big.Int) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.DepositWithExpiry(&_Thorchainrouter.TransactOpts, vault, asset, amount, memo, expiration)
}

// ReturnVaultAssets is a paid mutator transaction binding the contract method 0x2923e82e.
//
// Solidity: function returnVaultAssets(address router, address asgard, (address,uint256)[] coins, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterTransactor) ReturnVaultAssets(opts *bind.TransactOpts, router common.Address, asgard common.Address, coins []THORChainRouterCoin, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.contract.Transact(opts, "returnVaultAssets", router, asgard, coins, memo)
}

// ReturnVaultAssets is a paid mutator transaction binding the contract method 0x2923e82e.
//
// Solidity: function returnVaultAssets(address router, address asgard, (address,uint256)[] coins, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterSession) ReturnVaultAssets(router common.Address, asgard common.Address, coins []THORChainRouterCoin, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.ReturnVaultAssets(&_Thorchainrouter.TransactOpts, router, asgard, coins, memo)
}

// ReturnVaultAssets is a paid mutator transaction binding the contract method 0x2923e82e.
//
// Solidity: function returnVaultAssets(address router, address asgard, (address,uint256)[] coins, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterTransactorSession) ReturnVaultAssets(router common.Address, asgard common.Address, coins []THORChainRouterCoin, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.ReturnVaultAssets(&_Thorchainrouter.TransactOpts, router, asgard, coins, memo)
}

// TransferAllowance is a paid mutator transaction binding the contract method 0x1b738b32.
//
// Solidity: function transferAllowance(address router, address newVault, address asset, uint256 amount, string memo) returns()
func (_Thorchainrouter *ThorchainrouterTransactor) TransferAllowance(opts *bind.TransactOpts, router common.Address, newVault common.Address, asset common.Address, amount *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.contract.Transact(opts, "transferAllowance", router, newVault, asset, amount, memo)
}

// TransferAllowance is a paid mutator transaction binding the contract method 0x1b738b32.
//
// Solidity: function transferAllowance(address router, address newVault, address asset, uint256 amount, string memo) returns()
func (_Thorchainrouter *ThorchainrouterSession) TransferAllowance(router common.Address, newVault common.Address, asset common.Address, amount *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.TransferAllowance(&_Thorchainrouter.TransactOpts, router, newVault, asset, amount, memo)
}

// TransferAllowance is a paid mutator transaction binding the contract method 0x1b738b32.
//
// Solidity: function transferAllowance(address router, address newVault, address asset, uint256 amount, string memo) returns()
func (_Thorchainrouter *ThorchainrouterTransactorSession) TransferAllowance(router common.Address, newVault common.Address, asset common.Address, amount *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.TransferAllowance(&_Thorchainrouter.TransactOpts, router, newVault, asset, amount, memo)
}

// TransferOut is a paid mutator transaction binding the contract method 0x574da717.
//
// Solidity: function transferOut(address to, address asset, uint256 amount, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterTransactor) TransferOut(opts *bind.TransactOpts, to common.Address, asset common.Address, amount *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.contract.Transact(opts, "transferOut", to, asset, amount, memo)
}

// TransferOut is a paid mutator transaction binding the contract method 0x574da717.
//
// Solidity: function transferOut(address to, address asset, uint256 amount, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterSession) TransferOut(to common.Address, asset common.Address, amount *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.TransferOut(&_Thorchainrouter.TransactOpts, to, asset, amount, memo)
}

// TransferOut is a paid mutator transaction binding the contract method 0x574da717.
//
// Solidity: function transferOut(address to, address asset, uint256 amount, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterTransactorSession) TransferOut(to common.Address, asset common.Address, amount *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.TransferOut(&_Thorchainrouter.TransactOpts, to, asset, amount, memo)
}

// TransferOutAndCall is a paid mutator transaction binding the contract method 0x4039fd4b.
//
// Solidity: function transferOutAndCall(address target, address finalToken, address to, uint256 amountOutMin, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterTransactor) TransferOutAndCall(opts *bind.TransactOpts, target common.Address, finalToken common.Address, to common.Address, amountOutMin *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.contract.Transact(opts, "transferOutAndCall", target, finalToken, to, amountOutMin, memo)
}

// TransferOutAndCall is a paid mutator transaction binding the contract method 0x4039fd4b.
//
// Solidity: function transferOutAndCall(address target, address finalToken, address to, uint256 amountOutMin, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterSession) TransferOutAndCall(target common.Address, finalToken common.Address, to common.Address, amountOutMin *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.TransferOutAndCall(&_Thorchainrouter.TransactOpts, target, finalToken, to, amountOutMin, memo)
}

// TransferOutAndCall is a paid mutator transaction binding the contract method 0x4039fd4b.
//
// Solidity: function transferOutAndCall(address target, address finalToken, address to, uint256 amountOutMin, string memo) payable returns()
func (_Thorchainrouter *ThorchainrouterTransactorSession) TransferOutAndCall(target common.Address, finalToken common.Address, to common.Address, amountOutMin *big.Int, memo string) (*types.Transaction, error) {
	return _Thorchainrouter.Contract.TransferOutAndCall(&_Thorchainrouter.TransactOpts, target, finalToken, to, amountOutMin, memo)
}

// ThorchainrouterDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Thorchainrouter contract.
type ThorchainrouterDepositIterator struct {
	Event *ThorchainrouterDeposit // Event containing the contract specifics and raw log

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
func (it *ThorchainrouterDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThorchainrouterDeposit)
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
		it.Event = new(ThorchainrouterDeposit)
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
func (it *ThorchainrouterDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThorchainrouterDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThorchainrouterDeposit represents a Deposit event raised by the Thorchainrouter contract.
type ThorchainrouterDeposit struct {
	To     common.Address
	Asset  common.Address
	Amount *big.Int
	Memo   string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xef519b7eb82aaf6ac376a6df2d793843ebfd593de5f1a0601d3cc6ab49ebb395.
//
// Solidity: event Deposit(address indexed to, address indexed asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) FilterDeposit(opts *bind.FilterOpts, to []common.Address, asset []common.Address) (*ThorchainrouterDepositIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _Thorchainrouter.contract.FilterLogs(opts, "Deposit", toRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &ThorchainrouterDepositIterator{contract: _Thorchainrouter.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xef519b7eb82aaf6ac376a6df2d793843ebfd593de5f1a0601d3cc6ab49ebb395.
//
// Solidity: event Deposit(address indexed to, address indexed asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *ThorchainrouterDeposit, to []common.Address, asset []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _Thorchainrouter.contract.WatchLogs(opts, "Deposit", toRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThorchainrouterDeposit)
				if err := _Thorchainrouter.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xef519b7eb82aaf6ac376a6df2d793843ebfd593de5f1a0601d3cc6ab49ebb395.
//
// Solidity: event Deposit(address indexed to, address indexed asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) ParseDeposit(log types.Log) (*ThorchainrouterDeposit, error) {
	event := new(ThorchainrouterDeposit)
	if err := _Thorchainrouter.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThorchainrouterTransferAllowanceIterator is returned from FilterTransferAllowance and is used to iterate over the raw logs and unpacked data for TransferAllowance events raised by the Thorchainrouter contract.
type ThorchainrouterTransferAllowanceIterator struct {
	Event *ThorchainrouterTransferAllowance // Event containing the contract specifics and raw log

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
func (it *ThorchainrouterTransferAllowanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThorchainrouterTransferAllowance)
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
		it.Event = new(ThorchainrouterTransferAllowance)
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
func (it *ThorchainrouterTransferAllowanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThorchainrouterTransferAllowanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThorchainrouterTransferAllowance represents a TransferAllowance event raised by the Thorchainrouter contract.
type ThorchainrouterTransferAllowance struct {
	OldVault common.Address
	NewVault common.Address
	Asset    common.Address
	Amount   *big.Int
	Memo     string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferAllowance is a free log retrieval operation binding the contract event 0x05b90458f953d3fcb2d7fb25616a2fddeca749d0c47cc5c9832d0266b5346eea.
//
// Solidity: event TransferAllowance(address indexed oldVault, address indexed newVault, address asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) FilterTransferAllowance(opts *bind.FilterOpts, oldVault []common.Address, newVault []common.Address) (*ThorchainrouterTransferAllowanceIterator, error) {

	var oldVaultRule []interface{}
	for _, oldVaultItem := range oldVault {
		oldVaultRule = append(oldVaultRule, oldVaultItem)
	}
	var newVaultRule []interface{}
	for _, newVaultItem := range newVault {
		newVaultRule = append(newVaultRule, newVaultItem)
	}

	logs, sub, err := _Thorchainrouter.contract.FilterLogs(opts, "TransferAllowance", oldVaultRule, newVaultRule)
	if err != nil {
		return nil, err
	}
	return &ThorchainrouterTransferAllowanceIterator{contract: _Thorchainrouter.contract, event: "TransferAllowance", logs: logs, sub: sub}, nil
}

// WatchTransferAllowance is a free log subscription operation binding the contract event 0x05b90458f953d3fcb2d7fb25616a2fddeca749d0c47cc5c9832d0266b5346eea.
//
// Solidity: event TransferAllowance(address indexed oldVault, address indexed newVault, address asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) WatchTransferAllowance(opts *bind.WatchOpts, sink chan<- *ThorchainrouterTransferAllowance, oldVault []common.Address, newVault []common.Address) (event.Subscription, error) {

	var oldVaultRule []interface{}
	for _, oldVaultItem := range oldVault {
		oldVaultRule = append(oldVaultRule, oldVaultItem)
	}
	var newVaultRule []interface{}
	for _, newVaultItem := range newVault {
		newVaultRule = append(newVaultRule, newVaultItem)
	}

	logs, sub, err := _Thorchainrouter.contract.WatchLogs(opts, "TransferAllowance", oldVaultRule, newVaultRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThorchainrouterTransferAllowance)
				if err := _Thorchainrouter.contract.UnpackLog(event, "TransferAllowance", log); err != nil {
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

// ParseTransferAllowance is a log parse operation binding the contract event 0x05b90458f953d3fcb2d7fb25616a2fddeca749d0c47cc5c9832d0266b5346eea.
//
// Solidity: event TransferAllowance(address indexed oldVault, address indexed newVault, address asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) ParseTransferAllowance(log types.Log) (*ThorchainrouterTransferAllowance, error) {
	event := new(ThorchainrouterTransferAllowance)
	if err := _Thorchainrouter.contract.UnpackLog(event, "TransferAllowance", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThorchainrouterTransferOutIterator is returned from FilterTransferOut and is used to iterate over the raw logs and unpacked data for TransferOut events raised by the Thorchainrouter contract.
type ThorchainrouterTransferOutIterator struct {
	Event *ThorchainrouterTransferOut // Event containing the contract specifics and raw log

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
func (it *ThorchainrouterTransferOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThorchainrouterTransferOut)
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
		it.Event = new(ThorchainrouterTransferOut)
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
func (it *ThorchainrouterTransferOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThorchainrouterTransferOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThorchainrouterTransferOut represents a TransferOut event raised by the Thorchainrouter contract.
type ThorchainrouterTransferOut struct {
	Vault  common.Address
	To     common.Address
	Asset  common.Address
	Amount *big.Int
	Memo   string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransferOut is a free log retrieval operation binding the contract event 0xa9cd03aa3c1b4515114539cd53d22085129d495cb9e9f9af77864526240f1bf7.
//
// Solidity: event TransferOut(address indexed vault, address indexed to, address asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) FilterTransferOut(opts *bind.FilterOpts, vault []common.Address, to []common.Address) (*ThorchainrouterTransferOutIterator, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Thorchainrouter.contract.FilterLogs(opts, "TransferOut", vaultRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ThorchainrouterTransferOutIterator{contract: _Thorchainrouter.contract, event: "TransferOut", logs: logs, sub: sub}, nil
}

// WatchTransferOut is a free log subscription operation binding the contract event 0xa9cd03aa3c1b4515114539cd53d22085129d495cb9e9f9af77864526240f1bf7.
//
// Solidity: event TransferOut(address indexed vault, address indexed to, address asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) WatchTransferOut(opts *bind.WatchOpts, sink chan<- *ThorchainrouterTransferOut, vault []common.Address, to []common.Address) (event.Subscription, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Thorchainrouter.contract.WatchLogs(opts, "TransferOut", vaultRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThorchainrouterTransferOut)
				if err := _Thorchainrouter.contract.UnpackLog(event, "TransferOut", log); err != nil {
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

// ParseTransferOut is a log parse operation binding the contract event 0xa9cd03aa3c1b4515114539cd53d22085129d495cb9e9f9af77864526240f1bf7.
//
// Solidity: event TransferOut(address indexed vault, address indexed to, address asset, uint256 amount, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) ParseTransferOut(log types.Log) (*ThorchainrouterTransferOut, error) {
	event := new(ThorchainrouterTransferOut)
	if err := _Thorchainrouter.contract.UnpackLog(event, "TransferOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThorchainrouterTransferOutAndCallIterator is returned from FilterTransferOutAndCall and is used to iterate over the raw logs and unpacked data for TransferOutAndCall events raised by the Thorchainrouter contract.
type ThorchainrouterTransferOutAndCallIterator struct {
	Event *ThorchainrouterTransferOutAndCall // Event containing the contract specifics and raw log

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
func (it *ThorchainrouterTransferOutAndCallIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThorchainrouterTransferOutAndCall)
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
		it.Event = new(ThorchainrouterTransferOutAndCall)
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
func (it *ThorchainrouterTransferOutAndCallIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThorchainrouterTransferOutAndCallIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThorchainrouterTransferOutAndCall represents a TransferOutAndCall event raised by the Thorchainrouter contract.
type ThorchainrouterTransferOutAndCall struct {
	Vault        common.Address
	Target       common.Address
	Amount       *big.Int
	FinalAsset   common.Address
	To           common.Address
	AmountOutMin *big.Int
	Memo         string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferOutAndCall is a free log retrieval operation binding the contract event 0x8e5841bcd195b858d53b38bcf91b38d47f3bc800469b6812d35451ab619c6f6c.
//
// Solidity: event TransferOutAndCall(address indexed vault, address target, uint256 amount, address finalAsset, address to, uint256 amountOutMin, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) FilterTransferOutAndCall(opts *bind.FilterOpts, vault []common.Address) (*ThorchainrouterTransferOutAndCallIterator, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Thorchainrouter.contract.FilterLogs(opts, "TransferOutAndCall", vaultRule)
	if err != nil {
		return nil, err
	}
	return &ThorchainrouterTransferOutAndCallIterator{contract: _Thorchainrouter.contract, event: "TransferOutAndCall", logs: logs, sub: sub}, nil
}

// WatchTransferOutAndCall is a free log subscription operation binding the contract event 0x8e5841bcd195b858d53b38bcf91b38d47f3bc800469b6812d35451ab619c6f6c.
//
// Solidity: event TransferOutAndCall(address indexed vault, address target, uint256 amount, address finalAsset, address to, uint256 amountOutMin, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) WatchTransferOutAndCall(opts *bind.WatchOpts, sink chan<- *ThorchainrouterTransferOutAndCall, vault []common.Address) (event.Subscription, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Thorchainrouter.contract.WatchLogs(opts, "TransferOutAndCall", vaultRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThorchainrouterTransferOutAndCall)
				if err := _Thorchainrouter.contract.UnpackLog(event, "TransferOutAndCall", log); err != nil {
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

// ParseTransferOutAndCall is a log parse operation binding the contract event 0x8e5841bcd195b858d53b38bcf91b38d47f3bc800469b6812d35451ab619c6f6c.
//
// Solidity: event TransferOutAndCall(address indexed vault, address target, uint256 amount, address finalAsset, address to, uint256 amountOutMin, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) ParseTransferOutAndCall(log types.Log) (*ThorchainrouterTransferOutAndCall, error) {
	event := new(ThorchainrouterTransferOutAndCall)
	if err := _Thorchainrouter.contract.UnpackLog(event, "TransferOutAndCall", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThorchainrouterVaultTransferIterator is returned from FilterVaultTransfer and is used to iterate over the raw logs and unpacked data for VaultTransfer events raised by the Thorchainrouter contract.
type ThorchainrouterVaultTransferIterator struct {
	Event *ThorchainrouterVaultTransfer // Event containing the contract specifics and raw log

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
func (it *ThorchainrouterVaultTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThorchainrouterVaultTransfer)
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
		it.Event = new(ThorchainrouterVaultTransfer)
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
func (it *ThorchainrouterVaultTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThorchainrouterVaultTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThorchainrouterVaultTransfer represents a VaultTransfer event raised by the Thorchainrouter contract.
type ThorchainrouterVaultTransfer struct {
	OldVault common.Address
	NewVault common.Address
	Coins    []THORChainRouterCoin
	Memo     string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterVaultTransfer is a free log retrieval operation binding the contract event 0x281daef48d91e5cd3d32db0784f6af69cd8d8d2e8c612a3568dca51ded51e08f.
//
// Solidity: event VaultTransfer(address indexed oldVault, address indexed newVault, (address,uint256)[] coins, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) FilterVaultTransfer(opts *bind.FilterOpts, oldVault []common.Address, newVault []common.Address) (*ThorchainrouterVaultTransferIterator, error) {

	var oldVaultRule []interface{}
	for _, oldVaultItem := range oldVault {
		oldVaultRule = append(oldVaultRule, oldVaultItem)
	}
	var newVaultRule []interface{}
	for _, newVaultItem := range newVault {
		newVaultRule = append(newVaultRule, newVaultItem)
	}

	logs, sub, err := _Thorchainrouter.contract.FilterLogs(opts, "VaultTransfer", oldVaultRule, newVaultRule)
	if err != nil {
		return nil, err
	}
	return &ThorchainrouterVaultTransferIterator{contract: _Thorchainrouter.contract, event: "VaultTransfer", logs: logs, sub: sub}, nil
}

// WatchVaultTransfer is a free log subscription operation binding the contract event 0x281daef48d91e5cd3d32db0784f6af69cd8d8d2e8c612a3568dca51ded51e08f.
//
// Solidity: event VaultTransfer(address indexed oldVault, address indexed newVault, (address,uint256)[] coins, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) WatchVaultTransfer(opts *bind.WatchOpts, sink chan<- *ThorchainrouterVaultTransfer, oldVault []common.Address, newVault []common.Address) (event.Subscription, error) {

	var oldVaultRule []interface{}
	for _, oldVaultItem := range oldVault {
		oldVaultRule = append(oldVaultRule, oldVaultItem)
	}
	var newVaultRule []interface{}
	for _, newVaultItem := range newVault {
		newVaultRule = append(newVaultRule, newVaultItem)
	}

	logs, sub, err := _Thorchainrouter.contract.WatchLogs(opts, "VaultTransfer", oldVaultRule, newVaultRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThorchainrouterVaultTransfer)
				if err := _Thorchainrouter.contract.UnpackLog(event, "VaultTransfer", log); err != nil {
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

// ParseVaultTransfer is a log parse operation binding the contract event 0x281daef48d91e5cd3d32db0784f6af69cd8d8d2e8c612a3568dca51ded51e08f.
//
// Solidity: event VaultTransfer(address indexed oldVault, address indexed newVault, (address,uint256)[] coins, string memo)
func (_Thorchainrouter *ThorchainrouterFilterer) ParseVaultTransfer(log types.Log) (*ThorchainrouterVaultTransfer, error) {
	event := new(ThorchainrouterVaultTransfer)
	if err := _Thorchainrouter.contract.UnpackLog(event, "VaultTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
