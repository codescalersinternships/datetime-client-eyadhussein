package main

import (
	"log"
	"time"

	datetimeclient "github.com/codescalersinternships/datetime-client-eyadhussein/pkg"
)

func main() {
	c := datetimeclient.NewRealClient("http://localhost", "8080", time.Duration(1)*time.Second)

	data, err := c.GetCurrentDateTime()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(data))
}
