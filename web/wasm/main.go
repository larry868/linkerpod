// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

// this main package contains the web assembly source code.
// It's compiled into a '.wasm' file with "GOOS=js GOARCH=wasm go build -o ../webapp/main.wasm"
package main

import (
	"fmt"
	"os"
)

// the main func is required by the wasm GO builder
// prints will appears in the console of the browser
func main() {

	c := make(chan struct{})
	fmt.Println("Go/WASM loaded")

	// here start the code to customize
	// Welcome code
	// welcomeE := spa.GetElementById("welcome")
	// if welcomeE == nil {
	// 	panic("bad page, welcome missing")
	// }
	// welcomeE.SetInnerHTML("Loading...")

	if !ApiGetHealth() {
		fmt.Println("Go/WASM stopped")
		os.Exit(1)
	}
	// welcomeE.SetInnerHTML("Welcome Jack")

	// let's go
	fmt.Println("Go/WASM running")
	<-c
}
