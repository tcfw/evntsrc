// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storer.proto

/*
	Package evntsrc_storer is a generated protocol buffer package.

	It is generated from these files:
		storer.proto

	It has these top-level messages:
		AcknowledgeRequest
		AcknowledgeResponse
		ExtendTTLRequest
		ExtendTTLResponse
		QueryRequest
		QueryTTLExpired
*/
package evntsrc_storer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"
import evntsrc_event "github.com/tcfw/evntsrc/internal/event/protos"

import time "time"

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

import types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AcknowledgeRequest struct {
	Stream  int32  `protobuf:"varint,1,opt,name=Stream,proto3" json:"Stream,omitempty"`
	EventID string `protobuf:"bytes,2,opt,name=EventID,proto3" json:"EventID,omitempty"`
}

func (m *AcknowledgeRequest) Reset()                    { *m = AcknowledgeRequest{} }
func (m *AcknowledgeRequest) String() string            { return proto.CompactTextString(m) }
func (*AcknowledgeRequest) ProtoMessage()               {}
func (*AcknowledgeRequest) Descriptor() ([]byte, []int) { return fileDescriptorStorer, []int{0} }

func (m *AcknowledgeRequest) GetStream() int32 {
	if m != nil {
		return m.Stream
	}
	return 0
}

func (m *AcknowledgeRequest) GetEventID() string {
	if m != nil {
		return m.EventID
	}
	return ""
}

func (*AcknowledgeRequest) XXX_MessageName() string {
	return "evntsrc.storer.AcknowledgeRequest"
}

type AcknowledgeResponse struct {
	Time *time.Time `protobuf:"bytes,1,opt,name=Time,stdtime" json:"Time,omitempty"`
}

func (m *AcknowledgeResponse) Reset()                    { *m = AcknowledgeResponse{} }
func (m *AcknowledgeResponse) String() string            { return proto.CompactTextString(m) }
func (*AcknowledgeResponse) ProtoMessage()               {}
func (*AcknowledgeResponse) Descriptor() ([]byte, []int) { return fileDescriptorStorer, []int{1} }

func (m *AcknowledgeResponse) GetTime() *time.Time {
	if m != nil {
		return m.Time
	}
	return nil
}

func (*AcknowledgeResponse) XXX_MessageName() string {
	return "evntsrc.storer.AcknowledgeResponse"
}

type ExtendTTLRequest struct {
	Stream     int32      `protobuf:"varint,1,opt,name=Stream,proto3" json:"Stream,omitempty"`
	EventID    string     `protobuf:"bytes,2,opt,name=EventID,proto3" json:"EventID,omitempty"`
	CurrentTTL *time.Time `protobuf:"bytes,3,opt,name=CurrentTTL,stdtime" json:"CurrentTTL,omitempty"`
	TTLTime    *time.Time `protobuf:"bytes,4,opt,name=TTLTime,stdtime" json:"TTLTime,omitempty"`
}

func (m *ExtendTTLRequest) Reset()                    { *m = ExtendTTLRequest{} }
func (m *ExtendTTLRequest) String() string            { return proto.CompactTextString(m) }
func (*ExtendTTLRequest) ProtoMessage()               {}
func (*ExtendTTLRequest) Descriptor() ([]byte, []int) { return fileDescriptorStorer, []int{2} }

func (m *ExtendTTLRequest) GetStream() int32 {
	if m != nil {
		return m.Stream
	}
	return 0
}

func (m *ExtendTTLRequest) GetEventID() string {
	if m != nil {
		return m.EventID
	}
	return ""
}

func (m *ExtendTTLRequest) GetCurrentTTL() *time.Time {
	if m != nil {
		return m.CurrentTTL
	}
	return nil
}

func (m *ExtendTTLRequest) GetTTLTime() *time.Time {
	if m != nil {
		return m.TTLTime
	}
	return nil
}

func (*ExtendTTLRequest) XXX_MessageName() string {
	return "evntsrc.storer.ExtendTTLRequest"
}

type ExtendTTLResponse struct {
}

func (m *ExtendTTLResponse) Reset()                    { *m = ExtendTTLResponse{} }
func (m *ExtendTTLResponse) String() string            { return proto.CompactTextString(m) }
func (*ExtendTTLResponse) ProtoMessage()               {}
func (*ExtendTTLResponse) Descriptor() ([]byte, []int) { return fileDescriptorStorer, []int{3} }

func (*ExtendTTLResponse) XXX_MessageName() string {
	return "evntsrc.storer.ExtendTTLResponse"
}

type QueryRequest struct {
	Stream int32 `protobuf:"varint,1,opt,name=Stream,proto3" json:"Stream,omitempty"`
	Limit  int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	// Types that are valid to be assigned to Query:
	//	*QueryRequest_Ttl
	Query isQueryRequest_Query `protobuf_oneof:"query"`
}

func (m *QueryRequest) Reset()                    { *m = QueryRequest{} }
func (m *QueryRequest) String() string            { return proto.CompactTextString(m) }
func (*QueryRequest) ProtoMessage()               {}
func (*QueryRequest) Descriptor() ([]byte, []int) { return fileDescriptorStorer, []int{4} }

type isQueryRequest_Query interface {
	isQueryRequest_Query()
	MarshalTo([]byte) (int, error)
	Size() int
}

type QueryRequest_Ttl struct {
	Ttl *QueryTTLExpired `protobuf:"bytes,3,opt,name=ttl,oneof"`
}

func (*QueryRequest_Ttl) isQueryRequest_Query() {}

func (m *QueryRequest) GetQuery() isQueryRequest_Query {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *QueryRequest) GetStream() int32 {
	if m != nil {
		return m.Stream
	}
	return 0
}

func (m *QueryRequest) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *QueryRequest) GetTtl() *QueryTTLExpired {
	if x, ok := m.GetQuery().(*QueryRequest_Ttl); ok {
		return x.Ttl
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*QueryRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _QueryRequest_OneofMarshaler, _QueryRequest_OneofUnmarshaler, _QueryRequest_OneofSizer, []interface{}{
		(*QueryRequest_Ttl)(nil),
	}
}

func _QueryRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*QueryRequest)
	// query
	switch x := m.Query.(type) {
	case *QueryRequest_Ttl:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ttl); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("QueryRequest.Query has unexpected type %T", x)
	}
	return nil
}

func _QueryRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*QueryRequest)
	switch tag {
	case 3: // query.ttl
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(QueryTTLExpired)
		err := b.DecodeMessage(msg)
		m.Query = &QueryRequest_Ttl{msg}
		return true, err
	default:
		return false, nil
	}
}

func _QueryRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*QueryRequest)
	// query
	switch x := m.Query.(type) {
	case *QueryRequest_Ttl:
		s := proto.Size(x.Ttl)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func (*QueryRequest) XXX_MessageName() string {
	return "evntsrc.storer.QueryRequest"
}

type QueryTTLExpired struct {
}

func (m *QueryTTLExpired) Reset()                    { *m = QueryTTLExpired{} }
func (m *QueryTTLExpired) String() string            { return proto.CompactTextString(m) }
func (*QueryTTLExpired) ProtoMessage()               {}
func (*QueryTTLExpired) Descriptor() ([]byte, []int) { return fileDescriptorStorer, []int{5} }

func (*QueryTTLExpired) XXX_MessageName() string {
	return "evntsrc.storer.QueryTTLExpired"
}
func init() {
	proto.RegisterType((*AcknowledgeRequest)(nil), "evntsrc.storer.AcknowledgeRequest")
	proto.RegisterType((*AcknowledgeResponse)(nil), "evntsrc.storer.AcknowledgeResponse")
	proto.RegisterType((*ExtendTTLRequest)(nil), "evntsrc.storer.ExtendTTLRequest")
	proto.RegisterType((*ExtendTTLResponse)(nil), "evntsrc.storer.ExtendTTLResponse")
	proto.RegisterType((*QueryRequest)(nil), "evntsrc.storer.QueryRequest")
	proto.RegisterType((*QueryTTLExpired)(nil), "evntsrc.storer.QueryTTLExpired")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for StorerService service

type StorerServiceClient interface {
	Acknowledge(ctx context.Context, in *AcknowledgeRequest, opts ...grpc.CallOption) (*AcknowledgeResponse, error)
	ExtendTTL(ctx context.Context, in *ExtendTTLRequest, opts ...grpc.CallOption) (*ExtendTTLResponse, error)
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (StorerService_QueryClient, error)
}

type storerServiceClient struct {
	cc *grpc.ClientConn
}

func NewStorerServiceClient(cc *grpc.ClientConn) StorerServiceClient {
	return &storerServiceClient{cc}
}

func (c *storerServiceClient) Acknowledge(ctx context.Context, in *AcknowledgeRequest, opts ...grpc.CallOption) (*AcknowledgeResponse, error) {
	out := new(AcknowledgeResponse)
	err := grpc.Invoke(ctx, "/evntsrc.storer.StorerService/Acknowledge", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storerServiceClient) ExtendTTL(ctx context.Context, in *ExtendTTLRequest, opts ...grpc.CallOption) (*ExtendTTLResponse, error) {
	out := new(ExtendTTLResponse)
	err := grpc.Invoke(ctx, "/evntsrc.storer.StorerService/ExtendTTL", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storerServiceClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (StorerService_QueryClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_StorerService_serviceDesc.Streams[0], c.cc, "/evntsrc.storer.StorerService/Query", opts...)
	if err != nil {
		return nil, err
	}
	x := &storerServiceQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StorerService_QueryClient interface {
	Recv() (*evntsrc_event.Event, error)
	grpc.ClientStream
}

type storerServiceQueryClient struct {
	grpc.ClientStream
}

func (x *storerServiceQueryClient) Recv() (*evntsrc_event.Event, error) {
	m := new(evntsrc_event.Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for StorerService service

type StorerServiceServer interface {
	Acknowledge(context.Context, *AcknowledgeRequest) (*AcknowledgeResponse, error)
	ExtendTTL(context.Context, *ExtendTTLRequest) (*ExtendTTLResponse, error)
	Query(*QueryRequest, StorerService_QueryServer) error
}

func RegisterStorerServiceServer(s *grpc.Server, srv StorerServiceServer) {
	s.RegisterService(&_StorerService_serviceDesc, srv)
}

func _StorerService_Acknowledge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcknowledgeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorerServiceServer).Acknowledge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/evntsrc.storer.StorerService/Acknowledge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorerServiceServer).Acknowledge(ctx, req.(*AcknowledgeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorerService_ExtendTTL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtendTTLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorerServiceServer).ExtendTTL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/evntsrc.storer.StorerService/ExtendTTL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorerServiceServer).ExtendTTL(ctx, req.(*ExtendTTLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorerService_Query_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorerServiceServer).Query(m, &storerServiceQueryServer{stream})
}

type StorerService_QueryServer interface {
	Send(*evntsrc_event.Event) error
	grpc.ServerStream
}

type storerServiceQueryServer struct {
	grpc.ServerStream
}

func (x *storerServiceQueryServer) Send(m *evntsrc_event.Event) error {
	return x.ServerStream.SendMsg(m)
}

var _StorerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "evntsrc.storer.StorerService",
	HandlerType: (*StorerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Acknowledge",
			Handler:    _StorerService_Acknowledge_Handler,
		},
		{
			MethodName: "ExtendTTL",
			Handler:    _StorerService_ExtendTTL_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _StorerService_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "storer.proto",
}

func (m *AcknowledgeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AcknowledgeRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Stream != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintStorer(dAtA, i, uint64(m.Stream))
	}
	if len(m.EventID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintStorer(dAtA, i, uint64(len(m.EventID)))
		i += copy(dAtA[i:], m.EventID)
	}
	return i, nil
}

func (m *AcknowledgeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AcknowledgeResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Time != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintStorer(dAtA, i, uint64(types.SizeOfStdTime(*m.Time)))
		n1, err := types.StdTimeMarshalTo(*m.Time, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *ExtendTTLRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExtendTTLRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Stream != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintStorer(dAtA, i, uint64(m.Stream))
	}
	if len(m.EventID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintStorer(dAtA, i, uint64(len(m.EventID)))
		i += copy(dAtA[i:], m.EventID)
	}
	if m.CurrentTTL != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintStorer(dAtA, i, uint64(types.SizeOfStdTime(*m.CurrentTTL)))
		n2, err := types.StdTimeMarshalTo(*m.CurrentTTL, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.TTLTime != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintStorer(dAtA, i, uint64(types.SizeOfStdTime(*m.TTLTime)))
		n3, err := types.StdTimeMarshalTo(*m.TTLTime, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}

func (m *ExtendTTLResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExtendTTLResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *QueryRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Stream != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintStorer(dAtA, i, uint64(m.Stream))
	}
	if m.Limit != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintStorer(dAtA, i, uint64(m.Limit))
	}
	if m.Query != nil {
		nn4, err := m.Query.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn4
	}
	return i, nil
}

func (m *QueryRequest_Ttl) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Ttl != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintStorer(dAtA, i, uint64(m.Ttl.Size()))
		n5, err := m.Ttl.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	return i, nil
}
func (m *QueryTTLExpired) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTTLExpired) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func encodeVarintStorer(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *AcknowledgeRequest) Size() (n int) {
	var l int
	_ = l
	if m.Stream != 0 {
		n += 1 + sovStorer(uint64(m.Stream))
	}
	l = len(m.EventID)
	if l > 0 {
		n += 1 + l + sovStorer(uint64(l))
	}
	return n
}

func (m *AcknowledgeResponse) Size() (n int) {
	var l int
	_ = l
	if m.Time != nil {
		l = types.SizeOfStdTime(*m.Time)
		n += 1 + l + sovStorer(uint64(l))
	}
	return n
}

func (m *ExtendTTLRequest) Size() (n int) {
	var l int
	_ = l
	if m.Stream != 0 {
		n += 1 + sovStorer(uint64(m.Stream))
	}
	l = len(m.EventID)
	if l > 0 {
		n += 1 + l + sovStorer(uint64(l))
	}
	if m.CurrentTTL != nil {
		l = types.SizeOfStdTime(*m.CurrentTTL)
		n += 1 + l + sovStorer(uint64(l))
	}
	if m.TTLTime != nil {
		l = types.SizeOfStdTime(*m.TTLTime)
		n += 1 + l + sovStorer(uint64(l))
	}
	return n
}

func (m *ExtendTTLResponse) Size() (n int) {
	var l int
	_ = l
	return n
}

func (m *QueryRequest) Size() (n int) {
	var l int
	_ = l
	if m.Stream != 0 {
		n += 1 + sovStorer(uint64(m.Stream))
	}
	if m.Limit != 0 {
		n += 1 + sovStorer(uint64(m.Limit))
	}
	if m.Query != nil {
		n += m.Query.Size()
	}
	return n
}

func (m *QueryRequest_Ttl) Size() (n int) {
	var l int
	_ = l
	if m.Ttl != nil {
		l = m.Ttl.Size()
		n += 1 + l + sovStorer(uint64(l))
	}
	return n
}
func (m *QueryTTLExpired) Size() (n int) {
	var l int
	_ = l
	return n
}

func sovStorer(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozStorer(x uint64) (n int) {
	return sovStorer(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AcknowledgeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStorer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AcknowledgeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AcknowledgeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stream", wireType)
			}
			m.Stream = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Stream |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EventID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStorer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EventID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStorer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStorer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AcknowledgeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStorer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AcknowledgeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AcknowledgeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStorer
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Time == nil {
				m.Time = new(time.Time)
			}
			if err := types.StdTimeUnmarshal(m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStorer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStorer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ExtendTTLRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStorer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExtendTTLRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExtendTTLRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stream", wireType)
			}
			m.Stream = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Stream |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EventID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStorer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EventID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentTTL", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStorer
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CurrentTTL == nil {
				m.CurrentTTL = new(time.Time)
			}
			if err := types.StdTimeUnmarshal(m.CurrentTTL, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TTLTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStorer
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TTLTime == nil {
				m.TTLTime = new(time.Time)
			}
			if err := types.StdTimeUnmarshal(m.TTLTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStorer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStorer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ExtendTTLResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStorer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExtendTTLResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExtendTTLResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipStorer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStorer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStorer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stream", wireType)
			}
			m.Stream = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Stream |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			m.Limit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Limit |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ttl", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStorer
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &QueryTTLExpired{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Query = &QueryRequest_Ttl{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStorer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStorer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryTTLExpired) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStorer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryTTLExpired: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTTLExpired: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipStorer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStorer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipStorer(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStorer
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStorer
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthStorer
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowStorer
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipStorer(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthStorer = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStorer   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("storer.proto", fileDescriptorStorer) }

var fileDescriptorStorer = []byte{
	// 459 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x4d, 0x6f, 0xd3, 0x40,
	0x14, 0xec, 0xd2, 0xb8, 0x51, 0x5f, 0xcb, 0x47, 0xb7, 0x15, 0xb2, 0x2c, 0xe4, 0x04, 0x73, 0xe9,
	0x05, 0x1b, 0xb5, 0x5c, 0x40, 0x42, 0x82, 0x40, 0x10, 0x88, 0x1c, 0xc0, 0x59, 0x71, 0x4f, 0x9c,
	0x57, 0x63, 0x61, 0x7b, 0xd3, 0xf5, 0x73, 0x5a, 0xf8, 0x15, 0x1c, 0xf9, 0x39, 0x88, 0x53, 0x8f,
	0xfc, 0x03, 0x50, 0xf2, 0x27, 0x38, 0xa2, 0xec, 0xda, 0x25, 0x0d, 0xa8, 0x44, 0xdc, 0x3c, 0xde,
	0x99, 0xd9, 0x9d, 0xf7, 0x06, 0xb6, 0x0b, 0x92, 0x0a, 0x95, 0x3f, 0x56, 0x92, 0x24, 0xbf, 0x86,
	0x93, 0x9c, 0x0a, 0x15, 0xf9, 0xe6, 0xaf, 0x73, 0x37, 0x4e, 0xe8, 0x5d, 0x39, 0xf4, 0x23, 0x99,
	0x05, 0xb1, 0x8c, 0x65, 0xa0, 0x69, 0xc3, 0xf2, 0x48, 0x23, 0x0d, 0xf4, 0x97, 0x91, 0x3b, 0xad,
	0x58, 0xca, 0x38, 0xc5, 0xdf, 0x2c, 0x4a, 0x32, 0x2c, 0x68, 0x90, 0x8d, 0x2b, 0xc2, 0x83, 0x05,
	0x3f, 0x8a, 0x8e, 0x4e, 0x82, 0xea, 0xbe, 0x20, 0xc9, 0x09, 0x55, 0x3e, 0x48, 0x03, 0x9c, 0x60,
	0x4e, 0xc6, 0xa0, 0x30, 0xc0, 0x48, 0xbd, 0xe7, 0xc0, 0x9f, 0x44, 0xef, 0x73, 0x79, 0x92, 0xe2,
	0x28, 0xc6, 0x10, 0x8f, 0x4b, 0x2c, 0x88, 0xdf, 0x84, 0x8d, 0x3e, 0x29, 0x1c, 0x64, 0x36, 0x6b,
	0xb3, 0x7d, 0x2b, 0xac, 0x10, 0xb7, 0xa1, 0xd9, 0x9d, 0x8b, 0x5f, 0x3e, 0xb3, 0xaf, 0xb4, 0xd9,
	0xfe, 0x66, 0x58, 0x43, 0xef, 0x15, 0xec, 0x5e, 0xf0, 0x29, 0xc6, 0x32, 0x2f, 0x90, 0xdf, 0x87,
	0x86, 0x48, 0x32, 0xd4, 0x36, 0x5b, 0x07, 0x8e, 0x6f, 0x92, 0xf8, 0x75, 0x12, 0x5f, 0xd4, 0x49,
	0x3a, 0x8d, 0x4f, 0xdf, 0x5b, 0x2c, 0xd4, 0x6c, 0xef, 0x2b, 0x83, 0x1b, 0xdd, 0x53, 0xc2, 0x7c,
	0x24, 0x44, 0xef, 0xbf, 0xdf, 0xc4, 0x1f, 0x03, 0x3c, 0x2d, 0x95, 0xc2, 0x9c, 0x84, 0xe8, 0xd9,
	0xeb, 0x2b, 0x3e, 0x61, 0x41, 0xc3, 0x1f, 0x42, 0x53, 0x88, 0x9e, 0x4e, 0xd0, 0x58, 0x51, 0x5e,
	0x0b, 0xbc, 0x5d, 0xd8, 0x59, 0xc8, 0x60, 0xe6, 0xe1, 0x7d, 0x84, 0xed, 0x37, 0x25, 0xaa, 0x0f,
	0xff, 0x0a, 0xb5, 0x07, 0x56, 0x9a, 0x64, 0x09, 0xe9, 0x48, 0x56, 0x68, 0x00, 0x3f, 0x84, 0x75,
	0xa2, 0xb4, 0x4a, 0xd2, 0xf2, 0x2f, 0xb6, 0xca, 0xd7, 0xc6, 0x42, 0xf4, 0xba, 0xa7, 0xe3, 0x44,
	0xe1, 0xe8, 0xc5, 0x5a, 0x38, 0x67, 0x77, 0x9a, 0x60, 0x1d, 0xcf, 0x4f, 0xbc, 0x1d, 0xb8, 0xbe,
	0x44, 0x39, 0xf8, 0xc9, 0xe0, 0x6a, 0x5f, 0xab, 0xfb, 0xa8, 0x26, 0x49, 0x84, 0xfc, 0x2d, 0x6c,
	0x2d, 0xec, 0x91, 0x7b, 0xcb, 0x97, 0xfc, 0x59, 0x16, 0xe7, 0xce, 0xa5, 0x9c, 0xaa, 0x08, 0xaf,
	0x61, 0xf3, 0x7c, 0x1a, 0xbc, 0xbd, 0xac, 0x58, 0x5e, 0xb6, 0x73, 0xfb, 0x12, 0x46, 0xe5, 0xf8,
	0x08, 0x2c, 0x1d, 0x87, 0xdf, 0xfa, 0xeb, 0x20, 0x6a, 0xa7, 0xbd, 0xf3, 0x53, 0x53, 0x7b, 0x5d,
	0x8e, 0x7b, 0xac, 0x63, 0x9f, 0x4d, 0x5d, 0xf6, 0x6d, 0xea, 0xb2, 0x1f, 0x53, 0x97, 0x7d, 0x9e,
	0xb9, 0x6b, 0x5f, 0x66, 0x2e, 0x3b, 0x9b, 0xb9, 0x6c, 0xb8, 0xa1, 0x77, 0x7b, 0xf8, 0x2b, 0x00,
	0x00, 0xff, 0xff, 0x8e, 0xb9, 0xe9, 0x08, 0xc4, 0x03, 0x00, 0x00,
}
