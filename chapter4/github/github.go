package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const issueURL = "https://api.github.com/search/issues"

// https://docs.github.com/en/rest/search/search?apiVersion=2026-03-10#search-issues-and-pull-requests

// When json.Decoder reads the JSON, it processes it as a stream of key-value pairs — order doesn't matter at all. For each key it encounters, it looks up the matching struct field using this priority:

// 1. A field with an exact json tag match        → json:"total_count"
// 2. A field with a case-insensitive name match  → "TotalCount" matches "totalcount"
// 3. If no match found → the field is silently ignored

type IssueSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Titttle   string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func SearchIssues(terms []string) (*IssueSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	res, err := http.Get(issueURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("Search query failed: %s", res.Status)
	}

	var result IssueSearchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		res.Body.Close()
		return nil, err
	}
	res.Body.Close()
	return &result, err
}
