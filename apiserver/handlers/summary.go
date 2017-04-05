package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

//openGraphPrefix is the prefix used for Open Graph meta properties
const openGraphPrefix = "og:"

//openGraphProps represents a map of open graph property names and values
type openGraphProps map[string]string

func getPageSummary(url string) (openGraphProps, error) {
	//Get the URL
	//If there was an error, return it

	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error fetching the URL: %v", err)
	}

	//ensure that the response body stream is closed eventually
	//HINTS: https://gobyexample.com/defer
	//https://golang.org/pkg/net/http/#Response

	defer resp.Body.Close()

	//if the response StatusCode is >= 400
	//return an error, using the response's .Status
	//property as the error message

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("response status was %s", resp.Status)
	}

	//if the response's Content-Type header does not
	//start with "text/html", return an error noting
	//what the content type was and that you were
	//expecting HTML

	cType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(cType, "text/html") {
		log.Fatalf("response content type was %s and not text/html\n", cType)
		return nil, err
	}

	//create a new openGraphProps map instance to hold
	//the Open Graph properties you find
	//(see type definition above)

	body := make(openGraphProps)

	//tokenize the response body's HTML and extract
	//any Open Graph properties you find into the map,
	//using the Open Graph property name as the key, and the
	//corresponding content as the value.
	//strip the openGraphPrefix from the property name before
	//you add it as a new key, so that the key is just `title`
	//and not `og:title` (for example).

	//HINTS: https://info344-s17.github.io/tutorials/tokenizing/
	//https://godoc.org/golang.org/x/net/html

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
		//done iterating over the url and can leave the loop
		if tokenType == html.ErrorToken {
			break
		}

		//meta props begin with start tag or they are self closing
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			//this gets the whole tag
			token := tokenizer.Token()
			if "meta" == token.Data {
				tokenProp := token.Attr[0]
				if "property" == tokenProp.Key {
					//gets the meta prop and the value of it--breaks outta the loop
					prop := strings.Split(token.Attr[0].Val, ":")
					if "og" == prop[0] {
						switch prop[1] {
						case
							"url",
							"title",
							"description",
							"image":
							//ensures that og:image:width is not received
							if len(prop) == 2 {
								body[prop[1]] = token.Attr[1].Val
							}
						default:
						}
					}
				}
			}

		}
	}

	if len(body) != 0 {
		return body, nil
	}
	return nil, err

}

//SummaryHandler fetches the URL in the `url` query string parameter, extracts
//summary information about the returned page and sends those summary properties
//to the client as a JSON-encoded object.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	//Add the following header to the response
	//   Access-Control-Allow-Origin: *
	//this will allow JavaScript served from other origins
	//to call this API

	w.Header().Add("Access-Control-Allow-Origin", "*")

	//get the `url` query string parameter
	//if you use r.FormValue() it will also handle cases where
	//the client did POST with `url` as a form field
	//HINT: https://golang.org/pkg/net/http/#Request.FormValue

	URL := r.FormValue("url")

	//if no `url` parameter was provided, respond with
	//an http.StatusBadRequest error and return
	//HINT: https://golang.org/pkg/net/http/#Error

	if URL == "" {
		http.Error(w, "bad response no URL", http.StatusBadRequest)
		return
	}

	//call getPageSummary() passing the requested URL
	//and holding on to the returned openGraphProps map
	//(see type definition above)

	ogProps, err := getPageSummary(URL)

	//if you get back an error, respond to the client
	//with that error and an http.StatusBadRequest code
	if err != nil {
		http.Error(w, "bad request when getting summary", http.StatusBadRequest)
		return
	}
	//otherwise, respond by writing the openGrahProps
	//map as a JSON-encoded object
	//add the following headers to the response before
	//you write the JSON-encoded object:
	//   Content-Type: application/json; charset=utf-8
	//this tells the client that you are sending it JSON

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	jsonProp, err := json.Marshal(ogProps)
	if err != nil {
		http.Error(w, "error encoding JSON: "+err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonProp)
}
