package main

import (
	"fmt"
	"github.com/port-domain/cmd/app"
	"os"
)

func main() {
	err := app.New().Execute()

	if err != nil {
		fmt.Printf("\nexit: %s\n", err)
		os.Exit(1)
	}
}
