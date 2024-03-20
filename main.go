package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/clay-codes/ec2-vault/cloud"
	"github.com/clay-codes/ec2-vault/utils"
)

var runCleanup bool
var nodes int64
var keyPath string

func init() {

	// prompt user if they want to run cleanup
	fmt.Print("Would you like to run cleanup? ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	runCleanup = strings.ToLower(response) == "yes" || strings.ToLower(response) == "y"
	if !runCleanup {
		// ask user how many instances they want to create
		fmt.Print("How many instances would you like? ")
		_, err := fmt.Scanln(&nodes)
		if err != nil {
			log.Fatal(err)
		}
	}


	if nodes > 7 {
		fmt.Printf("warning: you have entered %d number of instances.  Are you sure you want to continue? ", nodes)
		_, err := fmt.Scanln(&response)
		if err != nil {
			log.Fatal(err)
		}
		if strings.ToLower(response) != "no" && strings.ToLower(response) != "n" {
			log.Fatal("exiting")
		}
	}

	keyPath, err = utils.GetKeyPath()
	if err != nil {
		log.Fatal(err)
	}

	// authenticate with AWS
	cloud.CheckAuth()

	// creating a session
	cloud.SetRegion()
	if err := cloud.CreateSession(); err != nil {
		log.Fatal(err)
	}

	// creating needed services from session
	if err := cloud.GetSession().CreateServices(); err != nil {
		log.Fatal(err)
	}
}

// build environment
func bootStrap() {
	key, err := cloud.CreateKP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nResources created    Resource ID")
	fmt.Println("----------------------------------------")
	fmt.Println("pem file             key.pem")
	fmt.Printf("key pair             %s", key)

	sgid, err := cloud.CreateSG()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nsecurity group       %s", sgid)
	err = cloud.CreateInstProf()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nrole                 ec2-admin-role-custom")
	fmt.Println("instance profile     ec2-InstProf-custom")
	// wait for instance profile to be created sometimes necessary to avoid not found error
	time.Sleep(5 * time.Second)

	dnsMap, err := cloud.BuildEC2(nodes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n\nEnvironment nearly ready. Vault server is being installed.")
	fmt.Println("\n\nAfter a few moments, run below cmds to ssh, and be auto-authed to Vault as root:")
	for instance_name, dns := range dnsMap {
		fmt.Printf("\n%s:\n  ssh -i %s -o StrictHostKeyChecking=no ec2-user@%s\n", instance_name, keyPath, dns)
	}
}

func CleanupCloud() {
	if err := os.Remove(keyPath); err != nil {
		fmt.Printf("key.pem file may not exist: %v\n", err)
	} else {
		fmt.Println("\nResources deleted     Resource ID")
		fmt.Println("----------------------------------------")
		fmt.Println("pem file              key.pem")
	}
	if err := cloud.TerminateEC2Instance(); err != nil {
		fmt.Printf("instance may not have been created: %v\n", err)
	}
	if err := cloud.DeleteKeyPair(); err != nil {
		fmt.Printf("key pair may not exist: %v\n", err)
	}

	if err := cloud.DetachPolicyFromRole(); err != nil {
		fmt.Println(err)
	}

	if err := cloud.DetachRoleFromInstanceProfile(); err != nil {
		fmt.Printf("error detaching role from instance profile: %v\n", err)
	}
	if err := cloud.DeleteInstanceProfile(); err != nil {
		fmt.Printf("error deleting instance profile: %v\n", err)
	}
	if err := cloud.DeleteRole(); err != nil {
		fmt.Printf("error deleting role: %v\n", err)
	}
	if err := cloud.DeleteSecurityGroup(); err != nil {
		fmt.Printf("error deleting security group: %v\n", err)
	}
}

func main() {
	if runCleanup {
		CleanupCloud()
	} else {
		bootStrap()
	}
}
