package cast2

import (
	"sort"
	"strings"
)

// InStrings 判断字符串是否在字符串数组中
func InStrings(target string, arr []string) bool {
	sort.Strings(arr)
	return InStringsSorted(target, arr)
}

// InStringsSorted 判断字符串是否在已排序的字符串数组中
func InStringsSorted(target string, arr []string) bool {
	index := sort.SearchStrings(arr, target)
	if index < len(arr) && arr[index] == target { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}

func ToSlice[T any](s, sep string) (list2 []T, err error) {
	list := strings.Split(s, sep)
	list2 = make([]T, len(list))

	for i, s2 := range list {
		value, err := To[T](s2)
		if err != nil {
			return nil, err
		}
		list2[i] = value
	}
	return list2, nil
}
