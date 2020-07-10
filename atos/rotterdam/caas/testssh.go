package caas

import (
	"fmt"
    sshclient "github.com/helloyi/go-sshclient"    
)

func Install() {
	client, err := sshclient.DialWithPasswd("192.168.1.101:22", "vagrant", "vagrant")
	if err != nil {
	  fmt.Println(err)
	}
	
	
	// CHECK CONNECTION
	err = client.Cmd("ls -la").Run();
	if err != nil {
	  fmt.Println(err)
	}
	fmt.Println("<< Connected to 192.168.1.101:22 (vagrant)>>")


	// INSTALL MICROK8S
	fmt.Println("Executing > sudo snap install microk8s --classic --channel=1.17/stable")
	out, err := client.Cmd("sudo snap install microk8s --classic --channel=1.17/stable").Output()
	if err != nil {
	  fmt.Println(err)
	}
	fmt.Println(string(out))
	
	
	// SUDO USERMOD
	fmt.Println("Executing > sudo usermod -a -G microk8s vagrant")
	out, err = client.Cmd("sudo usermod -a -G microk8s vagrant").Output()
	if err != nil {
	  fmt.Println(err)
	}
	fmt.Println(string(out))
	
	
	// SU -
	/*fmt.Println("Executing > su - vagrant")
	err = client.Cmd("{ sleep 3; echo 'vagrant'; } | script -q -c 'su - vagrant' /dev/null").Run();
	if err != nil {
	  fmt.Println(err)
	}*/
	
	
	
	client.Close()
	
	
	
	fmt.Println("Connecting again > su - vagrant")
	client, err = sshclient.DialWithPasswd("192.168.1.101:22", "vagrant", "vagrant")
	if err != nil {
	  fmt.Println(err)
	}
	// CHECK CONNECTION
	err = client.Cmd("ls -la").Run();
	if err != nil {
	  fmt.Println(err)
	}
	fmt.Println("<< Connected to 192.168.1.101:22 (vagrant)>>")
	
	
	// run one command
	fmt.Println("Executing script ...")
	//out, err = client.Cmd("microk8s.kubectl proxy --port=8001 --address='192.168.1.101' --accept-hosts='.*'").SmartOutput();
	//if err != nil {
	//  fmt.Println(err)
	//}
	//fmt.Println(string(out))
	
	script := `
	  #!/bin/bash
	  echo '>>>>>>>>>> microk8s.kubectl proxy...'
	  microk8s.kubectl proxy --port=8001 --address='192.168.1.101' --accept-hosts='.*' &>/dev/null &
	  
	`
	out, err = client.Script(script).Output()
	if err != nil {
	  fmt.Println(err)
	}
	// the 'out' is stdout output
	fmt.Println(string(out))
	
	
	fmt.Println("<<<<<<<<<< end")
	
	
	
	defer client.Close()
}