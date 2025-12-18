
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

type QueryGetJobRequest struct {
	JobId string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (m *QueryGetJobRequest) Reset()         { *m = QueryGetJobRequest{} }
func (m *QueryGetJobRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetJobRequest) ProtoMessage()    {}

type QueryGetJobResponse struct {
	Job *Job `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
}

func (m *QueryGetJobResponse) Reset()         { *m = QueryGetJobResponse{} }
func (m *QueryGetJobResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetJobResponse) ProtoMessage()    {}

type QueryListJobsRequest struct {
}

func (m *QueryListJobsRequest) Reset()         { *m = QueryListJobsRequest{} }
func (m *QueryListJobsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryListJobsRequest) ProtoMessage()    {}

type QueryListJobsResponse struct {
	Jobs []Job `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs"`
}

func (m *QueryListJobsResponse) Reset()         { *m = QueryListJobsResponse{} }
func (m *QueryListJobsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryListJobsResponse) ProtoMessage()    {}

type QueryGetTaskRequest struct {
	TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

func (m *QueryGetTaskRequest) Reset()         { *m = QueryGetTaskRequest{} }
func (m *QueryGetTaskRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetTaskRequest) ProtoMessage()    {}

type QueryGetTaskResponse struct {
	Task *Task `protobuf:"bytes,1,opt,name=task,proto3" json:"task,omitempty"`
}

func (m *QueryGetTaskResponse) Reset()         { *m = QueryGetTaskResponse{} }
func (m *QueryGetTaskResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetTaskResponse) ProtoMessage()    {}

type QueryGetTasksByJobRequest struct {
	JobId string `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (m *QueryGetTasksByJobRequest) Reset()         { *m = QueryGetTasksByJobRequest{} }
func (m *QueryGetTasksByJobRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetTasksByJobRequest) ProtoMessage()    {}

type QueryGetTasksByJobResponse struct {
	Tasks []Task `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks"`
}

func (m *QueryGetTasksByJobResponse) Reset()         { *m = QueryGetTasksByJobResponse{} }
func (m *QueryGetTasksByJobResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetTasksByJobResponse) ProtoMessage()    {}

type QueryClient interface {
	GetJob(ctx context.Context, in *QueryGetJobRequest, opts ...grpc.CallOption) (*QueryGetJobResponse, error)
	ListJobs(ctx context.Context, in *QueryListJobsRequest, opts ...grpc.CallOption) (*QueryListJobsResponse, error)
	GetTask(ctx context.Context, in *QueryGetTaskRequest, opts ...grpc.CallOption) (*QueryGetTaskResponse, error)
	GetTasksByJob(ctx context.Context, in *QueryGetTasksByJobRequest, opts ...grpc.CallOption) (*QueryGetTasksByJobResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) GetJob(ctx context.Context, in *QueryGetJobRequest, opts ...grpc.CallOption) (*QueryGetJobResponse, error) {
	out := new(QueryGetJobResponse)
	err := c.cc.Invoke(ctx, "/atlas.training.Query/GetJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListJobs(ctx context.Context, in *QueryListJobsRequest, opts ...grpc.CallOption) (*QueryListJobsResponse, error) {
	out := new(QueryListJobsResponse)
	err := c.cc.Invoke(ctx, "/atlas.training.Query/ListJobs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetTask(ctx context.Context, in *QueryGetTaskRequest, opts ...grpc.CallOption) (*QueryGetTaskResponse, error) {
	out := new(QueryGetTaskResponse)
	err := c.cc.Invoke(ctx, "/atlas.training.Query/GetTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetTasksByJob(ctx context.Context, in *QueryGetTasksByJobRequest, opts ...grpc.CallOption) (*QueryGetTasksByJobResponse, error) {
	out := new(QueryGetTasksByJobResponse)
	err := c.cc.Invoke(ctx, "/atlas.training.Query/GetTasksByJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type QueryServer interface {
	GetJob(context.Context, *QueryGetJobRequest) (*QueryGetJobResponse, error)
	ListJobs(context.Context, *QueryListJobsRequest) (*QueryListJobsResponse, error)
	GetTask(context.Context, *QueryGetTaskRequest) (*QueryGetTaskResponse, error)
	GetTasksByJob(context.Context, *QueryGetTasksByJobRequest) (*QueryGetTasksByJobResponse, error)
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_GetJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.training.Query/GetJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetJob(ctx, req.(*QueryGetJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.training.Query/ListJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListJobs(ctx, req.(*QueryListJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.training.Query/GetTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetTask(ctx, req.(*QueryGetTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetTasksByJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetTasksByJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetTasksByJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atlas.training.Query/GetTasksByJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetTasksByJob(ctx, req.(*QueryGetTasksByJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "atlas.training.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetJob",
			Handler:    _Query_GetJob_Handler,
		},
		{
			MethodName: "ListJobs",
			Handler:    _Query_ListJobs_Handler,
		},
		{
			MethodName: "GetTask",
			Handler:    _Query_GetTask_Handler,
		},
		{
			MethodName: "GetTasksByJob",
			Handler:    _Query_GetTasksByJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "atlas/training/query.proto",
}

