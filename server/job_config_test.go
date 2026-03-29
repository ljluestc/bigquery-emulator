package server_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/goccy/bigquery-emulator/server"
	"github.com/goccy/bigquery-emulator/types"
)

func TestJobTimeoutMsUnmarshal(t *testing.T) {
	ctx := context.Background()
	bqServer, err := server.New(server.TempStorage)
	if err != nil {
		t.Fatal(err)
	}
	if err := bqServer.Load(server.StructSource(types.NewProject("test"))); err != nil {
		t.Fatal(err)
	}
	testServer := bqServer.TestServer()
	defer func() {
		testServer.Close()
		bqServer.Stop(ctx)
	}()

	// Simulate the JSON body sent by Node.js client with jobTimeoutMs as number
	body := `{
		"configuration": {
			"query": {
				"query": "SELECT 1",
				"useLegacySql": false
			},
			"jobTimeoutMs": 60000
		}
	}`

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/bigquery/v2/projects/test/jobs", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}
}
