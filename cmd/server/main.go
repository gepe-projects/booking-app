package main

import (
	"booking/internal/bootstrap"
)

func main() {
	apps := bootstrap.InitializeApps()

	apps.Server.Run()
}
