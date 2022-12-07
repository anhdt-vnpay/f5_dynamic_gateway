package runtime

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type defaultConnectionService struct {
	dbHandler       IDbHandler
	endpointHandler IEndpointHandler
	config          ConnectionConfig
	poll            []*connectionInfo
	pollToRemove    []*EndpointInfo
}

func NewDefaultConnectionService(dbHanlder IDbHandler, endpointHandler IEndpointHandler) ConnectionService {
	connConfig := ConnectionConfig{
		MaxConn: 1,
	}
	service := defaultConnectionService{
		dbHandler:       dbHanlder,
		endpointHandler: endpointHandler,
		config:          connConfig,
		poll:            []*connectionInfo{},
		pollToRemove:    []*EndpointInfo{},
	}
	service.initConnections()
	go func() {
		service.manageRemoveConnection()
	}()
	return &service
}

func (s *defaultConnectionService) SetConfig(config ConnectionConfig, isForceApply bool) error {
	s.config = config
	return nil
}

func (s *defaultConnectionService) Add(serviceName string, endpoint string) error {
	endpointInfo, err := s.buildEndpointInfo(serviceName, endpoint)
	if err != nil {
		return err
	}

	s.saveDb(endpoint, serviceName)
	s.manageConnection(serviceName, endpointInfo)
	s.registerEndpoint(serviceName, endpointInfo)
	return nil
}

func newEndpointInfo(serviceName, endpoint string, conn *grpc.ClientConn) *EndpointInfo {
	return &EndpointInfo{
		Id:       "",
		Endpoint: endpoint,
		GrpcConn: conn,
	}
}

func (s *defaultConnectionService) buildEndpointInfo(serviceName string, endpoint string) (*EndpointInfo, error) {
	opts := []grpc.DialOption{
		grpc.WithAuthority(endpoint),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return newEndpointInfo(serviceName, endpoint, conn), nil
}

func (s *defaultConnectionService) manageConnection(serviceName string, endpointInfo *EndpointInfo) {
	isRegistedService := false
	for index := 0; index < len(s.poll); index++ {
		if s.poll[index].serviceName == serviceName {
			connInfo := s.poll[index]
			connInfo.endpoints = append(s.poll[index].endpoints, endpointInfo)
			s.poll[index].endpoints = s.manageLimitedConnnection(connInfo)
			isRegistedService = true
			break
		}
	}
	if !isRegistedService {
		connInfo := &connectionInfo{
			serviceName: serviceName,
			endpoints:   []*EndpointInfo{},
		}
		connInfo.endpoints = append(connInfo.endpoints, endpointInfo)
		s.poll = append(s.poll, connInfo)
	}
}

func (s *defaultConnectionService) manageLimitedConnnection(connInfo *connectionInfo) []*EndpointInfo {
	numberRemove := len(connInfo.endpoints) - s.config.MaxConn
	if numberRemove > 0 {
		for loop := 0; loop < numberRemove; loop++ {
			endpointInfo := connInfo.endpoints[loop]
			s.deleteDb(connInfo.serviceName, endpointInfo.Endpoint)
			s.unRegisterEndpoint(connInfo.serviceName, endpointInfo)
			s.pollToRemove = append(s.pollToRemove, endpointInfo)
		}
	}
	return connInfo.endpoints[numberRemove:]
}

func (s *defaultConnectionService) manageRemoveConnection() {
	for {
		if len(s.pollToRemove) <= 0 {
			time.Sleep(60 * time.Second)
			continue
		}
		processSlice := []*EndpointInfo{}
		for {
			if len(s.pollToRemove) <= 0 {
				break
			}
			conn := s.pollToRemove[0]
			s.pollToRemove = s.pollToRemove[1:]
			state := conn.GrpcConn.GetState()
			if state == connectivity.Idle {
				conn.GrpcConn.Close()
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

func (s *defaultConnectionService) deleteDb(serviceName, endpoint string) error {
	if s.dbHandler != nil {
		return s.dbHandler.Delete(endpoint, serviceName)
	}
	return nil
}

func (s *defaultConnectionService) loadAll() ([]*connectionInfo, error) {
	connectionInfos := []*connectionInfo{}
	if s.dbHandler != nil {
		serviceNames, err := s.dbHandler.LoadServiceName()
		if err != nil {
			return connectionInfos, err
		}
		for _, serviceName := range serviceNames {
			item, err := s.initConnectionInfos(serviceName)
			if err == nil {
				connectionInfos = append(connectionInfos, item)
			}
		}
	}
	return connectionInfos, nil
}

func (s *defaultConnectionService) initConnectionInfos(serviceName string) (*connectionInfo, error) {
	endpoints, err := s.dbHandler.Load(serviceName)
	if err != nil {
		return nil, err
	}
	endpointInfos := []*EndpointInfo{}
	for _, endpoint := range endpoints {
		item, err := s.buildEndpointInfo(serviceName, endpoint)
		if err == nil {
			endpointInfos = append(endpointInfos, item)
		}
	}
	return &connectionInfo{
		serviceName: serviceName,
		endpoints:   endpointInfos,
	}, nil
}

func (s *defaultConnectionService) initConnections() {
	connectionInfos, err := s.loadAll()
	if err != nil || len(connectionInfos) <= 0 {
		return
	}

	for _, connInfo := range connectionInfos {
		for _, endpoint := range connInfo.endpoints {
			s.registerEndpoint(connInfo.serviceName, endpoint)
		}
	}
}

func (s *defaultConnectionService) registerEndpoint(serviceName string, endpointInfo *EndpointInfo) error {
	if s.endpointHandler != nil {
		return s.endpointHandler.RegisterEndpoint(serviceName, endpointInfo)
	}
	return nil
}

func (s *defaultConnectionService) unRegisterEndpoint(serviceName string, endpointInfo *EndpointInfo) error {
	if s.endpointHandler != nil {
		return s.endpointHandler.UnRegisterEndpoint(serviceName, endpointInfo)
	}
	return nil
}
