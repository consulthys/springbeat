package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/consulthys/springbeat/beater"
)

func main() {
	err := beat.Run("springbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
