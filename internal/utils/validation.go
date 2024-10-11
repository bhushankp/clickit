package utils

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

// ValidateHeaders ensures the Excel file has the correct headers
func ValidateHeaders(headers []string) error {
	expectedHeaders := []string{
		"first_name", "last_name", "company_name", "address", "city", "county", "postal", "phone", "email", "web",
	}

	// Log expected headers
	log.Printf("Expected headers: %v", expectedHeaders)

	// Check if the number of headers matches
	if len(headers) != len(expectedHeaders) {
		return errors.New("Invalid number of columns in the Excel file: expected " +
			strconv.Itoa(len(expectedHeaders)) + " but got " + strconv.Itoa(len(headers)))
	}

	// Validate each header against the expected headers
	for i, header := range headers {
		// Log the headers being compared
		log.Printf("Comparing header: '%s' with expected: '%s'", header, expectedHeaders[i])

		// Trim and convert header to lowercase for comparison
		if strings.TrimSpace(strings.ToLower(header)) != expectedHeaders[i] {
			return errors.New("Invalid Excel file format: Incorrect header '" + header +
				"' at position " + strconv.Itoa(i) + ". Expected '" + expectedHeaders[i] + "'")
		}
		log.Printf("Header '%s' matches expected header.", header) // Log if the header matches
	}

	return nil // Return nil if all headers are valid
}
