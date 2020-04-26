## Description

**Still work in progress**

If you have a lot of accounts and VPCs, this is a simple way to check where or which resource on AWS are using a specific public or private IP.

## Install

```sh
$ brew install luizm/tap/aws-whois
```

Or download the binary from [github releases](https://github.com/luizm/aws-whois/releases)

## Usage

Basically use the `flag` --ip or `--dns`.

Examples:

`$ aws-whois --profile example --ip 54.X.X.X`

```log
{
  "Profile": "example",
  "ENI": {
    "OwnerId": "xxxxxx",
    "Description": "Interface for NAT Gateway nat-xxxxxx",
    "InterfaceType": "nat_gateway",
    "AvailabilityZone": "us-east-1c",
  }
  "EC2Associated": null
}
```

When use the flag `--dns` the address will be resolved before find out the resource.

`$ aws-whois --profile example --profile example2 --dns example.mydomain.local`

```log
{
   "Profile": "example",
   "Message": "ip 172.20.X.X not found"
}
{
   "Profile": "example2",
   "ENI": {
      "Status": "in-use",
      "VPCId": "vpc-xxxxx",
      "VPCName": "example",
      "Description": "Primary network interface",
      "InterfaceType": "interface",
      "AvailabilityZone": "us-east-1a",
      "EC2Associated": {
         "InstanceId": "i-xxxxx",
         "InstanceState": "running",
         "LaunchTime": "2020-04-20T19:59:40Z",
         "InstanceType": "t3.micro",
         "Tags": [
            {
               "Key": "Name",
               "Value": "example.mydomain.local"
            }
         ]
      },
      "RDSAssociated": null
   }
}
```

`$ aws-whois --profile example --profile example3 --ip 10.100.X.X`

```log
{
   "Profile": "example",
   "Message": "ip 10.100.X.X not found"
}
{
   "Profile": "example3",
   "ENI": {
      "Status": "in-use",
      "VPCId": "vpc-xxxxx",
      "VPCName": "production",
      "Description": "RDSNetworkInterface",
      "InterfaceType": "interface",
      "AvailabilityZone": "us-east-1c",
      "EC2Associated": null,
      "RDSAssociated": {
         "Identifier": "rds-example",
         "Engine": "postgres",
         "DBInstanceClass": "db.m5.large",
         "Endpoint": "rds-example.xxxxx.us-east-1.rds.amazonaws.com",
         "Status": "available"
      }
   }
}
```
