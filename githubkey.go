package githubkey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type githubKeys []GithubKey

func unmarshalGithubKeys(data []byte) (githubKeys, error) {
	var r githubKeys
	err := json.Unmarshal(data, &r)
	return r, err
}

func unmarshalGithubKey(data []byte) (GithubKey, error) {
	var r GithubKey
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *githubKeys) marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GithubKey) marshal() ([]byte, error) {
	return json.Marshal(r)
}

// GithubKey is the JSON representation of a Github deploy key.
type GithubKey struct {
	ID        int64  `json:"id"`
	Key       string `json:"key"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Verified  bool   `json:"verified"`
	CreatedAt string `json:"created_at"`
	ReadOnly  bool   `json:"read_only"`
}

// GetKeyError is an error wrapper for get operations
type GetKeyError struct {
	err string
}

func (e *GetKeyError) Error() string {
	return fmt.Sprintf("%s", e.err)
}

// DeleteKeyError is an error wrapper for delete
type DeleteKeyError struct {
	err string
}

func (e *DeleteKeyError) Error() string {
	return fmt.Sprintf("%s", e.err)
}

// CreateKeyError is an error wrapper for create operations
type CreateKeyError struct {
	err string
}

func (e *CreateKeyError) Error() string {
	return fmt.Sprintf("%s", e.err)
}

// Doer is an interface used to make testing easier
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// GetDeployKey returns the deploy key for a repo that matches the keyTitle parameter. It returns -1 if the key title is not found.
func GetDeployKey(client Doer, githubUsername, githubPassword, repo, keyTitle string) (GithubKey, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+githubUsername+"/"+repo+"/keys", nil)
	req.SetBasicAuth(githubUsername, githubPassword)
	resp, err := client.Do(req)
	if err != nil {
		return GithubKey{}, fmt.Errorf("unable to get keys from github: %w", &GetKeyError{err: err.Error()})
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GithubKey{}, fmt.Errorf("unable to read response: %w", &GetKeyError{err: err.Error()})
	}

	githubKeys, err := unmarshalGithubKeys(bodyBytes)
	if err != nil {
		return GithubKey{}, err
	}

	for _, key := range githubKeys {
		if key.Title == keyTitle {
			return key, nil
		}
	}

	return GithubKey{}, nil
}

// DeleteDeployKey deletes a GitHub deploy key matching the keyID parameter. It returns true after a successful delete and false if the delete is unsuccessful.
func DeleteDeployKey(client Doer, githubUsername, githubPassword, repo string, keyID int64) error {

	req, err := http.NewRequest("DELETE", "https://api.github.com/repos/"+githubUsername+"/"+repo+"/keys/"+strconv.FormatInt(keyID, 10), nil)
	req.SetBasicAuth(githubUsername, githubPassword)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to delete key id %d: %w", keyID, &DeleteKeyError{err: err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("unable to delete key id %d: %w", keyID, &DeleteKeyError{err: ""})
}

// CreateDeployKey creates a new GitHub deploy key.
func CreateDeployKey(client Doer, githubUsername, githubPassword, repo, keyTitle, newKey string, readOnly bool) (GithubKey, error) {
	githubKey := &GithubKey{
		Title:    keyTitle,
		Key:      newKey,
		ReadOnly: readOnly,
	}

	requestBody, err := githubKey.marshal()

	req, err := http.NewRequest("POST", "https://api.github.com/repos/"+githubUsername+"/"+repo+"/keys", bytes.NewBuffer(requestBody))
	req.SetBasicAuth(githubUsername, githubPassword)
	resp, err := client.Do(req)
	if err != nil {
		return GithubKey{}, fmt.Errorf("error creating key during http request: %w", &CreateKeyError{err: err.Error()})
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GithubKey{}, fmt.Errorf("error reading response: %w", &CreateKeyError{err: err.Error()})
	}
	newGithubKey, err := unmarshalGithubKey(bodyBytes)
	if err != nil {
		return GithubKey{}, fmt.Errorf("error unmarshalling response: %w", &CreateKeyError{err: err.Error()})
	}

	if resp.StatusCode == 201 {
		return newGithubKey, nil
	}

	return GithubKey{}, fmt.Errorf("error http status code %d: %w", resp.StatusCode, &CreateKeyError{err: ""})
}
