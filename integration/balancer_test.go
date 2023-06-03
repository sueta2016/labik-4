package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type IntegrationSuite struct{}

var _ = Suite(&IntegrationSuite{})

const (
	baseAddress = "http://balancer:8090"
	teamName    = "sueta2016"
)

var client = http.Client{
	Timeout: 3 * time.Second,
}

type RespBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *IntegrationSuite) TestLoadBalancer(c *C) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		c.Skip("Integration test is not enabled")
	}

	resp1, _ := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	c.Assert(resp1.StatusCode, Equals, http.StatusBadRequest)

	resp2, _ := client.Get(fmt.Sprintf("%s/api/v1/some-data?key=sueta", baseAddress))
	c.Assert(resp2.StatusCode, Equals, http.StatusNotFound)

	db, err := client.Get(fmt.Sprintf("%s/api/v1/some-data?key=sueta2016", baseAddress))
	c.Assert(err, IsNil)

	var body RespBody
	err = json.NewDecoder(db.Body).Decode(&body)
	c.Assert(err, IsNil)

	c.Assert(body.Key, Equals, teamName)
	c.Assert(body.Value, Not(Equals), "")

	db.Body.Close()
}

func (s *IntegrationSuite) BenchmarkLoadBalancer(c *C) {
	for i := 0; i < c.N; i++ {
		resp, _ := client.Get(fmt.Sprintf("%s/api/v1/some-data?key=sueta2016", baseAddress))
		c.Assert(resp.StatusCode, Equals, http.StatusOK)
	}
}