package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path    string
	Method  string
	Params  map[string]string
	Queries map[string]string

	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	c := Context{
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		Writer: w,
	}
	c.parseQuery()
	return &c
}

func (c *Context) parseQuery() {
	c.Queries = make(map[string]string)

	URL := c.Req.URL.String()
	index := strings.Index(URL, "?")
	if index == -1 || len(URL) <= index+1 {
		return
	}

	query := URL[index+1:]
	queries := strings.Split(query, "&")
	// simplified implementation, can cause bugs when req contains multiple '='s in a single subquery.
	for _, subQuery := range queries {
		q := strings.Split(subQuery, "=")
		if q[0] == "" {
			continue
		}
		if len(q) < 2 {
			if _, ok := c.Queries[q[0]]; !ok {
				c.Queries[q[0]] = ""
				continue
			}
		}
		c.Queries[q[0]] = q[1]
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status(code int) {
	c.Writer.Header().Add("Code", strconv.Itoa(code))
}

func (c *Context) HTML(statusCode int, file string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(statusCode)
	c.Writer.Write([]byte(file))
}

func (c *Context) String(statusCode int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(statusCode)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) JSON(statusCode int, v interface{}) error {
	c.SetHeader("Content-Type", "application/json")
	c.Status(statusCode)
	json, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = c.Writer.Write(json)
	return err
}

func (c *Context) Query(query string) (result string) {
	return c.Queries[query]
}

func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}
