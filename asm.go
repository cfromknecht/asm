package asm

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
)

type Direction int

const (
	LEFT Direction = iota
	RIGHT
)

type WitnessPath []WitnessNode

type WitnessNode struct {
	hash string
	dir  Direction
}

const (
	Bot = "-"
)

type AsyncAcc []string

func NewAsyncAcc() AsyncAcc {
	return AsyncAcc{Bot}
}

func (acc *AsyncAcc) Add(x string) (witPath WitnessPath) {
	newAcc := *acc

	d := 0
	z := base64SHA256(x)

	for newAcc[d] != Bot {
		if len(newAcc) < d+2 {
			newAcc = append(newAcc, Bot)
		}

		z = base64SHA256(newAcc[d] + z)
		witPath = append(witPath, WitnessNode{newAcc[d], LEFT})
		newAcc[d] = Bot

		d++
	}
	newAcc[d] = z

	*acc = newAcc
	return
}

func (acc AsyncAcc) Verify(x string, witPath WitnessPath) bool {
	for _, a := range getAncestors(x, witPath) {
		for _, root := range acc {
			if root == a {
				return true
			}
		}
	}
	return false
}

func UpdateWitness(y string, witPathY, witPathX WitnessPath) (newWitPathX WitnessPath) {
	dx := len(witPathX)
	dy := len(witPathY)
	// No updates to witness
	if dy < dx {
		return witPathX
	}

	ancestorsY := getAncestors(y, witPathY)
	newWitPathX = witPathX
	newWitPathX = append(newWitPathX, WitnessNode{ancestorsY[dx], RIGHT})
	if dx+1 < len(witPathY) {
		newWitPathX = append(newWitPathX, witPathY[dx+1:]...)
	}

	return
}

func getAncestors(x string, witPath WitnessPath) []string {
	c := base64SHA256(x)
	ancestors := []string{c}

	for _, node := range witPath {
		if node.dir == LEFT {
			c = base64SHA256(node.hash + c)
		} else {
			c = base64SHA256(c + node.hash)
		}
		ancestors = append(ancestors, c)
	}

	return ancestors
}

func base64SHA256(x string) string {
	h256 := sha256.New()
	io.WriteString(h256, x)
	return base64.StdEncoding.EncodeToString(h256.Sum(nil))
}
