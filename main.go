package main

import (
	"fmt"
	"log"
	"net/http"

	"bettersdkp/handlers"
	"bettersdkp/middleware"
)

func main() {
	http.HandleFunc("/api/login", middleware.EnableCORS(handlers.ForwardLoginRequest))
	http.HandleFunc("/api/courses", middleware.EnableCORS(handlers.ForwardCoursesRequest))
	http.HandleFunc("/api/download-project", middleware.EnableCORS(handlers.ForwardProjectDownload))
	http.HandleFunc("/api/upload-project", middleware.EnableCORS(handlers.ForwardProjectUpload))
	http.HandleFunc("/api/tasks", middleware.EnableCORS(handlers.ForwardTasksRequest))
	http.HandleFunc("/api/task-result", middleware.EnableCORS(handlers.ForwardTaskResult))
	http.HandleFunc("/api/course-tasks", middleware.EnableCORS(handlers.ForwardCourseTasks))

	fmt.Println("Server starting on port 8070...")
	log.Fatal(http.ListenAndServe(":8070", nil))
}
