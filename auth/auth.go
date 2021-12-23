package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
	"time"
)

const ProductionURL string = "https://api.ebay.com/identity/v1/oauth2/token"
const SandboxURL string = "https://api.sandbox.ebay.com/identity/v1/oauth2/token"

const (
	ScopeCredentialCommon string = "https://api.ebay.com/oauth/api_scope"
	// TODO more scopes here
)

// ApplicationToken represents eBay application token object.
type ApplicationToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
}

// Service is a service for creating ApplicationToken.
// HTTP client's default timeout is 10 seconds.
type Service struct {
	authString string
	URL        string
	scopes     string
	scopesMap  map[string]struct{}
	client     *resty.Client
}

func newService(authString, URL string) *Service {
	return &Service{
		authString: authString,
		URL:        URL,
		client:     resty.New().SetTimeout(10 * time.Second),
	}
}

// NewService creates new Service with empty credentials.
// Default API endpoint is ProductionURL.
func NewService() *Service {
	s := createB64String("", "")

	c := newService(s, ProductionURL)
	return c
}

// NewServiceProd creates new Service.
func NewServiceProd(clientID, clientSecret string) *Service {
	s := createB64String(clientID, clientSecret)

	return newService(s, ProductionURL)
}

// NewServiceSandbox creates new Service.
func NewServiceSandbox(clientID, clientSecret string) *Service {
	s := createB64String(clientID, clientSecret)

	return newService(s, SandboxURL)
}

// NewServiceCustom creates new Service.
func NewServiceCustom(clientID, clientSecret, url string) *Service {
	s := createB64String(clientID, clientSecret)

	return newService(s, url)
}

// WithURL changes endpoint for API calls. By default ProductionURL.
func (s *Service) WithURL(endpoint string) *Service {
	s.URL = endpoint
	return s
}

// WithCredentials changes credentials.
func (s *Service) WithCredentials(clientID, clientSecret string) *Service {
	authString := createB64String(clientID, clientSecret)
	s.authString = authString
	return s
}

// WithTimeout adds timeout to HTTP client.
// Default: 10 seconds.
func (s *Service) WithTimeout(timeout time.Duration) *Service {
	s.client.SetTimeout(timeout)
	return s
}

// WithScopes adds unique scopes to Service. Does not delete previous.
func (s *Service) WithScopes(scopes ...string) *Service {
	if len(s.scopesMap) == 0 {
		s.scopesMap = make(map[string]struct{})
	}
	for _, scope := range scopes {
		s.scopesMap[scope] = struct{}{}
	}

	var scopesEvaluated []string
	for key := range s.scopesMap {
		scopesEvaluated = append(scopesEvaluated, key)
	}
	scopesString := strings.Join(scopesEvaluated, " ")
	//scopesString = url.PathEscape(scopesString)

	s.scopes = scopesString
	return s
}

// GetScopes returns the list of scopes for Service
func (s *Service) GetScopes() []string {
	var scopes []string
	for key := range s.scopesMap {
		scopes = append(scopes, key)
	}
	return scopes
}

// GetAppToken creates ApplicationToken.
// Credentials and scopes are taken from Service.
// You can use it only if you set up credentials for Service with WithCredentials etc.
func (s *Service) GetAppToken() (ApplicationToken, error) {
	return s.getAppTokenWithAuth(s.authString)
}

// GetAppTokenWithCredentials creates ApplicationToken by credentials.
// This method doesn't change credentials inside Service.
// Scopes are taken from Service.
func (s *Service) GetAppTokenWithCredentials(clientID, clientSecret string) (ApplicationToken, error) {
	authString := createB64String(clientID, clientSecret)
	return s.getAppTokenWithAuth(authString)
}

func (s *Service) getAppTokenWithAuth(authString string) (ApplicationToken, error) {
	at := ApplicationToken{}
	res, err := s.client.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Authorization", authString).
		SetFormData(map[string]string{
			"grant_type": "client_credentials",
			"scope":      s.scopes,
		}).Post(s.URL)
	if err != nil {
		return at, err
	}
	if res.StatusCode() != 200 {
		return at, fmt.Errorf("response code is %d", res.StatusCode())
	}

	err = json.Unmarshal(res.Body(), &at)
	if err != nil {
		return ApplicationToken{}, fmt.Errorf("unable to unmarshall response body")
	}
	return at, nil
}

func createB64String(clientID, clientSecret string) string {
	s := strings.Join([]string{clientID, clientSecret}, ":")
	b64 := base64.StdEncoding.EncodeToString([]byte(s))
	s = strings.Join([]string{"Basic", b64}, " ")
	return s
}
