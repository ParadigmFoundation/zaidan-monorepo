package tests

type Quote struct {
	QuoteId               string `json:"quoteId"`
	MakerAssetAddress     string `json:"makerAssetAddress"`
	TakerAssetAddress     string `json:"takerAssetAddress"`
	MakerAssetSize        string `json:"makerAssetSize"`
	TakerAssetSize        string `json:"takerAssetSize"`
	Expiration            int64  `json:"expiration"`
	ServerTime            int64  `json:"serverTime"`
	ZeroExTransactionHash string `json:"zeroExTransactionHash"`
	ZeroExTransactionInfo struct {
		Order struct {
			Signature string `json:"signature"`
		} `json:"order"`
		Transaction struct {
			Data                  string `json:"data"`
			ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
			GasPrice              string `json:"gasPrice"`
			Salt                  string `json:"salt"`
			SignerAddress         string `json:"signerAddress"`
		} `json:"transaction"`
	} `json:"zeroExTransactionInfo"`
}
