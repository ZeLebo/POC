package DB

import (
	u "backend/pkg/User"
	"database/sql"
	"errors"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SQL struct {
	Store *sql.DB
}

func openDataBase() *sql.DB {
	db, err := sql.Open("sqlite3", "pkg/DB/StorageData.db")
	if err != nil {
		log.Error(err.Error())
	}
	return db
}

// NewSQLDataBase creates the database and connects to it
func NewSQLDataBase() *SQL {
	var database SQL
	database.Store = openDataBase()
	return &database
}

func (database *SQL) containsUser(userId int) bool {
	rows, err := database.Store.Query(`select count(UserId) from users where UserID = ?`, userId)
	if err != nil {
		log.Error(err.Error())
	}
	var contain int
	for rows.Next() {
		_ = rows.Scan(&contain)
	}
	if contain == 0 {
		return false
	}
	return true
}

func (database *SQL) containsTag(tagID string) bool {
	rows, err := database.Store.Query(`select count(TagID) from Notes where TagID = ?`, tagID)
	if err != nil {
		log.Error(err.Error())
	}
	var contain int
	for rows.Next() {
		_ = rows.Scan(&contain)
	}
	if contain == 0 {
		return false
	}
	return true
}

// GetUserID return UserID, if user with such login and password exist
func (database *SQL) GetUserID(login, password string) (int, error) {
	// get login
	rows, err := database.Store.Query(`
	select count(UserID) from Users where Login = ?;
	`, login)

	if err != nil {
		return 0, err
	}

	var cnt, UserID int

	for rows.Next() {
		err = rows.Scan(&cnt)
	}

	if cnt == 0 {
		return 0, errors.New("no such user")
	}

	rows, err = database.Store.Query(`
	select count(UserID), UserID from Users where Login = ? and Password = ?;
	`, login, password)

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		err = rows.Scan(&cnt, &UserID)
		// if no userID provided makes error with parsing NULL to int
		if err != nil {
			return 0, errors.New("login or password is incorrect")
		}
	}

	return UserID, nil
}

// GetAllUsers for getting all users from database
// Return []string for answering the request and error status
func (database *SQL) GetAllUsers() ([]string, error) {
	rows, err := database.Store.Query(`select UserID from "Users"`)
	if err != nil {
		return nil, err
	}

	var UserID int
	var items []u.User

	for rows.Next() {
		err = rows.Scan(&UserID)
		if err != nil {
			return nil, err
		}

		items = append(items, u.User{
			UserID: UserID,
		})
	}

	var response []string
	for _, id := range items {
		response = append(response, strconv.Itoa(id.UserID))
	}
	return response, nil
}

// CreateUser creates a user in db
// Returns created user and error status
func (database *SQL) CreateUser(login, password string) (u.User, error) {

	rows, err := database.Store.Query(`select count(UserID) from Users`)
	if err != nil {
		return u.User{}, err
	}
	var userId int

	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return u.User{}, err
		}
	}

	var cnt int

	rows, err = database.Store.Query(`select count(Login) from Users where Login = ?`, login)
	if err != nil {
		return u.User{}, err
	}

	for rows.Next() {
		err = rows.Scan(&cnt)
		if err != nil {
			return u.User{}, err
		}
	}

	if cnt != 0 {
		return u.User{}, errors.New("user with such login exists")
	}

	stmt, err := database.Store.Prepare(`insert into Users (UserID, Login, Password) values (?, ?, ?)`)
	if err != nil {
		return u.User{}, err
	}
	_, err = stmt.Exec(userId, login, password)
	if err != nil {
		return u.User{}, err
	}
	log.Info("New user: ", userId)
	return u.User{
		UserID: userId,
	}, nil
}

// ChangeLogin change login of user
func (database *SQL) ChangeLogin(userId int, login string) error {
	stmt, err := database.Store.Prepare(`update Users set Login = ? where UserID = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(login, userId)
	if err != nil {
		return err
	}
	return nil
}

// ChangePassword change password of user
func (database *SQL) ChangePassword(userId int, password string) error {
	stmt, err := database.Store.Prepare(`update Users set Password = ? where UserID = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(password, userId)
	if err != nil {
		return err
	}
	return nil
}

// GetUserTags get all tags from specific user
// Return []string for answering the request and error status
func (database *SQL) GetUserTags(userId int) ([]string, error) {
	if !database.containsUser(userId) {
		return nil, errors.New("no such user")
	}
	rows, err := database.Store.Query(`select distinct TagId from Notes where UserID = ?`, userId)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var tagID string
	var result []u.Tag

	for rows.Next() {
		err = rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}

		result = append(result, u.Tag{
			TagID: tagID,
		})
	}

	var response []string
	for _, tag := range result {
		response = append(response, tag.TagID)
	}
	return response, nil
}

// GetUserNotes get user notes from tag
// Return []string for answering the request and error status
func (database *SQL) GetUserNotes(userId int, tagId string) ([]u.Note, error) {
	if !database.containsUser(userId) {
		return nil, errors.New("no such user")
	}

	if !database.containsTag(tagId) {
		return nil, errors.New("no such tag")
	}

	rows, err := database.Store.Query(
		`select note, data from Notes where UserID = ? and notes.TagID = ?`, userId, tagId)
	if err != nil {
		return nil, err
	}
	var note u.Note
	var result []u.Note

	for rows.Next() {
		err = rows.Scan(&note.Note, &note.Time)
		if err != nil {
			return nil, err
		}
		result = append(result, note)
	}
	return result, nil
}

// AddNote creates a new note for tag
// Creates tag if not exist
// Return Tag object and error status
func (database *SQL) AddNote(userId int, tagId, noteInfo string) (u.Tag, error) {
	if !database.containsUser(userId) {
		return u.Tag{}, errors.New("no such user")
	}

	stmt, err := database.Store.Prepare(`insert into Notes (UserID, TagID, Note, Data)  values (?, ?, ?, ?)`)
	if err != nil {
		return u.Tag{}, err
	}

	if _, err = stmt.Exec(userId, tagId, noteInfo, time.Now()); err != nil {
		return u.Tag{}, err
	}

	rows, err := database.Store.Query(
		`select note, data from Notes where UserID = ? and notes.TagID = ?`, userId, tagId)
	if err != nil {
		return u.Tag{}, err
	}

	var note u.Note
	var result u.Tag

	result.TagID = tagId

	for rows.Next() {
		err = rows.Scan(&note.Note, &note.Time)
		if err != nil {
			return u.Tag{}, err
		}
		result.Notes = append(result.Notes, note)
	}

	return result, nil
}
