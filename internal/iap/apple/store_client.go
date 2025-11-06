package apple

import (
	"fmt"
	"os"

	"github.com/awa/go-iap/appstore/api"
)

var storeClient *api.StoreClient

func NewAppleIAPStoreClient(keyPath, bundleID, keyID, issuer string) error {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("reading file content failed: %w", err)
	}
	if bundleID == "" {
		return fmt.Errorf("invalid bundle id")
	}
	if keyID == "" {
		return fmt.Errorf("invalid key id")
	}
	if issuer == "" {
		return fmt.Errorf("invalid issuer")
	}

	c := &api.StoreConfig{
		KeyContent: keyBytes,
		KeyID:      keyID,
		BundleID:   bundleID,
		Issuer:     issuer,
		Sandbox:    os.Getenv("ENV") == "DEV",
	}
	storeClient = api.NewStoreClient(c)

	return nil
}

func IAPStoreClient() *api.StoreClient { return storeClient }
