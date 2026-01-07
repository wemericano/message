package dbcall

import (
	"database/sql"
	"messanger/model"
)

var db *sql.DB

// InitDB DB 연결 초기화
func InitDB(database *sql.DB) {
	db = database
}

// CreateUser 사용자 생성
func CreateUser(user *model.User) error {
	query := `INSERT INTO users (username, password, name, email) VALUES (@p1, @p2, @p3, @p4)`
	_, err := db.Exec(query, user.Username, user.Password, user.Name, user.Email)
	return err
}

// GetUserByUsername username으로 사용자 조회
func GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, password, name, email FROM users WHERE username = @p1`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// SearchUsersByName 이름으로 사용자 검색
func SearchUsersByName(name string) ([]*model.User, error) {
	query := `SELECT id, username, password, name, email FROM users WHERE name LIKE @p1`
	rows, err := db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
