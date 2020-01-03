package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

var (
	backAccountHistoriesCsvFilepath string
	expensesCsvFilepath             string
)

var (
	logger    = log.New(os.Stdout, "", log.LstdFlags)
	errLogger = log.New(os.Stderr, "[error] ", log.LstdFlags)
)

const (
	// ExitCodeOK ok
	ExitCodeOK int = iota
	// ExitCodeError error
	ExitCodeError
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	flag.StringVar(
		&backAccountHistoriesCsvFilepath,
		"bank-histories-file",
		"",
		"Specify bank account histories csv file downloaded from moneytree.")
	flag.StringVar(
		&expensesCsvFilepath,
		"expenses-file",
		"",
		"Specify expenses csv file downloaded from moneytree.")

	os.Args = args
	flag.Parse()

	if backAccountHistoriesCsvFilepath != "" {
		bankAccountHistoriesFile, err := os.Open(backAccountHistoriesCsvFilepath)
		if err != nil {
			errLogger.Printf("Failed to open back account histories csv file. err: %s", err)
			return ExitCodeError
		}
		defer bankAccountHistoriesFile.Close()

		backAccountHistories := MoneytreeBankAccountHistories{}
		if err := gocsv.UnmarshalFile(bankAccountHistoriesFile, &backAccountHistories); err != nil {
			errLogger.Printf("Failed to unmarshal bank account histories. err: %s", err)
			return ExitCodeError
		}
		if err := saveFreeeDealsToCsvFile(backAccountHistories.toFreeeDeals(), "income_deals"); err != nil {
			errLogger.Println(err)
			return ExitCodeError
		}
	}
	if expensesCsvFilepath != "" {
		expensesFile, err := os.Open(expensesCsvFilepath)
		if err != nil {
			errLogger.Printf("Failed to open back expenses csv file. err: %s", err)
			return ExitCodeError
		}
		defer expensesFile.Close()

		expenses := MoneytreeExpenses{}
		if err := gocsv.UnmarshalFile(expensesFile, &expenses); err != nil {
			errLogger.Printf("Failed to unmarshal expenses. err: %s", err)
			return ExitCodeError
		}
		if err := saveFreeeDealsToCsvFile(expenses.toFreeeDeals(), "expense_deals"); err != nil {
			errLogger.Println(err)
			return ExitCodeError
		}
	}

	return ExitCodeOK
}

// saveFreeeDealsToCsvFile save freee deals to csv file.
func saveFreeeDealsToCsvFile(deals *FreeeDeals, filenamePrefix string) error {
	destFile, err := os.Create(fmt.Sprintf("%s_%s.csv", filenamePrefix, time.Now().Format("20060102150405")))
	if err != nil {
		return fmt.Errorf("Failed to create a file. err: %s", err)
	}
	defer destFile.Close()

	if err = gocsv.MarshalFile(deals, destFile); err != nil {
		return fmt.Errorf("Failed to marshal file. err: %s", err)
	}

	return nil
}
