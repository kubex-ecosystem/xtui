package cli

import (
	"fmt"
	"strings"

	gl "github.com/kubex-ecosystem/logz"
	cp "github.com/kubex-ecosystem/xtui/components"
	wp "github.com/kubex-ecosystem/xtui/wrappers"
	"github.com/spf13/cobra"
)

func AppsCmdsList() []*cobra.Command {
	return []*cobra.Command{
		InstallApplicationsCommand(),
	}
}

func InstallApplicationsCommand() *cobra.Command {
	var depList []string
	var path string
	var yes, quiet bool

	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i", "ins", "add"},
		Annotations: GetDescriptions(
			[]string{
				"Install applications and dependencies",
				"Install applications from a file or a repository and add, them to the system"},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(depList) == 0 && len(args) == 0 {
				gl.Log("error", "Empty applications list: no applications to install")
				return fmt.Errorf("no applications to install")

			}
			newArgs := []string{strings.Join(depList, " "), path, fmt.Sprintf("%t", yes), fmt.Sprintf("%t", quiet)}
			args = append(args, newArgs...)

			availableProperties := getAvailableProperties()
			if len(availableProperties) > 0 {
				adaptedArgs := adaptArgsToProperties(args, availableProperties)
				return wp.InstallDependenciesWithUI(adaptedArgs...)
			}

			// Notification: Starting installation
			cp.DisplayNotification("Starting installation of applications", "info")

			err := wp.InstallDependenciesWithUI(args...)

			if err != nil {
				// Notification: Error during installation
				cp.DisplayNotification(fmt.Sprintf("Error during installation: %s", err.Error()), "error")
				return err
			}

			// Notification: Successful installation
			cp.DisplayNotification("Applications installed successfully", "info")

			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&depList, "application", "a", []string{}, "Applications list to install")
	cmd.Flags().StringVarP(&path, "path", "p", "", "Apps installation path")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic yes to prompts")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")

	return cmd
}

func NavigateAndExecuteCommand(cmd *cobra.Command, args []string) error {
	// Detect command and its flags
	commandName := cmd.Name()
	flags := cmd.Flags()

	// Display command selection and flag definition in a form
	formConfig := createFormConfig(commandName, flags)
	formResult, err := cp.ShowFormWithNotification(formConfig)
	if err != nil {
		return err
	}

	// Set flag values based on form input
	for key, value := range formResult {
		if err := cmd.Flags().Set(key, value); err != nil {
			return err
		}
	}

	// Execute the command
	return cmd.Execute()
}
