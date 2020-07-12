package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	rs "remote_shutdown/remote_shutdown"
	"remote_shutdown/utils"
)

type RSRequest struct {
	Hosts []rs.RShutdown `json:"hosts"`
}

type RSResponse struct {
	Data string `json:"data"`
	Error string	`json:"error"`
}

func handleShutdown(securityCode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := utils.New()
		if r.Method != http.MethodPost {
			u.Log.Printf("API accessed with method %s. Sending HTTP 405", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			//w.WriteHeader(404)
			return
		}

		q := r.URL
		u.Log.Infof("q = %#v", q)

		sKey := r.URL.Query().Get("s")
		u.Log.Infof("skey = %s", sKey)

		output := &RSResponse{
			Data:  "",
			Error: "",
		}

		bodyBytes, bodyBytesErr := ioutil.ReadAll(r.Body)
		if bodyBytesErr != nil {
			u.Log.Errorf("Unable to read request Body. %s", bodyBytesErr.Error())
			output.Error = bodyBytesErr.Error()
			u.SendResponseJSON(w, output)
			return
		}

		var req RSRequest

		inputErr := json.Unmarshal(bodyBytes, &req)
		if inputErr != nil {
			u.Log.Errorf("Unable to parse request data. %s", inputErr.Error())
			output.Error = inputErr.Error()
			u.SendResponseJSON(w, output)
			return
		}

		for key, hostData := range req.Hosts {
			u.Log.Infof("key = %d", key)
			u.Log.Infof("HostData = %#v", hostData)
			err := hostData.Execute(sKey)

			//err := cmdShutdown()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				u.Log.Errorf("Failed to shutdown system. %s", err.Error())
				output.Error = "Failed to shutdown system. " + err.Error()
				u.SendResponseJSON(w, output)
				return
			}

			output.Data = "System is going to shutdown soon"
			u.Log.Infof(output.Data)
			u.SendResponseJSON(w, output)
		}
	}
}

// This application shutdowns the **Linux server** with **shutdown** command.
// For additional security, you should genereate some secret code (string) and start your application with it.
// You should provide this secret code in the request to `GET /shutdown` in **s** parameter.
func main() {
		securityCode := flag.String("sec-code", "", "Security code")
		port         := flag.String("port", "9898", "Port to listen")
		flag.Parse()

		u := utils.New()
		u.Log.Infof("Security code: %s", *securityCode)
		http.HandleFunc("/shutdown", handleShutdown(*securityCode))
		log.Fatal(http.ListenAndServe(":"+*port, nil))
}
