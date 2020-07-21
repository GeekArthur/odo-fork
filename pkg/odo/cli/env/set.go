package env

import (
	"fmt"
	"strings"

	"github.com/openshift/odo/pkg/envinfo"
	"github.com/openshift/odo/pkg/log"

	"github.com/openshift/odo/pkg/odo/cli/ui"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const setCommandName = "set"

var (
	setLongDesc = ktemplates.LongDesc(`
	Set an individual value in the odo environment file
	`)

	setExample = ktemplates.Examples(`
  	# Set an individual value in the odo environment file
   	%[1]s %[2]s myNodejs
   	%[1]s %[3]s myProject
   	%[1]s %[4]s 8888
	`)
)

// SetOptions encapsulates the options for the command
type SetOptions struct {
	context    string
	cfg        *envinfo.EnvSpecificInfo
	paramName  string
	paramValue string
	forceFlag  bool
}

// NewSetOptions creates a new SetOptions instance
func NewSetOptions() *SetOptions {
	return &SetOptions{}
}

// Complete completes SetOptions after they've been created
func (o *SetOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	o.cfg, err = envinfo.NewEnvSpecificInfo(o.context)
	if err != nil {
		return errors.Wrap(err, "failed to load environment file")
	}

	o.paramName = args[0]
	o.paramValue = args[1]

	return nil
}

// Validate validates the SetOptions based on completed values
func (o *SetOptions) Validate() (err error) {
	if !o.cfg.EnvInfoFileExists() {
		return errors.Errorf("the context directory doesn't contain a component, please refer `odo create --help` to create a component")
	}

	return nil
}

// Run contains the logic for the command
func (o *SetOptions) Run() (err error) {
	if !o.forceFlag {
		if isSet := o.cfg.IsSet(o.paramName); isSet {
			if !ui.Proceed(fmt.Sprintf("%v is already set. Do you want to override it in the environment", o.paramName)) {
				log.Info("Aborted by the user")
				return nil
			}
		}
	}

	err = o.cfg.SetConfiguration(strings.ToLower(o.paramName), o.paramValue)
	if err != nil {
		return err
	}

	log.Info("Environment was successfully updated")
	return nil
}

// NewCmdSet implements the env set odo command
func NewCmdSet(name, fullName string) *cobra.Command {
	o := NewSetOptions()
	envSetCmd := &cobra.Command{
		Use:   name,
		Short: "Set a value in odo environment file",
		Long:  setLongDesc,
		Example: fmt.Sprintf(fmt.Sprint(setExample), fullName,
			envinfo.Name, envinfo.Namespace, envinfo.DebugPort),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("please provide a parameter name and value")
			} else if len(args) > 2 {
				return fmt.Errorf("only one value per parameter is allowed")
			} else {
				return nil
			}

		}, Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	envSetCmd.Flags().BoolVarP(&o.forceFlag, "force", "f", false, "Don't ask for confirmation, set the environment directly")
	envSetCmd.Flags().StringVar(&o.context, "context", "", "Use given context directory as a source for component settings")

	return envSetCmd
}
