package balance

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"
)

var instance []*Instance

func TestMain(m *testing.M) {

	for i := 0; i < 10; i++ {
		host := fmt.Sprintf("192.168.%s.%s", strconv.Itoa(rand.Intn(255)), strconv.Itoa(rand.Intn(rand.Intn(255))))
		port := rand.Intn(65536)
		weight := rand.Intn(10)
		instance = append(instance, NewInstance(host, port, weight))
	}
	os.Exit(m.Run())
}

func TestRoundRobin(t *testing.T) {
	size := len(instance)
	inst, err := StartBalance("roundrobin", instance)
	if err != nil {
		t.Error("do balance error", err)
		return
	}
	for i := 0; i < size-1; i++ {
		StartBalance("roundrobin", instance)
	}
	result, err := StartBalance("roundrobin", instance)
	if err != nil {
		t.Error("do balance error next stage", err)
		return
	}
	if !reflect.DeepEqual(inst, result) {
		t.Errorf("expected %v, result is %v", inst, result)
	}
}

func TestRandom(t *testing.T) {
	for i := 0; i < len(instance); i++ {
		balance, err := StartBalance("random", instance)
		if err != nil {
			t.Log("do balance error", err)
		}
		d, _ := json.Marshal(balance)
		t.Log("select instance " + string(d))
	}
}

func TestShuffle2(t *testing.T) {
	for i := 0; i < len(instance); i++ {
		balance, err := StartBalance("shuffle2", instance)
		if err != nil {
			t.Log("do balance error", err)
		}
		d, _ := json.Marshal(balance)
		t.Log("select instance " + string(d))
	}
}

func TestShuffle(t *testing.T) {
	for i := 0; i < len(instance); i++ {
		balance, err := StartBalance("shuffle", instance)
		if err != nil {
			t.Log("do balance error", err)
		}
		d, _ := json.Marshal(balance)
		t.Log("select instance " + string(d))
	}
}
func TestHashconsistent(t *testing.T) {
	for i := 0; i < len(instance); i++ {
		balance, err := StartBalance("hashconsistent", instance)
		if err != nil {
			t.Log("do balance error", err)
		}
		d, _ := json.Marshal(balance)
		t.Log("select instance " + string(d))
	}
}

func TestWrr(t *testing.T) {
	counterMap := make(map[string]int)
	instantceMap := make(map[string]*Instance)

	weightSum := 0
	for _, instance := range instance {
		weightSum += instance.Weight
		d, _ := json.Marshal(instance)
		instantceMap[string(d)] = instance
	}

	for i := 0; i < weightSum; i++ {
		inst, err := StartBalance("roundrobinweight", instance)
		if err != nil {
			t.Error(err)
			continue
		}
		d, _ := json.Marshal(inst)
		if _, ok := counterMap[string(d)]; !ok {
			counterMap[string(d)] = 0
		}
		t.Logf("select instance %v\n", inst)
		counterMap[string(d)]++
	}

	// equal
	for key, count := range counterMap {
		in := instantceMap[key]
		d, _ := json.Marshal(in)
		if in.Weight != count {
			t.Errorf("Sum : %d, Count: %d, Instance: %s", weightSum, count, string(d))
		}
	}
}
