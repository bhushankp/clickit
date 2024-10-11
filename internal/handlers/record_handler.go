package handlers

import (
	"clickit/internal/models"
	"clickit/internal/services"
	"clickit/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
)

// UploadExcel processes an uploaded Excel file
func UploadExcel(w http.ResponseWriter, r *http.Request) {
	// Retrieve the uploaded file from the request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		log.Printf("Error retrieving the file: %v", err) // Log the error for debugging
		return
	}
	defer file.Close() // Ensure the file is closed after processing

	// Open the uploaded Excel file using the excelize library
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Unable to read Excel file", http.StatusBadRequest)
		log.Printf("Error opening Excel file: %v", err) // Log the error for debugging
		return
	}

	// Get the list of sheet names in the Excel file
	sheetNames := excelFile.GetSheetList()
	if len(sheetNames) == 0 {
		http.Error(w, "No sheets found in the Excel file", http.StatusBadRequest)
		log.Println("No sheets found in the uploaded Excel file") // Log if no sheets are present
		return
	}

	// Use the first sheet for processing
	firstSheet := sheetNames[0]
	rows, err := excelFile.GetRows(firstSheet)
	if err != nil || len(rows) == 0 {
		http.Error(w, "Error reading rows from Excel file", http.StatusInternalServerError)
		log.Printf("Error reading rows: %v", err) // Log the error for debugging
		return
	}

	// Log the name of the sheet being processed and the number of rows
	log.Printf("Processing sheet: %s with %d rows", firstSheet, len(rows))
	log.Printf("Read headers: %v", rows[0]) // Log the headers read from the Excel file

	// Validate the headers against expected values
	if err := utils.ValidateHeaders(rows[0]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Header validation error: %v", err) // Log validation errors for debugging
		return
	}

	// Process the rows asynchronously to improve performance
	go services.ProcessExcelRows(rows[1:])

	// Respond to the client indicating the upload and processing was initiated
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded and processing initiated"))
}

// GetPaginatedRecords handles pagination for retrieving records
func GetPaginatedRecords(w http.ResponseWriter, r *http.Request) {
	// Retrieve query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	// Convert query parameters to integers
	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)

	// Default values for pagination
	if page == 0 {
		page = 1 // Set to first page if no page parameter provided
	}
	if perPage == 0 {
		perPage = 10 // Set default number of records per page
	}

	// Fetch paginated records from the service layer
	records, err := services.FetchPaginatedRecords(page, perPage)
	if err != nil {
		http.Error(w, "Error fetching records", http.StatusInternalServerError)
		log.Printf("Error fetching records: %v", err) // Log errors while fetching records
		return
	}

	// Set the response content type to JSON and encode the records
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// UpdateRecord updates a specific record in both MySQL and Redis
func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	// Get the record ID from the URL parameters
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	// Decode the request body into an updated record struct
	var updatedRecord models.Record
	if err := json.NewDecoder(r.Body).Decode(&updatedRecord); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err) // Log decoding errors for debugging
		return
	}

	// Update the record in both MySQL and Redis using the service layer
	err := services.UpdateRecord(id, updatedRecord)
	if err != nil {
		http.Error(w, "Error updating record", http.StatusInternalServerError)
		log.Printf("Error updating record with ID %d: %v", id, err) // Log update errors for debugging
		return
	}

	// Respond to the client indicating the record was updated successfully
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Record updated successfully"))
}

// DeleteRecord removes a record from MySQL and Redis
func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	// Get the record ID from the URL parameters
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	// Delete the record using the service layer
	err := services.DeleteRecord(id)
	if err != nil {
		http.Error(w, "Error deleting record", http.StatusInternalServerError)
		log.Printf("Error deleting record with ID %d: %v", id, err) // Log delete errors for debugging
		return
	}

	// Respond to the client indicating the record was deleted successfully
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Record deleted successfully"))
}
