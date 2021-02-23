package utils

func Find(arrPtr interface{}, filter func(index int, elem interface{}) bool) (int, interface{}) {
	array := arrPtr.([]interface{})
	for index, elem := range array {
		if filter(index, elem) {
			return index, elem
		}
	}
	return -1, nil
}
