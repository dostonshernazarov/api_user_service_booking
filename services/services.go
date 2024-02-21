package services

import (
	"fmt"

	"api_user_service_booking/config"
	pbu "api_user_service_booking/genproto/user_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() pbu.UserServiceClient
	//PostService() pbp.PostServiceClient
	//CommentService() pbc.CommentServiceClient
	//LikeService() pbl.LikeServiceClient
}

type serviceManager struct {
	userService pbu.UserServiceClient
	//postService    pbp.PostServiceClient
	//commentService pbc.CommentServiceClient
	//likeService    pbl.LikeServiceClient
}

func (s *serviceManager) UserService() pbu.UserServiceClient {
	return s.userService
}

//	func (s *serviceManager) PostService() pbp.PostServiceClient {
//		return s.postService
//	}
//
//	func (s *serviceManager) CommentService() pbc.CommentServiceClient {
//		return s.commentService
//	}
//
//	func (s *serviceManager) LikeService() pbl.LikeServiceClient {
//		return s.likeService
//	}
func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.UserServiceHost, conf.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	//connPost, err := grpc.Dial(
	//	fmt.Sprintf("%s:%d", conf.PostServiceHost, conf.PostServicePort),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	return nil, err
	//}
	//
	//connComment, err := grpc.Dial(
	//	fmt.Sprintf("%s:%d", conf.CommentServiceHost, conf.CommentServicePort),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	return nil, err
	//}
	//
	//connLike, err := grpc.Dial(
	//	fmt.Sprintf("%s:%d", conf.LikeServiceHost, conf.LikeServicePort),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	return nil, err
	//}

	serviceManager := &serviceManager{
		userService: pbu.NewUserServiceClient(connUser),
		//postService:    pbp.NewPostServiceClient(connPost),
		//commentService: pbc.NewCommentServiceClient(connComment),
		//likeService:    pbl.NewLikeServiceClient(connLike),
	}

	return serviceManager, nil
}
