package main

func decorator(f func(string)) func(string) {
	return func(s string) {
		println("Started")
		f(s)
		println("Done")
	}
}

func print(s string) {
	println(s)
}

func main() {
	// decorate the function
	decoratedPrint := decorator(print)
	decoratedPrint("Hello, World!")
}
