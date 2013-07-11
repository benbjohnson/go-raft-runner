package main

import (
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/benbjohnson/go-raft"
	"github.com/benbjohnson/go-raft-runner"
	"net"
	"net/http"
	"os"
	"path"
)

var logLevel string

func init() {
	flag.StringVar(&logLevel, "log-level", "", "log level (debug, trace)")
	raft.RegisterCommand(&runner.JoinCommand{})
}

//------------------------------------------------------------------------------
//
// Functions
//
//------------------------------------------------------------------------------

//--------------------------------------
// Main
//--------------------------------------

func main() {
	flag.Parse()
	switch logLevel {
	case "debug":
		raft.LogLevel = raft.Debug
	case "trace":
		raft.LogLevel = raft.Trace
	}

	var err error
	path := getPath()
	name, laddr := getName(path)
	
	// Setup new raft server.
	transporter := raft.NewHTTPTransporter("/raft")
	server, err := raft.NewServer(name, path, transporter, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	if err = server.Initialize(); err != nil {
		log.Fatalln(err)
	}
	
	if server.IsLogEmpty() {
		server.StartLeader()
		server.Do(&runner.JoinCommand{Name:name})
	} else {
		server.StartFollower()
	}

	// Setup HTTP server.
    mux := http.NewServeMux()
	transporter.Install(server, mux)
	http.Handle("/", mux)

	// Start server.
	fmt.Println(name)
	fmt.Println(path)
	fmt.Println("")
	log.Fatal(http.ListenAndServe(laddr, nil))
}

//--------------------------------------
// Utility
//--------------------------------------

// Retrieves the name and laddr of the server.
func getName(basepath string) (string, string) {
	var name string
	
	// Read name of server if it's already been set.
	if b, _ := ioutil.ReadFile(path.Join(basepath, "name")); len(b) > 0 {
		name = string(b)

	// Otherwise create a name based on the hostname and first available port.
	} else {
		hostname, _ := os.Hostname()
		port := 20000
		for {
			if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
				listener.Close()
				break
			} else {
				port++
			}
		}
		name = fmt.Sprintf("%s:%d", hostname, port)

		if err := ioutil.WriteFile(path.Join(basepath, "name"), []byte(name), 0644); err != nil {
			log.Fatalln(err)
		}
	}

	_, port, _ := net.SplitHostPort(name)
	return name, net.JoinHostPort("", port)
}

// Retrieves the path to save the log and configuration to. Uses the first
// parameter passed into the command line or creates a temporary directory if
// a path is not passed in.
func getPath() string {
	var path string
	if flag.NArg() == 0 {
		path, _ = ioutil.TempDir("", "go-raft-runner")
	} else {
		path = flag.Arg(0)
		if err := os.MkdirAll(path, 0744); err != nil {
			log.Fatalln(err)
		}
	}
	return path
}

