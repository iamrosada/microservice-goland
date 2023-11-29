package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/entity"
)

func GenerateNewID() uint {
	return uint(time.Now().UnixNano())
}

func FetchUsersFromUserAppliedMicroservice(url string) ([]int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch users, status: %d", resp.StatusCode)
	}

	rawBody, _ := io.ReadAll(resp.Body)
	log.Printf("Raw Response Body: %s", string(rawBody))

	var response struct {
		UserIDs []int `json:"user_ids"`
	}

	err = json.Unmarshal(rawBody, &response)
	if err != nil {
		return nil, err
	}

	log.Printf("userIDs: %+v", response.UserIDs)

	return response.UserIDs, nil
}

func FetchUsersFromUserMicroservice(url string) ([]entity.User, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch users, status: %d", resp.StatusCode)
	}

	var response struct {
		Users []entity.User `json:"users"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Users, nil
}
