package authentication

import (
	"context"
	"github.com/shabunin/cardia/proto"
)

type Server struct {
	svc Authenticator
	proto.UnimplementedAuthenticationServer
}

// TODO grpc errors with status code

func (s *Server) PasswordAuth(ctx context.Context, req *proto.AuthPasswordReq) (*proto.AuthPasswordRes, error) {
	user := req.GetAccount()
	pass := req.GetPassword()
	token, err := s.svc.AuthenticateWithPassword(user, pass)
	if err != nil {
		return nil, err
	}

	res := &proto.AuthPasswordRes{
		Result: &proto.AuthSuccess{
			Token: token}}
	return res, nil
}

func (s *Server) PubkeyAuth(srv proto.Authentication_PubkeyAuthServer) error {
	req, err := srv.Recv()
	if err != nil {
		return err
	}
	user := req.GetAccount()
	algo := req.GetPubkeyAlgorithm()
	_ = algo // TODO use all params
	pubk := req.GetPubkeyBlob()

	token, err := s.svc.AuthenticateWithPubkey(user, pubk,
		func(request []byte) []byte {
			err := srv.Send(&proto.AuthPubkeyRes{
				Payload: &proto.AuthPubkeyRes_SignRequest{
					SignRequest: request,
				},
			})
			if err != nil {
				return nil
			}

			sig, err := srv.Recv()
			if err != nil {
				return nil
			}

			return sig.GetSignature()
		})
	if err != nil {
		return err
	}

	return srv.Send(
		&proto.AuthPubkeyRes{
			Payload: &proto.AuthPubkeyRes_Result{
				Result: &proto.AuthSuccess{
					Token: token}}}) // =} =)
}
