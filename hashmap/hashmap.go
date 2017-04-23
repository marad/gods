package hashmap

/**
 *   Copyright (c) Marcin Radoszewski. All rights reserved.
 *   The use and distribution terms for this software are covered by the
 *   Eclipse Public License 1.0 (http://opensource.org/licenses/eclipse-1.0.php)
 *   which can be found in the file epl-v10.html at the root of this distribution.
 *   By using this software in any fashion, you are agreeing to be bound by
 * 	 the terms of this license.
 *   You must not remove this notice, or any other, from this software.
 **/

import (
	"errors"
	"fmt"
	. "gods"
)

type Node interface {
	Assoc(key uint32, val Value, part uint) Node
	Find(key uint32, part uint) Node
}

type ValueNode struct {
	Key       uint32
	BaseValue Value
}

type SubtreeNode struct {
	BitMapKey uint32
	Branches  [32]Node
}

type HashMap struct {
	root     SubtreeNode
	hashFunc func(Value) (uint32, error)
}

/*/
** Hash helper functions
/*/

func popcnt(x uint32) (n byte) {
	// bit population count, see
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	x -= (x >> 1) & 0x55555555
	x = (x>>2)&0x33333333 + x&0x33333333
	x += x >> 4
	x &= 0x0f0f0f0f
	x *= 0x01010101
	return byte(x >> 24)
}

func hashPart(hash uint32, part uint) uint32 {
	return (hash & (0x1F << (part * 5)) >> (part * 5))
}

func hashString(s string, level int) uint32 {
	var vHash, a, b uint32
	vHash, a, b = 0, 31415, 27183
	for _, c := range s {
		vHash = a*vHash*uint32(level) + uint32(c)
		a *= b
	}
	return vHash
}

func Hash(value Value) (uint32, error) {
	switch value.(type) {
	case string:
		return hashString(value.(string), 1), nil
	}
	return 0, errors.New(fmt.Sprintf("Don't know how to hash %v", value))
}

/*/
** Value Node Implementation
/*/

func (vn ValueNode) Assoc(key uint32, val Value, part uint) Node {
	resultNode := SubtreeNode{}
	var currentNode *SubtreeNode
	currentNode = &resultNode
	currentPart := part

	if key == vn.Key {
		return ValueNode{key, val}
	} else {
		for i := 0; i < 7; i++ {
			localPart := hashPart(vn.Key, currentPart)
			newPart := hashPart(key, currentPart)

			if localPart != newPart {
				currentNode.Branches[localPart] = &ValueNode{vn.Key, vn.BaseValue}
				currentNode.Branches[newPart] = &ValueNode{key, val}
				break
			} else {
				newNode := SubtreeNode{}
				currentNode.Branches[localPart] = &newNode
				currentNode = &newNode
				currentPart += 1
			}
		}
	}
	return resultNode
}

func (vn ValueNode) Find(key uint32, part uint) Node {
	if key == vn.Key {
		return vn
	} else {
		return nil
	}
}

/*/
** Subtree Node Implementation
/*/

func (sn SubtreeNode) Assoc(key uint32, val Value, part uint) Node {
	index := hashPart(key, part)
	node := sn.Branches[index]
	newNode := SubtreeNode{sn.BitMapKey, sn.Branches}
	if node == nil {
		newNode.Branches[index] = ValueNode{key, val}
		return newNode
	} else {
		newNode.Branches[index] = node.Assoc(key, val, part+1)
		return newNode
	}
}

func (sn SubtreeNode) Find(key uint32, part uint) Node {
	index := hashPart(key, part)
	node := sn.Branches[index]
	if node != nil {
		return node.Find(key, part+1)
	} else {
		return nil
	}
}

/*/
** Hash Map Trie
/*/

func New() *HashMap {
	hm := HashMap{}
	hm.hashFunc = Hash
	return &hm
}

func (hm *HashMap) Assoc(key Value, val Value) *HashMap {
	hash, _ := hm.hashFunc(key)
	newRoot, _ := hm.root.Assoc(hash, val, 0).(SubtreeNode)
	return &HashMap{newRoot, hm.hashFunc}
}

func (hm *HashMap) Find(key Value) Value {
	hash, _ := hm.hashFunc(key)
	result := hm.root.Find(hash, 0)
	if vn, ok := result.(ValueNode); ok {
		return vn.BaseValue
	} else {
		return nil
	}
}
