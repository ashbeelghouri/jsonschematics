package main

func main() {
	RegisterValidator("IsString", IsString)
	RegisterValidator("IsInt", IsInt)
}
