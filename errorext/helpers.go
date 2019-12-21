package errorext

func ToStringSlice(errs ...error) []string {
	list := make([]string, len(errs))
	for i, err := range errs {
		list[i] = err.Error()
	}
	return list
}
