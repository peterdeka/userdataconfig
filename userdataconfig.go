package userdataconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var DEFAULT_AWS_URL = "http://169.254.169.254/latest/user-data"

func GetVars(url *string) (map[string]string, error) {
	if url == nil {
		url = &DEFAULT_AWS_URL
	}
	resp, err := http.Get(*url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := string(htmlData)
	lines := strings.Split(str, "\n")
	varsLine := ""
	for i, v := range lines {
		if v == "#VARS" && i+1 < len(lines) {
			varsLine = lines[i+1]
			break
		}
	}
	if varsLine == "" {
		return nil, errors.New("VARS not found")
	}
	vars := map[string]string{}
	if err := json.Unmarshal([]byte(varsLine[1:]), &vars); err != nil {
		return nil, err
	}
	return vars, nil
}
