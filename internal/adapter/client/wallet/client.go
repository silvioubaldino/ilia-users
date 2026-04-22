package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
	"github.com/silvioubaldino/ilia-users/pkg/jwtutil"
)

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "CREDIT"
	TransactionTypeDebit  TransactionType = "DEBIT"
)

type Transaction struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Type      TransactionType `json:"type"`
	Amount    int64           `json:"amount"`
}

type Client struct {
	baseURL        string
	internalSecret string
	httpClient     *http.Client
}

func NewClient(baseURL, internalSecret string) *Client {
	return &Client{
		baseURL:        baseURL,
		internalSecret: internalSecret,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) generateToken(userID uuid.UUID) (string, error) {
	return jwtutil.GenerateToken(userID, "", c.internalSecret)
}

func (c *Client) GetBalance(ctx context.Context, userID uuid.UUID) (int64, error) {
	token, err := c.generateToken(userID)
	if err != nil {
		return 0, fmt.Errorf("wallet client: generate token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/balance", nil)
	if err != nil {
		return 0, fmt.Errorf("wallet client: new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("wallet client: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return 0, apperrors.ErrUnauthorized
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("wallet client: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Amount int64 `json:"amount"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("wallet client: decode response: %w", err)
	}

	return body.Amount, nil
}

func (c *Client) GetTransactions(ctx context.Context, userID uuid.UUID) ([]Transaction, error) {
	token, err := c.generateToken(userID)
	if err != nil {
		return nil, fmt.Errorf("wallet client: generate token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/transactions", nil)
	if err != nil {
		return nil, fmt.Errorf("wallet client: new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wallet client: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, apperrors.ErrUnauthorized
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet client: unexpected status %d", resp.StatusCode)
	}

	var transactions []Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, fmt.Errorf("wallet client: decode response: %w", err)
	}

	return transactions, nil
}
