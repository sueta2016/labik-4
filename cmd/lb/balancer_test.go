package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestBalancer(c *C) {
	first_host := getIndex("127.0.0.1:8080")
	second_host := getIndex("192.168.0.0:80")
	third_host := getIndex("26.143.218.9:80")

	c.Assert(first_host, Equals, 2)
	c.Assert(second_host, Equals, 0)
	c.Assert(third_host, Equals, 1)
}

func (s *TestSuite) TestHealth(c *C) {
	result := make([]string, len(serversPool))

	first_host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer first_host.Close()

	second_host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer second_host.Close()

	parsedURL1, _ := url.Parse(first_host.URL)
	hostURL1 := parsedURL1.Host

	parsedURL2, _ := url.Parse(first_host.URL)
	hostURL2 := parsedURL2.Host

	servers := []string{
		hostURL1,
		hostURL2,
		"third_host:8080",
	}

	healthCheck(servers, result)
	time.Sleep(12 * time.Second)

	c.Assert(result[0], Equals, hostURL1)
	c.Assert(result[1], Equals, hostURL2)
	c.Assert(result[2], Equals, "")
}