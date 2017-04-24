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
** Binary functions
/*/

// Count binary ones in uint32
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

// Check if the bit at given index is set in bitmap
func isBitSet(bitmap uint32, index uint) bool {
	return (bitmap & (1 << index)) > 0
}

// Set the bit at given index in bitmap
func setBit(bitmap uint32, index uint) uint32 {
	return (bitmap | (1 << index))
}

// Clear the bit at given index in bitmap
func clearBit(bitmap uint32, index uint) uint32 {
	return (bitmap &^ (1 << index))
}

/*/
** Hash functions
/*/

// Extract 5 bits of the hash. The part argument is the
// index of the 5-bit block to extract.
func hashPart(hash uint32, part uint) uint32 {
	shift := part * 5
	return (hash & (0x1F << shift) >> shift)
}

// Calculate the hash for given string
func hashString(s string, level int) uint32 {
	var vHash, a, b uint32
	vHash, a, b = 0, 31415, 27183
	for _, c := range s {
		vHash = a*vHash*uint32(level) + uint32(c)
		a *= b
	}
	return vHash
}

// Calculate the hash for a Value
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
		for {
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
				// TODO: rehash the key when part > 6 (otherwise it will overflow the 32 bit key)
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
	// TODO: use bitmaps for branch compression:
	// https://infoscience.epfl.ch/record/64398/files/idealhashtrees.pdf
	// https://idea.popcount.org/2012-07-25-introduction-to-hamt/
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

// TODO: Dissoc, Count?
