package rest

import (
	"log"

	"net/http"

	"github.com/SmartMeshFoundation/Atmosphere/service"
	"github.com/SmartMeshFoundation/Atmosphere/smapi"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/lightningnetwork/lnd/lnrpc"
)

// GetBalance : get both balance on smartraiden and lnd
func GetBalance(w rest.ResponseWriter, r *rest.Request) {
	type Balance struct {
		SmBalance       []*smapi.SmBalanceResponseDetailDTO `json:"smartraiden_balance,omitempty"`
		SmBalanceError  string                              `json:"smart_raiden_balance_error,omitempty"`
		LndBalance      *lnrpc.ChannelBalanceResponse       `json:"lnd_balance,omitempty"`
		LndBalanceError string                              `json:"lnd_balance_error,omitempty"`
	}
	balance := new(Balance)

	// get balance on lnd
	if service.LndAPI != nil {
		lndBalance, err := service.LndAPI.ChannelBalance()
		if err != nil {
			log.Println(err.Error())
			balance.LndBalance = nil
			balance.LndBalanceError = err.Error()
		} else {
			balance.LndBalance = lndBalance
			balance.LndBalanceError = ""
		}
	} else {
		balance.LndBalanceError = "no connect to lnd"
	}

	// get balance on smartraiden
	if service.SmAPI != nil {
		smBalance, err := service.SmAPI.GetBalanceByTokenAddress("")
		if err != nil {
			log.Println(err.Error())
			balance.SmBalance = nil
			balance.SmBalanceError = err.Error()
		} else {
			balance.SmBalance = smBalance
			balance.SmBalanceError = ""
		}
	} else {
		balance.SmBalanceError = "no connect to smartraiden"
	}
	err := w.WriteJson(balance)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
