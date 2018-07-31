package dto

import (
	"encoding/json"
	"fmt"
)

type StringList struct {
	value []string
	key   []byte
}

func NewStringListCache(key string) StringList {
	return StringList{
		key:   []byte(key),
		value: nil,
	}
}
func (sl *StringList) CacheKey() []byte {
	return sl.key
}

func (sl *StringList) SetValue(value []string) {
	sl.value = value
}

func (sl *StringList) GetValue() []string {
	return sl.value
}

func (sl *StringList) JSON() (data []byte, err error) {
	data, err = json.Marshal(sl.value)
	if err != nil {
		err = fmt.Errorf("StringList:Json:%s,%v", string(sl.key), err)
	}
	return
}

func (sl *StringList) LoadJSON(data []byte) error {
	var value []string
	err := json.Unmarshal(data, &value)
	if err == nil {
		sl.value = value
	}
	return err
}
func (sl *StringList) CacheExpireTime() int {
	return 3600 * 24
}
