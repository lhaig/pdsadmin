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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
	"github.com/lhaig/pdsadmin/internal"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List accounts",
	Long: `List accounts on the PDS
e.g. pdsadmin account list.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing Accounts on  " + viper.GetString("PDS_HOSTNAME"))
		listAccounts(viper.GetString("PDS_ADMIN_PASSWORD"))
	},
}

func init() {
	accountCmd.AddCommand(listCmd)
}

func listAccounts(AdminPassword string) {
	accounts := []internal.RepoDetailResponse{}
	listUrl := "https://" + viper.GetString("PDS_HOSTNAME") + "/xrpc/com.atproto.sync.listRepos?limit=100"
	method := "GET"

	listClient := &http.Client{}
	listReq, listErr := http.NewRequest(method, listUrl, nil)
	if listErr != nil {
		fmt.Println(listErr)
		return
	}
	listReq.Header.Add("Accept", "application/json")

	listRes, err := listClient.Do(listReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listRes.Body.Close()

	body, err := io.ReadAll(listRes.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var listResponseObject internal.ListResponse
	json.Unmarshal(body, &listResponseObject)

	for i := 0; i < len(listResponseObject.Repos); i++ {
		url := "https://" + viper.GetString("PDS_HOSTNAME") + "/xrpc/com.atproto.admin.getAccountInfo?did=" + listResponseObject.Repos[i].Did
		// method := "GET"

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
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
		var responseObject internal.RepoDetailResponse
		json.Unmarshal(body, &responseObject)
		accounts = append(accounts, responseObject)
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Handle", "Email", "DID")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, widget := range accounts {
		tbl.AddRow(widget.Handle, widget.Email, widget.Did)
	}

	tbl.Print()
}
