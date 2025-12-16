package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"s21/models"
	"slices"
)

func ParticipantsFromJSON() ([]string, error) {
	data, err := os.ReadFile(filepath.Join("jsons", "participants.json"))
	if err != nil {
		return nil, err
	}

	var response models.ParticipantsList
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response.Participants, nil
}

func SlicesDiff(s1, s2 []string) []string {
	var result []string
	for _, elem := range s1 {
		if !slices.Contains(s2, elem) {
			result = append(result, elem)
		}
	}
	return result
}

func S21Token() string {
	data, err := os.ReadFile("auth.json")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	return response.AccessToken
}
