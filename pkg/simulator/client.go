package simulator

import (
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

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error string `json:"error"`
}

// ClosestScootersResponse represents the response from GET /api/v1/scooters/closest
type ClosestScootersResponse struct {
	Scooters []APIScooter `json:"scooters"`
	Count    int          `json:"count"`
}

// ScooterListResponse represents the response from GET /api/v1/scooters
type ScooterListResponse struct {
	Scooters []APIScooter `json:"scooters"`
	Count    int          `json:"count"`
	Limit    int          `json:"limit"`
	Offset   int          `json:"offset"`
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

// GetAllScooters fetches all scooters from the server
func (c *APIClient) GetAllScooters(ctx context.Context) ([]APIScooter, error) {
	url := fmt.Sprintf("%s/api/v1/scooters", c.baseURL)

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

// GetClosestScooters finds the closest available scooters to a given location
func (c *APIClient) GetClosestScooters(ctx context.Context, lat, lng float64, radius int, limit int) ([]APIScooter, error) {
	url := fmt.Sprintf("%s/api/v1/scooters/closest?lat=%.6f&lng=%.6f&radius=%d&limit=%d&status=available",
		c.baseURL, lat, lng, radius, limit)

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

	var response ClosestScootersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Scooters, nil
}

// GetScootersInBounds finds scooters within geographic bounds
func (c *APIClient) GetScootersInBounds(ctx context.Context, minLat, maxLat, minLng, maxLng float64, limit int) ([]APIScooter, error) {
	url := fmt.Sprintf("%s/api/v1/scooters?min_lat=%.6f&max_lat=%.6f&min_lng=%.6f&max_lng=%.6f&limit=%d&status=available",
		c.baseURL, minLat, maxLat, minLng, maxLng, limit)

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

	var response ScooterListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Scooters, nil
}
