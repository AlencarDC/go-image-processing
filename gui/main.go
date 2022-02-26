package main

import "log"

func main() {

	app, err := NewApp("Photochopp v1.0", 1280, 720)
	if err != nil {
		log.Println("app: could not create new app", err)
	}

	app.Start()
}
