// main.go

package productapp

func main() {
	a := App{}
	a.Initialize()
	a.Run(":8080")
}
