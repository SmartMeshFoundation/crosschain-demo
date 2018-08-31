package smapi

import "math/big"

// SmBalanceResponseDetailDTO : account's balance of token on smartraiden
type SmBalanceResponseDetailDTO struct {
	TokenAddress string   `json:"token_address"`
	Balance      *big.Int `json:"balance"`
	LockedAmount *big.Int `json:"locked_amount"`
}
