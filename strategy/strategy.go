package main

type operation interface {
	apply(int, int) int
}

type add struct{}

func (add) apply(a, b int) int {
	return a + b
}

type subtract struct{}

func (subtract) apply(a, b int) int {
	return a - b
}

func perform(a, b int, op operation) int {
	return op.apply(a, b)
}

func main() {
	println(perform(1, 2, add{}))
	println(perform(1, 2, subtract{}))
}
