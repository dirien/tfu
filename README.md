[![codecov](https://codecov.io/gh/dirien/tfu/branch/main/graph/badge.svg?token=ZPXEUC4NFQ)](https://codecov.io/gh/dirien/tfu)
![VEXXHOST](https://img.shields.io/badge/Terraform-7B42BC?style=for-the-badge&logo=terraform&logoColor=white)


# tfu (speak 'TF-up')

`tfu` is a Terraform helper to update the providers of [Terraform](https://registry.terraform.io/browse/providers)

Works only starting from version Terraform 0.13+

_Nothing more nothing less._

## Why? 🤷

After a Terraform session with [Nico Meisenzahl](https://github.com/nmeisenzahl), we thought about a way to
automatically update the different provisioner.

Right now you need to call the [Terraform Registry](https://registry.terraform.io/) website and look them up. That is so
time-consuming. This was moment the idea 💡 `tfu` for was born.

### TL;DR 🚀

Install via homebrew:

```bash
brew tap dirien/homebrew-dirien
brew install tfu
```

Linux or Windows user, can directly download (or use `curl`/`wget`) the binary via
the [release page](https://github.com/dirien/tfu/releases).

### Usage ⚙

So simple that downloading the cli will take you longer!

#### GitHub Token

`tfu` supports private modules hosted on github. To not run into a rate limit:

```
403 API rate limit exceeded for xxxxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.) [rate reset in 16m17s]
```

Please set the env variable:

```
export GIT_TOKEN=xxx
```

For more details on module sources -> https://www.terraform.io/docs/language/modules/sources.html#github

For directory:

```bash
tfu update -d <directory> [--dry-run]
```

For a single file

```bash
tfu update -f <path to file> [--dry-run]
```

Example:

```bash
tfu update -d /Users/dirien/Tools/repos/stackit-minecraft/

🔎 Start scanning for TF providers...  ⢿ 
🎉 Scanning finished...   

   FILE   |                                          PROVIDER (P) / MODULE (M)                                          | USED VERSION | LATEST VERSION | UPDATABLE  
----------+-------------------------------------------------------------------------------------------------------------+--------------+----------------+------------
  main.tf | git@github.com:rackspace-infrastructure-automation/aws-terraform-asg_instance_replacement//?ref=v0.12.0 (M) | 0.12.0       | 0.12.1         | true       
  main.tf | git@github.com:rackspace-infrastructure-automation/aws-terraform-vpc_basenetwork//?ref=v0.12.1 (M)          | 0.12.1       | 0.12.7         | true       
  main.tf | git@github.com:rackspace-infrastructure-automation/aws-terraform-security_group//?ref=v0.12.0 (M)           | 0.12.0       | 0.12.3         | true       
  main.tf | git@github.com:rackspace-infrastructure-automation/aws-terraform-ec2_asg//?ref=v0.12.1 (M)                  | 0.12.1       | 0.12.15        | true       
  main.tf | git@github.com:rackspace-infrastructure-automation/aws-terraform-ec2_asg//?ref=v0.12.1 (M)                  | 0.12.1       | 0.12.15        | true       
  main.tf | hashicorp/consul/aws (M)                                                                                    | 0.1.0        | 0.11.0         | true       
  main.tf | hashicorp/oci (P)                                                                                           | 4.31.0       | 4.40.0         | true  
```

### Contributing 🤝

#### Contributing via GitHub

Feel free to join.

#### License

Apache License, Version 2.0

### Libraries & Tools 🔥

- https://github.com/fatih/color
- https://github.com/go-resty/resty
- https://github.com/hashicorp/go-version
- https://github.com/hashicorp/hcl
- https://github.com/olekukonko/tablewriter
- https://github.com/spf13/cobra