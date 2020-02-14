package main

import (
	"fmt"
	"github.com/Gorynychdo/pjgo/internal/app/pjgo2"
	"github.com/Gorynychdo/pjgo/internal/pjsua2"
	"os"
	"os/signal"
)

type SipUser struct {
	sipService *pjgo2.SipService
	callId     string
}

func (su *SipUser) OnRegState(userId string, isActive bool, code pjsua2.Pjsip_status_code) {
	fmt.Printf("[ OnRegState ] userId=%v, isActive=%v, code=%v\n", userId, isActive, code)
	if isActive {
		su.callId = su.sipService.MakeCall("test1", "test1")
	}
}

func (su *SipUser) OnIncomingCall(callIdString string, from string, to string) interface{} {
	su.callId = callIdString
	return "user"
}

func main() {

	sipUser := SipUser{}
	sipService := pjgo2.NewSipService(&sipUser)
	sipUser.sipService = sipService

	sipService.RegisterAccount("test1", "test1")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
