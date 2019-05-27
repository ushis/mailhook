package main

import (
	"flag"
	"fmt"
	"github.com/chrj/smtpd"
	"net"
	"os"
)

var (
	listenAddr string
	hookDir    string
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":25", "address to listen to")
	flag.StringVar(&hookDir, "hook-dir", "", "path to look for hook files")
}

func main() {
	flag.Parse()

	if len(hookDir) == 0 {
		fmt.Println("a hook directory (-hook-dir) is required")
		os.Exit(1)
	}
	hooks, err := NewHookDir(hookDir).Hooks()

	if err != nil {
		fmt.Println("could not read hook file:", err)
		os.Exit(1)
	}
	hostname, err := os.Hostname()

	if err != nil {
		fmt.Println("could not get hostname:", err)
		os.Exit(1)
	}
	l, err := net.Listen(netOfAddr(listenAddr), listenAddr)

	if err != nil {
		fmt.Println("could not setup listener:", err)
		os.Exit(1)
	}
	defer l.Close()

	smtp := &smtpd.Server{
		Hostname:       hostname,
		WelcomeMessage: fmt.Sprintf("%s ESMTP mailhook", hostname),
	}
	NewServer(smtp, hooks).Serve(l)
}

func netOfAddr(addr string) string {
	if len(addr) > 0 && addr[0] == '/' {
		return "unix"
	}
	return "tcp"
}
