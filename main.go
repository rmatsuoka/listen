package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

var (
	netflag = flag.String("net", "tcp", "network")
	aflag   = flag.String("a", "localhost:7000", "address")
	eflag   = flag.Bool("e", false, "send stderr of cmd to the client")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-net network] [-a addr] [-e] cmd [args...]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)
	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() == 0 {
		Usage()
		os.Exit(1)
	}

	if _, err := exec.LookPath(flag.Arg(0)); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen(*netflag, *aflag)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		go listen(conn, flag.Arg(0), flag.Args()[1:]...)
	}
}

func listen(conn net.Conn, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdin = conn
	c.Stdout = conn
	if *eflag {
		c.Stderr = conn
	}
	if err := c.Start(); err != nil {
		log.Print(err)
	}

	c.Wait()
}