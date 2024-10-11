package services

import (
	"clickit/internal/config"
	"clickit/internal/models"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

// ProcessExcelRows handles Excel row processing and inserts records into MySQL
func ProcessExcelRows(rows [][]string) {
	var records []models.Record
	// Loop through each row in the Excel data
	for _, row := range rows {
		// Create a new Record instance from the row data
		record := models.Record{
			FirstName: row[0],
			LastName:  row[1],
			Company:   row[2],
			Address:   row[3],
			City:      row[4],
			County:    row[5],
			Postal:    row[6],
			Phone:     row[7],
			Email:     row[8],
			Web:       row[9],
		}
		records = append(records, record) // Append the record to the slice
	}

	// Batch insert into MySQL
	if err := config.DB.Create(&records).Error; err != nil {
		log.Printf("Error inserting records into MySQL: %v", err) // Log any error during insertion
		return
	}

	// Cache the records in Redis
	for _, record := range records {
		cacheRecord(record)
	}
}

// cacheRecord stores a single record in Redis with a defined expiration time
func cacheRecord(record models.Record) {
	data, err := json.Marshal(record) // Convert the record to JSON
	if err != nil {
		log.Printf("Error marshaling record to JSON: %v", err) // Log any JSON marshaling errors
		return
	}
	cacheKey := "record:" + strconv.Itoa(int(record.ID))      // Generate the cache key
	config.RDB.Set(config.Ctx, cacheKey, data, 5*time.Minute) // Set the record in Redis with a 5-minute expiration
}

// FetchPaginatedRecords retrieves paginated records from MySQL
func FetchPaginatedRecords(page, perPage int) ([]models.Record, error) {
	var records []models.Record
	offset := (page - 1) * perPage                                      // Calculate the offset for pagination
	err := config.DB.Limit(perPage).Offset(offset).Find(&records).Error // Fetch the records with limit and offset
	return records, err                                                 // Return the records and any error encountered
}

// UpdateRecord updates a record in both MySQL and Redis
func UpdateRecord(id int, updatedRecord models.Record) error {
	var record models.Record
	// Find the existing record by ID
	if err := config.DB.First(&record, id).Error; err != nil {
		return err // Return error if the record is not found
	}

	// Update the record fields with the new values
	record.FirstName = updatedRecord.FirstName
	record.LastName = updatedRecord.LastName
	record.Company = updatedRecord.Company
	record.Address = updatedRecord.Address
	record.City = updatedRecord.City
	record.County = updatedRecord.County
	record.Postal = updatedRecord.Postal
	record.Phone = updatedRecord.Phone
	record.Email = updatedRecord.Email
	record.Web = updatedRecord.Web

	// Save the updated record back to MySQL
	if err := config.DB.Save(&record).Error; err != nil {
		return err // Return error if saving fails
	}

	// Cache the updated record in Redis
	cacheRecord(record)
	return nil // Return nil if update is successful
}

// DeleteRecord deletes a record from MySQL and Redis
func DeleteRecord(id int) error {
	var record models.Record
	// Find the existing record by ID
	if err := config.DB.First(&record, id).Error; err != nil {
		return err // Return error if the record is not found
	}

	// Delete the record from MySQL
	if err := config.DB.Delete(&record).Error; err != nil {
		return err // Return error if deleting fails
	}

	// Remove the record from Redis
	config.RDB.Del(config.Ctx, "record:"+strconv.Itoa(id))
	return nil // Return nil if deletion is successful
}
