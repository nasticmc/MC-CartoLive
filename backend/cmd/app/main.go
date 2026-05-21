package main

import (
	"fmt"
	"os"

	"meshcore-canada-live-map/backend/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
