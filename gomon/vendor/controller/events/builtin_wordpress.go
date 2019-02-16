package events

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func WPPage(url string) (state string, message string, metric string) {
	var pages interface{}

	resp, err := http.Get(url + "/wp-json/wp/v2/pages")
	if err != nil {
		return "CRITICAL", err.Error(), ""
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &pages)
		switch v := pages.(type) {
		case []interface{}:
			log.Printf("%v", len(v))
			if len(v) > 0 {
				return "OK", "", "pages=" + string(len(v))
			} else {
				return "WARNING", "No WP pages found", ""
			}
		default:
			return "WARNING", "No WP pages found", ""
		}
	}
	return "CRITICAL", "", ""
}
