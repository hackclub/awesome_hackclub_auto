package couchdb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	URI      string
	user     string
	password string
}

type Database struct {
	client *Client
	Name   string
}

func NewClient(uri, user, password string) *Client {
	return &Client{
		URI:      uri,
		user:     user,
		password: password,
	}
}

func (client *Client) Database(name string) *Database {
	return &Database{
		client: client,
		Name:   name,
	}
}

func (database *Database) Get(id string, v interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", database.client.URI, database.Name, id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", database.client.authHeader())
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return errors.New(string(result))
	}

	err = json.Unmarshal(result, v)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) authHeader() string {
	buf := bytes.NewBuffer(nil)
	writer := base64.NewEncoder(base64.StdEncoding, buf)
	_, err := writer.Write([]byte(client.user + ":" + client.password))
	if err != nil {
		return ""
	}
	writer.Close()
	return "Basic " + buf.String()
}
