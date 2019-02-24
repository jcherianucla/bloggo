package handlers

import (
	"context"

	"go.uber.org/fx"

	"github.com/jcherianucla/bloggo/.gen/idl/proto"
)

var Module = fx.Provide(New)

type bloggoHandler struct {
}

func New() proto.BloggoServer {
	return bloggoHandler{}
}

func (bh bloggoHandler) CreatePost(
	ctx context.Context,
	in *proto.CreatePostRequest,
) (*proto.CreatePostResponse, error) {
	return nil, nil
}
func (bh bloggoHandler) FetchPost(
	ctx context.Context,
	in *proto.FetchPostRequest,
) (*proto.FetchPostResponse, error) {
	return nil, nil
}
func (bh bloggoHandler) UpdatePost(
	ctx context.Context,
	in *proto.UpdatePostRequest,
) (*proto.UpdatePostResponse, error) {
	return nil, nil
}
func (bh bloggoHandler) DeletePost(
	ctx context.Context,
	in *proto.DeletePostRequest,
) (*proto.DeletePostResponse, error) {
	return nil, nil
}
