package main

import (
	"context"
	"fmt"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/rx"
)

const observedKey = 0x10

// Decode data by Y3 as Int64 type from YoMo-Zipper
var decode = func(v []byte) (interface{}, error) {
	val, err := y3.ToInt64(v)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return val, nil
}

// Calculate sum when every time new data arrived
var sum = func(_ context.Context, acc interface{}, elem interface{}) (interface{}, error) {
	if acc == nil {
		return elem, nil
	}
	return acc.(int64) + elem.(int64), nil
}

// Handler will handle data in Rx way
func Handler(rxstream rx.RxStream) rx.RxStream {
	stream := rxstream.
		Subscribe(observedKey).
		OnObserve(decode).
		Scan(sum).
		StdOut()

	return stream
}
