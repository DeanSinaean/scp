package scp_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/go.crypto/ssh/agent"

	"github.com/tmc/scp"
)

func getAgent() (agent.Agent, error) {
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	return agent.NewClient(agentConn), err
}

func ExampleCopyPath() {
	f, _ := ioutil.TempFile("", "")
	fmt.Fprintln(f, "hello world")
	f.Close()
	defer os.Remove(f.Name())
	defer os.Remove(f.Name() + "-copy")

	agent, err := getAgent()
	if err != nil {
		log.Println("Failed to connect to SSH_AUTH_SOCK:", err)
		os.Exit(1)
	}

	client, err := ssh.Dial("tcp", "127.0.0.1:22", &ssh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.Signers),
		},
	})
	if err != nil {
		log.Println("Failed to dial:", err)
		os.Exit(1)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Println("Failed to create session: " + err.Error())
		os.Exit(1)
	}

	dest := f.Name() + "-copy"
	err = scp.CopyPath(f.Name(), dest, session)
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", dest)
	} else {
		fmt.Println("success")
	}
	// output:
	// success
}