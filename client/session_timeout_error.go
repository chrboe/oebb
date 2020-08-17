package oebb

type SessionTimeoutError struct {
}

func (e *SessionTimeoutError) Error() string {
	return "Session timed out, please authenticate again."
}
