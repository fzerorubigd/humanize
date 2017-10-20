package humanize

type lateBinder interface {
	lateBind() error
}

func lateBind(t interface{}) error {
	if v, ok := t.(lateBinder); ok {
		return v.lateBind()
	}

	return nil
}
