package balance

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
)

type HashConsistentBalance struct {
	table *crc32.Table
}

func init() {
	RegisterBalance("hashconsistent", &HashConsistentBalance{table: crc32.MakeTable(crc32.IEEE)})
}
func (r *HashConsistentBalance) Balance(ins []*Instance, keys ...string) (*Instance, error) {
	m := len(ins)
	if m == 0 {
		return nil, fmt.Errorf("no instance found")
	}
	hashKey := strconv.Itoa(rand.Int())
	if len(keys) > 0 {
		hashKey = keys[0]
	}
	hashcode := crc32.Checksum([]byte(hashKey), r.table)
	index := int(hashcode) % m
	return ins[index], nil
}
