package models 

import "time"

type PingResponse struct {
	
	Date    time.Time `json:"date"`    
	Domain  string    `json:"domain"`  
	Latency int       `json:"latency"` 
	Status  int       `json:"status"`  
}
