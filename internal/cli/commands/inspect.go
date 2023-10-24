package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	// inspectCommand is the command for inspecting your Capital API account
	//
	// example curl:
	//   curl --location 'https://api-capital.backend-capital.com/api/v1/session' \
	//     --header 'X-CAP-API-KEY: PNSk969vsStoqHTU' \
	//     --header 'Content-Type: application/json'
	//     --data '{ "identifier": "apostolis.anastasiou.alpha@gmail.com", "password": "REu8aZW89p@$Fs#o" }'
	//
	// No, this is not my actual password, don't even try it idiot.
	//
	// example http response:
	//    {
	//       "accountType":"CFD",
	//       "accountInfo":{"balance":6609.28,"deposit":7151.98,"profitLoss":-542.70,"available":1.34},
	//       "currencyIsoCode":"USD","currencySymbol":"$","currentAccountId":"178572855993259204","streamingHost":"wss://api-streaming-capital.backend-capital.com/","accounts":[{"accountId":"178572855993259204","accountName":"","preferred":true,"accountType":"CFD"}],
	//       "clientId":"24984807",
	//       "timezoneOffset":3,
	//       "hasActiveDemoAccounts":true,
	//       "hasActiveLiveAccounts":true,
	//       "trailingStopsEnabled":true
	//    }
	//
	// example usage:
	//
	//    # With pass on terminal
	//    capcli inspect --email info@example --api-pass 1234567890abcdef
	//
	//    # With pass from stdin - hidden
	//    capcli inspect --email info@example
	inspectCommand = &cobra.Command{
		Use:     "inspect",
		Short:   "Inspect capital cli instruments",
		Long:    "Inspect capital cli instruments",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				panic(err)
				return
			}
		},
	}
	inspectAccount = &cobra.Command{
		Use:   "account",
		Short: "Inspect your Capital API account",
		Long:  "Inspect your Capital API account",
		ValidArgs: []string{
			"email",
			"api-key",
			"api-pass",
		},
		ArgAliases: []string{
			"e",
			"p",
		},
		Example: "capcli inspect account --email info@example.com --api-pass 1234567890abcdef",
		Run: func(cmd *cobra.Command, args []string) {
			var emailStr string
			var apiKeyStr string
			var apiPassStr string
			var err error

			email := cmd.Flag("email").Value.String()
			apiKey := cmd.Flag("api-key").Value.String()
			apiPass := cmd.Flag("api-pass").Value.String()

			if email == "" {
				cmd.PrintErrln(errors.New("email is required"))
				return
			}
			emailStr = email

			if apiKey == "" {
				cmd.PrintErrln(errors.New("api-key is required"))
				return
			}
			apiKeyStr = apiKey

			if apiPass == "" {
				cmd.Print("Enter API Pass: ")
				var apiPassBytes []byte

				apiPassBytes, err = terminal.ReadPassword(0)
				if err != nil {
					cmd.PrintErrln(fmt.Errorf("could not read password: %w", err))
					return
				} // This is to clear the password from the terminal
				apiPassStr = string(apiPassBytes)

				apiPassStr = strings.TrimSpace(apiPassStr)
			} else {
				apiPassStr = apiPass
			}

			httpBodyStr := fmt.Sprintf(`{ "identifier": "%s", "password": "%s" }`, emailStr, apiPassStr)
			httpBodyReader := strings.NewReader(httpBodyStr)

			httpReq, err := http.NewRequest(
				"POST",
				"https://api-capital.backend-capital.com/api/v1/session",
				httpBodyReader,
			)
			if err != nil {
				cmd.PrintErrln(fmt.Errorf("could not create http request: %w", err))
				return
			}

			httpReq.Header.Add("X-CAP-API-KEY", apiKeyStr)
			httpReq.Header.Add("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(httpReq)
			if err != nil {
				cmd.PrintErrln(fmt.Errorf("could not execute http request: %w", err))
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				cmd.PrintErrln(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					cmd.PrintErrln(fmt.Errorf("could not read response body (which had an error response): %w", err))
					return
				}

				bodyStr := string(bodyBytes)

				var bodyJsonMarshalled []byte
				bodyJsonMarshalled, err = json.MarshalIndent(bodyStr, "  ", "  ")

				var bodyJsonStr string
				bodyJsonStr = string(bodyJsonMarshalled)
				bodyJsonStr = strings.TrimSpace(bodyJsonStr)
				bodyJsonStr, err = strconv.Unquote(bodyJsonStr)
				if err != nil {
					cmd.PrintErrln(fmt.Errorf("could not unquote response body (which had an error response): %w", err))
					return
				}

				if bodyJsonStr == "{\"errorCode\":\"error.invalid.details\"}" {
					cmd.PrintErrln("Invalid email or password")
					return
				}

				cmd.Printf("Response error, API Response: \n%s\n", bodyJsonStr)
				return
			}

			type responseJsonType struct {
				AccountType string `json:"accountType"`
				AccountInfo struct {
					Balance    float64 `json:"balance"`
					Deposit    float64 `json:"deposit"`
					ProfitLoss float64 `json:"profitLoss"`
					Available  float64 `json:"available"`
				} `json:"accountInfo"`
				CurrencyIsoCode  string `json:"currencyIsoCode"`
				CurrencySymbol   string `json:"currencySymbol"`
				CurrentAccountId string `json:"currentAccountId"`
				StreamingHost    string `json:"streamingHost"`
				Accounts         []struct {
					AccountId   string `json:"accountId"`
					AccountName string `json:"accountName"`
					Preferred   bool   `json:"preferred"`
					AccountType string `json:"accountType"`
				} `json:"accounts"`
				ClientId              string `json:"clientId"`
				TimezoneOffset        int    `json:"timezoneOffset"`
				HasActiveDemoAccounts bool   `json:"hasActiveDemoAccounts"`
				HasActiveLiveAccounts bool   `json:"hasActiveLiveAccounts"`
				TrailingStopsEnabled  bool   `json:"trailingStopsEnabled"`
			}

			responseJson := &responseJsonType{}

			err = json.NewDecoder(resp.Body).Decode(responseJson)
			if err != nil {
				cmd.PrintErrln(fmt.Errorf("could not decode json: %w", err))
				return
			}

			cmd.Println("Account Type:", responseJson.AccountType)
			cmd.Println()
			cmd.Println("Account Info:")
			cmd.Println("  Balance:", responseJson.AccountInfo.Balance)
			cmd.Println("  Deposit:", responseJson.AccountInfo.Deposit)
			cmd.Println("  Profit Loss:", responseJson.AccountInfo.ProfitLoss)
			cmd.Println("  Available:", responseJson.AccountInfo.Available)
			cmd.Println()
			cmd.Println("Currency Iso Code:", responseJson.CurrencyIsoCode)
			cmd.Println()
			cmd.Println("Currency Symbol:", responseJson.CurrencySymbol)
			cmd.Println()
			cmd.Println("Current Account Id:", responseJson.CurrentAccountId)
			cmd.Println()
			cmd.Println("Streaming Host:", responseJson.StreamingHost)
			cmd.Println()
			cmd.Println("Accounts:", responseJson.Accounts)
			cmd.Println()
			cmd.Println("Client Id:", responseJson.ClientId)
			cmd.Println()
			cmd.Println("Timezone Offset:", responseJson.TimezoneOffset)
			cmd.Println()
			cmd.Println("Has Active Demo Accounts:", responseJson.HasActiveDemoAccounts)
			cmd.Println()
			cmd.Println("Has Active Live Accounts:", responseJson.HasActiveLiveAccounts)
			cmd.Println()
			cmd.Println("Trailing Stops Enabled:", responseJson.TrailingStopsEnabled)
		},
	}
)

func init() {
	inspectCommand.AddCommand(inspectAccount)

	inspectAccount.Flags().
		StringP("email", "e", "", "Your Capital API email")
	inspectAccount.Flags().
		StringP("api-key", "k", "", "Your Capital API key")
	inspectAccount.Flags().
		StringP("api-pass", "p", "", "Your Capital API password")
}
