package pjgo

import (
	"fmt"
	"github.com/Gorynychdo/pjgo/internal/pkg/pjsua2"
)

type SipCall struct {
	sipService *SipService
	call       pjsua2.Call
	recorder   pjsua2.AudioMediaRecorder
}

func NewSipCall(sipService *SipService) *SipCall {
	return &SipCall{sipService, nil, nil}
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

func (sc *SipCall) OnCallMediaState(prm pjsua2.OnCallMediaStateParam) {
	capMedDev := sc.sipService.endpoint.AudDevManager().GetCaptureDevMedia()

	if sc.recorder == nil {
		sc.recorder = pjsua2.NewAudioMediaRecorder()
		sc.recorder.CreateRecorder("media/record.wav")
		capMedDev.StartTransmit(sc.recorder)
		fmt.Println("[ SipCall ] Recorder is active")
	} else {
		sc.recorder.StopTransmit(capMedDev)
		sc.recorder = nil
		fmt.Println("[ SipCall ] Recorder is inactive")
	}
}
