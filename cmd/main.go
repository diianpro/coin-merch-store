package main

import (
	"diianpro/coin-merch-store/internal/app"
)

//go:generate swag init docs.go

func main() {
	app.Run()
}
