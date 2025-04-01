package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"bettersdkp/client"
	"bettersdkp/middleware"
	"bettersdkp/models"
)

func ForwardTasksRequest(w http.ResponseWriter, r *http.Request) {
	c := client.GetCustomClient()

	sessionID := middleware.ExtractSessionID(r)
	if sessionID == "" {
		http.Error(w, "Missing session ID", http.StatusUnauthorized)
		return
	}

	var data = strings.NewReader(`action=StudentTasks`)
	req, err := http.NewRequest("POST", "https://sdkp.pjwstk.edu.pl/sdkp", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	headers := SetCommonHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":       "JSESSIONID=" + sessionID,
	})
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bodyStr := string(bodyText)

	tasksResp := models.TasksResponse{
		Success: strings.Contains(bodyStr, "success"),
	}

	tasksRegex := regexp.MustCompile(`\['([^']*)',\s*'([^']*)']`)
	matches := tasksRegex.FindAllStringSubmatch(bodyStr, -1)

	tasksResp.Data.StudentTasks = make([][]string, len(matches))
	for i, match := range matches {
		if len(match) > 2 {
			tasksResp.Data.StudentTasks[i] = []string{match[1], match[2]}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tasksResp)
}

func ForwardTaskResult(w http.ResponseWriter, r *http.Request) {
	c := client.GetCustomClient()

	sessionID := middleware.ExtractSessionID(r)
	if sessionID == "" {
		http.Error(w, "No session ID provided", http.StatusBadRequest)
		return
	}

	taskID := r.URL.Query().Get("task")
	if taskID == "" {
		http.Error(w, "No task ID provided", http.StatusBadRequest)
		return
	}

	data := strings.NewReader(fmt.Sprintf("action=StudentTaskResult&comboCid=%s", taskID))

	req, err := http.NewRequest("POST", "https://sdkp.pjwstk.edu.pl/sdkp", data)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	headers := SetCommonHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":       fmt.Sprintf("JSESSIONID=%s", sessionID),
	})
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
