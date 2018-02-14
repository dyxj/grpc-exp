package main

import (
	"context"
	"testing"

	hw "github.com/dyxj/grpc-exp/helloworld/helloworld"
)

func TestSayHello(t *testing.T) {
	s := server{}

	// Test data
	ttd := []struct {
		name   string
		expect string
	}{

		{
			name:   "Doodie",
			expect: "Hello World!! Doodie",
		},
		{
			name:   "From Narnia",
			expect: "Hello World!! From Narnia",
		},
	}

	for _, td := range ttd {
		req := &hw.HelloRequest{Name: td.name}
		resp, err := s.SayHello(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error(TestSayHello, SayHello): %v", err)
		}
		if resp.Message != td.expect {
			t.Errorf("expected: %v, got: %v", td.expect, resp.Message)
		}
	}
}

func TestSayBye(t *testing.T) {
	s := server{}

	// Test data
	ttd := []struct {
		name   string
		expect string
	}{
		{
			name:   "Doodie",
			expect: "Bye World!! Doodie",
		},
		{
			name:   "From Narnia",
			expect: "Bye World!! From Narnia",
		},
	}

	for _, td := range ttd {
		req := &hw.HelloRequest{Name: td.name}
		resp, err := s.SayBye(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error(TestSayBye, SayBye): %v", err)
		}
		if resp.Message != td.expect {
			t.Errorf("expected: %v, got: %v", td.expect, resp.Message)
		}
	}
}
