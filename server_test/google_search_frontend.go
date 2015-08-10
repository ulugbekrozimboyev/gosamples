package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

var resultsTemplate = template.Must(template.New("results").Parse(`
<html>
<head/>
<body>
  <ol>
  {{range .Results}}
    <li>{{.Title}} - <a href="{{.URL}}">{{.URL}}</a></li>
  {{end}}
  </ol>
  <p>{{len .Results}} results in {{.Elapsed}}</p>
</body>
</html>
`))

func main() {
	http.HandleFunc("/search", handleSearch)
	fmt.Println("serving on http://localhost:8080/search")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

// handleSearch handles URLs like "/search?q=golang" by running a
// Google search for "golang" and writing the results as HTML to w.
func handleSearch(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)

	// Check the search query.
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, `missing "q" URL parameter`, http.StatusBadRequest)
		return
	}

	// Run the Google search.
	start := time.Now()
	results, err := Search(query)
	elapsed := time.Since(start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(query)

	if err := resultsTemplate.Execute(w, templateData{
		Results: results,
		Elapsed: elapsed,
	}); err != nil {
		log.Print(err)
		return
	}

}

func Search(query string) ([]Result, error) {

	// Prepare the Google Search API request.
	u, err := url.Parse("https://ajax.googleapis.com/ajax/services/search/web?v=1.0")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", query)
	u.RawQuery = q.Encode()

	// Issue the HTTP request and handle the response.
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var jsonResponse struct {
		ResponseData struct {
			Results []struct {
				TitleNoFormatting, URL string
			}
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		return nil, err
	}

	// Extract the Results from jsonResponse and return them.
	var results []Result
	for _, r := range jsonResponse.ResponseData.Results {
		results = append(results, Result{Title: r.TitleNoFormatting, URL: r.URL})
	}
	return results, nil
}

// Render the results.
type templateData struct {
	Results []Result
	Elapsed time.Duration
}

// A Result contains the title and URL of a search result.
type Result struct {
	Title, URL string
}

type error interface {
	Error() string // a useful human-readable error message
}

type Writer interface {
	Write(p []byte) (n int, err error)
}
