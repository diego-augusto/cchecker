package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// ContainerData represents information about a container for the dashboard
type ContainerData struct {
	ID          string
	Name        string
	Status      string
	Running     bool
	Health      string
	CssClass    string
	HealthClass string
}

// DashboardData holds all data for the dashboard template
type DashboardData struct {
	Containers      []ContainerData
	TotalContainers int
	HealthyCount    int
	UnhealthyCount  int
	NotRunningCount int
	Timestamp       string
}

type HealthHandler struct {
	// Create a Docker client
	CLI *client.Client
}

func (h *HealthHandler) MainHandler(w http.ResponseWriter, r *http.Request) {
	containers, err := h.CLI.ContainerList(r.Context(), container.ListOptions{All: true})
	if err != nil {
		http.Error(w, "Error listing containers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var dashboard DashboardData
	dashboard.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	for _, container := range containers {
		containerInfo, err := h.CLI.ContainerInspect(r.Context(), container.ID)
		if err != nil {
			continue
		}

		// Get container name (removing leading slash)
		containerName := containerInfo.Name
		if len(containerName) > 0 && containerName[0] == '/' {
			containerName = containerName[1:]
		}

		// Determine health status and CSS classes
		healthStatus := "unknown"
		healthClass := "unknown"
		cssClass := ""

		if !containerInfo.State.Running {
			healthStatus = "not running"
			healthClass = "not-running"
			cssClass = "not-running"
			dashboard.NotRunningCount++
		} else if containerInfo.State.Health != nil {
			healthStatus = containerInfo.State.Health.Status
			healthClass = healthStatus

			switch healthStatus {
			case "healthy":
				dashboard.HealthyCount++
				cssClass = "healthy"
			case "unhealthy":
				dashboard.UnhealthyCount++
				cssClass = "unhealthy"
			case "starting":
				cssClass = "starting"
			}
		} else {
			healthStatus = "running (no health check)"
			healthClass = "healthy"
			cssClass = "healthy"
			dashboard.HealthyCount++
		}

		containerData := ContainerData{
			ID:          container.ID[:12],
			Name:        containerName,
			Status:      container.Status,
			Running:     containerInfo.State.Running,
			Health:      healthStatus,
			CssClass:    cssClass,
			HealthClass: healthClass,
		}

		dashboard.Containers = append(dashboard.Containers, containerData)
	}

	dashboard.TotalContainers = len(dashboard.Containers)

	// Accept JSON if requested via Accept header or query parameter
	if r.URL.Query().Get("format") == "json" || r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Use json encoder to send JSON response
		json.NewEncoder(w).Encode(dashboard)
		return
	}

	// Otherwise render HTML template
	tmplPath := filepath.Join("templates", "dashboard.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, dashboard)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}
