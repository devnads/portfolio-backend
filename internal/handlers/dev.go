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

func CreateDev(w http.ResponseWriter, r *http.Request) {
	var dev models.Dev
	if err := json.NewDecoder(r.Body).Decode(&dev); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	dev.CreatedAt = time.Now()

	_, err := db.Conn.Exec(context.Background(), `
		INSERT INTO devs (id, username, profile_image, roles, address, github, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,dev.ID, dev.Username, dev.ProfileImage, dev.Roles, dev.Address, dev.Github, dev.CreatedAt)

	if err != nil {
		http.Error(w, "DB error while inserting", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message":"Successfully created a dev"})
}

func GetAllDevs(w http.ResponseWriter, r *http.Request){
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort_by")
	sortDir := r.URL.Query().Get("sort_dir")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	query := `SELECT id, username, profile_image, roles, discord, twitter, address, github, created_at FROM devs WHERE 1=1`
	args := []interface{}{}
	i := 1

	if search != "" {
		query += ` AND username ILIKE $` + fmt.Sprint(i)
		args = append(args, "%"+search+"%")
		i++
	}

	validSortFields := map[string]bool{
		"username":true, "created_at":true, "id": true,
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


	rows,err := db.Conn.Query(context.Background(), query, args...)

	if err != nil {
		log.Printf("DB error: %v | Query: %s | Args: %v", err, query, args)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	
	defer rows.Close()

	var devs []models.Dev 

	for rows.Next() {
		var dev models.Dev
		rows.Scan(&dev.ID,  &dev.Username, &dev.ProfileImage, &dev.Roles, &dev.Discord, &dev.Twitter, &dev.Address, &dev.Github, &dev.CreatedAt)
		devs = append(devs, dev)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devs)
}

func GetDev(w http.ResponseWriter, r *http.Request){
	devID := r.URL.Query().Get("dev_id")
	includeProjects := r.URL.Query().Get("include") == "projects"

	if devID == "" {
		http.Error(w, "Missing dev_id", http.StatusBadRequest)
		return
	}

	var dev models.Dev

	err := db.Conn.QueryRow(context.Background(), 
		`SELECT id, username, profile_image, roles, discord, twitter, address, github, created_at 
		FROM devs WHERE id = $1`, devID).Scan(
		&dev.ID, 
		&dev.Username, 
		&dev.ProfileImage, 
		&dev.Roles, 
		&dev.Discord, 
		&dev.Twitter,
		&dev.Address, 
		&dev.Github,
		&dev.CreatedAt,
	)

	if err != nil {
        http.Error(w, "Dev not found", http.StatusNotFound)
        return
    }
	if includeProjects {
		query := `SELECT id, dev_id, mission_id, name, image, categories, description, created_at FROM projects WHERE dev_id = $1`
		rows, err := db.Conn.Query(context.Background(), query, devID)
		if err != nil {
			http.Error(w, "Error fetching projects", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var projects []models.Project
		for rows.Next() {
			var project models.Project
			err := rows.Scan(
				&project.ID,
				&project.DevID,
				&project.MissionID,
				&project.Name,
				&project.Image,
				&project.Categories,
				&project.Description,
				&project.CreatedAt,
			)
			if err != nil {
				http.Error(w, "Error scanning project", http.StatusInternalServerError)
				return
			}
			projects = append(projects, project)
		}

		response := map[string]interface{}{
			"dev": dev,
			"projects": projects,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dev)
}



func GetDevProjects(w http.ResponseWriter, r *http.Request) {
	devID := r.URL.Query().Get("dev_id")
	search := r.URL.Query().Get("search")

	if devID == "" {
		http.Error(w, "Missing dev_id", http.StatusBadRequest)
		return
	} 

	query := `SELECT id, dev_id, mission_id, name, image, categories, description, created_at FROM projects WHERE dev_id = $1`
	args := []interface{}{devID}
	i := 2
	
	if search != "" {
	    query += ` AND (name ILIKE $` + fmt.Sprint(i) + ` OR $` + fmt.Sprint(i) + ` = ANY(categories))`
		args = append(args,search)
		i++
	}

	rows, err := db.Conn.Query(context.Background(), query, args...)
	if err != nil {
		http.Error(w, "Error while fetching data", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next(){
		var project models.Project
		rows.Scan(
			&project.ID,         
			&project.DevID,      
			&project.MissionID,  
			&project.Name,        
			&project.Image,       
			&project.Categories, 
			&project.Description,
			&project.CreatedAt, 
		)
		projects = append(projects, project)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&projects)

}
