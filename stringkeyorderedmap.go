package generics

import (
	"bytes"
	"encoding/json"
)

type KV[V any] struct {
	Key   string
	Value V
}

type StringKeyOrderedMap[V any] struct {
	kvList    []*KV[V]
	idxLookup map[string]int
}

func NewStringKeyOrderedMap[V any](kvList ...*KV[V]) *StringKeyOrderedMap[V] {
	result := &StringKeyOrderedMap[V]{
		idxLookup: make(map[string]int),
	}

	for i := 0; i < len(kvList); i++ {
		result.Set(kvList[i].Key, kvList[i].Value)
	}

	return result
}

func (m *StringKeyOrderedMap[V]) Set(key string, value V) *StringKeyOrderedMap[V] {
	idx, ok := m.idxLookup[key]
	if !ok {
		m.idxLookup[key] = len(m.kvList)
		m.kvList = append(m.kvList, &KV[V]{key, value})
		return m
	}

	m.kvList[idx].Value = value
	return m
}

func (m *StringKeyOrderedMap[V]) Get(key string) V {
	if idx, ok := m.idxLookup[key]; ok {
		return m.kvList[idx].Value
	}
	var v V
	return v
}

func (m *StringKeyOrderedMap[V]) Exists(key string) bool {
	_, ok := m.idxLookup[key]
	return ok
}

func (m *StringKeyOrderedMap[V]) Delete(key string) {
	if idx, ok := m.idxLookup[key]; ok {
		delete(m.idxLookup, key)
		m.kvList[idx] = nil
	}
}

func (m *StringKeyOrderedMap[V]) GetKeys() []string {
	keys := make([]string, 0, len(m.idxLookup))
	for idx := 0; idx < len(m.kvList); idx++ {
		if m.kvList[idx] == nil {
			continue
		}

		keys = append(keys, m.kvList[idx].Key)
	}
	return keys
}

func (m *StringKeyOrderedMap[V]) GetList() []KV[V] {
	kvList := make([]KV[V], 0, len(m.idxLookup))
	for idx := 0; idx < len(m.kvList); idx++ {
		if m.kvList[idx] == nil {
			continue
		}

		kvList = append(kvList, *m.kvList[idx])
	}
	return kvList
}

func (m *StringKeyOrderedMap[V]) Append(newOm *StringKeyOrderedMap[V], overwrite bool) *StringKeyOrderedMap[V] {
	for _, kv := range newOm.GetList() {
		if !overwrite && m.Exists(kv.Key) {
			continue
		}

		m.Set(kv.Key, kv.Value)
	}
	return m
}

func (m *StringKeyOrderedMap[V]) Len() int {
	return len(m.idxLookup)
}

func (m *StringKeyOrderedMap[V]) String() string {
	data, _ := m.MarshalJSON()
	return string(data)
}

func (m *StringKeyOrderedMap[V]) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	first := true
	for idx := 0; idx < len(m.kvList); idx++ {
		if m.kvList[idx] == nil {
			continue
		}

		if !first {
			buffer.WriteRune(',')
		}

		key, err := json.Marshal(m.kvList[idx].Key)
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(m.kvList[idx].Value)
		if err != nil {
			return nil, err
		}

		buffer.Write(key)
		buffer.WriteByte(58)
		buffer.Write(value)

		first = false
	}

	buffer.WriteRune('}')
	return buffer.Bytes(), nil
}

func (kv KV[V]) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	key, err := json.Marshal(kv.Key)
	if err != nil {
		return nil, err
	}
	value, err := json.Marshal(kv.Value)
	if err != nil {
		return nil, err
	}

	buffer.Write(key)
	buffer.WriteByte(58)
	buffer.Write(value)

	buffer.WriteRune('}')
	return buffer.Bytes(), nil
}
