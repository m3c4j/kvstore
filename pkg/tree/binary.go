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

package tree

import (
	"fmt"

	"go.uber.org/zap"
)

type node struct {
	key   string
	value string

	left  *node
	right *node
	p     *node
}

// Binary implement Tree interface for BinarySearchTree
type Binary struct {
	root *node
}

// Put key-value into tree
// ref Introduction to Algorithms, 3rd Edition 12.3, P294
func (b *Binary) Put(key, value string) error {
	z := &node{
		key:   key,
		value: value,
	}

	var y *node
	x := b.root

	for x != nil {
		y = x
		if z.key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}

	z.p = y
	if y == nil {
		b.root = z
	} else if z.key < y.key {
		y.left = z
	} else {
		y.right = z
	}

	return nil
}

// GetRecursion implementation
func (b *Binary) GetRecursion(k string) (string, error) {
	n := b.getRecursion(b.root, k)
	if n == nil {
		return "", fmt.Errorf("Key %s not found", k)
	}
	return n.value, nil
}

func (b *Binary) getRecursion(r *node, k string) *node {
	if r == nil || k == r.key {
		return r
	}
	if k < r.key {
		return b.getRecursion(r.left, k)
	}
	return b.getRecursion(r.right, k)
}

// Get value by key from tree
func (b *Binary) Get(k string) (string, error) {
	r := b.root
	for r != nil && k != r.key {
		if k < r.key {
			r = r.left
		} else {
			r = r.right
		}
	}

	if r != nil {
		return r.value, nil
	}
	return "", fmt.Errorf("Key %s not found", k)
}

// Delete key-value from tree if exist
func (b *Binary) Delete(key string) error {
	return nil
}

// Walk all nodes in tree
func (b *Binary) Walk() error {
	zap.L().Info("Walking")
	b.InorderTreeWalk(b.root)
	return nil
}

// InorderTreeWalk ref L288
func (b *Binary) InorderTreeWalk(x *node) {
	if x != nil {
		b.InorderTreeWalk(x.left)
		zap.L().Info(x.key)
		b.InorderTreeWalk(x.right)
	}
}
