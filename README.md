# aws-whois

WIP - Just for fun

`aws-whois --profile example --ip 172.20.X.X`

```log
Checking in the profile: example
{
  "OwnerId": "xxxxxxxx",
  "Description": "Interface for NAT Gateway nat-xxxxxxxxxxxx",
  "InterfaceType": "nat_gateway",
  "AvailabilityZone": "us-east-1c"
}
```

`aws-whois --profile example --profile example2 --ip 54.X.X.X`

```log
Checking in the profile: example2
ip 54.X.X.X not found

Checking in the profile: example
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
