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
	//todo: add more validation logic, this is just an example.
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
	// (must check for nils: if something is nil, create it!!!!!) - TODO why??
	//if Creatives is nil, create create an array of 1 creative and
	// the create a linear & TrackingEvent -
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

// TODO: 2 options:
// 1. send full  *VideoClicks
// 2. update separately: ClickThroughs, ClickTrackings, CustomClicks
func (ad *Ad) AddClickTrackings(videoClicks ...VideoClick) {
	//todo: similar to AddTrackingEvents
	//todo: must check for nils: if something is nil, create it!
	if inline := ad.InLine; inline != nil {
		for i := range inline.Creatives {
			c := &inline.Creatives[i]
			linear := c.Linear
			if linear == nil {
				continue
			}
			//linear.VideoClicks = append(linear.VideoClicks, videoClicks...)
		}
	}
	if wrapper := ad.Wrapper; wrapper != nil {
		for i := range wrapper.Creatives {
			c := &wrapper.Creatives[i]
			linear := c.Linear
			if linear == nil {
				continue
			}
			//linear.VideoClicks = append(linear.TrackingEvents, trackingEvents...)
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
