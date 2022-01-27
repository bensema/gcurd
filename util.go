package gcurd

func checkIn(dis []string, k string) bool {
	for i, _ := range dis {
		if dis[i] == k {
			return true
		}
	}
	return false
}
