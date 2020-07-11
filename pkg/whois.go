package whois

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Result struct {
	Profile string
	Message string        `json:",omitempty"`
	Result  *ResultOutput `json:",omitempty"`
}

type ResultOutput struct {
	Status           *string
	VPCId            *string
	VPCName          *string
	Description      *string
	InterfaceType    *string
	AvailabilityZone *string
	EC2Associated    *EC2Output `json:",omitempty"`
	RDSAssociated    *RDSOutput `json:",omitempty"`
}

type EC2Output struct {
	InstanceId    *string
	InstanceState *string
	LaunchTime    *time.Time
	InstanceType  *string
	Tags          []*ec2.Tag
}

type RDSEndpoint struct {
	Identifier *string
	Endpoint   *string
}

type RDSOutput struct {
	Identifier      *string
	Engine          *string
	DBInstanceClass *string
	Endpoint        *string
	Status          *string
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
		return nil, err
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

func FindIP(profile, region, ip string) (*Result, error) {
	var ec2 *EC2Output
	var rds *RDSOutput

	sess, err := newSession(profile, region)
	if err != nil {
		return &Result{}, err
	}

	eni, err := getENIInfo(sess, ip)
	if err != nil {
		return &Result{}, err
	}
	if eni.NetworkInterfaces == nil {
		result := &Result{
			Profile: profile,
			Message: fmt.Sprintf("ip %v not found", ip),
		}
		return result, nil
	}

	VPCName, err := getVPCName(sess, eni.NetworkInterfaces[0].VpcId)
	if err != nil {
		return &Result{}, err
	}

	// Getting information about EC2
	instanceID := eni.NetworkInterfaces[0].Attachment.InstanceId
	if instanceID != nil {
		ec2, _ = getEC2Info(sess, instanceID)
	}

	// Getting information about RDS
	if *eni.NetworkInterfaces[0].Description == "RDSNetworkInterface" {
		RDSEndpoints, err := getRDSEndpoints(sess)
		if err != nil {
			return &Result{}, err
		}

		for _, r := range RDSEndpoints {
			i, _ := ResolvDNS(string(*r.Endpoint))
			if ip == string(i[0]) {
				rds, _ = getRDSInfo(sess, r.Identifier)
				break
			}
		}
	}

	result := &Result{
		Profile: profile,
		Result: &ResultOutput{
			Status:           eni.NetworkInterfaces[0].Status,
			VPCId:            eni.NetworkInterfaces[0].VpcId,
			VPCName:          VPCName,
			Description:      eni.NetworkInterfaces[0].Description,
			InterfaceType:    eni.NetworkInterfaces[0].InterfaceType,
			AvailabilityZone: eni.NetworkInterfaces[0].AvailabilityZone,
			EC2Associated:    ec2,
			RDSAssociated:    rds,
		},
	}

	return result, nil
}
