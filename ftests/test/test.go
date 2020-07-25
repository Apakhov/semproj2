package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

var Host = "http://127.0.0.1:8181"

type TestCase struct {
	Path      string
	Method    string
	GetParams map[string]string
	Body      interface{}
	Response  WithError
}

func ToMap(s interface{}) map[string]interface{} {
	bt, _ := json.Marshal(&s)
	mp := map[string]interface{}{}
	json.Unmarshal(bt, &mp)
	return mp
}

type WithError struct {
	Err string      `json:"error"`
	Val interface{} `json:"value"`
}

func (tc *TestCase) Run(t *testing.T, dst interface{}) WithError {
	path := tc.Path
	cnt := 0
	for k, v := range tc.GetParams {
		cnt++
		if cnt == 1 {
			path += "?"
		}
		path += k + "=" + v
		if cnt != len(tc.GetParams) {
			path += "&"
		}
	}
	body := []byte{}
	var err error
	if tc.Body != nil {
		body, err = json.Marshal(tc.Body)
		assert.NoError(t, err, "http cr")
	}
	req, err := http.NewRequest(tc.Method, Host+path, bytes.NewBuffer(body))
	assert.NoError(t, err, "http cr")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "http req")
	defer resp.Body.Close()

	respjs, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "http read body")

	res := WithError{}
	fmt.Println(string(respjs))
	err = json.Unmarshal(respjs, &res)
	assert.NoError(t, err, "unp json")
	mp1 := ToMap(tc.Response)
	mp2 := ToMap(res)

	_, ok := (mp1["value"]).(map[string]interface{})
	if ok {
		(mp1["value"]).(map[string]interface{})["id"] = (mp2["value"]).(map[string]interface{})["id"]
	}
	fmt.Println(mp1)
	fmt.Println(mp2)

	diffs := deep.Equal(mp1, mp2)
	assert.Equal(t, []string(nil), diffs, "bad response")
	bt, _ := json.Marshal(res.Val)
	if dst != nil {
		json.Unmarshal(bt, dst)
	}
	return res
}

func UStr() string {
	id, _ := uuid.DefaultGenerator.NewV4()
	return id.String()
}

func UEmail() string {
	return UStr() + "@mail.ru"
}
