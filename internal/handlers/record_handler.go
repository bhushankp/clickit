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

func UploadExcel(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		log.Printf("Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Unable to read Excel file", http.StatusBadRequest)
		log.Printf("Error opening Excel file: %v", err)
		return
	}

	sheetNames := excelFile.GetSheetList()
	if len(sheetNames) == 0 {
		http.Error(w, "No sheets found in the Excel file", http.StatusBadRequest)
		log.Println("No sheets found in the uploaded Excel file")
		return
	}

	firstSheet := sheetNames[0]
	rows, err := excelFile.GetRows(firstSheet)
	if err != nil || len(rows) == 0 {
		http.Error(w, "Error reading rows from Excel file", http.StatusInternalServerError)
		log.Printf("Error reading rows: %v", err)
		return
	}

	log.Printf("Processing sheet: %s with %d rows", firstSheet, len(rows))
	log.Printf("Read headers: %v", rows[0])

	if err := utils.ValidateHeaders(rows[0]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Header validation error: %v", err)
		return
	}

	go services.ProcessExcelRows(rows[1:])

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded and processing initiated"))
}

func GetPaginatedRecords(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)

	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = 10
	}

	records, err := services.FetchPaginatedRecords(page, perPage)
	if err != nil {
		http.Error(w, "Error fetching records", http.StatusInternalServerError)
		log.Printf("Error fetching records: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updatedRecord models.Record
	if err := json.NewDecoder(r.Body).Decode(&updatedRecord); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	err := services.UpdateRecord(id, updatedRecord)
	if err != nil {
		http.Error(w, "Error updating record", http.StatusInternalServerError)
		log.Printf("Error updating record with ID %d: %v", id, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Record updated successfully"))
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	err := services.DeleteRecord(id)
	if err != nil {
		http.Error(w, "Error deleting record", http.StatusInternalServerError)
		log.Printf("Error deleting record with ID %d: %v", id, err) // Log delete errors for debugging
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Record deleted successfully"))
}
