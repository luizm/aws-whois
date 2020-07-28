[![Go Report Card](https://goreportcard.com/badge/github.com/luizm/aws-whois)](https://goreportcard.com/report/github.com/luizm/aws-whois)
[![codecov](https://codecov.io/gh/luizm/aws-whois/branch/master/graph/badge.svg)](https://codecov.io/gh/luizm/aws-whois)

## Description

If you have a lot of accounts and VPCs, this is a simple way to check based on DNS or IP, where is it running, and which resource on AWS is it.

## Install

```sh
$ brew install luizm/tap/aws-whois
```

Or download the binary from [github releases](https://github.com/luizm/aws-whois/releases)

## Usage

By default, all AWS profiles configured will be used to verify a specific address.

Parameter available:

```
GLOBAL OPTIONS:
   --region value, -r value          the region to use. Overrides config/env settings. (default: "us-east-1")
   --profile value, -p value         use a specific profile from your credential file.
   --ignore-profile value, -i value  ignore a specific profile from your credential file, can be used multiple times.
   --address value, -a value         the ip or dns address to find the resource associated
   --help, -h                        show help (default: false)
   --version, -v                     print the version (default: false)
```

Example:

`$ aws-whois --address 54.X.X.X`

```log
{
   "Profile": "example1",
   "Message": "ip 54.X.X.X not found"
}
{
   "Profile": "example2",
   "Message": "ip 54.X.X.X not found"
}
{
   "Profile": "example3",
   "Result": {
      "Status": "in-use",
      "VPCId": "vpc-xxxxx",
      "VPCName": "vpc-example",
      "Description": "Primary network interface",
      "InterfaceType": "interface",
      "AvailabilityZone": "us-east-1a",
      "EC2Associated": {
         "InstanceId": "i-xxxxx",
         "InstanceState": "running",
         "LaunchTime": "2020-03-14T19:59:40Z",
         "InstanceType": "t3.micro",
         "Tags": [
            {
               "Key": "Name",
               "Value": "intance-example3.local"
            }
         ]
      }
   }
}
```
