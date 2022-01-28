package gcurd

func CheckIn[T comparable](dis []T, k T) bool {
	for i, _ := range dis {
		if dis[i] == k {
			return true
		}
	}
	return false
}
