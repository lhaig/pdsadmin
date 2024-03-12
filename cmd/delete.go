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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/lhaig/pdsadmin/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <DID>",
	Short: "Delete an account",
	Long: `Delete an account by providing the did number
	For example:

pdsadmin account delete did:plc:hcf3ftdudjlrzluwune22aar.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires DID number as argument")
		}
		if internal.IsValidDid(args[0]) {
			return nil
		}
		return fmt.Errorf("ERROR: DID parameter must start with did:: yours is %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		ok := internal.PromptYesNo("Are you sure you'd like to delete "+args[0]+"?", false)
		if ok {
			fmt.Println("Deleting Account")
			deleteAccount(viper.GetString("PDS_ADMIN_PASSWORD"), args[0])
		} else {
			fmt.Println("Command Canceled?")
		}
	},
}

func init() {
	accountCmd.AddCommand(deleteCmd)
}

// deleteAccount takes an AdminPassword string and Did number and then deletes an account
func deleteAccount(AdminPassword string, Did string) {
	pdsurl := "https://" + viper.GetString("PDS_HOSTNAME") + "/xrpc/com.atproto.admin.deleteAccount"

	var accToDelete internal.DeleteDid
	accToDelete.Did = Did

	params, err := json.Marshal(accToDelete)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", pdsurl, bytes.NewBuffer(params))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("admin", AdminPassword)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var responseObject internal.CreateResponseData
	json.Unmarshal(body, &responseObject)
	fmt.Println("-----------------------------")
	fmt.Println("Account deleted successfully!")
	fmt.Println("-----------------------------")
}
