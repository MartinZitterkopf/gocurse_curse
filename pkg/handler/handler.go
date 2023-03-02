package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MartinZitterkopf/gocurse_library_response/response"
	"github.com/MartinZitterkopf/gocurse_microservice_curse/internal/curse"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewUserHTTPServer(ctx context.Context, endpoints curse.Endpoints) http.Handler {

	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("/curses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeStoreCurse,
		encodeResponse,
		opts...,
	)).Methods("POST")

	router.Handle("/curses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/curses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetByID),
		decodeGetByIDCurse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/curses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCurse,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	router.Handle("/curses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCurse,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return router
}

func decodeStoreCurse(_ context.Context, r *http.Request) (interface{}, error) {

	var req curse.CreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeCreateCurse(_ context.Context, r *http.Request) (interface{}, error) {

	var req curse.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v", err.Error()))
	}

	return req, nil
}

func decodeGetByIDCurse(_ context.Context, r *http.Request) (interface{}, error) {

	p := mux.Vars(r)
	req := curse.GetByIDReq{
		ID: p["id"],
	}

	return req, nil
}

func decodeGetAllUser(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := curse.GetAllReq{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil
}

func decodeUpdateCurse(_ context.Context, r *http.Request) (interface{}, error) {

	var req curse.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v", err.Error()))
	}

	path := mux.Vars(r)
	req.ID = path["id"]

	return req, nil
}

func decodeDeleteCurse(_ context.Context, r *http.Request) (interface{}, error) {

	path := mux.Vars(r)
	req := curse.DeleteReq{
		ID: path["id"],
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {

	r := resp.(response.Response)
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
