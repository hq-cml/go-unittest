package common

import (
	"errors"
)

//接口定义，模拟某种存储的操作，稍后通过gomock进行替身变换
type StorageClient interface {
	Set(k, v string) error
	Get(k string) (string, bool)
}

var DataMap map[string]string

func init () {
	DataMap = map[string]string{}
}

//具体实现
type RealClient struct {
}

func NewStorageClient() StorageClient {
	return &RealClient{}
}

func (m *RealClient)Get(k string) (string, bool) {
	if DataMap == nil {
		return "", false
	}
	if v, ok := DataMap[k]; ok {
		return v, true
	} else {
		return "", false
	}
}

func (m *RealClient)Set(k, v string) error {
	if DataMap == nil {
		return errors.New("DataMap is nil")
	}
	DataMap[k] = v
	return nil
}