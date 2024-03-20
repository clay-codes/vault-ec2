package cloud

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const EncodedUserData = "IyEvYmluL2Jhc2gKCmZ1bmN0aW9uIGluc3RhbGxfZGVwcyB7CiAgICB5dW0gaW5zdGFsbCAteSB5dW0tdXRpbHMKICAgIHl1bS1jb25maWctbWFuYWdlciAtLWFkZC1yZXBvIGh0dHBzOi8vcnBtLnJlbGVhc2VzLmhhc2hpY29ycC5jb20vQW1hem9uTGludXgvaGFzaGljb3JwLnJlcG8KICAgIHl1bSAteSBpbnN0YWxsIHZhdWx0LWVudGVycHJpc2UgPCIvZGV2L251bGwiCiAgICB5dW0gLXkgaW5zdGFsbCBqcSA8Ii9kZXYvbnVsbCIKICAgIHl1bSAteSBpbnN0YWxsIG5jIDwiL2Rldi9udWxsIgogICAgeXVtIHVwZGF0ZSAteQogICAgeXVtIGluc3RhbGwgYXdzY2xpIC15CiAgICBJTlNUQU5DRV9JRD0kKGN1cmwgLXMgaHR0cDovLzE2OS4yNTQuMTY5LjI1NC9sYXRlc3QvbWV0YS1kYXRhL2luc3RhbmNlLWlkKQogICAgUFJJVkFURV9JUD0kKGN1cmwgLXMgaHR0cDovLzE2OS4yNTQuMTY5LjI1NC9sYXRlc3QvbWV0YS1kYXRhL2xvY2FsLWlwdjQpCiAgICBSRUdJT049JChjdXJsIC1zIGh0dHA6Ly8xNjkuMjU0LjE2OS4yNTQvbGF0ZXN0L21ldGEtZGF0YS9wbGFjZW1lbnQvYXZhaWxhYmlsaXR5LXpvbmUgfCBzZWQgJ3MvW2Etel0kLy8nKQogICAgVEFHX05BTUU9JChhd3MgZWMyIGRlc2NyaWJlLXRhZ3MgLS1maWx0ZXJzICJOYW1lPXJlc291cmNlLWlkLFZhbHVlcz0kSU5TVEFOQ0VfSUQiICJOYW1lPWtleSxWYWx1ZXM9TmFtZSIgLS1yZWdpb24gJFJFR0lPTiAtLW91dHB1dCB0ZXh0IHwgY3V0IC1mNSkKICAgIGVjaG8gImV4cG9ydCBQUzE9J1skVEFHX05BTUVAJFBSSVZBVEVfSVAgXFddXCQgJyIgPj4gL2hvbWUvZWMyLXVzZXIvLmJhc2hfcHJvZmlsZQp9CgojIGluaXRpYWxpemVzIGEgc2luZ2xlIHNlcnZlciB2YXVsdCBpbnN0YW5jZSByYWZ0CmZ1bmN0aW9uIGluaXRfdmF1bHQgewogICAgZWNobyAiUEFTVEVfTElDRU5TRV9IRVJFIiA+L2V0Yy92YXVsdC5kL3ZhdWx0LmhjbGljCiAgICBjYXQgPDxFT0YxID4vZXRjL3ZhdWx0LmQvdmF1bHQuaGNsCnN0b3JhZ2UgInJhZnQiIHsKICBwYXRoICAgID0gIi9vcHQvdmF1bHQvZGF0YSIKICBub2RlX2lkID0gIiQoaG9zdG5hbWUpIgp9CgpsaXN0ZW5lciAidGNwIiB7CiAgYWRkcmVzcyAgICAgICAgID0gIjAuMC4wLjA6ODIwMCIKICB0bHNfZGlzYWJsZSAgICAgPSB0cnVlCn0KCmxpY2Vuc2VfcGF0aCA9ICIvZXRjL3ZhdWx0LmQvdmF1bHQuaGNsaWMiCmFwaV9hZGRyID0gImh0dHA6Ly8kKGN1cmwgLXMgaHR0cDovLzE2OS4yNTQuMTY5LjI1NC9sYXRlc3QvbWV0YS1kYXRhL2xvY2FsLWlwdjQpOjgyMDAiCmNsdXN0ZXJfYWRkciA9ICJodHRwOi8vJChob3N0bmFtZSk6ODIwMSIKbG9nX2xldmVsID0gInRyYWNlIgpFT0YxCiAgICBlY2hvICdleHBvcnQgVkFVTFRfQUREUj1odHRwOi8vMTI3LjAuMC4xOjgyMDAnID4+L2V0Yy9lbnZpcm9ubWVudAogICAgZWNobyAiZXhwb3J0IEFXU19ERUZBVUxUX1JFR0lPTj0kKGN1cmwgLXMgaHR0cDovLzE2OS4yNTQuMTY5LjI1NC9sYXRlc3QvZHluYW1pYy9pbnN0YW5jZS1pZGVudGl0eS9kb2N1bWVudCB8IC91c3IvYmluL2pxIC1yICcucmVnaW9uJykiID4+L2V0Yy9lbnZpcm9ubWVudAogICAgZXhwb3J0IFZBVUxUX0FERFI9aHR0cDovLzEyNy4wLjAuMTo4MjAwCiAgICBzeXN0ZW1jdGwgc3RhcnQgdmF1bHQKICAgIHZhdWx0IG9wZXJhdG9yIGluaXQgLWtleS1zaGFyZXM9MSAta2V5LXRocmVzaG9sZD0xID4vaG9tZS9lYzItdXNlci9rZXlzCiAgICBlY2hvICQoZ3JlcCAnS2V5IDE6JyAvaG9tZS9lYzItdXNlci9rZXlzIHwgYXdrICd7cHJpbnQgJE5GfScpID4vaG9tZS9lYzItdXNlci91bnNlYWwKICAgIHZhdWx0IG9wZXJhdG9yIHVuc2VhbCAkKGNhdCBob21lL2VjMi11c2VyL3Vuc2VhbCkKICAgIGVjaG8gJChncmVwICdJbml0aWFsIFJvb3QgVG9rZW46JyAvaG9tZS9lYzItdXNlci9rZXlzIHwgYXdrICd7cHJpbnQgJE5GfScpID4vaG9tZS9lYzItdXNlci9yb290CiAgICBybSAvaG9tZS9lYzItdXNlci9rZXlzCiAgICBjYXQgPDxFT0YyID4+L2hvbWUvZWMyLXVzZXIvLmJhc2hfcHJvZmlsZQpmdW5jdGlvbiBsb2dpbiAoKSB7CiAgICB2YXVsdCBsb2dpbiAtPHJvb3QKfQpsb2dpbgpFT0YyCn0KCmluc3RhbGxfZGVwcwppbml0X3ZhdWx0"
func GetImgID() (string, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String("/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"),
	}

	result, err := svc.ssm.GetParameter(input)
	//aws-specific error library https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/handling-errors.html
	if err != nil {
		return "", err
	}

	return *result.Parameter.Value, nil
}

func CreateKP() (string, error) {
	//if key pair already exists, return the key name
	descResult, err := svc.ec2.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{
		KeyNames: []*string{
			aws.String("vault-EC2-kp"),
		},
	})
	if err == nil {
		return *descResult.KeyPairs[0].KeyName, nil
	}
	// Create the key pair
	input := &ec2.CreateKeyPairInput{
		KeyName: aws.String("vault-EC2-kp"),
		KeyType: aws.String("rsa"),
	}

	result, err := svc.ec2.CreateKeyPair(input)
	if err != nil {
		return "", fmt.Errorf("error creating key pair: %w", err)
	}

	// Write the key material to a file
	file, err := os.Create("key.pem")
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// write key material to file
	_, err = file.WriteString(*result.KeyMaterial)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	// modify key.pem permissions to be read-only
	if err = os.Chmod("key.pem", 0400); err != nil {
		return "", fmt.Errorf("error changing file permissions: %w", err)
	}
	return *result.KeyName, nil
}

func GetVPC() (string, error) {
	vpcs, err := svc.ec2.DescribeVpcs(nil)
	if err != nil {
		return "", fmt.Errorf("error when calling ec2.DescribeVpcs: %w", err)
	}

	// Select the first VPC
	vpcID := vpcs.Vpcs[0].VpcId

	return *vpcID, nil
}

// checkAndReturnSG checks if the security group already exists and returns it
func checkAndReturnSG() (string, error) {
	sgid, err := GetSGID()
	if err != nil {
		return "", fmt.Errorf("error getting security group ID: %v", err)
	}

	if len(sgid) == 1 {
		return sgid[0], nil
	}

	if len(sgid) > 1 {
		fmt.Printf("more than one security group found with the name vault-EC2-sg, using %v", sgid[0])
		fmt.Printf("consider inspecting, then deleting the rest: %v", sgid[1:])
		fmt.Printf("can inspect via: aws ec2 describe-security-groups --group-id <sgid>")
		fmt.Printf("delete via AWS CLI with the following command: aws ec2 delete-security-group --group-id <sgid>")
		return sgid[0], nil
	}

	return "", nil
}

// CreateSG creates a security group and authorizes all inbound traffic
func CreateSG() (string, error) {
	// Check if the security group already exists
	sg, err := checkAndReturnSG()
	if err != nil {
		return "", fmt.Errorf("error checking security group: %v", err)
	}

	// If the security group exists, return its ID
	if sg != "" {
		return sg, nil
	}

	vpcID, err := GetVPC()
	if err != nil {
		return "", fmt.Errorf("error getting VPC ID: %v", err)
	}

	// Define the security group parameters
	createSGInput := &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String("EC2-Vault-SG"),
		Description: aws.String("sg for vault instance"),
		VpcId:       aws.String(vpcID), // Replace with your VPC ID
	}

	createSGOutput, err := svc.ec2.CreateSecurityGroup(createSGInput)
	if err != nil {
		return "", fmt.Errorf("error creating security group: %v", err)
	}

	// Authorize all inbound traffic
	authorizeIngressInput := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: createSGOutput.GroupId,
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("-1"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("for ec2-vault-sg-ingress"),
					},
				},
			},
		},
	}

	_, err = svc.ec2.AuthorizeSecurityGroupIngress(authorizeIngressInput)
	if err != nil {
		return "", fmt.Errorf("error authorizing security group ingress: %v", err)
	}
	// NOTE: AWS ALREADY HAS A DEFAULT EGRESS RULE ALLOWING ALL TRAFFIC, SO NO NEED TO AUTHORIZE EGRESS
	return *createSGOutput.GroupId, nil
}

func GetSGID() ([]string, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String("EC2-Vault-SG")},
			},
		},
	}

	result, err := svc.ec2.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	var groupIds []string
	for _, group := range result.SecurityGroups {
		groupIds = append(groupIds, *group.GroupId)
	}

	return groupIds, nil
}

func CreateInstProf() error {
	// Check if the instance profile already exists
	_, err := svc.iam.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
	})

	// If the instance profile exists, return its ARN
	if err == nil {
		fmt.Println("\ninstance profile `ec2-InstProf-custom` already exists, using it.")
		return err
	}

	// Define the trust policy document
	policyDocument := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "ec2.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`

	// Create the role
	createRoleInput := &iam.CreateRoleInput{
		RoleName:                 aws.String("ec2-admin-role-custom"),
		AssumeRolePolicyDocument: aws.String(policyDocument),
	}

	// check icf role already exists
	// _, err = svc.iam.GetRole(&iam.GetRoleInput{
	// 	RoleName: aws.String("ec2-admin-role-custom"),
	// })
	// if err == nil {
	// 	fmt.Println("role already exists, using 'ec2-admin-role-custom'")

	// } else {}

	_, err = svc.iam.CreateRole(createRoleInput)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	input := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonSSMAutomationRole"),
		RoleName:  aws.String("ec2-admin-role-custom"),
	}

	_, err = svc.iam.AttachRolePolicy(input)
	if err != nil {
		return err
	}
	// Create the instance profile
	createInstanceProfileInput := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
	}
	_, err = svc.iam.CreateInstanceProfile(createInstanceProfileInput)
	if err != nil {
		return fmt.Errorf("error creating instance profile: %w", err)
	}

	// wait for instance profile to be created
	time.Sleep(5 * time.Second)

	// Attach the role to the instance profile
	addRoleToInstanceProfileInput := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
		RoleName:            aws.String("ec2-admin-role-custom"),
	}
	_, err = svc.iam.AddRoleToInstanceProfile(addRoleToInstanceProfileInput)
	if err != nil {
		return fmt.Errorf("error adding role to instance profile: %w", err)
	}

	return nil
}

func GetSubnetID() (string, error) {
	vpcID, err := GetVPC()
	if err != nil {
		return "", err
	}
	// Describe subnets with the specified VPC ID
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := svc.ec2.DescribeSubnets(input)
	if err != nil {
		return "", fmt.Errorf("error describing subnets: %w", err)
	}
	// Check if there is at least one subnet and get its ID
	if len(result.Subnets) == 0 {
		return "", fmt.Errorf("no subnets found for given VPC ID: %s", vpcID)
	}
	return *result.Subnets[0].SubnetId, nil
}

func GetEC2ID() (string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String("vault-ent-server-raft")},
			},
		},
	}

	result, err := svc.ec2.DescribeInstances(input)
	if err != nil {
		return "", fmt.Errorf("error describing instances: %v", err)
	}

	return *result.Reservations[0].Instances[0].InstanceId, nil
}

func GetPublicDNS(instanceID *string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			instanceID,
		},
	}

	result, err := svc.ec2.DescribeInstances(input)
	if err != nil {
		return "", fmt.Errorf("error describing instances: %v", err)
	}

	return *result.Reservations[0].Instances[0].PublicDnsName, nil
}

//	use this function in developemnt to customize the user-data.yaml file
//
// `go build` will not compile the user-data.yaml file to binary--must be hardcoded
//
//	func encodeUserData() (string, error) {
//		// Read user data from file
//		userData, err := os.ReadFile("user-data.yaml")
//		if err != nil {
//			return nil, err
//		}
//		return base64.StdEncoding.EncodeToString(userData), nil
//	}
func BuildEC2(nodes int64) (map[string]string, error) {
	encodedUserData := EncodedUserData
	// can use instead of EncodedUserData if you would like to customize the user-data.yaml file =
	// encodedUserData, err := encodeUserData()
	// if err != nil {
	// 	return nil, fmt.Errorf("error encoding user data: %v", err)
	// }

	imageID, err := GetImgID()
	if err != nil {
		return nil, fmt.Errorf("error getting image ID: %v", err)
	}

	sgID, err := GetSGID()
	if err != nil {
		return nil, fmt.Errorf("error getting security group ID: %v", err)
	}

	subnetID, err := GetSubnetID()
	if err != nil {
		return nil, fmt.Errorf("error getting subnet ID: %v", err)
	}

	input := &ec2.RunInstancesInput{
		ImageId:          aws.String(imageID),
		InstanceType:     aws.String("t2.micro"),
		KeyName:          aws.String("vault-EC2-kp"),
		SecurityGroupIds: aws.StringSlice(sgID),
		SubnetId:         aws.String(subnetID),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String("ec2-InstProf-custom"),
		},
		UserData: aws.String(encodedUserData),
		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(nodes),
	}
	result, err := svc.ec2.RunInstances(input)
	if err != nil {
		return nil, fmt.Errorf("error running instances: %v", err)
	}

	instanceIds := make([]*string, nodes)
	for i, instance := range result.Instances {
		instanceIds[i] = instance.InstanceId
	}

	// assign unique "Name" tags to each instance
	for i, instanceId := range instanceIds {
		tagName := fmt.Sprintf("vault-ent-node-%d", i+1) // Unique name for each instance
		_, err := svc.ec2.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{instanceId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(tagName),
				},
			},
		})
		if err != nil {
			fmt.Printf("error tagging instance %s: %s\n", *instanceId, err)
			continue // Attempt to tag the next instance
		}
		fmt.Printf("%v     %v\n", tagName, *instanceId)
	}

	err = svc.ec2.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		return nil, fmt.Errorf("error waiting for instance to run: %v", err)
	}

	info, err := svc.ec2.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		return nil, fmt.Errorf("error describing instances: %v", err)
	}

	// Create a map to hold the instance names and their public DNS addresses
	instMap := make(map[string]string)

	// for range of instanceIds, get the public DNS and name of each instance
	for _, reservation := range info.Reservations {
		for _, instance := range reservation.Instances {
			pubDNS := *instance.PublicDnsName
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
			instMap[name] = pubDNS
		}
	}

	return instMap, nil
}
