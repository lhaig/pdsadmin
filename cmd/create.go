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
	"fmt"
	"io"
	"net/http"

	"github.com/lhaig/pdsadmin/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	emailFlag  string
	handleFlag string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create --email=<EMAIL> --handle=<HANDLE>",
	Short: "Create a new account",
	Long: `Command to create a new account on your PDS
You will need to provide an email address and a handle to create the account.
A random password will be created for the new user`,

	Run: func(cmd *cobra.Command, args []string) {
		inviteCode := getInviteCode(viper.GetString("PDS_ADMIN_PASSWORD"))
		genPassword := internal.GeneratePassword(24, true, false, true)
		NewAccount := internal.CreateAccount{
			Email:      emailFlag,
			Handle:     handleFlag,
			Password:   genPassword,
			InviteCode: inviteCode,
		}
		createNewAccount(viper.GetString("PDS_ADMIN_PASSWORD"), NewAccount)
	},
}

func init() {
	createCmd.Flags().StringVar(&emailFlag, "email", "", "Email address for the new account")
	createCmd.MarkFlagRequired("email")
	createCmd.Flags().StringVar(&handleFlag, "handle", "", "Handle for the new account")
	createCmd.MarkFlagRequired("handle")
	accountCmd.AddCommand(createCmd)
}

// getInviteCode takes an AdminPassword string and returns an invite code
func getInviteCode(AdminPassword string) string {
	url := "https://" + viper.GetString("PDS_HOSTNAME") + "/xrpc/com.atproto.server.createInviteCode"
	invRequestBody := []byte(`{"useCount": 1}`)
	postBody := bytes.NewReader(invRequestBody)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, postBody)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("admin", AdminPassword)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var responseObject internal.InviteCodeResponse
	json.Unmarshal(body, &responseObject)
	fmt.Printf("Invite Code: %v \n", responseObject.Code)

	return responseObject.Code
}

// createNewAccount takes an AdminPassword string and NewAccount object and then creates an account
func createNewAccount(AdminPassword string, NewAccount internal.CreateAccount) {
	pdsurl := "https://" + viper.GetString("PDS_HOSTNAME") + "/xrpc/com.atproto.server.createAccount"

	params, err := json.Marshal(NewAccount)
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
	fmt.Println("Account created successfully!")
	fmt.Println("-----------------------------")
	fmt.Printf("Handle   : %v \n", responseObject.Handle)
	fmt.Printf("DID      : %v \n", responseObject.Did)
	fmt.Printf("Password : %v \n", NewAccount.Password)
	fmt.Println("-----------------------------")
	fmt.Println("Save this password, it will not be displayed again.")
}
