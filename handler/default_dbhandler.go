package handler

import (
	"fmt"
	"strings"

	"github.com/anhdt-vnpay/f5_dynamic_gateway/domain/handler"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	SEPERATOR    = ",-,"
	SERVICES_KEY = "DataForServiceNames---"
)

type defaultDbHandler struct {
	db *leveldb.DB
}

func NewDefaultDbHandler(dbFile string) (handler.IDbHandler, error) {
	db, err := leveldb.OpenFile(dbFile, nil)
	return &defaultDbHandler{
		db: db,
	}, err
}

func (handler *defaultDbHandler) LoadServiceName() ([]string, error) {
	data, err := handler.db.Get([]byte(SERVICES_KEY), nil)
	if err != nil {
		fmt.Println("get saved data error ", err.Error())
		return []string{}, err
	}
	dataStr := string(data)
	fmt.Println("saved service name data ", dataStr)
	result := strings.Split(dataStr, SEPERATOR)
	if len(result) == 0 && len(dataStr) > 0 {
		result = append(result, dataStr)
	}
	return result, nil
}

func (handler *defaultDbHandler) Load(serviceName string) ([]string, error) {
	data, err := handler.db.Get([]byte(serviceName), nil)
	if err != nil {
		return []string{}, err
	}
	dataStr := string(data)
	fmt.Println("saved endpoint data ", dataStr)
	result := strings.Split(dataStr, SEPERATOR)
	if len(result) == 0 && len(dataStr) > 0 {
		result = append(result, dataStr)
	}
	return result, nil
}

func (handler *defaultDbHandler) Store(endpoint string, serviceName string) error {
	err := handler.updateServiceName(serviceName)
	if err != nil {
		fmt.Println("save service name err = ", err.Error())
		return err
	}
	endpoints, err := handler.Load(serviceName)
	if err != nil && err != leveldb.ErrNotFound {
		fmt.Println("load service name err = ", err.Error())
		return err
	}
	endpoints = append(endpoints, endpoint)
	return handler.saveData(endpoints, serviceName)
}

func (handler *defaultDbHandler) Delete(endpoint string, serviceName string) error {
	fmt.Println("delete endpoint = ", endpoint)
	err := handler.updateServiceName(serviceName)
	if err != nil {
		return err
	}
	endpoints, err := handler.Load(serviceName)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	itemIndex := -1
	for index, item := range endpoints {
		fmt.Println("item = ", item)
		if item == endpoint {
			itemIndex = index
			break
		}
	}
	if itemIndex >= 0 {
		fmt.Println("itemIndex = ", itemIndex)
		endpoints = append(endpoints[:itemIndex], endpoints[itemIndex+1:]...)
		return handler.saveData(endpoints, serviceName)
	}
	return nil
}

func (handler *defaultDbHandler) updateServiceName(serviceName string) error {
	serviceNames, err := handler.LoadServiceName()
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	if !containString(serviceNames, serviceName) {
		serviceNames = append(serviceNames, serviceName)
		return handler.saveData(serviceNames, SERVICES_KEY)
	}
	return nil
}

func (handler *defaultDbHandler) saveData(data []string, key string) error {
	savedData := strings.Join(data, SEPERATOR)
	fmt.Println("saved data ", savedData, " for key ", key)
	err := handler.db.Put([]byte(key), []byte(savedData), nil)
	if err != nil {
		fmt.Println("save data for key ", key, " err = ", err.Error())
	}
	return err
}

func containString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
