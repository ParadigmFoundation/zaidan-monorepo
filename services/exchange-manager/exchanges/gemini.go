package exchanges

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/golang/protobuf/ptypes/empty"
)

type GeminiConf struct {
	BaseURL    string
	Key        string
	Secret     string
	APIVersion string
}

type Gemini struct {
	client *http.Client
	cfg    GeminiConf
}

func NewGemini(cfg GeminiConf) *Gemini {
	if cfg.APIVersion == "" {
		cfg.APIVersion = "v1"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.gemini.com"
	}
	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL = cfg.BaseURL + "/"
	}

	return &Gemini{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (g *Gemini) Post(ctx context.Context, path string, params map[string]interface{}, out interface{}) error {
	reqMap := map[string]interface{}{
		"request": fmt.Sprintf("/%s/%s", g.cfg.APIVersion, path),
		"nonce":   time.Now().UnixNano(),
	}

	for key, val := range params {
		reqMap[key] = val
	}

	jsonPayload, err := json.Marshal(reqMap)
	if err != nil {
		return err
	}
	payload := base64.StdEncoding.EncodeToString(jsonPayload)

	url := g.cfg.BaseURL + reqMap["request"].(string)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Content-Length":     "0",
		"Content-Type":       "text/plain",
		"X-GEMINI-APIKEY":    g.cfg.Key,
		"X-GEMINI-PAYLOAD":   payload,
		"X-GEMINI-SIGNATURE": hex.EncodeToString(g.hmac(payload)),
		"Cache-Control":      "no-cache",
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	resp, err := g.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		resp := &GeminiError{}
		if err := dec.Decode(resp); err != nil {
			return err
		}
		return fmt.Errorf("%s: %s", resp.Reason, resp.Message)
	}

	return dec.Decode(out)
}

func (g *Gemini) hmac(input string) []byte {
	hasher := hmac.New(sha512.New384, []byte(g.cfg.Secret))
	hasher.Write([]byte(input))
	return hasher.Sum(nil)
}

func (g *Gemini) CreateOrder(ctx context.Context, order *types.ExchangeOrder) error {
	req := map[string]interface{}{
		"symbol": g.toSymbol(order.Symbol),
		"amount": order.Amount,
		"price":  order.Price,
		"side":   strings.ToLower(order.Side.String()),
		"type":   "exchange limit",
	}

	resp := &GeminiOrder{}
	if err := g.Post(ctx, "order/new", req, resp); err != nil {
		return err
	}

	order.Id = fmt.Sprintf("%d", resp.ID)
	return nil
}

func (g *Gemini) GetOrder(ctx context.Context, id string) (*types.ExchangeOrderResponse, error) {
	req := map[string]interface{}{
		"order_id": id,
	}

	resp := &GeminiOrder{}
	if err := g.Post(ctx, "order/status", req, resp); err != nil {
		return nil, err
	}
	if resp.IsCancelled == true {
		return nil, errors.New("Order was cancelled")
	}

	return g.NewOrderResponse(resp)
}

func (g *Gemini) GetOpenOrders(ctx context.Context) (*types.ExchangeOrderArrayResponse, error) {
	var resp []*GeminiOrder
	if err := g.Post(ctx, "orders", nil, &resp); err != nil {
		return nil, err
	}

	var orders = &types.ExchangeOrderArrayResponse{}
	for _, item := range resp {
		order, err := g.NewOrderResponse(item)
		if err != nil {
			return nil, err
		}
		orders.Array = append(orders.Array, order)
	}

	return orders, nil
}

func (g *Gemini) CancelOrder(ctx context.Context, id string) (*empty.Empty, error) {
	req := map[string]interface{}{
		"order_id": id,
	}

	resp := &GeminiOrder{}
	if err := g.Post(ctx, "order/cancel", req, resp); err != nil {
		return nil, err
	}

	return nil, nil
}

func (g *Gemini) NewOrderResponse(order *GeminiOrder) (*types.ExchangeOrderResponse, error) {
	infoBytes, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	resp := &types.ExchangeOrderResponse{
		Order: &types.ExchangeOrder{
			Id:     fmt.Sprintf("%d", order.ID),
			Price:  fmt.Sprintf("%f", order.Price),
			Symbol: g.fromSymbol(order.Symbol),
			Amount: fmt.Sprintf("%f", order.OriginalAmount),
			Side:   SideFromString(order.Side),
		},
		Status: &types.ExchangeOrderStatus{
			Timestamp: order.Timestamp,
			Filled:    fmt.Sprintf("%f", order.ExecutedAmount),
			Info:      infoBytes,
		},
	}
	return resp, nil
}

func (g *Gemini) toSymbol(s string) string { return strings.ToLower(strings.Replace(s, "/", "", 1)) }

func (g *Gemini) fromSymbol(s string) string {
	if len(s) < 3 {
		return s
	}

	sym := s[:3] + "/" + s[3:]
	return strings.ToUpper(sym)
}
