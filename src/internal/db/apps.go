package db

type AppItem struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Logo  []byte `json:"logo"`
}

func GetApps() ([]AppItem, error) {
	rows, err := db.Query("SELECT id, name, color, logo FROM apps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []AppItem{}
	for rows.Next() {
		var app AppItem
		err = rows.Scan(&app.Id, &app.Name, &app.Color, &app.Logo)
		if err != nil {
			return nil, err
		}
		result = append(result, app)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func CreateApp(app AppItem) error {
	stmt, err := db.Prepare("INSERT INTO apps (name, color, logo) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(app.Name, app.Color, app.Logo)
	if err != nil {
		return err
	}

	return nil
}

func UpdateApp(app AppItem) error {
	stmt, err := db.Prepare("UPDATE apps SET name = ?, color = ?, logo = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(app.Name, app.Color, app.Logo, app.Id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteApp(id int) error {
	stmt, err := db.Prepare("DELETE FROM apps WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
