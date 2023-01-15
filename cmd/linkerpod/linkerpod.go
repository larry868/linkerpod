// linkerpod CLI
//
// Run the werserver
package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/sunraylab/linkerpod/pkg/spa"
)

func main() {
	// get --env flag
	strenv := "dev"
	env := flag.String("env", "dev", ".env environement file to load, with the path and without the extension. dev by default.")
	flag.Parse()
	if env != nil {
		strenv = *env
	}

	// load environment variables
	err := godotenv.Load(strenv + ".env")
	if err != nil {
		log.Fatalf("Error loading .env variables: %s", err)
	}

	// Make a web server a add APIs route handlers
	spaws := spa.MakeWebserver()
	//spaws.ApiRouter.HandleFunc("/login", api.ServeLogin())

	// Let's start the server and listen requests
	spaws.Run()
}
