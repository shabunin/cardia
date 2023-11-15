package authentication

import (
	"context"
	"errors"
	"github.com/shabunin/cardia/proto"
)

type Server struct {
	svc Authenticator
	proto.UnimplementedAuthenticationServer
}

func (s *Server) PasswordAuth(ctx context.Context, req *proto.AuthPasswordReq) (*proto.AuthPasswordRes, error) {
	user := req.GetAccount()
	pass := req.GetPassword()
	token, err := s.svc.AuthenticateWithPassword(user, pass)
	if err != nil {
		return nil, err
	}

	res := &proto.AuthPasswordRes{Result: &proto.AuthSuccess{Token: token}}
	return res, nil
}
func (s *Server) PubkeyAuth(srv proto.Authentication_PubkeyAuthServer) error {
	req, err := srv.Recv()
	if err != nil {
		return err
	}
	user := req.

	return errors.New("implement me")
}
