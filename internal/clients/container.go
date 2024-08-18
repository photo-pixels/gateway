package clients

import (
	"fmt"
	userAccount "github.com/photo-pixels/gateway/pkg/gen/user_account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	UserAccountTarget string `yaml:"user_account_target"`
}

type ServiceClientsContainer struct {
	cfg        Config
	authClient userAccount.AuthServiceClient
	userClient userAccount.UserServiceClient
	connect    []*grpc.ClientConn
}

func NewServiceClientsContainer(cfg Config) (*ServiceClientsContainer, error) {
	s := ServiceClientsContainer{}

	conn, err := s.createConnect(cfg.UserAccountTarget)
	if err != nil {
		return nil, fmt.Errorf("s.createConnect: %w", err)
	}
	s.authClient = userAccount.NewAuthServiceClient(conn)

	conn, err = s.createConnect(cfg.UserAccountTarget)
	if err != nil {
		return nil, fmt.Errorf("s.createConnect: %w", err)
	}
	s.userClient = userAccount.NewUserServiceClient(conn)

	return &s, nil
}

func (s *ServiceClientsContainer) Close() error {
	for _, conn := range s.connect {
		err := conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ServiceClientsContainer) createConnect(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}
	s.connect = append(s.connect, conn)
	return conn, nil
}

func (s *ServiceClientsContainer) GetAuthClient() userAccount.AuthServiceClient {
	return s.authClient
}

func (s *ServiceClientsContainer) GetUserClient() userAccount.UserServiceClient {
	return s.userClient
}
