package githubkey

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ClientMock struct {
	do func(req *http.Request) (*http.Response, error)
}

func (mock *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return mock.do(req)
}

func Test_GetDeployKey(t *testing.T) {
	tests := map[string]struct {
		client           Doer
		githubUsername   string
		githubPassword   string
		keyTitle         string
		repo             string
		expectedResponse GithubKey
		expectedError    error
	}{
		"error": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{}, fmt.Errorf("an error")
				},
			},
			githubUsername:   "test",
			githubPassword:   "test",
			keyTitle:         "testKeyTitle",
			repo:             "test",
			expectedResponse: GithubKey{},
			expectedError:    fmt.Errorf("an error"),
		},
		"success": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					testKeys := &githubKeys{
						GithubKey{
							ID:    42,
							Title: "testKeyTitle",
						},
					}
					testKeysBytes, _ := testKeys.marshal()
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader(testKeysBytes)),
						StatusCode: 200,
					}, nil
				},
			},
			githubUsername: "test",
			githubPassword: "test",
			keyTitle:       "testKeyTitle",
			repo:           "test",
			expectedResponse: GithubKey{
				ID:    42,
				Title: "testKeyTitle",
			},
			expectedError: nil,
		},
		"success_not_found": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					testKeys := &githubKeys{
						GithubKey{
							ID:    42,
							Title: "testKeyTitle",
						},
					}
					testKeysBytes, _ := testKeys.marshal()
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader(testKeysBytes)),
						StatusCode: 200,
					}, nil
				},
			},
			githubUsername:   "test",
			githubPassword:   "test",
			keyTitle:         "testKeyTitle1",
			repo:             "test",
			expectedResponse: GithubKey{},
			expectedError:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("Running test case: %s", name)
			response, err := GetDeployKey(test.client, test.githubUsername, test.githubPassword, test.repo, test.keyTitle)
			assert.Equal(t, test.expectedResponse, response)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func Test_DeleteDeployKey(t *testing.T) {
	tests := map[string]struct {
		client         Doer
		githubUsername string
		githubPassword string
		repo           string
		keyID          int64
		expectedError  error
	}{
		"error": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{}, fmt.Errorf("an error")
				},
			},
			githubUsername: "test",
			githubPassword: "test",
			repo:           "test",
			keyID:          1,
			expectedError:  fmt.Errorf("an error"),
		},
		"success": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					testKeys := &githubKeys{
						GithubKey{
							ID:    42,
							Title: "testKeyTitle",
						},
					}
					testKeysBytes, _ := testKeys.marshal()
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader(testKeysBytes)),
						StatusCode: 204,
					}, nil
				},
			},
			githubUsername: "test",
			githubPassword: "test",
			repo:           "test",
			keyID:          1,
			expectedError:  nil,
		},
		"error_403": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
						StatusCode: 403,
					}, nil
				},
			},
			githubUsername: "test",
			githubPassword: "test",
			repo:           "test",
			keyID:          1,
			expectedError:  fmt.Errorf("could not delete KeyID 1"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("Running test case: %s", name)
			err := DeleteDeployKey(test.client, test.githubUsername, test.githubPassword, test.repo, test.keyID)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func Test_CreateDeployKey(t *testing.T) {
	tests := map[string]struct {
		client           Doer
		githubUsername   string
		githubPassword   string
		keyTitle         string
		repo             string
		newKey           string
		readOnly         bool
		expectedResponse GithubKey
		expectedError    error
	}{
		"error": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{}, fmt.Errorf("an error")
				},
			},
			githubUsername:   "test",
			githubPassword:   "test",
			keyTitle:         "testKeyTitle",
			repo:             "test",
			newKey:           "test",
			readOnly:         true,
			expectedResponse: GithubKey{},
			expectedError:    fmt.Errorf("an error"),
		},
		"success": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					testKeys := GithubKey{
						ID:    42,
						Title: "testKeyTitle",
					}
					testKeysBytes, _ := testKeys.marshal()
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader(testKeysBytes)),
						StatusCode: 201,
					}, nil
				},
			},
			githubUsername: "test",
			githubPassword: "test",
			keyTitle:       "testKeyTitle",
			repo:           "test",
			expectedResponse: GithubKey{
				ID:    42,
				Title: "testKeyTitle",
			},
			expectedError: nil,
		},
		"error_response_code": {
			client: &ClientMock{
				do: func(req *http.Request) (*http.Response, error) {
					testKeys := GithubKey{}
					testKeysBytes, _ := testKeys.marshal()
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader(testKeysBytes)),
						StatusCode: 403,
					}, nil
				},
			},
			githubUsername:   "test",
			githubPassword:   "test",
			keyTitle:         "testKeyTitle",
			repo:             "test",
			expectedResponse: GithubKey{},
			expectedError:    fmt.Errorf("http status code: 403"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("Running test case: %s", name)
			response, err := CreateDeployKey(test.client, test.githubUsername, test.githubPassword, test.repo, test.keyTitle, test.newKey, test.readOnly)
			assert.Equal(t, test.expectedResponse, response)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
