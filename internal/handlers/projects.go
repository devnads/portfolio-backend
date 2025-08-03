package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"monad-indexer/internal/db"
	"monad-indexer/internal/models"
	"net/http"
	"strconv"
	"time"
)

func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	devID := r.URL.Query().Get("dev_id")
	missionID := r.URL.Query().Get("mission_id")
	category := r.URL.Query().Get("categories")
	name := r.URL.Query().Get("name")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort_by")
	sortDir := r.URL.Query().Get("sort_dir")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	query := `SELECT id, dev_id, mission_id, name, image, categories, description, created_at FROM projects WHERE 1=1`
	args := []interface{}{}
	i := 1

	if search != "" {
		query += ` AND name ILIKE $` + fmt.Sprint(i)
		args = append(args, search)
		i++
	}

	if devID != "" {
		query += ` AND dev_id = $` + fmt.Sprint(i)
		args = append(args, devID)
		i++
	}

	if missionID != "" {
		query += ` AND mission_id = $` + fmt.Sprint(i)
		args = append(args, missionID)
		i++
	}

	if category != "" {
		query += ` AND $` + fmt.Sprint(i) + `= ANY(categories)`
		args = append(args, category)
		i++
	}

	if name != "" {
		query += `AND name = $` + fmt.Sprint(i) 
		args = append(args, name)
		i++
	}

	validSortFields := map[string]bool{
		"name":true, "dev_id":true, "created_at": true,
	}

	if sortBy != "" && validSortFields[sortBy] {
		direction := "ASC"
		if sortDir == "desc" {
			direction = "DESC"
		}
		query += ` ORDER BY `+ sortBy + ` ` + direction
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
		http.Error(w,"DB Error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var project models.Project
		rows.Scan(&project.ID, &project.DevID, &project.MissionID, &project.Name, &project.Image, &project.Categories, &project.Description, project.CreatedAt)
		projects = append(projects, project)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project 
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w,"Error while decoding body", http.StatusBadRequest)
		return
	}
	
	project.CreatedAt = time.Now()
	_, err := db.Conn.Exec(context.Background(),
		`INSERT INTO projects (name, creator_id, image, categories, description, created_at) VALUES ($1 $2 $3 $4 $5 $6)`,
		&project.Name, &project.Image, &project.Categories, &project.Description, &project.CreatedAt,
	)

	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message":"Successfully created a project"})
}

