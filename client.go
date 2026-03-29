package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type QuickChartClient struct {
	BaseURL string
	Client  *http.Client
}

func NewQuickChartClient() *QuickChartClient {
	return &QuickChartClient{
		BaseURL: "https://quickchart.io", 
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}


type ChartConfig struct {
	Type    string                 `json:"type"`
	Data    map[string]interface{} `json:"data"` 
	Options map[string]interface{} `json:"options"`
}

func (c *QuickChartClient) GenerateChart(config ChartConfig, width, height int) ([]byte, error) {
	configJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s/chart?c=%s&w=%d&h=%d&bkg=white",
		c.BaseURL,
		url.QueryEscape(string(configJson)),
		width,
		height,
	)

	resp, err := c.Client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}


func (c *QuickChartClient) CreateLineChartConfig(title string, labels []string, values []float64) ChartConfig {
	return ChartConfig{
		Type: "line",
		Data: map[string]interface{}{
			"labels": labels,
			"datasets": []interface{}{
				map[string]interface{}{
					"label":           title,
					"data":            values,
					"borderColor":     "rgb(54, 162, 235)",
					"fill":            false,
				},
			},
		},
		Options: map[string]interface{}{
			"title": map[string]interface{}{"display": true, "text": title},
		},
	}
}
