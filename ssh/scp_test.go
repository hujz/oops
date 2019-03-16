package ssh

import (
	"fmt"
	"os"
	"testing"
)

func TestScp(t *testing.T) {
	fmt.Println("aaaaa")
	client, err := OpenSSHSession("hujz", "123", "127.0.0.1:22")
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Stderr = os.Stdout
	client.Stdout = os.Stdout

	client.Run("scp")

}

func TestGetPerm(t *testing.T) {
	file, _ := os.Open("/opt/")
	getFileMode(*file)
	file, _ = os.Open("/opt/test/upload/intro.txt")
	getFileMode(*file)
}

func TestScpUload(t *testing.T) {
	client, err := OpenSSHSession("hujz", "123", "127.0.0.1:22")
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Stderr = os.Stdout
	client.Stdout = os.Stdout

	ScpPushFile(client, "/opt/test/haproxy/", "/opt/test/upload/")
}

func TestScpDownload(t *testing.T) {
	client, err := OpenSSHSession("hujz", "123", "127.0.0.1:22")
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	ScpPullFile(client, "/opt/test/haproxy/configuration.txt", "/opt/test/download2/")
}

func TestSSH(t *testing.T) {
	client, err := OpenSSHSession("hujz", "123", "127.0.0.1:22")
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Stderr = os.Stdout
	client.Stdout = os.Stdout
	client.Stdin = os.Stdout
	client.Run("scp -v")
}

func TestOct2Dec(t *testing.T) {
	fmt.Println(Oct2Dec("750"))
	fmt.Println(Dec2Oct(488))
	fmt.Println(fmt.Sprintf("%o", 511))
}
