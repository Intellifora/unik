package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/emc-advanced-dev/unik/pkg/config"
	"errors"
	"io/ioutil"
	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

var show bool

var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "Configure unik daemon URL for cli client commands",
	Long: `Sets the host url of the unik daemon for cli commands.
	If running unik locally, use 'unik target --host localhost'

	args:
	--host: <string, required>: host/ip address of the host running the unik daemon
	--port: <int, optional>: port the daemon is running on (default: 3000)

	--show: <bool,optional>: shows the current target that is set`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if show {
				if err := readClientConfig(); err != nil {
					return err
				}
				logrus.Infof("Current target: %s", clientConfig.Host)
				return nil
			}
			if host == "" {
				return errors.New("--host must be set for target")
			}
			if err := setClientConfig(host, port); err != nil {
				return errors.New(fmt.Sprintf("failed to save target to config file: %v", err))
			}
			logrus.Infof("target set: %s:%v", host, port)
			return nil
		}(); err != nil {
			logrus.Errorf("failed running target: %v", err)
			os.Exit(-1)
		}
	},
}

func setClientConfig(host string, port int) error {
	data, err := yaml.Marshal(config.ClientConfig{Host: fmt.Sprintf("%s:%v", host, port)})
	if err != nil {
		return errors.New("failed to convert config to yaml string: "+err.Error())
	}
	if err := ioutil.WriteFile(clientConfigFile, data, 0644); err !=nil {
		return errors.New("failed writing config to file "+ clientConfigFile +": "+err.Error())
	}
	return nil
}

func init() {
	RootCmd.AddCommand(targetCmd)
	targetCmd.Flags().BoolVar(&show, "show", false, "<bool,optional>: shows the current target that is set")
}
