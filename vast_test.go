package vast

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func loadFixture(path string) (*VAST, string, string, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, "", "", err
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var v VAST
	err = xml.Unmarshal(b, &v)

	res, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, "", "", err

	}

	return &v, string(b), string(res), err
}

func TestCreativeExtensions(t *testing.T) {
	v, _, _, err := loadFixture("testdata/creative_extensions.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "3.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "abc123", ad.ID)
		if assert.NotNil(t, ad.InLine) {
			if assert.Len(t, ad.InLine.Creatives, 1) {
				if assert.Len(t, ad.InLine.Creatives[0].CreativeExtensions, 4) {
					var ext Extension
					// asserting first extension
					ext = ad.InLine.Creatives[0].CreativeExtensions[0]
					assert.Equal(t, "geo", ext.Type)
					assert.Empty(t, ext.CustomTracking)
					assert.Equal(t, "\n              <Country>US</Country>\n              <Bandwidth>3</Bandwidth>\n              <BandwidthKbps>1680</BandwidthKbps>\n            ", string(ext.Data))
					// asserting second extension
					ext = ad.InLine.Creatives[0].CreativeExtensions[1]
					assert.Equal(t, "activeview", ext.Type)
					if assert.Len(t, ext.CustomTracking, 2) {
						// first tracker
						assert.Equal(t, "viewable_impression", ext.CustomTracking[0].Event)
						assert.Equal(t, "https://pubads.g.doubleclick.net/pagead/conversion/?ai=test&label=viewable_impression&acvw=[VIEWABILITY]&gv=[GOOGLE_VIEWABILITY]&ad_mt=[AD_MT]", ext.CustomTracking[0].URI)
						// second tracker
						assert.Equal(t, "abandon", ext.CustomTracking[1].Event)
						assert.Equal(t, "https://pubads.g.doubleclick.net/pagead/conversion/?ai=test&label=video_abandon&acvw=[VIEWABILITY]&gv=[GOOGLE_VIEWABILITY]", ext.CustomTracking[1].URI)
					}
					assert.Empty(t, string(ext.Data))
					// asserting third extension
					ext = ad.InLine.Creatives[0].CreativeExtensions[2]
					assert.Equal(t, "DFP", ext.Type)
					assert.Empty(t, ext.CustomTracking)
					assert.Equal(t, "\n              <SkippableAdType>Generic</SkippableAdType>\n            ", string(ext.Data))
					// asserting fourth extension
					ext = ad.InLine.Creatives[0].CreativeExtensions[3]
					assert.Equal(t, "metrics", ext.Type)
					assert.Empty(t, ext.CustomTracking)
					assert.Equal(t, "\n              <FeEventId>MubmWKCWLs_tiQPYiYrwBw</FeEventId>\n              <AdEventId>CIGpsPCTkdMCFdN-Ygod-xkCKQ</AdEventId>\n            ", string(ext.Data))
				}
			}
		}
	}
}

func TestInlineExtensions(t *testing.T) {
	v, _, _, err := loadFixture("testdata/inline_extensions.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "3.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "708365173", ad.ID)
		if assert.NotNil(t, ad.InLine) {
			if assert.NotNil(t, ad.InLine.Extensions) {
				if assert.Len(t, ad.InLine.Extensions, 4) {
					var ext Extension
					// asserting first extension
					ext = ad.InLine.Extensions[0]
					assert.Equal(t, "geo", ext.Type)
					assert.Empty(t, ext.CustomTracking)
					assert.Equal(t, "\n          <Country>US</Country>\n          <Bandwidth>3</Bandwidth>\n          <BandwidthKbps>1680</BandwidthKbps>\n        ", string(ext.Data))
					// asserting second extension
					ext = ad.InLine.Extensions[1]
					assert.Equal(t, "activeview", ext.Type)
					if assert.Len(t, ext.CustomTracking, 2) {
						// first tracker
						assert.Equal(t, "viewable_impression", ext.CustomTracking[0].Event)
						assert.Equal(t, "https://pubads.g.doubleclick.net/pagead/conversion/?ai=test&label=viewable_impression&acvw=[VIEWABILITY]&gv=[GOOGLE_VIEWABILITY]&ad_mt=[AD_MT]", ext.CustomTracking[0].URI)
						// second tracker
						assert.Equal(t, "abandon", ext.CustomTracking[1].Event)
						assert.Equal(t, "https://pubads.g.doubleclick.net/pagead/conversion/?ai=test&label=video_abandon&acvw=[VIEWABILITY]&gv=[GOOGLE_VIEWABILITY]", ext.CustomTracking[1].URI)
					}
					assert.Empty(t, string(ext.Data))
					// asserting third extension
					ext = ad.InLine.Extensions[2]
					assert.Equal(t, "DFP", ext.Type)
					assert.Equal(t, "\n          <SkippableAdType>Generic</SkippableAdType>\n        ", string(ext.Data))
					assert.Empty(t, ext.CustomTracking)
					// asserting fourth extension
					ext = ad.InLine.Extensions[3]
					assert.Equal(t, "metrics", ext.Type)
					assert.Equal(t, "\n          <FeEventId>MubmWKCWLs_tiQPYiYrwBw</FeEventId>\n          <AdEventId>CIGpsPCTkdMCFdN-Ygod-xkCKQ</AdEventId>\n        ", string(ext.Data))
					assert.Empty(t, ext.CustomTracking)
				}
			}
		}
	}
}

func TestInlineLinear(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_inline_linear.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "601364", ad.ID)
		assert.Nil(t, ad.Wrapper)
		assert.Equal(t, 0, ad.Sequence)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "Acudeo Compatible", inline.AdSystem.Name)
			assert.Equal(t, "1.0", inline.AdSystem.Version)
			assert.Equal(t, "VAST 2.0 Instream Test 1", inline.AdTitle.CDATA)
			assert.Equal(t, "VAST 2.0 Instream Test 1", inline.Description.CDATA)
			if assert.Len(t, inline.Errors, 2) {
				assert.Equal(t, "http://myErrorURL/error", inline.Errors[0].CDATA)
				assert.Equal(t, "http://myErrorURL/error2", inline.Errors[1].CDATA)
			}
			if assert.Len(t, inline.Impressions, 2) {
				assert.Equal(t, "http://myTrackingURL/impression", inline.Impressions[0].URI)
				assert.Equal(t, "http://myTrackingURL/impression2", inline.Impressions[1].URI)
				assert.Equal(t, "foo", inline.Impressions[1].ID)
			}
			if assert.Len(t, inline.Creatives, 2) {
				crea1 := inline.Creatives[0]
				assert.Equal(t, "601364", crea1.AdID)
				assert.Nil(t, crea1.NonLinearAds)
				assert.Nil(t, crea1.CompanionAds)
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear
					assert.Equal(t, Duration(30*time.Second), linear.Duration)
					if assert.Len(t, linear.TrackingEvents, 6) {
						assert.Equal(t, linear.TrackingEvents[0].Event, "creativeView")
						assert.Equal(t, linear.TrackingEvents[0].URI, "http://myTrackingURL/creativeView")
						assert.Equal(t, linear.TrackingEvents[1].Event, "start")
						assert.Equal(t, linear.TrackingEvents[1].URI, "http://myTrackingURL/start")
					}
					if assert.NotNil(t, linear.VideoClicks) {
						if assert.Len(t, linear.VideoClicks.ClickThroughs, 1) {
							assert.Equal(t, linear.VideoClicks.ClickThroughs[0].URI, "http://www.tremormedia.com")
						}
						if assert.Len(t, linear.VideoClicks.ClickTrackings, 1) {
							assert.Equal(t, linear.VideoClicks.ClickTrackings[0].URI, "http://myTrackingURL/click")
						}
						assert.Len(t, linear.VideoClicks.CustomClicks, 0)
					}
					if assert.Len(t, linear.MediaFiles, 1) {
						mf := linear.MediaFiles[0]
						assert.Equal(t, "progressive", mf.Delivery)
						assert.Equal(t, "video/x-flv", mf.Type)
						assert.Equal(t, 500, mf.Bitrate)
						assert.Equal(t, 400, mf.Width)
						assert.Equal(t, 300, mf.Height)
						assert.Equal(t, true, mf.Scalable)
						assert.Equal(t, true, mf.MaintainAspectRatio)
						assert.Equal(t, "http://cdnp.tremormedia.com/video/acudeo/Carrot_400x300_500kb.flv", mf.URI)
					}
				}

				crea2 := inline.Creatives[1]
				assert.Equal(t, "601364-Companion", crea2.AdID)
				assert.Nil(t, crea2.NonLinearAds)
				assert.Nil(t, crea2.Linear)
				if assert.NotNil(t, crea2.CompanionAds) {
					assert.Equal(t, "all", crea2.CompanionAds.Required)
					if assert.Len(t, crea2.CompanionAds.Companions, 2) {
						comp1 := crea2.CompanionAds.Companions[0]
						assert.Equal(t, 300, comp1.Width)
						assert.Equal(t, 250, comp1.Height)
						if assert.NotNil(t, comp1.StaticResource) {
							assert.Equal(t, "image/jpeg", comp1.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/Blistex1.jpg", comp1.StaticResource.URI)
						}
						if assert.Len(t, comp1.TrackingEvents, 1) {
							assert.Equal(t, "creativeView", comp1.TrackingEvents[0].Event)
							assert.Equal(t, "http://myTrackingURL/firstCompanionCreativeView", comp1.TrackingEvents[0].URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp1.CompanionClickThrough.CDATA)

						comp2 := crea2.CompanionAds.Companions[1]
						assert.Equal(t, 728, comp2.Width)
						assert.Equal(t, 90, comp2.Height)
						if assert.NotNil(t, comp2.StaticResource) {
							assert.Equal(t, "image/jpeg", comp2.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/728x90_banner1.jpg", comp2.StaticResource.URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp2.CompanionClickThrough.CDATA)
					}
				}
			}
		}
	}
}

func TestInlineLinearDurationUndefined(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_inline_linear-duration_undefined.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			if assert.Len(t, inline.Creatives, 1) {
				crea1 := inline.Creatives[0]
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear
					assert.Equal(t, Duration(0), linear.Duration)
				}
			}
		}
	}
}

func TestInlineNonLinear(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_inline_nonlinear.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "602678", ad.ID)
		assert.Nil(t, ad.Wrapper)
		assert.Equal(t, 0, ad.Sequence)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "Acudeo Compatible", inline.AdSystem.Name)
			assert.Equal(t, "NonLinear Test Campaign 1", inline.AdTitle.CDATA)
			assert.Equal(t, "NonLinear Test Campaign 1", inline.Description.CDATA)
			assert.Equal(t, "http://mySurveyURL/survey", inline.Survey.CDATA)
			if assert.Len(t, inline.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/error", inline.Errors[0].CDATA)
			}
			if assert.Len(t, inline.Impressions, 1) {
				assert.Equal(t, "http://myTrackingURL/impression", inline.Impressions[0].URI)
			}
			if assert.Len(t, inline.Creatives, 2) {
				crea1 := inline.Creatives[0]
				assert.Equal(t, "602678-NonLinear", crea1.AdID)
				assert.Nil(t, crea1.Linear)
				assert.Nil(t, crea1.CompanionAds)
				if assert.NotNil(t, crea1.NonLinearAds) {
					nonlin := crea1.NonLinearAds
					if assert.Len(t, nonlin.TrackingEvents, 5) {
						assert.Equal(t, nonlin.TrackingEvents[0].Event, "creativeView")
						assert.Equal(t, nonlin.TrackingEvents[0].URI, "http://myTrackingURL/nonlinear/creativeView")
						assert.Equal(t, nonlin.TrackingEvents[1].Event, "expand")
						assert.Equal(t, nonlin.TrackingEvents[1].URI, "http://myTrackingURL/nonlinear/expand")
					}
					if assert.Len(t, nonlin.NonLinears, 2) {
						assert.Equal(t, "image/jpeg", nonlin.NonLinears[0].StaticResource.CreativeType)
						assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/50x300_static.jpg", strings.TrimSpace(nonlin.NonLinears[0].StaticResource.URI))
						assert.Equal(t, "image/jpeg", nonlin.NonLinears[1].StaticResource.CreativeType)
						assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/50x450_static.jpg", strings.TrimSpace(nonlin.NonLinears[1].StaticResource.URI))
						assert.Equal(t, "http://www.tremormedia.com", strings.TrimSpace(nonlin.NonLinears[1].NonLinearClickThrough.CDATA))
					}
				}

				crea2 := inline.Creatives[1]
				assert.Equal(t, "602678-Companion", crea2.AdID)
				assert.Nil(t, crea2.NonLinearAds)
				assert.Nil(t, crea2.Linear)
				if assert.NotNil(t, crea2.CompanionAds) {
					if assert.Len(t, crea2.CompanionAds.Companions, 2) {
						comp1 := crea2.CompanionAds.Companions[0]
						assert.Equal(t, 300, comp1.Width)
						assert.Equal(t, 250, comp1.Height)
						if assert.NotNil(t, comp1.StaticResource) {
							assert.Equal(t, "application/x-shockwave-flash", comp1.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/300x250_companion_1.swf", comp1.StaticResource.URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp1.CompanionClickThrough.CDATA)

						comp2 := crea2.CompanionAds.Companions[1]
						assert.Equal(t, 728, comp2.Width)
						assert.Equal(t, 90, comp2.Height)
						if assert.NotNil(t, comp2.StaticResource) {
							assert.Equal(t, "image/jpeg", comp2.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/728x90_banner1.jpg", comp2.StaticResource.URI)
						}
						if assert.Len(t, comp2.TrackingEvents, 1) {
							assert.Equal(t, "creativeView", comp2.TrackingEvents[0].Event)
							assert.Equal(t, "http://myTrackingURL/secondCompanion", comp2.TrackingEvents[0].URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp2.CompanionClickThrough.CDATA)
					}
				}
			}
		}
	}
}

func TestWrapperLinear(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_wrapper_linear_1.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "602833", ad.ID)
		assert.Equal(t, 0, ad.Sequence)
		assert.Nil(t, ad.InLine)
		if assert.NotNil(t, ad.Wrapper) {
			wrapper := ad.Wrapper
			assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/vast_inline_linear.xml", wrapper.VASTAdTagURI.CDATA)
			assert.Equal(t, true, wrapper.FallbackOnNoAd)
			assert.Equal(t, true, wrapper.AllowMultipleAds)
			assert.Equal(t, true, wrapper.FollowAdditionalWrappers)
			assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/vast_inline_linear.xml", wrapper.VASTAdTagURI.CDATA)
			assert.Equal(t, "Acudeo Compatible", wrapper.AdSystem.Name)
			if assert.Len(t, wrapper.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/wrapper/error", wrapper.Errors[0].CDATA)
			}
			if assert.Len(t, wrapper.Impressions, 1) {
				assert.Equal(t, "http://myTrackingURL/wrapper/impression", wrapper.Impressions[0].URI)
			}

			if assert.Len(t, wrapper.Creatives, 3) {
				crea1 := wrapper.Creatives[0]
				assert.Equal(t, "602833", crea1.AdID)
				assert.Nil(t, crea1.NonLinearAds)
				assert.Nil(t, crea1.CompanionAds)
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear
					if assert.Len(t, linear.TrackingEvents, 11) {
						assert.Equal(t, linear.TrackingEvents[0].Event, "creativeView")
						assert.Equal(t, linear.TrackingEvents[0].URI, "http://myTrackingURL/wrapper/creativeView")
						assert.Equal(t, linear.TrackingEvents[1].Event, "start")
						assert.Equal(t, linear.TrackingEvents[1].URI, "http://myTrackingURL/wrapper/start")
					}
					assert.Nil(t, linear.VideoClicks)
				}

				crea2 := wrapper.Creatives[1]
				assert.Equal(t, "", crea2.AdID)
				assert.Nil(t, crea2.CompanionAds)
				assert.Nil(t, crea2.NonLinearAds)
				if assert.NotNil(t, crea2.Linear) {
					if assert.Len(t, crea2.Linear.VideoClicks.ClickTrackings, 1) {
						assert.Equal(t, "http://myTrackingURL/wrapper/click", crea2.Linear.VideoClicks.ClickTrackings[0].URI)
					}
				}

				crea3 := wrapper.Creatives[2]
				assert.Equal(t, "602833-NonLinearTracking", crea3.AdID)
				assert.Nil(t, crea3.CompanionAds)
				assert.Nil(t, crea3.Linear)
				if assert.NotNil(t, crea3.NonLinearAds) {
					if assert.Len(t, crea3.NonLinearAds.TrackingEvents, 1) {
						assert.Equal(t, "creativeView", crea3.NonLinearAds.TrackingEvents[0].Event)
						assert.Equal(t, "http://myTrackingURL/wrapper/creativeView", crea3.NonLinearAds.TrackingEvents[0].URI)
					}
				}
			}
		}
	}
}

func TestWrapperNonLinear(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_wrapper_nonlinear_1.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "602867", ad.ID)
		assert.Equal(t, 0, ad.Sequence)
		assert.Nil(t, ad.InLine)
		if assert.NotNil(t, ad.Wrapper) {
			wrapper := ad.Wrapper
			assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/vast_inline_nonlinear2.xml", wrapper.VASTAdTagURI.CDATA)
			assert.Equal(t, "Acudeo Compatible", wrapper.AdSystem.Name)
			if assert.Len(t, wrapper.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/wrapper/error", wrapper.Errors[0].CDATA)
			}
			if assert.Len(t, wrapper.Impressions, 1) {
				assert.Equal(t, "http://myTrackingURL/wrapper/impression", wrapper.Impressions[0].URI)
			}

			if assert.Len(t, wrapper.Creatives, 2) {
				crea1 := wrapper.Creatives[0]
				assert.Equal(t, "602867", crea1.AdID)
				assert.Nil(t, crea1.NonLinearAds)
				assert.Nil(t, crea1.CompanionAds)
				assert.NotNil(t, crea1.Linear)

				crea2 := wrapper.Creatives[1]
				assert.Equal(t, "602867-NonLinearTracking", crea2.AdID)
				assert.Nil(t, crea2.CompanionAds)
				assert.Nil(t, crea2.Linear)
				if assert.NotNil(t, crea2.NonLinearAds) {
					if assert.Len(t, crea2.NonLinearAds.TrackingEvents, 5) {
						assert.Equal(t, "creativeView", crea2.NonLinearAds.TrackingEvents[0].Event)
						assert.Equal(t, "http://myTrackingURL/wrapper/nonlinear/creativeView/creativeView", crea2.NonLinearAds.TrackingEvents[0].URI)
					}
				}
			}
		}
	}
}

func TestSpotXVpaid(t *testing.T) {
	v, _, _, err := loadFixture("testdata/spotx_vpaid.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "1130507-1818483", ad.ID)
		assert.Nil(t, ad.Wrapper)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "SpotXchange", inline.AdSystem.Name)
			assert.Equal(t, "1.0", inline.AdSystem.Version)
			assert.Equal(t, "IntegralAds_VAST_2_0_Ad_Wrapper", inline.AdTitle.CDATA)
			assert.Equal(t, "", inline.Description.CDATA)

			if assert.Len(t, inline.Creatives, 2) {
				crea1 := inline.Creatives[0]
				assert.Equal(t, 1, crea1.Sequence)
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear
					adParam, err := os.Open("testdata/spotx_adparameters.txt")
					if err != nil {
						assert.FailNow(t, "Cannot open adparams file")
					}
					defer adParam.Close()
					b, _ := ioutil.ReadAll(adParam)
					assert.Equal(t, linear.AdParameters.Parameters, string(b))
					if assert.Len(t, crea1.Linear.MediaFiles, 1) {
						media1 := crea1.Linear.MediaFiles[0]
						assert.Equal(t, "progressive", media1.Delivery)
						assert.Equal(t, "application/javascript", media1.Type)
						assert.Equal(t, 300, media1.Width)
						assert.Equal(t, 250, media1.Height)
						assert.Equal(t, "https://cdn.spotxcdn.com/integration/instreamadbroker/v1/instreamadbroker/beta.js", media1.URI)
					}
				}
				crea2 := inline.Creatives[1]
				assert.Equal(t, 1, crea2.Sequence)
				if assert.NotNil(t, crea2.CompanionAds) {
					if assert.Len(t, crea2.CompanionAds.Companions, 3) {
						companionAds1 := crea2.CompanionAds.Companions[0]
						assert.Equal(t, 300, companionAds1.Width)
						assert.Equal(t, 250, companionAds1.Height)
						assert.Equal(t, "medium_rectangle", companionAds1.ID)
						if assert.NotNil(t, companionAds1.IFrameResource) {
							assert.Equal(t, "https://search.spotxchange.com/banner?_a=137530&amp;_p=spotx&amp;_z=1&amp;_m=eNpVz99vgjAQB%2FD%2BLTyP0rM%2FLCZ7WOKWbFF5GJr4RAotUBUwgkHZ9rdvVB8W%2B3BJv3eXT45hQhBCwAnhEkKQKA2lZCHLhJYSwExTSgSkqcpBCEIIIKBTTglaf6JQ4glleAISAxOo7LpjOwuCvu9xrip7uJaq1tdK1ThrqgB9eefWnJLibLU385hRoTQ8FylPnSNzMgUtc6EpdY73dB%2B3bVKbflwYg6xp9tYkuqmUrceoNeqUlbg9Nt0lG7HCOGkcVGdtTZ2Z5IHyneU7zHea%2F8D93A6bcPR7e2i1XizuBW4R4oJcKPDx%2B99y5TuKD%2BU23rPlrhiW849dFG%2FstnqrVsPLJYrfYTm88mi%2BZ6uheP4DQppvcA%3D%3D&amp;_l=eNplj01rwkAYhN%2FfsleLvJv9yEbooUXoxU1BDKW5yGazNUk1SZMNtVr%2Fe9VgofQyh5mHGYYqqrhiEEYykIyDABWgAi4wlAAUKKdcRIpD3zZ%2BD5MJSBpB3dQOEI6EkdmRZGVOZgSniEySO9Kar2bwNwcvVt%2B6%2BpeJrtCQbUu7dntbmHrj1m%2FOXXOKiOfYDp3xLl%2FvTPfufLs11v1n8cKeThAni8UoCDhVAparR%2FDdsGuLxrttP7XNDlBxRjORW5pxhRYjIYULc8nDyDLm8vH3%2BRxy%2BBhM7a3p%2FdhKWSgC8XcGvuNK8%2FQp2b%2BuiiKdb1AfNNfVw6eukkBXcfk8T4L0ZVnqQ3L%2FAyjAYKo%3D&amp;_t=eNotj9Fu4yAQRf0tPFdbwK5dE%2B1b1aiVnGpXlVaJKlkDTGxaYyxMtsmm%2BfeFuDzA3MO9wyBhHNFn9Z5WQCtd6vt7xrCSOa2YlLBnZUkpZZnqwYyZ8x2MRmVySZ2JRW2gNZoIorWUpeSasbq4y%2BOOulaFrrGGUtccyQ3ZO28hRO%2F66TFKayy24TRhJMZCh2%2B3ndnHi0%2BjQx9hTmlUPZquTyl%2Bd5XBDq1ydhrQ4ph4hAZka%2BwhijTRwbYeVYCxG9KzM4JXfTtPqIg4kwk8WAzo56Q0%2FjUKU4XHkI55cuHYLriFafqekF3iwE7jQAS%2FIS5m%2BSWuLGM%2FKM1AVOI8i1KQbnASBrIygq6uYHJzWACrGY8sv%2F42gbyoKlpEVAiiTDh9m9gS9NgZNyZWsjqiKprcYQz%2B6uM8X3rppVfJ2eqSfW3Wv05b%2Fvy%2Be2j4br352NqN3f55Ys2%2Fj8%2FNQ3d8WTfHhv8eXl6bn%2F8BtbOdiA%3D%3D&amp;_b=eNozYEjRB8KUpCSzJKMUQ0NLE1NjIJmaYplskmKZaplolmJplKqXUZKbw%2BAX6uODIGp83T0NI7MysiOzHCv8wj0rI3MDq3yNwnKAtEFkuKtBVDiQrvI1jqxKtwUAEJMfUw%3D%3D&amp;resource_type=iframe", companionAds1.IFrameResource.CDATA)
						}
						companionAds2 := crea2.CompanionAds.Companions[1]
						assert.Equal(t, 300, companionAds2.Width)
						assert.Equal(t, 250, companionAds2.Height)
						assert.Equal(t, "medium_rectangle", companionAds2.ID)
						if assert.NotNil(t, companionAds2.HTMLResource) {
							htmlResource, err := os.Open("testdata/spotx_html_resource.html")
							if err != nil {
								assert.FailNow(t, "Cannot open spotx html resource file")
							}
							defer htmlResource.Close()
							b, _ := ioutil.ReadAll(htmlResource)
							assert.Equal(t, companionAds2.HTMLResource.HTML, string(b))
						}
						companionAds3 := crea2.CompanionAds.Companions[2]
						assert.Equal(t, 300, companionAds3.Width)
						assert.Equal(t, 250, companionAds3.Height)
						assert.Equal(t, "medium_rectangle", companionAds3.ID)
						if assert.NotNil(t, companionAds3.StaticResource) {
							assert.Equal(t, "image/gif", companionAds3.StaticResource.CreativeType)
							assert.Equal(t, "https://search.spotxchange.com/banner?_a=137530&amp;_p=spotx&amp;_z=1&amp;_m=eNpVz99vgjAQB%2FD%2BLTyP0rM%2FLCZ7WOKWbFF5GJr4RAotUBUwgkHZ9rdvVB8W%2B3BJv3eXT45hQhBCwAnhEkKQKA2lZCHLhJYSwExTSgSkqcpBCEIIIKBTTglaf6JQ4glleAISAxOo7LpjOwuCvu9xrip7uJaq1tdK1ThrqgB9eefWnJLibLU385hRoTQ8FylPnSNzMgUtc6EpdY73dB%2B3bVKbflwYg6xp9tYkuqmUrceoNeqUlbg9Nt0lG7HCOGkcVGdtTZ2Z5IHyneU7zHea%2F8D93A6bcPR7e2i1XizuBW4R4oJcKPDx%2B99y5TuKD%2BU23rPlrhiW849dFG%2FstnqrVsPLJYrfYTm88mi%2BZ6uheP4DQppvcA%3D%3D&amp;_l=eNplj01rwkAYhN%2FfsleLvJv9yEbooUXoxU1BDKW5yGazNUk1SZMNtVr%2Fe9VgofQyh5mHGYYqqrhiEEYykIyDABWgAi4wlAAUKKdcRIpD3zZ%2BD5MJSBpB3dQOEI6EkdmRZGVOZgSniEySO9Kar2bwNwcvVt%2B6%2BpeJrtCQbUu7dntbmHrj1m%2FOXXOKiOfYDp3xLl%2FvTPfufLs11v1n8cKeThAni8UoCDhVAparR%2FDdsGuLxrttP7XNDlBxRjORW5pxhRYjIYULc8nDyDLm8vH3%2BRxy%2BBhM7a3p%2FdhKWSgC8XcGvuNK8%2FQp2b%2BuiiKdb1AfNNfVw6eukkBXcfk8T4L0ZVnqQ3L%2FAyjAYKo%3D&amp;_t=eNotj9Fu4yAQRf0tPFdbwK5dE%2B1b1aiVnGpXlVaJKlkDTGxaYyxMtsmm%2BfeFuDzA3MO9wyBhHNFn9Z5WQCtd6vt7xrCSOa2YlLBnZUkpZZnqwYyZ8x2MRmVySZ2JRW2gNZoIorWUpeSasbq4y%2BOOulaFrrGGUtccyQ3ZO28hRO%2F66TFKayy24TRhJMZCh2%2B3ndnHi0%2BjQx9hTmlUPZquTyl%2Bd5XBDq1ydhrQ4ph4hAZka%2BwhijTRwbYeVYCxG9KzM4JXfTtPqIg4kwk8WAzo56Q0%2FjUKU4XHkI55cuHYLriFafqekF3iwE7jQAS%2FIS5m%2BSWuLGM%2FKM1AVOI8i1KQbnASBrIygq6uYHJzWACrGY8sv%2F42gbyoKlpEVAiiTDh9m9gS9NgZNyZWsjqiKprcYQz%2B6uM8X3rppVfJ2eqSfW3Wv05b%2Fvy%2Be2j4br352NqN3f55Ys2%2Fj8%2FNQ3d8WTfHhv8eXl6bn%2F8BtbOdiA%3D%3D&amp;_b=eNpFxl0LgjAUgGF%2FkSc%2FEhZ0EUje6AFFibxznuXmsklbGNKPL6%2FihYd350nnZnsA6Onh29m49za9mWASpDpw8jVxC%2BapBqAt4jzhIQUBi%2FfRT0Gsj4kJ1iXEQuEP6uZhk%2Bd%2FPsVaKUxPK6Z3haOO26xZ2gvKoi6Xa30eMas0rmWEtT5%2BAZ5CLzA%3D", companionAds3.StaticResource.URI)
						}
						if assert.NotNil(t, companionAds3.CompanionClickThrough) {
							assert.Equal(t, "https://search.spotxchange.com/click?_a=137530&amp;_p=spotx&amp;_z=1&amp;_m=eNpVz99vgjAQB%2FD%2BLTyP0rM%2FLCZ7WOKWbFF5GJr4RAotUBUwgkHZ9rdvVB8W%2B3BJv3eXT45hQhBCwAnhEkKQKA2lZCHLhJYSwExTSgSkqcpBCEIIIKBTTglaf6JQ4glleAISAxOo7LpjOwuCvu9xrip7uJaq1tdK1ThrqgB9eefWnJLibLU385hRoTQ8FylPnSNzMgUtc6EpdY73dB%2B3bVKbflwYg6xp9tYkuqmUrceoNeqUlbg9Nt0lG7HCOGkcVGdtTZ2Z5IHyneU7zHea%2F8D93A6bcPR7e2i1XizuBW4R4oJcKPDx%2B99y5TuKD%2BU23rPlrhiW849dFG%2FstnqrVsPLJYrfYTm88mi%2BZ6uheP4DQppvcA%3D%3D&amp;_l=eNplj01rwkAYhN%2FfsleLvJv9yEbooUXoxU1BDKW5yGazNUk1SZMNtVr%2Fe9VgofQyh5mHGYYqqrhiEEYykIyDABWgAi4wlAAUKKdcRIpD3zZ%2BD5MJSBpB3dQOEI6EkdmRZGVOZgSniEySO9Kar2bwNwcvVt%2B6%2BpeJrtCQbUu7dntbmHrj1m%2FOXXOKiOfYDp3xLl%2FvTPfufLs11v1n8cKeThAni8UoCDhVAparR%2FDdsGuLxrttP7XNDlBxRjORW5pxhRYjIYULc8nDyDLm8vH3%2BRxy%2BBhM7a3p%2FdhKWSgC8XcGvuNK8%2FQp2b%2BuiiKdb1AfNNfVw6eukkBXcfk8T4L0ZVnqQ3L%2FAyjAYKo%3D&amp;_t=eNotj0tv3CAUhf1bWEcJYMeOPeomqhplMVOp8iwcVUI87tgkYCzMNPPo%2FPeAHRZwz8c5l4s0Wn5k9QFXHFeqVE9PhEAlclwRIfiBlCXGmGRy4HrMnO%2F5qGUm%2BDiCZ3LJXpEFpTnTCjVIKSFKQRUhdfGYxx1ULQtVQ81LVVNAd%2BjgvOUhel9ef0VptQUWzhNEoi3v4e9Drw%2Fx4lOrMESYYxzVALofUoo%2BLjJYw6SzkwELY%2BIRai6Ytsco0kRHyzzIwMfepGdn4F4ObJ5AouaKJu65hQB%2BTkrBPy0hVXAK6ZgnF05sxYxP0%2FeE5BYHdgoMaugdcjFLb3FlGbnHOONN1VznpmxQb5zgBm10gzcLmNwcVkBqQiPLl98mkBdVhYuIigZJHc7fJrIGPfTajYmVpI6oiiZ3HINffJTmay%2B19iop2dyy%2F93lj%2Bnandm1%2B3P3PrzvXl4%2Fu8v2tG3fzNbuL7tW6a59Hn7%2F3P%2F4Ap5ioKM%3D&amp;_b=eNozYMgoKSkottLXLy7IL6lIzkjMS0%2FVS87PZfAL9fFhsEwzsLQ0NzQ0S7GwMDRMNU8yNjA3TEpKTDM0M3MDArCqmih3V5PIrGxDX5d0g8gqr0x%2Fl0hDfxe3bN8sTwO%2FrLDMyKycDD8XEJ1uCwCmYiLM", companionAds3.CompanionClickThrough.CDATA)
						}
						if assert.NotNil(t, companionAds3.AltText) {
							assert.Equal(t, "IntegralAds_VAST_2_0_Ad_Wrapper", companionAds3.AltText)
						}

					}
				}

			}
			if assert.Len(t, inline.Extensions, 2) {
				ext1 := inline.Extensions[0]
				assert.Equal(t, "LR-Pricing", ext1.Type)
				assert.Equal(t, "<Price model=\"CPM\" currency=\"USD\" source=\"spotxchange\"><![CDATA[3.06]]></Price>", strings.TrimSpace(string(ext1.Data)))
				ext2 := inline.Extensions[1]
				assert.Equal(t, "SpotX-Count", ext2.Type)
				assert.Equal(t, "<total_available><![CDATA[1]]></total_available>", strings.TrimSpace(string(ext2.Data)))
			}
		}
	}
}

func TestExtraSpacesVpaid(t *testing.T) {
	v, _, _, err := loadFixture("testdata/extraspaces_vpaid.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "1130507-1818483", ad.ID)
		assert.Nil(t, ad.Wrapper)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "SpotXchange", inline.AdSystem.Name)
			assert.Equal(t, "1.0", inline.AdSystem.Version)
			assert.Equal(t, "IntegralAds_VAST_2_0_Ad_Wrapper", inline.AdTitle.CDATA)
			assert.Equal(t, "", inline.Description.CDATA)

			if assert.Len(t, inline.Creatives, 1) {
				crea1 := inline.Creatives[0]
				assert.Equal(t, 1, crea1.Sequence)
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear

					assert.Equal(t, "        \n                  <VAST></VAST>\n                  \n                  ", linear.AdParameters.Parameters)
					if assert.Len(t, crea1.Linear.MediaFiles, 1) {
						media1 := crea1.Linear.MediaFiles[0]
						assert.Equal(t, "progressive", media1.Delivery)
						assert.Equal(t, "application/javascript", media1.Type)
						assert.Equal(t, 300, media1.Width)
						assert.Equal(t, 250, media1.Height)
						assert.Equal(t, "\n                     https://dummy.com/dummmy.js             \n                     ", media1.URI)
					}
				}
			}
		}
	}
}
