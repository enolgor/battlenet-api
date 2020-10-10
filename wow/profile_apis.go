package wow

type ProfileAPI interface {
	GetProfile(token string) (*Profile, error)
}
