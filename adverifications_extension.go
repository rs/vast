package vast

type AdVerifications struct {
	Verifications []Verification `xml:"Verification"`
}

type Verification struct {
	Vendor                 string                  `xml:"vendor,attr,omitempty"`
	JavaScriptResource     JavaScriptResource      `xml:"JavaScriptResource"`
	TrackingEvents         *[]Tracking             `xml:"TrackingEvents>Tracking,omitempty"`
	VerificationParameters *VerificationParameters `xml:"VerificationParameters"`
}

type JavaScriptResource struct {
	APIFramework    string `xml:"apiFramework,attr,omitempty"`
	BrowserOptional bool   `xml:"browserOptional,attr,omitempty"`
	URI             string `xml:",cdata"`
}

type VerificationParameters struct {
	Params string `xml:",cdata"`
}
