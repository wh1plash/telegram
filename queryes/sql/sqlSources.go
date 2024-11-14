package sql

import "io/ioutil"

func ReadSQLFile(file string) string {
	content, _ := ioutil.ReadFile(file)
	return string(content)
}
