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
	if ad.InLine != nil {
		//add companion to a creative and append the creative
	}
	if ad.Wrapper != nil {
		//add WrapperCompanion to a WrapperCreative and append the WrapperCreative

	}
}

func (ad *Ad) AddImpression() {

}

func (ad *Ad) AddClick() {

}

func (ad *Ad) AddX() {

}
