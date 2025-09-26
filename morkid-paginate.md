This file is a merged representation of the entire codebase, combined into a single document by Repomix.
The content has been processed where security check has been disabled.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Security check has been disabled - content may contain sensitive information
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
.circleci/
  config.yml
.github/
  workflows/
    go.yml
.gitignore
.travis.yml
go.mod
LICENSE
paginate_test.go
paginate.go
README.md
```

# Files

## File: .circleci/config.yml
````yaml
version: 2
jobs:
  build:
    docker:
      - image: golang:1.16-buster
    working_directory: /go/src/github.com/morkid/paginate
    steps:
      - checkout
      - run: go test -v
````

## File: .github/workflows/go.yml
````yaml
name: Go

on:
  push:
    branches: [ master ]
    paths:
    - '**.go'
    - '**.mod'
    - '**.sum'
  pull_request:
    branches: [ master ]

jobs:

  build:
    name: Build
    runs-on: ubuntu-latest
    steps:

    - name: Set up Go 1.x
      uses: actions/setup-go@v2
      with:
        go-version: ^1.13

    - name: Check out code into the Go module directory
      uses: actions/checkout@v2

    - name: Test
      run: go test -v
````

## File: .gitignore
````
/*
!.gitignore
!LICENSE
!*.go
!*.mod
!*.sum
!*.md
!*.yml
!/.circleci/
!/.github/
````

## File: .travis.yml
````yaml
language: go
os: linux
dist: xenial
go: 1.x
script: go test -v
````

## File: go.mod
````
module github.com/morkid/paginate

go 1.16

require (
	github.com/iancoleman/strcase v0.1.3
	github.com/klauspost/compress v1.11.12 // indirect
	github.com/morkid/gocache v1.0.0
	github.com/valyala/fasthttp v1.22.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.3
)
````

## File: LICENSE
````
MIT License

Copyright (c) 2020 morkid

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

## File: paginate_test.go
````go
package paginate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/morkid/gocache"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var format = "%s doesn't match. Expected: %v, Result: %v"

func TestGetNetHttp(t *testing.T) {
	size := int64(20)
	page := int64(1)
	sort := "user.name,-id"
	avg := "seventy %"

	queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"]]`, avg)
	query := fmt.Sprintf(`page=%d&size=%d&sort=%s&filters=%s`, page, size, sort, url.QueryEscape(queryFilter))

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}

	parsed := parseRequest(req, Config{})
	if parsed.Size != size {
		t.Errorf(format, "Size", size, parsed.Size)
	}
	if parsed.Page != page {
		t.Errorf(format, "Page", page, parsed.Page)
	}
	if len(parsed.Sorts) != 2 {
		t.Errorf(format, "Sort length", 2, len(parsed.Sorts))
	} else {
		if parsed.Sorts[0].Column != "user.name" {
			t.Errorf(format, "Sort field 0", "user.name", parsed.Sorts[0].Column)
		}
		if parsed.Sorts[0].Direction != "ASC" {
			t.Errorf(format, "Sort direction 0", "ASC", parsed.Sorts[0].Direction)
		}
		if parsed.Sorts[1].Column != "id" {
			t.Errorf(format, "Sort field 1", "id", parsed.Sorts[1].Column)
		}
		if parsed.Sorts[1].Direction != "DESC" {
			t.Errorf(format, "Sort direction 1", "DESC", parsed.Sorts[1].Direction)
		}
	}

	filters, ok := parsed.Filters.Value.([]pageFilters)
	if ok {
		if filters[0].Column != "user.average_point" {
			t.Errorf(format, "Filter field for user.average_point", "user.average_point", filters[0].Column)
		}
		if filters[0].Operator != "LIKE" {
			t.Errorf(format, "Filter operator for user.average_point", "LIKE", filters[0].Operator)
		}
		value, isValid := filters[0].Value.(string)
		expected := "%" + avg + "%"
		if !isValid || value != expected {
			t.Errorf(format, "Filter operator for user.average_point", expected, value)
		}
	} else {
		t.Log(parsed.Filters)
		t.Errorf(format, "pageFilters class", "paginate.pageFilters", "null")
	}
}
func TestGetFastHttp(t *testing.T) {
	size := int64(20)
	page := int64(1)
	sort := "user.name,-id"
	avg := "seventy %"

	queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"]]`, avg)
	query := fmt.Sprintf(`page=%d&size=%d&sort=%s&filters=%s`, page, size, sort, url.QueryEscape(queryFilter))

	req := &fasthttp.Request{}
	req.Header.SetMethod("GET")
	req.URI().SetQueryString(query)

	parsed := parseRequest(req, Config{})
	if parsed.Size != size {
		t.Errorf(format, "Size", size, parsed.Size)
	}
	if parsed.Page != page {
		t.Errorf(format, "Page", page, parsed.Page)
	}
	if len(parsed.Sorts) != 2 {
		t.Errorf(format, "Sort length", 2, len(parsed.Sorts))
	} else {
		if parsed.Sorts[0].Column != "user.name" {
			t.Errorf(format, "Sort field 0", "user.name", parsed.Sorts[0].Column)
		}
		if parsed.Sorts[0].Direction != "ASC" {
			t.Errorf(format, "Sort direction 0", "ASC", parsed.Sorts[0].Direction)
		}
		if parsed.Sorts[1].Column != "id" {
			t.Errorf(format, "Sort field 1", "id", parsed.Sorts[1].Column)
		}
		if parsed.Sorts[1].Direction != "DESC" {
			t.Errorf(format, "Sort direction 1", "DESC", parsed.Sorts[1].Direction)
		}
	}

	filters, ok := parsed.Filters.Value.([]pageFilters)
	if ok {
		if filters[0].Column != "user.average_point" {
			t.Errorf(format, "Filter field for user.average_point", "user.average_point", filters[0].Column)
		}
		if filters[0].Operator != "LIKE" {
			t.Errorf(format, "Filter operator for user.average_point", "LIKE", filters[0].Operator)
		}
		value, isValid := filters[0].Value.(string)
		expected := "%" + avg + "%"
		if !isValid || value != expected {
			t.Errorf(format, "Filter operator for user.average_point", expected, value)
		}
	} else {
		t.Log(parsed.Filters)
		t.Errorf(format, "pageFilters class", "paginate.pageFilters", "null")
	}
}

func TestPostNetHttp(t *testing.T) {
	size := int64(20)
	page := int64(1)
	sort := "user.name,-id"
	avg := "seventy %"

	data := `
		{
			"page": %d,
			"size": %d,
			"sort": "%s",
			"filters": %s
		}
	`

	queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"]]`, avg)
	query := fmt.Sprintf(data, page, size, sort, queryFilter)

	body := io.NopCloser(bytes.NewReader([]byte(query)))

	req := &http.Request{
		Method: "POST",
		Body:   body,
	}

	parsed := parseRequest(req, Config{})
	if parsed.Size != size {
		t.Errorf(format, "Size", size, parsed.Size)
	}
	if parsed.Page != page {
		t.Errorf(format, "Page", page, parsed.Page)
	}
	if len(parsed.Sorts) != 2 {
		t.Errorf(format, "Sort length", 2, len(parsed.Sorts))
	} else {
		if parsed.Sorts[0].Column != "user.name" {
			t.Errorf(format, "Sort field 0", "user.name", parsed.Sorts[0].Column)
		}
		if parsed.Sorts[0].Direction != "ASC" {
			t.Errorf(format, "Sort direction 0", "ASC", parsed.Sorts[0].Direction)
		}
		if parsed.Sorts[1].Column != "id" {
			t.Errorf(format, "Sort field 1", "id", parsed.Sorts[1].Column)
		}
		if parsed.Sorts[1].Direction != "DESC" {
			t.Errorf(format, "Sort direction 1", "DESC", parsed.Sorts[1].Direction)
		}
	}

	filters, ok := parsed.Filters.Value.([]pageFilters)
	if ok {
		if filters[0].Column != "user.average_point" {
			t.Errorf(format, "Filter field for user.average_point", "user.average_point", filters[0].Column)
		}
		if filters[0].Operator != "LIKE" {
			t.Errorf(format, "Filter operator for user.average_point", "LIKE", filters[0].Operator)
		}
		value, isValid := filters[0].Value.(string)
		expected := "%" + avg + "%"
		if !isValid || value != expected {
			t.Errorf(format, "Filter operator for user.average_point", expected, value)
		}
	} else {
		t.Log(parsed.Filters)
		t.Errorf(format, "pageFilters class", "paginate.pageFilters", "null")
	}
}
func TestPostFastHttp(t *testing.T) {
	size := int64(20)
	page := int64(1)
	sort := "user.name,-id"
	avg := "seventy %"

	data := `
		{
			"page": %d,
			"size": %d,
			"sort": "%s",
			"filters": %s
		}
	`

	queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"]]`, avg)
	query := fmt.Sprintf(data, page, size, sort, queryFilter)

	req := &fasthttp.Request{}
	req.Header.SetMethod("POST")
	req.SetBodyString(query)

	parsed := parseRequest(req, Config{})
	if parsed.Size != size {
		t.Errorf(format, "Size", size, parsed.Size)
	}
	if parsed.Page != page {
		t.Errorf(format, "Page", page, parsed.Page)
	}
	if len(parsed.Sorts) != 2 {
		t.Errorf(format, "Sort length", 2, len(parsed.Sorts))
	} else {
		if parsed.Sorts[0].Column != "user.name" {
			t.Errorf(format, "Sort field 0", "user.name", parsed.Sorts[0].Column)
		}
		if parsed.Sorts[0].Direction != "ASC" {
			t.Errorf(format, "Sort direction 0", "ASC", parsed.Sorts[0].Direction)
		}
		if parsed.Sorts[1].Column != "id" {
			t.Errorf(format, "Sort field 1", "id", parsed.Sorts[1].Column)
		}
		if parsed.Sorts[1].Direction != "DESC" {
			t.Errorf(format, "Sort direction 1", "DESC", parsed.Sorts[1].Direction)
		}
	}

	filters, ok := parsed.Filters.Value.([]pageFilters)
	if ok {
		if filters[0].Column != "user.average_point" {
			t.Errorf(format, "Filter field for user.average_point", "user.average_point", filters[0].Column)
		}
		if filters[0].Operator != "LIKE" {
			t.Errorf(format, "Filter operator for user.average_point", "LIKE", filters[0].Operator)
		}
		value, isValid := filters[0].Value.(string)
		expected := "%" + avg + "%"
		if !isValid || value != expected {
			t.Errorf(format, "Filter operator for user.average_point", expected, value)
		}
	} else {
		t.Errorf(format, "pageFilters class", "paginate.pageFilters", "null")
	}
}

func TestProgrammaticallyPaginate(t *testing.T) {
	size := int64(20)
	page := int64(1)
	sort := "user.name,-id"
	avg := "seventy %"

	req := &Request{
		Page: page,
		Size: size,
		Sort: sort,
		Filters: []interface{}{
			[]interface{}{"user.average_point", "like", avg},
			[]interface{}{"and"},
			[]interface{}{"user.average_point", "is not", nil},
		},
	}

	parsed := parseRequest(req, Config{})
	if parsed.Size != size {
		t.Errorf(format, "Size", size, parsed.Size)
	}
	if parsed.Page != page {
		t.Errorf(format, "Page", page, parsed.Page)
	}
	if len(parsed.Sorts) != 2 {
		t.Errorf(format, "Sort length", 2, len(parsed.Sorts))
	} else {
		if parsed.Sorts[0].Column != "user.name" {
			t.Errorf(format, "Sort field 0", "user.name", parsed.Sorts[0].Column)
		}
		if parsed.Sorts[0].Direction != "ASC" {
			t.Errorf(format, "Sort direction 0", "ASC", parsed.Sorts[0].Direction)
		}
		if parsed.Sorts[1].Column != "id" {
			t.Errorf(format, "Sort field 1", "id", parsed.Sorts[1].Column)
		}
		if parsed.Sorts[1].Direction != "DESC" {
			t.Errorf(format, "Sort direction 1", "DESC", parsed.Sorts[1].Direction)
		}
	}

	t.Log(parsed.Filters)

	filters, ok := parsed.Filters.Value.([]pageFilters)
	if ok {
		if filters[0].Column != "user.average_point" {
			t.Errorf(format, "Filter field for user.average_point", "user.average_point", filters[0].Column)
		}
		if filters[0].Operator != "LIKE" {
			t.Errorf(format, "Filter operator for user.average_point", "LIKE", filters[0].Operator)
		}
		value, isValid := filters[0].Value.(string)
		expected := "%" + avg + "%"
		if !isValid || value != expected {
			t.Errorf(format, "Filter operator for user.average_point", expected, value)
		}
	} else {
		t.Log(parsed.Filters)
		t.Errorf(format, "pageFilters class", "paginate.pageFilters", "null")
	}

}

func TestPaginate(t *testing.T) {
	type User struct {
		gorm.Model
		Name         string `json:"name"`
		AveragePoint string `json:"average_point"`
	}

	type Article struct {
		gorm.Model
		Title   string `json:"title"`
		Content string `json:"content"`
		UserID  uint   `json:"-"`
		User    User   `json:"user"`
	}

	// dsn := "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=postgres sslmode=disable TimeZone=Asia/Jakarta"
	// dsn := "gorm.db"
	dsn := "file::memory:?cache=shared"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})
	db.Exec("PRAGMA case_sensitive_like=ON;")
	db.AutoMigrate(&User{}, &Article{})

	users := []User{{Name: "John doe", AveragePoint: "Seventy %"}, {Name: "Jane doe", AveragePoint: "one hundred %"}}
	articles := []Article{}

	// add massive data
	for i := 0; i < 50; i++ {
		articles = append(articles, Article{
			Title:   fmt.Sprintf("Written by john %d", i),
			Content: fmt.Sprintf("Example by john %d", i),
			UserID:  1,
		})
		articles = append(articles, Article{
			Title:   fmt.Sprintf("Written by jane %d", i),
			Content: fmt.Sprintf("Example by jane %d", i),
			UserID:  2,
		})
	}

	if nil != err {
		t.Error(err.Error())
		return
	}

	tx := db.Begin()

	if err := tx.Create(&users).Error; nil != err {
		tx.Rollback()
		t.Error(err.Error())
		return
	} else if err := tx.Create(&articles).Error; nil != err {
		tx.Rollback()
		t.Error(err.Error())
		return
	} else if err := tx.Commit().Error; nil != err {
		tx.Rollback()
		t.Error(err.Error())
		return
	}

	// wait for transaction to finish
	time.Sleep(1 * time.Second)

	size := 1
	page := 0
	sort := "user.name,-id"
	avg := "y %"
	data := "page=%v&size=%d&sort=%s&filters=%s"

	queryFilter := fmt.Sprintf(`[["user.average_point","like","%s"],["AND"],["user.name","IS NOT",null]]`, avg)
	query := fmt.Sprintf(data, page, size, sort, url.QueryEscape(queryFilter))

	request := &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}
	response := []Article{}

	stmt := db.Joins("User").Model(&Article{})
	result := New(&Config{LikeAsIlikeDisabled: true}).With(stmt).Request(request).Response(&response)

	_, err = json.MarshalIndent(result, "", "  ")
	expectNil(t, err)
	expect(t, result.Page, int64(0), "Invalid page")
	expect(t, result.Total, int64(50), "Invalid total result")
	expect(t, result.TotalPages, int64(50), "Invalid total pages")
	expect(t, result.MaxPage, int64(49), "Invalid max page")
	expectTrue(t, result.First, "Invalid first page")
	expectFalse(t, result.Last, "Invalid last page")

	queryFilter = fmt.Sprintf(`[["users.average_point","like","%s"],["AND"],["user.name","IS NOT",null],["id","like","1"]]`, avg)
	query = fmt.Sprintf(data, page, size, sort, url.QueryEscape(queryFilter))

	request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}
	response = []Article{}

	stmt = db.Joins("User").Model(&Article{})
	result = New(&Config{ErrorEnabled: true}).With(stmt).Request(request).Response(&response)
	expectTrue(t, result.Error, "Failed to get error message")

	page = 1
	size = 100
	pageStart := int64(1)
	query = fmt.Sprintf(data, page, size, sort, "")

	request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}
	response = []Article{}

	stmt = db.Joins("User").Model(&Article{})
	result = New(&Config{PageStart: pageStart}).With(stmt).Request(request).Response(&response)
	expect(t, result.Page, int64(1), "Invalid page start")
	expect(t, result.MaxPage, int64(1), "Invalid max page")
	expect(t, len(response), 100, "Invalid total items")
	expect(t, result.Total, int64(100), "Invalid total result")
	expect(t, result.TotalPages, int64(1), "Invalid total pages")
	expectTrue(t, result.First, "Invalid value first")
	expectTrue(t, result.Last, "Invalid value last")

	queryFilter = `[["user.average_point","like","y %"],["AND"],["user.name,title","LIKE","john"]]`
	query = fmt.Sprintf(data, page, size, sort, url.QueryEscape(queryFilter))

	request = &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: query,
		},
	}
	response = []Article{}

	stmt = db.Joins("User").Model(&Article{})
	result = New(&Config{Operator: "AND", PageStart: pageStart, ErrorEnabled: true}).
		With(stmt).Request(request).Response(&response)
	expectFalse(t, result.Error, "An error occurred")
	expect(t, result.Page, int64(1), "Invalid page start")
	expect(t, result.MaxPage, int64(1), "Invalid max page")
	expect(t, result.Total, int64(50), "Invalid max page")
}

type noOpAdapter struct {
	keyValues          map[string]string
	T                  *testing.T
	clearCounter       int
	clearPrefixCounter int
}

func (n *noOpAdapter) Get(key string) (string, error) {
	n.T.Log(key)
	if v, ok := n.keyValues[key]; ok {
		n.T.Log("OK, Cache found! serving data from cache")
		return v, nil
	}

	n.T.Log("Cache not found")

	return "", errors.New("Cache not found")
}
func (n *noOpAdapter) Set(key string, value string) error {
	if _, ok := n.keyValues[key]; !ok {
		n.keyValues = map[string]string{}
	}
	n.keyValues[key] = value
	n.T.Log("Writing cache")
	return nil
}
func (n *noOpAdapter) IsValid(key string) bool {
	if _, ok := n.keyValues[key]; ok {
		n.T.Log("Cache exists and not expired")
		return false
	}
	n.T.Log("Cache doesn't exists or expired")
	return true
}
func (n *noOpAdapter) Clear(key string) error {
	return nil
}
func (n *noOpAdapter) ClearPrefix(keyPrefix string) error {
	if n.clearPrefixCounter > 2 {
		return errors.New("maximum clear")
	}
	n.clearPrefixCounter = n.clearPrefixCounter + 1
	return nil
}
func (n *noOpAdapter) ClearAll() error {
	if n.clearCounter > 0 {
		return errors.New("maximum clear")
	}
	n.clearCounter = n.clearCounter + 1
	return nil
}

func TestCache(t *testing.T) {
	type User struct {
		gorm.Model
		Name         string `json:"name"`
		AveragePoint string `json:"average_point"`
	}

	type Category struct {
		gorm.Model
		Name string `json:"name"`
	}

	type Article struct {
		gorm.Model
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		UserID     uint     `json:"-"`
		CategoryID uint     `json:"-"`
		User       User     `json:"user"`
		Category   Category `json:"category"`
	}
	dsn := "file::memory:"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})
	if nil != err {
		t.Error(err.Error())
		return
	}
	db.AutoMigrate(&User{}, &Article{})
	categories := []Category{{Name: "Blog"}}
	users := []User{{Name: "John doe", AveragePoint: "Seventy %"}, {Name: "Jane doe", AveragePoint: "one hundred %"}}
	articles := []Article{}
	articles = append(articles, Article{Title: "Written by john", Content: "Example by john", UserID: 1, CategoryID: 1})
	articles = append(articles, Article{Title: "Written by jane", Content: "Example by jane", UserID: 2, CategoryID: 1})
	db.Create(&categories)
	db.Create(&users)
	db.Create(&articles)
	request := &http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: "page=0&size=10&fields=id",
		},
	}

	var adapter gocache.AdapterInterface = &noOpAdapter{T: t}
	config := &Config{
		CacheAdapter:         &adapter,
		FieldSelectorEnabled: true,
	}
	pg := New(config)
	// set cache
	stmt1 := db.Joins("User").Model(&Article{}).Preload(`Category`)
	page1 := pg.With(stmt1).
		Request(request).
		Fields([]string{"id"}).
		Cache("cache_prefix").
		Response(&[]Article{})

	// get cache
	var cached []Article
	stmt2 := db.Joins("User").Model(&Article{})
	page2 := pg.With(stmt2).Request(request).Cache("cache_prefix").Response(&cached)

	if len(cached) < 1 {
		t.Error("Cache pointer not working perfectly")
	}

	if page1.Total != page2.Total {
		t.Error("Total doesn't match")
	}

	pg.ClearCache("cache", "cache_")
	pg.ClearCache("cache", "cache_")
	pg.ClearAllCache()
	pg.ClearAllCache()
}

func expect(t *testing.T, expected interface{}, actual interface{}, message ...string) {
	if expected != actual {
		t.Errorf("%s: Expected %s(%v), got %s(%v)",
			strings.Join(message, " "),
			reflect.TypeOf(expected), expected,
			reflect.TypeOf(actual), actual)
		t.Fail()
	}
}

func expectFalse(t *testing.T, actual bool, message ...string) {
	expect(t, false, actual, message...)
}

func expectTrue(t *testing.T, actual bool, message ...string) {
	expect(t, true, actual, message...)
}

func expectNil(t *testing.T, actual interface{}, message ...string) {
	expect(t, nil, actual, message...)
}

func expectNotNil(t *testing.T, actual interface{}, message ...string) {
	expect(t, false, actual == nil, message...)
}

func TestArrayFilter(t *testing.T) {
	jsonString := `[
		["name,email,address", "like", "abc"]
	]`
	var jsonData []interface{}
	json.Unmarshal([]byte(jsonString), &jsonData)
	filters := arrayToFilter(jsonData, Config{})

	expectNotNil(t, filters)
	expectNotNil(t, filters.Value)

	subFilters, ok := filters.Value.([]pageFilters)
	expectTrue(t, ok)
	expect(t, 1, len(subFilters))

	subFilterValues, ok := subFilters[0].Value.([]pageFilters)
	expectTrue(t, ok)
	expect(t, 1, len(subFilterValues))

	contents, ok := subFilterValues[0].Value.([]pageFilters)
	expectTrue(t, ok)
	expect(t, 5, len(contents))

	expect(t, "name", contents[0].Column)
	expect(t, "LIKE", contents[0].Operator)
	expect(t, "%abc%", contents[0].Value)

	expect(t, "OR", contents[1].Operator)

	expect(t, "email", contents[2].Column)
	expect(t, "LIKE", contents[2].Operator)
	expect(t, "%abc%", contents[2].Value)

	expect(t, "OR", contents[3].Operator)

	expect(t, "address", contents[4].Column)
	expect(t, "LIKE", contents[4].Operator)
	expect(t, "%abc%", contents[4].Value)
}

func TestGenerateWhereCauses(t *testing.T) {
	jsonString := `[
		["name,email,address", "like", "abc"],
		["id", ">", 1]
	]`
	var jsonData []interface{}
	json.Unmarshal([]byte(jsonString), &jsonData)
	filters := arrayToFilter(jsonData, Config{})
	wheres, params := generateWhereCauses(filters, Config{})

	where := strings.Join(wheres, " ")
	where = strings.ReplaceAll(where, "( ", "(")
	where = strings.ReplaceAll(where, " )", ")")
	expect(t, "((((name LIKE ? OR email LIKE ? OR address LIKE ?))) OR (id > ?))", where)
	expect(t, 4, len(params))
}
````

## File: paginate.go
````go
package paginate

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/morkid/gocache"
	"gorm.io/gorm"

	"github.com/valyala/fasthttp"
)

// ResponseContext interface
type ResponseContext interface {
	Cache(string) ResponseContext
	Fields([]string) ResponseContext
	Response(interface{}) Page
}

// RequestContext interface
type RequestContext interface {
	Request(interface{}) ResponseContext
}

// Pagination gorm paginate struct
type Pagination struct {
	Config *Config
}

// With func
func (p *Pagination) With(stmt *gorm.DB) RequestContext {
	return reqContext{
		Statement:  stmt,
		Pagination: p,
	}
}

// ClearCache clear cache contains prefix
func (p Pagination) ClearCache(keyPrefixes ...string) {
	if len(keyPrefixes) > 0 && nil != p.Config && nil != p.Config.CacheAdapter {
		adapter := *p.Config.CacheAdapter
		for i := range keyPrefixes {
			if err := adapter.ClearPrefix(keyPrefixes[i]); nil != err {
				log.Println(err)
			}
		}
	}
}

// ClearAllCache clear all existing cache
func (p Pagination) ClearAllCache() {
	if nil != p.Config && nil != p.Config.CacheAdapter {
		adapter := *p.Config.CacheAdapter
		if err := adapter.ClearAll(); nil != err {
			log.Println(err)
		}
	}
}

type reqContext struct {
	Statement  *gorm.DB
	Pagination *Pagination
}

func (r reqContext) Request(req interface{}) ResponseContext {
	var response ResponseContext = &resContext{
		Statement:  r.Statement,
		Request:    req,
		Pagination: r.Pagination,
	}

	return response
}

type resContext struct {
	Pagination  *Pagination
	Statement   *gorm.DB
	Request     interface{}
	cachePrefix string
	fieldList   []string
}

func (r *resContext) Cache(prefix string) ResponseContext {
	r.cachePrefix = prefix
	return r
}

func (r *resContext) Fields(fields []string) ResponseContext {
	r.fieldList = fields
	return r
}

func (r resContext) Response(res interface{}) Page {
	p := r.Pagination
	query := r.Statement
	p.Config = defaultConfig(p.Config)
	p.Config.Statement = query.Statement
	if p.Config.DefaultSize == 0 {
		p.Config.DefaultSize = 10
	}
	if p.Config.PageStart < 0 {
		p.Config.PageStart = 0
	}

	defaultWrapper := "LOWER(%s)"
	wrappers := map[string]string{
		"sqlite":   defaultWrapper,
		"mysql":    defaultWrapper,
		"postgres": "LOWER((%s)::text)",
	}

	if p.Config.LikeAsIlikeDisabled {
		defaultWrapper := "%s"
		wrappers = map[string]string{
			"sqlite":   defaultWrapper,
			"mysql":    defaultWrapper,
			"postgres": "(%s)::text",
		}
	}

	if p.Config.FieldWrapper == "" && p.Config.ValueWrapper == "" {
		p.Config.FieldWrapper = defaultWrapper
		if wrapper, ok := wrappers[query.Dialector.Name()]; ok {
			p.Config.FieldWrapper = wrapper
		}
	}

	page := Page{}
	pr := parseRequest(r.Request, *p.Config)
	causes := createCauses(pr)
	cKey := ""
	var adapter gocache.AdapterInterface
	var hasAdapter bool = false

	if nil != p.Config.CacheAdapter {
		cKey = createCacheKey(r.cachePrefix, pr)
		adapter = *p.Config.CacheAdapter
		hasAdapter = true
		if cKey != "" && adapter.IsValid(cKey) {
			if cache, err := adapter.Get(cKey); nil == err {
				page.Items = res
				if err := p.Config.JSONUnmarshal([]byte(cache), &page); nil == err {
					return page
				}
			}
		}
	}

	dbs := query.Statement.DB.Session(&gorm.Session{NewDB: true})
	var selects []string
	if len(r.fieldList) > 0 {
		if len(pr.Fields) > 0 && p.Config.FieldSelectorEnabled {
			for i := range pr.Fields {
				for j := range r.fieldList {
					if r.fieldList[j] == pr.Fields[i] {
						fname := query.Statement.Quote("s." + fieldName(pr.Fields[i]))
						if !contains(selects, fname) {
							selects = append(selects, fname)
						}
						break
					}
				}
			}
		} else {
			for i := range r.fieldList {
				fname := query.Statement.Quote("s." + fieldName(r.fieldList[i]))
				if !contains(selects, fname) {
					selects = append(selects, fname)
				}
			}
		}
	} else if len(pr.Fields) > 0 && p.Config.FieldSelectorEnabled {
		for i := range pr.Fields {
			fname := query.Statement.Quote("s." + fieldName(pr.Fields[i]))
			if !contains(selects, fname) {
				selects = append(selects, fname)
			}
		}
	}

	result := dbs.
		Unscoped().
		Table("(?) AS s", query)

	if len(selects) > 0 {
		result = result.Select(selects)
	}

	if len(causes.Params) > 0 || len(causes.WhereString) > 0 {
		result = result.Where(causes.WhereString, causes.Params...)
	}

	result = result.Count(&page.Total).
		Limit(int(causes.Limit)).
		Offset(int(causes.Offset))

	page.RawError = result.Error

	if result.Error != nil && p.Config.ErrorEnabled {
		page.Error = true
		page.ErrorMessage = result.Error.Error()
	}

	if nil != query.Statement.Preloads {
		for table, args := range query.Statement.Preloads {
			result = result.Preload(table, args...)
		}
	}
	if len(causes.Sorts) > 0 {
		for _, sort := range causes.Sorts {
			result = result.Order(sort.Column + " " + sort.Direction)
		}
	}

	rs := result.Find(res)
	if nil == page.RawError {
		page.RawError = rs.Error
	}

	if rs.Error != nil && p.Config.ErrorEnabled && !page.Error {
		page.Error = true
		page.ErrorMessage = rs.Error.Error()
	}

	page.Items = res
	f := float64(page.Total) / float64(causes.Limit)
	if math.Mod(f, 1.0) > 0 {
		f = f + 1
	}
	f = math.Max(f, 1)

	page.TotalPages = int64(f)
	page.MaxPage = page.TotalPages - 1 + p.Config.PageStart
	page.Page = int64(pr.Page)
	page.Size = int64(pr.Size)
	page.Visible = rs.RowsAffected

	if page.Total < 1 {
		page.MaxPage = p.Config.PageStart
		page.TotalPages = 0
	}
	page.First = causes.Offset < 1
	page.Last = page.Page >= page.MaxPage

	if hasAdapter && cKey != "" {
		if cache, err := p.Config.JSONMarshal(page); nil == err {
			if err := adapter.Set(cKey, string(cache)); err != nil {
				log.Println(err)
			}
		}
	}

	return page
}

// New Pagination instance
func New(params ...interface{}) *Pagination {
	if len(params) >= 1 {
		var config *Config
		for _, param := range params {
			c, isConfig := param.(*Config)
			if isConfig {
				config = c
				continue
			}
		}

		return &Pagination{Config: defaultConfig(config)}
	}

	return &Pagination{Config: defaultConfig(nil)}
}

// parseRequest func
func parseRequest(r interface{}, config Config) pageRequest {
	pr := pageRequest{
		Config: *defaultConfig(&config),
	}
	if netHTTP, isNetHTTP := r.(http.Request); isNetHTTP {
		parsingNetHTTPRequest(&netHTTP, &pr)
	} else {
		if netHTTPp, isNetHTTPp := r.(*http.Request); isNetHTTPp {
			parsingNetHTTPRequest(netHTTPp, &pr)
		} else {
			if fastHTTPp, isFastHTTPp := r.(*fasthttp.Request); isFastHTTPp {
				parsingFastHTTPRequest(fastHTTPp, &pr)
			} else {
				if request, isRequest := r.(*Request); isRequest {
					parsingQueryString(request, &pr)
				}
			}
		}
	}

	return pr
}

// createFilters func
func createFilters(filterParams interface{}, p *pageRequest) {
	f, ok := filterParams.([]interface{})
	s, ok2 := filterParams.(string)
	if ok {
		p.Filters = arrayToFilter(f, p.Config)
		p.Filters.Fields = p.Fields
	} else if ok2 {
		iface := []interface{}{}
		if e := p.Config.JSONUnmarshal([]byte(s), &iface); nil == e && len(iface) > 0 {
			p.Filters = arrayToFilter(iface, p.Config)
		}
		p.Filters.Fields = p.Fields
	}
}

// createCauses func
func createCauses(p pageRequest) requestQuery {
	query := requestQuery{}
	wheres, params := generateWhereCauses(p.Filters, p.Config)
	sorts := []sortOrder{}

	for _, so := range p.Sorts {
		so.Column = fieldName(so.Column)
		if nil != p.Config.Statement {
			so.Column = p.Config.Statement.Quote(so.Column)
		}
		sorts = append(sorts, so)
	}

	query.Limit = p.Size
	query.Offset = (p.Page - p.Config.PageStart) * p.Size
	query.Wheres = wheres
	query.WhereString = strings.Join(wheres, " ")
	query.Sorts = sorts
	query.Params = params

	return query
}

// parsingNetHTTPRequest func
func parsingNetHTTPRequest(r *http.Request, p *pageRequest) {
	param := &Request{}
	if r.Method == "" {
		r.Method = "GET"
	}
	if strings.ToUpper(r.Method) == "POST" {
		body, err := io.ReadAll(r.Body)
		if nil != err {
			body = []byte("{}")
		}
		defer r.Body.Close()
		if !p.Config.CustomParamEnabled {
			var postData Request
			if err := p.Config.JSONUnmarshal(body, &postData); nil == err {
				param = &postData
			} else {
				log.Println(err.Error())
			}
		} else {
			var postData map[string]string
			if err := p.Config.JSONUnmarshal(body, &postData); nil == err {
				generateParams(param, p.Config, func(key string) string {
					value, exists := postData[key]
					if !exists {
						value = ""
					}
					return value
				})
			} else {
				log.Println(err.Error())
			}
		}
	} else if strings.ToUpper(r.Method) == "GET" {
		query := r.URL.Query()
		if !p.Config.CustomParamEnabled {
			param.Size, _ = strconv.ParseInt(query.Get("size"), 10, 64)
			param.Page, _ = strconv.ParseInt(query.Get("page"), 10, 64)
			param.Sort = query.Get("sort")
			param.Order = query.Get("order")
			param.Filters = query.Get("filters")
			param.Fields = strings.Split(query.Get("fields"), ",")
		} else {
			generateParams(param, p.Config, func(key string) string {
				return query.Get(key)
			})
		}
	}

	parsingQueryString(param, p)
}

// parsingFastHTTPRequest func
func parsingFastHTTPRequest(r *fasthttp.Request, p *pageRequest) {
	param := &Request{}
	if r.Header.IsPost() {
		b := r.Body()
		if !p.Config.CustomParamEnabled {
			var postData Request
			if err := p.Config.JSONUnmarshal(b, &postData); nil == err {
				param = &postData
			} else {
				log.Println(err.Error())
			}
		} else {
			var postData map[string]string
			if err := p.Config.JSONUnmarshal(b, &postData); nil == err {
				generateParams(param, p.Config, func(key string) string {
					value, exists := postData[key]
					if !exists {
						value = ""
					}
					return value
				})
			} else {
				log.Println(err.Error())
			}
		}
	} else if r.Header.IsGet() {
		query := r.URI().QueryArgs()
		if !p.Config.CustomParamEnabled {
			param.Size, _ = strconv.ParseInt(string(query.Peek("size")), 10, 64)
			param.Page, _ = strconv.ParseInt(string(query.Peek("page")), 10, 64)
			param.Sort = string(query.Peek("sort"))
			param.Order = string(query.Peek("order"))
			param.Filters = string(query.Peek("filters"))
			param.Fields = strings.Split(string(query.Peek("fields")), ",")
		} else {
			generateParams(param, p.Config, func(key string) string {
				return string(query.Peek(key))
			})
		}
	}

	parsingQueryString(param, p)
}

func parsingQueryString(param *Request, p *pageRequest) {
	p.Size = param.Size
	if p.Size == 0 {
		if p.Config.DefaultSize > 0 {
			p.Size = p.Config.DefaultSize
		} else {
			p.Size = 10
		}
	}

	p.Page = param.Page
	if p.Page < p.Config.PageStart {
		p.Page = p.Config.PageStart
	}

	if param.Sort != "" {
		sorts := strings.Split(param.Sort, ",")
		for _, col := range sorts {
			if col == "" {
				continue
			}

			so := sortOrder{
				Column:    col,
				Direction: "ASC",
			}
			if strings.ToUpper(param.Order) == "DESC" {
				so.Direction = "DESC"
			}

			if string(col[0]) == "-" {
				so.Column = string(col[1:])
				so.Direction = "DESC"
			}

			p.Sorts = append(p.Sorts, so)
		}
	}

	if len(param.Fields) > 0 {
		re := regexp.MustCompile(`[^A-z0-9_\.,]+`)
		for _, field := range param.Fields {
			fieldName := re.ReplaceAllString(field, "")
			if fieldName != "" {
				p.Fields = append(p.Fields, fieldName)
			}
		}
	}

	createFilters(param.Filters, p)
}

func generateParams(param *Request, config Config, getValue func(string) string) {
	findValue := func(keys []string, defaultKey string) string {
		found := false
		value := ""
		for _, key := range keys {
			value = getValue(key)
			if value != "" {
				found = true
				break
			}
		}
		if !found {
			return getValue(defaultKey)
		}
		return value
	}

	param.Sort = findValue(config.SortParams, "sort")
	param.Page, _ = strconv.ParseInt(findValue(config.PageParams, "page"), 10, 64)
	param.Size, _ = strconv.ParseInt(findValue(config.SizeParams, "size"), 10, 64)
	param.Order = findValue(config.OrderParams, "order")
	param.Filters = findValue(config.FilterParams, "filters")
	param.Fields = strings.Split(findValue(config.FieldsParams, "fields"), ",")
}

func arrayToFilter(arr []interface{}, config Config) pageFilters {
	filters := pageFilters{
		Single: false,
	}

	operatorEscape := regexp.MustCompile(`[^A-z=\<\>\-\+\^/\*%&! ]+`)
	arrayLen := len(arr)
	defaultOperator := config.Operator
	if defaultOperator == "" {
		defaultOperator = "OR"
	}

	if len(arr) > 0 {
		subFilters := []pageFilters{}
		for k, i := range arr {
			iface, ok := i.([]interface{})
			if ok && !filters.Single {
				subFilters = append(subFilters, arrayToFilter(iface, config))
			} else if arrayLen == 1 {
				operator, ok := i.(string)
				if ok {
					operator = operatorEscape.ReplaceAllString(operator, "")
					filters.Operator = strings.ToUpper(operator)
					filters.IsOperator = true
					filters.Single = true
				}
			} else if arrayLen == 2 {
				if k == 0 {
					if column, ok := i.(string); ok {
						filters.Column = column
						filters.Operator = "="
						filters.Single = true
					}
				} else if k == 1 {
					filters.Value = i
					if nil == i {
						filters.Operator = "IS"
					}
					if strings.Contains(filters.Column, ",") {
						subFilters = filterToSubFilter(&filters, i, config)
						continue
					}
				}
			} else if arrayLen == 3 {
				if k == 0 {
					if column, ok := i.(string); ok {
						filters.Column = column
						filters.Single = true
					}
				} else if k == 1 {
					if operator, ok := i.(string); ok {
						operator = operatorEscape.ReplaceAllString(operator, "")
						filters.Operator = strings.ToUpper(operator)
						filters.Single = true
					}
				} else if k == 2 {
					if strings.Contains(filters.Column, ",") {
						subFilters = filterToSubFilter(&filters, i, config)
						continue
					}
					switch filters.Operator {
					case "LIKE", "ILIKE", "NOT LIKE", "NOT ILIKE":
						escapeString := ""
						escapePattern := `(%|\\)`
						if nil != config.Statement {
							driverName := config.Statement.Dialector.Name()
							switch driverName {
							case "sqlite", "sqlserver", "postgres":
								escapeString = `\`
								filters.ValueSuffix = "ESCAPE '\\'"
							case "mysql":
								escapeString = `\`
								filters.ValueSuffix = `ESCAPE '\\'`
							}
						}
						value := fmt.Sprintf("%v", i)
						re := regexp.MustCompile(escapePattern)
						value = re.ReplaceAllString(value, escapeString+`$1`)
						if config.SmartSearchEnabled {
							re := regexp.MustCompile(`[\s]+`)
							value = re.ReplaceAllString(value, "%")
						}
						filters.Value = fmt.Sprintf("%s%s%s", "%", value, "%")
					default:
						filters.Value = i
					}
				}
			}
		}
		if len(subFilters) > 0 {
			separatedSubFilters := []pageFilters{}
			hasOperator := false
			for k, s := range subFilters {
				if s.IsOperator && len(subFilters) == (k+1) {
					break
				}
				if !hasOperator && !s.IsOperator && k > 0 {
					separatedSubFilters = append(separatedSubFilters, pageFilters{
						Operator:   defaultOperator,
						IsOperator: true,
						Single:     true,
					})
				}
				hasOperator = s.IsOperator
				separatedSubFilters = append(separatedSubFilters, s)
			}
			filters.Value = separatedSubFilters
			filters.Single = false
			filters.IsOperator = false
		}
	}

	return filters
}

func filterToSubFilter(filters *pageFilters, value interface{}, config Config) []pageFilters {
	subFilters := []pageFilters{}
	re := regexp.MustCompile(`[^A-z0-9\._,]+`)
	colString := re.ReplaceAllString(filters.Column, "")
	columns := strings.Split(colString, ",")
	columnRepeat := []interface{}{}
	for _, col := range columns {
		columnRepeat = append(columnRepeat, []interface{}{col, filters.Operator, value})
	}

	filters.Column = ""
	filters.Single = false
	filters.Operator = ""
	filters.IsOperator = false
	subFilters = append(subFilters, arrayToFilter(columnRepeat, config))

	return subFilters
}

//gocyclo:ignore
func generateWhereCauses(f pageFilters, config Config) ([]string, []interface{}) {
	wheres := []string{}
	params := []interface{}{}

	if !f.Single && !f.IsOperator {
		ifaces, ok := f.Value.([]pageFilters)
		if ok && len(ifaces) > 0 {
			wheres = append(wheres, "(")
			hasOpen := false
			for _, i := range ifaces {
				subs, isSub := i.Value.([]pageFilters)
				regular, isNotSub := i.Value.(pageFilters)
				if isSub && len(subs) > 0 {
					wheres = append(wheres, "(")
					for _, s := range subs {
						subWheres, subParams := generateWhereCauses(s, config)
						wheres = append(wheres, subWheres...)
						params = append(params, subParams...)
					}
					wheres = append(wheres, ")")
				} else if isNotSub {
					subWheres, subParams := generateWhereCauses(regular, config)
					wheres = append(wheres, subWheres...)
					params = append(params, subParams...)
				} else {
					if !hasOpen && !i.IsOperator {
						wheres = append(wheres, "(")
						hasOpen = true
					}
					subWheres, subParams := generateWhereCauses(i, config)
					wheres = append(wheres, subWheres...)
					params = append(params, subParams...)
				}
			}
			if hasOpen {
				wheres = append(wheres, ")")
			}
			wheres = append(wheres, ")")
		}
	} else if f.Single {
		if f.IsOperator {
			wheres = append(wheres, f.Operator)
		} else {
			fname := fieldName(f.Column)
			if nil != config.Statement {
				fname = config.Statement.Quote(fname)
			}
			switch f.Operator {
			case "IS", "IS NOT":
				if nil == f.Value {
					wheres = append(wheres, fname, f.Operator, "NULL")
				} else {
					if strValue, isStr := f.Value.(string); isStr && strings.ToLower(strValue) == "null" {
						wheres = append(wheres, fname, f.Operator, "NULL")
					} else {
						wheres = append(wheres, fname, f.Operator, "?")
						params = append(params, f.Value)
					}
				}
			case "BETWEEN":
				if values, ok := f.Value.([]interface{}); ok && len(values) >= 2 {
					wheres = append(wheres, "(", fname, f.Operator, "? AND ?", ")")
					params = append(params, valueFixer(values[0]), valueFixer(values[1]))
				}
			case "IN", "NOT IN":
				if values, ok := f.Value.([]interface{}); ok {
					wheres = append(wheres, fname, f.Operator, "?")
					params = append(params, valueFixer(values))
				}
			case "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
				if config.FieldWrapper != "" {
					fname = fmt.Sprintf(config.FieldWrapper, fname)
				}
				wheres = append(wheres, fname, f.Operator, "?")
				if f.ValueSuffix != "" {
					wheres = append(wheres, f.ValueSuffix)
				}
				value, isStrValue := f.Value.(string)
				if isStrValue {
					if config.ValueWrapper != "" {
						value = fmt.Sprintf(config.ValueWrapper, value)
					} else if !config.LikeAsIlikeDisabled {
						value = strings.ToLower(value)
					}
					params = append(params, value)
				} else {
					params = append(params, f.Value)
				}
			default:
				wheres = append(wheres, fname, f.Operator, "?")
				params = append(params, valueFixer(f.Value))
			}
		}
	}

	return wheres, params
}

func valueFixer(n interface{}) interface{} {
	var values []interface{}
	if rawValues, ok := n.([]interface{}); ok {
		for i := range rawValues {
			values = append(values, valueFixer(rawValues[i]))
		}

		return values
	}
	if nil != n && reflect.TypeOf(n).Name() == "float64" {
		strValue := fmt.Sprintf("%v", n)
		if match, e := regexp.Match(`^[0-9]+$`, []byte(strValue)); nil == e && match {
			v, err := strconv.ParseInt(strValue, 10, 64)
			if nil == err {
				return v
			}
		}
	}

	return n
}

func contains(source []string, value string) bool {
	found := false
	for i := range source {
		if source[i] == value {
			found = true
			break
		}
	}

	return found
}

func fieldName(field string) string {
	slices := strings.Split(field, ".")
	if len(slices) == 1 {
		return field
	}
	newSlices := []string{}
	if len(slices) > 0 {
		newSlices = append(newSlices, strcase.ToCamel(slices[0]))
		for k, s := range slices {
			if k > 0 {
				newSlices = append(newSlices, s)
			}
		}
	}
	if len(newSlices) == 0 {
		return field
	}
	return strings.Join(newSlices, "__")

}

// Config for customize pagination result
type Config struct {
	Operator             string
	FieldWrapper         string
	ValueWrapper         string
	DefaultSize          int64
	PageStart            int64
	LikeAsIlikeDisabled  bool
	SmartSearchEnabled   bool
	Statement            *gorm.Statement `json:"-"`
	CustomParamEnabled   bool
	SortParams           []string
	PageParams           []string
	OrderParams          []string
	SizeParams           []string
	FilterParams         []string
	FieldsParams         []string
	FieldSelectorEnabled bool
	CacheAdapter         *gocache.AdapterInterface              `json:"-"`
	JSONMarshal          func(v interface{}) ([]byte, error)    `json:"-"`
	JSONUnmarshal        func(data []byte, v interface{}) error `json:"-"`
	ErrorEnabled         bool
}

// pageFilters struct
type pageFilters struct {
	Column      string
	Operator    string
	Value       interface{}
	ValuePrefix string
	ValueSuffix string
	Single      bool
	IsOperator  bool
	Fields      []string
}

// Page result wrapper
type Page struct {
	Items        interface{} `json:"items"`
	Page         int64       `json:"page"`
	Size         int64       `json:"size"`
	MaxPage      int64       `json:"max_page"`
	TotalPages   int64       `json:"total_pages"`
	Total        int64       `json:"total"`
	Last         bool        `json:"last"`
	First        bool        `json:"first"`
	Visible      int64       `json:"visible"`
	Error        bool        `json:"error,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
	RawError     error       `json:"-"`
}

// Request struct
type Request struct {
	Page    int64       `json:"page"`
	Size    int64       `json:"size"`
	Sort    string      `json:"sort"`
	Order   string      `json:"order"`
	Fields  []string    `json:"fields"`
	Filters interface{} `json:"filters"`
}

// query struct
type requestQuery struct {
	WhereString string
	Wheres      []string
	Params      []interface{}
	Sorts       []sortOrder
	Limit       int64
	Offset      int64
}

// pageRequest struct
type pageRequest struct {
	Size    int64
	Page    int64
	Sorts   []sortOrder
	Filters pageFilters
	Config  Config `json:"-"`
	Fields  []string
}

// sortOrder struct
type sortOrder struct {
	Column    string
	Direction string
}

func createCacheKey(cachePrefix string, pr pageRequest) string {
	key := ""
	if bte, err := pr.Config.JSONMarshal(pr); nil == err && cachePrefix != "" {
		key = fmt.Sprintf("%s%x", cachePrefix, md5.Sum(bte))
	}

	return key
}

func defaultConfig(c *Config) *Config {
	if nil == c {
		return &Config{
			JSONMarshal:   json.Marshal,
			JSONUnmarshal: json.Unmarshal,
		}
	}

	if nil == c.JSONMarshal {
		c.JSONMarshal = json.Marshal
	}

	if nil == c.JSONUnmarshal {
		c.JSONUnmarshal = json.Unmarshal
	}

	return c
}
````

## File: README.md
````markdown
# paginate - Gorm Pagination

[![Go Reference](https://pkg.go.dev/badge/github.com/morkid/paginate.svg)](https://pkg.go.dev/github.com/morkid/paginate)
[![Github Actions](https://github.com/morkid/paginate/workflows/Go/badge.svg)](https://github.com/morkid/paginate/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/morkid/paginate)](https://goreportcard.com/report/github.com/morkid/paginate)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/morkid/paginate)](https://github.com/morkid/paginate/releases)

Simple way to paginate [Gorm](https://github.com/go-gorm/gorm) result. **paginate** is compatible with [net/http](https://golang.org/pkg/net/http/) and [fasthttp](https://github.com/valyala/fasthttp). This library also supports many net/http or fasthttp based frameworks.

## Table Of Contents
- [Installation](#installation)
- [Configuration](#configuration)
- [Pagination Result](#pagination-result)
- [Paginate using http request](#paginate-using-http-request)
- [Example usage](#example-usage)
  - [net/http](#nethttp-example)
  - [Fasthttp](#fasthttp-example)
  - [Mux Router](#mux-router-example)
  - [Fiber](#fiber-example)
  - [Echo](#echo-example)
  - [Gin](#gin-example)
  - [Martini](#martini-example)
  - [Beego](#beego-example)
  - [jQuery DataTable Integration](#jquery-datatable-integration)
  - [jQuery Select2 Integration](#jquery-select2-integration)
  - [Programmatically Pagination](#programmatically-pagination)
- [Filter format](#filter-format)
- [Customize default configuration](#customize-default-configuration)
- [Override results](#override-results)
- [Field Selector](#field-selector)
- [Dynamic Field Selector](#dynamic-field-selector)
- [Speed up response with cache](#speed-up-response-with-cache)
  - [In Memory Cache](#in-memory-cache)
  - [Disk Cache](#disk-cache)
  - [Redis Cache](#redis-cache)
  - [Elasticsearch Cache](#elasticsearch-cache)
  - [Custom cache](#custom-cache)
  - [Clean up cache](#clean-up-cache)
- [Limitations](#limitations)
- [License](#license)

## Installation

```bash
go get -u github.com/morkid/paginate
```

## Configuration

```go
var db *gorm.DB = ...
var req *http.Request = ...
// or
// var req *fasthttp.Request

stmt := db.Where("id > ?", 1).Model(&Article{})
pg := paginate.New()
page := pg.With(stmt).Request(req).Response(&[]Article{})

log.Println(page.Total)
log.Println(page.Items)
log.Println(page.First)
log.Println(page.Last)

```
you can customize config with `paginate.Config` struct.  
```go
pg := paginate.New(&paginate.Config{
    DefaultSize: 50,
})
```
see more about [customize default configuration](#customize-default-configuration).

## Pagination Result

```js
{
    // the result items
    "items": *[]any, 
    
    // total results
    // including next pages
    "total": number,   

    // Current page
    // (provided by request parameter, eg: ?page=1)
    // note: page is always start from 0
    "page": number,
    
    // Current size
    // (provided by request parameter, eg: ?size=10)
    // note: negative value means unlimited
    "size": number,    

    // Total Pages
    "total_pages": number,

    // Max Page
    // start from 0 until last index
    // example: 
    //   if you have 3 pages (page0, page1, page2)
    //   max_page is 2 not 3
    "max_page": number,

    // Last Page is true if the page 
    // has reached the end of the page
    "last": bool,

    // First Page is true if the page is 0
    "first": bool,

    // Visible
    // total visible items
    "visible": number,

    // Error
    // true if an error has occurred and
    // paginate.Config.ErrorEnabled is true
    "error": bool,

    // Error Message
    // current error if available and
    // paginate.Config.ErrorEnabled is true
    "error_message": string,
}
```
## Paginate using http request
example paging, sorting and filtering:  
1. `http://localhost:3000/?size=10&page=0&sort=-name`  
    produces:
    ```sql
    SELECT * FROM user ORDER BY name DESC LIMIT 10 OFFSET 0
    ```
    `JSON` response:  
    ```js
    {
        // result items
        "items": [
            {
                "id": 1,
                "name": "john",
                "age": 20
            }
        ],
        "page": 0, // current selected page
        "size": 10, // current limit or size per page
        "max_page": 0, // maximum page
        "total_pages": 1, // total pages
        "total": 1, // total matches including next page
        "visible": 1, // total visible on current page
        "last": true, // if response is first page
        "first": true // if response is last page
    }
    ```
2. `http://localhost:3000/?size=10&page=1&sort=-name,id`  
    produces:
    ```sql
    SELECT * FROM user ORDER BY name DESC, id ASC LIMIT 10 OFFSET 10
    ```
3. `http://localhost:3000/?filters=["name","john"]`  
    produces:
    ```sql
    SELECT * FROM user WHERE name = 'john' LIMIT 10 OFFSET 0
    ```
4. `http://localhost:3000/?filters=["name","like","john"]`  
    produces:
    ```sql
    SELECT * FROM user WHERE name LIKE '%john%' LIMIT 10 OFFSET 0
    ```
5. `http://localhost:3000/?filters=["age","between",[20, 25]]`  
    produces:
     ```sql
    SELECT * FROM user WHERE ( age BETWEEN 20 AND 25 ) LIMIT 10 OFFSET 0
    ```
6. `http://localhost:3000/?filters=[["name","like","john%25"],["OR"],["age","between",[20, 25]]]`  
    produces:
     ```sql
    SELECT * FROM user WHERE (
        (name LIKE '%john\%%' ESCAPE '\') OR (age BETWEEN (20 AND 25))
    ) LIMIT 10 OFFSET 0
    ```
7. `http://localhost:3000/?filters=[[["name","like","john"],["AND"],["name","not like","doe"]],["OR"],["age","between",[20, 25]]]`  
    produces:
     ```sql
    SELECT * FROM user WHERE (
        (
            (name LIKE '%john%')
                    AND
            (name NOT LIKE '%doe%')
        ) 
        OR 
        (age BETWEEN (20 AND 25))
    ) LIMIT 10 OFFSET 0
    ```
8. `http://localhost:3000/?filters=["name","IS NOT",null]`  
    produces:
    ```sql
    SELECT * FROM user WHERE name IS NOT NULL LIMIT 10 OFFSET 0
    ```
9. Using `POST` method:  
   ```bash
   curl -X POST \
   -H 'Content-type: application/json' \
   -d '{"page":1,"size":20,"sort":"-name","filters":["name","john"]}' \
   http://localhost:3000/
   ```  
10. You can bypass HTTP Request with [Custom Request](#programmatically-pagination).

## Example usage

### NetHTTP Example

```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        stmt := db.Joins("User").Model(&Article{})
        page := pg.With(stmt).Request(r).Response(&[]Article{})
        j, _ := json.Marshal(page)
        w.Header().Set("Content-type", "application/json")
        w.Write(j)
    })

    log.Fatal(http.ListenAndServe(":3000", nil))
}
```

### Fasthttp Example

```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()

    fasthttp.ListenAndServe(":3000", func(ctx *fasthttp.RequestCtx) {
        stmt := db.Joins("User").Model(&Article{})
        page := pg.With(stmt).Request(&ctx.Request).Response(&[]Article{})
        j, _ := json.Marshal(page)
        ctx.SetContentType("application/json")
        ctx.SetBody(j)
    })
}
```

### Mux Router Example
```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()
    app := mux.NewRouter()
    app.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        stmt := db.Joins("User").Model(&Article{})
        page := pg.With(stmt).Request(req).Response(&[]Article{})
        j, _ := json.Marshal(page)
        w.Header().Set("Content-type", "application/json")
        w.Write(j)
    }).Methods("GET")
    http.Handle("/", app)
    http.ListenAndServe(":3000", nil)
}
```

### Fiber example

```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()
    app := fiber.New()
    app.Get("/", func(c *fiber.Ctx) error {
        stmt := db.Joins("User").Model(&Article{})
        page := pg.With(stmt).Request(c.Request()).Response(&[]Article{})
        return c.JSON(page)
    })

    app.Listen(":3000")
}
```

### Echo example

```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()
    app := echo.New()
    app.GET("/", func(c echo.Context) error {
        stmt := db.Joins("User").Model(&Article{})
        page := pg.With(stmt).Request(c.Request()).Response(&[]Article{})
        return c.JSON(200, page)
    })

    app.Logger.Fatal(app.Start(":3000"))
}
```

### Gin Example

```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()
    app := gin.Default()
    app.GET("/", func(c *gin.Context) {
        stmt := db.Joins("User").Model(&Article{})
        page := pg.With(stmt).Request(c.Request).Response(&[]Article{})
        c.JSON(200, page)
    })
    app.Run(":3000")
}

```

### Martini Example

```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()
    app := martini.Classic()
    app.Use(render.Renderer())
    app.Get("/", func(req *http.Request, r render.Render) {
        stmt := db.Joins("User").Model(&Article{})
        page := pg.With(stmt).Request(req).Response(&[]Article{})
        r.JSON(200, page)
    })
    app.Run()
}
```
### Beego Example

```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()
    web.Get("/", func(c *context.Context) {
        stmt := db.Joins("User").Model(&Article{})
        page := pg.With(stmt).Request(c.Request).Response(&[]Article{})
        c.Output.JSON(page, false, false)
    })
    web.Run(":3000")
}
```

### jQuery DataTable Integration

```js
var logicalOperator = "OR"

$('#myTable').DataTable({

    columns: [
        {
            title: "Author",
            data: "user.name"
        }, {
            title: "Title",
            data: "title"
        }
    ],

    processing: true,
    
    serverSide: true,

    ajax: {
        cache: true,
        url: "http://localhost:3000/articles",
        dataSrc: function(json) {
            json.recordsTotal = json.visible
            json.recordsFiltered = json.total
            return json.items
        },
        data: function(params) {
            var custom = {
                page: !params.start ? 0 : Math.round(params.start / params.length),
                size: params.length
            }

            if (params.order.length > 0) {
                var sorts = []
                for (var o in params.order) {
                    var order = params.order[o]
                    if (params.columns[order.column].orderable != false) {
                        var sort = order.dir != 'desc' ? '' : '-'
                        sort += params.columns[order.column].data
                        sorts.push(sort)
                    }
                }
                custom.sort = sorts.join()
            }

            if (params.search.value) {
                var columns = []
                for (var c in params.columns) {
                    var col = params.columns[c]
                    if (col.searchable == false) {
                        continue
                    }
                    columns.push(JSON.stringify([col.data, "like", encodeURIComponent(params.search.value.toLowerCase())]))
                }
                custom.filters = '[' + columns.join(',["' + logicalOperator + '"],') + ']'
            }

            return custom
        }
    },
})
```

### jQuery Select2 Integration

```js
$('#mySelect').select2({
    ajax: {
        url: "http://localhost:3000/users",
        processResults: function(json) {
            json.items.forEach(function(item) {
                item.text = item.name
            })
            // optional
            if (json.first) json.items.unshift({id: 0, text: 'All'})

            return {
                results: json.items,
                pagination: {
                    more: json.last == false
                }
            }
        },
        data: function(params) {
            var filters = [
                ["name", "like", params.term]
            ]

            return {
                filters: params.term ? JSON.stringify(filters) : "",
                sort: "name",
                page: params.page && params.page - 1 ? params.page - 1 : 0
            }
        },
    }
})
```

### Programmatically Pagination

```go
package main

import (
    "github.com/morkid/paginate"
    ...
)

func main() {
    // var db *gorm.DB
    pg := paginate.New()
    req := &paginate.Request{
        Page: 2,
        Size: 20,
        Sort: "-publish_date",
        filters: []interface{}{
            []interface{}{"user.name", "like", "john"},
            []interface{}{"and"},
            []interface{}{"publish_date", ">=", "2025-12-31"},
            []interface{}{"and"},
            []interface{}{"user.active", "=", true},
            []interface{}{"and"},
            []interface{}{"user.last_login", "is not", nil},
        }
    }

    stmt := db.Joins("User").Model(&Article{})
    page := pg.With(stmt).
        Request(req).
        Response(&[]Article{})

}
```


## Filter format

Paginate Filters was inspired by [Frappe Framework API](https://docs.frappe.io/framework/user/en/api/rest#listing-documents). This feature is very powerful to support deep search and keep it safe. The format of filter param is a json encoded of multidimensional array.  
Maximum array members is three, first index is `column_name`, second index is `operator` and third index is `values`, you can also pass array to values.  

```js
// Format:
["column_name", "operator", "values"]

// Example:
["age", "=", 20]
// Shortcut:
["age", 20]

// Produces:
// WHERE age = 20
```

Single array member is known as **Logical Operator**.
```js
// Example
[["age", "=", 20],["or"],["age", "=", 25]]

// Produces:
// WHERE age = 20 OR age = 25
```


You are allowed to send array inside a value.  
```js
["age", "between", [20, 30] ]
// Produces:
// WHERE age BETWEEN 20 AND 30

["age", "not in", [20, 21, 22, 23, 24, 25, 26, 26] ]
// Produces:
// WHERE age NOT IN(20, 21, 22, 23, 24, 25, 26, 26)
```

Define chain columns with same value separated by comma.
```js
// Example 1
["price,discount", ">", 10]
// Produces:
// WHERE price > 10 OR discount > 10

// Example 2
["deleted_at,expiration_date", null]
// Produces:
// WHERE deleted_at IS NULL OR expiration_date IS NULL
```

You can filter nested condition with deep array.  
```js
[
    [
        ["age", ">", 20],
        ["and"]
        ["age", "<", 30]
    ],
    ["and"],
    ["name", "like", "john"],
    ["and"],
    ["name", "like", "doe"]
]
// Produces:
// WHERE ( (age > 20 AND age < 20) and name like '%john%' and name like '%doe%' )
```

For `null` value, you can send string `"null"` or `null` value, *(lower)*
```js
// Wrong request
[ "age", "is", NULL ]
[ "age", "is", Null ]
[ "age", "is not", NULL ]
[ "age", "is not", Null ]

// Right request
[ "age", "is", "NULL" ]
[ "age", "is", "Null" ]
[ "age", "is", "null" ]
[ "age", "is", null ]
[ "age", null ]
[ "age", "is not", "NULL" ]
[ "age", "is not", "Null" ]
[ "age", "is not", "null" ]
[ "age", "is not", null ]
```

## Customize default configuration

You can customize the default configuration with `paginate.Config` struct. 

```go
pg := paginate.New(&paginate.Config{
    DefaultSize: 50,
})
```

Config             | Type       | Default               | Description
------------------ | ---------- | --------------------- | -------------
Operator           | `string`   | `OR`                  | Default conditional operator if no operator specified.<br>For example<br>`GET /user?filters=[["name","like","jo"],["age",">",20]]`,<br>produces<br>`SELECT * FROM user where name like '%jo' OR age > 20`
FieldWrapper       | `string`   | `LOWER(%s)`           | FieldWrapper for `LIKE` operator *(for postgres default is: `LOWER((%s)::text)`)*
DefaultSize        | `int64`    | `10`                  | Default size or limit per page
PageStart          | `int64`    | `0`                   | Set start page, default `0` if not set. `total_pages` , `max_page` and `page` variable will be affected if you set `PageStart` greater than `0` 
LikeAsIlikeDisabled | `bool`    | `false`               | By default, paginate using Case Insensitive on `LIKE` operator. Instead of using `ILIKE`, you can use `LIKE` operator to find what you want. You can set `LikeAsIlikeDisabled` to `true` if you need this feature to be disabled.
SmartSearchEnabled | `bool`     | `false`               | Enable smart search *(Experimental feature)*
CustomParamEnabled | `bool`     | `false`               | Enable custom request parameter
FieldSelectorEnabled | `bool`   | `false`               | Enable partial response with specific fields. Comma separated per field. eg: `?fields=title,user.name`
SortParams         | `[]string` | `[]string{"sort"}`    | if `CustomParamEnabled` is `true`,<br>you can set the `SortParams` with custom parameter names.<br>For example: `[]string{"sorting", "ordering", "other_alternative_param"}`.<br>The following requests will capture same result<br>`?sorting=-name`<br>or `?ordering=-name`<br>or `?other_alternative_param=-name`<br>or `?sort=-name`
PageParams         | `[]string` | `[]string{"page"}`    | if `CustomParamEnabled` is `true`,<br>you can set the `PageParams` with custom parameter names.<br>For example:<br>`[]string{"number", "num", "other_alternative_param"}`.<br>The following requests will capture same result `?number=0`<br>or `?num=0`<br>or `?other_alternative_param=0`<br>or `?page=0`
SizeParams         | `[]string` | `[]string{"size"}`    | if `CustomParamEnabled` is `true`,<br>you can set the `SizeParams` with custom parameter names.<br>For example:<br>`[]string{"limit", "max", "other_alternative_param"}`.<br>The following requests will capture same result `?limit=50`<br>or `?limit=50`<br>or `?other_alternative_param=50`<br>or `?max=50`
OrderParams         | `[]string` | `[]string{"order"}`    | if `CustomParamEnabled` is `true`,<br>you can set the `OrderParams` with custom parameter names.<br>For example:<br>`[]string{"order", "direction", "other_alternative_param"}`.<br>The following requests will capture same result `?order=desc`<br>or `?direction=desc`<br>or `?other_alternative_param=desc`
FilterParams       | `[]string` | `[]string{"filters"}` | if `CustomParamEnabled` is `true`,<br>you can set the `FilterParams` with custom parameter names.<br>For example:<br>`[]string{"search", "find", "other_alternative_param"}`.<br>The following requests will capture same result<br>`?search=["name","john"]`<br>or `?find=["name","john"]`<br>or `?other_alternative_param=["name","john"]`<br>or `?filters=["name","john"]`
FieldsParams       | `[]string` | `[]string{"fields"}`  | if `FieldSelectorEnabled` and `CustomParamEnabled` is `true`,<br>you can set the `FieldsParams` with custom parameter names.<br>For example:<br>`[]string{"fields", "columns", "other_alternative_param"}`.<br>The following requests will capture same result `?fields=title,user.name`<br>or `?columns=title,user.name`<br>or `?other_alternative_param=title,user.name`
CacheAdapter       | `*gocache.AdapterInterface` | `nil` | the cache adapter, see more about [cache config](#speed-up-response-with-cache).
ErrorEnabled       | `bool` | `false` | Show error message in pagination result.

## Override results

You can override result with custom function.  

```go
// var db = *gorm.DB
// var httpRequest ... net/http or fasthttp instance
// Example override function
override := func(article *Article) {
    if article.UserID > 0 {
        article.Title = fmt.Sprintf(
            "%s written by %s", article.Title, article.User.Name)
    }
}

var articles []Article
stmt := db.Joins("User").Model(&Article{})

pg := paginate.New()
page := pg.With(stmt).Request(httpRequest).Response(&articles)
for index := range articles {
    override(&articles[index])
}

log.Println(page.Items)

```

## Field selector
To implement a custom field selector, struct properties must have a json tag with omitempty.

```go
// real gorm model
type User {
    gorm.Model
    Name string `json:"name"`
    Age  int64  `json:"age"`
}

// fake gorm model
type UserNullable {
    ID        *string    `json:"id,omitempty"`
    CreatedAt *time.Time `json:"created_at,omitempty"`
    UpdatedAt *time.Time `json:"updated_at,omitempty"`
    Name      *string    `json:"name,omitempty"`
    Age       *int64     `json:"age,omitempty"`
}
```

```go
// usage
nameAndIDOnly := []string{"name","id"}
stmt := db.Model(&User{})

page := pg.With(stmt).
   Request(req).
   Fields(nameAndIDOnly).
   Response([]&UserNullable{})
```

```javascript
// response
{
    "items": [
        {
            "id": 1,
            "name": "John"
        }
    ],
    ...
}
```
## Dynamic field selector
If the request contains query parameter `fields` (eg: `?fieilds=name,id`), then the response will show only `name` and `id`. To activate this feature, please set `FieldSelectorEnabled` to `true`.
```go
config := paginate.Config{
    FieldSelectorEnabled: true,
}

pg := paginate.New(config)
```

## Speed up response with cache
You can speed up results without looking database directly with cache adapter. See more about [cache adapter](https://github.com/morkid/gocache).

### In memory cache
in memory cache is not recommended for production environment:
```go
import (
    "github.com/morkid/gocache"
    ...
)

func main() {
    ...
    adapterConfig := gocache.InMemoryCacheConfig{
        ExpiresIn: 1 * time.Hour,
    }
    pg := paginate.New(&paginate.Config{
        CacheAdapter: gocache.NewInMemoryCache(adapterConfig),
    })

    page := pg.With(stmt).
               Request(req).
               Cache("article"). // set cache name
               Response(&[]Article{})
    ...
}
```

### Disk cache
Disk cache will create a file for every single request. You can use disk cache if you don't care about inode.
```go
import (
    "github.com/morkid/gocache"
    ...
)

func main() {
    adapterConfig := gocache.DiskCacheConfig{
        Directory: "/writable/path/to/my-cache-dir",
        ExpiresIn: 1 * time.Hour,
    }
    pg := paginate.New(&paginate.Config{
        CacheAdapter: gocache.NewDiskCache(adapterConfig),
    })

    page := pg.With(stmt).
               Request(req).
               Cache("article"). // set cache name
               Response(&[]Article{})
    ...
}
```

### Redis cache
Redis cache require [redis client](https://github.com/go-redis/redis) for golang.
```go
import (
    cache "github.com/morkid/gocache-redis/v8"
    "github.com/go-redis/redis/v8"
    ...
)

func main() {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    adapterConfig := cache.RedisCacheConfig{
        Client:    client,
        ExpiresIn: 1 * time.Hour,
    }
    pg := paginate.New(&paginate.Config{
        CacheAdapter: cache.NewRedisCache(adapterConfig),
    })

    page := pg.With(stmt).
               Request(req).
               Cache("article").
               Response(&[]Article{})
    ...
}
```
> if your code already adopts another redis client, you can implement the [redis adapter](https://github.com/morkid/gocache-redis) according to its version. See more about [redis adapter](https://github.com/morkid/gocache-redis).

### Elasticsearch cache
Elasticsearch cache require official [elasticsearch client](https://github.com/elastic/go-elasticsearch) for golang.
```go
import (
    cache "github.com/morkid/gocache-elasticsearch/v7"
    "github.com/elastic/go-elasticsearch/v7"
    ...
)

func main() {
    config := elasticsearch.Config{
        Addresses: []string{
            "http://localhost:9200",
        },
    }
    es, err := elasticsearch.NewClient(config)
    if nil != err {
        panic(err)
    }

    adapterConfig := cache.ElasticCacheConfig{
        Client:    es,
        Index:     "exampleproject",
        ExpiresIn: 1 * time.Hour,
    }
    pg := paginate.New(&paginate.Config{
        CacheAdapter: cache.NewElasticCache(adapterConfig),
    })

    page := pg.With(stmt).
               Request(req).
               Cache("article").
               Response(&[]Article{})
    ...
}
```
> if your code already adopts another elasticsearch client, you can implement the [elasticsearch adapter](https://github.com/morkid/gocache-elasticsearch) according to its version. See more about [elasticsearch adapter](https://github.com/morkid/gocache-elasticsearch).

### Custom cache
Create your own cache adapter by implementing [gocache AdapterInterface](https://github.com/morkid/gocache/blob/master/gocache.go). See more about [cache adapter](https://github.com/morkid/gocache).
```go
// AdapterInterface interface
type AdapterInterface interface {
    // Set cache with key
    Set(key string, value string) error
    // Get cache by key
    Get(key string) (string, error)
    // IsValid check if cache is valid
    IsValid(key string) bool
    // Clear clear cache by key
    Clear(key string) error
    // ClearPrefix clear cache by key prefix
    ClearPrefix(keyPrefix string) error
    // Clear all cache
    ClearAll() error
}
```

### Clean up cache
Clear cache by cache name
```go
pg.ClearCache("article")
```
Clear multiple cache
```go
pg.ClearCache("cache1", "cache2", "cache3")
```

Clear all cache
```go
pg.ClearAllCache()
```


## Limitations

Paginate doesn't support has many relationship. You can make API with separated endpoints for parent and child:
```javascript
GET /users

{
    "items": [
        {
            "id": 1,
            "name": "john",
            "age": 20,
            "addresses": [...] // doesn't support
        }
    ],
    ...
}
```

Best practice:

```javascript
GET /users
{
    "items": [
        {
            "id": 1,
            "name": "john",
            "age": 20
        }
    ],
    ...
}

GET /users/1/addresses
{
    "items": [
        {
            "id": 1,
            "name": "home",
            "street": "home street"
            "user": {
                "id": 1,
                "name": "john",
                "age": 20
            }
        }
    ],
    ...
}
```

Paginate doesn't support for customized json or table field name.  
Make sure your struct properties have same name with gorm column and json property before you expose them.  

Example bad configuration:  

```go

type User struct {
    gorm.Model
    UserName       string `gorm:"column:nickname" json:"name"`
    UserAddress    string `gorm:"column:user_address" json:"address"`
}

// request: GET /path/to/endpoint?sort=-name,address
// response: "items": [] with sql error (column name not found)
```

Best practice:
```go
type User struct {
    gorm.Model
    Name       string `gorm:"column:name" json:"name"`
    Address    string `gorm:"column:address" json:"address"`
}

```

## License

Published under the [MIT License](https://github.com/morkid/paginate/blob/master/LICENSE).
````
