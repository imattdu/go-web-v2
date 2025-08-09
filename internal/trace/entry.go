package trace

func (e *entry) K() string {
	if e == nil {
		return ""
	}
	return e.k
}

func (e *entry) V() string {
	if e == nil {
		return ""
	}
	return e.v
}
