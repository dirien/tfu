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

                                          FILE                                          |                PROVIDER                | VERSION | REGISTRY VERSION | UPDATE  
----------------------------------------------------------------------------------------+----------------------------------------+---------+------------------+---------
  /Users/dirien/Tools/repos/stackit-minecraft/minecraft/main.tf                         | hashicorp/azurerm                      | 2.69.0  | 2.69.0           | false   
  /Users/dirien/Tools/repos/stackit-minecraft/minecraft/modules/minecraft-infra/main.tf | hashicorp/local                        | 2.1.0   | 2.1.0            | false   
  /Users/dirien/Tools/repos/stackit-minecraft/minecraft/modules/minecraft-infra/main.tf | terraform-provider-openstack/openstack | 1.43.0  | 1.43.0           | false   
  /Users/dirien/Tools/repos/stackit-minecraft/porter/terraform/main.tf                  | civo/civo                              | 0.10.6  | 0.10.6           | false   
  /Users/dirien/Tools/repos/stackit-minecraft/porter/terraform/main.tf                  | hashicorp/local                        | 2.1.0   | 2.1.0            | false   
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