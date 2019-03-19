package main

import (
	"encoding/json"
	"github.com/linghaihui/kuafu"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type Hello struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func SayHello(ctx *kuafu.Context) {
	age := ctx.GetParam("age", "10")
	ageInt, _ := strconv.Atoi(age)
	ctx.Json(200, Hello{
		Code: 0,
		Msg:  kuafu.FormatString("Hello %s", ctx.GetParam("name")),
		Name: ctx.GetParam("name"),
		Age:  ageInt,
	})
}

func TestKuafu(t *testing.T) {
	server := kuafu.NewServer()
	server.SetDebugMode()
	server.Use(kuafu.PrintRequest, kuafu.Handler404)
	registry := server.NewRegistry()
	registry.GET("/hello/<:name>/<:age>", SayHello)
	registry.POST("/hello/<:name>/<:age>", SayHello)
	group := server.NewGroup("group", "/group")
	group.GET("/hello", SayHello)
	// server.Run("127.0.0.1:9999")
	ts := httptest.NewServer(server)
	defer ts.Close()
	req, _ := http.NewRequest("GET", kuafu.FormatString("%s/hello/haihui/25", ts.URL), nil)
	resp, _ := ts.Client().Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header["X-Server-Framework"][0], kuafu.FormatString("Kuafu/%s", kuafu.Version))
	body, _ := ioutil.ReadAll(resp.Body)
	h := Hello{}
	json.Unmarshal(body, &h)
	assert.Equal(t, h.Name, "haihui")
	assert.Equal(t, h.Age, 25)
	// test POST
	req, _ = http.NewRequest("POST", kuafu.FormatString("%s/hello/linghaihui/2", ts.URL), nil)
	resp, _ = ts.Client().Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header["X-Server-Framework"][0], kuafu.FormatString("Kuafu/%s", kuafu.Version))
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &h)
	assert.Equal(t, h.Name, "linghaihui")
	assert.Equal(t, h.Age, 2)
	// test group
	req, _ = http.NewRequest("GET", kuafu.FormatString("%s/group/hello", ts.URL), nil)
	resp, _ = ts.Client().Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header["X-Server-Framework"][0], kuafu.FormatString("Kuafu/%s", kuafu.Version))
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &h)
	assert.Equal(t, h.Name, "")
	assert.Equal(t, h.Age, 10)
	req, _ = http.NewRequest("GET", kuafu.FormatString("%s/group/hello/world", ts.URL), nil)
	resp, _ = ts.Client().Do(req)
	assert.Equal(t, resp.StatusCode, 404)
}
