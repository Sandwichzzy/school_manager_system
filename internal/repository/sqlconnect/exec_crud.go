package sqlconnect

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/Sandwichzzy/REST_API_GO/internal/models"
	"github.com/Sandwichzzy/REST_API_GO/pkg/utils"
	"github.com/go-mail/mail/v2"
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
		newExec.Password, err = utils.HashPassword(newExec.Password)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error in hashing password")
		}

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

func GetUserByUsername(username string) (*models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal server error")
	}
	defer db.Close()

	user := &models.Exec{}
	err = db.QueryRow(`SELECT id,first_name,last_name,email,username,password,role,inactive_status FROM execs WHERE username=?`, username).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.Role, &user.InactiveStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrorHandler(err, "user not found")
		}
		return nil, utils.ErrorHandler(err, "database query error")
	}
	return user, nil
}

// POST /execs/{id}/updatepassword
func UpdatePasswordInDb(userId int, currentPassword, newPassword string) (bool, error) {
	db, err := ConnectDb()
	if err != nil {
		return false, utils.ErrorHandler(err, "database connection error")
	}
	defer db.Close()

	var username string
	var userPassword string
	var userRole string
	// 1. validate user existence
	err = db.QueryRow("SELECT username, password, role FROM execs WHERE id=?", userId).Scan(&username, &userPassword, &userRole)
	if err != nil {
		return false, utils.ErrorHandler(err, "user not found")
	}

	err = utils.VerifyPassword(currentPassword, userPassword)
	if err != nil {
		return false, utils.ErrorHandler(err, "current password is incorrect")
	}

	hashedNewPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return false, utils.ErrorHandler(err, "error hashing new password")
	}

	currentTime := time.Now().Format(time.RFC3339)
	_, err = db.Exec("UPDATE execs SET password =?, password_changed_at =? WHERE id =?", hashedNewPassword, currentTime, userId)
	if err != nil {
		return false, utils.ErrorHandler(err, "error updating password")
	}

	// token, err := utils.SignToken(userId, username, userRole)
	// if err != nil {
	// 	utils.ErrorHandler(err, "Password UPdated, error create token")
	// 	return
	// }
	return true, nil
}

func ForgotPasswordDbHandler(emailId string) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "database connection error")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow("SELECT id FROM execs WHERE email=?", emailId).Scan(&exec.ID)
	if err != nil {
		return utils.ErrorHandler(err, "user not found")
	}

	duration, err := time.ParseDuration(os.Getenv("RESET_TOKEN_EXP_DURATION"))
	if err != nil {
		return utils.ErrorHandler(err, "failed to send password reset email")
	}

	expiry := time.Now().Add(duration).Format(time.RFC3339)

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return utils.ErrorHandler(err, "failed to send password reset email")
	}

	token := hex.EncodeToString(tokenBytes)

	hashedToken := sha256.Sum256(tokenBytes)

	hashedTokenString := hex.EncodeToString(hashedToken[:])

	_, err = db.Exec("UPDATE execs SET password_reset_token =?, password_token_expires =? WHERE id =?", hashedTokenString, expiry, exec.ID)
	if err != nil {
		return utils.ErrorHandler(err, "failed to send password reset email")
	}

	//Send the reset email
	resetURL := fmt.Sprintf("https://localhost:3000/execs/resetpassword/reset/%s", token)
	message := fmt.Sprintf("Forgot your password? Reset your password using the following link:\n%s\n If you didn't request a password reset, please ignore this email. this Link is only valid for %d minutes", resetURL, int(duration.Minutes()))

	m := mail.NewMessage()
	m.SetHeader("From", "schooladmin@school.com") //
	m.SetHeader("To", emailId)
	m.SetHeader("Subject", "Your password reset link")
	m.SetBody("text/plain", message)

	d := mail.NewDialer("localhost", 1025, "", "")
	err = d.DialAndSend(m)
	if err != nil {
		return utils.ErrorHandler(err, "failed to send password reset email")
	}
	return nil
}

func ResetPasswordDbHandler(token, newPassword string) error {
	bytes, err := hex.DecodeString(token)
	if err != nil {
		return utils.ErrorHandler(err, "invalid reset token")
	}
	hashedToken := sha256.Sum256(bytes)
	hashedTokenString := hex.EncodeToString(hashedToken[:])

	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "database connection error")
	}
	defer db.Close()

	var user models.Exec

	query := "SELECT id, email FROM execs WHERE password_reset_token = ? AND password_token_expires > ?"
	err = db.QueryRow(query, hashedTokenString, time.Now().Format(time.RFC3339)).Scan(&user.ID, &user.Email)
	if err != nil {
		return utils.ErrorHandler(err, "invalid or expired reset token")
	}

	//Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to hash new password")
	}
	//Update the password in the database
	updateQuery := "UPDATE execs SET password =?, password_changed_at =?, password_reset_token = NULL, password_token_expires = NULL WHERE id = ?"
	_, err = db.Exec(updateQuery, hashedPassword, time.Now().Format(time.RFC3339), user.ID)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to update password")
	}
	return nil
}
