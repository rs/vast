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
	//todo: must check for nils: if something is nil, create it!
	if ad.InLine != nil {
		//add companion to a creative and append the creative
	}
	if ad.Wrapper != nil {
		//add WrapperCompanion to a WrapperCreative and append the WrapperCreative

	}
}

func (ad *Ad) AddImpressions(impression ...Impression) {
	//todo: add more validation logic, this is just an example.
	if ad.InLine != nil {
		if ad.InLine.Impressions == nil {
			ad.InLine.Impressions = make([]Impression, 0, len(impression))
		}
		ad.InLine.Impressions = append(ad.InLine.Impressions, impression...)
	}
	if ad.Wrapper != nil {
		if ad.Wrapper.Impressions == nil {
			ad.Wrapper.Impressions = make([]Impression, 0, len(impression))
		}
		ad.Wrapper.Impressions = append(ad.Wrapper.Impressions, impression...)
	}
}

func (ad *Ad) AddErrors(error ...Error) {
	//todo: similar to AddImpressions
}

func (ad *Ad) AddTrackingEvents(trackingEvents ...Tracking) {
	//todo: need to append to both inline & wrapper, to each creative! (must check for nils: if something is nil, create it!!!!!)
	//if Creatives is nil, create create an array of 1 creative and the create a linear & TrackingEvents
	for _, c := range ad.InLine.Creatives {
		//must check for nils: if something is nil, create it!
		c.Linear.TrackingEvents = append(c.Linear.TrackingEvents, trackingEvents...)
	}
}

func (ad *Ad) AddClickTrackings(videoClicks ...VideoClick) {
	//todo: similar to AddTrackingEvents
	//todo: must check for nils: if something is nil, create it!
}
