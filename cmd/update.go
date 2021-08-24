package tfu

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/dirien/tfu/pkg/git"

	"github.com/dirien/tfu/pkg/hcl"
	"github.com/dirien/tfu/pkg/registry"
	"github.com/hashicorp/go-version"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	tfuCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("directory", "d", ".", "directory to check")
	updateCmd.Flags().StringP("file", "f", "", "Single file")
	updateCmd.Flags().BoolP("dry-run", "", false, "Dry Run Update")
}

var updateCmd = &cobra.Command{
	Use:           "update",
	Short:         "Updates the terraform provider and modules",
	RunE:          runList,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runList(cmd *cobra.Command, _ []string) error {
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return errors.Wrap(err, "failed to get 'file' value.")
	}
	directory, err := cmd.Flags().GetString("directory")
	if err != nil {
		return errors.Wrap(err, "failed to get 'directory' value.")
	}
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		return errors.Wrap(err, "failed to get 'dry-run' value.")
	}
	fmt.Println("")
	if !git.CheckGitTokenIsSet() {
		fmt.Println("🔴 GIT_TOKEN not set, skipping private module checks")
	}
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.HideCursor = true
	s.Prefix = "🔎 Start scanning for TF providers...  "
	s.FinalMSG = "\n🎉 Scanning finished...\n"
	s.Start()

	var tfProvisioner []TFProvisioner
	if len(file) > 0 {
		tfProvisioner, err = inspectFileForOutdatedProvisioner(file, dryRun)
		if err != nil {
			return err
		}
	} else if len(directory) > 0 {
		err = filepath.Walk(directory, func(path string, info os.FileInfo, fileErr error) error {
			if filepath.Ext(path) == ".tf" {
				provisionersPerFile, err := inspectFileForOutdatedProvisioner(path, dryRun)
				if err != nil {
					return err
				}
				tfProvisioner = append(tfProvisioner, provisionersPerFile...)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	s.Stop()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"FILE", "PROVIDER (P) / MODULE (M)", "USED VERSION", "LATEST VERSION", "UPDATABLE"})
	table.SetAutoWrapText(false)
	for _, provisioner := range tfProvisioner {
		update := provisioner.OldVersion.LessThan(provisioner.NewVersion)
		var color []tablewriter.Colors
		if update {
			color = []tablewriter.Colors{{tablewriter.Normal}, {tablewriter.Normal}, {tablewriter.FgHiWhiteColor, tablewriter.BgRedColor}, {tablewriter.Normal, tablewriter.BgGreenColor}}
		}
		table.Rich([]string{provisioner.Filename, provisioner.Provider, provisioner.OldVersion.String(), provisioner.NewVersion.String(), strconv.FormatBool(update)}, color)
	}
	table.SetBorder(false)
	table.Render()

	return nil
}

func inspectFileForOutdatedProvisioner(file string, dryRun bool) ([]TFProvisioner, error) {
	hclFile, err := hcl.NewHCLFileParser(file)
	if err != nil {
		return nil, err
	}
	tfProvisioner, err := checkProviderVersions(file, hclFile, dryRun)
	if err != nil {
		return nil, err
	}
	return tfProvisioner, nil
}

type TFProvisioner struct {
	Filename   string
	Provider   string
	OldVersion *version.Version
	NewVersion *version.Version
}

func checkProviderVersions(filename string, tfFile *hcl.TFFile, dryRun bool) ([]TFProvisioner, error) {
	var tfProvisioner []TFProvisioner
	// check for version in module block
	if git.CheckGitTokenIsSet() {
		for _, provider := range tfFile.Module {
			repo := git.ParseGithubInfos(provider.Source)
			if repo != nil {
				latestVersion, err := repo.GetLatestVersion()
				if err != nil {
					return nil, err
				}
				if latestVersion != nil {
					tfProvisioner = append(tfProvisioner, TFProvisioner{
						Filename:   filename,
						Provider:   fmt.Sprintf("%s (M)", provider.Source),
						OldVersion: repo.Version,
						NewVersion: latestVersion,
					})
					providerSourceBase := strings.SplitAfter(provider.Source, "ref=")[0]
					newSourceBase := fmt.Sprintf("%s%s", providerSourceBase, latestVersion)
					err = updateHCLFile(filename, provider.Source, newSourceBase, dryRun)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	// check for Terraform registry modules
	for _, provider := range tfFile.Module {
		if len(provider.Version) != 0 {
			registryProvider, err := registry.GetRegistryDetails(provider.Source, registry.Modules)
			if err != nil {
				return tfProvisioner, err
			}

			oldV, err := version.NewVersion(provider.Version)
			if err != nil {
				return nil, err
			}
			newV, err := version.NewVersion(registryProvider.Version)
			if err != nil {
				return nil, err
			}

			tfProvisioner = append(tfProvisioner, TFProvisioner{
				Filename:   filename,
				Provider:   fmt.Sprintf("%s (M)", provider.Source),
				OldVersion: oldV,
				NewVersion: newV,
			})

			err = updateHCLFile(filename, fmt.Sprintf("version = \"%s\"", provider.Version), fmt.Sprintf("version = \"%s\"", registryProvider.Version), dryRun)
			if err != nil {
				return nil, err
			}
		}
	}
	// check required_provider in terraform block
	for _, provider := range tfFile.Terraform.RequiredProviders.Providers {

		registryProvider, err := registry.GetRegistryDetails(provider["source"], registry.Providers)
		if err != nil {
			return tfProvisioner, err
		}

		oldV, err := version.NewVersion(provider["version"])
		if err != nil {
			return nil, err
		}
		newV, err := version.NewVersion(registryProvider.Version)
		if err != nil {
			return nil, err
		}
		tfProvisioner = append(tfProvisioner, TFProvisioner{
			Filename:   filename,
			Provider:   fmt.Sprintf("%s (P)", provider["source"]),
			OldVersion: oldV,
			NewVersion: newV,
		})

		err = updateHCLFile(filename, fmt.Sprintf("version = \"%s\"", provider["version"]), fmt.Sprintf("version = \"%s\"", registryProvider.Version), dryRun)
		if err != nil {
			return nil, err
		}
	}
	return tfProvisioner, nil
}

func updateHCLFile(filename, oldVersion, newVersion string, dryRun bool) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	str := string(b)
	str = strings.Replace(str, oldVersion, newVersion, -1)
	if !dryRun {
		err = ioutil.WriteFile(filename, []byte(str), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
