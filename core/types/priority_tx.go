// Copyright Electroneum 2023
package types

import (
	"math/big"

	"github.com/electroneum/electroneum-sc/common"
)

type PriorityTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int
	Data       []byte
	AccessList AccessList

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// Electroneum Signature values
	PriorityV *big.Int `json:"priorityV" gencodec:"required"`
	PriorityR *big.Int `json:"priorityR" gencodec:"required"`
	PriorityS *big.Int `json:"priorityS" gencodec:"required"`
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *PriorityTx) copy() TxData {
	cpy := &PriorityTx{
		Nonce: tx.Nonce,
		To:    copyAddressPtr(tx.To),
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		Value:      new(big.Int),
		ChainID:    new(big.Int),
		GasTipCap:  new(big.Int),
		GasFeeCap:  new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
		PriorityV:  new(big.Int),
		PriorityR:  new(big.Int),
		PriorityS:  new(big.Int),
	}
	copy(cpy.AccessList, tx.AccessList)
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasTipCap != nil {
		cpy.GasTipCap.Set(tx.GasTipCap)
	}
	if tx.GasFeeCap != nil {
		cpy.GasFeeCap.Set(tx.GasFeeCap)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	if tx.PriorityV != nil {
		cpy.PriorityV.Set(tx.PriorityV)
	}
	if tx.PriorityR != nil {
		cpy.PriorityR.Set(tx.PriorityR)
	}
	if tx.PriorityS != nil {
		cpy.PriorityS.Set(tx.PriorityS)
	}
	return cpy
}

// accessors for innerTx.
func (tx *PriorityTx) txType() byte           { return PriorityTxType }
func (tx *PriorityTx) chainID() *big.Int      { return tx.ChainID }
func (tx *PriorityTx) accessList() AccessList { return tx.AccessList }
func (tx *PriorityTx) data() []byte           { return tx.Data }
func (tx *PriorityTx) gas() uint64            { return tx.Gas }
func (tx *PriorityTx) gasFeeCap() *big.Int    { return tx.GasFeeCap }
func (tx *PriorityTx) gasTipCap() *big.Int    { return tx.GasTipCap }
func (tx *PriorityTx) gasPrice() *big.Int     { return tx.GasFeeCap }
func (tx *PriorityTx) value() *big.Int        { return tx.Value }
func (tx *PriorityTx) nonce() uint64          { return tx.Nonce }
func (tx *PriorityTx) to() *common.Address    { return tx.To }

func (tx *PriorityTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}
func (tx *PriorityTx) rawPrioritySignatureValues() (v, r, s *big.Int) {
	return tx.PriorityV, tx.PriorityR, tx.PriorityS
}

func (tx *PriorityTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}

func (tx *PriorityTx) setPrioritySignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.PriorityV, tx.PriorityR, tx.PriorityS = chainID, v, r, s
}
