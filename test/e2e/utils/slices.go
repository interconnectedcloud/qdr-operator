// Package utils provides methods for dealing with common Slices operations
package utils

import "sort"

// ContainsAll can be used to compare two sorted slices and validate
// both are not nil and all elements from given model are present on
// the target instance.
func ContainsAll(model, target []interface{}) bool {
	if model == nil || target == nil {
		return false
	}

	if len(model) > len(target) {
		return false
	}

	ti := 0
	for _, vm := range model {
		found := false
		for i := ti; i < len(target); i++ {
			if vm == target[i] {
				found = true
				if ti == i {
					ti++
				} else {
					ti = i
				}
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// FromInts returns an interface array with sorted elements from int array
func FromInts(m []int) []interface{} {
	res := make([]interface{}, len(m))
	for i, v := range m {
		res[i] = v
	}
	sort.Slice(res, func(i1, i2 int) bool {
		return res[i1].(int) < res[i2].(int)
	})
	return res
}

// FromStrings returns an interface array with sorted elements from string array
func FromStrings(m []string) []interface{} {
	res := make([]interface{}, len(m))
	for i, v := range m {
		res[i] = v
	}
	sort.Slice(res, func(i1, i2 int) bool {
		return res[i1].(string) < res[i2].(string)
	})
	return res
}
