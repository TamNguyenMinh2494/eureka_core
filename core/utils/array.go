package utils

import "reflect"

func Find(arrPtr interface{}, filter func(index int, elem interface{}) bool) (int, interface{}) {
	switch reflect.TypeOf(arrPtr).Kind() {
	case reflect.Slice, reflect.Ptr:
		arr := reflect.Indirect(reflect.ValueOf(arrPtr))
		if arr.Len() == 0 {
			return -1, nil
		}
		elem := reflect.Indirect(reflect.ValueOf(arrPtr))

		if elem.Kind() != arr.Index(0).Kind() {
			return -1, nil
		}
		for i := 0; i < arr.Len(); i++ {
			if filter(i, arr.Index(i).Interface()) {
				return i, arr.Index(i).Interface()
			}
		}
	}
	return -1, nil
}
