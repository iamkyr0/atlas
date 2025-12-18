
package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

type QueryGetModelRequest struct {
	ModelId string `protobuf:"bytes,1,opt,name=model_id,json=modelId,proto3" json:"model_id,omitempty"`
}

func (m *QueryGetModelRequest) Reset()         { *m = QueryGetModelRequest{} }
func (m *QueryGetModelRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetModelRequest) ProtoMessage()    {}

type QueryGetModelResponse struct {
	Model *ModelProto `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
}

type ModelProto struct {
	Id        string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Version   string            `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Cid       string            `protobuf:"bytes,4,opt,name=cid,proto3" json:"cid,omitempty"`
	CreatedAt int64             `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Metadata  map[string]string `protobuf:"bytes,6,rep,name=metadata,proto3" json:"metadata,omitempty"`
}

func (m *ModelProto) Reset()         { *m = ModelProto{} }
func (m *ModelProto) String() string { return "" }
func (*ModelProto) ProtoMessage()    {}

func (m *QueryGetModelResponse) Reset()         { *m = QueryGetModelResponse{} }
func (m *QueryGetModelResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetModelResponse) ProtoMessage()    {}

type QueryListModelsRequest struct {
}

func (m *QueryListModelsRequest) Reset()         { *m = QueryListModelsRequest{} }
func (m *QueryListModelsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryListModelsRequest) ProtoMessage()    {}

type QueryListModelsResponse struct {
	Models []ModelProto `protobuf:"bytes,1,rep,name=models,proto3" json:"models"`
}

func (m *QueryListModelsResponse) Reset()         { *m = QueryListModelsResponse{} }
func (m *QueryListModelsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryListModelsResponse) ProtoMessage()    {}

type QueryClient interface {
	GetModel(ctx context.Context, in *QueryGetModelRequest, opts ...grpc.CallOption) (*QueryGetModelResponse, error)
	ListModels(ctx context.Context, in *QueryListModelsRequest, opts ...grpc.CallOption) (*QueryListModelsResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) GetModel(ctx context.Context, in *QueryGetModelRequest, opts ...grpc.CallOption) (*QueryGetModelResponse, error) {
	out := new(QueryGetModelResponse)
	err := c.cc.Invoke(ctx, "/atlas.model.Query/GetModel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListModels(ctx context.Context, in *QueryListModelsRequest, opts ...grpc.CallOption) (*QueryListModelsResponse, error) {
	out := new(QueryListModelsResponse)
	err := c.cc.Invoke(ctx, "/atlas.model.Query/ListModels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type QueryServer interface {
	GetModel(context.Context, *QueryGetModelRequest) (*QueryGetModelResponse, error)
	ListModels(context.Context, *QueryListModelsRequest) (*QueryListModelsResponse, error)
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_GetModel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetModelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetModel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.model.Query/GetModel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetModel(ctx, req.(*QueryGetModelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListModels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListModelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListModels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.model.Query/ListModels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListModels(ctx, req.(*QueryListModelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "atlas.model.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetModel",
			Handler:    _Query_GetModel_Handler,
		},
		{
			MethodName: "ListModels",
			Handler:    _Query_ListModels_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "atlas/model/query.proto",
}

