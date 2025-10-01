package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T){
	rr:=httptest.NewRecorder()
	req,err:=http.NewRequest(http.MethodGet,"/healthcheck",nil)
	if err!=nil{
		t.Fatal(err)
	}
	testApp.healthcheck(rr,req)
	if rr.Code!=http.StatusOK{
		t.Errorf("want %d; got %d",http.StatusOK,rr.Code)
	}
}