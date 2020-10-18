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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"google.golang.org/grpc"

	pb "github.com/kvstore/api/kvserverpb"
)

var (
	version string
	url     string

	command string
	key     string
	value   string
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No command specified")
		os.Exit(0)
	}

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.StringVar(&version, "version", "", "Print version code")
	flagSet.StringVar(&url, "url", "127.0.0.1:8000", "Server address and port")
	flagSet.Parse(os.Args)

	var err error

	switch os.Args[1] {
	case "put":
		if len(os.Args) < 4 {
			fmt.Printf("No enough values\n\te.g. %s put a b\n", os.Args[0])
			os.Exit(0)
		}

		if err = put(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("Put failed", err.Error())
		}

	case "get":
		if len(os.Args) < 3 {
			fmt.Printf("No enough values\n\te.g. %s get a\n", os.Args[0])
			os.Exit(0)
		}

		if err = get(os.Args[2]); err != nil {
			fmt.Println("get failed", err.Error())
		}

	default:
		fmt.Println("The command is not vaild, e.g. put get del")
		os.Exit(0)
	}

}

func put(key, value string) error {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dail grpc failed", url)
		return err
	}
	defer conn.Close()

	client := pb.NewKVClient(conn)

	ctx := context.Background()

	_, err = client.Put(ctx, &pb.PutRequest{
		Key:   []byte(key),
		Value: []byte(value),
	})
	if err != nil {
		fmt.Println("Put request failed")
		return err
	}

	// fmt.Println(string(putRep.GetValue()))
	fmt.Println("ok")

	return nil
}

func get(key string) error {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dail grpc failed", url)
		return err
	}
	defer conn.Close()

	client := pb.NewKVClient(conn)

	ctx := context.Background()

	getRep, err := client.Get(ctx, &pb.GetRequest{
		Key: []byte(key),
	})

	if err != nil {
		fmt.Println("Get request failed")
		return err
	}

	fmt.Println(string(getRep.GetValue()))

	return nil
}
