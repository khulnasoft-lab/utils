package khulnasoft

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/khulnasoft-lab/retryablehttp-go"
	"github.com/khulnasoft-lab/utils/env"
	fileutil "github.com/khulnasoft-lab/utils/file"
	folderutil "github.com/khulnasoft-lab/utils/folder"
	urlutil "github.com/khulnasoft-lab/utils/url"
)

var (
	KHULNASOFTDir      = filepath.Join(folderutil.HomeDirOrDefault(""), ".khulnasoft")
	KHULNASOFTCredFile = filepath.Join(KHULNASOFTDir, "credentials.yaml")
	ErrNoCreds         = fmt.Errorf("no credentials found in %s", KHULNASOFTDir)
)

const (
	userProfileURL   = "https://%s/v1/user?utm_source=%s"
	apiKeyEnv        = "KHULNASOFT_API_KEY"
	apiServerEnv     = "KHULNASOFT_API_SERVER"
	ApiKeyHeaderName = "X-Api-Key"
	dashBoardEnv     = "KHULNASOFT_DASHBOARD_URL"
)

type KHULNASOFTCredentials struct {
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
	APIKey   string `yaml:"api-key"`
	Server   string `yaml:"server"`
}

type KHULNASOFTUserProfileResponse struct {
	UserName string `json:"name"`
	Email    string `json:"email"`
	// there are more fields but we don't need them
	/// below fields are added later on and not part of the response
}

// KHULNASOFTCredHandler is interface for adding / retrieving khulnasoft credentials
// from file system
type KHULNASOFTCredHandler struct{}

// GetCreds retrieves the credentials from the file system or environment variables
func (p *KHULNASOFTCredHandler) GetCreds() (*KHULNASOFTCredentials, error) {
	credsFromEnv := p.getCredsFromEnv()
	if credsFromEnv != nil {
		return credsFromEnv, nil
	}
	if !fileutil.FolderExists(KHULNASOFTDir) || !fileutil.FileExists(KHULNASOFTCredFile) {
		return nil, ErrNoCreds
	}
	bin, err := os.Open(KHULNASOFTCredFile)
	if err != nil {
		return nil, err
	}
	// for future use-cases
	var creds []KHULNASOFTCredentials
	err = yaml.NewDecoder(bin).Decode(&creds)
	if err != nil {
		return nil, err
	}
	if len(creds) == 0 {
		return nil, ErrNoCreds
	}
	return &creds[0], nil
}

// getCredsFromEnv retrieves the credentials from the environment
// if not or incomplete credentials are found it return nil
func (p *KHULNASOFTCredHandler) getCredsFromEnv() *KHULNASOFTCredentials {
	apiKey := env.GetEnvOrDefault(apiKeyEnv, "")
	apiServer := env.GetEnvOrDefault(apiServerEnv, DefaultApiServer)
	if apiKey == "" || apiServer == "" {
		return nil
	}
	return &KHULNASOFTCredentials{APIKey: apiKey, Server: apiServer}
}

// SaveCreds saves the credentials to the file system
func (p *KHULNASOFTCredHandler) SaveCreds(resp *KHULNASOFTCredentials) error {
	if resp == nil {
		return fmt.Errorf("invalid response")
	}
	if !fileutil.FolderExists(KHULNASOFTDir) {
		_ = fileutil.CreateFolder(KHULNASOFTDir)
	}
	bin, err := yaml.Marshal([]*KHULNASOFTCredentials{resp})
	if err != nil {
		return err
	}
	return os.WriteFile(KHULNASOFTCredFile, bin, 0600)
}

// ValidateAPIKey validates the api key and retrieves associated user metadata like username
// from given api server/host
func (p *KHULNASOFTCredHandler) ValidateAPIKey(key string, host string, toolName string) (*KHULNASOFTCredentials, error) {
	// get address from url
	urlx, err := urlutil.Parse(host)
	if err != nil {
		return nil, err
	}
	req, err := retryablehttp.NewRequest("GET", fmt.Sprintf(userProfileURL, urlx.Host, toolName), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(ApiKeyHeaderName, key)
	resp, err := retryablehttp.DefaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	bin, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var profile KHULNASOFTUserProfileResponse
	err = json.Unmarshal(bin, &profile)
	if err != nil {
		return nil, err
	}
	if profile.Email == "" {
		return nil, fmt.Errorf("invalid response from server got %v", string(bin))
	}
	return &KHULNASOFTCredentials{Username: profile.UserName, Email: profile.Email, APIKey: key, Server: host}, nil
}

func init() {
	DashBoardURL = env.GetEnvOrDefault("KHULNASOFT_DASHBOARD_URL", DashBoardURL)
}
