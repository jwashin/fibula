// main.go

package main

func main() {
	a := productapp.App{}
	a.Initialize()
	a.Run(":8080")
}
