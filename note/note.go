package note

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func MyIP() map[string]string {
	m := make(map[string]string)
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		//return err.Error()
		fmt.Sprintln(err)
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		//return err.Error()
		fmt.Sprintln("error")
	}
	//var ip IP
	json.Unmarshal(body, &m)
	return m
}
