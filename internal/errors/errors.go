package errors

type Kind string

const (
	NotFound   Kind = "not_found"
	Invalid    Kind = "invalid"
	Internal   Kind = "internal"
	Unexpected Kind = "unexpected"
)

type Optional func(e *apiError)

func WithKind(kind Kind) Optional {
	return func(e *apiError) {
		e.kind = kind
	}
}

func WithError(err error) Optional {
	return func(e *apiError) {
		e.err = err
	}
}

type ApiError interface {
	Error() string
	Unwrap() error
}

type apiError struct {
	kind    Kind
	message string
	err     error
}

func (e *apiError) Kind() Kind    { return e.kind }
func (e *apiError) Unwrap() error { return e.err }

func (e *apiError) Error() string {
	msg := e.message
	if e.err != nil {
		msg += ". cause: " + e.err.Error()
	}
	return msg
}

func NewApiError(message string, optional ...Optional) ApiError {
	apiErr := &apiError{}
	for _, opt := range optional {
		opt(apiErr)
	}
	return apiErr
}
