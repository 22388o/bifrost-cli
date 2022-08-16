# Bifrost-cli

Bifrost-cli is a command-line interface for interacting with a BIfrost service.

## Install
```bash
$ git clone https://github.com/CashierPay/bifrost-cli.git
$ cd bifrost-cli
$ go build .
$ sudo mv bifrost-cli /usr/bin
$ bifrost-cli
```

```bash
Usage:
  bifrost-cli [command]

Available Commands:
  address     Get Bitcoin address.
  auth        Authenticate to a Bifrost services.
  balances    Show the balance sheet of all assets.
  completion  Generate the autocompletion script for the specified shell
  connect     Connect to the Bifrost services of your choice.
  help        Help about any command
  invoice     Generate a new Lightning invoice.
  sell        Sell bitcoin for fiat currency.
  tickets     List the price of bitcoin in fiat currencies.

Flags:
  -h, --help   help for bifrost-cli

Use "bifrost-cli [command] --help" for more information about a command.
```
