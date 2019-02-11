package util

func ContainsString(array []string, search string) bool {
	for _, v := range array {
		if v == search {
			return true
		}
	}
	return false
}
