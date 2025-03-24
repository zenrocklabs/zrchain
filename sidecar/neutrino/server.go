package neutrino

//
//import (
//	"context"
//	"fmt"
//	"github.com/lightninglabs/neutrino"
//	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
//	"google.golang.org/grpc"
//	"log"
//	"net"
//)
//
//func StartGRPCServer(node *neutrino.ChainService, port int) {
//	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//
//	s := grpc.NewServer()
//	api.RegisterSidecarServiceServer(s, &neutrinoService{node: node})
//
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}
//
//type neutrinoService struct {
//	api.UnimplementedSidecarServiceServer
//	node *neutrino.ChainService
//}
//
//func (n *neutrinoService) GetSidecarState(ctx context.Context, req *api.SidecarStateRequest) (*api.SidecarStateResponse, error) {
//	return nil, nil
//}
//
//func (n *neutrinoService) GetSidecarStateByID(ctx context.Context, req *api.SidecarStateByIDRequest) (*api.SidecarStateResponse, error) {
//	return nil, nil
//}
//
////func (n *neutrinoService) GetBestBlock(ctx context.Context, req *api.BlockRequest) (*api.SidecarStateResponse, error) {
////	return nil, nil
////}
