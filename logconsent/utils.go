package logconsent

import (
	"log"
)

func Handle(err error) {
	if err != nil {
		log.Print(err)
	}
}
