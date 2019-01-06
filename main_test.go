package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RequestHandler(res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal("{'name':'roushan','lat':13.63229, 'long':77.228, 'cabtype':'normal'}")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func ConfirmRequestHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Authorization", "Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJjaWQiOjUxLCJpYXQiOjE1NDY3NjkwMTcsImp0aSI6IjE1NDY3NjkwMTdmdWJlciIsInVpZCI6NDB9.vMjLLNbrkNMXr7FVpFHI-TYIjf2XgWXSaD8dceewdgFoPVxogXavXq1Fb-5uRE8R5Ve33hnBJla4AjUhIX1vgDiQIZBOxOMsdzo_Cp6pVNamfiC5uV6DzkAi7Gfny6AuI5kFEC_3N46T6qcgN-3su203J0AQ9qevwdFzttgy648f18Q")

}

func EndTripHandler(res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal("{'lat':13.444, 'long':77.222}")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Authorization", "Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJjaWQiOjUxLCJpYXQiOjE1NDY3NjkwMTcsImp0aSI6IjE1NDY3NjkwMTdmdWJlciIsInVpZCI6NDB9.vMjLLNbrkNMXr7FVpFHI-TYIjf2XgWXSaD8dceewdgFoPVxogXavXq1Fb-5uRE8R5Ve33hnBJla4AjUhIX1vgDiQIZBOxOMsdzo_Cp6pVNamfiC5uV6DzkAi7Gfny6AuI5kFEC_3N46T6qcgN-3su203J0AQ9qevwdFzttgy648f18Q")

	res.Write(data)
}

func TestHandlerRequestCab(t *testing.T) {
	request, _ := http.NewRequest("POST", "/cabrequest", nil)
	response := httptest.NewRecorder()

	RequestHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
	}
}

func TestHandlerConfirmCab(t *testing.T) {
	request, _ := http.NewRequest("GET", "/confirmride", nil)
	response := httptest.NewRecorder()

	ConfirmRequestHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
	}
}
func TestHandlerEndTrip(t *testing.T) {
	request, _ := http.NewRequest("POST", "/endride", nil)
	response := httptest.NewRecorder()

	EndTripHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
	}
}
