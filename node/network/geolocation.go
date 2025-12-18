package network

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Geolocation struct {
	IP      string
	Country string
	Region  string
	City    string
	Lat     float64
	Lon     float64
}

func GetGeolocation(ctx context.Context) (*Geolocation, error) {
	ip, err := getPublicIP(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get public IP: %w", err)
	}
	
	geo, err := getLocationFromIP(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("failed to get geolocation: %w", err)
	}
	
	geo.IP = ip
	return geo, nil
}

func getPublicIP(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.ipify.org?format=text", nil)
	if err != nil {
		return "", err
	}
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}

func getLocationFromIP(ctx context.Context, ip string) (*Geolocation, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Status      string  `json:"status"`
		Country     string  `json:"country"`
		RegionName  string  `json:"regionName"`
		City        string  `json:"city"`
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	if result.Status != "success" {
		return nil, fmt.Errorf("geolocation API returned status: %s", result.Status)
	}
	
	return &Geolocation{
		Country: result.Country,
		Region:  result.RegionName,
		City:   result.City,
		Lat:    result.Lat,
		Lon:    result.Lon,
	}, nil
}

