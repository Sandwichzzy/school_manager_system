package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Sandwichzzy/REST_API_GO/internal/models"
	"github.com/Sandwichzzy/REST_API_GO/internal/repository/sqlconnect"
	"github.com/Sandwichzzy/REST_API_GO/pkg/utils"
)

// GET
// /execs/ or /execs?first_name=John&sortby=last_name:ASC&sortby=class:DESC
func GetExecsHandler(w http.ResponseWriter, r *http.Request) {

	var execs []models.Exec
	execs, err := sqlconnect.GetExecsDbHandler(execs, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execs),
		Data:   execs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /execs/{id}
func GetOneExecHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	Exec, err := sqlconnect.GetExecByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Exec)
}

// POST /execs/
func AddExecsHandler(w http.ResponseWriter, r *http.Request) {

	var newExecs []models.Exec
	var rawExecs []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &rawExecs)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 获取 Exec 结构体的字段名
	fields := GetFieldNames(models.Exec{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, Exec := range rawExecs {
		for key := range Exec {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, fmt.Sprintf("Field %s is not allowed", key), http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(body, &newExecs)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// 检查必填字段是否为空
	for _, Exec := range newExecs {
		err := CheckBlankFields(Exec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedExecs, err := sqlconnect.AddExecsDBHandler(newExecs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	json.NewEncoder(w).Encode(response)
}

// PATCH /execs
func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchExecs(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//发送状态码和头部
	w.WriteHeader(http.StatusNoContent)

}

// PATCH /execs/{id}
func PatchOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updateExec, err := sqlconnect.PatchOneExec(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateExec)
}

// DELETE /execs/{id}
func DeleteOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneExec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusNoContent)
	// response body
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Exec successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

// POST /execs/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Exec
	//Data Validation
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	//Search for user if user actually exists
	user, err := sqlconnect.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	// is user active
	if user.InactiveStatus {
		http.Error(w, "user is inactive", http.StatusForbidden)
		return
	}

	//verify password
	err = utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	//generate JWT token
	tokenString, err := utils.SignToken(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}
	//Send token as a response or a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	//response
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}
	json.NewEncoder(w).Encode(response)
}

// POST /execs/logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}

// PUT /execs/{id}/password
func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
		return
	}

	var req models.UpdatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// Validate request
	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "Current password and new password are required", http.StatusBadRequest)
		return
	}

	_, err = sqlconnect.UpdatePasswordInDb(userId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//send token as a response or as a cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "Bearer",
	// 	Value:    token,
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	Expires:  time.Now().Add(24 * time.Hour),
	// 	SameSite: http.SameSiteStrictMode,
	// })

	//Response body
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Password updated successfully",
	}
	// response := models.UpdatePasswordResponse{
	// 	Token:           token,
	// 	PasswordUpdated: true,
	// }
	json.NewEncoder(w).Encode(response)

}

// POST /execs/forgotpassword
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = sqlconnect.ForgotPasswordDbHandler(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//respond with success confirmation
	fmt.Fprintf(w, "Password reset email sent to %s", req.Email)

}

// POST /execs/resetpassword/reset/{resetcode}
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("resetcode")

	var req models.ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//TODO:Data validation
	if req.NewPassword == "" || req.ConfirmPassword == "" {
		http.Error(w, "New password and confirm password are required", http.StatusBadRequest)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "New password and confirm password do not match", http.StatusBadRequest)
		return
	}

	err = sqlconnect.ResetPasswordDbHandler(token, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Password has been reset successfully")
}
