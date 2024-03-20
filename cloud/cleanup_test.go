package cloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminateEC2Instance(t *testing.T) {
	CheckAuth()
	err := TerminateEC2Instance()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
}
func TestDeleteKeyPair(t *testing.T) {
	CheckAuth()
	err := DeleteKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Remove("../key.pem"); err != nil {
		t.Fatal(err)
	}
}
func TestDeleSecurityGroup(t *testing.T) {
	CheckAuth()
	CreateSession()
	err := GetSession().CreateServices("ec2")
	if err != nil {
		t.Fatal("Error:", err)
	}

	err = DeleteSecurityGroup()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}

func TestDetachRoleFromInstanceProfile(t *testing.T) {
	CheckAuth()
	err := DetachRoleFromInstanceProfile()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
func TestDeleteInstanceProfile(t *testing.T) {
	CheckAuth()
	err := DeleteInstanceProfile()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
func TestDeleteRole(t *testing.T) {
	CheckAuth()
	err := DeleteRole()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}

func TestBootPrint(t *testing.T) {
	key:= "key"
	instanceId:= "i-1234567890abcdef0"
	tagName:= "vault-ent-node-1"
	sgid:= "sgid"
	fmt.Println("\nResources created    Resource ID")
	fmt.Println("----------------------------------------")
	fmt.Println("pem file             key.pem")
	fmt.Printf("key pair             %s", key)
	fmt.Printf("\nsecurity group       %s", sgid)
	fmt.Println("\nrole                 ec2-admin-role-custom")
	fmt.Println("instance profile     ec2-InstProf-custom")
	fmt.Printf("%s     %s\n", tagName, instanceId)
	fmt.Println("\n\nEnvironment nearly ready. Vault server is being installed on the EC2 instance.")
	fmt.Println("\n\nAfter a few moments, run below cmd to ssh.")
	fmt.Println("You will be automatically logged into Vault: ")
	fmt.Printf("\nssh -i key.pem -o StrictHostKeyChecking=no ec2-user@pubDNS")
	fmt.Println("\n\nRun this to login to Vault: ")
	fmt.Println("\nvault login -<root")
}

func TestCleanupPrint(t *testing.T) {
	instanceID:= "i-1234567890abcdef0"
	sgID:= "sg-1234567890abcdef0"
	fmt.Println("\nResources deleted     Resource ID")
	fmt.Println("----------------------------------------")
	fmt.Println("\nEC2 instance         ", instanceID)
	fmt.Println("Key pair              vault-EC2-kp")
	fmt.Println("Key file              key.pem")
	fmt.Println("Instance profile      ec2-InstProf-custom")
	fmt.Println("Custom role           ec2-admin-role-custom")
	fmt.Println("Security group       ", sgID)

	
	
}