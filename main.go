package main

import (
	"bitbucket.org/chrj/smtpd"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
)

var (
	listenAddr string
	hookFile   string
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":25", "address to listen to")
	flag.StringVar(&hookFile, "hook-file", "", "path to the hook file")
}

func main() {
	flag.Parse()

	if len(hookFile) == 0 {
		fmt.Println("a hook file (-hook-file) is required")
		os.Exit(1)
	}
	hooks, err := readHookFile(hookFile)

	if err != nil {
		fmt.Println("could not read hook file:", err)
		os.Exit(1)
	}
	l, err := net.Listen(netOfAddr(listenAddr), listenAddr)

	if err != nil {
		fmt.Println("could not setup listener:", err)
		os.Exit(1)
	}
	defer l.Close()

	NewServer(&smtpd.Server{}, hooks).Serve(l)
}

func readHookFile(path string) (Hooks, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}
	hooks := Hooks{}

	if err := yaml.Unmarshal(buf, &hooks); err != nil {
		return nil, err
	}
	return hooks, nil
}

func netOfAddr(addr string) string {
	if len(addr) > 0 && addr[0] == '/' {
		return "unix"
	}
	return "tcp"
}
