package service

import (
	"testing"

	"math/big"

	"time"

	"fmt"

	"github.com/SmartMeshFoundation/Atmosphere/utils"
	"github.com/ethereum/go-ethereum/common"
)

var BobSmAddress = common.HexToAddress("0x2b0C1545DBBEC6BFe7B26c699b74EB3513e52724")
var BobLndAddress = "03ffe99920feaa06a13d20544fb14650833cc8664485c317fbc675dbc6d4bff642"
var BobSmHost = "http://localhost:2001"
var BobLndHost = "localhost:10002"

var AliceSmAddress = common.HexToAddress("0xaaAA7F676a677c0B3C8E4Bb14aEC7Be61365acfE")
var AliceSmHost = "http://localhost:3001"
var AliceLndHost = "localhost:10001"

var smTokenAddress = common.HexToAddress("0x2D5D1FD0509eBEDc96aB825B1e4D4104AcA493be")
var smAmount = big.NewInt(2)
var lndAmount = big.NewInt(1000)
var secret = utils.NewRandomHash()
var lockSecretHash = utils.ShaSecret(secret.Bytes())
var lockSecretHash2 = common.HexToHash("0x4454f0779f09e9509a15c3076e92d2225a85294a071369f2fb17cd32662cf77e")

func TestRegisterExchangeStateBySmSender(t *testing.T) {
	// Bob
	fmt.Println("secret", secret.String())
	fmt.Println("lockSecretHash", lockSecretHash.String())
	InitAPI(BobSmHost, BobLndHost)
	RegisterExchangeStateBySmSender(AliceSmAddress, smTokenAddress, smAmount, lndAmount, secret)
	time.Sleep(1000 * time.Second)
}

func TestRegisterExchangeStateByLndSender(t *testing.T) {
	// Alice
	InitAPI(AliceSmHost, AliceLndHost)
	RegisterExchangeStateByLndSender(BobLndAddress, smTokenAddress, smAmount, lndAmount, lockSecretHash2)
	time.Sleep(1000 * time.Second)
}
