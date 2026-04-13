package jsonrpc

import (
	"encoding/json"
	"testing"

	"github.com/boolw/go-web3"
	"github.com/stretchr/testify/assert"
)

type recorderTransport struct {
	method string
	params []interface{}
}

func (r *recorderTransport) Call(method string, out interface{}, params ...interface{}) error {
	r.method = method
	r.params = params

	if raw, ok := out.(*json.RawMessage); ok {
		*raw = json.RawMessage(`{"ok":true}`)
	}
	return nil
}

func (r *recorderTransport) Close() error {
	return nil
}

func TestDebugTraceTransaction(t *testing.T) {
	c, err := NewClient("http://localhost")
	assert.NoError(t, err)
	defer c.Close()

	r := &recorderTransport{}
	c.SetTransport(r)

	h := web3.Hash{0x1}

	_, err = c.Debug().TraceTransaction(h, nil)
	assert.NoError(t, err)
	assert.Equal(t, "debug_traceTransaction", r.method)
	assert.Len(t, r.params, 1)
	assert.Equal(t, h, r.params[0])

	opts := map[string]interface{}{"disableStack": true}
	_, err = c.Debug().TraceTransaction(h, opts)
	assert.NoError(t, err)
	assert.Equal(t, "debug_traceTransaction", r.method)
	assert.Len(t, r.params, 2)
	assert.Equal(t, h, r.params[0])
	assert.Equal(t, opts, r.params[1])
}

func TestDebugTraceCall(t *testing.T) {
	c, err := NewClient("http://localhost")
	assert.NoError(t, err)
	defer c.Close()

	r := &recorderTransport{}
	c.SetTransport(r)

	msg := &web3.CallMsg{
		From: web3.Address{0x1},
		To:   web3.Address{0x2},
		Data: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	_, err = c.Debug().TraceCall(msg, web3.Latest, nil)
	assert.NoError(t, err)
	assert.Equal(t, "debug_traceCall", r.method)
	assert.Len(t, r.params, 2)
	assert.Equal(t, msg, r.params[0])
	assert.Equal(t, web3.Latest.String(), r.params[1])

	opts := map[string]interface{}{"tracer": "callTracer"}
	_, err = c.Debug().TraceCall(msg, web3.Latest, opts)
	assert.NoError(t, err)
	assert.Equal(t, "debug_traceCall", r.method)
	assert.Len(t, r.params, 3)
	assert.Equal(t, msg, r.params[0])
	assert.Equal(t, web3.Latest.String(), r.params[1])
	assert.Equal(t, opts, r.params[2])
}
