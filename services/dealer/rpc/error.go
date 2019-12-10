package rpc

// Standard JSONRPC errors
var (
	ErrorParseJSON = NewError("unable to parse: invalid JSON", -32700)
)

// Error represents a JSONRPC error and satisfies the rpc.Error interface (go-ethereum/rpc)
type Error struct {
	message string
	code    int
}

// NewError creates a new Error with given message and code
func NewError(message string, code int) Error { return Error{message, code} }

func (e Error) Error() string  { return e.message }
func (e Error) ErrorCode() int { return e.code }
