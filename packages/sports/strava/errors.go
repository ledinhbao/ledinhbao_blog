package strava

type TokenError struct {
	string
}

func (e *TokenError) Error() string {
	return e.string
}
