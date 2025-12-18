
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

type MsgRegisterNode struct {
	Creator   string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	NodeId    string `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Address   string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	CpuCores  int32  `protobuf:"varint,4,opt,name=cpu_cores,json=cpuCores,proto3" json:"cpu_cores,omitempty"`
	GpuCount  int32  `protobuf:"varint,5,opt,name=gpu_count,json=gpuCount,proto3" json:"gpu_count,omitempty"`
	MemoryGb  int32  `protobuf:"varint,6,opt,name=memory_gb,json=memoryGb,proto3" json:"memory_gb,omitempty"`
	StorageGb int32  `protobuf:"varint,7,opt,name=storage_gb,json=storageGb,proto3" json:"storage_gb,omitempty"`
}

func (m *MsgRegisterNode) Reset()         { *m = MsgRegisterNode{} }
func (m *MsgRegisterNode) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterNode) ProtoMessage()    {}

type MsgRegisterNodeResponse struct {
}

func (m *MsgRegisterNodeResponse) Reset()         { *m = MsgRegisterNodeResponse{} }
func (m *MsgRegisterNodeResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterNodeResponse) ProtoMessage()    {}

type MsgUpdateHeartbeat struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	NodeId  string `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (m *MsgUpdateHeartbeat) Reset()         { *m = MsgUpdateHeartbeat{} }
func (m *MsgUpdateHeartbeat) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateHeartbeat) ProtoMessage()    {}

type MsgUpdateHeartbeatResponse struct {
}

func (m *MsgUpdateHeartbeatResponse) Reset()         { *m = MsgUpdateHeartbeatResponse{} }
func (m *MsgUpdateHeartbeatResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateHeartbeatResponse) ProtoMessage()    {}

type MsgClient interface {
	RegisterNode(ctx context.Context, in *MsgRegisterNode, opts ...grpc.CallOption) (*MsgRegisterNodeResponse, error)
	UpdateHeartbeat(ctx context.Context, in *MsgUpdateHeartbeat, opts ...grpc.CallOption) (*MsgUpdateHeartbeatResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RegisterNode(ctx context.Context, in *MsgRegisterNode, opts ...grpc.CallOption) (*MsgRegisterNodeResponse, error) {
	out := new(MsgRegisterNodeResponse)
	err := c.cc.Invoke(ctx, "/atlas.compute.Msg/RegisterNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateHeartbeat(ctx context.Context, in *MsgUpdateHeartbeat, opts ...grpc.CallOption) (*MsgUpdateHeartbeatResponse, error) {
	out := new(MsgUpdateHeartbeatResponse)
	err := c.cc.Invoke(ctx, "/atlas.compute.Msg/UpdateHeartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type MsgServer interface {
	RegisterNode(context.Context, *MsgRegisterNode) (*MsgRegisterNodeResponse, error)
	UpdateHeartbeat(context.Context, *MsgUpdateHeartbeat) (*MsgUpdateHeartbeatResponse, error)
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_RegisterNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.compute.Msg/RegisterNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterNode(ctx, req.(*MsgRegisterNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateHeartbeat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.compute.Msg/UpdateHeartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateHeartbeat(ctx, req.(*MsgUpdateHeartbeat))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "atlas.compute.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterNode",
			Handler:    _Msg_RegisterNode_Handler,
		},
		{
			MethodName: "UpdateHeartbeat",
			Handler:    _Msg_UpdateHeartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "atlas/compute/tx.proto",
}

