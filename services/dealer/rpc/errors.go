package rpc

// Standard JSONRPC errors
var (
	ErrParseJSON         = NewError("unable to parse: invalid JSON", -32700)
	ErrInvalidRequest    = NewError("invalid JSONRPC request", -32600)
	ErrInvalidParameters = NewError("invalid or missing method parameters", -32602)
	ErrInternal          = NewError("unknown internal error encountered", -32603)
)

// Zaidan Dealer JSONRPC errors
// https://github.com/ParadigmFoundation/zaidan-dealer-specification#error-codes
var (
	ErrInvalidFilter              = NewError("invalid or conflicting filter parameter(s)", -42002)
	ErrInvalidAddress             = NewError("invalid Ethereum address", -42003)
	ErrInvalidAssetData           = NewError("invalid token asset data", -42004)
	ErrTwoSizeRequests            = NewError("two size requests; expecting either maker or taker size", -42005)
	ErrUnauthorizedTaker          = NewError("taker address not authorized for trading", -42006)
	ErrInvalidSide                = NewError("unsupported trade side for requested market", -42007)
	ErrTemporaryRestriction       = NewError("temporary restriction on trading; try again later", -42008)
	ErrUnsupportedMarket          = NewError("unsupported maker asset", -42009)
	ErrUnsupportedTakerAsset      = NewError("unsupported taker asset for maker market", -42010)
	ErrQuoteTooLarge              = NewError("quote request size too large", -42011)
	ErrQuoteTooSmall              = NewError("quote request size too small", -42012)
	ErrQuoteUnavailable           = NewError("quotes unavailable at this time; try again later", -42013)
	ErrQuoteExpired               = NewError("quote has expired; request a new one", -42014)
	ErrUnknownQuoteID             = NewError("no known quote associated with provided UUID", -42015)
	ErrOrderFilled                = NewError("requested order has already been filled", -42016)
	ErrFillValidationFailed       = NewError("0x fill transaction validation failed; unable to execute", -42017)
	ErrInsufficientTakerBalance   = NewError("taker has insufficient token balance for fill", -42018)
	ErrInsufficientTakerAllowance = NewError("taker has insufficient 0x proxy allowance for fill", -42019)
	ErrQuoteValidationFailed      = NewError("quote validation failed; fill rejected", -42020)
	ErrInvalidTransactionID       = NewError("invalid or malformed Ethereum transaction hash", -42021)
	ErrInvalidOrderHash           = NewError("invalid or malformed 0x order hash", -42022)
	ErrInvalidUUID                = NewError("invalid or malformed UUID string", -42023)
	ErrRateLimitReached           = NewError("rate limit reached; try again later", -42024)
)

// Error represents a JSONRPC error
type Error struct {
	message string
	code    int
}

// NewError creates a new Error with given message and code
func NewError(message string, code int) Error { return Error{message, code} }

// Error implements rpc.Error (go-ethereum/rpc)
func (e Error) Error() string { return e.message }

// ErrorCode implements rpc.Error (go-ethereum/rpc)
func (e Error) ErrorCode() int { return e.code }
