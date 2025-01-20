package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/crypto-org-chain/cronos/x/biol1/types"
)

type QueryServer struct {
	Keeper
}

var _ types.QueryServer = QueryServer{}

func NewQueryServer(k Keeper) QueryServer {
	return QueryServer{Keeper: k}
}

func (qs QueryServer) GCContent(ctx context.Context, req *types.GCContentRequest) (*types.GCContentResponse, error) {
	if qs.sidecarURL == "" {
		return nil, fmt.Errorf("Biopython sidecar URL is empty")
	}

	bodyData, err := json.Marshal(map[string]string{"sequence": req.Sequence})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// POST to sidecar
	httpReq, err := http.NewRequest("POST", qs.sidecarURL+"/gc-content", bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed calling sidecar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sidecar returned %d", resp.StatusCode)
	}

	var result struct {
		GCContent float64 `json:"gc_content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed decoding response: %w", err)
	}

	return &types.GCContentResponse{GcContent: result.GCContent}, nil
}
