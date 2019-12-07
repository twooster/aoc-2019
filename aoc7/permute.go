package main

func permute(a []int, size int, ch chan<- []int) {
	if size == 1 {
		c := make([]int, len(a))
		copy(c, a)
		ch <- c
	} else {
		for i := 0; i < size; i++ {
			permute(a, size-1, ch)
			if size&1 == 1 {
				a[0], a[size-1] = a[size-1], a[0]
			} else {
				a[i], a[size-1] = a[size-1], a[i]
			}
		}
	}
}

func Permutate(a []int) <-chan []int {
	ch := make(chan []int)
	go func() {
		permute(a, len(a), ch)
		close(ch)
	}()
	return ch
}
