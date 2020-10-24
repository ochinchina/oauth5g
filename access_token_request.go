package main

import (
	"bytes"
	"encoding/json"
	"github.com/ajg/form"
	log "github.com/sirupsen/logrus"
	"io"
)

var ALL_NFTYPES map[string]bool = map[string]bool{"NRF": true,
	"UDM":    true,
	"AMF":    true,
	"SMF":    true,
	"AUSF":   true,
	"NEF":    true,
	"PCF":    true,
	"SMSF":   true,
	"NSSF":   true,
	"UDR":    true,
	"LMF":    true,
	"GMLC":   true,
	"5G_EIR": true,
	"SEPP":   true,
	"UPF":    true,
	"N3IWF":  true,
	"AF":     true,
	"UDSF":   true,
	"BSF":    true,
	"CHF":    true,
	"NWDAF":  true,
	"PCSCF":  true,
	"HSS":    true,
	"UCMF":   true,
	"SOR_AF": true,
	"SPAF":   true,
	"MME":    true,
	"SCSAS":  true,
	"SCEF":   true,
	"SCP":    true,
	"NSSAAF": true,
	"ICSCF":  true,
	"SCSCF":  true,
}

// TS 29.510 Clause 6.1.6.3.11
var ALL_SERVICE_NAMES map[string]bool = map[string]bool{
	"nnrf-nfm":                     true, //Nnrf_NFManagement Service offered by the NRF
	"nnrf-disc":                    true, //Nnrf_NFDiscovery Service offered by the NRF
	"nnrf-oauth2":                  true, //Nnrf_AccessToken Service offered by the NRF
	"nudm-sdm":                     true, //Nudm_SubscriberDataManagement Service offered by the UDM
	"nudm-uecm":                    true, //Nudm_UEContextManagement Service offered by the UDM
	"nudm-ueau":                    true, //Nudm_UEAuthentication Service offered by the UDM
	"nudm-ee":                      true, //Nudm_EventExposure Service offered by the UDM
	"nudm-pp":                      true, //Nudm_ParameterProvision Service offered by the UDM
	"nudm-niddau":                  true, //Nudm_NIDDAuthorization Service offered by the UDM
	"nudm-mt":                      true, //Nudm_MT Service offered by the UDM
	"namf-comm":                    true, //Namf_Communication Service offered by the AMF
	"namf-evts":                    true, //Namf_EventExposure Service offered by the AMF
	"namf-mt":                      true, //Namf_MT Service offered by the AMF
	"namf-loc":                     true, //amf_Location Service offered by the AMF
	"nsmf-pdusession":              true, //Nsmf_PDUSession Service offered by the SMF
	"nsmf-event-exposure":          true, //Nsmf_EventExposure Service offered by the SMF
	"nsmf-nidd":                    true, //Nsmf_NIDD Service offered by the SMF
	"nausf-auth":                   true, //Nausf_UEAuthentication Service offered by the AUSF
	"nausf-sorprotection":          true, //Nausf_SoRProtection Service offered by the AUSF
	"nausf-upuprotection":          true, //Nausf_UPUProtection Service offered by the AUSF
	"nnef-pfdmanagement":           true, //Nnef_PFDManagement offered by the NEF
	"nnef-smcontext":               true, //Nnef_SMContext Service offered by the NEF
	"nnef-eventexposure":           true, //Nnef_EventExposure Service offered by the NEF
	"npcf-am-policy-control":       true, //Npcf_AMPolicyControl Service offered by the PCF
	"npcf-smpolicycontrol":         true, //Npcf_SMPolicyControl Service offered by the PCF
	"npcf-policyauthorization":     true, //Npcf_PolicyAuthorization Service offered by the PCF
	"npcf-bdtpolicycontrol":        true, //Npcf_BDTPolicyControl Service offered by the PCF
	"npcf-eventexposure":           true, //Npcf_EventExposure Service offered by the PCF
	"npcf-ue-policy-control":       true, //Npcf_UEPolicyControl Service offered by the PCF
	"nsmsf-sms":                    true, //Nsmsf_SMService Service offered by the SMSF
	"nnssf-nsselection":            true, //Nnssf_NSSelection Service offered by the NSSF
	"nnssf-nssaiavailability":      true, //Nnssf_NSSAIAvailability Service offered by the NSSF
	"nudr-dr":                      true, //Nudr_DataRepository Service offered by the UDR
	"nudr-group-id-map":            true, //Nudr_GroupIDmap Service offered by the UDR
	"nlmf-loc":                     true, //Nlmf_Location Service offered by the LMF
	"n5g-eir-eica":                 true, //N5g-eir_EquipmentIdentityCheck Service offered by the 5G-EIR
	"nbsf-management":              true, //Nbsf_Management Service offered by the BSF
	"nchf-spendinglimitcontrol":    true, //Nchf_SpendingLimitControl Service offered by the CHF
	"nchf-convergedcharging":       true, //Nchf_Converged_Charging Service offered by the CHF
	"nchf-offlineonlycharging":     true, //Nchf_OfflineOnlyCharging Service offered by the CHF
	"nnwdaf-eventssubscription":    true, //Nnwdaf_EventsSubscription Service offered by the NWDAF
	"nnwdaf-analyticsinfo":         true, //Nnwdaf_AnalyticsInfo Service offered by the NWDAF
	"ngmlc-loc":                    true, //Ngmlc_Location Service offered by GMLC
	"nucmf-provisioning":           true, //Nucmf_Provisioning Service offered by UCMF
	"nucmf-uecapabilitymanagement": true, //Nucmf_UECapabilityManagement Service offered by UCMF
	"nhss-sdm":                     true, //Nhss_SubscriberDataManagement Service offered by the HSS
	"nhss-uecm":                    true, //Nhss_UEContextManagement Service offered by the HSS
	"nhss-ueau":                    true, //Nhss_UEAuthentication Service offered by the HSS
	"nhss-ee":                      true, //Nhss_EventExposure Service offered by the HSS
	"nhss-ims-sdm":                 true, //Nhss_imsSubscriberDataManagement Service offered by the HSS
	"nhss-ims-uecm":                true, //Nhss_imsUEContextManagement Service offered by the HSS
	"nhss-ims-ueau":                true, //Nhss_imsUEAuthentication Service offered by the HSS
	"nsepp-telescopic":             true, //Nsepp_Telescopic_FQDN_Mapping Service offered by the SEPP
	"nsoraf-sor":                   true, //Nsoraf_SteeringOfRoaming Service offered by the SOR-AF
	"nspaf-secured-packed":         true, //Nspaf_SecuredPacket Service offered by the SP-AF
	"nudsf-dr":                     true, //Nudsf Data Repository service offered by the UDSF.
	"nnssaaf-nssaa":                true, //Nnssaaf_NSSAA service offered by the NSSAAF.

}

func IsValidNFType(nfType string) bool {
	_, ok := ALL_NFTYPES[nfType]
	return ok
}
func IsValidServiceName(serviceName string) bool {
	_, ok := ALL_SERVICE_NAMES[serviceName]
	return ok
}

type PlmnId struct {
	// 3 digital
	Mcc string `json:"mcc" form:"mcc"`
	// 2 or 2 digital
	Mnc string `json:"mnc" form:"mnc"`
}

type Snssai struct {
	// [0,255]
	Sst int32  `json:"sst" form:"sst"`
	Sd  string `json:"sd,omitempty"`
}

type PlmnIdNid struct {
	Mcc string `json:"mcc" form:"mcc"`
	Mnc string `json:"mnc" form:"mnc"`
	// pattern: ^[A-Fa-f0-9]{11}$
	Nid string `json:"nid,omitempty"`
}

type AccessTokenRequest struct {
	GrantType string `json:"grant_type" form:"grant_type"`
	// in uuid format
	NfInstanceId string `json:"nfInstanceId" form:"nfInstanceId"`
	// NFType in TS29501_Nnrf_NFManagement.yaml clause 6.1.6.3.3
	NfType       string `json:"nfType,omitempty" form:"nfType,omitempty"`
	TargetNfType string `json:"targetNfType,omitempty" form:"targetNfType,omitempty"`
	// in uuid format
	TargetNfInstanceId string `json:"targetNfInstanceId,emitempty" form:"targetNfInstanceId,omitempty"`
	// defined in TS 29.510 6.1.6.3.11
	Scope                string       `json:"scope" form:"scope,required"`
	RequesterPlmn        *PlmnId      `json:"targetNfInstanceId,omitempty" form:"targetNfInstanceId,omitempty"`
	RequesterPlmnList    []*PlmnId    `json:"requesterPlmnList,omitempty" form:"requesterPlmnList,omitempty"`
	RequesterSnssaiList  []*Snssai    `json:"requesterSnssaiList,omitempty" form:"requesterSnssaiList,omitempty"`
	RequesterFqdn        string       `json:"requesterFqdn,omitempty" form:"requesterFqdn,omitempty"`
	RequesterSnpnList    []*PlmnIdNid `json:"requesterSnpnList,omitempty" form:"requesterSnpnList,omitempty"`
	TargetPlmn           *PlmnId      `json:"targetPlmn,omitempty" form:"targetPlmn,omitempty"`
	TargetSnssaiList     []*Snssai    `json:"targetSnssaiList,omitempty" form:"targetSnssaiList,omitempty"`
	TargetNsiList        []string     `json:"targetNsiList,omitempty" form:"targetNsiList,omitempty"`
	TargetNfSetId        string       `json:"targetNfSetId,omitempty" form:"targetNfSetId,omitempty"`
	TargetNfServiceSetId string       `json:"targetNfServiceSetId,omitempty" form:"targetNfServiceSetId,omitempty"`
}

func NewAccessTokenRequest() *AccessTokenRequest {
	return &AccessTokenRequest{}
}

func (atr *AccessTokenRequest) ToJson() ([]byte, error) {
	return json.Marshal(atr)
}

func (atr *AccessTokenRequest) ToX3WFormEncoding() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0))
	encoder := form.NewEncoder(w)
	err := encoder.Encode(atr)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (atr *AccessTokenRequest) FromX3WFormEncoding(r io.Reader) error {
	decoder := form.NewDecoder(r)
	return decoder.Decode(atr)
}

func (atr *AccessTokenRequest) FromJson(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(atr)
}

func (atr *AccessTokenRequest) isValid() bool {
	if atr.GrantType != "client_credentials" {
		log.Error("the grant_type ", atr.GrantType, " is not client_credentials")
		return false
	}
	if len(atr.NfInstanceId) <= 0 {
		log.Error("Missing nfInstanceId")
		return false
	}
	if !IsValidServiceName(atr.Scope) {
		log.Error("Not valid scope ", atr.Scope)
		return false
	}
	return true

}

func (atr *AccessTokenRequest) isRequestByType() bool {
	if len(atr.NfType) > 0 && len(atr.TargetNfType) > 0 {
		if !IsValidNFType(atr.NfType) {
			log.Error("Invalid nfType ", atr.NfType)
			return false
		}
		if !IsValidNFType(atr.TargetNfType) {
			log.Error("Invalid targetNfType", atr.TargetNfType)
			return false
		}
		return true
	}
	return false
}
