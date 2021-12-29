package main

import "sre-exporter/infra/app"

// @Title sre-exporter
// @version v0.1.x
// @Description A metadata service for changes.
// @Accept  json
// @Produce  json
func main() {
	app.Bootstrap()
}
