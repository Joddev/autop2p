package main

import (
	"github.com/Joddev/autop2p/honestfund"
	"net/http"
)

var Client = &http.Client{}

var HonestfundApi = honestfund.NewApi(Client)
var HonestfundService = honestfund.NewService(HonestfundApi)

