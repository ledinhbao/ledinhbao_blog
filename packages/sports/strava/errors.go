package strava

type TokenError struct {
	string
}

type stravaFault struct {
	Message string `json:"message"`
}

func (e *TokenError) Error() string {
	return e.string
}

func (sf *stravaFault) Error() string {
	return sf.Message
}
