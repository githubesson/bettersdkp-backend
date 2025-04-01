package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"bettersdkp/client"
	"bettersdkp/middleware"
	"bettersdkp/models"

	"golang.org/x/net/html"
)

func ForwardCoursesRequest(w http.ResponseWriter, r *http.Request) {
	c := client.GetCustomClient()

	sessionID := middleware.ExtractSessionID(r)
	if sessionID == "" {
		http.Error(w, "Missing session ID", http.StatusUnauthorized)
		return
	}

	var data = strings.NewReader("action=CurrentCourses")
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

	coursesResp := models.CoursesResponse{
		Success: strings.Contains(bodyStr, "success"),
	}

	coursesRegex := regexp.MustCompile(`\['([^']*)',\s*'([^']*)']`)
	matches := coursesRegex.FindAllStringSubmatch(bodyStr, -1)

	coursesResp.Data.StudentCourses = make([][]string, len(matches))
	for i, match := range matches {
		if len(match) > 2 {
			coursesResp.Data.StudentCourses[i] = []string{match[1], match[2]}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(coursesResp)
}

func ForwardCourseTasks(w http.ResponseWriter, r *http.Request) {
	c := client.GetCustomClient()

	taskURL := r.URL.Query().Get("url")
	if taskURL == "" {
		http.Error(w, "No task URL provided", http.StatusBadRequest)
		return
	}

	resp, err := c.Get(taskURL)
	if err != nil {
		log.Printf("Error fetching task page: %v", err)
		http.Error(w, "Failed to fetch task page", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		http.Error(w, "Failed to parse task page", http.StatusInternalServerError)
		return
	}

	tasks := extractCourseTasksFromHTML(doc, c)

	response := models.CourseTasksResponse{
		Success: true,
	}
	response.Data.Tasks = tasks

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func extractCourseTasksFromHTML(doc *html.Node, c *http.Client) []models.CourseTask {
	var tasks []models.CourseTask
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href, text string
			for _, a := range n.Attr {
				if a.Key == "href" {
					href = a.Val
					break
				}
			}
			if n.FirstChild != nil {
				text = n.FirstChild.Data
			}
			if href != "" && text != "" {

				taskName := text
				if idx := strings.Index(text, "(max."); idx != -1 {
					taskName = strings.TrimSpace(text[:idx])
				}

				task := models.CourseTask{
					Name:   taskName,
					URL:    href,
					Points: ExtractPoints(text),
				}

				taskResp, err := c.Get(href)
				if err == nil {
					defer taskResp.Body.Close()
					if contentBytes, err := io.ReadAll(taskResp.Body); err == nil {
						task.Content = string(contentBytes)
					} else {
						log.Printf("Error reading task content for %s: %v", href, err)
					}
				} else {
					log.Printf("Error fetching task content for %s: %v", href, err)
				}

				tasks = append(tasks, task)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return tasks
}
