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
// /teachers/ or /teachers?first_name=John&sortby=last_name:ASC&sortby=class:DESC
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var teachers []models.Teacher
	teachers, err := sqlconnect.GetTeachersDbHandler(teachers, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /teachers/{id}
func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	teacher, err := sqlconnect.GetTeacherByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// POST /teachers/
func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var newTeachers []models.Teacher
	var rawTeachers []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &rawTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 获取 Teacher 结构体的字段名
	fields := GetFieldNames(models.Teacher{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, teacher := range rawTeachers {
		for key := range teacher {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, fmt.Sprintf("Field %s is not allowed", key), http.StatusBadRequest)
				return
			}
		}
	}

	err = json.Unmarshal(body, &newTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// 检查必填字段是否为空
	for _, teacher := range newTeachers {
		err := CheckBlankFields(teacher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedTeachers, err := sqlconnect.AddTeachersDBHandler(newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}

// PUT /teachers/{id}
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacher(id, updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDB)
}

// PATCH /teachers
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchTeachers(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//发送状态码和头部
	w.WriteHeader(http.StatusNoContent)

}

// PATCH /teachers/{id}
func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updateTeacher, err := sqlconnect.PatchOneTeacher(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateTeacher)
}

// DELETE /teachers/{id}
func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneTeacher(id)
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
		Status: "Teacher successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

// DELETE /teachers/
func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachers(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "Teachers successfully deleted",
		DeletedIDs: deletedIds,
	}
	json.NewEncoder(w).Encode(response)

}

// GET /teachers/{id}/students
func GetStudentsByTeacherIdHandler(w http.ResponseWriter, r *http.Request) {
	teacherId := r.PathValue("id")

	var students []models.Student

	students, err := sqlconnect.GetStudentsByTeacherId(teacherId, students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(students),
		Data:   students,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /students/{id}/studentcount
func GetStudentsCountByTeacherIdHandler(w http.ResponseWriter, r *http.Request) {
	teacherId := r.PathValue("id")

	studentCount, err := GetStudentsCountByTeacherId(teacherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Status: "success",
		Count:  studentCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetStudentsCountByTeacherId(teacherId string) (int, error) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		return 0, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := `SELECT COUNT(*) FROM students WHERE class =(SELECT class FROM teachers WHERE id = ?)`
	var studentCount int
	err = db.QueryRow(query, teacherId).Scan(&studentCount)
	if err != nil {
		return 0, utils.ErrorHandler(err, "error query data")
	}
	return studentCount, nil
}
