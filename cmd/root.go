/*
Copyright © 2023 teaf-vigoli
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// ARK holdings
var Holdings = []string{"ARKK", "ARKW", "ARKF", "ARKQ", "ARKG", "ARKX", "PRNT"}

type HoldingsReply struct {
	Symbol   string `json:"symbol"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Holdings []struct {
		Fund        string  `json:"fund"`
		Date        string  `json:"date"`
		Ticker      string  `json:"ticker"`
		Company     string  `json:"company"`
		Cusip       string  `json:"cusip"`
		Shares      int     `json:"shares"`
		MarketValue float64 `json:"market_value"`
		SharePrice  float64 `json:"share_price"`
		Weight      float64 `json:"weight"`
		WeightRank  int     `json:"weight_rank"`
	} `json:"holdings"`
}

var rootCmd = &cobra.Command{
	Use:   "goNoah",
	Short: "Query arkfunds.io to get the latest snapshot of ARK funds",
	Run: func(cmd *cobra.Command, args []string) {

		// Query holdings for the current date
		requestURL := fmt.Sprintf("https://arkfunds.io/api/v2/etf/holdings?symbol=%s", strings.Join(Holdings, ","))

		// Make the request and marshall the response back
		resp, err := http.Get(requestURL)
		if err != nil {
			fmt.Printf("Couldn't connect to endpoint")
			os.Exit(1)
		}
		defer resp.Body.Close()

		var funds HoldingsReply
		body, err_read := ioutil.ReadAll(resp.Body)
		err_unmarshall := json.Unmarshal(body, &funds)
		if err_read != nil || err_unmarshall != nil {
			fmt.Println("Error reading response body")
			os.Exit(1)
		}

		// Only print the top asset per holding for now
		var holdings_top = make(map[string]string)
		for _, holding := range funds.Holdings {
			if holding.WeightRank == 1 {
				holdings_top[holding.Fund] = holding.Ticker
			}
		}

		fmt.Println(holdings_top)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goNoah.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
