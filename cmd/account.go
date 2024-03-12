/*
Copyright © 2024 Lance haig <lnhaig@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "This command is used to manage accounts on the PDS",
	Long: `This command is used for managing accounts on the PDS instance.
For example:

pdsadmin account list
pdsadmin account create --email=<EMAIL> --handle=<HANDLE>
pdsadmin account delete <DID>
pdsadmin account takedown <DID>
pdsadmin account untakedown <DID>
pdsadmin account password-reset <DID>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
