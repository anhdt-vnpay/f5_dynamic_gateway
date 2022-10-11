package services

import (
	"fmt"
	"time"

	"github.com/anhdt-vnpay/f5_dynamic_gateway/v1/domain/handler"
	"github.com/anhdt-vnpay/f5_dynamic_gateway/v1/domain/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type connectionInfo struct {
	serviceName string
	endpoints   []string
	connPool    []*grpc.ClientConn
}

type defaultConnectionService struct {
	dbHandler       handler.IDbHandler
	endpointHandler handler.IEndpointHandler
	config          services.ConnectionConfig
	poll            []*connectionInfo
	pollToRemove    []*grpc.ClientConn
}

func NewDefaultConnectionService(dbHanlder handler.IDbHandler, endpointHandler handler.IEndpointHandler) services.ConnectionService {
	connConfig := services.ConnectionConfig{
		MaxConn: 1,
	}
	service := defaultConnectionService{
		dbHandler:       dbHanlder,
		endpointHandler: endpointHandler,
		config:          connConfig,
		poll:            []*connectionInfo{},
		pollToRemove:    []*grpc.ClientConn{},
	}
	service.initConnections()
	go func() {
		service.manageRemoveConnection()
	}()
	return &service

}

func (s *defaultConnectionService) SetConfig(config services.ConnectionConfig, isForceApply bool) error {
	s.config = config
	return nil
}

func (s *defaultConnectionService) Add(endpoint string, serviceName string) error {
	conn, err := s.buildConnection(endpoint)
	if err != nil {
		return err
	}

	s.saveDb(endpoint, serviceName)
	s.manageConnection(conn, endpoint, serviceName)
	s.registerEndpoint(conn, endpoint, serviceName)
	return nil
}

func (s *defaultConnectionService) buildConnection(endpoint string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(), // Block when calling Dial until the connection is really established
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return grpc.Dial(endpoint, opts...)
}

func (s *defaultConnectionService) manageConnection(conn *grpc.ClientConn, endpoint string, serviceName string) {
	isRegistedService := false
	for index := 0; index < len(s.poll); index++ {
		if s.poll[index].serviceName == serviceName {
			connInfo := s.poll[index]
			connInfo.endpoints = append(s.poll[index].endpoints, endpoint)
			connInfo.connPool = append(s.poll[index].connPool, conn)
			s.poll[index].connPool, s.poll[index].endpoints = s.manageLimitedConnnection(*connInfo)
			isRegistedService = true
			break
		}
	}
	if !isRegistedService {
		connInfo := &connectionInfo{
			serviceName: serviceName,
			connPool:    []*grpc.ClientConn{},
			endpoints:   []string{},
		}
		connInfo.connPool = append(connInfo.connPool, conn)
		connInfo.endpoints = append(connInfo.endpoints, endpoint)
		s.poll = append(s.poll, connInfo)
	}
}

func (s *defaultConnectionService) manageLimitedConnnection(connInfo connectionInfo) ([]*grpc.ClientConn, []string) {
	numberRemove := len(connInfo.connPool) - s.config.MaxConn
	if numberRemove > 0 {
		fmt.Println("number conn over = ", numberRemove)
		for loop := 0; loop < numberRemove; loop++ {
			endpoint := connInfo.endpoints[loop]
			conn := connInfo.connPool[loop]
			s.deleteDb(endpoint, connInfo.serviceName)
			s.unRegisterEndpoint(conn, endpoint, connInfo.serviceName)
			s.pollToRemove = append(s.pollToRemove, conn)
		}
	}
	return connInfo.connPool[numberRemove:], connInfo.endpoints[numberRemove:]
}

func (s *defaultConnectionService) manageRemoveConnection() {
	for {
		if len(s.pollToRemove) <= 0 {
			fmt.Println("No conn to close -> sleep 5 second")
			time.Sleep(60 * time.Second)
			continue
		}
		processSlice := []*grpc.ClientConn{}
		for {
			if len(s.pollToRemove) <= 0 {
				break
			}
			conn := s.pollToRemove[0]
			s.pollToRemove = s.pollToRemove[1:]
			state := conn.GetState()
			if state == connectivity.Idle {
				conn.Close()
				continue
			}
			processSlice = append(processSlice, conn)
		}
		s.pollToRemove = append(s.pollToRemove, processSlice...)
		time.Sleep(60 * time.Second)
	}
}

func (s *defaultConnectionService) saveDb(endpoint string, serviceName string) error {
	if s.dbHandler != nil {
		return s.dbHandler.Store(endpoint, serviceName)
	}
	return nil
}

func (s *defaultConnectionService) deleteDb(endpoint string, serviceName string) error {
	if s.dbHandler != nil {
		return s.dbHandler.Delete(endpoint, serviceName)
	}
	return nil
}

func (s *defaultConnectionService) loadAll() ([]connectionInfo, error) {
	connectionInfos := []connectionInfo{}
	if s.dbHandler != nil {
		serviceNames, err := s.dbHandler.LoadServiceName()
		if err != nil {
			fmt.Println("load service name err " + err.Error())
			return connectionInfos, err
		}
		fmt.Println("load service name count ", len(serviceNames))
		for index := 0; index < len(serviceNames); index++ {
			serviceName := serviceNames[index]
			endpoinds, err := s.dbHandler.Load(serviceName)
			fmt.Println("endpoint count = ", len(endpoinds))
			if err != nil {
				continue
			}
			connectionInfo := connectionInfo{
				serviceName: serviceName,
				endpoints:   endpoinds,
				connPool:    []*grpc.ClientConn{},
			}
			connectionInfos = append(connectionInfos, connectionInfo)
		}
	}
	return connectionInfos, nil
}

func (s *defaultConnectionService) initConnections() {
	connectionInfos, err := s.loadAll()
	if err != nil || len(connectionInfos) <= 0 {
		return
	}

	for _, connInfo := range connectionInfos {
		for _, endpoint := range connInfo.endpoints {
			conn, err := s.buildConnection(endpoint)
			if err == nil {
				connInfo.connPool = append(connInfo.connPool, conn)
				s.endpointHandler.RegisterEndpoint(conn, endpoint, connInfo.serviceName)
			}
		}
		if len(connInfo.connPool) > 0 {
			s.poll = append(s.poll, &connInfo)
		}
	}
}

func (s *defaultConnectionService) registerEndpoint(conn *grpc.ClientConn, endpoint string, serviceName string) error {
	if s.endpointHandler != nil {
		return s.endpointHandler.RegisterEndpoint(conn, endpoint, serviceName)
	}
	return nil
}

func (s *defaultConnectionService) unRegisterEndpoint(conn *grpc.ClientConn, endpoint string, serviceName string) error {
	if s.endpointHandler != nil {
		return s.endpointHandler.UnRegisterEndpoint(conn, endpoint, serviceName)
	}
	return nil
}
