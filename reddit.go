// Package reddit implements a basic client for the Reddit API.

package reddit

import (
    "log"
	"fmt"
	"errors"
    "net/http"
	"encoding/json"
)

// response describes the JSON-formatted response returned by the Reddit API.
type response struct {
    Data struct {
        Children []struct {
            Data Item
        }
    }
}

// Item describes a Reddit item.
type Item struct {
    Title 	 string
    URL   	 string
	Comments int `json:"num_comments"`
}

// String returns a formatted string representation of an Item.
func (i Item) String() string {
	com := ""
	switch i.Comments {
	case 0:
		// nothing
	case 1:
		com = " (1 comment)"
	default:
		com = fmt.Sprintf(" (%d comments)", i.Comments)
	}
	return fmt.Sprintf("%s%s\n%s", i.Title, com, i.URL)
}

// Get fetches the most recent Items posted to a specific subreddit.
func Get(reddit string) ([]Item, error) {
    url := fmt.Sprintf("http://reddit.com/r/%s.json", reddit)
	resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
	defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New(resp.Status)
    }

	r := new(response)
    err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	items := make([]Item, len(r.Data.Children))

	for i, item := range r.Data.Children {
		items[i] = item.Data
	}

	return items, nil
}
