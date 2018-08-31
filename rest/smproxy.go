package rest

import (
	"log"

	"encoding/json"

	"net/http"

	"github.com/SmartMeshFoundation/Atmosphere/service"
	"github.com/ant0ine/go-json-rest/rest"
)

// GetBalanceOnSm :
func GetBalanceOnSm(w rest.ResponseWriter, r *rest.Request) {
	statueCode, body, err := service.SmAPI.GetBalanceByTokenAddress(r.PathParam("tokenaddress"))
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var buf interface{}
	err = json.Unmarshal(body, &buf)
	w.WriteHeader(statueCode)
	err = w.WriteJson(buf)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
