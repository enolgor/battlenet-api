package wow

type Profile map[string]interface{}

func (wc *wowClientImpl) GetProfile(token string) (*Profile, error) {
	var profile Profile = make(map[string]interface{})
	err := wc.getProfileData("/profile/user/wow", token, &profile)
	return &profile, err
}
