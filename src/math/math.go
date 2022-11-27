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

func RemoveFrom(elems []int, in int) []int {
	for index, value := range elems {
		if in == value {
			elems[index] = elems[len(elems)-1]
			return elems[:len(elems)-1]
		}
	}

	fmt.Println("Item not deleted")
	return elems
}
