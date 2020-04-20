package whois

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Result struct {
	Profile string `json:"Profile"`
	ENI     *DefaultOutput
}

type ResultNotFond struct {
	Profile string `json:"Profile"`
	Message string
}

type DefaultOutput struct {
	OwnerId          *string `json:"OwnerId"`
	Description      *string `json:"Description"`
	InterfaceType    *string `json:"InterfaceType"`
	AvailabilityZone *string `json:"AvailabilityZone"`
	EC2Associated    *EC2Output
}

type EC2Output struct {
	VpcId        *string    `json:"VpcId"`
	InstanceId   *string    `json:"InstanceId"`
	InstanceType *string    `json:"InstanceType"`
	Tags         []*ec2.Tag `json:"Tags"`
}

func isPrivateIP(ip string) (bool, error) {
	private := false
	IP := net.ParseIP(ip)
	if IP == nil {
		return false, errors.New("invalid IP address")
	}
	_, prefix08, _ := net.ParseCIDR("10.0.0.0/8")
	_, prefix12, _ := net.ParseCIDR("172.16.0.0/12")
	_, prefix16, _ := net.ParseCIDR("192.168.0.0/16")
	private = prefix08.Contains(IP) || prefix12.Contains(IP) || prefix16.Contains(IP)

	return private, nil
}

func getInstanceInfo(sess *session.Session, i *string) (*EC2Output, error) {

	svc := ec2.New(sess)

	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{i},
	})
	if err != nil {
		return nil, err
	}

	return &EC2Output{
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

func getInterfaceInfo(sess *session.Session, ip string) (*ec2.DescribeNetworkInterfacesOutput, error) {
	svc := ec2.New(sess)
	filter := "association.public-ip"

	isPrivate, err := isPrivateIP(ip)
	if err != nil {
		return nil, err
	}

	if isPrivate {
		filter = "private-ip-address"
	}

	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String(filter),
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

func FindIP(profile, region, ip string) ([]byte, error) {
	var e *EC2Output

	sess, err := newSession(profile, region)
	if err != nil {
		return nil, err
	}

	result, err := getInterfaceInfo(sess, ip)
	if err != nil {
		return nil, err
	}
	if result.NetworkInterfaces == nil {
		r := &ResultNotFond{
			Profile: profile,
			Message: fmt.Sprintf("ip %v not found", ip),
		}
		b, err := json.MarshalIndent(r, "", "   ")
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	instanceID := result.NetworkInterfaces[0].Attachment.InstanceId
	if instanceID != nil {
		e, _ = getInstanceInfo(sess, instanceID)
	}

	r := Result{
		Profile: profile,
		ENI: &DefaultOutput{
			OwnerId:          result.NetworkInterfaces[0].OwnerId,
			Description:      result.NetworkInterfaces[0].Description,
			InterfaceType:    result.NetworkInterfaces[0].InterfaceType,
			AvailabilityZone: result.NetworkInterfaces[0].AvailabilityZone,
			EC2Associated:    e,
		},
	}
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return nil, err
	}
	return b, nil
}
