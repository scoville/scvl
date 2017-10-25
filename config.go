package main

import (
	"os"
)

var envProd = "production"
var envDev = "development"

var config = struct {
	env         string
	baseTplPath string
}{
	env:         envDev,
	baseTplPath: "templates",
}

func initConfig() {
	if os.Getenv("ENV") == envProd {
		config.env = envProd
	}
}
