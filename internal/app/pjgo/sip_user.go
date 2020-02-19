package pjgo

import (
	"fmt"
	"github.com/Gorynychdo/pjgo/internal/pkg/pjsua2"
)

type SipUser struct {
	SipService *SipService
	callId     string
}

func NewSipUser() *SipUser {
	return &SipUser{}
}

func (su *SipUser) OnRegState(userId string, isActive bool, code pjsua2.Pjsip_status_code) {
	fmt.Printf("[ OnRegState ] userId=%v, isActive=%v, code=%v\n", userId, isActive, code)
	//if isActive {
	//	su.callId = su.SipService.MakeCall("test1", "test1")
	//}
}

func (su *SipUser) OnIncomingCall(callIdString string, from string, to string) interface{} {
	su.callId = callIdString
	return "user"
}
