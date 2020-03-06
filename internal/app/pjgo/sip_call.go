package pjgo

import (
	"fmt"
	"github.com/Gorynychdo/pjgo/internal/pkg/pjsua2"
)

type SipCall struct {
	sipService *SipService
	call       pjsua2.Call
}

func NewSipCall(sipService *SipService) *SipCall {
	return &SipCall{sipService, nil}
}

func (sc *SipCall) OnCallState(prm pjsua2.OnCallStateParam) {
	ci := sc.call.GetInfo()

	fmt.Printf("[ SipCall ] onCallState %v, aor = %v\n", ci.GetStateText(), ci.GetRemoteUri())

	if ci.GetState() == pjsua2.PJSIP_INV_STATE_DISCONNECTED {
		fmt.Printf("[ SipCall ] Call Closed, CallId=%v, AOR=%v, reason=%v, lastStatusCode=%v\n",
			ci.GetCallIdString(), ci.GetRemoteUri(),
			ci.GetLastReason(), ci.GetLastStatusCode())

		sc.sipService.call = nil
		pjsua2.DeleteCall(sc.call)
	}
}
