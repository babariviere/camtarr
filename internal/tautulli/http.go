package tautulli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func getActivity(uri, apiKey string) (respGetActivity, error) {
	u, err := url.Parse(uri + "/api/v2")
	if err != nil {
		log.Fatal("invalid uri for tautulli:", err)
	}

	q := u.Query()
	q.Set("apikey", apiKey)
	q.Set("cmd", "get_activity")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return respGetActivity{}, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respGetActivity{}, fmt.Errorf("invalid http status code: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var d respGetActivity
	if err = dec.Decode(&d); err != nil {
		return d, fmt.Errorf("cannot decode get_activity response: %w", err)
	}

	return d, nil
}
