package main

import ( 
    "net/http"
    "encoding/json"
)

func listConfiguration(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(configuration); err != nil {
        panic(err)
    }
	
}
