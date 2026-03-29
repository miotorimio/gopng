package transform

import (
	"microgo/models"
)

type ChartData struct {
	Labels []string
	Values []int
}


func GroupByDomain(responses []models.PingResponse) map[string]ChartData {
	grouped := make(map[string]ChartData)

	for _, res := range responses {
		data := grouped[res.Domain]
		data.Labels = append(data.Labels, res.Date.Format("15:04:05"))
		data.Values = append(data.Values, res.Latency)
		grouped[res.Domain] = data
	}
	return grouped
}