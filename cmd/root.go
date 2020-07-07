/*
Copyright © 2020 Thomas Fini Hansen <xen@xen.dk>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var host string
var stream string
var topic string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "tozulip [message to send]",
	Short: "Send a message to Zulip chat",
	Long: `Send a message to Zulip chat.

Create a Zulip bot and send messages to streams from the command line. Handy
for deployment messages, or other things you might want to send from the CLI.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Message is required")
		} else if len(args) > 1 {
			return errors.New("Only one message, please")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := url.URL{
			User:   url.UserPassword(viper.GetString("mail"), viper.GetString("apikey")),
			Host:   viper.GetString("host"),
			Path:   "/api/v1/messages",
			Scheme: "https",
		}

		payload := url.Values{}
		payload.Set("type", "stream")
		payload.Set("to", viper.GetString("stream"))
		payload.Set("topic", viper.GetString("topic"))
		payload.Set("content", args[0])

		resp, err := http.Post(apiURL.String(), "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if resp.StatusCode != 200 {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			bodyString := string(bodyBytes)

			fmt.Println("Error from Zulip: ", string(bodyString))
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.tozulip.yaml)")
	rootCmd.PersistentFlags().StringP("mail", "m", "", "bot email")
	rootCmd.PersistentFlags().StringP("apikey", "k", "", "API key of bot")
	rootCmd.PersistentFlags().StringP("host", "H", "", "hostname of Zulip server")
	rootCmd.PersistentFlags().StringP("stream", "s", "", "stream to send message to")
	rootCmd.PersistentFlags().StringP("topic", "t", "", "topic of message")
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tozulip" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tozulip")
	}

	viper.SetEnvPrefix("tozulip")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
