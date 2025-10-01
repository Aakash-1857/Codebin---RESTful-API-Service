package main

import (
	"os"
	"testing"
)

var testApp application

func TestMain(m *testing.M){
	testApp=application{}
	os.Exit(m.Run())
}