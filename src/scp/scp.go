package scp

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

var SSH_CONNEXION_TO_TARGET *ssh.Client

func ConnectToHost(user, host string) (*ssh.Client /* , *ssh.Session"*/, error) {
	var pass string
	fmt.Print("Password: ")
	fmt.Scanf("%s\n", &pass)

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, err
	}

	// session, err := client.NewSession()
	// if err != nil {
	// 	client.Close()
	// 	return nil, nil, err
	// }

	return client /*,  session*/, nil
}
