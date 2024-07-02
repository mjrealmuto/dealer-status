package sshclient

import (
	"fmt"
	"os"
	"reflect"

	config "dealer-status/internal/config"

	"github.com/helloyi/go-sshclient"
)

type sshClient struct {
	Host string
	Port string
	User string
	
}

func SetCredentials() sshClient {
	sc := &sshClient{}
	
	attrs := reflect.TypeOf(*sc)

	for i := 0 ; i < attrs.NumField() ; i++ {
		attrName := attrs.Field(i).Name
		attrVal := config.ValidateEnvCred(attrName, "SSH")
		sc.setProperty(attrName, attrVal)
	}

	return *sc
}

func TunnelIn(sc sshClient) *sshclient.Client {

	fmt.Println("Creating the SSH Tunnel...")

	client, err := sshclient.DialWithKey(
		fmt.Sprintf("%v:%v", sc.Host, sc.Port), 
		sc.User, 
		"/root/.ssh/id_rsa",
	)

	if err != nil {
		fmt.Println("SSH Tunnel failed: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("SSH Tunnel has been created; Connecting to DB...")

	return client
}

func (s *sshClient) setProperty(name string, value string) {
	reflect.ValueOf(s).Elem().FieldByName(name).Set(reflect.ValueOf(value))
}
