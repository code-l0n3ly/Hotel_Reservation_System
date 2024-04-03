// main.go
package main

import (
	App "GraduationProject.com/m/cmd/api"
)

func main() {
	app := App.App{}
	app.Initialize("root", "wgLCfSQUYtKqCGBfviHSyMRtIloljyqm", "viaduct.proxy.rlwy.net:38199", "Hotel")
	app.Run("127.0.0.1:8080")
}
