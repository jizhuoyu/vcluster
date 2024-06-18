/*
 (c) Copyright [2023-2024] Open Text.
 Licensed under the Apache License, Version 2.0 (the "License");
 You may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertica/vcluster/vclusterops"
	"github.com/vertica/vcluster/vclusterops/vlog"
)

/* CmdConfigShow
 *
 * A subcommand printing the YAML config file
 * in the default or a specified directory.
 *
 * Implements ClusterCommand interface
 */
type CmdConfigShow struct {
	sOptions vclusterops.DatabaseOptions
	CmdBase
}

func makeCmdConfigShow() *cobra.Command {
	newCmd := &CmdConfigShow{}

	cmd := makeBasicCobraCmd(
		newCmd,
		configShowSubCmd,
		"Show the content of the config file",
		`This subcommand prints the content of the config file.

Examples:
  # Show the cluster config file in the default location
  vcluster config show

  # Show the contents of the config file at /tmp/vertica_cluster.yaml
  vcluster config show --config /tmp/vertica_cluster.yaml
`,
		[]string{configFlag},
	)

	return cmd
}

func (c *CmdConfigShow) Parse(inputArgv []string, logger vlog.Printer) error {
	c.argv = inputArgv
	logger.LogArgParse(&c.argv)

	return nil
}

func (c *CmdConfigShow) Run(vcc vclusterops.ClusterCommands) error {
	fileBytes, err := os.ReadFile(dbOptions.ConfigPath)
	if err != nil {
		return fmt.Errorf("fail to read config file, details: %w", err)
	}
	fmt.Printf("%s", string(fileBytes))
	vcc.DisplayInfo("Successfully read the config file %s", dbOptions.ConfigPath)

	return nil
}

// SetDatabaseOptions will assign a vclusterops.DatabaseOptions instance
func (c *CmdConfigShow) SetDatabaseOptions(opt *vclusterops.DatabaseOptions) {
	c.sOptions = *opt
}