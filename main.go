package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var DB = Database{}

func main() {
	http.HandleFunc("/", webCacheHandler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

type RequestObj struct {
	ReqObj string `json:"req_obj"`
}

func webCacheHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cache server")
	d := json.NewDecoder(r.Body)
	ro := &RequestObj{}
	err := d.Decode(ro)
	if err != nil {
		log.Println(err)
	}
	cache, has := DB.Has(ro.ReqObj)
	log.Println("has: ", has)
	if has {
		body, valid, err := validateCache(cache, ro.ReqObj)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		if valid {
			w.Write(cache.obj)
		} else {
			w.Write(body)
		}
	} else {
		body, err := addToCache(ro.ReqObj)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}

}

func addToCache(raw string) (body []byte, err error) {
	resp, err := requestToServer(raw, -1)
	if err != nil {
		log.Println("error when requesting to server: ", err)
		return nil, err
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error when reading from server response")
		return body, nil
	}
	d := resp.Header.Get("Last-Modified")
	unixTime, err := time.Parse(http.TimeFormat, d)
	if err != nil {
		log.Println("error when parsing Last-Modified: skipped adding to cache")
		return body, nil
	}
	c := &Cache{
		obj:          body,
		lastModified: UnixTime(unixTime.Unix()),
	}
	DB.Add(raw, c)
	return body, nil
}

func revalidateCache(url string, body []byte, r *http.Response) {
	d := r.Header.Get("Last-Modified")
	unixTime, err := time.Parse(http.TimeFormat, d)
	if err != nil {
		return //do not add to cache: skip
	}
	c := &Cache{
		obj:          body,
		lastModified: UnixTime(unixTime.Unix()),
	}
	DB.Update(url, c)
}

func validateCache(c *Cache, url string) (body []byte, valid bool, err error) {
	lastModified := c.lastModified
	r, err := requestToServer(url, lastModified)
	if err != nil {
		return nil, false, err
	}

	if r.StatusCode == http.StatusNotModified { //respond from cache
		return nil, true, nil
	}

	//revalidate cache and return body of response
	body, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, false, err
	}
	defer r.Body.Close()
	revalidateCache(url, body, r)
	return body, false, nil
}

func requestToServer(raw string, u UnixTime) (*http.Response, error) {
	headers := map[string][]string{}
	if u != -1 {
		t := time.Unix(int64(u), 0).UTC()
		headers["If-Modified-Since"] = []string{t.Format(http.TimeFormat)}
	}
	log.Println("headers[\"If-Modified-Since\"] ", headers["If-Modified-Since"])
	link, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: "GET",
		URL:    link,
		Header: headers,
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
