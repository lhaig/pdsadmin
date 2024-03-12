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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// invitecodeCmd represents the invitecode command
var createInviteCodeCmd = &cobra.Command{
	Use:   "create-invite-code",
	Short: "Create an Invite Code",
	Long: `Create an Invite Code on the PDS
e.g. pdsadmin create-invite-code.`,
	Run: func(cmd *cobra.Command, args []string) {
		inviteCode := getInviteCode(viper.GetString("PDS_ADMIN_PASSWORD"))
		fmt.Println("-----------------------------------")
		fmt.Println("Invite Code generated successfully!")
		fmt.Println(inviteCode)
		fmt.Println("-----------------------------------")
	},
}

func init() {
	accountCmd.AddCommand(createInviteCodeCmd)
}
