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
	"net/http"

	"github.com/lhaig/pdsadmin/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relayHosts []string

// requestCrawlCmd represents the requestCrawl command
var requestCrawlCmd = &cobra.Command{
	Use:   "request-crawl --rh=<SERVER1>",
	Short: "Request a crawl from the PDS crawlers",
	Long: `Request crawling from specific relay servers or the default PDS crawlers.
For example:

pdsadmin pds request-crawl --rh=https://relay.domain.example
pdsadmin pds request-crawl --rh="https://SERVER1,https://SERVER2"
If no flag is provided the default PDS crawlers will be notified.`,
	Run: func(cmd *cobra.Command, args []string) {
		requestCrawl(viper.GetString("PDS_HOSTNAME"), relayHosts)
	},
}

func init() {
	requestCrawlCmd.Flags().StringSliceVar(&relayHosts, "rh", []string{}, "pdsadmin pds request-crawl --rh=server1 --rh=server2")
	pdsCmd.AddCommand(requestCrawlCmd)
}

func requestCrawl(pdsHostName string, relayHosts []string) {
	var crawlRequestBody internal.CrawlHost
	if len(relayHosts) != 0 {
		// do nothing
	} else {
		relayHosts = append(relayHosts, viper.GetString("PDS_HOSTNAME"))
	}
	for _, rh := range relayHosts {
		fmt.Printf("Requesting crawl from %v \n", rh)
		crawlRequestBody.HostName = pdsHostName
		url := rh + "/xrpc/com.atproto.sync.requestCrawl"
		params, err := json.Marshal(crawlRequestBody)
		if err != nil {
			fmt.Println(err)
			return
		}
		postBody := bytes.NewReader(params)
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, url, postBody)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
		fmt.Println(res.StatusCode)
		fmt.Println(rh)

	}
}
