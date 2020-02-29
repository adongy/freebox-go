package main

import (
	"fmt"
	"os"

	freebox "github.com/adongy/freebox-go"
)

func main() {
	client, err := freebox.NewClient(freebox.WithApp("freebox-go", os.Args[1], "0.0.1"))
	if err != nil {
		panic(err)
	}

	if err := client.Login(); err != nil {
		panic(err)
	}

	resp, err := client.ListPortForwarding()
	if err != nil {
		panic(err)
	}

	for _, conf := range resp.Result {
		fmt.Printf("%+v\n", conf)
	}
}
