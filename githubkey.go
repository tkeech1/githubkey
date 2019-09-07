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
		return GithubKey{}, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GithubKey{}, err
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
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("could not delete KeyID %d", keyID)
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
		return GithubKey{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GithubKey{}, err
	}
	newGithubKey, err := unmarshalGithubKey(bodyBytes)
	if err != nil {
		return GithubKey{}, err
	}

	if resp.StatusCode == 201 {
		return newGithubKey, nil
	}

	return GithubKey{}, fmt.Errorf("http status code: %d", resp.StatusCode)
}
