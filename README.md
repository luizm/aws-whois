# aws-whois

WIP - Just for fun

`aws-whois --profile example --ip 172.20.202.101`

```log
{
  "OwnerId": "xxxxxxxx",
  "Description": "Interface for NAT Gateway nat-xxxxxxxxxxxx",
  "InterfaceType": "nat_gateway",
  "AvailabilityZone": "us-east-1c"
}
```

`aws-whois --profile example --ip 172.20.202.10`

```log
{
  "OwnerId": "xxxxxxxx",
  "VpcId": "vpc-xxxxxxxx",
  "InstanceId": "i-xxxxxxxxxx",
  "InstanceType": "t3.micro",
  "Tags": [
   {
    "Key": "Name",
    "Value": "instance-example.local"
   }
  ]
}
```
