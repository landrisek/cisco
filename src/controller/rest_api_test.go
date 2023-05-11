package controller

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/landrisek/cisco/src/repository"
)

func TestRestAPI(t *testing.T) {
	// dereferencing for server structure
	node := *repository.NewNode().SetName("root").SetChildren([]repository.GNode{
		*repository.NewNode().SetName("child1").SetChildren([]repository.GNode{
			*repository.NewNode().SetName("grandchild1"),
			*repository.NewNode().SetName("grandchild2"),
		}),
		*repository.NewNode().SetName("child2").SetChildren([]repository.GNode{
			*repository.NewNode().SetName("grandchild3"),
		}),
	})

	go RestAPI(node)

	time.Sleep(100 * time.Millisecond) // wait for server to start

	token := repository.GetValidToken()

	tests := []struct {
		name           string
		url            string
		token          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing tag parameter",
			url:            "http://localhost:8080/taggedContent?token=" + token,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing 'tag' parameter\n",
		},
		{
			name:           "unauthorized request",
			url:            "http://localhost:8080/taggedContent?tag=child1&token=invalid",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized\n",
		},
		{
			name:           "not found",
			url:            "http://localhost:8080/taggedContent?tag=unknown&token=" + token,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Tag unknown was not found\n",
		},
		{
			name:           "success",
			url:            "http://localhost:8080/taggedContent?tag=child1&token=" + token,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"name":"child1","children":[{"name":"grandchild1","children":null},{"name":"grandchild2","children":null}]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.url+"", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("token", tc.token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, but got %d", tc.expectedStatus, resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if string(body) != tc.expectedBody {
				t.Errorf("Expected response body %q, but got %q", tc.expectedBody, string(body))
			}
		})
	}

	// Wait for server to shut down
	time.Sleep(100 * time.Millisecond)
}
