package cloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

// getEC2ID gets the IDs of all EC2 instances with the given prefix in their name
func getEC2ID() ([]*string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance.group-name"),
				Values: []*string{aws.String("EC2-Vault-SG")},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}

	result, err := svc.ec2.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("error describing instances: %v", err)
	}

	var instances []*ec2.Instance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" && strings.HasPrefix(*tag.Value, "vault-ent-node-") {
					instances = append(instances, instance)
				}
			}
		}
	}

	instanceIDs := make([]*string, len(instances))
	for i, instance := range instances {
		instanceIDs[i] = instance.InstanceId
	}
	return instanceIDs, nil
}

func waitForInstanceTermination(instanceID string) error {
	for {
		input := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(instanceID)},
		}

		result, err := svc.ec2.DescribeInstances(input)
		if err != nil {
			return err
		}

		state := result.Reservations[0].Instances[0].State.Name
		if *state == ec2.InstanceStateNameTerminated {
			break
		}

		time.Sleep(15 * time.Second) // Wait for 15 seconds before checking again
	}
	return nil
}

// terminateEC2Instance terminates the specified EC2 instance

func TerminateEC2Instance() error {
	instanceIDs, err := getEC2ID()
	if err != nil {
		return fmt.Errorf("error in getEC2ID(): %v", err)
	}
	_, err = svc.ec2.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error terminating EC2 instance: %v", err)
	}
	for _, instanceID := range instanceIDs {
		waitForInstanceTermination(*instanceID)
		fmt.Println("EC2 instance         ", *instanceID)
	}
	return nil
}

// deleteKeyPair deletes the specified key pair
func DeleteKeyPair() error {
	_, err := svc.ec2.DeleteKeyPair(&ec2.DeleteKeyPairInput{
		KeyName: aws.String("vault-EC2-kp"),
	})
	if err != nil {
		return fmt.Errorf("error deleting key pair: %v", err)
	}
	fmt.Println("Key pair              vault-EC2-kp")
	return nil
}

func DetachPolicyFromRole() error {
	_, err := svc.iam.DetachRolePolicy(&iam.DetachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonSSMAutomationRole"),
		RoleName:  aws.String("ec2-admin-role-custom"),
	})
	if err != nil {
		return fmt.Errorf("error detaching policy from role: %v", err)
	}
	return nil
}

// detachRoleFromInstanceProfile detaches the specified role from the instance profile
func DetachRoleFromInstanceProfile() error {
	_, err := svc.iam.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
		RoleName:            aws.String("ec2-admin-role-custom"),
	})
	if err != nil {
		return fmt.Errorf("error detaching role from instance profile: %v", err)
	}
	return nil
}

// deleteInstanceProfile deletes the specified instance profile
func DeleteInstanceProfile() error {
	_, err := svc.iam.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
	})
	if err != nil {
		return fmt.Errorf("error deleting instance profile: %v", err)
	}
	fmt.Println("Instance profile      ec2-InstProf-custom")
	return nil
}

// deleteRole deletes the specified IAM role
func DeleteRole() error {
	_, err := svc.iam.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String("ec2-admin-role-custom"),
	})
	if err != nil {
		return fmt.Errorf("error deleting role: %v", err)
	}
	fmt.Println("Custom role           ec2-admin-role-custom")
	return nil
}

// deleteSecurityGroup deletes the specified security group
func DeleteSecurityGroup() error {
	sgID, err := GetSGID()
	if err != nil {
		return fmt.Errorf("error in getSGID(): %v", err)
	}
	for _, id := range sgID {
		_, err = svc.ec2.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(id),
		})
		if err != nil {
			// Handle error
			fmt.Println("Error deleting security group:", err)
		}
	}
	fmt.Println("Security group       ", sgID[0])
	return nil
}
