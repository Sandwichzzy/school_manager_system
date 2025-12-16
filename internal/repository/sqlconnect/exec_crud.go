package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/Sandwichzzy/REST_API_GO/internal/models"
	"github.com/Sandwichzzy/REST_API_GO/pkg/utils"
)

// GET /execs/{id}
func GetExecByID(id int) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "error connecting DB")
	}
	defer db.Close()

	var Exec models.Exec
	err = db.QueryRow("SELECT id,first_name,last_name,email,username,user_created_at,inactive_status, role FROM execs WHERE id=?", id).Scan(
		&Exec.ID, &Exec.FirstName, &Exec.LastName,
		&Exec.Email, &Exec.Username, &Exec.UserCreatedAt,
		&Exec.InactiveStatus, &Exec.Role)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.ErrorHandler(err, "Exec not found")
	} else if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "error retrieving data")
	}
	return Exec, nil
}

// GET
// /execs/ or /execs?first_name=John&sortby=last_name:ASC&sortby=xxx:DESC
func GetExecsDbHandler(execs []models.Exec, r *http.Request) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := "SELECT id,first_name,last_name,email,username,user_created_at,inactive_status,role FROM execs WHERE 1=1"
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	query = utils.AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "error query data")
	}
	defer rows.Close()

	// ExecList := make([]models.Exec, 0)
	var Exec models.Exec
	for rows.Next() {
		// var Exec models.Exec
		err := rows.Scan(&Exec.ID, &Exec.FirstName, &Exec.LastName, &Exec.Email,
			&Exec.Username, &Exec.UserCreatedAt, &Exec.InactiveStatus, &Exec.Role)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error retrieving data")
		}
		execs = append(execs, Exec)
	}
	return execs, nil
}

// POST /execs/
func AddExecsDBHandler(newExecs []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", models.Exec{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error in preparing SQL query")
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))
	for i, newExec := range newExecs {
		values := utils.GetStructValues(newExec)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error in inserting data into database")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error in fetching last insert ID")
		}
		newExec.ID = int(lastID)
		addedExecs[i] = newExec
	}
	return addedExecs, nil
}

// PATCH /execs/
func PatchExecs(updates []map[string]interface{}) error {
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
		fmt.Println("idStr:", idStr)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid or missing Exec ID in update")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid Exec ID format")
		}

		var ExecFromDb models.Exec
		err = db.QueryRow("SELECT id,first_name,last_name,email,username FROM execs WHERE id=?", id).Scan(
			&ExecFromDb.ID, &ExecFromDb.FirstName, &ExecFromDb.LastName, &ExecFromDb.Email, &ExecFromDb.Username)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err, "Exec not found")
			}
			return utils.ErrorHandler(err, "Error retrieving Exec data")
		}
		//apply updates using reflection
		ExecVal := reflect.ValueOf(&ExecFromDb).Elem()
		ExecType := ExecVal.Type() // models.Exec type

		for k, v := range update {
			if k == "id" {
				continue //skip ID field
			}
			for i := 0; i < ExecVal.NumField(); i++ {
				field := ExecType.Field(i)

				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := ExecVal.Field(i)
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
		_, err = tx.Exec("UPDATE execs SET first_name=?, last_name=?, email=?, username = ? WHERE id=?",
			ExecFromDb.FirstName, ExecFromDb.LastName, ExecFromDb.Email, ExecFromDb.Username, ExecFromDb.ID)

		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating Exec")
		}
	}
	//commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Error committing transaction")
	}
	return nil
}

// PATCH /execs/{id}
func PatchOneExec(id int, updates map[string]interface{}) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	var existingExec models.Exec
	err = db.QueryRow("SELECT id,first_name,last_name,email,username FROM execs WHERE id=?", id).Scan(
		&existingExec.ID, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.ErrorHandler(err, "Exec not found")
	} else if err != nil {
		return models.
			//apply updates using reflect
			Exec{}, utils.ErrorHandler(err, "Unable to retrieve data")
	}

	ExecVal := reflect.ValueOf(&existingExec).Elem()
	ExecType := ExecVal.Type() // models.Exec type

	for k, v := range updates {
		for i := 0; i < ExecVal.NumField(); i++ {
			field := ExecType.Field(i)
			//id,omitempty  first_name,omitempty ...
			if field.Tag.Get("json") == k+",omitempty" {
				if ExecVal.Field(i).CanSet() {
					fieldVal := ExecVal.Field(i)
					fmt.Println("fieldVal:", fieldVal)
					fmt.Println("ExecVal.Field(i).Type():", ExecVal.Field(i).Type())
					fmt.Println("reflect.valueof(v):", reflect.ValueOf(v))
					fieldVal.Set(reflect.ValueOf(v).Convert(ExecVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name=?, last_name=?, email=?, username=? WHERE id=?",
		existingExec.FirstName, existingExec.LastName, existingExec.Email, existingExec.Username, existingExec.ID)

	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error updating Exec")
	}
	return existingExec, nil
}

// DELETE /execs/{id}
func DeleteOneExec(id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id=?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Error executing delete query")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error retrieving delete result")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(nil, fmt.Sprintf("Exec ID %d not found", id))
	}
	return nil
}
