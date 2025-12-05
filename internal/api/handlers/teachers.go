package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/Sandwichzzy/REST_API_GO/internal/models"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextID   = 1
)

// initialize some dummy data
func init() {
	teachers[nextID] = models.Teacher{ID: nextID, FirstName: "John", LastName: "Doe", Class: "10A", Subject: "Math"}
	nextID++
	teachers[nextID] = models.Teacher{ID: nextID, FirstName: "Jane", LastName: "Smith", Class: "10B", Subject: "Science"}
	nextID++
	teachers[nextID] = models.Teacher{ID: nextID, FirstName: "Jane", LastName: "Doe", Class: "10C", Subject: "History"}
	nextID++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	//r.URL.Path 只包含路径部分，不包括查询参数（? 后面的内容），也不包括协议、域名和端口。
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println("IDStr is ", idStr)

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		teacherList := make([]models.Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}
		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	//handle path parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(teacher)
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextID
		teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		nextID++
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

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//teachers/{id}
	//teachers?key=value&query=value2&sortby=email&sortorder=ASC
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodPut:
		w.Write([]byte("hello PUT METHOD on Teachers route"))
	case http.MethodPatch:
		w.Write([]byte("hello PATCH METHOD on Teachers route"))
	case http.MethodDelete:
		w.Write([]byte("hello DEL METHOD on Teachers route"))
	}

	//w.Write([]byte("hello Teachers route"))
	//fmt.Println("Hello Teachers Route")
}
