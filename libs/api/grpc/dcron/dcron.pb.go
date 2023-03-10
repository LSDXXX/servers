// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dcron.proto

package dcron

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RspCode int32

const (
	RspCode_StateSuccess RspCode = 0
	RspCode_StateFailed  RspCode = 1
)

var RspCode_name = map[int32]string{
	0: "StateSuccess",
	1: "StateFailed",
}

var RspCode_value = map[string]int32{
	"StateSuccess": 0,
	"StateFailed":  1,
}

func (x RspCode) String() string {
	return proto.EnumName(RspCode_name, int32(x))
}

func (RspCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_93aa2ed5a7a309a6, []int{0}
}

type CronTaskReq struct {
	Spec                 string               `protobuf:"bytes,1,opt,name=spec,proto3" json:"spec,omitempty"`
	TaskId               string               `protobuf:"bytes,2,opt,name=taskId,proto3" json:"taskId,omitempty"`
	ProjectId            int32                `protobuf:"varint,3,opt,name=projectId,proto3" json:"projectId,omitempty"`
	StartTime            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime              *timestamp.Timestamp `protobuf:"bytes,5,opt,name=endTime,proto3" json:"endTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CronTaskReq) Reset()         { *m = CronTaskReq{} }
func (m *CronTaskReq) String() string { return proto.CompactTextString(m) }
func (*CronTaskReq) ProtoMessage()    {}
func (*CronTaskReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_93aa2ed5a7a309a6, []int{0}
}

func (m *CronTaskReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CronTaskReq.Unmarshal(m, b)
}
func (m *CronTaskReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CronTaskReq.Marshal(b, m, deterministic)
}
func (m *CronTaskReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CronTaskReq.Merge(m, src)
}
func (m *CronTaskReq) XXX_Size() int {
	return xxx_messageInfo_CronTaskReq.Size(m)
}
func (m *CronTaskReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CronTaskReq.DiscardUnknown(m)
}

var xxx_messageInfo_CronTaskReq proto.InternalMessageInfo

func (m *CronTaskReq) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

func (m *CronTaskReq) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *CronTaskReq) GetProjectId() int32 {
	if m != nil {
		return m.ProjectId
	}
	return 0
}

func (m *CronTaskReq) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *CronTaskReq) GetEndTime() *timestamp.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

type RemoveTaskReq struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=taskId,proto3" json:"taskId,omitempty"`
	ProjectId            int32    `protobuf:"varint,2,opt,name=projectId,proto3" json:"projectId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveTaskReq) Reset()         { *m = RemoveTaskReq{} }
func (m *RemoveTaskReq) String() string { return proto.CompactTextString(m) }
func (*RemoveTaskReq) ProtoMessage()    {}
func (*RemoveTaskReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_93aa2ed5a7a309a6, []int{1}
}

func (m *RemoveTaskReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveTaskReq.Unmarshal(m, b)
}
func (m *RemoveTaskReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveTaskReq.Marshal(b, m, deterministic)
}
func (m *RemoveTaskReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveTaskReq.Merge(m, src)
}
func (m *RemoveTaskReq) XXX_Size() int {
	return xxx_messageInfo_RemoveTaskReq.Size(m)
}
func (m *RemoveTaskReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveTaskReq.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveTaskReq proto.InternalMessageInfo

func (m *RemoveTaskReq) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *RemoveTaskReq) GetProjectId() int32 {
	if m != nil {
		return m.ProjectId
	}
	return 0
}

type ResultRsp struct {
	Code                 RspCode  `protobuf:"varint,1,opt,name=code,proto3,enum=dcron.RspCode" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResultRsp) Reset()         { *m = ResultRsp{} }
func (m *ResultRsp) String() string { return proto.CompactTextString(m) }
func (*ResultRsp) ProtoMessage()    {}
func (*ResultRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_93aa2ed5a7a309a6, []int{2}
}

func (m *ResultRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResultRsp.Unmarshal(m, b)
}
func (m *ResultRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResultRsp.Marshal(b, m, deterministic)
}
func (m *ResultRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResultRsp.Merge(m, src)
}
func (m *ResultRsp) XXX_Size() int {
	return xxx_messageInfo_ResultRsp.Size(m)
}
func (m *ResultRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_ResultRsp.DiscardUnknown(m)
}

var xxx_messageInfo_ResultRsp proto.InternalMessageInfo

func (m *ResultRsp) GetCode() RspCode {
	if m != nil {
		return m.Code
	}
	return RspCode_StateSuccess
}

type TriggerTaskReq struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=taskId,proto3" json:"taskId,omitempty"`
	ProjectId            int32    `protobuf:"varint,2,opt,name=projectId,proto3" json:"projectId,omitempty"`
	Data                 *any.Any `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TriggerTaskReq) Reset()         { *m = TriggerTaskReq{} }
func (m *TriggerTaskReq) String() string { return proto.CompactTextString(m) }
func (*TriggerTaskReq) ProtoMessage()    {}
func (*TriggerTaskReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_93aa2ed5a7a309a6, []int{3}
}

func (m *TriggerTaskReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TriggerTaskReq.Unmarshal(m, b)
}
func (m *TriggerTaskReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TriggerTaskReq.Marshal(b, m, deterministic)
}
func (m *TriggerTaskReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TriggerTaskReq.Merge(m, src)
}
func (m *TriggerTaskReq) XXX_Size() int {
	return xxx_messageInfo_TriggerTaskReq.Size(m)
}
func (m *TriggerTaskReq) XXX_DiscardUnknown() {
	xxx_messageInfo_TriggerTaskReq.DiscardUnknown(m)
}

var xxx_messageInfo_TriggerTaskReq proto.InternalMessageInfo

func (m *TriggerTaskReq) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *TriggerTaskReq) GetProjectId() int32 {
	if m != nil {
		return m.ProjectId
	}
	return 0
}

func (m *TriggerTaskReq) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("dcron.RspCode", RspCode_name, RspCode_value)
	proto.RegisterType((*CronTaskReq)(nil), "dcron.CronTaskReq")
	proto.RegisterType((*RemoveTaskReq)(nil), "dcron.RemoveTaskReq")
	proto.RegisterType((*ResultRsp)(nil), "dcron.ResultRsp")
	proto.RegisterType((*TriggerTaskReq)(nil), "dcron.TriggerTaskReq")
}

func init() { proto.RegisterFile("dcron.proto", fileDescriptor_93aa2ed5a7a309a6) }

var fileDescriptor_93aa2ed5a7a309a6 = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x51, 0xc1, 0x8e, 0x94, 0x40,
	0x10, 0x95, 0x95, 0x19, 0x42, 0xa1, 0x23, 0xa9, 0xac, 0x06, 0x89, 0x89, 0x13, 0x4e, 0xc4, 0x18,
	0x30, 0xe8, 0xc1, 0x8b, 0x87, 0xcd, 0x46, 0x93, 0xbd, 0xf6, 0xf0, 0x03, 0x3d, 0x74, 0x49, 0x70,
	0x07, 0x1a, 0xbb, 0x1b, 0x93, 0xfd, 0x42, 0x7f, 0xcb, 0xd8, 0xc0, 0x30, 0xce, 0xc4, 0x78, 0xf0,
	0xd6, 0xf5, 0xea, 0x3d, 0x78, 0xef, 0x15, 0x04, 0xa2, 0x52, 0xb2, 0xcb, 0x7a, 0x25, 0x8d, 0xc4,
	0x95, 0x1d, 0xe2, 0xd7, 0xb5, 0x94, 0xf5, 0x81, 0x72, 0x0b, 0xee, 0x87, 0xaf, 0xb9, 0x69, 0x5a,
	0xd2, 0x86, 0xb7, 0xfd, 0xc8, 0x8b, 0x5f, 0x9e, 0x13, 0x78, 0xf7, 0x30, 0xae, 0x92, 0x9f, 0x0e,
	0x04, 0xb7, 0x4a, 0x76, 0x25, 0xd7, 0xf7, 0x8c, 0xbe, 0x23, 0x82, 0xab, 0x7b, 0xaa, 0x22, 0x67,
	0xeb, 0xa4, 0x3e, 0xb3, 0x6f, 0x7c, 0x01, 0x6b, 0xc3, 0xf5, 0xfd, 0x9d, 0x88, 0xae, 0x2c, 0x3a,
	0x4d, 0xf8, 0x0a, 0xfc, 0x5e, 0xc9, 0x6f, 0x54, 0x99, 0x3b, 0x11, 0x3d, 0xde, 0x3a, 0xe9, 0x8a,
	0x2d, 0x00, 0x7e, 0x04, 0x5f, 0x1b, 0xae, 0x4c, 0xd9, 0xb4, 0x14, 0xb9, 0x5b, 0x27, 0x0d, 0x8a,
	0x38, 0x1b, 0x8d, 0x64, 0xb3, 0x91, 0xac, 0x9c, 0x9d, 0xb2, 0x85, 0x8c, 0x1f, 0xc0, 0xa3, 0x4e,
	0x58, 0xdd, 0xea, 0x9f, 0xba, 0x99, 0x9a, 0x7c, 0x86, 0xa7, 0x8c, 0x5a, 0xf9, 0x83, 0xe6, 0x28,
	0x8b, 0x6d, 0xe7, 0xef, 0xb6, 0xaf, 0xce, 0x6c, 0x27, 0x39, 0xf8, 0x8c, 0xf4, 0x70, 0x30, 0x4c,
	0xf7, 0x98, 0x80, 0x5b, 0x49, 0x41, 0xf6, 0x03, 0x9b, 0x62, 0x93, 0x8d, 0xe5, 0x33, 0xdd, 0xdf,
	0x4a, 0x41, 0xcc, 0xee, 0x92, 0x1e, 0x36, 0xa5, 0x6a, 0xea, 0x9a, 0xd4, 0x7f, 0xfd, 0x18, 0x53,
	0x70, 0x05, 0x37, 0xdc, 0x16, 0x19, 0x14, 0xd7, 0x17, 0x91, 0x6f, 0xba, 0x07, 0x66, 0x19, 0x6f,
	0xde, 0x82, 0x37, 0x59, 0xc0, 0x10, 0x9e, 0xec, 0x0c, 0x37, 0xb4, 0x1b, 0xaa, 0x8a, 0xb4, 0x0e,
	0x1f, 0xe1, 0x33, 0x08, 0x2c, 0xf2, 0x85, 0x37, 0x07, 0x12, 0xa1, 0x53, 0x74, 0xe0, 0xef, 0x86,
	0x7d, 0xdb, 0x18, 0x43, 0x0a, 0x73, 0xf0, 0x6e, 0x84, 0xf8, 0x7d, 0x70, 0xc4, 0x29, 0xcd, 0xc9,
	0xf5, 0xe3, 0x70, 0x4e, 0x78, 0x6c, 0xe0, 0x1d, 0xac, 0xc7, 0x56, 0xf1, 0xfa, 0xb8, 0x3b, 0x29,
	0xf9, 0x52, 0x51, 0x7c, 0x02, 0x6f, 0xea, 0x03, 0x8b, 0xe5, 0xf9, 0x7c, 0xe2, 0xfd, 0x59, 0xd5,
	0xa5, 0x7c, 0xbf, 0xb6, 0x81, 0xdf, 0xff, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x8e, 0xd8, 0x95,
	0xe9, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SubmitterClient is the client API for Submitter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SubmitterClient interface {
	AddCron(ctx context.Context, in *CronTaskReq, opts ...grpc.CallOption) (*ResultRsp, error)
	Remove(ctx context.Context, in *RemoveTaskReq, opts ...grpc.CallOption) (*ResultRsp, error)
}

type submitterClient struct {
	cc *grpc.ClientConn
}

func NewSubmitterClient(cc *grpc.ClientConn) SubmitterClient {
	return &submitterClient{cc}
}

func (c *submitterClient) AddCron(ctx context.Context, in *CronTaskReq, opts ...grpc.CallOption) (*ResultRsp, error) {
	out := new(ResultRsp)
	err := c.cc.Invoke(ctx, "/dcron.Submitter/AddCron", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *submitterClient) Remove(ctx context.Context, in *RemoveTaskReq, opts ...grpc.CallOption) (*ResultRsp, error) {
	out := new(ResultRsp)
	err := c.cc.Invoke(ctx, "/dcron.Submitter/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubmitterServer is the server API for Submitter service.
type SubmitterServer interface {
	AddCron(context.Context, *CronTaskReq) (*ResultRsp, error)
	Remove(context.Context, *RemoveTaskReq) (*ResultRsp, error)
}

func RegisterSubmitterServer(s *grpc.Server, srv SubmitterServer) {
	s.RegisterService(&_Submitter_serviceDesc, srv)
}

func _Submitter_AddCron_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CronTaskReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubmitterServer).AddCron(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dcron.Submitter/AddCron",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubmitterServer).AddCron(ctx, req.(*CronTaskReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Submitter_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTaskReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubmitterServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dcron.Submitter/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubmitterServer).Remove(ctx, req.(*RemoveTaskReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Submitter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dcron.Submitter",
	HandlerType: (*SubmitterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddCron",
			Handler:    _Submitter_AddCron_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Submitter_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dcron.proto",
}

// TriggerClient is the client API for Trigger service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TriggerClient interface {
	Trigger(ctx context.Context, in *TriggerTaskReq, opts ...grpc.CallOption) (*ResultRsp, error)
}

type triggerClient struct {
	cc *grpc.ClientConn
}

func NewTriggerClient(cc *grpc.ClientConn) TriggerClient {
	return &triggerClient{cc}
}

func (c *triggerClient) Trigger(ctx context.Context, in *TriggerTaskReq, opts ...grpc.CallOption) (*ResultRsp, error) {
	out := new(ResultRsp)
	err := c.cc.Invoke(ctx, "/dcron.Trigger/Trigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TriggerServer is the server API for Trigger service.
type TriggerServer interface {
	Trigger(context.Context, *TriggerTaskReq) (*ResultRsp, error)
}

func RegisterTriggerServer(s *grpc.Server, srv TriggerServer) {
	s.RegisterService(&_Trigger_serviceDesc, srv)
}

func _Trigger_Trigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerTaskReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TriggerServer).Trigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dcron.Trigger/Trigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TriggerServer).Trigger(ctx, req.(*TriggerTaskReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Trigger_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dcron.Trigger",
	HandlerType: (*TriggerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Trigger",
			Handler:    _Trigger_Trigger_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dcron.proto",
}
