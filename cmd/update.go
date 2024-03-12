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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update PDS to the latest version",
	Long: `Update PDS to the latest version
	e.g. pdsadmin update.`,
	Run: func(cmd *cobra.Command, args []string) {
		updatePds()
	},
}

func init() {
	pdsCmd.AddCommand(updateCmd)
}

func areFilesEqual(file1, file2 string) (bool, error) {
	// Read the first file
	data1, err := os.ReadFile(file1)
	if err != nil {
		return false, err
	}

	// Read the second file
	data2, err := os.ReadFile(file2)
	if err != nil {
		return false, err
	}

	// Compare the content of both files
	return bytes.Equal(data1, data2), nil
}

func updatePds() {
	composeUrl := "https://raw.githubusercontent.com/bluesky-social/pds/main/compose.yaml"
	pdsDataDir := "/pds"
	composeFile := filepath.Join(pdsDataDir, "compose.yaml")
	backupComposeFile := filepath.Join(pdsDataDir, "compose.yaml.bkup")
	tmpFile := filepath.Join(os.TempDir(), "tempcompose.yaml")

	// Create the temp file
	out, err := os.Create(tmpFile)
	if err != nil {
		fmt.Printf("Error creating tempfile %v", err)
	}
	defer out.Close()

	// Get the compose file
	resp, err := http.Get(composeUrl)
	if err != nil {
		fmt.Printf("Error downloading the compose file %v", err)
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error copying the compose file %v", err)
	}

	areEqual, err := areFilesEqual(composeFile, tmpFile)
	if err != nil {
		fmt.Println("Error reading files:", err)
		os.Exit(1)
	}

	if areEqual {
		fmt.Println("PDS already up to date.")
		err = os.Remove(tmpFile)
		if err != nil {
			fmt.Printf("Error removing the temporary file %v", err)
		}
	} else {
		fmt.Println("PDS needs updating.")
		// Backup old composefile
		err = os.Rename(composeFile, backupComposeFile)
		if err != nil {
			fmt.Printf("Error renaming current compose file %v", err)
		}
		// Move new composefile
		err = os.Rename(tmpFile, composeFile)
		if err != nil {
			fmt.Printf("Error moving new compose file %v", err)
		}

		cmd := exec.Command("systemctl", "restart", "pds")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to restart PDS with error %s\n", err)
		}
		fmt.Println("PDS has been updated")
		fmt.Println("---------------------")
		fmt.Println("Check systemd logs: journalctl --unit pds")
		fmt.Println("Check container logs: docker logs pds")
	}
}
