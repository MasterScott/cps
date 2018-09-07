package properties

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	kv "cps/pkg/kv"

	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	gjson "github.com/tidwall/gjson"
)

func init() {
	// logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

type Error struct {
	Status string `json:"status"`
}

func GetProperties(w http.ResponseWriter, r *http.Request, account, region string) {
	vars := mux.Vars(r)
	scope := strings.Split(vars["scope"], "/")
	service := scope[0]
	fullPath := scope[1:len(scope)]
	log.Println(kv.Cache)

	var path bytes.Buffer
	path.WriteString("account/")
	path.WriteString(account)
	path.WriteString("/kubernetes/")
	path.WriteString(region)
	path.WriteString("/service/")
	path.WriteString(service)

	jsoni := kv.GetProperty(path.String())
	jb := jsoni.([]byte)

	b := new(bytes.Buffer)
	if err := json.Compact(b, jb); err != nil {
		log.Error(err)
	}

	j := []byte(b.Bytes())

	if len(fullPath) > 0 {
		f := strings.Join(fullPath, ".")
		p := gjson.GetBytes(j, "properties")
		selected := gjson.GetBytes([]byte(p.String()), f)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(strings.TrimSpace(selected.String())))
	} else {
		w.Header().Set("Content-Type", "application/json")
		p := gjson.GetBytes(j, "properties")
		w.Write([]byte(strings.TrimSpace(p.String())))
	}
}
