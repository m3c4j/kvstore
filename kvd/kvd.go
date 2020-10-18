// MIT License
//
// Copyright (c) 2020 The KVStore Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package kvd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/kvstore/api/kvserverpb"
	"github.com/kvstore/pkg/tree"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Main kvd server entry
func Main(args []string) error {
	fmt.Println("Hello there")

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		zap.L().Fatal("Cannot bind on port")
	}

	cfg := newConfig()
	cfg.parse(os.Args[1:])
	server := newKVServer(cfg)

	grpcServer := grpc.NewServer()
	pb.RegisterKVServer(grpcServer, server)

	// mux := http.NewServeMux()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := runtime.NewServeMux()
	gwAddr := ":8000" // internal grpc port, aka, lis port
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterKVHandlerFromEndpoint(ctx, gwmux, gwAddr, dopts)
	if err != nil {
		zap.L().Warn("Serve handler failed", zap.Error(err))
		return err
	}

	// start rest api
	go http.ListenAndServe(":9000", gwmux)

	return grpcServer.Serve(lis)
}

type kvServer struct {
	tree tree.Tree
}

func (s *kvServer) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	zap.L().Debug("Trigger Put method")
	s.tree.Put(string(req.GetKey()), string(req.GetValue()))
	s.tree.Walk()
	return &pb.PutResponse{Value: req.GetValue()}, nil
}

func (s *kvServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	zap.L().Debug("Trigger Get method")
	v, e := s.tree.Get(string(req.GetKey()))
	if e != nil {
		return &pb.GetResponse{}, e
	}
	return &pb.GetResponse{Value: []byte(v)}, nil
}

func (s *kvServer) Del(ctx context.Context, req *pb.DelRequest) (*pb.DelResponse, error) {
	zap.L().Debug("Trigger Del method")
	return &pb.DelResponse{}, errors.New("Key not found")
}

func newKVServer(cfg *config) *kvServer {
	s := kvServer{}

	switch cfg.backend {
	case "binary":
		s.tree = &tree.Binary{}
		zap.L().Info("Using binary search tree as backend")
	default:
		s.tree = &tree.Binary{}
		zap.L().Info("No backend specified, using binary search tree as default")
	}

	return &s
}
