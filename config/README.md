# Configuration

This directory contains the configuration files for the project.

## config.json

The `config.json` file is used to configure the AWS deployment parameters:

```json
{
  "aws_region": "eu-central-1",
  "vpc_id": "vpc-xxxx",
  "subnets": ["subnet-xxxx", "subnet-xxxx", "subnet-xxxx"],
  "key_name": "devops",
  "docker_image": "sciclon2/api:latest",
  "ami_id": "ami-xxxx"
}
```
