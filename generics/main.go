package main

import "fmt"

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7}
	characters := []string{"a", "b", "c", "d", "e", "f", "g"}

	reverse(numbers)
	reverse(characters)

	fmt.Println(numbers)
	fmt.Println(characters)

	// ================

	forEach(numbers, func(e int) { fmt.Println(e) })

	// ================

}

func reverse[T any](s []T) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - 1 - i
		s[i], s[j] = s[j], s[i]
	}
}

func forEach[T any](s []T, fn func(e T)) {
	for _, e := range s {
		fn(e)
	}
}
