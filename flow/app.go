package main

import (
	"context"
	"fmt"

	"github.com/google/martian/log"
	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/rx"
)

var printer = func(_ context.Context, i interface{}) (interface{}, error) {
	value := i.(int64)
	fmt.Println("serverless get value:", value)
	return value, nil
}

var callback = func(v []byte) (interface{}, error) {
	res, _, _, err := y3.DecodePrimitivePacket(v)
	if err != nil {
		log.Errorf("y3 err: %v", err)
	}

	val, err := res.ToInt64()
	if err != nil {
		log.Errorf("y3 toint64 err: %v", err)
	}
	return val, nil
}

// Handler will handle data in Rx way
func Handler(rxstream rx.RxStream) rx.RxStream {
	log.Debugf("-----")
	stream := rxstream.
		Subscribe(0x01).
		OnObserve(callback).
		Map(printer).
		StdOut()

	return stream
}
