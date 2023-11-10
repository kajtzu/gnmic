// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	gitURL  = ""
)

var downloadURL = "https://github.com/openconfig/gnmic/raw/main/install.sh"

func (a *App) VersionRun(cmd *cobra.Command, args []string) {
	if a.Config.Format != "json" {
		fmt.Printf("version : %s\n", version)
		fmt.Printf(" commit : %s\n", commit)
		fmt.Printf("   date : %s\n", date)
		fmt.Printf(" gitURL : %s\n", gitURL)
		fmt.Printf("   docs : https://gnmic.openconfig.net\n")
		return
	}
	b, err := json.Marshal(map[string]string{
		"version": version,
		"commit":  commit,
		"date":    date,
		"gitURL":  gitURL,
		"docs":    "https://gnmic.openconfig.net",
	}) // need indent? use jq
	if err != nil {
		a.Logger.Printf("failed: %v", err)
		if !a.Config.Log {
			fmt.Printf("failed: %v\n", err)
		}
		return
	}
	fmt.Println(string(b))
}

func (a *App) VersionUpgradeRun(cmd *cobra.Command, args []string) error {
	f, err := os.CreateTemp("", "gnmic")
	defer os.Remove(f.Name())
	if err != nil {
		return err
	}
	err = downloadFile(downloadURL, f)
	if err != nil {
		return err
	}

	var c *exec.Cmd
	switch a.Config.LocalFlags.UpgradeUsePkg {
	case true:
		c = exec.Command("bash", f.Name(), "--use-pkg")
	case false:
		c = exec.Command("bash", f.Name())
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		return err
	}
	return nil
}

// downloadFile will download a file from a URL and write its content to a file
func downloadFile(url string, file *os.File) error {
	client := http.Client{Timeout: 30 * time.Second}
	// Get the data
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
