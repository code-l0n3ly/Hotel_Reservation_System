// main.go
package main

import (
	App "GraduationProject.com/m/cmd/api"
)

func main() {
	app := App.App{}
	app.Initialize("root", "", "mydb")
	app.Run(":8080")
}
