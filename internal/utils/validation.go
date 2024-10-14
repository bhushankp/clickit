package utils

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func ValidateHeaders(headers []string) error {
	expectedHeaders := []string{
		"first_name", "last_name", "company_name", "address", "city", "county", "postal", "phone", "email", "web",
	}

	log.Printf("Expected headers: %v", expectedHeaders)

	if len(headers) != len(expectedHeaders) {
		return errors.New("Invalid number of columns in the Excel file: expected " +
			strconv.Itoa(len(expectedHeaders)) + " but got " + strconv.Itoa(len(headers)))
	}

	for i, header := range headers {
		log.Printf("Comparing header: '%s' with expected: '%s'", header, expectedHeaders[i])

		if strings.TrimSpace(strings.ToLower(header)) != expectedHeaders[i] {
			return errors.New("Invalid Excel file format: Incorrect header '" + header +
				"' at position " + strconv.Itoa(i) + ". Expected '" + expectedHeaders[i] + "'")
		}
		log.Printf("Header '%s' matches expected header.", header)
	}

	return nil
}
