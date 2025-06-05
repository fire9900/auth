package client

import (
	"context"
	"fmt"
	auth "github.com/fire9900/auth/pkg/api/g_rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn    *grpc.ClientConn
	service auth.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		conn:    conn,
		service: auth.NewAuthServiceClient(conn),
	}, nil
}

func (c *AuthClient) Close() {
	c.conn.Close()
}

func (c *AuthClient) ValidateToken(token string) (bool, error) {
	fmt.Println(token)
	resp, err := c.service.ValidateToken(context.Background(), &auth.TokenRequest{
		Token: token,
	})

	if err != nil {
		return false, err
	}

	return resp.Valid, nil
}

func (c *AuthClient) GetUserID(token string) (int32, error) {
	resp, err := c.service.GetUserID(context.Background(), &auth.TokenRequest{
		Token: token,
	})

	if err != nil {
		return 0, err
	}

	return resp.UserId, nil
}
