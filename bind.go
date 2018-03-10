package humanize

type lateBinder interface {
	lateBind() error
}

func lateBind(t ...interface{}) error {
	for i := range t {
		if v, ok := t[i].(lateBinder); ok {
			err := v.lateBind()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
