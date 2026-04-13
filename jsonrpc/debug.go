package jsonrpc

import (
	"encoding/json"

	"github.com/boolw/go-web3"
)

// Debug is the debug namespace
type Debug struct {
	c *Client
}

// Debug returns the reference to the debug namespace
func (c *Client) Debug() *Debug {
	return c.endpoints.d
}

// TraceCall executes the given call and returns the structured execution trace.
func (d *Debug) TraceCall(msg *web3.CallMsg, block web3.BlockNumber, options map[string]interface{}) (json.RawMessage, error) {
	var out json.RawMessage
	if options == nil {
		err := d.c.Call("debug_traceCall", &out, msg, block.String())
		return out, err
	}
	err := d.c.Call("debug_traceCall", &out, msg, block.String(), options)
	return out, err
}

// TraceTransaction returns the structured execution trace for the given transaction hash.
func (d *Debug) TraceTransaction(hash web3.Hash, options map[string]interface{}) (json.RawMessage, error) {
	var out json.RawMessage
	if options == nil {
		err := d.c.Call("debug_traceTransaction", &out, hash)
		return out, err
	}
	err := d.c.Call("debug_traceTransaction", &out, hash, options)
	return out, err
}
