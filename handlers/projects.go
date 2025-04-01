package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"bettersdkp/client"
	"bettersdkp/middleware"
)

func ForwardProjectDownload(w http.ResponseWriter, r *http.Request) {
	c := client.GetCustomClient()

	courseStr := r.URL.Query().Get("course")
	if courseStr == "" {
		http.Error(w, "Missing course parameter", http.StatusBadRequest)
		return
	}

	sessionID := middleware.ExtractSessionID(r)
	if sessionID == "" {
		http.Error(w, "Missing session ID", http.StatusUnauthorized)
		return
	}

	projectUrl := fmt.Sprintf("https://sdkp.pjwstk.edu.pl/sdkp?action=LoadProject&comboLab=%s", courseStr)
	req, err := http.NewRequest("GET", projectUrl, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	headers := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language":           "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":                "keep-alive",
		"Referer":                   "https://sdkp.pjwstk.edu.pl/",
		"Sec-Fetch-Dest":            "iframe",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "same-origin",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"sec-ch-ua":                 `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        `"Windows"`,
		"Cookie":                    "JSESSIONID=" + sessionID,
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to download project: %d", resp.StatusCode), resp.StatusCode)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error streaming file: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func ForwardProjectUpload(w http.ResponseWriter, r *http.Request) {
	c := client.GetCustomClient()

	sessionID := middleware.ExtractSessionID(r)
	if sessionID == "" {
		http.Error(w, "Missing session ID", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	studentLabs := r.FormValue("studentLabs")
	if studentLabs == "" {
		http.Error(w, "Missing studentLabs parameter", http.StatusBadRequest)
		return
	}

	var projectFile multipart.File
	var projectFileName string
	file, fileHeader, err := r.FormFile("prjFile")
	if err == nil && fileHeader != nil {
		projectFile = file
		projectFileName = fileHeader.Filename
		defer projectFile.Close()
	}

	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)

	formField, err := writer.CreateFormField("studentLabs")
	if err != nil {
		http.Error(w, "Error creating form field: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = formField.Write([]byte(studentLabs))
	if err != nil {
		http.Error(w, "Error writing to form field: "+err.Error(), http.StatusInternalServerError)
		return
	}

	formField, err = writer.CreateFormField("action")
	if err != nil {
		http.Error(w, "Error creating form field: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = formField.Write([]byte("UploadProject"))
	if err != nil {
		http.Error(w, "Error writing to form field: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if projectFile != nil {
		fileField, err := writer.CreateFormFile("prjFile", projectFileName)
		if err != nil {
			http.Error(w, "Error creating file field: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(fileField, projectFile)
		if err != nil {
			http.Error(w, "Error copying file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {

		formField, err = writer.CreateFormField("prjFile")
		if err != nil {
			http.Error(w, "Error creating form field: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = formField.Write([]byte(""))
		if err != nil {
			http.Error(w, "Error writing to form field: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = writer.Close()
	if err != nil {
		http.Error(w, "Error closing form writer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "https://sdkp.pjwstk.edu.pl/sdkp", form)
	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	headers := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language":           "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7",
		"Cache-Control":             "max-age=0",
		"Connection":                "keep-alive",
		"Origin":                    "https://sdkp.pjwstk.edu.pl",
		"Referer":                   "https://sdkp.pjwstk.edu.pl/",
		"Sec-Fetch-Dest":            "iframe",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "same-origin",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"sec-ch-ua":                 `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        `"Windows"`,
		"Cookie":                    "JSESSIONID=" + sessionID,
		"Content-Type":              writer.FormDataContentType(),
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		http.Error(w, "Error sending request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
