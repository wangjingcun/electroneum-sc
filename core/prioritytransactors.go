package core

import (
	"math/big"
	"strings"

	"github.com/electroneum/electroneum-sc/accounts/abi"
	"github.com/electroneum/electroneum-sc/common"
	"github.com/electroneum/electroneum-sc/contracts/prioritytransactors"
	"github.com/electroneum/electroneum-sc/core/vm"
	"github.com/electroneum/electroneum-sc/params"
)

func GetPriorityTransactors(blockNumber *big.Int, config *params.ChainConfig, evm *vm.EVM) (common.PriorityTransactorMap, error) {
	var (
		address  = config.GetPriorityTransactorsContractAddress(blockNumber)
		contract = vm.AccountRef(address)
		method   = "getTransactors"
		result   = make(common.PriorityTransactorMap)
	)

	if address != (common.Address{}) {
		// Check if contract code exists at the address. If it doesn't. We haven't deployed the contract yet, so no error needed.
		byteCode := evm.StateDB.GetCode(address)
		if len(byteCode) == 0 {
			return result, nil
		}

		contractABI, _ := abi.JSON(strings.NewReader(prioritytransactors.ETNPriorityTransactorsInterfaceABI))
		input, _ := contractABI.Pack(method)
		output, _, err := evm.StaticCall(contract, address, input, params.MaxGasLimit)
		if err != nil {
			return result, err
		}
		unpackResult, err := contractABI.Unpack(method, output)
		if err != nil {
			return result, err
		}
		transactorsMeta := abi.ConvertType(unpackResult[0], new([]prioritytransactors.ETNPriorityTransactorsInterfaceTransactorMeta)).(*[]prioritytransactors.ETNPriorityTransactorsInterfaceTransactorMeta)
		for _, t := range *transactorsMeta {
			result[common.HexToPublicKey(t.PublicKey)] = common.PriorityTransactor{
				IsGasPriceWaiver: t.IsGasPriceWaiver,
				EntityName:       t.Name,
			}
		}
	}
	return result, nil
}

func getPriorityTransactorByKey(blockNumber *big.Int, publicKey common.PublicKey, config *params.ChainConfig, evm *vm.EVM) (common.PriorityTransactor, bool) {
	var (
		address  = config.GetPriorityTransactorsContractAddress(blockNumber)
		contract = vm.AccountRef(address)
		method   = "getTransactorByKey"
	)

	if address != (common.Address{}) {
		contractABI, _ := abi.JSON(strings.NewReader(prioritytransactors.ETNPriorityTransactorsInterfaceABI))
		input, _ := contractABI.Pack(method, publicKey.ToUnprefixedHexString())
		output, _, err := evm.StaticCall(contract, address, input, params.MaxGasLimit)
		if err != nil {
			return common.PriorityTransactor{}, false
		}
		unpackResult, err := contractABI.Unpack(method, output)
		if err != nil {
			return common.PriorityTransactor{}, false
		}
		transactorMeta := abi.ConvertType(unpackResult[0], new(prioritytransactors.ETNPriorityTransactorsInterfaceTransactorMeta)).(*prioritytransactors.ETNPriorityTransactorsInterfaceTransactorMeta)

		return common.PriorityTransactor{
			EntityName:       transactorMeta.Name,
			IsGasPriceWaiver: transactorMeta.IsGasPriceWaiver,
		}, true

	}
	return common.PriorityTransactor{}, false
}
