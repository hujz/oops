package main

import (
	"fmt"
	"os"
)

func main() {
	client, err := OpenSSHSession("hujz", "123", "127.0.0.1:22")
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Stderr = os.Stdout
	client.Stdout = os.Stdout

	client.Run("scp --help")
}
