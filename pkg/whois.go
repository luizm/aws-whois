package whois

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type DefaultOutput struct {
	OwnerId          *string `json:"OwnerId"`
	Description      *string `json:"Description"`
	InterfaceType    *string `json:"InterfaceType"`
	AvailabilityZone *string `json:"AvailabilityZone"`
}

type EC2Output struct {
	OwnerId      *string    `json:"OwnerId"`
	VpcId        *string    `json:"VpcId"`
	InstanceId   *string    `json:"InstanceId"`
	InstanceType *string    `json:"InstanceType"`
	Tags         []*ec2.Tag `json:"Tags"`
}

func describeInstace(sess *session.Session, i *string) (*EC2Output, error) {

	svc := ec2.New(sess)

	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{i},
	})
	if err != nil {
		return nil, err
	}

	return &EC2Output{
		OwnerId:      result.Reservations[0].OwnerId,
		VpcId:        result.Reservations[0].Instances[0].VpcId,
		InstanceId:   result.Reservations[0].Instances[0].InstanceId,
		InstanceType: result.Reservations[0].Instances[0].InstanceType,
		Tags:         result.Reservations[0].Instances[0].Tags,
	}, nil
}

func newSession(profile, region string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,

		Config: aws.Config{
			Region: aws.String(region),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})
	return sess, err
}

func describeNetworkInterface(sess *session.Session, ip string) (*ec2.DescribeNetworkInterfacesOutput, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("private-ip-address"),
				Values: []*string{
					aws.String(ip),
				},
			},
		},
	}

	result, err := svc.DescribeNetworkInterfaces(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Find(profile, region, ip string) ([]byte, error) {
	sess, err := newSession(profile, region)
	if err != nil {
		return nil, err
	}

	result, err := describeNetworkInterface(sess, ip)
	if err != nil {
		return nil, err
	}
	if result.NetworkInterfaces == nil {
		return []byte(fmt.Sprintf("Check in the %v profile, ip %v not found", profile, ip)), nil
	}

	instanceID := result.NetworkInterfaces[0].Attachment.InstanceId
	if instanceID != nil {
		r, _ := describeInstace(sess, instanceID)
		b, err := json.MarshalIndent(r, " ", " ")
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	b, err := json.MarshalIndent(&DefaultOutput{
		OwnerId:          result.NetworkInterfaces[0].OwnerId,
		Description:      result.NetworkInterfaces[0].Description,
		InterfaceType:    result.NetworkInterfaces[0].InterfaceType,
		AvailabilityZone: result.NetworkInterfaces[0].AvailabilityZone}, " ", " ")

	if err != nil {
		return nil, err
	}
	return b, nil
}
