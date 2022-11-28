package math

import "fmt"

func Max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func Min(x, y int) int {
	return x + y - Max(x, y)
}

func Contains(elems []int, i int) bool {
	for _, value := range elems {
		if i == value {
			return true
		}
	}
	return false
}
func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func RemoveFrom(elems []int, in int) []int {
	fmt.Printf("deleting: %v from %v\n", in, elems)
	for index, value := range elems {
		if in == value {
			return RemoveIndex(elems, index)
		}
	}

	fmt.Println("Item not deleted")
	return elems
}
