package pjgo

import (
	"fmt"
	"github.com/Gorynychdo/pjgo/internal/pkg/pjsua2"
)

type SipAccount struct {
	userId     string
	sipService *SipService
}

func NewSipAccount(userId string, sipService *SipService) *SipAccount {
	ac := &SipAccount{userId, sipService}
	return ac
}

func (sa *SipAccount) OnRegState(prm pjsua2.OnRegStateParam) {
	sa.sipService.checkThread()
	info := sa.sipService.account.GetInfo()

	var regState string
	if info.GetRegIsActive() {
		regState = "REGISTER"
	} else {
		regState = "UNREGISTER"
	}
	fmt.Printf("[ SipAccount ] %v : code = %v\n", regState, prm.GetCode())

	fmt.Printf("[ OnRegState ] userId=%v, isActive=%v, code=%v\n",
		info.GetUri(), info.GetRegIsActive(), prm.GetCode())

	if sa.sipService.config.MakeCall && info.GetRegIsActive() {
		sa.sipService.MakeCall(sa.sipService.config.RemoteUri)
	}
}

func (sa *SipAccount) OnIncomingCall(prm pjsua2.OnIncomingCallParam) {
	account := sa.sipService.account
	sipCall := NewSipCall(sa.sipService)
	call := pjsua2.NewDirectorCall(sipCall, account, prm.GetCallId())
	sipCall.call = call
	callInfo := call.GetInfo()

	fmt.Printf("[ SipAccount ] IncomingCall %v\n"+
		"...remoteUri = %v, localUri = %v\n",
		prm.GetRdata().GetInfo(), callInfo.GetRemoteUri(), callInfo.GetLocalUri())

	sa.sipService.call = call

	callOpParam := pjsua2.NewCallOpParam()
	callOpParam.SetStatusCode(pjsua2.PJSIP_SC_OK)
	callOpParam.GetOpt().SetAudioCount(1)
	callOpParam.GetOpt().SetVideoCount(0)
	call.Answer(callOpParam)
}
