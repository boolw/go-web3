package jsonrpc

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/boolw/go-web3"
)

// Eth is the eth namespace
type Eth struct {
	c *Client
}

// Eth returns the reference to the eth namespace
func (c *Client) Eth() *Eth {
	return c.endpoints.e
}

// Accounts returns a list of addresses owned by client.
func (e *Eth) Accounts() ([]web3.Address, error) {
	out := make([]web3.Address, 0)
	if err := e.c.Call("eth_accounts", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BlockNumber returns the number of most recent block.
func (e *Eth) BlockNumber() (uint64, error) {
	var out string
	if err := e.c.Call("eth_blockNumber", &out); err != nil {
		return 0, err
	}
	return parseUint64orHex(out)
}

// GetBlockByNumber returns information about a block by block number.
func (e *Eth) GetBlockByNumber(i web3.BlockNumber, full bool) (*web3.Block, error) {
	b := new(web3.Block)
	if err := e.c.Call("eth_getBlockByNumber", b, i.String(), full); err != nil {
		return nil, err
	}
	return b, nil
}

// GetTransactionByHash returns information about a block by hash.
func (e *Eth) GetTransactionByHash(hash web3.Hash) (*web3.Transaction, error) {
	b := new(web3.Transaction)
	if err := e.c.Call("eth_getTransactionByHash", b, hash); err != nil {
		return nil, err
	}
	return b, nil
}

// GetBlockByHash returns information about a block by hash.
func (e *Eth) GetBlockByHash(hash web3.Hash, full bool) (*web3.Block, error) {
	b := new(web3.Block)
	if err := e.c.Call("eth_getBlockByHash", b, hash, full); err != nil {
		return nil, err
	}
	return b, nil
}

// SendTransaction creates new message call transaction or a contract creation.
func (e *Eth) SendTransaction(txn *web3.Transaction) (web3.Hash, error) {
	var hash web3.Hash
	err := e.c.Call("eth_sendTransaction", &hash, txn)
	return hash, err
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash.
func (e *Eth) GetTransactionReceipt(hash web3.Hash) (*web3.Receipt, error) {
	receipt := new(web3.Receipt)
	err := e.c.Call("eth_getTransactionReceipt", receipt, hash)
	return receipt, err
}

// GetNonce returns the nonce of the account
func (e *Eth) GetNonce(addr web3.Address, blockNumber web3.BlockNumber) (uint64, error) {
	var nonce string
	if err := e.c.Call("eth_getTransactionCount", &nonce, addr, blockNumber.String()); err != nil {
		return 0, err
	}
	return parseUint64orHex(nonce)
}

// GetBalance returns the balance of the account of given address.
func (e *Eth) GetBalance(addr web3.Address, blockNumber web3.BlockNumber) (*big.Int, error) {
	var out string
	if err := e.c.Call("eth_getBalance", &out, addr, blockNumber.String()); err != nil {
		return nil, err
	}
	b, ok := new(big.Int).SetString(out[2:], 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert to big.int")
	}
	return b, nil
}

// GasPrice returns the current price per gas in wei.
func (e *Eth) GasPrice() (uint64, error) {
	var out string
	if err := e.c.Call("eth_gasPrice", &out); err != nil {
		return 0, err
	}
	return parseUint64orHex(out)
}

// Call executes a new message call immediately without creating a transaction on the block chain.
func (e *Eth) Call(msg *web3.CallMsg, block web3.BlockNumber) (string, error) {
	var out string
	if err := e.c.Call("eth_call", &out, msg, block.String()); err != nil {
		return "", err
	}
	return out, nil
}

// EstimateGasContract estimates the gas to deploy a contract
func (e *Eth) EstimateGasContract(bin []byte) (uint64, error) {
	var out string
	msg := map[string]interface{}{
		"data": "0x" + hex.EncodeToString(bin),
	}
	if err := e.c.Call("eth_estimateGas", &out, msg); err != nil {
		return 0, err
	}
	return parseUint64orHex(out)
}

// EstimateGas generates and returns an estimate of how much gas is necessary to allow the transaction to complete.
func (e *Eth) EstimateGas(msg *web3.CallMsg) (uint64, error) {
	var out string
	if err := e.c.Call("eth_estimateGas", &out, msg); err != nil {
		return 0, err
	}
	return parseUint64orHex(out)
}

// GetLogs returns an array of all logs matching a given filter object
func (e *Eth) GetLogs(filter *web3.LogFilter) ([]*web3.Log, error) {
	var out []*web3.Log
	if err := e.c.Call("eth_getLogs", &out, filter); err != nil {
		return nil, err
	}
	return out, nil
}

// ChainID returns the id of the chain
func (e *Eth) ChainID() (*big.Int, error) {
	var out string
	if err := e.c.Call("eth_chainId", &out); err != nil {
		return nil, err
	}
	return parseBigInt(out), nil
}

func (e *Eth) GetStorageAt(addr web3.Address, hash web3.Hash, blockNumber web3.BlockNumber) (string, error) {
	var out string
	if err := e.c.Call("eth_getStorageAt", &out, addr, hash, blockNumber.String()); err != nil {
		return "", err
	}
	return out, nil
}
