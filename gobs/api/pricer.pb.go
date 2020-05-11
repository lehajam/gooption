// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pricer.proto

package api_pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type PriceRequest struct {
	Pricingdate          float64  `protobuf:"fixed64,1,opt,name=pricingdate,proto3" json:"pricingdate,omitempty"`
	Strike               float64  `protobuf:"fixed64,2,opt,name=strike,proto3" json:"strike,omitempty"`
	Expiry               float64  `protobuf:"fixed64,3,opt,name=expiry,proto3" json:"expiry,omitempty"`
	PutCall              string   `protobuf:"bytes,4,opt,name=put_call,json=putCall,proto3" json:"put_call,omitempty"`
	Spot                 float64  `protobuf:"fixed64,5,opt,name=spot,proto3" json:"spot,omitempty"`
	Vol                  float64  `protobuf:"fixed64,6,opt,name=vol,proto3" json:"vol,omitempty"`
	Rate                 float64  `protobuf:"fixed64,7,opt,name=rate,proto3" json:"rate,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PriceRequest) Reset()         { *m = PriceRequest{} }
func (m *PriceRequest) String() string { return proto.CompactTextString(m) }
func (*PriceRequest) ProtoMessage()    {}
func (*PriceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fec020e99f204f1f, []int{0}
}

func (m *PriceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PriceRequest.Unmarshal(m, b)
}
func (m *PriceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PriceRequest.Marshal(b, m, deterministic)
}
func (m *PriceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PriceRequest.Merge(m, src)
}
func (m *PriceRequest) XXX_Size() int {
	return xxx_messageInfo_PriceRequest.Size(m)
}
func (m *PriceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PriceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PriceRequest proto.InternalMessageInfo

func (m *PriceRequest) GetPricingdate() float64 {
	if m != nil {
		return m.Pricingdate
	}
	return 0
}

func (m *PriceRequest) GetStrike() float64 {
	if m != nil {
		return m.Strike
	}
	return 0
}

func (m *PriceRequest) GetExpiry() float64 {
	if m != nil {
		return m.Expiry
	}
	return 0
}

func (m *PriceRequest) GetPutCall() string {
	if m != nil {
		return m.PutCall
	}
	return ""
}

func (m *PriceRequest) GetSpot() float64 {
	if m != nil {
		return m.Spot
	}
	return 0
}

func (m *PriceRequest) GetVol() float64 {
	if m != nil {
		return m.Vol
	}
	return 0
}

func (m *PriceRequest) GetRate() float64 {
	if m != nil {
		return m.Rate
	}
	return 0
}

type PriceResponse struct {
	Value                float64  `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
	ValueType            string   `protobuf:"bytes,2,opt,name=value_type,json=valueType,proto3" json:"value_type,omitempty"`
	Error                string   `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PriceResponse) Reset()         { *m = PriceResponse{} }
func (m *PriceResponse) String() string { return proto.CompactTextString(m) }
func (*PriceResponse) ProtoMessage()    {}
func (*PriceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fec020e99f204f1f, []int{1}
}

func (m *PriceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PriceResponse.Unmarshal(m, b)
}
func (m *PriceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PriceResponse.Marshal(b, m, deterministic)
}
func (m *PriceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PriceResponse.Merge(m, src)
}
func (m *PriceResponse) XXX_Size() int {
	return xxx_messageInfo_PriceResponse.Size(m)
}
func (m *PriceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PriceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PriceResponse proto.InternalMessageInfo

func (m *PriceResponse) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *PriceResponse) GetValueType() string {
	if m != nil {
		return m.ValueType
	}
	return ""
}

func (m *PriceResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*PriceRequest)(nil), "gobs.PriceRequest")
	proto.RegisterType((*PriceResponse)(nil), "gobs.PriceResponse")
}

func init() { proto.RegisterFile("pricer.proto", fileDescriptor_fec020e99f204f1f) }

var fileDescriptor_fec020e99f204f1f = []byte{
	// 316 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0xcf, 0x4a, 0x33, 0x31,
	0x14, 0xc5, 0x49, 0xff, 0x4c, 0xbf, 0xb9, 0xed, 0x87, 0xf6, 0x2a, 0x12, 0x8b, 0x42, 0x99, 0x55,
	0x71, 0xd1, 0x8a, 0xee, 0x74, 0xa7, 0x0b, 0xb7, 0x32, 0xba, 0x90, 0x6e, 0x4a, 0x5a, 0xc3, 0x10,
	0x0c, 0x93, 0x6b, 0x92, 0x29, 0x76, 0xeb, 0x2b, 0xf8, 0x2c, 0x3e, 0x89, 0xaf, 0xe0, 0x83, 0x48,
	0x92, 0x0a, 0x75, 0x77, 0xce, 0x2f, 0xe7, 0x86, 0xc3, 0xbd, 0x30, 0x20, 0xab, 0x56, 0xd2, 0x4e,
	0xc9, 0x1a, 0x6f, 0xb0, 0x53, 0x99, 0xa5, 0x1b, 0x9d, 0x54, 0xc6, 0x54, 0x5a, 0xce, 0x04, 0xa9,
	0x99, 0xa8, 0x6b, 0xe3, 0x85, 0x57, 0xa6, 0x76, 0x29, 0x53, 0x7c, 0x32, 0x18, 0xdc, 0x87, 0xa1,
	0x52, 0xbe, 0x36, 0xd2, 0x79, 0x1c, 0x43, 0x3f, 0x7c, 0xa2, 0xea, 0xea, 0x59, 0x78, 0xc9, 0xd9,
	0x98, 0x4d, 0x58, 0xb9, 0x8b, 0xf0, 0x08, 0x32, 0xe7, 0xad, 0x7a, 0x91, 0xbc, 0x15, 0x1f, 0xb7,
	0x2e, 0x70, 0xf9, 0x46, 0xca, 0x6e, 0x78, 0x3b, 0xf1, 0xe4, 0xf0, 0x18, 0xfe, 0x51, 0xe3, 0x17,
	0x2b, 0xa1, 0x35, 0xef, 0x8c, 0xd9, 0x24, 0x2f, 0x7b, 0xd4, 0xf8, 0x5b, 0xa1, 0x35, 0x22, 0x74,
	0x1c, 0x19, 0xcf, 0xbb, 0x71, 0x20, 0x6a, 0xdc, 0x87, 0xf6, 0xda, 0x68, 0x9e, 0x45, 0x14, 0x64,
	0x48, 0xd9, 0xd0, 0xa5, 0x97, 0x52, 0x41, 0x17, 0x73, 0xf8, 0xbf, 0xad, 0xed, 0xc8, 0xd4, 0x4e,
	0xe2, 0x21, 0x74, 0xd7, 0x42, 0x37, 0xbf, 0x8d, 0x93, 0xc1, 0x53, 0x80, 0x28, 0x16, 0x7e, 0x43,
	0xa9, 0x6f, 0x5e, 0xe6, 0x91, 0x3c, 0x6e, 0x28, 0x0e, 0x49, 0x6b, 0x8d, 0x8d, 0x8d, 0xf3, 0x32,
	0x99, 0x8b, 0xa7, 0xed, 0xdf, 0xf6, 0x41, 0xda, 0xb5, 0x5a, 0x49, 0xbc, 0x83, 0x6e, 0x04, 0x88,
	0xd3, 0xb0, 0xd2, 0xe9, 0xee, 0xc2, 0x46, 0x07, 0x7f, 0x58, 0x6a, 0x53, 0x0c, 0xdf, 0xbf, 0xbe,
	0x3f, 0x5a, 0xfd, 0x22, 0x9b, 0xc5, 0x8b, 0x5c, 0xb1, 0xb3, 0x09, 0x3b, 0x67, 0x37, 0xc3, 0xf9,
	0x5e, 0x08, 0x87, 0x5b, 0x5c, 0x0b, 0x52, 0x0b, 0x5a, 0x2e, 0xb3, 0x78, 0x87, 0xcb, 0x9f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xd1, 0x46, 0x9d, 0xb3, 0xbb, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PricerServiceClient is the client API for PricerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PricerServiceClient interface {
	Price(ctx context.Context, opts ...grpc.CallOption) (PricerService_PriceClient, error)
}

type pricerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPricerServiceClient(cc grpc.ClientConnInterface) PricerServiceClient {
	return &pricerServiceClient{cc}
}

func (c *pricerServiceClient) Price(ctx context.Context, opts ...grpc.CallOption) (PricerService_PriceClient, error) {
	stream, err := c.cc.NewStream(ctx, &_PricerService_serviceDesc.Streams[0], "/gobs.PricerService/Price", opts...)
	if err != nil {
		return nil, err
	}
	x := &pricerServicePriceClient{stream}
	return x, nil
}

type PricerService_PriceClient interface {
	Send(*PriceRequest) error
	Recv() (*PriceResponse, error)
	grpc.ClientStream
}

type pricerServicePriceClient struct {
	grpc.ClientStream
}

func (x *pricerServicePriceClient) Send(m *PriceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *pricerServicePriceClient) Recv() (*PriceResponse, error) {
	m := new(PriceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PricerServiceServer is the server API for PricerService service.
type PricerServiceServer interface {
	Price(PricerService_PriceServer) error
}

// UnimplementedPricerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPricerServiceServer struct {
}

func (*UnimplementedPricerServiceServer) Price(srv PricerService_PriceServer) error {
	return status.Errorf(codes.Unimplemented, "method Price not implemented")
}

func RegisterPricerServiceServer(s *grpc.Server, srv PricerServiceServer) {
	s.RegisterService(&_PricerService_serviceDesc, srv)
}

func _PricerService_Price_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PricerServiceServer).Price(&pricerServicePriceServer{stream})
}

type PricerService_PriceServer interface {
	Send(*PriceResponse) error
	Recv() (*PriceRequest, error)
	grpc.ServerStream
}

type pricerServicePriceServer struct {
	grpc.ServerStream
}

func (x *pricerServicePriceServer) Send(m *PriceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *pricerServicePriceServer) Recv() (*PriceRequest, error) {
	m := new(PriceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _PricerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gobs.PricerService",
	HandlerType: (*PricerServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Price",
			Handler:       _PricerService_Price_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pricer.proto",
}
