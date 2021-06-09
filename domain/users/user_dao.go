package users

import (
	"fmt"

	"github.com/bygirish/bookstore_users-api/datasources/mysql/users_db"
	"github.com/bygirish/bookstore_users-api/logger"
	"github.com/bygirish/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error while trying to prepare get user statement ", err)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {

		logger.Error("error while trying to get user by id", getErr)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error while trying to prepare save user statement ", err)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error while trying to exectue save user statement ", saveErr)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return mysql_utils.ParseError(saveErr)
	}

	// Other way of executing query on db straight way without preparing
	// result, err := users_db.Client.Exec(queryInsertUser, user.FirstName, user.LastName, user.Email, user.DateCreated)

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error while trying to get last interested id of save user result ", err)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return mysql_utils.ParseError(err)
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error while trying to prepare update user statement ", err)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error while trying to execture update user statement ", err)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error while trying to prepare delete user statement ", err)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error while trying to execture delete user statement ", err)
		return errors.NewInternalServerError(errors.NewError("database error").Error())
		// return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindBystatus(status string) ([]User, *errors.RestErr) {

	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error while trying to prepare find user by status statement ", err)
		return nil, errors.NewInternalServerError(errors.NewError("database error").Error())
		// return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error while trying to query find user by status statement ", err)
		return nil, errors.NewInternalServerError(errors.NewError("database error").Error())
		// return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error while trying to scan rows of find user by status query result", err)
			return nil, errors.NewInternalServerError(errors.NewError("database error").Error())
			// return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}
