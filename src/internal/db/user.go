package db

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func GetUserByName(username string) (*User, error) {
	var user User

	err := db.QueryRow("SELECT id, username, password from users WHERE username = ? LIMIT 1", username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserById(userId int) (*User, error) {
	var user User

	err := db.QueryRow("SELECT id, username, password from users WHERE id = ?", userId).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(userId int, username string) error {
	stmt, err := db.Prepare("UPDATE users SET username = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, userId)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePassword(userId int, password string) error {
	stmt, err := db.Prepare("UPDATE users SET password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(password, userId)
	if err != nil {
		return err
	}

	return nil
}
