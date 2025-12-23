package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Sandwichzzy/REST_API_GO/internal/models"
	"github.com/Sandwichzzy/REST_API_GO/internal/repository/sqlconnect"
	"github.com/Sandwichzzy/REST_API_GO/pkg/utils"
)

// GET
// /students/ or /students?first_name=John&sortby=last_name:ASC&sortby=class:DESC
func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var students []models.Student

	// Pagination
	limit, page := utils.GetPaginationParams(r)

	students, totalStudents, err := sqlconnect.GetStudentsDbHandler(students, r, limit, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status   string           `json:"status"`
		Count    int              `json:"count"`
		Page     int              `json:"page"`
		PageSize int              `json:"page_size"`
		Data     []models.Student `json:"data"`
	}{
		Status:   "success",
		Count:    totalStudents,
		Data:     students,
		Page:     page,
		PageSize: limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /students/{id}
func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	Student, err := sqlconnect.GetStudentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Student)
}

// POST /students/
func AddStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var newStudents []models.Student
	var rawStudents []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &rawStudents)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 获取 Student 结构体的字段名
	fields := GetFieldNames(models.Student{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, Student := range rawStudents {
		for key := range Student {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, fmt.Sprintf("Field %s is not allowed", key), http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(body, &newStudents)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// 检查必填字段是否为空
	for _, Student := range newStudents {
		err := CheckBlankFields(Student)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedStudents, err := sqlconnect.AddStudentsDBHandler(newStudents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}
	json.NewEncoder(w).Encode(response)
}

// PUT /students/{id}
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedStudentFromDB, err := sqlconnect.UpdateStudent(id, updatedStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudentFromDB)
}

// PATCH /students
func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchStudents(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//发送状态码和头部
	w.WriteHeader(http.StatusNoContent)

}

// PATCH /students/{id}
func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updateStudent, err := sqlconnect.PatchOneStudent(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateStudent)
}

// DELETE /students/{id}
func DeleteOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneStudent(id)
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
		Status: "Student successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

// DELETE /students/
func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteStudents(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "students successfully deleted",
		DeletedIDs: deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}
