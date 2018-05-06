package main

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Get_Index(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	index(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
	if strings.Contains(string(body), "Stranger") == false {
		t.Errorf("Body does not contain Stranger")
	}
}
