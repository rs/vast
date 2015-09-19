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

func loadFixture(path string) (*VAST, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var v VAST
	err = xml.Unmarshal(b, &v)
	return &v, err
}

func TestInlineLinear(t *testing.T) {
	v, err := loadFixture("testdata/vast_inline_linear.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ad, 1) {
		ad := v.Ad[0]
		assert.Equal(t, "601364", ad.ID)
		assert.Nil(t, ad.Wrapper)
		assert.Nil(t, ad.Sequence)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "Acudeo Compatible", inline.AdSystem.Name)
			assert.Equal(t, "1.0", inline.AdSystem.Version)
			assert.Equal(t, "VAST 2.0 Instream Test 1", inline.AdTitle)
			assert.Equal(t, "VAST 2.0 Instream Test 1", inline.Description)
			if assert.Len(t, inline.Errors, 2) {
				assert.Equal(t, "http://myErrorURL/error", inline.Errors[0])
				assert.Equal(t, "http://myErrorURL/error2", inline.Errors[1])
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
						assert.Equal(t, "http://www.tremormedia.com", comp1.CompanionClickThrough)

						comp2 := crea2.CompanionAds.Companions[1]
						assert.Equal(t, 728, comp2.Width)
						assert.Equal(t, 90, comp2.Height)
						if assert.NotNil(t, comp2.StaticResource) {
							assert.Equal(t, "image/jpeg", comp2.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/728x90_banner1.jpg", comp2.StaticResource.URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp2.CompanionClickThrough)
					}
				}
			}
		}
	}
}

func TestInlineNonLinear(t *testing.T) {
	v, err := loadFixture("testdata/vast_inline_nonlinear.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ad, 1) {
		ad := v.Ad[0]
		assert.Equal(t, "602678", ad.ID)
		assert.Nil(t, ad.Wrapper)
		assert.Nil(t, ad.Sequence)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "Acudeo Compatible", inline.AdSystem.Name)
			assert.Equal(t, "NonLinear Test Campaign 1", inline.AdTitle)
			assert.Equal(t, "NonLinear Test Campaign 1", inline.Description)
			assert.Equal(t, "http://mySurveyURL/survey", inline.Survey)
			if assert.Len(t, inline.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/error", inline.Errors[0])
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
						assert.Equal(t, "http://www.tremormedia.com", strings.TrimSpace(nonlin.NonLinears[1].NonLinearClickThrough))
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
						assert.Equal(t, "http://www.tremormedia.com", comp1.CompanionClickThrough)

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
						assert.Equal(t, "http://www.tremormedia.com", comp2.CompanionClickThrough)
					}
				}
			}
		}
	}
}

func TestWrapperLinear(t *testing.T) {
	v, err := loadFixture("testdata/vast_wrapper_linear_1.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ad, 1) {
		ad := v.Ad[0]
		assert.Equal(t, "602833", ad.ID)
		assert.Nil(t, ad.Sequence)
		assert.Nil(t, ad.InLine)
		if assert.NotNil(t, ad.Wrapper) {
			wrapper := ad.Wrapper
			assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/vast_inline_linear.xml", wrapper.VASTAdTagURI)
			assert.Equal(t, "Acudeo Compatible", wrapper.AdSystem.Name)
			if assert.Len(t, wrapper.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/wrapper/error", wrapper.Errors[0])
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
	v, err := loadFixture("testdata/vast_wrapper_nonlinear_1.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ad, 1) {
		ad := v.Ad[0]
		assert.Equal(t, "602867", ad.ID)
		assert.Nil(t, ad.Sequence)
		assert.Nil(t, ad.InLine)
		if assert.NotNil(t, ad.Wrapper) {
			wrapper := ad.Wrapper
			assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/vast_inline_nonlinear2.xml", wrapper.VASTAdTagURI)
			assert.Equal(t, "Acudeo Compatible", wrapper.AdSystem.Name)
			if assert.Len(t, wrapper.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/wrapper/error", wrapper.Errors[0])
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
