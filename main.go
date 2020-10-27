package main

import (
	"fmt"
	"log"
	"net/http"

	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/mrcoding"
)

func main() {
	// Get config at ./config.yaml
	c, err := config.New()
	if err != nil {
		log.Fatalf("Read config.yaml file fail, err: %v", err)
	}

	bot, err := mrcoding.New(c)
	if err != nil {
		log.Fatalf("Create linebot fail, err: %v", err)
	}

	http.HandleFunc("/callback", mrcoding.Handler(bot))

	err = http.ListenAndServe(":1225", nil)
	if err != nil {
		log.Fatalf("Listen and serve fail, err: %v", err)
	}

	fmt.Println("serve on :1225")
}
