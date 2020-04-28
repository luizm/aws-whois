package whois

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Result struct {
	Profile string `json`
	ENI     *DefaultOutput
}

type ResultNotFond struct {
	Profile string `json`
	Message string
}

type DefaultOutput struct {
	Status           *string `json`
	VPCId            *string `json`
	VPCName          *string `json`
	Description      *string `json`
	InterfaceType    *string `json`
	AvailabilityZone *string `json`
	EC2Associated    *EC2Output
	RDSAssociated    *RDSOutput
}

type RDSEndpoint struct {
	Identifier *string
	Endpoint   *string
}

type EC2Output struct {
	InstanceId    *string    `json`
	InstanceState *string    `json`
	LaunchTime    *time.Time `json`
	InstanceType  *string    `json`
	Tags          []*ec2.Tag `json`
}

type RDSOutput struct {
	Identifier      *string `json`
	Engine          *string `json`
	DBInstanceClass *string `json`
	Endpoint        *string `json`
	Status          *string `json`
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

func getVPCName(sess *session.Session, v *string) (*string, error) {
	svc := ec2.New(sess)
	var n *string

	r, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(*v),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, t := range r.Vpcs[0].Tags {
		if *t.Key == "Name" {
			n = t.Value
			break
		}
	}
	return n, nil
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

func getRDSEndpoints(sess *session.Session) ([]RDSEndpoint, error) {
	var e []RDSEndpoint
	svc := rds.New(sess)

	result, err := svc.DescribeDBInstances(&rds.DescribeDBInstancesInput{})
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	for _, r := range result.DBInstances {
		e = append(e, RDSEndpoint{
			Identifier: r.DBInstanceIdentifier,
			Endpoint:   r.Endpoint.Address,
		})
	}
	return e, nil
}

func getRDSInfo(sess *session.Session, r *string) (*RDSOutput, error) {
	svc := rds.New(sess)

	result, err := svc.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: r,
	})
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	return &RDSOutput{
		Identifier:      result.DBInstances[0].DBInstanceIdentifier,
		Engine:          result.DBInstances[0].Engine,
		DBInstanceClass: result.DBInstances[0].DBInstanceClass,
		Endpoint:        result.DBInstances[0].Endpoint.Address,
		Status:          result.DBInstances[0].DBInstanceStatus,
	}, nil
}

func getEC2Info(sess *session.Session, i *string) (*EC2Output, error) {
	svc := ec2.New(sess)

	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{i},
	})
	if err != nil {
		return nil, err
	}
	return &EC2Output{
		InstanceId:    result.Reservations[0].Instances[0].InstanceId,
		InstanceType:  result.Reservations[0].Instances[0].InstanceType,
		InstanceState: result.Reservations[0].Instances[0].State.Name,
		LaunchTime:    result.Reservations[0].Instances[0].LaunchTime,
		Tags:          result.Reservations[0].Instances[0].Tags,
	}, nil
}

func getENIInfo(sess *session.Session, ip string) (*ec2.DescribeNetworkInterfacesOutput, error) {
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

func ResolvDNS(dns string) ([]string, error) {
	var ips []string
	i, err := net.LookupIP(dns)
	if err != nil {
		return nil, err
	}

	for _, i := range i {
		ips = append(ips, i.String())
	}
	return ips, nil
}

func FindIP(profile, region, ip string) ([]byte, error) {
	var e *EC2Output
	var r *RDSOutput

	sess, err := newSession(profile, region)
	if err != nil {
		return nil, err
	}

	result, err := getENIInfo(sess, ip)
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

	VPCName, err := getVPCName(sess, result.NetworkInterfaces[0].VpcId)
	if err != nil {
		return nil, err
	}
	if *result.NetworkInterfaces[0].Description == "RDSNetworkInterface" {
		endpoints, err := getRDSEndpoints(sess)
		if err != nil {
			return nil, err
		}

		for _, e := range endpoints {
			i, _ := ResolvDNS(string(*e.Endpoint))
			if ip == string(i[0]) {
				r, _ = getRDSInfo(sess, e.Identifier)
				break
			}
		}
	}

	instanceID := result.NetworkInterfaces[0].Attachment.InstanceId
	if instanceID != nil {
		e, _ = getEC2Info(sess, instanceID)
	}
	o := Result{
		Profile: profile,
		ENI: &DefaultOutput{
			Status:           result.NetworkInterfaces[0].Status,
			VPCId:            result.NetworkInterfaces[0].VpcId,
			VPCName:          VPCName,
			Description:      result.NetworkInterfaces[0].Description,
			InterfaceType:    result.NetworkInterfaces[0].InterfaceType,
			AvailabilityZone: result.NetworkInterfaces[0].AvailabilityZone,
			EC2Associated:    e,
			RDSAssociated:    r,
		},
	}
	b, err := json.MarshalIndent(o, "", "   ")
	if err != nil {
		return nil, err
	}
	return b, nil
}
