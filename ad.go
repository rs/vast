package vast

func (ad *Ad) GetExtensions() *Extensions {
	if ad.InLine != nil {
		return ad.InLine.Extensions
	}
	if ad.Wrapper != nil {
		return ad.Wrapper.Extensions
	}
	return nil
}

func (ad *Ad) AddExtension(extension *Extension) {
	if extension == nil {
		return
	}
	if ad.InLine != nil {
		ad.InLine.Extensions = addVastExtension(ad.InLine.Extensions, extension)
	}
	if ad.Wrapper != nil {
		ad.Wrapper.Extensions = addVastExtension(ad.Wrapper.Extensions, extension)
	}
}

func addVastExtension(extensions *Extensions, extension *Extension) *Extensions {
	if extensions == nil {
		extensions = new(Extensions)
	}
	extensions.Extensions = append(extensions.Extensions, *extension)
	return extensions
}

func (ad *Ad) AddCompanion(companion *Companion) {
	if companion == nil {
		return
	}
	if inLine := ad.InLine; inLine != nil {
		//add companion to a creative and append the creative
		inLine.Creatives = append(inLine.Creatives, Creative{
			CompanionAds: &CompanionAds{
				Companions: []Companion{*companion},
			},
		})
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		//add WrapperCompanion to a WrapperCreative and append the WrapperCreative
		wrapper.Creatives = append(wrapper.Creatives, CreativeWrapper{
			CompanionAds: &CompanionAdsWrapper{
				Companions: []CompanionWrapper{{
					StaticResource:        companion.StaticResource,
					CompanionClickThrough: companion.CompanionClickThrough.URI,
					TrackingEvents:        companion.TrackingEvents,
				}},
			},
		})
	}
}

func (ad *Ad) AddImpressions(impressions ...Impression) {
	if len(impressions) == 0 {
		return
	}
	if inline := ad.InLine; inline != nil {
		if inline.Impressions == nil {
			inline.Impressions = make([]Impression, 0, len(impressions))
		}
		inline.Impressions = append(inline.Impressions, impressions...)
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		if wrapper.Impressions == nil {
			wrapper.Impressions = make([]Impression, 0, len(impressions))
		}
		wrapper.Impressions = append(wrapper.Impressions, impressions...)
	}
}

func (ad *Ad) AddErrors(errors ...Error) {
	if len(errors) == 0 {
		return
	}
	if inline := ad.InLine; inline != nil {
		if ad.InLine.Error == nil {
			ad.InLine.Error = make([]Error, 0, len(errors))
		}
		ad.InLine.Error = append(ad.InLine.Error, errors...)
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		if wrapper.Error == nil {
			wrapper.Error = make([]Error, 0, len(errors))
		}
		wrapper.Error = append(wrapper.Error, errors...)
	}
}

func (ad *Ad) AddTrackingEvents(trackingEvents ...Tracking) {
	if len(trackingEvents) == 0 {
		return
	}
	if inline := ad.InLine; inline != nil {
		for i := range inline.Creatives {
			c := &inline.Creatives[i]
			linear := c.Linear
			if linear == nil {
				continue
			}
			linear.TrackingEvents = append(linear.TrackingEvents, trackingEvents...)
		}
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		if len(wrapper.Creatives) == 0 {
			wrapper.Creatives = []CreativeWrapper{{Linear: &LinearWrapper{}}}
		}
		for i := range wrapper.Creatives {
			c := &wrapper.Creatives[i]
			linear := c.Linear
			if linear == nil {
				continue
			}
			linear.TrackingEvents = append(linear.TrackingEvents, trackingEvents...)
		}
	}
}

func (ad *Ad) AddCompanionTrackingEvents(trackingEvents ...Tracking) {
	if len(trackingEvents) == 0 {
		return
	}
	if inline := ad.InLine; inline != nil {
		for i := range inline.Creatives {
			c := &inline.Creatives[i]
			cAds := c.CompanionAds
			if cAds == nil || len(cAds.Companions) == 0 {
				continue
			}
			for j := range cAds.Companions {
				cAd := &cAds.Companions[j]
				cAdsTrackingEvents := cAd.TrackingEvents
				if len(cAdsTrackingEvents) == 0 {
					cAdsTrackingEvents = []Tracking{}
				}
				cAdsTrackingEvents = append(cAdsTrackingEvents, trackingEvents...)
				cAd.TrackingEvents = cAdsTrackingEvents
			}
		}
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		for i := range wrapper.Creatives {
			c := &wrapper.Creatives[i]
			cAds := c.CompanionAds
			if cAds == nil || len(cAds.Companions) == 0 {
				continue
			}
			for j := range cAds.Companions {
				cAd := &cAds.Companions[j]
				cAdsTrackingEvents := cAd.TrackingEvents
				if len(cAdsTrackingEvents) == 0 {
					cAdsTrackingEvents = []Tracking{}
				}
				cAdsTrackingEvents = append(cAdsTrackingEvents, trackingEvents...)
				cAd.TrackingEvents = cAdsTrackingEvents
			}
		}
	}
}

func (ad *Ad) AddClickTrackings(clickTrackings ...VideoClick) {
	if len(clickTrackings) == 0 {
		return
	}
	if inline := ad.InLine; inline != nil {
		for i := range inline.Creatives {
			c := &inline.Creatives[i]
			linear := c.Linear
			if linear == nil {
				continue
			}
			videoClicks := linear.VideoClicks
			if videoClicks == nil {
				videoClicks = &VideoClicks{}
				linear.VideoClicks = videoClicks
			}
			videoClicks.ClickTrackings = append(videoClicks.ClickTrackings, clickTrackings...)
		}
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		if len(wrapper.Creatives) == 0 {
			wrapper.Creatives = []CreativeWrapper{{Linear: &LinearWrapper{}}}
		}
		for i := range wrapper.Creatives {
			c := &wrapper.Creatives[i]
			linear := c.Linear
			if linear == nil {
				continue
			}
			videoClicks := linear.VideoClicks
			if videoClicks == nil {
				videoClicks = &VideoClicks{}
				c.Linear.VideoClicks = videoClicks
			}
			c.Linear.VideoClicks.ClickTrackings = append(c.Linear.VideoClicks.ClickTrackings, clickTrackings...)
		}
	}
}

func (ad *Ad) AddClickThrough(clickThroughs ...VideoClick) {
	if len(clickThroughs) == 0 {
		return
	}
	if inline := ad.InLine; inline != nil {
		for i := range inline.Creatives {
			c := &inline.Creatives[i]
			linear := c.Linear
			if linear == nil {
				continue
			}
			videoClicks := linear.VideoClicks
			if videoClicks == nil {
				videoClicks = &VideoClicks{}
				linear.VideoClicks = videoClicks
			}
			videoClicks.ClickThroughs = append(videoClicks.ClickThroughs, clickThroughs...)
		}
	}
}

func (ad *Ad) AddCompanionClickTrackings(CompanionClickTrackings ...*CompanionClickTracking) {
	if len(CompanionClickTrackings) == 0 {
		return
	}
	if inline := ad.InLine; inline != nil {
		for i := range inline.Creatives {
			c := &inline.Creatives[i]
			cAds := c.CompanionAds
			if cAds == nil || len(cAds.Companions) == 0 {
				continue
			}
			for j := range cAds.Companions {
				cAd := &cAds.Companions[j]
				cAdsClickTrackings := cAd.CompanionClickTracking
				if len(cAdsClickTrackings) == 0 {
					cAdsClickTrackings = []*CompanionClickTracking{}
				}
				cAdsClickTrackings = append(cAdsClickTrackings, CompanionClickTrackings...)
				cAd.CompanionClickTracking = cAdsClickTrackings
			}
		}
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		for i := range wrapper.Creatives {
			c := &wrapper.Creatives[i]
			cAds := c.CompanionAds
			if cAds == nil || len(cAds.Companions) == 0 {
				continue
			}
			for j := range cAds.Companions {
				cAd := &cAds.Companions[j]
				cAdsClickTrackings := cAd.CompanionClickTracking
				if len(cAdsClickTrackings) == 0 {
					cAdsClickTrackings = []string{}
				}
				CompanionClickTrackingStrings := make([]string, 0, len(CompanionClickTrackings))
				for i, companionClickTracking := range CompanionClickTrackings {
					CompanionClickTrackingStrings[i] = companionClickTracking.URI
				}
				cAdsClickTrackings = append(cAdsClickTrackings, CompanionClickTrackingStrings...)
				cAd.CompanionClickTracking = cAdsClickTrackings
			}
		}
	}
}

func (ad *Ad) SetAdSystem(name, version string) {
	if ad.InLine != nil {
		ad.InLine.AdSystem = &AdSystem{Name: name, Version: version}
	}
	if ad.Wrapper != nil {
		ad.Wrapper.AdSystem = &AdSystem{Name: name, Version: version}
	}
}
