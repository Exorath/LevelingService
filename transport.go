package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"strconv"
	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(ctx context.Context, s LevelingService) http.Handler {

	r := mux.NewRouter()
	r.Methods("GET").Path("/players/{uuid}").Handler(httptransport.NewServer(
		ctx,
		makeGetAccountEndpoint(s),
		decodeGetAccountRequest,
		encodeResponse,
	))

	r.Methods("POST").Path("/players/{uuid}/xp/{xp}").Handler(httptransport.NewServer(
		ctx,
		makeAddExperienceEndpoint(s),
		decodeAddExperienceRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/players/{uuid}/consume/{lvl}").Handler(httptransport.NewServer(
		ctx,
		makeConsumeLevelEndpoint(s),
		decodeConsumeLevelRequest,
		encodeResponse,
	))
	return r
}
func makeAddExperienceEndpoint(svc LevelingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addExperienceRequest)
		success, err := svc.addExperience(req.Uuid, req.Experience)
		return addExperienceResponse{success, err}, nil
	}
}
func decodeAddExperienceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	xp, err := strconv.Atoi(vars["xp"])
	var request addExperienceRequest =
		addExperienceRequest{vars["uuid"], xp}
	return request, err
}

func makeGetAccountEndpoint(svc LevelingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountRequest)
		acc, err := svc.getAccount(req.Uuid)
		return getAccountResponse{acc, err}, nil
	}
}

func decodeGetAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var request getAccountRequest =
		getAccountRequest{vars["uuid"]}
	return request, nil
}

func makeConsumeLevelEndpoint(svc LevelingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(consumeLevelRequest)
		success, err := svc.consumeLevel(req.Uuid, req.Level)
		return consumeLevelResponse{success, err}, nil
	}
}

func decodeConsumeLevelRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	lvl, errLvl := strconv.Atoi(vars["lvl"])
	if(errLvl != nil){
		return nil, errLvl
	}
	var request consumeLevelRequest =
		consumeLevelRequest{uuid, lvl}
	return request, nil
}



func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

//Requests and responses
//addExperience
type addExperienceRequest struct {
	Uuid string `json:"uuid"`
	Experience int `json:"xp"`
}

type addExperienceResponse struct {
	Success bool `json:"s"`
	Err error `json:"err,omitempty"` // errors don't define JSON marshaling
}
//getAccount
type getAccountRequest struct {
	Uuid string `json:"uuid"`
}

type getAccountResponse struct {
	Account LevelAccount `json:"acc"`
	Err error `json:"err,omitempty"` // errors don't define JSON marshaling
}

//consumeLevel
type consumeLevelRequest struct {
	Uuid string `json:"uuid"`
	Level int `json:"lvl"`
}

type consumeLevelResponse struct {
	Success bool `json:"s"`
	Err error `json:"err,omitempty"` // errors don't define JSON marshaling
}