// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bridge.proto

/*
	Package evntsrc_bridge is a generated protocol buffer package.

	It is generated from these files:
		bridge.proto

	It has these top-level messages:
		PublishRequest
		SubscribeRequest
		GeneralResponse
*/
package evntsrc_bridge

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import evntsrc_event "github.com/tcfw/evntsrc/pkg/event/protos"

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PublishRequest struct {
	Stream int32                `protobuf:"varint,1,opt,name=Stream,proto3" json:"Stream,omitempty"`
	Event  *evntsrc_event.Event `protobuf:"bytes,2,opt,name=event" json:"event,omitempty"`
}

func (m *PublishRequest) Reset()                    { *m = PublishRequest{} }
func (m *PublishRequest) String() string            { return proto.CompactTextString(m) }
func (*PublishRequest) ProtoMessage()               {}
func (*PublishRequest) Descriptor() ([]byte, []int) { return fileDescriptorBridge, []int{0} }

func (m *PublishRequest) GetStream() int32 {
	if m != nil {
		return m.Stream
	}
	return 0
}

func (m *PublishRequest) GetEvent() *evntsrc_event.Event {
	if m != nil {
		return m.Event
	}
	return nil
}

type SubscribeRequest struct {
	Stream  int32  `protobuf:"varint,1,opt,name=Stream,proto3" json:"Stream,omitempty"`
	Channel string `protobuf:"bytes,2,opt,name=Channel,proto3" json:"Channel,omitempty"`
}

func (m *SubscribeRequest) Reset()                    { *m = SubscribeRequest{} }
func (m *SubscribeRequest) String() string            { return proto.CompactTextString(m) }
func (*SubscribeRequest) ProtoMessage()               {}
func (*SubscribeRequest) Descriptor() ([]byte, []int) { return fileDescriptorBridge, []int{1} }

func (m *SubscribeRequest) GetStream() int32 {
	if m != nil {
		return m.Stream
	}
	return 0
}

func (m *SubscribeRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

type GeneralResponse struct {
}

func (m *GeneralResponse) Reset()                    { *m = GeneralResponse{} }
func (m *GeneralResponse) String() string            { return proto.CompactTextString(m) }
func (*GeneralResponse) ProtoMessage()               {}
func (*GeneralResponse) Descriptor() ([]byte, []int) { return fileDescriptorBridge, []int{2} }

func init() {
	proto.RegisterType((*PublishRequest)(nil), "evntsrc.bridge.PublishRequest")
	proto.RegisterType((*SubscribeRequest)(nil), "evntsrc.bridge.SubscribeRequest")
	proto.RegisterType((*GeneralResponse)(nil), "evntsrc.bridge.GeneralResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for BridgeService service

type BridgeServiceClient interface {
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*GeneralResponse, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (BridgeService_SubscribeClient, error)
}

type bridgeServiceClient struct {
	cc *grpc.ClientConn
}

func NewBridgeServiceClient(cc *grpc.ClientConn) BridgeServiceClient {
	return &bridgeServiceClient{cc}
}

func (c *bridgeServiceClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := grpc.Invoke(ctx, "/evntsrc.bridge.BridgeService/Publish", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bridgeServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (BridgeService_SubscribeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_BridgeService_serviceDesc.Streams[0], c.cc, "/evntsrc.bridge.BridgeService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &bridgeServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BridgeService_SubscribeClient interface {
	Recv() (*evntsrc_event.Event, error)
	grpc.ClientStream
}

type bridgeServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *bridgeServiceSubscribeClient) Recv() (*evntsrc_event.Event, error) {
	m := new(evntsrc_event.Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for BridgeService service

type BridgeServiceServer interface {
	Publish(context.Context, *PublishRequest) (*GeneralResponse, error)
	Subscribe(*SubscribeRequest, BridgeService_SubscribeServer) error
}

func RegisterBridgeServiceServer(s *grpc.Server, srv BridgeServiceServer) {
	s.RegisterService(&_BridgeService_serviceDesc, srv)
}

func _BridgeService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BridgeServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/evntsrc.bridge.BridgeService/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BridgeServiceServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BridgeService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BridgeServiceServer).Subscribe(m, &bridgeServiceSubscribeServer{stream})
}

type BridgeService_SubscribeServer interface {
	Send(*evntsrc_event.Event) error
	grpc.ServerStream
}

type bridgeServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *bridgeServiceSubscribeServer) Send(m *evntsrc_event.Event) error {
	return x.ServerStream.SendMsg(m)
}

var _BridgeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "evntsrc.bridge.BridgeService",
	HandlerType: (*BridgeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _BridgeService_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _BridgeService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "bridge.proto",
}

func (m *PublishRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PublishRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Stream != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintBridge(dAtA, i, uint64(m.Stream))
	}
	if m.Event != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintBridge(dAtA, i, uint64(m.Event.Size()))
		n1, err := m.Event.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *SubscribeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubscribeRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Stream != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintBridge(dAtA, i, uint64(m.Stream))
	}
	if len(m.Channel) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintBridge(dAtA, i, uint64(len(m.Channel)))
		i += copy(dAtA[i:], m.Channel)
	}
	return i, nil
}

func (m *GeneralResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GeneralResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func encodeVarintBridge(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *PublishRequest) Size() (n int) {
	var l int
	_ = l
	if m.Stream != 0 {
		n += 1 + sovBridge(uint64(m.Stream))
	}
	if m.Event != nil {
		l = m.Event.Size()
		n += 1 + l + sovBridge(uint64(l))
	}
	return n
}

func (m *SubscribeRequest) Size() (n int) {
	var l int
	_ = l
	if m.Stream != 0 {
		n += 1 + sovBridge(uint64(m.Stream))
	}
	l = len(m.Channel)
	if l > 0 {
		n += 1 + l + sovBridge(uint64(l))
	}
	return n
}

func (m *GeneralResponse) Size() (n int) {
	var l int
	_ = l
	return n
}

func sovBridge(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozBridge(x uint64) (n int) {
	return sovBridge(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PublishRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridge
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
			return fmt.Errorf("proto: PublishRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PublishRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stream", wireType)
			}
			m.Stream = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return fmt.Errorf("proto: wrong wireType = %d for field Event", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Event == nil {
				m.Event = &evntsrc_event.Event{}
			}
			if err := m.Event.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBridge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBridge
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
func (m *SubscribeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridge
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
			return fmt.Errorf("proto: SubscribeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubscribeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stream", wireType)
			}
			m.Stream = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return fmt.Errorf("proto: wrong wireType = %d for field Channel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Channel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBridge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBridge
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
func (m *GeneralResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridge
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
			return fmt.Errorf("proto: GeneralResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GeneralResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipBridge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBridge
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
func skipBridge(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBridge
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
					return 0, ErrIntOverflowBridge
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
					return 0, ErrIntOverflowBridge
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
				return 0, ErrInvalidLengthBridge
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowBridge
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
				next, err := skipBridge(dAtA[start:])
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
	ErrInvalidLengthBridge = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBridge   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("bridge.proto", fileDescriptorBridge) }

var fileDescriptorBridge = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x2a, 0xca, 0x4c,
	0x49, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4b, 0x2d, 0xcb, 0x2b, 0x29, 0x2e,
	0x4a, 0xd6, 0x83, 0x88, 0x4a, 0x99, 0xa4, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7,
	0xea, 0x97, 0x24, 0xa7, 0x95, 0xeb, 0x43, 0xe5, 0xf5, 0x0b, 0xb2, 0xd3, 0xf5, 0x53, 0xcb, 0x52,
	0xf3, 0x4a, 0xf4, 0xc1, 0x1a, 0x8b, 0x21, 0x1c, 0x88, 0x29, 0x4a, 0x21, 0x5c, 0x7c, 0x01, 0xa5,
	0x49, 0x39, 0x99, 0xc5, 0x19, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x62, 0x5c, 0x6c,
	0xc1, 0x25, 0x45, 0xa9, 0x89, 0xb9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xac, 0x41, 0x50, 0x9e, 0x90,
	0x16, 0x17, 0x2b, 0x58, 0xa3, 0x04, 0x93, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x88, 0x1e, 0xcc, 0x7e,
	0x88, 0x71, 0xae, 0x20, 0x32, 0x08, 0xa2, 0x44, 0xc9, 0x85, 0x4b, 0x20, 0xb8, 0x34, 0xa9, 0x38,
	0xb9, 0x28, 0x33, 0x29, 0x95, 0x90, 0xb9, 0x12, 0x5c, 0xec, 0xce, 0x19, 0x89, 0x79, 0x79, 0xa9,
	0x39, 0x60, 0x93, 0x39, 0x83, 0x60, 0x5c, 0x25, 0x41, 0x2e, 0x7e, 0xf7, 0xd4, 0xbc, 0xd4, 0xa2,
	0xc4, 0x9c, 0xa0, 0xd4, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0xa3, 0x65, 0x8c, 0x5c, 0xbc, 0x4e,
	0x60, 0xff, 0x06, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0x0a, 0xf9, 0x70, 0xb1, 0x43, 0x3d, 0x20,
	0x24, 0xa7, 0x87, 0x1a, 0x24, 0x7a, 0xa8, 0x3e, 0x93, 0x92, 0x47, 0x97, 0x47, 0x33, 0x5d, 0x89,
	0x41, 0xc8, 0x9d, 0x8b, 0x13, 0xee, 0x70, 0x21, 0x05, 0x74, 0xf5, 0xe8, 0x7e, 0x92, 0xc2, 0x1a,
	0x08, 0x4a, 0x0c, 0x06, 0x8c, 0x4e, 0x02, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8,
	0xe0, 0x91, 0x1c, 0xe3, 0x8c, 0xc7, 0x72, 0x0c, 0x49, 0x6c, 0xe0, 0x00, 0x37, 0x06, 0x04, 0x00,
	0x00, 0xff, 0xff, 0xa2, 0x93, 0x61, 0xdf, 0xc6, 0x01, 0x00, 0x00,
}
