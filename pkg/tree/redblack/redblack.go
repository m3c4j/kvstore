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

package redblack

import (
	"fmt"

	"go.uber.org/zap"
)

type color bool

const (
	red   color = false
	black color = true
)

type node struct {
	key   string
	value string

	left  *node
	right *node
	p     *node

	color color
}

// RedBlack implement Tree interface for RedBlackTree
type RedBlack struct {
	root     *node
	sentinel *node
}

// Put key-value into tree
func (r *RedBlack) Put(k, v string) error {
	z := &node{
		key:   k,
		value: v,

		left:  r.sentinel,
		right: r.sentinel,
		p:     r.sentinel,

		color: red,
	}

	y := r.sentinel
	x := r.root

	for x != r.sentinel {
		y = x
		if z.key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}

	z.p = y
	if y == r.sentinel {
		r.root = z
	} else if z.key < y.key {
		y.left = z
	} else {
		y.right = z
	}

	z.left = r.sentinel
	z.right = r.sentinel
	z.color = red

	r.insertFixup(z)

	return nil
}

func (r *RedBlack) insertFixup(z *node) {
	for z.p.color == red {
		if z.p == z.p.p.left {
			y := z.p.p.right
			if y.color == red {
				z.p.color = black
				y.color = black
				z.p.p.color = red
				z = z.p.p
			} else {
				if z == z.p.right {
					z = z.p
					r.leftRotate(z)
				}
				z.p.color = black
				z.p.p.color = red
				r.rightRotate(z.p.p)
			}
		} else {
			y := z.p.p.left
			if y.color == red {
				z.p.color = black
				y.color = black
				z.p.p.color = red
				z = z.p.p
			} else {
				if z == z.p.left {
					z = z.p
					r.rightRotate(z)
				}
				z.p.color = black
				z.p.p.color = red
				r.leftRotate(z.p.p)
			}
		}
	}
	r.root.color = black
}

// Get value by key from tree
func (r *RedBlack) Get(k string) (string, error) {
	t := r.root
	for t != r.sentinel && k != t.key {
		if k < t.key {
			t = t.left
		} else {
			t = t.right
		}
	}

	if t != r.sentinel {
		return t.value, nil
	}

	return "", fmt.Errorf("Key %s not found", k)
}

// Del key-value from tree if exist
func (r *RedBlack) Del(k string) error {
	return nil
}

// Walk all nodes in tree
func (r *RedBlack) Walk() error {
	// zap.L().Info("Walking")
	r.inorderTreeWalk(r.root)
	return nil
}

func (r *RedBlack) inorderTreeWalk(x *node) {
	if x != nil {
		r.inorderTreeWalk(x.left)
		zap.L().Info(x.key)
		r.inorderTreeWalk(x.right)
	}
}

func (r *RedBlack) leftRotate(x *node) {
	y := x.right
	x.right = y.left
	if y.left != r.sentinel {
		y.left.p = x
	}
	y.p = x.p
	if x.p == r.sentinel {
		r.root = y
	} else if x == x.p.left {
		x.p.left = y
	} else {
		x.p.right = y
	}
	y.left = x
	x.p = y
}

func (r *RedBlack) rightRotate(y *node) {
	x := y.left
	y.left = x.right
	if x.right != r.sentinel {
		x.right.p = y
	}
	x.p = y.p
	if y.p == r.sentinel {
		r.root = x
	} else if y == y.p.left {
		y.p.left = x
	} else {
		y.p.right = x
	}
	x.right = y
	y.p = x
}

func makeDummy() *node {
	return &node{
		key:   "",
		value: "",

		left:  nil,
		right: nil,
		p:     nil,

		color: black,
	}
}

// New create a new RedBlack tree
func New() *RedBlack {
	sentinel := makeDummy()
	return &RedBlack{
		root:     sentinel,
		sentinel: sentinel,
	}
}
