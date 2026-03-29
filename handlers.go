package handlers

import (
	"microgo/models"
	"microgo/transform"
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ChartHandler struct{}

func (h *ChartHandler) HandleSyncChart(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения", 500)
		return
	}

	var responses []models.PingResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		http.Error(w, "Ошибка JSON: нужен массив", 400)
		return
	}

	
	groupedData := transform.GroupByDomain(responses)

	
	zipBuf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuf)

	
	for domain, data := range groupedData {
		
		chartConfig := map[string]interface{}{
			"type": "line",
			"data": map[string]interface{}{
				"labels": data.Labels,
				"datasets": []map[string]interface{}{
					{
						"label":       fmt.Sprintf("Latency for %s (ms)", domain),
						"data":        data.Values,
						"borderColor": "rgb(0, 255, 0)",
						"fill":        false,
					},
				},
			},
		}

		configBytes, _ := json.Marshal(chartConfig)
		postData, _ := json.Marshal(map[string]interface{}{
			"chart":  string(configBytes),
			"format": "png",
		})

		
		resp, err := http.Post("https://quickchart.io/chart", "application/json", bytes.NewBuffer(postData))
		if err != nil || resp.StatusCode != 200 {
			log.Printf("Ошибка получения графика для %s", domain)
			continue
		}

		
		fileName := fmt.Sprintf("chart_%s.png", domain)
		f, err := zipWriter.Create(fileName)
		if err != nil {
			continue
		}

		
		io.Copy(f, resp.Body)
		resp.Body.Close()
	}

	
	zipWriter.Close()

	
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"charts.zip\"")
	w.Write(zipBuf.Bytes())

	log.Printf("ZIP отосланы: обработано %d доменов", len(groupedData))
}