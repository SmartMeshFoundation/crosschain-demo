package smapi

import (
	"log"
	"math/big"

	"fmt"

	"errors"

	"encoding/json"

	"github.com/SmartMeshFoundation/Atmosphere/httphelper"
	"github.com/SmartMeshFoundation/Atmosphere/utils"
	"github.com/ethereum/go-ethereum/common"
)

// SmAPI : api for smartraiden
type SmAPI struct {
	Host           string
	AccountAddress common.Address
}

// NewSmAPI :
func NewSmAPI(host string) *SmAPI {
	sm := new(SmAPI)
	sm.Host = host
	_, err := sm.GetNodeAddress()
	if err != nil {
		panic(err)
	}
	return sm
}

// SendTransferWithSecret :
func (sm *SmAPI) SendTransferWithSecret(target common.Address, tokenAddress common.Address, amount *big.Int, secret common.Hash) {
	type TransferWithSecretRequestPayload struct {
		Target string   `json:"target_address"`
		Token  string   `json:"token_address"`
		Amount *big.Int `json:"amount"`
		Secret string   `json:"secret"` // 当用户想使用自己指定的密码,而非随机密码时使用
	}
	payload := &TransferWithSecretRequestPayload{
		Target: target.String(),
		Token:  tokenAddress.String(),
		Amount: amount,
		Secret: secret.String(),
	}
	fullURL := sm.Host + fmt.Sprintf("/api/1/transfers/%s/%s", tokenAddress.String(), target.String())
	payloadStr := utils.ToJSONString(payload)
	log.Printf("SmAPI---> SendTransferWithSecret : POST url=%s payload=%s\n", fullURL, payloadStr)
	statusCode, body, err := httphelper.PostJSON(fullURL, payloadStr)
	if err != nil {
		log.Printf("SmAPI---> SendTransferWithSecret FAILED : err=%s\n", err.Error())
		return
	}
	if statusCode != 200 {
		log.Printf("SmAPI---> SendTransferWithSecret FAILED : statusCode=%d body=\n%s\n", statusCode, string(body))
		err = errors.New(string(body))
		return
	}
	return
}

// SendTransferWithSecretAsync :
func (sm *SmAPI) SendTransferWithSecretAsync(target common.Address, tokenAddress common.Address, amount *big.Int, secret common.Hash) {
	go sm.SendTransferWithSecret(target, tokenAddress, amount, secret)
	return
}

// GetUnfinishedReceivedTransfer :
func (sm *SmAPI) GetUnfinishedReceivedTransfer(tokenAddress common.Address, lockSecretHash common.Hash) (jsonStr string, err error) {
	fullURL := sm.Host + fmt.Sprintf("/api/1/getunfinishedreceivedtransfer/%s/%s", tokenAddress.String(), lockSecretHash.String())
	log.Printf("SmAPI---> GetUnfinishedReceivedTransfer : GET url=%s \n", fullURL)
	statusCode, body, err := httphelper.Get(fullURL)
	if err != nil {
		log.Printf("SmAPI---> GetUnfinishedReceivedTransfer FAILED : err=%s\n", err.Error())
		return
	}
	if statusCode != 200 {
		log.Printf("SmAPI---> GetUnfinishedReceivedTransfer FAILED : statusCode=%d body=\n%s\n", statusCode, string(body))
		err = errors.New(string(body))
		return
	}
	jsonStr = string(body)
	return
}

// AllowRevealSecret :
func (sm *SmAPI) AllowRevealSecret(tokenAddress common.Address, lockSecretHash common.Hash) (err error) {
	type AllowRevealSecretRequestPayload struct {
		LockSecretHash string `json:"lock_secret_hash"`
		TokenAddress   string `json:"token_address"`
	}
	payload := &AllowRevealSecretRequestPayload{
		LockSecretHash: lockSecretHash.String(),
		TokenAddress:   tokenAddress.String(),
	}
	fullURL := sm.Host + "/api/1/transfers/allowrevealsecret"
	payloadStr := utils.ToJSONString(payload)
	log.Printf("SmAPI---> AllowRevealSecret : POST url=%s payload=%s\n", fullURL, payloadStr)
	statusCode, body, err := httphelper.PostJSON(fullURL, payloadStr)
	if err != nil {
		log.Printf("SmAPI---> AllowRevealSecret FAILED : err=%s\n", err.Error())
		return
	}
	if statusCode != 200 {
		log.Printf("SmAPI---> AllowRevealSecret FAILED : statusCode=%d body=\n%s\n", statusCode, string(body))
		err = errors.New(string(body))
		return
	}
	return
}

// GetNodeAddress :
func (sm *SmAPI) GetNodeAddress() (nodeAddress common.Address, err error) {
	type GetNodeAddressResponse struct {
		Address string `json:"our_address"`
	}
	fullURL := sm.Host + "/api/1/address"
	statusCode, body, err := httphelper.Get(fullURL)
	if err != nil {
		log.Printf("SmAPI---> GetNodeAddress FAILED : err=%s\n", err.Error())
		return
	}
	if statusCode != 200 {
		log.Printf("SmAPI---> GetNodeAddress FAILED : statusCode=%d body=\n%s\n", statusCode, string(body))
		err = errors.New(string(body))
		return
	}
	var resp GetNodeAddressResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}
	nodeAddress = common.HexToAddress(resp.Address)
	sm.AccountAddress = nodeAddress
	return
}

//GetBalanceByTokenAddress : proxy
func (sm *SmAPI) GetBalanceByTokenAddress(tokenAddressStr string) (statusCode int, body []byte, err error) {
	return httphelper.Get(sm.Host + "/api/1/balance/" + tokenAddressStr)
}
