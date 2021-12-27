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
	// Client Credential Grant Type Scopes
	ScopeCredentialCommon                 string = "https://api.ebay.com/oauth/api_scope"
	ScopeCredentialBuyGuestOrder          string = "https://api.ebay.com/oauth/api_scope/buy.guest.order"
	ScopeCredentialBuyMarketing           string = "https://api.ebay.com/oauth/api_scope/buy.marketing"
	ScopeCredentialBuyProductFeed         string = "https://api.ebay.com/oauth/api_scope/buy.product.feed"
	ScopeCredentialBuyMarketplaceInsights string = "https://api.ebay.com/oauth/api_scope/buy.marketplace.insights"
	ScopeCredentialBuyProxyGuestOrder     string = "https://api.ebay.com/oauth/api_scope/buy.proxy.guest.order"
	ScopeCredentialBuyItemBulk            string = "https://api.ebay.com/oauth/api_scope/buy.item.bulk"

	// Authorization Code Grant Type Scopes
	ScopeAuthorizationCommon                                   string = "https://api.ebay.com/oauth/api_scope"
	ScopeAuthorizationBuyOrderReadonly                         string = "https://api.ebay.com/oauth/api_scope/buy.order.readonly"
	ScopeAuthorizationBuyGuestOrder                            string = "https://api.ebay.com/oauth/api_scope/buy.guest.order"
	ScopeAuthorizationSellMarketingReadonly                    string = "https://api.ebay.com/oauth/api_scope/sell.marketing.readonly"
	ScopeAuthorizationSellMarketing                            string = "https://api.ebay.com/oauth/api_scope/sell.marketing"
	ScopeAuthorizationSellInventoryReadonly                    string = "https://api.ebay.com/oauth/api_scope/sell.inventory.readonly"
	ScopeAuthorizationSellInventory                            string = "https://api.ebay.com/oauth/api_scope/sell.inventory"
	ScopeAuthorizationSellAccountReadonly                      string = "https://api.ebay.com/oauth/api_scope/sell.account.readonly"
	ScopeAuthorizationSellAccount                              string = "https://api.ebay.com/oauth/api_scope/sell.account"
	ScopeAuthorizationSellFulfillmentReadonly                  string = "https://api.ebay.com/oauth/api_scope/sell.fulfillment.readonly"
	ScopeAuthorizationSellFulfillment                          string = "https://api.ebay.com/oauth/api_scope/sell.fulfillment"
	ScopeAuthorizationSellAnalyticsReadonly                    string = "https://api.ebay.com/oauth/api_scope/sell.analytics.readonly"
	ScopeAuthorizationSellMarketplaceInsightsReadonly          string = "https://api.ebay.com/oauth/api_scope/sell.marketplace.insights.readonly"
	ScopeAuthorizationCommerceCatalogReadonly                  string = "https://api.ebay.com/oauth/api_scope/commerce.catalog.readonly"
	ScopeAuthorizationBuyShoppingCart                          string = "https://api.ebay.com/oauth/api_scope/buy.shopping.cart"
	ScopeAuthorizationBuyOfferAuction                          string = "https://api.ebay.com/oauth/api_scope/buy.offer.auction"
	ScopeAuthorizationCommerceIdentityReadonly                 string = "https://api.ebay.com/oauth/api_scope/commerce.identity.readonly"
	ScopeAuthorizationCommerceIdentityEmailReadonly            string = "https://api.ebay.com/oauth/api_scope/commerce.identity.email.readonly"
	ScopeAuthorizationCommerceIdentityPhoneReadonly            string = "https://api.ebay.com/oauth/api_scope/commerce.identity.phone.readonly"
	ScopeAuthorizationCommerceIdentityAddressReadonly          string = "https://api.ebay.com/oauth/api_scope/commerce.identity.address.readonly"
	ScopeAuthorizationCommerceIdentityNameReadonly             string = "https://api.ebay.com/oauth/api_scope/commerce.identity.name.readonly"
	ScopeAuthorizationCommerceIdentityStatusReadonly           string = "https://api.ebay.com/oauth/api_scope/commerce.identity.status.readonly"
	ScopeAuthorizationSellFinances                             string = "https://api.ebay.com/oauth/api_scope/sell.finances"
	ScopeAuthorizationSellItemDraft                            string = "https://api.ebay.com/oauth/api_scope/sell.item.draft"
	ScopeAuthorizationSellPaymentDispute                       string = "https://api.ebay.com/oauth/api_scope/sell.payment.dispute"
	ScopeAuthorizationSellItem                                 string = "https://api.ebay.com/oauth/api_scope/sell.item"
	ScopeAuthorizationSellReputation                           string = "https://api.ebay.com/oauth/api_scope/sell.reputation"
	ScopeAuthorizationSellReputationReadonly                   string = "https://api.ebay.com/oauth/api_scope/sell.reputation.readonly"
	ScopeAuthorizationCommerceNotificationSubscription         string = "https://api.ebay.com/oauth/api_scope/commerce.notification.subscription"
	ScopeAuthorizationCommerceNotificationSubscriptionReadonly string = "https://api.ebay.com/oauth/api_scope/commerce.notification.subscription.readonly"
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
