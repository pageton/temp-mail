// Package utils contains the utility functions for the application.
package utils

import (
	"log"
	"net/mail"
)

func ParseEmailAddresses(header string) []string {
	if header == "" {
		return nil
	}

	addrs, err := mail.ParseAddressList(header)
	if err != nil {
		log.Println("Error parsing addresses:", err)
		return nil
	}

	var result []string
	for _, addr := range addrs {
		result = append(result, addr.Address)
	}

	return result
}
