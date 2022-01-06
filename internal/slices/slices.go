package slices

func DeleteSlice(slice1 []string, slice2 []string) []string {
	var returnSlice []string
	for _, v := range slice1 {
		if !InSlice(v, slice2) {
			returnSlice = append(returnSlice, v)
		}
	}
	return returnSlice
}

func InSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func DeleteNil(slice []string) []string {
	var returnSlice []string
	for _, v := range slice {
		if v != "" {
			returnSlice = append(returnSlice, v)
		}
	}
	return returnSlice
}
