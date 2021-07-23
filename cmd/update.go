package tfu

import (
	"fmt"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	Short:         "Updates the terrafrom provider",
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
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.HideCursor = true
	s.Prefix = "🔎 Start scanning for TF providers...  "
	s.FinalMSG = "\n🎉 Scanning finished...\n"
	s.Start()

	var tfProvisioner []TFProvisioner

	if len(file) > 0 {
		tfProvisioner, err = inspectFileForOutdatedProvisioner(file, tfProvisioner, dryRun)
		if err != nil {
			return err
		}
	} else if len(directory) > 0 {
		err = filepath.Walk(directory, func(path string, info os.FileInfo, fileErr error) error {
			if filepath.Ext(path) == ".tf" {
				tfProvisioner, err = inspectFileForOutdatedProvisioner(path, tfProvisioner, dryRun)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	s.Stop()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"FILE", "PROVIDER", "VERSION", "REGISTRY VERSION", "UPDATE"})
	for _, provisioner := range tfProvisioner {
		update := provisioner.OldVersion.LessThan(provisioner.NewVersion)
		table.Append([]string{provisioner.Filename, provisioner.Provider, provisioner.OldVersion.String(), provisioner.NewVersion.String(), strconv.FormatBool(update)})
	}
	table.SetBorder(false)
	table.Render()

	return nil
}

func inspectFileForOutdatedProvisioner(file string, tfProvisioner []TFProvisioner, dryRun bool) ([]TFProvisioner, error) {
	providers, err := hcl.NewHCLRequiredProvidersParser(file)
	if err != nil {
		return nil, err
	}
	tfProvisioner, err = checkProviderVersions(file, providers, dryRun, tfProvisioner)
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

func checkProviderVersions(filename string, providers *hcl.RequiredProviders, dryRun bool, tfProvisioner []TFProvisioner) ([]TFProvisioner, error) {
	for _, provider := range providers.Providers {

		registryProvider, err := registry.GetRegistryProvider(provider["source"])
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
			Provider:   provider["source"],
			OldVersion: oldV,
			NewVersion: newV,
		})

		err = updateHCLFile(filename, provider["version"], registryProvider.Version, dryRun)
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
