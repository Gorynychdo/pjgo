package pjgo

import (
	"fmt"
	"github.com/Gorynychdo/pjgo/internal/pjsua2"
	"strings"
	"sync"
)

type SipService struct {
	endpoint       pjsua2.Endpoint
	activeAccounts map[string]pjsua2.Account
	activeCalls    map[string]pjsua2.Call
	sipUser        ISipService // application callback
}

var (
	mutex     sync.Mutex
	logWriter = pjsua2.NewDirectorLogWriter(new(SipLogWriter))
)

func NewSipService(sipUser ISipService) *SipService {
	sipService := SipService{}
	sipService.sipUser = sipUser
	sipService.init()
	return &sipService
}

func (ss *SipService) init() {
	ss.endpoint = pjsua2.NewEndpoint()
	ss.activeAccounts = make(map[string]pjsua2.Account)
	ss.activeCalls = make(map[string]pjsua2.Call)

	// Create endpoint
	ss.endpoint.LibCreate()

	// Init library
	epConfig := pjsua2.NewEpConfig()
	epConfig.GetLogConfig().SetLevel(4)
	epConfig.GetLogConfig().SetWriter(logWriter)
	ss.endpoint.LibInit(epConfig)
	ss.endpoint.AudDevManager().SetNullDev()

	// Transport
	transportConfig := pjsua2.NewTransportConfig()
	transportConfig.SetPort(5060)
	ss.endpoint.TransportCreate(pjsua2.PJSIP_TRANSPORT_UDP, transportConfig)

	// Start library
	ss.endpoint.LibStart()

	fmt.Printf("[ SipService ] Available codecs:\n")
	for i := 0; i < int(ss.endpoint.CodecEnum().Size()); i++ {
		c := ss.endpoint.CodecEnum().Get(i)
		fmt.Printf("\t - %s (priority: %d)\n", c.GetCodecId(), c.GetPriority())
	}

	fmt.Printf("[ SipService ] PJSUA2 STARTED ***\n")
}

func (ss *SipService) RegisterAccount(config *Config) string {
	ss.checkThread()
	fmt.Printf("[ SipService ] Registration start, user=%v\n", config.Login)
	account := ss.createLocalAccount(config)
	ss.activeAccounts[config.Login] = account

	return config.Login
}

func (ss *SipService) Unregister(accountId string) {
	ss.checkThread()

	account := ss.activeAccounts[accountId]
	if account == nil {
		return
	}

	fmt.Printf("[ SipService ] Un-Registration start, user=%v\n", accountId)
	account.SetRegistration(false)
}

func (ss *SipService) MakeCall(accountId string, remoteUser string) string {
	ss.checkThread()

	account := ss.activeAccounts[accountId]
	if account == nil {
		fmt.Printf("[ SipService ] makeCall error : first use create_account or register_account\n")
		return ""
	}

	return ss.makeCallWithAccount(account, remoteUser)
}

func (ss *SipService) makeCallWithAccount(account pjsua2.Account, remoteUser string) string {
	ss.checkThread()

	remoteUri := ss.getRemoteURI(remoteUser)

	// Make outgoing call
	sipCall := NewSipCall(ss)
	call := pjsua2.NewDirectorCall(sipCall, account)
	sipCall.call = call
	callOpParam := pjsua2.NewCallOpParam(true)
	callOpParam.GetOpt().SetAudioCount(1)

	call.MakeCall(remoteUri, callOpParam)
	ci := call.GetInfo()
	ss.activeCalls[ci.GetCallIdString()] = call
	fmt.Printf("[ SipService ] Make Call, From = %v, To = %v, callId = %v\n",
		account.GetInfo().GetUri(), remoteUri, ci.GetCallIdString())
	return ci.GetCallIdString()
}

func (ss *SipService) createLocalAccount(config *Config) pjsua2.Account {
	sipAccount := pjsua2.NewDirectorAccount(NewSipAccount(config.Login, ss))

	accountConfig := pjsua2.NewAccountConfig()
	accountConfig.SetIdUri(config.Id)
	accountConfig.GetRegConfig().SetRegistrarUri(config.Uri)
	cred := pjsua2.NewAuthCredInfo("digest", "*", config.Login, 0, config.Password)
	accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	sipAccount.Create(accountConfig)

	fmt.Printf("[ SipService ] Account Created, URI = %v\n", sipAccount.GetInfo().GetUri())

	return sipAccount
}

func (ss *SipService) getRemoteURI(remoteUser string) string {
	// remoteURI
	remoteUri := strings.Builder{}

	remoteUri.WriteString("sip:")
	remoteUri.WriteString(remoteUser)
	remoteUri.WriteString("@pjsip.org:5060;transport=udp")

	return remoteUri.String()
}

func (ss *SipService) getAccount(user string) pjsua2.Account {
	return ss.activeAccounts[user]
}

func (ss *SipService) addCall(callIdString string, call pjsua2.Call) {
	ss.activeCalls[callIdString] = call
}

func (ss *SipService) removeCall(callIdString string) {
	call := ss.activeCalls[callIdString]
	if call != nil {
		fmt.Printf("[ SipService ] Remove Call, callId = %v\n", callIdString)
		delete(ss.activeCalls, callIdString)
		fmt.Printf("[ SipService ] Active Calls = %v\n", len(ss.activeCalls))
	}
}

func (ss *SipService) onRegState(uri string, isActive bool, code pjsua2.Pjsip_status_code) {
	ss.sipUser.OnRegState(uri, isActive, code)
}

func (ss *SipService) onIncomingCall(callIdString string, from string, to string) interface{} {
	return ss.sipUser.OnIncomingCall(callIdString, from, to)
}

func (ss *SipService) checkThread() {
	mutex.Lock()
	defer mutex.Unlock()

	if !ss.endpoint.LibIsThreadRegistered() {
		ss.endpoint.LibRegisterThread("")
	}
}
