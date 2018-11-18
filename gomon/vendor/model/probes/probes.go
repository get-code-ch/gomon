package probes

import (
	"controller"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
)

func GetProbes() ([]map[string]interface{}, error) {
	var (
		probes    []map[string]interface{}
		data      string
		jsonArray string
	)

	db, err := sql.Open("postgres", controller.Config.Database)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("select data from probe")
	defer rows.Close()

	jsonArray = "["
	row := rows.Next()
	for row {
		err := rows.Scan(&data)
		if err != nil {
			return nil, err
		}
		row = rows.Next()
		if row {
			jsonArray += data + ","
		} else {
			jsonArray += data + "]"
		}
	}
	json.Unmarshal([]byte(jsonArray), &probes)

	if err != nil {
		probes = nil
	}
	return probes, err
}
