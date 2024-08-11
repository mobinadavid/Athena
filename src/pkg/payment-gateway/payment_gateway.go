package payment_gateway

import (
	"athena/src/config"
	"athena/src/pkg/payment-gateway/drivers/fiat"
	"athena/src/pkg/vault"
	"context"
	"errors"
	"fmt"
	"strings"
)

// NewIPG creates a new IPG driver based on the provided driver name.
func NewIPG(driver ...string) (fiat.IIPGDriver, error) {
	driverName := getDriverName(driver...)

	secrets, err := vault.GetInstance().GetVault().KVv2("kv-v2").Get(context.Background(), fmt.Sprintf("%s/ipg/%s", config.GetInstance().Get("APP_NAME"), driverName))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secrets from vault: %w", err)
	}

	if secrets == nil {
		return nil, errors.New("failed to retrieve secrets from vault")
	}

	apiKey, ok := secrets.Data["api-key"].(string)
	if !ok {
		return nil, errors.New("failed to retrieve secrets from vault")
	}

	switch driverName {
	case "sep":
		return fiat.NewSep(apiKey)
	default:
		return nil, fmt.Errorf("unsupported IPG driver: %s", driverName)
	}
}

func getDriverName(driver ...string) string {
	if len(driver) > 0 && driver[0] != "" {
		return strings.ToLower(driver[0])
	}

	envDriver := config.GetInstance().Get("IPG_ACTIVE_GATEWAY")
	if envDriver != "" {
		return strings.ToLower(envDriver)
	}

	return "sep"
}
