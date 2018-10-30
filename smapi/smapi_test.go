package smapi

import (
	"testing"

	"fmt"

	"math/big"

	"time"

	"github.com/SmartMeshFoundation/crosschain-demo/utils"
	"github.com/ethereum/go-ethereum/common"
)

var target = common.HexToAddress("0x2b0C1545DBBEC6BFe7B26c699b74EB3513e52724")
var tokenAddress = common.HexToAddress("0x80E010c563024EDfc3a4efB4fB0dE54E728AcC9d")
var secret = common.HexToHash("0x64e604787cbf194841e7b68d7cd28786f6c9a0a3ab9f8b0a0e87cb4387ab0107")
var lockSecretHash = common.HexToHash("0x55c6593b6ab2d834d1e9b89cc9cdd3c866bf319b04ef6b08548dd85c031376a8")
var sm = NewSmAPI("http://localhost:2001")

func TestSmAPI_GetUnfinishedReceivedTransfer(t *testing.T) {
	jsonStr, err := sm.GetUnfinishedReceivedTransfer(tokenAddress, lockSecretHash)
	fmt.Println("err", err)
	fmt.Println(jsonStr)
	fmt.Println(jsonStr == "null")
}

func TestSmAPI_AllowRevealSecret(t *testing.T) {
	err := sm.AllowRevealSecret(tokenAddress, lockSecretHash)
	fmt.Println("err", err)
}

func TestSmAPI_SendTransferWithSecretAsync(t *testing.T) {
	sm.SendTransferWithSecretAsync(target, tokenAddress, big.NewInt(5), secret)
	time.Sleep(3 * time.Second)
}

func TestSmAPI_GetNodeAddress(t *testing.T) {
	address, err := sm.GetNodeAddress()
	fmt.Println("err", err)
	fmt.Println("address", address.String())
}

func TestSmAPI_GetBalanceByTokenAddress(t *testing.T) {
	resp, err := sm.GetBalanceByTokenAddress(utils.EmptyAddress.String())
	fmt.Println("err", err)
	fmt.Println(utils.ToFormatJSONString(resp))
}
