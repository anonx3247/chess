package main

func main() {
}

func check(errors ...error) {
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}
}
