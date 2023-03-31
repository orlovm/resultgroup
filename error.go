package resultgroup

// errorWithUnwrap is an interface that represents an error with the ability to
// unwrap the underlying errors. This interface is compatible with Go 1.20
// wrapped errors.
type errorWithUnwrap interface {
	error
	Unwrap() []error
}

type multiError struct {
	errs []error
}

func (me *multiError) Error() string {
	var b []byte
	for i, err := range me.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

func (me *multiError) Unwrap() []error {
	return me.errs
}
