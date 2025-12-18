
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

type MsgSubmitJob struct {
	Creator    string                 `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	ModelId    string                 `protobuf:"bytes,2,opt,name=model_id,json=modelId,proto3" json:"model_id,omitempty"`
	DatasetCid string                 `protobuf:"bytes,3,opt,name=dataset_cid,json=datasetCid,proto3" json:"dataset_cid,omitempty"`
	Config     map[string]interface{} `protobuf:"bytes,4,rep,name=config,proto3" json:"config,omitempty"`
}

func (m *MsgSubmitJob) Reset()         { *m = MsgSubmitJob{} }
func (m *MsgSubmitJob) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitJob) ProtoMessage()    {}

type MsgSubmitJobResponse struct {
	JobId string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (m *MsgSubmitJobResponse) Reset()         { *m = MsgSubmitJobResponse{} }
func (m *MsgSubmitJobResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitJobResponse) ProtoMessage()    {}

type MsgCreateTask struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	JobId   string `protobuf:"bytes,2,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	ShardId string `protobuf:"bytes,3,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	NodeId  string `protobuf:"bytes,4,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (m *MsgCreateTask) Reset()         { *m = MsgCreateTask{} }
func (m *MsgCreateTask) String() string { return proto.CompactTextString(m) }
func (*MsgCreateTask) ProtoMessage()    {}

type MsgCreateTaskResponse struct {
	TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

func (m *MsgCreateTaskResponse) Reset()         { *m = MsgCreateTaskResponse{} }
func (m *MsgCreateTaskResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCreateTaskResponse) ProtoMessage()    {}

type MsgUpdateTaskStatus struct {
	Creator       string  `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	TaskId        string  `protobuf:"bytes,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Status        string  `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Progress      float64 `protobuf:"fixed64,4,opt,name=progress,proto3" json:"progress,omitempty"`
	CheckpointCid string  `protobuf:"bytes,5,opt,name=checkpoint_cid,json=checkpointCid,proto3" json:"checkpoint_cid,omitempty"`
}

func (m *MsgUpdateTaskStatus) Reset()         { *m = MsgUpdateTaskStatus{} }
func (m *MsgUpdateTaskStatus) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateTaskStatus) ProtoMessage()    {}

type MsgUpdateTaskStatusResponse struct {
}

func (m *MsgUpdateTaskStatusResponse) Reset()         { *m = MsgUpdateTaskStatusResponse{} }
func (m *MsgUpdateTaskStatusResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateTaskStatusResponse) ProtoMessage()    {}

type MsgClient interface {
	SubmitJob(ctx context.Context, in *MsgSubmitJob, opts ...grpc.CallOption) (*MsgSubmitJobResponse, error)
	CreateTask(ctx context.Context, in *MsgCreateTask, opts ...grpc.CallOption) (*MsgCreateTaskResponse, error)
	UpdateTaskStatus(ctx context.Context, in *MsgUpdateTaskStatus, opts ...grpc.CallOption) (*MsgUpdateTaskStatusResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SubmitJob(ctx context.Context, in *MsgSubmitJob, opts ...grpc.CallOption) (*MsgSubmitJobResponse, error) {
	out := new(MsgSubmitJobResponse)
	err := c.cc.Invoke(ctx, "/atlas.training.Msg/SubmitJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CreateTask(ctx context.Context, in *MsgCreateTask, opts ...grpc.CallOption) (*MsgCreateTaskResponse, error) {
	out := new(MsgCreateTaskResponse)
	err := c.cc.Invoke(ctx, "/atlas.training.Msg/CreateTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateTaskStatus(ctx context.Context, in *MsgUpdateTaskStatus, opts ...grpc.CallOption) (*MsgUpdateTaskStatusResponse, error) {
	out := new(MsgUpdateTaskStatusResponse)
	err := c.cc.Invoke(ctx, "/atlas.training.Msg/UpdateTaskStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type MsgServer interface {
	SubmitJob(context.Context, *MsgSubmitJob) (*MsgSubmitJobResponse, error)
	CreateTask(context.Context, *MsgCreateTask) (*MsgCreateTaskResponse, error)
	UpdateTaskStatus(context.Context, *MsgUpdateTaskStatus) (*MsgUpdateTaskStatusResponse, error)
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SubmitJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitJob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.training.Msg/SubmitJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitJob(ctx, req.(*MsgSubmitJob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.training.Msg/CreateTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateTask(ctx, req.(*MsgCreateTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateTaskStatus)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.training.Msg/UpdateTaskStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateTaskStatus(ctx, req.(*MsgUpdateTaskStatus))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "atlas.training.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitJob",
			Handler:    _Msg_SubmitJob_Handler,
		},
		{
			MethodName: "CreateTask",
			Handler:    _Msg_CreateTask_Handler,
		},
		{
			MethodName: "UpdateTaskStatus",
			Handler:    _Msg_UpdateTaskStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "atlas/training/tx.proto",
}

