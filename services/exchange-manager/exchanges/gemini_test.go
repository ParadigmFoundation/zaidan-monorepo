package exchanges

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestGemini(t *testing.T) {
	cfg := GeminiConf{
		BaseURL: "https://api.sandbox.gemini.com",
		Key:     os.Getenv("EM_GEMINI_TEST_KEY"),
		Secret:  os.Getenv("EM_GEMINI_TEST_SECRET"),
	}
	if cfg.Key == "" && cfg.Secret == "" {
		t.Skip("EM_GEMINI_TEST_{KEY|SECRET} env not defined")
	}

	exchange := NewGemini(cfg)
	suite.Run(t, &ExchangesSuite{exchange: exchange})
}
