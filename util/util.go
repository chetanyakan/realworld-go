package util

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
)

// DumpRequest returns the dump of request
func DumpRequest(r *http.Request) string {
	d, err := httputil.DumpRequest(r, r.Method == http.MethodPost)
	if err != nil {
		d = []byte{}
	}

	return string(d)
}

// writeJSON writes a JSON Content-Type header and a JSON-encoded object to the
// http.ResponseWriter.
func writeJSON(w http.ResponseWriter, v interface{}) error {
	// Indent the JSON so it's easier to read for hackers.
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	return err
}
