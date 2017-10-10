// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sync"
	"net/http"
	"time"
	"os/exec"
	"errors"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command with Cronitor monitoring",
	Long: ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("A unique monitor code and cli command are required")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		wg.Add(1)

		//var verbose Flag
		//verbose = cmd.PersistentFlags().Lookup("verbose")
		//if verbose {
		//	fmt.Println("is verbose")
		//} else {
		//	fmt.Println("is verbose")
		//}

		verbose := true
		go sendPing("run", args[0], verbose, &wg)

		wrappedCommand := exec.Command("sh", "-c", args[1])
		err := wrappedCommand.Run()

		if err == nil {
			wg.Add(1)
			go sendPing("complete", args[0], verbose, &wg)
		} else {
			fmt.Println(err)
			wg.Add(1)
			go sendPing("fail", args[0], verbose, &wg)
		}

		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(execCmd)
	RootCmd.Flags()
}

func sendPing(endpoint string, uniqueIdentifier string, verbose bool, group *sync.WaitGroup) {
	if verbose {
		fmt.Printf("Sending %s ping", endpoint)
	}

	Client := &http.Client{
		Timeout: time.Second * 3,
	}

	for i:=1; i<=5; i++  {
		_, err := Client.Get( fmt.Sprintf("https://cronitor.link/%s/%s?try=%d", uniqueIdentifier, endpoint, i))
		if err == nil {
			break
		}
	}

	group.Done()
}