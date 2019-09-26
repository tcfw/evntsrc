// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: emails.proto

/*
	Package evntsrc_emails is a generated protocol buffer package.

	It is generated from these files:
		emails.proto

	It has these top-level messages:
		Recipient
		Email
		Attachment
		EmailResponse
*/
package evntsrc_emails

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

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

type Recipient struct {
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *Recipient) Reset()                    { *m = Recipient{} }
func (m *Recipient) String() string            { return proto.CompactTextString(m) }
func (*Recipient) ProtoMessage()               {}
func (*Recipient) Descriptor() ([]byte, []int) { return fileDescriptorEmails, []int{0} }

func (m *Recipient) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Recipient) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (*Recipient) XXX_MessageName() string {
	return "evntsrc.emails.Recipient"
}

type Email struct {
	From        string            `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To          []*Recipient      `protobuf:"bytes,2,rep,name=to" json:"to,omitempty"`
	Subject     string            `protobuf:"bytes,3,opt,name=subject,proto3" json:"subject,omitempty"`
	PlainText   string            `protobuf:"bytes,4,opt,name=plain_text,json=plainText,proto3" json:"plain_text,omitempty"`
	Html        string            `protobuf:"bytes,5,opt,name=html,proto3" json:"html,omitempty"`
	Attachments []*Attachment     `protobuf:"bytes,6,rep,name=attachments" json:"attachments,omitempty"`
	Headers     map[string]string `protobuf:"bytes,7,rep,name=headers" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *Email) Reset()                    { *m = Email{} }
func (m *Email) String() string            { return proto.CompactTextString(m) }
func (*Email) ProtoMessage()               {}
func (*Email) Descriptor() ([]byte, []int) { return fileDescriptorEmails, []int{1} }

func (m *Email) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *Email) GetTo() []*Recipient {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Email) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Email) GetPlainText() string {
	if m != nil {
		return m.PlainText
	}
	return ""
}

func (m *Email) GetHtml() string {
	if m != nil {
		return m.Html
	}
	return ""
}

func (m *Email) GetAttachments() []*Attachment {
	if m != nil {
		return m.Attachments
	}
	return nil
}

func (m *Email) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (*Email) XXX_MessageName() string {
	return "evntsrc.emails.Email"
}

type Attachment struct {
	Filename string `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	// Types that are valid to be assigned to Type:
	//	*Attachment_Data
	//	*Attachment_Uri
	Type isAttachment_Type `protobuf_oneof:"type"`
}

func (m *Attachment) Reset()                    { *m = Attachment{} }
func (m *Attachment) String() string            { return proto.CompactTextString(m) }
func (*Attachment) ProtoMessage()               {}
func (*Attachment) Descriptor() ([]byte, []int) { return fileDescriptorEmails, []int{2} }

type isAttachment_Type interface {
	isAttachment_Type()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Attachment_Data struct {
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3,oneof"`
}
type Attachment_Uri struct {
	Uri string `protobuf:"bytes,3,opt,name=uri,proto3,oneof"`
}

func (*Attachment_Data) isAttachment_Type() {}
func (*Attachment_Uri) isAttachment_Type()  {}

func (m *Attachment) GetType() isAttachment_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Attachment) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *Attachment) GetData() []byte {
	if x, ok := m.GetType().(*Attachment_Data); ok {
		return x.Data
	}
	return nil
}

func (m *Attachment) GetUri() string {
	if x, ok := m.GetType().(*Attachment_Uri); ok {
		return x.Uri
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Attachment) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Attachment_OneofMarshaler, _Attachment_OneofUnmarshaler, _Attachment_OneofSizer, []interface{}{
		(*Attachment_Data)(nil),
		(*Attachment_Uri)(nil),
	}
}

func _Attachment_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Attachment)
	// type
	switch x := m.Type.(type) {
	case *Attachment_Data:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		_ = b.EncodeRawBytes(x.Data)
	case *Attachment_Uri:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.Uri)
	case nil:
	default:
		return fmt.Errorf("Attachment.Type has unexpected type %T", x)
	}
	return nil
}

func _Attachment_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Attachment)
	switch tag {
	case 2: // type.data
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.Type = &Attachment_Data{x}
		return true, err
	case 3: // type.uri
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Type = &Attachment_Uri{x}
		return true, err
	default:
		return false, nil
	}
}

func _Attachment_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Attachment)
	// type
	switch x := m.Type.(type) {
	case *Attachment_Data:
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Data)))
		n += len(x.Data)
	case *Attachment_Uri:
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Uri)))
		n += len(x.Uri)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func (*Attachment) XXX_MessageName() string {
	return "evntsrc.emails.Attachment"
}

type EmailResponse struct {
}

func (m *EmailResponse) Reset()                    { *m = EmailResponse{} }
func (m *EmailResponse) String() string            { return proto.CompactTextString(m) }
func (*EmailResponse) ProtoMessage()               {}
func (*EmailResponse) Descriptor() ([]byte, []int) { return fileDescriptorEmails, []int{3} }

func (*EmailResponse) XXX_MessageName() string {
	return "evntsrc.emails.EmailResponse"
}
func init() {
	proto.RegisterType((*Recipient)(nil), "evntsrc.emails.Recipient")
	proto.RegisterType((*Email)(nil), "evntsrc.emails.Email")
	proto.RegisterType((*Attachment)(nil), "evntsrc.emails.Attachment")
	proto.RegisterType((*EmailResponse)(nil), "evntsrc.emails.EmailResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EmailService service

type EmailServiceClient interface {
	Send(ctx context.Context, in *Email, opts ...grpc.CallOption) (*EmailResponse, error)
	SendRaw(ctx context.Context, in *Email, opts ...grpc.CallOption) (*EmailResponse, error)
}

type emailServiceClient struct {
	cc *grpc.ClientConn
}

func NewEmailServiceClient(cc *grpc.ClientConn) EmailServiceClient {
	return &emailServiceClient{cc}
}

func (c *emailServiceClient) Send(ctx context.Context, in *Email, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := grpc.Invoke(ctx, "/evntsrc.emails.EmailService/Send", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailServiceClient) SendRaw(ctx context.Context, in *Email, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := grpc.Invoke(ctx, "/evntsrc.emails.EmailService/SendRaw", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EmailService service

type EmailServiceServer interface {
	Send(context.Context, *Email) (*EmailResponse, error)
	SendRaw(context.Context, *Email) (*EmailResponse, error)
}

func RegisterEmailServiceServer(s *grpc.Server, srv EmailServiceServer) {
	s.RegisterService(&_EmailService_serviceDesc, srv)
}

func _EmailService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/evntsrc.emails.EmailService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).Send(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailService_SendRaw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).SendRaw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/evntsrc.emails.EmailService/SendRaw",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).SendRaw(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

var _EmailService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "evntsrc.emails.EmailService",
	HandlerType: (*EmailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _EmailService_Send_Handler,
		},
		{
			MethodName: "SendRaw",
			Handler:    _EmailService_SendRaw_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "emails.proto",
}

func (m *Recipient) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Recipient) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Email) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintEmails(dAtA, i, uint64(len(m.Email)))
		i += copy(dAtA[i:], m.Email)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintEmails(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	return i, nil
}

func (m *Email) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Email) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.From) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintEmails(dAtA, i, uint64(len(m.From)))
		i += copy(dAtA[i:], m.From)
	}
	if len(m.To) > 0 {
		for _, msg := range m.To {
			dAtA[i] = 0x12
			i++
			i = encodeVarintEmails(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Subject) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintEmails(dAtA, i, uint64(len(m.Subject)))
		i += copy(dAtA[i:], m.Subject)
	}
	if len(m.PlainText) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintEmails(dAtA, i, uint64(len(m.PlainText)))
		i += copy(dAtA[i:], m.PlainText)
	}
	if len(m.Html) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintEmails(dAtA, i, uint64(len(m.Html)))
		i += copy(dAtA[i:], m.Html)
	}
	if len(m.Attachments) > 0 {
		for _, msg := range m.Attachments {
			dAtA[i] = 0x32
			i++
			i = encodeVarintEmails(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Headers) > 0 {
		for k, _ := range m.Headers {
			dAtA[i] = 0x3a
			i++
			v := m.Headers[k]
			mapSize := 1 + len(k) + sovEmails(uint64(len(k))) + 1 + len(v) + sovEmails(uint64(len(v)))
			i = encodeVarintEmails(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintEmails(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintEmails(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	return i, nil
}

func (m *Attachment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Attachment) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Filename) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintEmails(dAtA, i, uint64(len(m.Filename)))
		i += copy(dAtA[i:], m.Filename)
	}
	if m.Type != nil {
		nn1, err := m.Type.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	return i, nil
}

func (m *Attachment_Data) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Data != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintEmails(dAtA, i, uint64(len(m.Data)))
		i += copy(dAtA[i:], m.Data)
	}
	return i, nil
}
func (m *Attachment_Uri) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	dAtA[i] = 0x1a
	i++
	i = encodeVarintEmails(dAtA, i, uint64(len(m.Uri)))
	i += copy(dAtA[i:], m.Uri)
	return i, nil
}
func (m *EmailResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EmailResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func encodeVarintEmails(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Recipient) Size() (n int) {
	var l int
	_ = l
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovEmails(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovEmails(uint64(l))
	}
	return n
}

func (m *Email) Size() (n int) {
	var l int
	_ = l
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovEmails(uint64(l))
	}
	if len(m.To) > 0 {
		for _, e := range m.To {
			l = e.Size()
			n += 1 + l + sovEmails(uint64(l))
		}
	}
	l = len(m.Subject)
	if l > 0 {
		n += 1 + l + sovEmails(uint64(l))
	}
	l = len(m.PlainText)
	if l > 0 {
		n += 1 + l + sovEmails(uint64(l))
	}
	l = len(m.Html)
	if l > 0 {
		n += 1 + l + sovEmails(uint64(l))
	}
	if len(m.Attachments) > 0 {
		for _, e := range m.Attachments {
			l = e.Size()
			n += 1 + l + sovEmails(uint64(l))
		}
	}
	if len(m.Headers) > 0 {
		for k, v := range m.Headers {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovEmails(uint64(len(k))) + 1 + len(v) + sovEmails(uint64(len(v)))
			n += mapEntrySize + 1 + sovEmails(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *Attachment) Size() (n int) {
	var l int
	_ = l
	l = len(m.Filename)
	if l > 0 {
		n += 1 + l + sovEmails(uint64(l))
	}
	if m.Type != nil {
		n += m.Type.Size()
	}
	return n
}

func (m *Attachment_Data) Size() (n int) {
	var l int
	_ = l
	if m.Data != nil {
		l = len(m.Data)
		n += 1 + l + sovEmails(uint64(l))
	}
	return n
}
func (m *Attachment_Uri) Size() (n int) {
	var l int
	_ = l
	l = len(m.Uri)
	n += 1 + l + sovEmails(uint64(l))
	return n
}
func (m *EmailResponse) Size() (n int) {
	var l int
	_ = l
	return n
}

func sovEmails(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozEmails(x uint64) (n int) {
	return sovEmails(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Recipient) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEmails
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
			return fmt.Errorf("proto: Recipient: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Recipient: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Email", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Email = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEmails(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEmails
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
func (m *Email) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEmails
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
			return fmt.Errorf("proto: Email: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Email: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field To", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.To = append(m.To, &Recipient{})
			if err := m.To[len(m.To)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subject", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Subject = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PlainText", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PlainText = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Html", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Html = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attachments", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Attachments = append(m.Attachments, &Attachment{})
			if err := m.Attachments[len(m.Attachments)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Headers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Headers == nil {
				m.Headers = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowEmails
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowEmails
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthEmails
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowEmails
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthEmails
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipEmails(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthEmails
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Headers[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEmails(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEmails
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
func (m *Attachment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEmails
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
			return fmt.Errorf("proto: Attachment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Attachment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filename", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Filename = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := make([]byte, postIndex-iNdEx)
			copy(v, dAtA[iNdEx:postIndex])
			m.Type = &Attachment_Data{v}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uri", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEmails
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
				return ErrInvalidLengthEmails
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = &Attachment_Uri{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEmails(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEmails
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
func (m *EmailResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEmails
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
			return fmt.Errorf("proto: EmailResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EmailResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipEmails(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEmails
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
func skipEmails(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEmails
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
					return 0, ErrIntOverflowEmails
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
					return 0, ErrIntOverflowEmails
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
				return 0, ErrInvalidLengthEmails
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowEmails
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
				next, err := skipEmails(dAtA[start:])
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
	ErrInvalidLengthEmails = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEmails   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("emails.proto", fileDescriptorEmails) }

var fileDescriptorEmails = []byte{
	// 432 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x4d, 0x6b, 0xdb, 0x40,
	0x10, 0xf5, 0xca, 0xb2, 0x5d, 0x4f, 0xdc, 0x0f, 0x96, 0x14, 0xb6, 0x82, 0x88, 0xa0, 0x53, 0x7a,
	0xa8, 0x02, 0x29, 0x85, 0x12, 0x42, 0x21, 0x81, 0x80, 0xcf, 0x9b, 0x9e, 0x7a, 0x29, 0x6b, 0x79,
	0x6c, 0xa9, 0x95, 0xb4, 0x42, 0x1a, 0xb9, 0xf1, 0xcf, 0xe8, 0xad, 0x3f, 0xa7, 0xc7, 0x40, 0x2f,
	0xfd, 0x09, 0xc5, 0xfe, 0x23, 0x45, 0xa3, 0xc8, 0x4d, 0x43, 0x4e, 0xb9, 0xcd, 0x9b, 0x99, 0xf7,
	0xe6, 0xed, 0x63, 0x61, 0x82, 0x99, 0x49, 0xd2, 0x2a, 0x2c, 0x4a, 0x4b, 0x56, 0x3e, 0xc3, 0x55,
	0x4e, 0x55, 0x19, 0x85, 0x6d, 0xd7, 0x7b, 0xb3, 0x4c, 0x28, 0xae, 0x67, 0x61, 0x64, 0xb3, 0xe3,
	0xa5, 0x5d, 0xda, 0x63, 0x5e, 0x9b, 0xd5, 0x0b, 0x46, 0x0c, 0xb8, 0x6a, 0xe9, 0xc1, 0x3b, 0x18,
	0x6b, 0x8c, 0x92, 0x22, 0xc1, 0x9c, 0xe4, 0x3e, 0x0c, 0x58, 0x45, 0x89, 0x43, 0x71, 0x34, 0xd6,
	0x2d, 0x90, 0x12, 0xdc, 0xdc, 0x64, 0xa8, 0x1c, 0x6e, 0x72, 0x1d, 0xfc, 0x72, 0x60, 0x70, 0xd9,
	0x4d, 0x17, 0xa5, 0xcd, 0x6e, 0x29, 0x5c, 0xcb, 0xd7, 0xe0, 0x90, 0x55, 0xce, 0x61, 0xff, 0x68,
	0xef, 0xe4, 0x55, 0xf8, 0xbf, 0xc1, 0x70, 0x77, 0x4e, 0x3b, 0x64, 0xa5, 0x82, 0x51, 0x55, 0xcf,
	0xbe, 0x60, 0x44, 0xaa, 0xcf, 0x0a, 0x1d, 0x94, 0x07, 0x00, 0x45, 0x6a, 0x92, 0xfc, 0x33, 0xe1,
	0x35, 0x29, 0x97, 0x87, 0x63, 0xee, 0x7c, 0xc4, 0x6b, 0x6a, 0xee, 0xc6, 0x94, 0xa5, 0x6a, 0xd0,
	0xde, 0x6d, 0x6a, 0x79, 0x06, 0x7b, 0x86, 0xc8, 0x44, 0x71, 0x86, 0x39, 0x55, 0x6a, 0xc8, 0x06,
	0xbc, 0xfb, 0x06, 0xce, 0x77, 0x2b, 0xfa, 0xee, 0xba, 0x3c, 0x83, 0x51, 0x8c, 0x66, 0x8e, 0x65,
	0xa5, 0x46, 0xcc, 0x0c, 0xee, 0x33, 0xf9, 0xc5, 0xe1, 0xb4, 0x5d, 0xba, 0xcc, 0xa9, 0x5c, 0xeb,
	0x8e, 0xe2, 0x9d, 0xc2, 0xe4, 0xee, 0x40, 0xbe, 0x80, 0xfe, 0x57, 0x5c, 0xdf, 0xc6, 0xd2, 0x94,
	0x4d, 0xba, 0x2b, 0x93, 0xd6, 0x5d, 0x90, 0x2d, 0x38, 0x75, 0xde, 0x8b, 0xe0, 0x13, 0xc0, 0x3f,
	0x53, 0xd2, 0x83, 0x27, 0x8b, 0x24, 0x45, 0xce, 0xbc, 0xa5, 0xef, 0xb0, 0xdc, 0x07, 0x77, 0x6e,
	0xc8, 0xb0, 0xc4, 0x64, 0xda, 0xd3, 0x8c, 0xa4, 0x84, 0x7e, 0x5d, 0x26, 0x6d, 0x80, 0xd3, 0x9e,
	0x6e, 0xc0, 0xc5, 0x10, 0x5c, 0x5a, 0x17, 0x18, 0x3c, 0x87, 0xa7, 0x6c, 0x5b, 0x63, 0x55, 0xd8,
	0xbc, 0xc2, 0x93, 0xef, 0x02, 0x26, 0xdc, 0xb9, 0xc2, 0x72, 0x95, 0x44, 0x28, 0x3f, 0x80, 0x7b,
	0x85, 0xf9, 0x5c, 0xbe, 0x7c, 0xf0, 0xb9, 0xde, 0xc1, 0x83, 0xed, 0x4e, 0x2e, 0xe8, 0xc9, 0x73,
	0x18, 0x35, 0x7c, 0x6d, 0xbe, 0x3d, 0x56, 0xe2, 0x42, 0xdd, 0x6c, 0x7c, 0xf1, 0x7b, 0xe3, 0x8b,
	0x3f, 0x1b, 0x5f, 0xfc, 0xd8, 0xfa, 0xbd, 0x9f, 0x5b, 0x5f, 0xdc, 0x6c, 0x7d, 0x31, 0x1b, 0xf2,
	0x37, 0x7d, 0xfb, 0x37, 0x00, 0x00, 0xff, 0xff, 0x1a, 0xd4, 0x8b, 0x59, 0xf5, 0x02, 0x00, 0x00,
}
