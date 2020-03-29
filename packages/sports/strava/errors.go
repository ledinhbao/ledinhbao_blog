package strava

type TokenError struct {
	string
}

type stravaFault struct {
	Error   string
	Message string `json:"message"`
}

func (e *TokenError) Error() string {
	return e.string
}

func (sf *stravaFault) Error() string {
	return sf.Message
}
