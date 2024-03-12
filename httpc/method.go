package httpc

import (
	"encoding/json"
)

type H map[string]string

// GET

func Get(url string, headers H) ([]byte, error) {

	c := Client{"GET", url, "", headers, 0}
	return c.Request()

}

func TextGet(url string, headers H) (string, error) {

	c := Client{"GET", url, "", headers, 0}
	return c.TextRequest()

}

// POST

func Post(url, query string, headers H) ([]byte, error) {

	c := Client{"POST", url, query, headers, 0}
	return c.Request()

}

func JsonPost(url string, query any, headers H) ([]byte, error) {

	data, err := json.Marshal(query)

	if err != nil {
		return nil, err
	}

	c := Client{"POST", url, string(data), headers, 0}
	return c.JsonRequest()

}

func TextPost(url string, query string, headers H) (string, error) {

	c := Client{"POST", url, query, headers, 0}
	return c.TextRequest()

}
