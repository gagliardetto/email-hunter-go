# email-hunter-go

[![GoDoc](https://godoc.org/github.com/gagliardetto/email-hunter-go?status.svg)](https://godoc.org/github.com/gagliardetto/email-hunter-go)
[![GitHub license](https://img.shields.io/github/license/gagliardetto/email-hunter-go.svg)](https://github.com/gagliardetto/email-hunter-go/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gagliardetto/email-hunter-go)](https://goreportcard.com/report/github.com/gagliardetto/email-hunter-go)

## Description

Go wrapper for Email Hunter API. Unofficial. Pre-Alpha.

## How to get an API key

You can get an API key from https://emailhunter.co/api_keys

## Installation

Run the following command to install the package:

```
go get -u github.com/gagliardetto/email-hunter-go
```

## Examples

#### Get API info

```go
package main

import (
	"fmt"

	"github.com/gagliardetto/email-hunter-go"
)

var APIKey string = "[API_KEY]"

func main() {
	client, err := emailHunter.NewClient(APIKey)
	if err != nil {
		fmt.Println("client setup error: ", err)
		return
	}

	domainSearchOptions := emailHunter.DomainSearchOptions{
		Domain: "stripe.com",
	}
	domainSearchResults, err := client.DomainSearch(domainSearchOptions)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", domainSearchResults)

	emailFinderOptions := emailHunter.EmailFinderOptions{
		Domain:    "asana.com",
		FirstName: "Dustin",
		LastName:  "Moskovitz",
	}
	emailFinderResults, err := client.EmailFinder(emailFinderOptions)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", emailFinderResults)

	emailVerificationResults, err := client.EmailVerification("steli@close.io")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", emailVerificationResults)

	emailCountResults, err := client.EmailCount("stripe.com")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", emailCountResults)

	accountInformationResults, err := client.AccountInformation()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", accountInformationResults)
}

```