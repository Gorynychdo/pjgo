package pjgo

import (
	"fmt"
	"github.com/Gorynychdo/pjgo/internal/pkg/pjsua2"
	"sync"
)

type SipService struct {
	endpoint pjsua2.Endpoint
	account  pjsua2.Account
	call     pjsua2.Call
	config   *Config
}

var (
	mutex     sync.Mutex
	logWriter = pjsua2.NewDirectorLogWriter(new(SipLogWriter))
)

func NewSipService(config *Config) *SipService {
	sipService := SipService{
		config: config,
	}
	sipService.init()
	return &sipService
}

func (ss *SipService) init() {
	ss.endpoint = pjsua2.NewEndpoint()
	ss.endpoint.LibCreate()

	epConfig := pjsua2.NewEpConfig()
	epConfig.GetLogConfig().SetLevel(ss.config.LogLevel)
	epConfig.GetLogConfig().SetWriter(logWriter)

	transportConfig := pjsua2.NewTransportConfig()
	transportConfig.SetPort(ss.config.Port)

	ss.endpoint.LibInit(epConfig)
	ss.endpoint.AudDevManager().SetNullDev()
	ss.endpoint.TransportCreate(pjsua2.PJSIP_TRANSPORT_UDP, transportConfig)
	ss.endpoint.LibStart()

	fmt.Printf("[ SipService ] Available codecs:\n")
	for i := 0; i < int(ss.endpoint.CodecEnum().Size()); i++ {
		c := ss.endpoint.CodecEnum().Get(i)
		fmt.Printf("\t - %s (priority: %d)\n", c.GetCodecId(), c.GetPriority())
	}

	fmt.Printf("[ SipService ] PJSUA2 STARTED ***\n")
}

func (ss *SipService) RegisterAccount() {
	ss.checkThread()
	fmt.Printf("[ SipService ] Registration start, user=%v\n", ss.config.Login)

	ss.account = pjsua2.NewDirectorAccount(NewSipAccount(ss.config.Login, ss))

	accountConfig := pjsua2.NewAccountConfig()
	accountConfig.SetIdUri(ss.config.Id)
	accountConfig.GetRegConfig().SetRegistrarUri(ss.config.Uri)
	cred := pjsua2.NewAuthCredInfo("digest", "*", ss.config.Login, 0, ss.config.Password)
	accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	ss.account.Create(accountConfig)

	fmt.Printf("[ SipService ] Account Created, URI = %v\n", ss.account.GetInfo().GetUri())
}

func (ss *SipService) MakeCall(remoteUri string) string {
	ss.checkThread()

	// Make outgoing call
	sipCall := NewSipCall(ss)
	ss.call = pjsua2.NewDirectorCall(sipCall, ss.account)
	sipCall.call = ss.call
	callOpParam := pjsua2.NewCallOpParam(true)
	callOpParam.GetOpt().SetAudioCount(1)

	ss.call.MakeCall(remoteUri, callOpParam)
	ci := ss.call.GetInfo()
	fmt.Printf("[ SipService ] Make Call, From = %v, To = %v, callId = %v\n",
		ss.account.GetInfo().GetUri(), remoteUri, ci.GetCallIdString())
	return ci.GetCallIdString()
}

func (ss *SipService) checkThread() {
	mutex.Lock()
	defer mutex.Unlock()

	if !ss.endpoint.LibIsThreadRegistered() {
		ss.endpoint.LibRegisterThread("")
	}
}
