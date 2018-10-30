package service

import (
	"github.com/SmartMeshFoundation/crosschain-demo/lndapi"
	"github.com/SmartMeshFoundation/crosschain-demo/photonapi"
)

// SmAPI api of photon
var SmAPI *photonapi.SmAPI

// LndAPI api of lnd
var LndAPI *lndapi.LndAPI

// InitAPI :
func InitAPI(smHost, lndHost string) {
	SmAPI = photonapi.NewSmAPI(smHost)
	LndAPI = lndapi.NewLndAPI(lndHost)
}
