package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/creepfmd/jsonpath"
	"github.com/gorilla/mux"
)

// our main function
func main() {
	log.Println("Setting router...")
	router := mux.NewRouter()
	router.HandleFunc("/webhook/replaceValue/", replaceValue).
		Queries("param1", "{param1}").
		Queries("param2", "{param2").
		Methods("POST")
	router.HandleFunc("/webhook/replaceKey/", replaceKey).
		Queries("param1", "{param1}").
		Queries("param2", "{param2").
		Methods("POST")
	router.HandleFunc("/webhook/addPrefix/", addPrefix).
		Queries("param1", "{param1}").
		Queries("param2", "{param2").
		Methods("POST")
	router.HandleFunc("/webhook/addSuffix/", addSuffix).
		Queries("param1", "{param1}").
		Queries("param2", "{param2").
		Methods("POST")
	router.HandleFunc("/webhook/calculate/", calculate).
		Queries("param1", "{param1}").
		Queries("param2", "{param2").
		Queries("param3", "{param3").
		Methods("POST")
	log.Println("Listening 8081...")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func replaceValue(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	jsonData := getJsonObject(w, r)

	jpath, _ := jsonpath.Compile((string)(params["param1"]))
	jpathSteps := jpath.GetSteps()

	res, _ := jpath.Lookup(jsonData)

	clearBody, _ := json.Marshal(jsonData)
	switch x := res.(type) {
	case []interface{}:
		nodeName := jpathSteps[len(jpathSteps)-1]
		rgxp := regexp.MustCompile(`"` + nodeName + `":\[.*\]`)
		var messageParts []string
		for _, e := range x {
			dummy, _ := json.Marshal(e)
			messageParts = append(messageParts, rgxp.ReplaceAllString((string)(clearBody[:]), `"`+nodeName+`":`+(string)(dummy[:])))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[` + strings.Join(messageParts, `,`) + `]`))
	default:
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(clearBody))
	}
}

func replaceKey(w http.ResponseWriter, r *http.Request) {
}

func addPrefix(w http.ResponseWriter, r *http.Request) {
}

func addSuffix(w http.ResponseWriter, r *http.Request) {
}

func calculate(w http.ResponseWriter, r *http.Request) {
}

func getJsonObject(w http.ResponseWriter, r *http.Request) interface{} {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return nil
	}

	var jsonData interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return nil
	}

	return jsonData
}
