package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/CashierPay/bifrost-go"
	"github.com/GrayFinance/expanduser"
	"github.com/GrayFinance/pretty-print"
	"github.com/spf13/cobra"
)

type Config struct {
	Service string `json:"services"`
	Token   string `json:"token"`
}

var (
	path string = expanduser.ExpandUser("~/.bifrost")
	file string = path + "/config.json"
)

var bfrost *bifrost.Bifrost

func init() {
	os.Mkdir(path, os.ModePerm)
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		data, err := json.Marshal(Config{})
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(file, data, 0644)
	}

	file, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	json.Unmarshal(data, &config)

	bfrost = &bifrost.Bifrost{Service: config.Service, Token: config.Token}
}

func main() {
	connect := &cobra.Command{
		Use:   "connect [URL]",
		Short: "Connect to the Bifrost services of your choice.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var config Config
			config.Service = args[0]
			config.Token = ""

			data, err := json.Marshal(config)
			if err != nil {
				log.Fatal(err)
			}
			ioutil.WriteFile(file, data, 0644)
			bfrost.Service = args[0]
		},
	}

	auth := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate to a Bifrost services.",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				username string
				password string
			)

			fmt.Print("Enter your username: ")
			fmt.Scanf("%s", &username)

			fmt.Print("Enter your password: ")
			fmt.Scanf("%s", &password)

			res, err := bfrost.Auth(username, password)
			if err != nil {
				log.Fatal(err)
			}

			config := Config{
				Service: bfrost.Service,
				Token:   res.Get("token").String(),
			}
			data, err := json.Marshal(config)
			if err != nil {
				log.Fatal(err)
			}

			ioutil.WriteFile(file, data, 0644)
			pretty.PrettyPrint("{\"message\":\"You have successfully logged in.\"}")
		},
	}

	balances := &cobra.Command{
		Use:   "balances",
		Short: "Show the balance sheet of all assets.",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			balances, err := bfrost.GetBalances()
			if err != nil {
				log.Fatal(err)
			}
			pretty.PrettyPrint(balances.Raw)
		},
	}

	invoice := &cobra.Command{
		Use:   "invoice <amount>",
		Short: "Generate a new Lightning invoice.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("You did not specify an amount.")
			}

			amount, _ := strconv.ParseFloat(args[0], 64)
			invoice, err := bfrost.CreateInvoice(amount)
			if err != nil {
				log.Fatal(err)
			}
			pretty.PrettyPrint(invoice.Raw)
		},
	}

	address := &cobra.Command{
		Use:   "address",
		Short: "Get Bitcoin address.",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			address, err := bfrost.GetAddress()
			if err != nil {
				log.Fatal(err)
			}
			pretty.PrettyPrint(address.Raw)
		},
	}

	tickets := &cobra.Command{
		Use:   "tickets",
		Short: "List the price of bitcoin in fiat currencies.",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			tickets, err := bfrost.GetTickets()
			if err != nil {
				log.Fatal(err)
			}
			pretty.PrettyPrint(tickets.Raw)
		},
	}

	sell := &cobra.Command{
		Use:   "sell <amount> <pair>",
		Short: "Sell ​​bitcoin for fiat currency.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			amount, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				log.Fatal(err)
			}

			quote := string(args[1])
			offer, err := bfrost.CreateOffer(amount, "BTC", quote, "SELL")
			if err != nil {
				log.Fatal(err)
			}
			pretty.PrettyPrint(offer.Raw)

			var confirm string
			fmt.Print("Do you want to sell ", amount, " [Y/N]: ")
			fmt.Scanf("%s", &confirm)

			if strings.ToUpper(confirm) != "Y" {
				return
			}

			offer, err = bfrost.ConfirmOffer(offer.Get("id").String())
			if err != nil {
				log.Fatal(err)
			}
			pretty.PrettyPrint(offer.Raw)
		},
	}

	root := &cobra.Command{Use: "bifrost-cli"}
	root.AddCommand(connect, auth, invoice, balances, address, tickets, sell)
	root.Execute()
}
