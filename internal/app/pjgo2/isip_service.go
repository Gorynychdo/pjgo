package pjgo2

import "github.com/Gorynychdo/pjgo/internal/pjsua2"

type ISipService interface {
	OnRegState(userId string, isActive bool, code pjsua2.Pjsip_status_code)
	OnIncomingCall(callIdString string, from string, to string) interface{}
}
