package handler

import "net/http"

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}
