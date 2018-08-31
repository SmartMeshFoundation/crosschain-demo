package service

import (
	"github.com/SmartMeshFoundation/Atmosphere/lndapi"
	"github.com/SmartMeshFoundation/Atmosphere/smapi"
)

var smAPI *smapi.SmAPI
var lndAPI *lndapi.LndAPI

// InitAPI :
func InitAPI(smHost, lndHost string) {
	smAPI = smapi.NewSmAPI(smHost)
	lndAPI = lndapi.NewLndAPI(lndHost)
}
