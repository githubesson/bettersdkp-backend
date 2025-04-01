package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"bettersdkp/client"
	"bettersdkp/models"
)

func ForwardLoginRequest(w http.ResponseWriter, r *http.Request) {
	client := client.GetCustomClient()

	err := r.ParseMultipartForm(10 << 20)
	if err != nil && err != http.ErrNotMultipart {

		http.Error(w, "Error parsing multipart form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err == http.ErrNotMultipart {
		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	usr := r.FormValue("usr")
	pwd := r.FormValue("pwd")

	fmt.Printf("Received login request - Username: %s", usr)

	if usr == "" || pwd == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	var data = strings.NewReader(fmt.Sprintf("action=Login&Usr=%s&Pwd=%s", usr, pwd))
	req, err := http.NewRequest("POST", "https://sdkp.pjwstk.edu.pl/sdkp", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	headers := SetCommonHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	})
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
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

	cookies := resp.Cookies()
	var sessionID string
	for _, cookie := range cookies {
		if cookie.Name == "JSESSIONID" {
			sessionID = cookie.Value
			break
		}
	}

	bodyStr := string(bodyText)

	loginResp := models.LoginResponse{
		Success:   strings.Contains(bodyStr, "success"),
		IsStudent: true,
	}

	userIdRegex := regexp.MustCompile(`UserId:'([^']*)'`)
	userIdMatch := userIdRegex.FindStringSubmatch(bodyStr)
	if len(userIdMatch) > 1 {
		loginResp.Data.UserId = userIdMatch[1]
	}

	userNameRegex := regexp.MustCompile(`UserName:'([^']*)'`)
	userNameMatch := userNameRegex.FindStringSubmatch(bodyStr)
	if len(userNameMatch) > 1 {
		loginResp.Data.UserName = userNameMatch[1]
	}

	w.Header().Set("Content-Type", "application/json")
	if sessionID != "" {

		cookie := &http.Cookie{
			Name:     "JSESSIONID",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			MaxAge:   3600 * 24,
		}

		http.SetCookie(w, cookie)

		loginResp.Data.SessionID = sessionID
	}

	json.NewEncoder(w).Encode(loginResp)
}
