package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	"strings"
)

var DB *sql.DB

func init() {

	viper.Set("MYSQL_DSN", "root:root@tcp(159.65.153.232:3306)/eventackle")

	var err error
	DB, err = sql.Open("mysql", viper.GetString("MYSQL_DSN"))
	if err != nil {
		logrus.Errorln("failed to open connection to MYSQL")
		logrus.Errorln("error: ", err)
	}
}

func GetPositions(numberOfEmails int, emails []interface{}) (map[string]string, error) {

	query := `SELECT  CONCAT (J.title, ', ', J.company) as position, U.email as email
					FROM users U
					INNER JOIN user_profiles P 
					ON P.user_id = U.id
					INNER JOIN user_jobs J
					ON J.id = P.job_id
					WHERE U.email IN (?` + strings.Repeat(",?", numberOfEmails-1) + `)
				`
	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(emails...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var position, email string
	positions := make(map[string]string, 0)

	for rows.Next() {
		err := rows.Scan(&position, &email)
		if err != nil {
			return nil, err
		}
		positions[email] = position
	}
	return positions, nil
}

func GetOrganizerEmail(id string)(string, error){
	query := `SELECT users.email
			FROM organizer_profiles
			INNER JOIN users
			ON organizer_profiles.user_id = users.id
			WHERE organizer_profiles.id = ?
	`
	var email string
	err := DB.QueryRow(query, id).Scan(&email)
	if err != nil {
		logrus.Error("error:", err)
		return email, err
	}
	return email, nil

}