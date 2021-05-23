package main

import (
	"github.com/Joddev/autop2p/honestfund"
	"github.com/Joddev/autop2p/peoplefund"
	"net/http"
)

var Client = &http.Client{}

var HonestfundApi = honestfund.NewApi(Client)
var HonestfundService = honestfund.NewService(HonestfundApi)

var PeoplefundApi = peoplefund.NewApi(Client)
var PeoplefundService = peoplefund.NewService(PeoplefundApi)
