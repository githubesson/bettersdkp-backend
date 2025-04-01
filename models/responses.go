package models

type LoginResponse struct {
	Success   bool `json:"success"`
	IsStudent bool `json:"isStudent"`
	Data      struct {
		UserId    string `json:"UserId"`
		UserName  string `json:"UserName"`
		SessionID string `json:"SessionID"`
	} `json:"data"`
}

type CoursesResponse struct {
	Success bool `json:"success"`
	Data    struct {
		StudentCourses [][]string `json:"StudentCourses"`
	} `json:"data"`
}

type TasksResponse struct {
	Success bool `json:"success"`
	Data    struct {
		StudentTasks [][]string `json:"StudentTasks"`
	} `json:"data"`
}

type CourseTasksResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Tasks []CourseTask `json:"tasks"`
	} `json:"data"`
}

type CourseTask struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Points  string `json:"points"`
	Content string `json:"content"`
}
