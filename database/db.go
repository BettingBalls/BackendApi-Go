package database

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"fmt"
)
func SupabaseRequest(method, path string, body interface{}) (*http.Response, error) {

	fmt.Println("SUPABASE_KEY:", os.Getenv("SUPABASE_KEY"))
	fmt.Println("SUPABASE_URL:", os.Getenv("SUPABASE_URL"))

	var buffer *bytes.Buffer

	if body != nil {
		b, _ := json.Marshal(body)
		buffer = bytes.NewBuffer(b)
	} else {
		buffer = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(
		method,
		os.Getenv("SUPABASE_URL")+"/rest/v1"+path,
		buffer,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	return http.DefaultClient.Do(req)
}