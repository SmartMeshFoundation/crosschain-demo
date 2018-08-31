package service

import (
	"github.com/SmartMeshFoundation/Atmosphere/lndapi"
	"github.com/SmartMeshFoundation/Atmosphere/smapi"
)

// api of smartraiden
var SmAPI *smapi.SmAPI

// api of lnd
var LndAPI *lndapi.LndAPI

// InitAPI :
func InitAPI(smHost, lndHost string) {
	SmAPI = smapi.NewSmAPI(smHost)
	LndAPI = lndapi.NewLndAPI(lndHost)
}
