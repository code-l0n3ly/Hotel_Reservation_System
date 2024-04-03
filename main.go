// main.go
package main

import (
	"os"

	App "GraduationProject.com/m/cmd/api"
)

func main() {
	app := App.App{}
	app.Initialize("root", "wgLCfSQUYtKqCGBfviHSyMRtIloljyqm", "viaduct.proxy.rlwy.net:38199", "Hotel")
	port := os.Getenv("PORT") // Get the PORT environment variable
	if port == "" {
		port = "8080" // Default to 8080 if not specified
	}
	addr := "0.0.0.0:" + port
	app.Run(addr)
}
