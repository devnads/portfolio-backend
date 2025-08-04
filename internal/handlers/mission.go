package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"monad-indexer/internal/db"
	"monad-indexer/internal/models"
	"net/http"
	"strconv"
	"time"
)

func CreateMission(w http.ResponseWriter, r *http.Request) {
	var mission models.Mission
	if err := json.NewDecoder(r.Body).Decode(&mission); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	mission.CreatedAt = time.Now()
}
func GetAllMissions(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort_by")
	sortDir := r.URL.Query().Get("sort_dir")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	query := `SELECT id, name, start_time, end_time, round, description, created_at FROM missions WHERE 1=1`
	args := []interface{}{}
	i := 1

	if search != "" {
		query += ` AND name ILIKE $` + fmt.Sprint(i)
		args = append(args, "%"+search+"%")
		i++
	}

	validSortFields := map[string]bool{
		"name": true, "created_at": true,
	}

	if sortBy != "" && validSortFields[sortBy] {
		direction := "ASC"
		if sortDir == "desc" {
			direction = "DESC"
		}
		query += ` ORDER BY ` + sortBy + ` ` + direction
	} else {
		query += ` ORDER BY created_at DESC`
	}

	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	query += ` LIMIT $` + fmt.Sprint(i)
	args = append(args, limit)
	i++

	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	query += ` OFFSET $` + fmt.Sprint(i)
	args = append(args, offset)

	rows, err := db.Conn.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("DB Error (missions): %+v\n", err)
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var missions []models.Mission

	for rows.Next() {
		var mission models.Mission
		err := rows.Scan(&mission.ID, &mission.Name, &mission.StartTime, &mission.EndTime, &mission.Round, &mission.Description, &mission.CreatedAt)
		if err != nil {
			log.Printf("Error scanning mission: %+v\n", err)
			continue
		}

		projectRows, err := db.Conn.Query(context.Background(), `
			SELECT id, name, image, categories, description, created_at
			FROM projects
			WHERE mission_id = $1
		`, mission.ID)
		if err != nil {
			log.Printf("DB Error (projects): %+v\n", err)
			http.Error(w, "DB Error", http.StatusInternalServerError)
			return
		}

		var projects []models.Project
		for projectRows.Next() {
			var project models.Project
			err := projectRows.Scan(&project.ID, &project.Name, &project.Image, &project.Categories, &project.Description, &project.CreatedAt)
			if err != nil {
				log.Printf("Error scanning project: %+v\n", err)
				continue
			}
			projects = append(projects, project)
		}
		projectRows.Close()

		mission.Projects = projects
		missions = append(missions, mission)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(missions)
}
