package util

func ContainsString(stack []string, needle string) bool {
	for _, s := range stack {
		if s == needle {
			return true
		}
	}
	return false
}

func ContainsOneOfStrings(stack []string, needles []string) bool {
	for _, n := range needles {
		for _, s := range stack {
			if s == n {
				return true
			}
		}
	}
	return false
}
