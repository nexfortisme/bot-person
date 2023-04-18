package logging

import (
	"log"
)

func LogError(err string) {
	log.Fatalf(err)
}
