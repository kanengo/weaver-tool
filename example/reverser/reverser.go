package main

import (
	"context"
	"github.com/ServiceWeaver/weaver"
)

type Reverser interface {
	Reverse(context.Context, string) (string, error)
	GetInfo(ctx context.Context, info Info) (Info, error)
}

type reverser struct {
	weaver.Implements[Reverser]
}

type Info struct {
	weaver.AutoMarshal
	Id   int64
	Name string
	Age  uint32
}

func (*reverser) Reverse(_ context.Context, s string) (string, error) {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}

	return string(runes), nil
}

func (*reverser) GetInfo(ctx context.Context, info Info) (Info, error) {
	return info, nil
}
