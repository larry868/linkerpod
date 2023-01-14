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
	env := flag.String("env", "dev", "{file}.env environement file to load. dev by default.")
	if env != nil {
		strenv = *env
	}
	strenv = "./configs/" + strenv + ".env"

	// load environment variables
	err := godotenv.Load(strenv)
	if err != nil {
		log.Fatalf("Error loading .env variables: %s", err)
	}

	// run the web server
	spaws := spa.MakeWebserver()
	spaws.Run()

}
