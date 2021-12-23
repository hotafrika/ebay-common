# ebay-common
The repo contains common tools for using my other eBay libs.

### 1. Auth token generation

```go
package main

import (
	"fmt"
	"github.com/hotafrika/ebay-common/auth"
)

func main() {
	clientID := "xxx"
	clientSecret := "xxxx"

	service := auth.NewService().
		WithScopes(auth.ScopeCredentialCommon).
		WithCredentials(clientID, clientSecret)
	
	token, err := service.GetAppToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	
	// or make reusable service for generating tokens by different credentials
	reusableService := auth.NewService().WithScopes(auth.ScopeCredentialCommon)
	token1, err := reusableService.GetAppTokenWithCredentials("aaa", "aaaaa")
	token2, err := reusableService.GetAppTokenWithCredentials("bbb", "bbbb")
	fmt.Println(token1, token2)
}
```

### 2. Encoding/decoding eBay datetime values
In progress...