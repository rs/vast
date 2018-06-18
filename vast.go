// Package vast implements IAB VAST 3.0 specification http://www.iab.net/media/file/VASTv3.0.pdf
package vast

import "encoding/xml"

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
	Errors []CDATAString `xml:"Error,omitempty"`
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

// CDATAString ...
type CDATAString struct {
	CDATA string `xml:",cdata"`
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
	AdTitle CDATAString
	// One or more URIs that directs the video player to a tracking resource file that the
	// video player should request when the first frame of the ad is displayed
	Impressions []Impression `xml:"Impression"`
	// The container for one or more <Creative> elements
	Creatives []Creative `xml:"Creatives>Creative"`
	// A string value that provides a longer description of the ad.
	Description CDATAString `xml:",omitempty"`
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
	Survey CDATAString `xml:",omitempty"`
	// A URI representing an error-tracking pixel; this element can occur multiple
	// times.
	Errors []CDATAString `xml:"Error,omitempty"`
	// Provides a value that represents a price that can be used by real-time bidding
	// (RTB) systems. VAST is not designed to handle RTB since other methods exist,
	// but this element is offered for custom solutions if needed.
	Pricing *Pricing `xml:",omitempty"`
	// XML node for custom extensions, as defined by the ad server. When used, a
	// custom element should be nested under <Extensions> to help separate custom
	// XML elements from VAST elements. The following example includes a custom
	// xml element within the Extensions element.
	Extensions *[]Extension `xml:"Extensions>Extension,omitempty"`
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
	Value string `xml:",cdata"`
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
	VASTAdTagURI CDATAString
	// One or more URIs that directs the video player to a tracking resource file that the
	// video player should request when the first frame of the ad is displayed
	Impressions []Impression `xml:"Impression"`
	// A URI representing an error-tracking pixel; this element can occur multiple
	// times.
	Errors []CDATAString `xml:"Error,omitempty"`
	// The container for one or more <Creative> elements
	Creatives []CreativeWrapper `xml:"Creatives>Creative"`
	// XML node for custom extensions, as defined by the ad server. When used, a
	// custom element should be nested under <Extensions> to help separate custom
	// XML elements from VAST elements. The following example includes a custom
	// xml element within the Extensions element.
	Extensions []Extension `xml:"Extensions>Extension,omitempty"`

	FallbackOnNoAd           *bool `xml:"fallbackOnNoAd,attr,omitempty"`
	AllowMultipleAds         *bool `xml:"allowMultipleAds,attr,omitempty"`
	FollowAdditionalWrappers *bool `xml:"followAdditionalWrappers,attr,omitempty"`
}

// AdSystem contains information about the system that returned the ad
type AdSystem struct {
	Version string `xml:"version,attr,omitempty"`
	Name    string `xml:",cdata"`
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
	// If defined, defines non linear creatives
	NonLinearAds *NonLinearAds `xml:",omitempty"`
	// If present, provides a VAST 4.x universal ad id
	UniversalAdID *UniversalAdID `xml:"UniversalAdId,omitempty"`
	// When an API framework is needed to execute creative, a
	// <CreativeExtensions> element can be added under the <Creative>. This
	// extension can be used to load an executable creative with or without using
	// a media file.
	// A <CreativeExtension> element is nested under the <CreativeExtensions>
	// (plural) element so that any xml extensions are separated from VAST xml.
	// Additionally, any xml used in this extension should identify an xml name
	// space (xmlns) to avoid confusing any of the xml element names with those
	// of VAST.
	// The nested <CreativeExtension> includes an attribute for type, which
	// specifies the MIME type needed to execute the extension.
	CreativeExtensions *[]Extension `xml:"CreativeExtensions>CreativeExtension,omitempty"`
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
	// If defined, defines companions creatives
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
	Duration       Duration
	AdParameters   *AdParameters `xml:",omitempty"`
	Icons          *Icons
	TrackingEvents []Tracking   `xml:"TrackingEvents>Tracking,omitempty"`
	VideoClicks    *VideoClicks `xml:",omitempty"`
	MediaFiles     []MediaFile  `xml:"MediaFiles>MediaFile,omitempty"`
}

// LinearWrapper defines a wrapped linear creative
type LinearWrapper struct {
	Icons          *Icons
	TrackingEvents []Tracking   `xml:"TrackingEvents>Tracking,omitempty"`
	VideoClicks    *VideoClicks `xml:",omitempty"`
}

// Companion defines a companion ad
type Companion struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty"`
	// Pixel dimensions of companion slot.
	Width int `xml:"width,attr,omitempty"`
	// Pixel dimensions of companion slot.
	Height int `xml:"height,attr,omitempty"`
	// Pixel dimensions of the companion asset.
	AssetWidth int `xml:"assetWidth,attr,omitempty"`
	// Pixel dimensions of the companion asset.
	AssetHeight int `xml:"assetHeight,attr,omitempty"`
	// Pixel dimensions of expanding companion ad when in expanded state.
	ExpandedWidth int `xml:"expandedWidth,attr,omitempty"`
	// Pixel dimensions of expanding companion ad when in expanded state.
	ExpandedHeight int `xml:"expandedHeight,attr,omitempty"`
	// The apiFramework defines the method to use for communication with the companion.
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// Used to match companion creative to publisher placement areas on the page.
	AdSlotID string `xml:"adSlotId,attr,omitempty"`
	// URL to open as destination page when user clicks on the the companion banner ad.
	CompanionClickThrough CDATAString `xml:",omitempty"`
	// URLs to ping when user clicks on the the companion banner ad.
	CompanionClickTracking []CDATAString `xml:",omitempty"`
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
	IFrameResource CDATAString `xml:",omitempty"`
	// HTML to display the companion element
	HTMLResource *HTMLResource `xml:",omitempty"`
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
	ExpandedHeight int `xml:"expandedHeight,attr"`
	// The apiFramework defines the method to use for communication with the companion.
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// Used to match companion creative to publisher placement areas on the page.
	AdSlotID string `xml:"adSlotId,attr,omitempty"`
	// URL to open as destination page when user clicks on the the companion banner ad.
	CompanionClickThrough CDATAString `xml:",omitempty"`
	// URLs to ping when user clicks on the the companion banner ad.
	CompanionClickTracking []CDATAString `xml:",omitempty"`
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
	IFrameResource CDATAString `xml:",omitempty"`
	// HTML to display the companion element
	HTMLResource *HTMLResource `xml:",omitempty"`
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
	ExpandedHeight int `xml:"expandedHeight,attr"`
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
	NonLinearClickTracking []CDATAString `xml:",omitempty"`
	// URL to open as destination page when user clicks on the non-linear ad unit.
	NonLinearClickThrough CDATAString `xml:",omitempty"`
	// Data to be passed into the video ad.
	AdParameters *AdParameters `xml:",omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource CDATAString `xml:",omitempty"`
	// HTML to display the companion element
	HTMLResource *HTMLResource `xml:",omitempty"`
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
	ExpandedHeight int `xml:"expandedHeight,attr"`
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
	NonLinearClickTracking []CDATAString `xml:",omitempty"`
}

type Icons struct {
	XMLName xml.Name `xml:"Icons,omitempty"`
	Icon    []Icon   `xml:"Icon,omitempty"`
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
	Duration Duration `xml:"duration,attr"`
	// The apiFramework defines the method to use for communication with the icon element
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	// URL to open as destination page when user clicks on the icon.
	IconClickThrough CDATAString `xml:"IconClicks>IconClickThrough,omitempty"`
	// URLs to ping when user clicks on the the icon.
	IconClickTrackings []CDATAString `xml:"IconClicks>IconClickTracking,omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource CDATAString `xml:",omitempty"`
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
	HTML       string `xml:",cdata"`
}

// AdParameters defines arbitrary ad parameters
type AdParameters struct {
	// Specifies whether the parameters are XML-encoded
	XMLEncoded bool   `xml:"xmlEncoded,attr,omitempty"`
	Parameters string `xml:",cdata"`
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

// UniversalAdID describes a VAST 4.x universal ad id.
type UniversalAdID struct {
	IDRegistry string `xml:"idRegistry,attr"`
	IDValue    string `xml:"idValue,attr"`
	ID         string `xml:",cdata"`
}
