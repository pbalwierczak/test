package simulator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient handles HTTP communication with the server
type APIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// APIScooter represents a scooter from the API
type APIScooter struct {
	ID        string  `json:"id"`
	Status    string  `json:"status"`
	Latitude  float64 `json:"current_latitude"`
	Longitude float64 `json:"current_longitude"`
}

// ScootersResponse represents the response from GET /api/v1/scooters
type ScootersResponse struct {
	Scooters []APIScooter `json:"scooters"`
}

// TripStartRequest represents the request to start a trip
type TripStartRequest struct {
	UserID         string  `json:"user_id"`
	StartLatitude  float64 `json:"start_latitude"`
	StartLongitude float64 `json:"start_longitude"`
}

// TripStartResponse represents the response from starting a trip
type TripStartResponse struct {
	TripID string `json:"trip_id"`
}

// TripEndRequest represents the request to end a trip
type TripEndRequest struct {
	UserID string `json:"user_id"`
}

// LocationUpdateRequest represents the request to update location
type LocationUpdateRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetAvailableScooters fetches available scooters from the server
func (c *APIClient) GetAvailableScooters(ctx context.Context) ([]APIScooter, error) {
	url := fmt.Sprintf("%s/api/v1/scooters?status=available", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var response ScootersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Scooters, nil
}

// StartTrip starts a trip for a scooter
func (c *APIClient) StartTrip(ctx context.Context, scooterID, userID string, startLat, startLng float64) (*TripStartResponse, error) {
	url := fmt.Sprintf("%s/api/v1/scooters/%s/trip/start", c.baseURL, scooterID)

	requestBody := TripStartRequest{
		UserID:         userID,
		StartLatitude:  startLat,
		StartLongitude: startLng,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var response TripStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// EndTrip ends a trip for a scooter
func (c *APIClient) EndTrip(ctx context.Context, scooterID, userID string) error {
	url := fmt.Sprintf("%s/api/v1/scooters/%s/trip/end", c.baseURL, scooterID)

	requestBody := TripEndRequest{UserID: userID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// UpdateLocation updates the location of a scooter
func (c *APIClient) UpdateLocation(ctx context.Context, scooterID string, lat, lng float64) error {
	url := fmt.Sprintf("%s/api/v1/scooters/%s/location", c.baseURL, scooterID)

	requestBody := LocationUpdateRequest{
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}
