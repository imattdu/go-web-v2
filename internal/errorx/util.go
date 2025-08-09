package errorx

func Errs2Msg(errs []error) string {
	var msg string
	for _, err := range errs {
		if err == nil {
			continue
		}
		msg += err.Error()
	}
	return ""
}
