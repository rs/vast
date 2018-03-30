package vast

import "encoding/xml"

// Extension represent arbitrary XML provided by the platform to extend the
// VAST response or by custom trackers.
type Extension struct {
	Type           string     `xml:"type,attr,omitempty"`
	CustomTracking []Tracking `xml:"CustomTracking>Tracking,omitempty"`
	Data           []byte     `xml:",innerxml"`
}

// the extension type as a middleware in the encoding process.
type extension Extension

type extensionNoCT struct {
	Type string `xml:"type,attr,omitempty"`
	Data []byte `xml:",innerxml"`
}

// MarshalXML implements xml.Marshaler interface.
func (e Extension) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	// create a temporary element from a wrapper Extension, copy what we need to
	// it and return it's encoding.
	var e2 interface{}
	// if we have custom trackers, we should ignore the data, if not, then we
	// should consider only the data.
	if len(e.CustomTracking) > 0 {
		e2 = extension{Type: e.Type, CustomTracking: e.CustomTracking}
	} else {
		e2 = extensionNoCT{Type: e.Type, Data: e.Data}
	}

	return enc.EncodeElement(e2, start)
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (e *Extension) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	// decode the extension into a temporary element from a wrapper Extension,
	// copy what we need over.
	var e2 extension
	if err := dec.DecodeElement(&e2, &start); err != nil {
		return err
	}
	// copy the type and the customTracking
	e.Type = e2.Type
	e.CustomTracking = e2.CustomTracking
	// copy the data only of customTracking is empty
	if len(e.CustomTracking) == 0 {
		e.Data = e2.Data
	}
	return nil
}
