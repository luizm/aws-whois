## Description

If you have a lot of accounts, this is a simple way to check which resource on AWS are using a specific public or private IP.

## Install 

```sh 
$ brew install luizm/tap/aws-whois
```

Or download the binary from [github releases](https://github.com/luizm/aws-whois/releases)

## Usage

`$ aws-whois --profile example --ip 172.20.X.X`

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

`$ aws-whois --profile example --profile example2 --ip 54.X.X.X`

```log
{
   "Profile": "example",
   "Message": "ip 54.X.X.X not found"
}
{
   "Profile": "example2",
   "ENI": {
      "OwnerId": "xxxxxx",
      "Description": "Primary network interface",
      "InterfaceType": "interface",
      "AvailabilityZone": "us-east-1a",
      "EC2Associated": {
         "VpcId": "vpc-xxxxxx",
         "InstanceId": "i-xxxxxx",
         "InstanceType": "t3.micro",
         "Tags": [
            {
               "Key": "Name",
               "Value": "example.local"
            }
         ]
      }
   }
}
```
