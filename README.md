![example workflow](https://github.com/hotafrika/ebay-common/actions/workflows/autotests.yml/badge.svg)


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

```go
package main

import (
	"fmt"
	"github.com/hotafrika/ebay-common/datetime"
)

func main() {
	ebayDatetime := "2007-07-24T21:05:05.781Z"
	normalDatetime, _ := datetime.FromEbayDateTime(ebayDatetime)

	againEbayDatetime := datetime.ToEbayDateTime(normalDatetime)
	fmt.Println(againEbayDatetime) // "2007-07-24T21:05:05.781Z"
}
```