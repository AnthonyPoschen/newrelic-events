package events

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// 950kb (newrelic is 1MB max
// no sane person would have a single 50kb message???
// TODO: allow crazy things, because we are in a crazy world
const maxSize = 950000

///////////////////////////////////////////////////////////////////////////

type dataStore struct {
	*sync.Mutex
	Data string
}

///////////////////////////////////////////////////////////////////////////

func New(AccountID string, APIKey string) *Events {
	return &Events{
		Poster: StandardPost(http.DefaultClient),
		URL:    fmt.Sprintf("https://insights-collector.newrelic.com/v1/accounts/%s/events", AccountID),
		data: dataStore{
			Mutex: &sync.Mutex{},
			Data:  "",
		},
		key: APIKey,
	}
}

///////////////////////////////////////////////////////////////////////////

type Events struct {
	Poster func(req *http.Request) error

	data dataStore
	URL  string
	key  string
}

// Record will add the event to the queue of events that is thread safe, you can go Record
func (n *Events) Record(Name string, in map[string]interface{}) error {
	if Name == "" {
		return errors.New("No Event Name")
	}
	if in == nil {
		return errors.New("data is nil")
	}
	in["eventType"] = Name
	n.data.Lock()
	defer n.data.Unlock()
	leaderKey := ""
	if len(n.data.Data) > 0 {
		leaderKey = ","
	}
	marshledData, err := json.Marshal(in)
	if err != nil {
		return err
	}
	n.data.Data += fmt.Sprintf("%s%s", leaderKey, marshledData)

	if len(n.data.Data) > maxSize {
		// copy data into function so we can safely reuse the memory incase post is Async
		err = n._Post(n.data.Data)
		n.data.Data = ""
	}
	return err
}

// _Post is in charge of building the http Request and passing it on to the designated poster
func (n *Events) _Post(data string) error {
	// wrap the hand made json array correctly for posting (don't know a faster way to perform this logic)
	data = fmt.Sprintf("[%s]", data)
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	// reduce memory buffer usage by syncing through a channel as the content is read
	// to perform the request
	go func() {
		zipper := gzip.NewWriter(w)
		zipper.Write([]byte(data))
		zipper.Flush()
		w.Close()
		zipper.Close()
	}()
	req, err := http.NewRequest("POST", n.URL, r)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Insert-Key", n.key)
	req.Header.Add("Content-Encoding", "gzip")
	return n.Poster(req)
}

///////////////////////////////////////////////////////////////////////////

// Sync performs a force Post to newrelic disregarding waiting for max buffer size
func (n *Events) Sync() error {
	n.data.Lock()
	defer n.data.Unlock()
	err := n._Post(n.data.Data)
	if err != nil {
		return err
	}
	n.data.Data = ""
	return nil
}

///////////////////////////////////////////////////////////////////////////

func StandardPost(client *http.Client) func(*http.Request) error {
	return func(req *http.Request) error {
		ctx, canFunc := context.WithTimeout(context.Background(), time.Second*30)
		defer canFunc()
		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("Bad Response: %d - %s", resp.StatusCode, resp.Status)
		}
		return nil
	}
}

///////////////////////////////////////////////////////////////////////////

// AsyncPost is an en example poster that will run in a go routine and callback a function with the status if provided
// as the standard error won't be valuable
func AsyncPost(ctx context.Context, client http.Client, callback func(error)) func(*http.Request) error {
	return func(req *http.Request) error {
		req = req.WithContext(ctx)
		go func() {
			resp, err := client.Do(req)
			if err != nil {
				if callback != nil {
					callback(err)
				}
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				if callback != nil {
					callback(err)
				}
			}
			if callback != nil {
				callback(nil)
			}
			return
		}()
		return nil
	}
}
