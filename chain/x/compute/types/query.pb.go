
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

type QueryGetNodeRequest struct {
	NodeId string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (m *QueryGetNodeRequest) Reset()         { *m = QueryGetNodeRequest{} }
func (m *QueryGetNodeRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetNodeRequest) ProtoMessage()    {}

type QueryGetNodeResponse struct {
	Node *Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
}

func (m *QueryGetNodeResponse) Reset()         { *m = QueryGetNodeResponse{} }
func (m *QueryGetNodeResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetNodeResponse) ProtoMessage()    {}

type QueryListNodesRequest struct {
}

func (m *QueryListNodesRequest) Reset()         { *m = QueryListNodesRequest{} }
func (m *QueryListNodesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryListNodesRequest) ProtoMessage()    {}

type QueryListNodesResponse struct {
	Nodes []Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes"`
}

func (m *QueryListNodesResponse) Reset()         { *m = QueryListNodesResponse{} }
func (m *QueryListNodesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryListNodesResponse) ProtoMessage()    {}

type QueryClient interface {
	GetNode(ctx context.Context, in *QueryGetNodeRequest, opts ...grpc.CallOption) (*QueryGetNodeResponse, error)
	ListNodes(ctx context.Context, in *QueryListNodesRequest, opts ...grpc.CallOption) (*QueryListNodesResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) GetNode(ctx context.Context, in *QueryGetNodeRequest, opts ...grpc.CallOption) (*QueryGetNodeResponse, error) {
	out := new(QueryGetNodeResponse)
	err := c.cc.Invoke(ctx, "/atlas.compute.Query/GetNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListNodes(ctx context.Context, in *QueryListNodesRequest, opts ...grpc.CallOption) (*QueryListNodesResponse, error) {
	out := new(QueryListNodesResponse)
	err := c.cc.Invoke(ctx, "/atlas.compute.Query/ListNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type QueryServer interface {
	GetNode(context.Context, *QueryGetNodeRequest) (*QueryGetNodeResponse, error)
	ListNodes(context.Context, *QueryListNodesRequest) (*QueryListNodesResponse, error)
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_GetNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.compute.Query/GetNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetNode(ctx, req.(*QueryGetNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.compute.Query/ListNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListNodes(ctx, req.(*QueryListNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "atlas.compute.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNode",
			Handler:    _Query_GetNode_Handler,
		},
		{
			MethodName: "ListNodes",
			Handler:    _Query_ListNodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "atlas/compute/query.proto",
}

