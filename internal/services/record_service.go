package services

import (
	"clickit/internal/config"
	"clickit/internal/models"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

func ProcessExcelRows(rows [][]string) {
	var records []models.Record
	for _, row := range rows {
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
		records = append(records, record)
	}

	if err := config.DB.Create(&records).Error; err != nil {
		log.Printf("Error inserting records into MySQL: %v", err)
		return
	}

	for _, record := range records {
		cacheRecord(record)
	}
}

func cacheRecord(record models.Record) {
	data, err := json.Marshal(record)
	if err != nil {
		log.Printf("Error marshaling record to JSON: %v", err)
		return
	}
	cacheKey := "record:" + strconv.Itoa(int(record.ID))
	config.RDB.Set(config.Ctx, cacheKey, data, 5*time.Minute)
}

func FetchPaginatedRecords(page, perPage int) ([]models.Record, error) {
	var records []models.Record
	offset := (page - 1) * perPage
	err := config.DB.Limit(perPage).Offset(offset).Find(&records).Error
	return records, err
}

func UpdateRecord(id int, updatedRecord models.Record) error {
	var record models.Record
	if err := config.DB.First(&record, id).Error; err != nil {
		return err
	}

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

	if err := config.DB.Save(&record).Error; err != nil {
		return err
	}

	cacheRecord(record)
	return nil
}

func DeleteRecord(id int) error {
	var record models.Record
	if err := config.DB.First(&record, id).Error; err != nil {
		return err
	}

	if err := config.DB.Delete(&record).Error; err != nil {
		return err
	}

	config.RDB.Del(config.Ctx, "record:"+strconv.Itoa(id))
	return nil
}
