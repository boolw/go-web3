package abi

import (
	"bytes"
	"fmt"
	"github.com/boolw/go-web3"
	"reflect"
	"testing"
)

func TestAbi(t *testing.T) {
	cases := []struct {
		Input  string
		Output *ABI
	}{
		{
			Input: `[
				{
					"name": "abc",
					"type": "function"
				},
				{
                    "anonymous": false,
                    "inputs": [],
                    "name": "Transfer",
                    "type": "event"
                }
			]`,
			Output: &ABI{
				Methods: map[string]*Method{
					"abc": {
						Name:    "abc",
						Inputs:  &Type{kind: KindTuple, raw: "tuple", tuple: []*TupleElem{}},
						Outputs: &Type{kind: KindTuple, raw: "tuple", tuple: []*TupleElem{}},
						id:      []byte{146, 39, 121, 51},
					},
				},
				Events: map[string]*Event{
					"Transfer": {
						Name:      "Transfer",
						Anonymous: false,
						Inputs:    &Type{kind: KindTuple, size: 0, raw: "tuple", tuple: []*TupleElem{
							//{
							//	Name:    "from",
							//	Elem:    &Type{kind: KindAddress, size: 0, raw: "address", t: reflect.TypeOf(web3.Address{}), tuple: []*TupleElem{}},
							//	Indexed: true,
							//},
							//{
							//	Name:    "to",
							//	Elem:    &Type{kind: KindAddress, size: 0, raw: "address", t: reflect.TypeOf(web3.Address{}), tuple: []*TupleElem{}},
							//	Indexed: true,
							//},
							//{
							//	Name:    "value",
							//	Elem:    &Type{kind: KindUInt, size: 256, raw: "uint256", t: reflect.TypeOf(big.NewInt(0)), tuple: []*TupleElem{}},
							//	Indexed: false,
							//},
						}},
						id: web3.HexToHash("0x406dade31f7ae4b5dbc276258c28dde5ae6d5c2773c5745802c493a2360e55e0"),
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			abi, err := NewABI(c.Input)
			if err != nil {
				t.Fatal(err)
			}
			for k, evt := range abi.Events {
				if evt.ID() != c.Output.Events[k].id {
					t.Fatal("bad")
				}
				fmt.Println(evt.ID())
			}
			for k, method := range abi.Methods {
				fmt.Println(method.ID())
				if bytes.Compare(method.ID(), c.Output.Methods[k].id) != 0 {
					t.Fatal("bad")
				}
			}
			if !reflect.DeepEqual(abi, c.Output) {
				t.Fatal("bad")
			}
		})
	}
}
