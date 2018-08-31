package rest

import (
	"log"
	"math/big"
	"net/http"

	"github.com/SmartMeshFoundation/Atmosphere/service"
	"github.com/SmartMeshFoundation/Atmosphere/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

const (
	roleSmSender  = "sm-sender"
	roleLndSender = "lnd-sender"
)

// RegisterExchange :
func RegisterExchange(w rest.ResponseWriter, r *rest.Request) {
	type RegisterExchangeRequestPayload struct {
		Role              string   `json:"role"`
		PartnerSmAddress  string   `json:"partner_sm_address"`
		PartnerLndAddress string   `json:"partner_lnd_address"`
		SmTokenAddress    string   `json:"sm_token_address"`
		SmAmount          *big.Int `json:"sm_amount"`
		LndAmount         *big.Int `json:"lnd_amount"`
		Secret            string   `json:"secret"`
		LockSecretHash    string   `json:"lock_secret_hash"`
	}
	payload := &RegisterExchangeRequestPayload{}
	err := r.DecodeJsonPayload(payload)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	partnerSmAddress := common.HexToAddress(payload.PartnerSmAddress)
	smTokenAddress := common.HexToAddress(payload.SmTokenAddress)

	// 参数校验
	if payload.SmAmount.Cmp(big.NewInt(0)) <= 0 {
		rest.Error(w, "Invalid sm_amount", http.StatusBadRequest)
		return
	}
	if payload.LndAmount.Cmp(big.NewInt(0)) <= 0 {
		rest.Error(w, "Invalid lnd_amount", http.StatusBadRequest)
		return
	}

	if roleSmSender == payload.Role {
		secret := common.HexToHash(payload.Secret)
		if secret == utils.EmptyHash {
			rest.Error(w, "Invalid secret", http.StatusBadRequest)
			return
		}
		_, err = service.RegisterExchangeStateBySmSender(partnerSmAddress, smTokenAddress, payload.SmAmount, payload.LndAmount, secret)
		if err != nil {
			log.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if roleLndSender == payload.Role {
		lockSecretHash := common.HexToHash(payload.LockSecretHash)
		if lockSecretHash == utils.EmptyHash {
			rest.Error(w, "Invalid lockSecretHash", http.StatusBadRequest)
			return
		}
		_, err = service.RegisterExchangeStateByLndSender(payload.PartnerLndAddress, smTokenAddress, payload.SmAmount, payload.LndAmount, lockSecretHash)
		if err != nil {
			log.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		rest.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}
}
