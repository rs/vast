package vast

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	verificationExtensionXML                   = []byte(`<AdVerifications><Verification vendor="company.com-omid"><JavaScriptResource apiFramework="omid" browserOptional="true"><![CDATA[https://verification.com/omid_verification.js]]></JavaScriptResource><TrackingEvents><Tracking event="verificationNotExecuted"><![CDATA[https://verification.com/trackingurl]]></Tracking></TrackingEvents><VerificationParameters><![CDATA[{"foo":"bar"}]]></VerificationParameters></Verification></AdVerifications>`)
	verificationExtensionXMLExecutableResource = []byte(`<AdVerifications><Verification vendor="company.com-omid"><ExecutableResource apiFramework="omid"><![CDATA[https://verification.com/omid_verification.bin]]></ExecutableResource><TrackingEvents><Tracking event="verificationNotExecuted"><![CDATA[https://verification.com/trackingurl]]></Tracking></TrackingEvents><VerificationParameters><![CDATA[{"foo":"bar"}]]></VerificationParameters></Verification></AdVerifications>`)
	twoVerificationExtensionsXML               = []byte(`<AdVerifications><Verification vendor="company.com-omid"><JavaScriptResource apiFramework="omid" browserOptional="true"><![CDATA[https://verification.com/omid_verification.js]]></JavaScriptResource><TrackingEvents><Tracking event="verificationNotExecuted"><![CDATA[https://verification.com/trackingurl]]></Tracking></TrackingEvents><VerificationParameters><![CDATA[{"foo":"bar"}]]></VerificationParameters></Verification><Verification vendor="company1.com-omid"><JavaScriptResource apiFramework="omid" browserOptional="true"><![CDATA[https://verification1.com/omid_verification.js]]></JavaScriptResource><TrackingEvents><Tracking event="verificationNotExecuted"><![CDATA[https://verification1.com/trackingurl]]></Tracking></TrackingEvents><VerificationParameters><![CDATA[{"bar":"foo"}]]></VerificationParameters></Verification></AdVerifications>`)
	verificationExtensionXMLNoTracking         = []byte(`<AdVerifications><Verification vendor="company.com-omid"><JavaScriptResource apiFramework="omid" browserOptional="true"><![CDATA[https://verification.com/omid_verification.js]]></JavaScriptResource><VerificationParameters><![CDATA[{"foo":"bar"}]]></VerificationParameters></Verification></AdVerifications>`)
	verificationExtensionNoParams              = []byte(`<AdVerifications><Verification vendor="company.com-omid"><JavaScriptResource apiFramework="omid" browserOptional="true"><![CDATA[https://verification.com/omid_verification.js]]></JavaScriptResource><TrackingEvents><Tracking event="verificationNotExecuted"><![CDATA[https://verification.com/trackingurl]]></Tracking></TrackingEvents></Verification></AdVerifications>`)
)

func TestUnmarshal(t *testing.T) {
	adverifications := AdVerifications{}
	xml.Unmarshal(verificationExtensionXML, &adverifications)
	assert.Equal(t, "company.com-omid", adverifications.Verifications[0].Vendor)
	assert.Equal(t, "omid", adverifications.Verifications[0].JavaScriptResource.APIFramework)
	assert.Equal(t, true, adverifications.Verifications[0].JavaScriptResource.BrowserOptional)
	assert.Equal(t, "https://verification.com/omid_verification.js", adverifications.Verifications[0].JavaScriptResource.URI)
	assert.NotNil(t, *(adverifications.Verifications[0].TrackingEvents))
	assert.Len(t, *(adverifications.Verifications[0].TrackingEvents), 1)
	assert.Equal(t, "verificationNotExecuted", (*adverifications.Verifications[0].TrackingEvents)[0].Event)
	assert.Equal(t, "https://verification.com/trackingurl", (*adverifications.Verifications[0].TrackingEvents)[0].URI)
	assert.Nil(t, (*adverifications.Verifications[0].TrackingEvents)[0].Offset)
	assert.NotNil(t, adverifications.Verifications[0].VerificationParameters)
	assert.Equal(t, "{\"foo\":\"bar\"}", adverifications.Verifications[0].VerificationParameters.Params)
}

func TestUnmarshalExecutableResource(t *testing.T) {
	adverifications := AdVerifications{}
	xml.Unmarshal(verificationExtensionXMLExecutableResource, &adverifications)
	assert.Equal(t, "company.com-omid", adverifications.Verifications[0].Vendor)
	assert.Equal(t, "omid", adverifications.Verifications[0].ExecutableResource.APIFramework)
	assert.Equal(t, "https://verification.com/omid_verification.bin", adverifications.Verifications[0].ExecutableResource.URI)
	assert.NotNil(t, *(adverifications.Verifications[0].TrackingEvents))
	assert.Len(t, *(adverifications.Verifications[0].TrackingEvents), 1)
	assert.Equal(t, "verificationNotExecuted", (*adverifications.Verifications[0].TrackingEvents)[0].Event)
	assert.Equal(t, "https://verification.com/trackingurl", (*adverifications.Verifications[0].TrackingEvents)[0].URI)
	assert.Nil(t, (*adverifications.Verifications[0].TrackingEvents)[0].Offset)
	assert.NotNil(t, adverifications.Verifications[0].VerificationParameters)
	assert.Equal(t, "{\"foo\":\"bar\"}", adverifications.Verifications[0].VerificationParameters.Params)
}

func TestNoTrackingEventsUnmarshal(t *testing.T) {
	adverifications := AdVerifications{}
	xml.Unmarshal(verificationExtensionXMLNoTracking, &adverifications)
	assert.Nil(t, adverifications.Verifications[0].TrackingEvents)
}

func TestNoParamsUnmarshal(t *testing.T) {
	adverifications := AdVerifications{}
	xml.Unmarshal(verificationExtensionNoParams, &adverifications)
	assert.Nil(t, adverifications.Verifications[0].VerificationParameters)
}

func TestTwoVerificationsUnmarshal(t *testing.T) {
	adverifications := AdVerifications{}
	xml.Unmarshal(twoVerificationExtensionsXML, &adverifications)
	assert.Equal(t, "company.com-omid", adverifications.Verifications[0].Vendor)
	assert.Equal(t, "company1.com-omid", adverifications.Verifications[1].Vendor)
	assert.Equal(t, "https://verification.com/omid_verification.js", adverifications.Verifications[0].JavaScriptResource.URI)
	assert.Equal(t, "https://verification1.com/omid_verification.js", adverifications.Verifications[1].JavaScriptResource.URI)
	assert.Equal(t, "https://verification.com/trackingurl", (*adverifications.Verifications[0].TrackingEvents)[0].URI)
	assert.Equal(t, "https://verification1.com/trackingurl", (*adverifications.Verifications[1].TrackingEvents)[0].URI)
	assert.Equal(t, "{\"foo\":\"bar\"}", adverifications.Verifications[0].VerificationParameters.Params)
	assert.Equal(t, "{\"bar\":\"foo\"}", adverifications.Verifications[1].VerificationParameters.Params)
}

var (
	verificationExtension = AdVerifications{
		Verifications: []Verification{
			{
				Vendor: "company.com",
				JavaScriptResource: &JavaScriptResource{
					APIFramework:    "omid",
					BrowserOptional: true,
					URI:             "https://dummy.com",
				},
				TrackingEvents: &[]Tracking{
					{
						Event: "verificationNotExecuted",
						URI:   "https://dummy.com/track",
					},
				},
				VerificationParameters: &VerificationParameters{
					Params: "{\"foo\":\"bar\"}",
				},
			},
		},
	}
	verificationExtensionExecutableResource = AdVerifications{
		Verifications: []Verification{
			{
				Vendor: "company.com",
				ExecutableResource: &ExecutableResource{
					APIFramework: "omid",
					URI:          "https://dummy.bin",
				},
				TrackingEvents: &[]Tracking{
					{
						Event: "verificationNotExecuted",
						URI:   "https://dummy.com/track",
					},
				},
				VerificationParameters: &VerificationParameters{
					Params: "{\"foo\":\"bar\"}",
				},
			},
		},
	}
	verificationExtensionNoTrackingOrParams = AdVerifications{
		Verifications: []Verification{
			{
				Vendor: "company.com",
				JavaScriptResource: &JavaScriptResource{
					APIFramework:    "omid",
					BrowserOptional: true,
					URI:             "https://dummy.com",
				},
			},
		},
	}
)

func TestMarshal(t *testing.T) {
	adVerificationsXML, err := xml.Marshal(verificationExtension)
	assert.Nil(t, err)
	assert.NotNil(t, adVerificationsXML)
	assert.Equal(t, []byte(`<AdVerifications><Verification vendor="company.com"><JavaScriptResource apiFramework="omid" browserOptional="true"><![CDATA[https://dummy.com]]></JavaScriptResource><TrackingEvents><Tracking event="verificationNotExecuted"><![CDATA[https://dummy.com/track]]></Tracking></TrackingEvents><VerificationParameters><![CDATA[{"foo":"bar"}]]></VerificationParameters></Verification></AdVerifications>`), adVerificationsXML)
}

func TestExecutableResourceMarshal(t *testing.T) {
	adVerificationsXML, err := xml.Marshal(verificationExtensionExecutableResource)
	assert.Nil(t, err)
	assert.NotNil(t, adVerificationsXML)
	assert.Equal(t, string([]byte(`<AdVerifications><Verification vendor="company.com"><ExecutableResource apiFramework="omid"><![CDATA[https://dummy.bin]]></ExecutableResource><TrackingEvents><Tracking event="verificationNotExecuted"><![CDATA[https://dummy.com/track]]></Tracking></TrackingEvents><VerificationParameters><![CDATA[{"foo":"bar"}]]></VerificationParameters></Verification></AdVerifications>`)), string(adVerificationsXML))
}

func TestNoTrackingOrParamsMarshal(t *testing.T) {
	adVerificationsXML, err := xml.Marshal(verificationExtensionNoTrackingOrParams)
	assert.Nil(t, err)
	assert.NotNil(t, adVerificationsXML)
	assert.Equal(t, []byte(`<AdVerifications><Verification vendor="company.com"><JavaScriptResource apiFramework="omid" browserOptional="true"><![CDATA[https://dummy.com]]></JavaScriptResource></Verification></AdVerifications>`), adVerificationsXML)
}
