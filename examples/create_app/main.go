package main

import (
	"fmt"

	freebox "github.com/adongy/freebox-go"
)

func main() {
	client, err := freebox.NewClient()
	if err != nil {
		panic(err)
	}

	token, err := client.Authorize(&freebox.TokenRequestPayload{
		AppID:      "freebox-go",
		AppName:    "freebox go",
		AppVersion: "0.0.1",
		DeviceName: "fbx",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("app token:", token)
}
