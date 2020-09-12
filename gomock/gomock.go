package gomock

import (
	"errors"
	"github.com/hq-cml/go-unittest/common"
)

//client句柄从外部传入
//典型的依赖注入场景
func CheckItemKey1(client common.StorageClient, key string) (bool, error) {
	v, ok := client.Get(key)
	if !ok {
		return false, errors.New("Shoud Exist!")
	}

	if v == "Hello world" {
		return true, nil
	} else {
		return false, errors.New("Value Wrong!")
	}
}

//client内部生成
//需要内部打桩
func CheckItemKey2(key string) (bool, error) {
	client := common.NewStorageClient()
	v, ok := client.Get(key)
	if !ok {
		return false, errors.New("Shoud Exist!")
	}

	if v == "Hello world" {
		return true, nil
	} else {
		return false, errors.New("Value Wrong!")
	}
}

//如果存在则返回，否则用默认值设置
//测试行为的保序
func Replace(client common.StorageClient, key,def string) (string, error) {
	v, ok := client.Get(key)
	if ok {
		return v, nil
	}

	err := client.Set(key, def)
	if err != nil {
		return "", err
	}

	v, ok = client.Get(key)
	if !ok {
		return "", errors.New("Shoud Exist")
	}

	if v != def {
		return "", errors.New("Shoud =" + def)
	} else {
		return def, nil
	}
}