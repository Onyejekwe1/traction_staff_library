package config

type Settings struct {
	SiteFullName, SiteSlogan, SiteBaseURL, SiteTopMenuLogo, SiteProperDomainName, SiteShortName, SiteEmail, SitePhoneNumbers, SiteCompanyAddress string
	SiteYear                                                                                                                                     int
}

var SiteSettings = Settings{
	SiteFullName:         SiteFullName,
	SiteSlogan:           SiteSlogan,
	SiteBaseURL:          SiteBaseURL,
	SiteTopMenuLogo:      SiteTopMenuLogo,
	SiteProperDomainName: SiteProperDomainName,
	SiteShortName:        SiteShortName,
	SiteEmail:            SiteEmail,
	SitePhoneNumbers:     SitePhoneNumbers,
	SiteCompanyAddress:   SiteCompanyAddress,
	SiteYear:             SiteYear,
}
