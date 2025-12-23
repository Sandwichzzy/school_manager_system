package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/Sandwichzzy/REST_API_GO/internal/models"
	"github.com/Sandwichzzy/REST_API_GO/pkg/utils"
)

// GET /students/{id}
func GetStudentByID(id int) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "error connecting DB")
	}
	defer db.Close()

	var Student models.Student
	err = db.QueryRow("SELECT id,first_name,last_name,email,class FROM students WHERE id=?", id).Scan(&Student.ID, &Student.FirstName, &Student.LastName, &Student.Email, &Student.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "Student not found")
	} else if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "error retrieving data")
	}
	return Student, nil
}

// GET
// /students/ or /students?first_name=John&sortby=last_name:ASC&sortby=class:DESC
func GetStudentsDbHandler(students []models.Student, r *http.Request, limit, page int) ([]models.Student, int, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := "SELECT id,first_name,last_name,email,class FROM students WHERE 1=1 "
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	//Add pagination
	offset := (page - 1) * limit
	query += "LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	query = utils.AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, 0, utils.ErrorHandler(err, "error query data")
	}
	defer rows.Close()

	// StudentList := make([]models.Student, 0)
	var Student models.Student
	for rows.Next() {
		// var Student models.Student
		err := rows.Scan(&Student.ID, &Student.FirstName, &Student.LastName, &Student.Email, &Student.Class)
		if err != nil {
			return nil, 0, utils.ErrorHandler(err, "error retrieving data")
		}
		students = append(students, Student)
	}
	// Get the total count of students
	var totalStudents int
	err = db.QueryRow("SELECT COUNT(*) FROM students").Scan(&totalStudents)
	if err != nil {
		utils.ErrorHandler(err, "error retrieving total student count")
		totalStudents = 0
	}
	return students, totalStudents, nil
}

// POST /students/
func AddStudentsDBHandler(newStudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO students(first_name, last_name, email, class) VALUES(?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("students", models.Student{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error in preparing SQL query")
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(newStudents))
	for i, newStudent := range newStudents {
		// res, err := stmt.Exec(newStudent.FirstName, newStudent.LastName, newStudent.Email, newStudent.Class)
		values := utils.GetStructValues(newStudent)
		res, err := stmt.Exec(values...)
		if err != nil {
			if strings.Contains(err.Error(), " a foreign key constraint fails (`school`.`students`, CONSTRAINT `students_ibfk_1` FOREIGN KEY (`class`) REFERENCES `teachers` (`class`))") {
				return nil, utils.ErrorHandler(err, "class/class teacher does not exist")
			}
			return nil, utils.ErrorHandler(err, "Error in inserting data into database")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error in fetching last insert ID")
		}
		newStudent.ID = int(lastID)
		addedStudents[i] = newStudent
	}
	return addedStudents, nil
}

// PUT /students/{id}
func UpdateStudent(id int, updatedStudent models.Student) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {

		return models.Student{}, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id,first_name,last_name,email,class FROM students WHERE id=?", id).Scan(
		&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email,
		&existingStudent.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "Student not found")
	} else if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Unable to retrieve data")
	}

	updatedStudent.ID = id
	_, err = db.Exec("UPDATE students SET first_name=?, last_name=?, email=?, class=? WHERE id=?",
		updatedStudent.FirstName, updatedStudent.LastName, updatedStudent.Email,
		updatedStudent.Class, id)

	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error updating Student")
	}
	return updatedStudent, nil
}

// PATCH /students/
func PatchStudents(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "Database transaction error")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid or missing Student ID in update")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid Student ID format")
		}

		var StudentFromDb models.Student
		err = db.QueryRow("SELECT id,first_name,last_name,email,class FROM students WHERE id=?", id).Scan(
			&StudentFromDb.ID, &StudentFromDb.FirstName, &StudentFromDb.LastName, &StudentFromDb.Email,
			&StudentFromDb.Class)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err, "Student not found")
			}
			return utils.ErrorHandler(err, "Error retrieving Student data")
		}
		//apply updates using reflection
		StudentVal := reflect.ValueOf(&StudentFromDb).Elem()
		StudentType := StudentVal.Type() // models.Student type

		for k, v := range update {
			if k == "id" {
				continue //skip ID field
			}
			for i := 0; i < StudentVal.NumField(); i++ {
				field := StudentType.Field(i)

				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := StudentVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("cannot convert %v to %v", val.Type(), fieldVal.Type())
							return err
						}
					}
					break
				}
			}
		}
		_, err = tx.Exec("UPDATE students SET first_name=?, last_name=?, email=?, class=? WHERE id=?",
			StudentFromDb.FirstName, StudentFromDb.LastName, StudentFromDb.Email,
			StudentFromDb.Class, id)

		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating Student")
		}
	}
	//commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Error committing transaction")
	}
	return nil
}

// PATCH /students/{id}
func PatchOneStudent(id int, updates map[string]interface{}) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id,first_name,last_name,email,class FROM students WHERE id=?", id).Scan(
		&existingStudent.ID, &existingStudent.FirstName, &existingStudent.LastName, &existingStudent.Email,
		&existingStudent.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "Student not found")
	} else if err != nil {
		return models.
			//apply updates using reflect
			Student{}, utils.ErrorHandler(err, "Unable to retrieve data")
	}

	StudentVal := reflect.ValueOf(&existingStudent).Elem()
	StudentType := StudentVal.Type() // models.Student type

	for k, v := range updates {
		for i := 0; i < StudentVal.NumField(); i++ {
			field := StudentType.Field(i)
			//id,omitempty  first_name,omitempty ...
			if field.Tag.Get("json") == k+",omitempty" {
				if StudentVal.Field(i).CanSet() {
					fieldVal := StudentVal.Field(i)
					fmt.Println("fieldVal:", fieldVal)
					fmt.Println("StudentVal.Field(i).Type():", StudentVal.Field(i).Type())
					fmt.Println("reflect.valueof(v):", reflect.ValueOf(v))
					fieldVal.Set(reflect.ValueOf(v).Convert(StudentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE students SET first_name=?, last_name=?, email=?, class=? WHERE id=?",
		existingStudent.FirstName, existingStudent.LastName, existingStudent.Email,
		existingStudent.Class, id)

	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "Error updating Student")
	}
	return existingStudent, nil
}

// DELETE /students/{id}
func DeleteOneStudent(id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM students WHERE id=?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Error executing delete query")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error retrieving delete result")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(nil, fmt.Sprintf("Student ID %d not found", id))
	}
	return nil
}

// DELETE /students/
func DeleteStudents(ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database transaction error")
	}
	//批量/重复数据库操作：预处理语句快很多（SQL 只在数据库编译一次）
	stmt, err := tx.Prepare("DELETE FROM students WHERE id=?")
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Error in preparing delete query")
	}
	defer stmt.Close()

	deletedIds := []int{}
	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Error executing delete query")
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("Error retrieving delete result for ID %d", id))
		}
		//if Student was deleted then add ID to the deletedIds slice
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			return nil, utils.ErrorHandler(nil, fmt.Sprintf("Student ID %d not found", id))
		}
	}

	//commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error committing transaction")
	}

	if len(deletedIds) < 1 {
		return nil, utils.ErrorHandler(nil, "IDs not exist")
	}
	return deletedIds, nil
}
