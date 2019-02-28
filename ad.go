package vast

func (ad *Ad) AddExtension(extension *Extension) {
	if ad.InLine != nil {
		ad.InLine.Extensions = addVastExtension(ad.InLine.Extensions, extension)
	}
	if ad.Wrapper != nil {
		ad.Wrapper.Extensions = addVastExtension(ad.Wrapper.Extensions, extension)
	}
}

func addVastExtension(extensions *Extensions, extension *Extension) *Extensions {
	if extension == nil {
		return nil
	}
	if extensions == nil {
		extensions = &Extensions{}
	}
	extensions.Extensions = append(extensions.Extensions, *extension)
	return extensions
}

func (ad *Ad) AddCompanion(companion *Companion) {
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
		if ad.InLine.Errors == nil {
			ad.InLine.Errors = make([]Error, 0, len(errors))
		}
		ad.InLine.Errors = append(ad.InLine.Errors, errors...)
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		if wrapper.Errors == nil {
			wrapper.Errors = make([]Error, 0, len(errors))
		}
		wrapper.Errors = append(wrapper.Errors, errors...)
	}
}

func (ad *Ad) AddTrackingEvents(trackingEvents ...Tracking) {
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
			}
			videoClicks.ClickTrackings = append(videoClicks.ClickTrackings, clickTrackings...)
		}
	}
}
