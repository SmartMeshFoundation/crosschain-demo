package service

import (
	"github.com/SmartMeshFoundation/crosschain-demo/lndapi"
	"github.com/SmartMeshFoundation/crosschain-demo/smapi"
)

// SmAPI api of photon
var SmAPI *smapi.SmAPI

// LndAPI api of lnd
var LndAPI *lndapi.LndAPI

// InitAPI :
func InitAPI(smHost, lndHost string) {
	SmAPI = smapi.NewSmAPI(smHost)
	LndAPI = lndapi.NewLndAPI(lndHost)
}
