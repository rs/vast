// Package vast implements IAB VAST 3.0 specification http://www.iab.net/media/file/VASTv3.0.pdf
package vast

import (
	"bytes"
	"encoding/xml"
	"github.com/pkg/errors"
	"strings"
)

const (
	WrapperCreativeType = "w"
	VastCreativeType    = "v"
	VpaidCreativeType   = "vp"
)

// MarshalXML is a custom XML marshalling method, with some fixes on top of the native encoding/xml package
func (v *VAST) MarshalXML() ([]byte, error) {
	v.formatVastStrings()
	data, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	strXML := string(data)
	strXML = strings.Replace(strXML, "_xmlns", "xmlns", -1)
	strXML = strings.Replace(strXML, ` xmlns:xmlns="xmlns"`, "", -1)
	strXML = xml.Header + strXML
	strXML = strings.Replace(strXML, "\n", "", -1)
	strXML = strings.Replace(strXML, "\t", "", -1)
	strXML = strings.TrimSpace(strXML)
	return []byte(strXML), nil
}

func (v *VAST) GetMediaFile() (*MediaFile, string) {
	if len(v.Ads) == 0 {
		return nil, ""
	}
	for _, ad := range v.Ads {
		if ad.InLine == nil {
			return nil, ""
		}
		inline := ad.InLine
		for _, creative := range inline.Creatives {
			if creative.Linear != nil {
				linear := creative.Linear
				for i := range linear.MediaFiles {
					return &linear.MediaFiles[i], linear.Duration
				}
			}
		}
	}
	return nil, ""
}

func (v *VAST) GetClickTroughURL() string {
	if len(v.Ads) == 0 {
		return ""
	}
	for _, ad := range v.Ads {
		if ad.InLine == nil {
			continue
		}
		inline := ad.InLine
		for _, creative := range inline.Creatives {
			linear := creative.Linear
			if linear == nil {
				continue
			}
			vc := linear.VideoClicks
			if vc == nil {
				continue
			}
			clickThrough := vc.ClickThroughs
			for _, ct := range clickThrough {
				if uri := ct.URI; uri != "" {
					return uri
				}
			}
		}
	}
	return ""
}

func (v *VAST) RemovePixels() {
	v.Errors = nil
	for _, ad := range v.Ads {
		if inline := ad.InLine; inline != nil {
			inline.Impressions = nil
			inline.Error = nil
			for i := range inline.Creatives {
				cr := &inline.Creatives[i]
				if linear := cr.Linear; linear != nil {
					linear.TrackingEvents = nil
				}
				if ca := cr.CompanionAds; ca != nil {
					for c := range ca.Companions {
						ca.Companions[c].TrackingEvents = nil
					}
				}
			}
		}
	}
}

func (v *VAST) GetVastType() (string, error) {
	ads := v.Ads
	if len(ads) == 0 {
		return "", errors.New("empty Ads attribute")
	}
	ad := ads[0]

	if wrapper := ad.Wrapper; wrapper != nil && wrapper.VASTAdTagURI != "" {
		return WrapperCreativeType, nil
	}
	inLine := ad.InLine
	if inLine == nil {
		return "", errors.New("empty InLine attribute")
	}
	creatives := inLine.Creatives
	if len(creatives) == 0 {
		return "", errors.New("empty Creative attribute")
	}
	linear := creatives[0].Linear
	if linear == nil {
		return "", errors.New("empty Linear attribute")
	}
	mediafiles := linear.MediaFiles
	if len(mediafiles) == 0 {
		return "", errors.New("empty MediaFiles attribute")
	}
	if strings.EqualFold(mediafiles[0].APIFramework, "VPAID") {
		return VpaidCreativeType, nil
	}
	return VastCreativeType, nil
}

// FromXML is a custom XML unmarshalling method, with some fixes on top of the native encoding/xml package
func FromXML(xmlStr []byte) (*VAST, error) {
	xmlStr = bytes.Replace(xmlStr, []byte("\n"), []byte(""), -1)
	xmlStr = bytes.Replace(xmlStr, []byte("\t"), []byte(""), -1)
	var v VAST
	if err := xml.Unmarshal(xmlStr, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// VAST is the root <VAST> tag
type VAST struct {
	// The version of the VAST spec (should be either "2.0" or "3.0")
	Version string `xml:"version,attr"`
	// One or more Ad elements. Advertisers and video content publishers may
	// associate an <Ad> element with a line item video ad defined in contract
	// documentation, usually an insertion order. These line item ads typically
	// specify the creative to display, price, delivery schedule, targeting,
	// and so on.
	Ads []Ad `xml:"Ad"`
	// Contains a URI to a tracking resource that the video player should request
	// upon receiving a “no ad” response
	Errors []string `xml:"Error"`
	// Defaults for media file's width/height attributes in case they're "0" or empty
	mediaFileDefaultWidth  int
	mediaFileDefaultHeight int
}

// Ad represent an <Ad> child tag in a VAST document
//
// Each <Ad> contains a single <InLine> element or <Wrapper> element (but never both).
type Ad struct {
	// An ad server-defined identifier string for the ad
	ID string `xml:"id,attr,omitempty"`
	// A number greater than zero (0) that identifies the sequence in which
	// an ad should play; all <Ad> elements with sequence values are part of
	// a pod and are intended to be played in sequence
	Sequence int      `xml:"sequence,attr,omitempty"`
	InLine   *InLine  `xml:",omitempty"`
	Wrapper  *Wrapper `xml:",omitempty"`
}

// InLine is a vast <InLine> ad element containing actual ad definition
//
// The last ad server in the ad supply chain serves an <InLine> element.
// Within the nested elements of an <InLine> element are all the files and
// URIs necessary to display the ad.
type InLine struct {
	// The name of the ad server that returned the ad
	AdSystem *AdSystem
	// The common name of the ad
	AdTitle *AdTitle

	// One or more URIs that directs the video player to a tracking resource file that the
	// video player should request when the first frame of the ad is displayed
	Impressions []Impression `xml:"Impression"`
	// The container for one or more <Creative> elements
	Creatives []Creative `xml:"Creatives>Creative"`
	// A string value that provides a longer description of the ad.
	Description string `xml:",omitempty"`
	// The name of the advertiser as defined by the ad serving party.
	// This element can be used to prevent displaying ads with advertiser
	// competitors. Ad serving parties and publishers should identify how
	// to interpret values provided within this element. As with any optional
	// elements, the video player is not required to support it.
	Advertiser string `xml:",omitempty"`
	// A URI to a survey vendor that could be the survey, a tracking pixel,
	// or anything to do with the survey. Multiple survey elements can be provided.
	// A type attribute is available to specify the MIME type being served.
	// For example, the attribute might be set to type=”text/javascript”.
	// Surveys can be dynamically inserted into the VAST response as long as
	// cross-domain issues are avoided.
	Survey string `xml:",omitempty"`
	// A URI representing an error-tracking pixel; this element can occur multiple
	// times.
	Error []Error `xml:",omitempty"`
	// Provides a value that represents a price that can be used by real-time bidding
	// (RTB) systems. VAST is not designed to handle RTB since other methods exist,
	// but this element is offered for custom solutions if needed.
	Pricing string `xml:",omitempty"`
	// XML node for custom extensions, as defined by the ad server. When used, a
	// custom element should be nested under <Extensions> to help separate custom
	// XML elements from VAST elements. The following example includes a custom
	// xml element within the Extensions element.
	Extensions      *Extensions      `xml:",omitempty"`
	AdVerifications *AdVerifications `xml:",omitempty"`
}

type Error struct {
	URI string `xml:",cdata"`
}

// Impression is a URI that directs the video player to a tracking resource file that
// the video player should request when the first frame of the ad is displayed
type Impression struct {
	ID  string `xml:"id,attr,omitempty"`
	URI string `xml:",cdata"`
}

// Pricing provides a value that represents a price that can be used by real-time
// bidding (RTB) systems. VAST is not designed to handle RTB since other methods
// exist,  but this element is offered for custom solutions if needed.
type Pricing struct {
	// Identifies the pricing model as one of "cpm", "cpc", "cpe" or "cpv".
	Model string `xml:"model,attr"`
	// The 3 letter ISO-4217 currency symbol that identifies the currency of
	// the value provided
	Currency string `xml:"currency,attr"`
	// If the value provided is to be obfuscated/encoded, publishers and advertisers
	// must negotiate the appropriate mechanism to do so. When included as part of
	// a VAST Wrapper in a chain of Wrappers, only the value offered in the first
	// Wrapper need be considered.
	Value string `xml:",chardata"`
}

// Wrapper element contains a URI reference to a vendor ad server (often called
// a third party ad server). The destination ad server either provides the ad
// files within a VAST <InLine> ad element or may provide a secondary Wrapper
// ad, pointing to yet another ad server. Eventually, the final ad server in
// the ad supply chain must contain all the necessary files needed to display
// the ad.
type Wrapper struct {
	// The name of the ad server that returned the ad
	AdSystem *AdSystem
	// URL of ad tag of downstream Secondary Ad Server
	VASTAdTagURI string
	// One or more URIs that directs the video player to a tracking resource file that the
	// video player should request when the first frame of the ad is displayed
	Impressions []Impression `xml:"Impression"`
	// A URI representing an error-tracking pixel; this element can occur multiple
	// times.
	Error []Error `xml:",omitempty"`
	// The container for one or more <Creative> elements
	Creatives []CreativeWrapper `xml:"Creatives>Creative"`
	// XML node for custom extensions, as defined by the ad server. When used, a
	// custom element should be nested under <Extensions> to help separate custom
	// XML elements from VAST elements. The following example includes a custom
	// xml element within the Extensions element.
	Extensions      *Extensions      `xml:",omitempty"`
	AdVerifications *AdVerifications `xml:",omitempty"`
}

// AdSystem contains information about the system that returned the ad
type AdSystem struct {
	Version string `xml:"version,attr,omitempty"`
	Name    string `xml:",cdata"`
}

type AdTitle struct {
	Name string `xml:",cdata"`
}

// Creative is a file that is part of a VAST ad.
type Creative struct {
	// An ad server-defined identifier for the creative
	ID string `xml:"id,attr,omitempty"`
	// The preferred order in which multiple Creatives should be displayed
	Sequence int `xml:"sequence,attr,omitempty"`
	// Identifies the ad with which the creative is served
	AdID string `xml:"AdID,attr,omitempty"`
	// The technology used for any included API
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// If present, defines a linear creative
	Linear *Linear `xml:",omitempty"`
	// If defined, defins companions creatives
	CompanionAds *CompanionAds `xml:",omitempty"`
	// If defined, defins non linear creatives
	NonLinearAds *NonLinearAds `xml:",omitempty"`
}

// CompanionAds contains companions creatives
type CompanionAds struct {
	// Provides information about which companion creative to display.
	// All means that the player must attempt to display all. Any means the player
	// must attempt to play at least one. None means all companions are optional
	Required   string      `xml:"required,attr,omitempty"`
	Companions []Companion `xml:"Companion,omitempty"`
}

// NonLinearAds contains non linear creatives
type NonLinearAds struct {
	TrackingEvents []Tracking `xml:"TrackingEvents>Tracking,omitempty"`
	// Non linear creatives
	NonLinears []NonLinear `xml:"NonLinear,omitempty"`
}

// CreativeWrapper defines wrapped creative's parent trackers
type CreativeWrapper struct {
	// An ad server-defined identifier for the creative
	ID string `xml:"id,attr,omitempty"`
	// The preferred order in which multiple Creatives should be displayed
	Sequence int `xml:"sequence,attr,omitempty"`
	// Identifies the ad with which the creative is served
	AdID string `xml:"AdID,attr,omitempty"`
	// If present, defines a linear creative
	Linear *LinearWrapper `xml:",omitempty"`
	// If defined, defins companions creatives
	CompanionAds *CompanionAdsWrapper `xml:"CompanionAds,omitempty"`
	// If defined, defines non linear creatives
	NonLinearAds *NonLinearAdsWrapper `xml:"NonLinearAds,omitempty"`
}

// CompanionAdsWrapper contains companions creatives in a wrapper
type CompanionAdsWrapper struct {
	// Provides information about which companion creative to display.
	// All means that the player must attempt to display all. Any means the player
	// must attempt to play at least one. None means all companions are optional
	Required   string             `xml:"required,attr,omitempty"`
	Companions []CompanionWrapper `xml:"Companion,omitempty"`
}

// NonLinearAdsWrapper contains non linear creatives in a wrapper
type NonLinearAdsWrapper struct {
	TrackingEvents []Tracking `xml:"TrackingEvents>Tracking,omitempty"`
	// Non linear creatives
	NonLinears []NonLinearWrapper `xml:"NonLinear,omitempty"`
}

// Linear is the most common type of video advertisement trafficked in the
// industry is a “linear ad”, which is an ad that displays in the same area
// as the content but not at the same time as the content. In fact, the video
// player must interrupt the content before displaying a linear ad.
// Linear ads are often displayed right before the video content plays.
// This ad position is called a “pre-roll” position. For this reason, a linear
// ad is often called a “pre-roll.”
type Linear struct {
	// To specify that a Linear creative can be skipped, the ad server must
	// include the skipoffset attribute in the <Linear> element. The value
	// for skipoffset is a time value in the format HH:MM:SS or HH:MM:SS.mmm
	// or a percentage in the format n%. The .mmm value in the time offset
	// represents milliseconds and is optional. This skipoffset value
	// indicates when the skip control should be provided after the creative
	// begins playing.
	SkipOffset *Offset `xml:"skipoffset,attr,omitempty"`
	// Duration in standard time format, hh:mm:ss
	Duration           string
	AdParameters       *AdParameters `xml:",omitempty"`
	Icons              []Icon
	TrackingEvents     []Tracking          `xml:"TrackingEvents>Tracking,omitempty"`
	VideoClicks        *VideoClicks        `xml:",omitempty"`
	MediaFiles         []MediaFile         `xml:"MediaFiles>MediaFile,omitempty"`
	CreativeExtensions *CreativeExtensions `xml:",omitempty"`
}

// LinearWrapper defines a wrapped linear creative
type LinearWrapper struct {
	Icons              []Icon
	TrackingEvents     []Tracking          `xml:"TrackingEvents>Tracking,omitempty"`
	VideoClicks        *VideoClicks        `xml:",omitempty"`
	CreativeExtensions *CreativeExtensions `xml:",omitempty"`
}

// Companion defines a companion ad
type Companion struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty"`
	// Pixel dimensions of companion slot.
	Width int `xml:"width,attr"`
	// Pixel dimensions of companion slot.
	Height int `xml:"height,attr"`
	// Pixel dimensions of the companion asset.
	AssetWidth int `xml:"assetWidth,attr"`
	// Pixel dimensions of the companion asset.
	AssetHeight int `xml:"assetHeight,attr"`
	// Pixel dimensions of expanding companion ad when in expanded state.
	ExpandedWidth int `xml:"expandedWidth,attr"`
	// Pixel dimensions of expanding companion ad when in expanded state.
	ExpandeHeight int `xml:"expandedHeight,attr"`
	// The apiFramework defines the method to use for communication with the companion.
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// Used to match companion creative to publisher placement areas on the page.
	AdSlotID string `xml:"adSlotId,attr,omitempty"`
	// URL to open as destination page when user clicks on the the companion banner ad.
	CompanionClickThrough  *CompanionClickThrough    `xml:",omitempty"`
	CompanionClickTracking []*CompanionClickTracking `xml:",omitempty"`
	// Alt text to be displayed when companion is rendered in HTML environment.
	AltText string `xml:",omitempty"`
	// The creativeView should always be requested when present. For Companions
	// creativeView is the only supported event.
	TrackingEvents []Tracking `xml:"TrackingEvents>Tracking,omitempty"`
	// Data to be passed into the companion ads. The apiFramework defines the method
	// to use for communication (e.g. “FlashVar”)
	AdParameters *AdParameters `xml:",omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource string `xml:",omitempty"`
	// HTML to display the companion element
	HTMLResource *HTMLResource `xml:",omitempty"`
	// Extensions
	CreativeExtensions *CreativeExtensions `xml:",omitempty"`
}

// CompanionWrapper defines a companion ad in a wrapper
type CompanionWrapper struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty"`
	// Pixel dimensions of companion slot.
	Width int `xml:"width,attr"`
	// Pixel dimensions of companion slot.
	Height int `xml:"height,attr"`
	// Pixel dimensions of the companion asset.
	AssetWidth int `xml:"assetWidth,attr"`
	// Pixel dimensions of the companion asset.
	AssetHeight int `xml:"assetHeight,attr"`
	// Pixel dimensions of expanding companion ad when in expanded state.
	ExpandedWidth int `xml:"expandedWidth,attr"`
	// Pixel dimensions of expanding companion ad when in expanded state.
	ExpandeHeight int `xml:"expandedHeight,attr"`
	// The apiFramework defines the method to use for communication with the companion.
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// Used to match companion creative to publisher placement areas on the page.
	AdSlotID string `xml:"adSlotId,attr,omitempty"`
	// URL to open as destination page when user clicks on the the companion banner ad.
	CompanionClickThrough string `xml:",omitempty"`
	// URLs to ping when user clicks on the the companion banner ad.
	CompanionClickTracking []string `xml:",omitempty"`
	// Alt text to be displayed when companion is rendered in HTML environment.
	AltText string `xml:",omitempty"`
	// The creativeView should always be requested when present. For Companions
	// creativeView is the only supported event.
	TrackingEvents []Tracking `xml:"TrackingEvents>Tracking,omitempty"`
	// Data to be passed into the companion ads. The apiFramework defines the method
	// to use for communication (e.g. “FlashVar”)
	AdParameters *AdParameters `xml:",omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource string `xml:",omitempty"`
	// HTML to display the companion element
	HTMLResource       *HTMLResource       `xml:",omitempty"`
	CreativeExtensions *CreativeExtensions `xml:",omitempty"`
}

// NonLinear defines a non linear ad
type NonLinear struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty"`
	// Pixel dimensions of companion.
	Width int `xml:"width,attr"`
	// Pixel dimensions of companion.
	Height int `xml:"height,attr"`
	// Pixel dimensions of expanding nonlinear ad when in expanded state.
	ExpandedWidth int `xml:"expandedWidth,attr"`
	// Pixel dimensions of expanding nonlinear ad when in expanded state.
	ExpandeHeight int `xml:"expandedHeight,attr"`
	// Whether it is acceptable to scale the image.
	Scalable bool `xml:"scalable,attr,omitempty"`
	// Whether the ad must have its aspect ratio maintained when scales.
	MaintainAspectRatio bool `xml:"maintainAspectRatio,attr,omitempty"`
	// Suggested duration to display non-linear ad, typically for animation to complete.
	// Expressed in standard time format hh:mm:ss.
	MinSuggestedDuration *Duration `xml:"minSuggestedDuration,attr,omitempty"`
	// The apiFramework defines the method to use for communication with the nonlinear element.
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// URLs to ping when user clicks on the the non-linear ad.
	NonLinearClickTracking []string `xml:",omitempty"`
	// URL to open as destination page when user clicks on the non-linear ad unit.
	NonLinearClickThrough string `xml:",omitempty"`
	// Data to be passed into the video ad.
	AdParameters *AdParameters `xml:",omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource string `xml:",omitempty"`
	// HTML to display the companion element
	HTMLResource       *HTMLResource       `xml:",omitempty"`
	CreativeExtensions *CreativeExtensions `xml:",omitempty"`
}

// NonLinearWrapper defines a non linear ad in a wrapper
type NonLinearWrapper struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty"`
	// Pixel dimensions of companion.
	Width int `xml:"width,attr"`
	// Pixel dimensions of companion.
	Height int `xml:"height,attr"`
	// Pixel dimensions of expanding nonlinear ad when in expanded state.
	ExpandedWidth int `xml:"expandedWidth,attr"`
	// Pixel dimensions of expanding nonlinear ad when in expanded state.
	ExpandeHeight int `xml:"expandedHeight,attr"`
	// Whether it is acceptable to scale the image.
	Scalable bool `xml:"scalable,attr,omitempty"`
	// Whether the ad must have its aspect ratio maintained when scales.
	MaintainAspectRatio bool `xml:"maintainAspectRatio,attr,omitempty"`
	// Suggested duration to display non-linear ad, typically for animation to complete.
	// Expressed in standard time format hh:mm:ss.
	MinSuggestedDuration *Duration `xml:"minSuggestedDuration,attr,omitempty"`
	// The apiFramework defines the method to use for communication with the nonlinear element.
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// The creativeView should always be requested when present.
	TrackingEvents []Tracking `xml:"TrackingEvents>Tracking,omitempty"`
	// URLs to ping when user clicks on the the non-linear ad.
	NonLinearClickTracking []string            `xml:",omitempty"`
	CreativeExtensions     *CreativeExtensions `xml:",omitempty"`
}

// Icon represents advertising industry initiatives like AdChoices.
type Icon struct {
	// Identifies the industry initiative that the icon supports.
	Program string `xml:"program,attr"`
	// Pixel dimensions of icon.
	Width int `xml:"width,attr"`
	// Pixel dimensions of icon.
	Height int `xml:"height,attr"`
	// The horizontal alignment location (in pixels) or a specific alignment.
	// Must match ([0-9]*|left|right)
	XPosition string `xml:"xPosition,attr"`
	// The vertical alignment location (in pixels) or a specific alignment.
	// Must match ([0-9]*|top|bottom)
	YPosition string `xml:"yPosition,attr"`
	// Start time at which the player should display the icon. Expressed in standard time format hh:mm:ss.
	Offset Offset `xml:"offset,attr"`
	// duration for which the player must display the icon. Expressed in standard time format hh:mm:ss.
	Duration string `xml:"duration,attr"`
	// The apiFramework defines the method to use for communication with the icon element
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// URL to open as destination page when user clicks on the icon.
	IconClickThrough string `xml:"IconClicks>IconClickThrough,omitempty"`
	// URLs to ping when user clicks on the the icon.
	IconClickTrackings []string `xml:"IconClicks>IconClickTracking,omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource string `xml:",omitempty"`
	// HTML to display the companion element
	HTMLResource *HTMLResource `xml:",omitempty"`
}

// Tracking defines an event tracking URL
type Tracking struct {
	// The name of the event to track for the element. The creativeView should
	// always be requested when present.
	//
	// Possible values are creativeView, start, firstQuartile, midpoint, thirdQuartile,
	// complete, mute, unmute, pause, rewind, resume, fullscreen, exitFullscreen, expand,
	// collapse, acceptInvitation, close, skip, progress.
	Event string `xml:"event,attr"`
	// The time during the video at which this url should be pinged. Must be present for
	// progress event. Must match (\d{2}:[0-5]\d:[0-5]\d(\.\d\d\d)?|1?\d?\d(\.?\d)*%)
	Offset *Offset `xml:"offset,attr,omitempty"`
	URI    string  `xml:",cdata"`
}

// StaticResource is the URL to a static file, such as an image or SWF file
type StaticResource struct {
	// Mime type of static resource
	CreativeType string `xml:"creativeType,attr,omitempty"`
	// URL to a static file, such as an image or SWF file
	URI string `xml:",cdata"`
}

// HTMLResource is a container for HTML data
type HTMLResource struct {
	// Specifies whether the HTML is XML-encoded
	XMLEncoded bool   `xml:"xmlEncoded,attr,omitempty"`
	HTML       []byte `xml:",cdata"`
}

// AdParameters defines arbitrary ad parameters
type AdParameters struct {
	// Specifies whether the parameters are XML-encoded
	XMLEncoded bool   `xml:"xmlEncoded,attr,omitempty"`
	Parameters []byte `xml:",cdata"`
}

// VideoClicks contains types of video clicks
type VideoClicks struct {
	ClickThroughs  []VideoClick `xml:"ClickThrough,omitempty"`
	ClickTrackings []VideoClick `xml:"ClickTracking,omitempty"`
	CustomClicks   []VideoClick `xml:"CustomClick,omitempty"`
}

// VideoClick defines a click URL for a linear creative
type VideoClick struct {
	ID  string `xml:"id,attr,omitempty"`
	URI string `xml:",cdata"`
}

// MediaFile defines a reference to a linear creative asset
type MediaFile struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty"`
	// Method of delivery of ad (either "streaming" or "progressive")
	Delivery string `xml:"delivery,attr"`
	// MIME type. Popular MIME types include, but are not limited to
	// “video/x-ms-wmv” for Windows Media, and “video/x-flv” for Flash
	// Video. Image ads or interactive ads can be included in the
	// MediaFiles section with appropriate Mime types
	Type string `xml:"type,attr"`
	// The codec used to produce the media file.
	Codec string `xml:"codec,attr,omitempty"`
	// Bitrate of encoded video in Kbps. If bitrate is supplied, MinBitrate
	// and MaxBitrate should not be supplied.
	Bitrate int `xml:"bitrate,attr,omitempty"`
	// Minimum bitrate of an adaptive stream in Kbps. If MinBitrate is supplied,
	// MaxBitrate must be supplied and Bitrate should not be supplied.
	MinBitrate int `xml:"minBitrate,attr,omitempty"`
	// Maximum bitrate of an adaptive stream in Kbps. If MaxBitrate is supplied,
	// MinBitrate must be supplied and Bitrate should not be supplied.
	MaxBitrate int `xml:"maxBitrate,attr,omitempty"`
	// Pixel dimensions of video.
	Width int `xml:"width,attr"`
	// Pixel dimensions of video.
	Height int `xml:"height,attr"`
	// Whether it is acceptable to scale the image.
	Scalable bool `xml:"scalable,attr,omitempty"`
	// Whether the ad must have its aspect ratio maintained when scales.
	MaintainAspectRatio bool `xml:"maintainAspectRatio,attr,omitempty"`
	// The APIFramework defines the method to use for communication if the MediaFile
	// is interactive. Suggested values for this element are “VPAID”, “FlashVars”
	// (for Flash/Flex), “initParams” (for Silverlight) and “GetVariables” (variables
	// placed in key/value pairs on the asset request).
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	URI          string `xml:",cdata"`
}

// Extensions defines extensions
type Extensions struct {
	Extensions []Extension `xml:"Extension,omitempty"`
}

// CreativeExtensions defines extensions for creatives
type CreativeExtensions struct {
	Extensions []Extension `xml:"CreativeExtension,omitempty"`
}

// Extension represent aribtrary XML provided by the platform to extend the VAST response
type Extension struct {
	Type string `xml:"type,attr,omitempty"`
	Data []byte `xml:",innerxml"`
}

type CompanionClickThrough struct {
	// URL to a static file, such as an image or SWF file
	URI string `xml:",cdata"`
}

type CompanionClickTracking struct {
	// URL to a static file, such as an image or SWF file
	URI string `xml:",cdata"`
}

//The <AdVerifications> element contains one or more <Verification> elements, which list the resources and metadata
// required to execute third-party measurement code in order to verify creative playback.
type AdVerifications struct {
	AdVerifications []Verification `xml:"Verification,omitempty"`
}

// The <Verification> element contains the executable and bootstrapping data required to run the measurement code for
// a single verification vendor.
type Verification struct {
	Vendor                 string                  `xml:"vendor,attr"`
	JavaScriptResource     *JavaScriptResource     `xml:",omitempty"`
	ExecutableResource     *ExecutableResource     `xml:",omitempty"`
	TrackingEvents         []Tracking              `xml:"TrackingEvents>Tracking,omitempty"`
	VerificationParameters *VerificationParameters `xml:",omitempty"`
}

// JavaScriptResource is a container for the URI to the JavaScript file used to collect verification data.
type JavaScriptResource struct {
	// The name of the API framework used to execute the AdVerification code (required)
	APIFramework string `xml:"apiFramework,attr"`
	// If "true", this resource is optimized and able to execute in an environment without DOM and other browser
	// built-ins (required)
	BrowserOptional bool `xml:"browserOptional,attr"`
	// URL to the JavaScript used to collect data
	URI string `xml:",cdata"`
}

// ExecutableResource is a reference to a non-JavaScript or custom-integration resource intended for collecting
// verification data via the listed apiFramework.
type ExecutableResource struct {
	// The name of the API framework used to execute the AdVerification code (required)
	APIFramework string `xml:"apiFramework,attr"`
	// The type of executable resource provided
	Type string `xml:"type,attr,omitempty"`
	// A CDATA-wrapped reference to the resource
	URI string `xml:",cdata"`
}

// VerificationParameters contains a CDATA-wrapped string intended for bootstrapping the verification code and providing
// metadata about the current impression.
type VerificationParameters struct {
	// CDATA-wrapped metadata string for the verification executable
	URI string `xml:",cdata"`
}

func (v *VAST) GetDuration() (string, error) {
	ads := v.Ads
	if len(ads) == 0 {
		return "", errors.New("no Ads in Vast")
	}
	inline := ads[0].InLine
	if inline == nil {
		return "", errors.New("no Inline in Vast")
	}
	creatives := inline.Creatives
	if len(creatives) == 0 {
		return "", errors.New("no Creatives in Vast")
	}
	linear := creatives[0].Linear
	if linear == nil {
		return "", errors.New("no Linear object in first creative")
	}
	return linear.Duration, nil
}

func formatVastString(url string) string {
	return strings.TrimSpace(url)
}

func (v *VAST) SetMediaFileDefaultSize(mediaFileDefaultHeight, mediaFileDefaultWidth int) {
	v.mediaFileDefaultHeight = mediaFileDefaultHeight
	v.mediaFileDefaultWidth = mediaFileDefaultWidth
}

func (mf *MediaFile) fixMediaFileAttrs(mediaFileDefaultHeight, mediaFileDefaultWidth int) {
	if mf.Height == 0 || mf.Width == 0 {
		mf.Height = mediaFileDefaultHeight
		mf.Width = mediaFileDefaultWidth
	}
}

func (v *VAST) formatVastStrings() {
	if len(v.Ads) == 0 {
		return
	}

	for _, ad := range v.Ads {
		if ad.InLine == nil {
			return
		}
		inline := ad.InLine
		if inline.AdSystem != nil {
			ad.InLine.AdSystem.Name = formatVastString(inline.AdSystem.Name)
		}
		if inline.AdTitle != nil {
			inline.AdTitle.Name = formatVastString(inline.AdTitle.Name)
		}
		for i := range inline.Impressions {
			imp := &inline.Impressions[i]
			imp.URI = formatVastString(imp.URI)
		}
		for _, creative := range ad.InLine.Creatives {
			if creative.Linear != nil {
				linear := creative.Linear
				for i := range linear.TrackingEvents {
					track := &linear.TrackingEvents[i]
					track.URI = formatVastString(track.URI)
				}
				if linear.VideoClicks != nil {
					for i := range linear.VideoClicks.ClickThroughs {
						click := &linear.VideoClicks.ClickThroughs[i]
						click.URI = formatVastString(click.URI)
					}
					for i := range linear.VideoClicks.ClickTrackings {
						click := &linear.VideoClicks.ClickTrackings[i]
						click.URI = formatVastString(click.URI)
					}

					for i := range linear.VideoClicks.ClickTrackings {
						click := &linear.VideoClicks.ClickTrackings[i]
						click.URI = formatVastString(click.URI)
					}
				}
				for i := range linear.MediaFiles {
					media := &linear.MediaFiles[i]
					media.fixMediaFileAttrs(v.mediaFileDefaultHeight, v.mediaFileDefaultWidth)
					media.URI = formatVastString(media.URI)
				}
			}
			if creative.NonLinearAds != nil {
				nonLinears := creative.NonLinearAds
				for i := range nonLinears.TrackingEvents {
					track := &nonLinears.TrackingEvents[i]
					track.URI = formatVastString(track.URI)
				}
			}
		}
		for i := range inline.Error {
			err := &inline.Error[i]
			err.URI = formatVastString(err.URI)
		}
	}
}
