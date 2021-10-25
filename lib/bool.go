package lib

func Bool(bool bool) *bool {
	return &bool
}

func UnWrapBool(b *bool) bool {
	if b == nil {
		return false
	}

	return *b
}
