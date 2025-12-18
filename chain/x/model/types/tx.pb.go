
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

type MsgRegisterModel struct {
	Creator  string                 `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Name     string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Version  string                 `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Cid      string                 `protobuf:"bytes,4,opt,name=cid,proto3" json:"cid,omitempty"`
	Metadata map[string]interface{} `protobuf:"bytes,5,rep,name=metadata,proto3" json:"metadata,omitempty"`
}

func (m *MsgRegisterModel) Reset()         { *m = MsgRegisterModel{} }
func (m *MsgRegisterModel) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterModel) ProtoMessage()    {}

type MsgRegisterModelResponse struct {
	ModelId string `protobuf:"bytes,1,opt,name=model_id,json=modelId,proto3" json:"model_id,omitempty"`
}

func (m *MsgRegisterModelResponse) Reset()         { *m = MsgRegisterModelResponse{} }
func (m *MsgRegisterModelResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterModelResponse) ProtoMessage()    {}

type MsgClient interface {
	RegisterModel(ctx context.Context, in *MsgRegisterModel, opts ...grpc.CallOption) (*MsgRegisterModelResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RegisterModel(ctx context.Context, in *MsgRegisterModel, opts ...grpc.CallOption) (*MsgRegisterModelResponse, error) {
	out := new(MsgRegisterModelResponse)
	err := c.cc.Invoke(ctx, "/atlas.model.Msg/RegisterModel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type MsgServer interface {
	RegisterModel(context.Context, *MsgRegisterModel) (*MsgRegisterModelResponse, error)
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_RegisterModel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterModel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterModel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.model.Msg/RegisterModel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterModel(ctx, req.(*MsgRegisterModel))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "atlas.model.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterModel",
			Handler:    _Msg_RegisterModel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "atlas/model/tx.proto",
}

