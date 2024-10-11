// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.26.1
// source: proto/telegram.proto

package protogen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Portfolio struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Portfolio) Reset() {
	*x = Portfolio{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_telegram_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Portfolio) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Portfolio) ProtoMessage() {}

func (x *Portfolio) ProtoReflect() protoreflect.Message {
	mi := &file_proto_telegram_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Portfolio.ProtoReflect.Descriptor instead.
func (*Portfolio) Descriptor() ([]byte, []int) {
	return file_proto_telegram_proto_rawDescGZIP(), []int{0}
}

func (x *Portfolio) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Portfolio) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type PortfolioListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatId string `protobuf:"bytes,1,opt,name=chatId,proto3" json:"chatId,omitempty"`
}

func (x *PortfolioListRequest) Reset() {
	*x = PortfolioListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_telegram_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortfolioListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortfolioListRequest) ProtoMessage() {}

func (x *PortfolioListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_telegram_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortfolioListRequest.ProtoReflect.Descriptor instead.
func (*PortfolioListRequest) Descriptor() ([]byte, []int) {
	return file_proto_telegram_proto_rawDescGZIP(), []int{1}
}

func (x *PortfolioListRequest) GetChatId() string {
	if x != nil {
		return x.ChatId
	}
	return ""
}

type PortfolioRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ChatId string `protobuf:"bytes,2,opt,name=chatId,proto3" json:"chatId,omitempty"`
}

func (x *PortfolioRequest) Reset() {
	*x = PortfolioRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_telegram_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortfolioRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortfolioRequest) ProtoMessage() {}

func (x *PortfolioRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_telegram_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortfolioRequest.ProtoReflect.Descriptor instead.
func (*PortfolioRequest) Descriptor() ([]byte, []int) {
	return file_proto_telegram_proto_rawDescGZIP(), []int{2}
}

func (x *PortfolioRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PortfolioRequest) GetChatId() string {
	if x != nil {
		return x.ChatId
	}
	return ""
}

type PortfolioListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Portfolios []*Portfolio `protobuf:"bytes,1,rep,name=portfolios,proto3" json:"portfolios,omitempty"`
}

func (x *PortfolioListResponse) Reset() {
	*x = PortfolioListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_telegram_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortfolioListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortfolioListResponse) ProtoMessage() {}

func (x *PortfolioListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_telegram_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortfolioListResponse.ProtoReflect.Descriptor instead.
func (*PortfolioListResponse) Descriptor() ([]byte, []int) {
	return file_proto_telegram_proto_rawDescGZIP(), []int{3}
}

func (x *PortfolioListResponse) GetPortfolios() []*Portfolio {
	if x != nil {
		return x.Portfolios
	}
	return nil
}

type PortfolioSummaryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *PortfolioSummaryResponse) Reset() {
	*x = PortfolioSummaryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_telegram_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortfolioSummaryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortfolioSummaryResponse) ProtoMessage() {}

func (x *PortfolioSummaryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_telegram_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortfolioSummaryResponse.ProtoReflect.Descriptor instead.
func (*PortfolioSummaryResponse) Descriptor() ([]byte, []int) {
	return file_proto_telegram_proto_rawDescGZIP(), []int{4}
}

func (x *PortfolioSummaryResponse) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

var File_proto_telegram_proto protoreflect.FileDescriptor

var file_proto_telegram_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x09, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f,
	0x6c, 0x69, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2e, 0x0a, 0x14, 0x50, 0x6f, 0x72, 0x74, 0x66,
	0x6f, 0x6c, 0x69, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x22, 0x3a, 0x0a, 0x10, 0x50, 0x6f, 0x72, 0x74, 0x66,
	0x6f, 0x6c, 0x69, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x63,
	0x68, 0x61, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x68, 0x61,
	0x74, 0x49, 0x64, 0x22, 0x43, 0x0a, 0x15, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x0a,
	0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x52, 0x0a, 0x70, 0x6f,
	0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x73, 0x22, 0x2e, 0x0a, 0x18, 0x50, 0x6f, 0x72, 0x74,
	0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x32, 0x97, 0x01, 0x0a, 0x06, 0x54, 0x65, 0x6c,
	0x65, 0x67, 0x61, 0x12, 0x41, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f,
	0x6c, 0x69, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x15, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f,
	0x6c, 0x69, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x72,
	0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x11, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f,
	0x6c, 0x69, 0x6f, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x0e, 0x5a, 0x0c, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67,
	0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_telegram_proto_rawDescOnce sync.Once
	file_proto_telegram_proto_rawDescData = file_proto_telegram_proto_rawDesc
)

func file_proto_telegram_proto_rawDescGZIP() []byte {
	file_proto_telegram_proto_rawDescOnce.Do(func() {
		file_proto_telegram_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_telegram_proto_rawDescData)
	})
	return file_proto_telegram_proto_rawDescData
}

var file_proto_telegram_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_telegram_proto_goTypes = []interface{}{
	(*Portfolio)(nil),                // 0: Portfolio
	(*PortfolioListRequest)(nil),     // 1: PortfolioListRequest
	(*PortfolioRequest)(nil),         // 2: PortfolioRequest
	(*PortfolioListResponse)(nil),    // 3: PortfolioListResponse
	(*PortfolioSummaryResponse)(nil), // 4: PortfolioSummaryResponse
}
var file_proto_telegram_proto_depIdxs = []int32{
	0, // 0: PortfolioListResponse.portfolios:type_name -> Portfolio
	1, // 1: Telega.GetPortfolioList:input_type -> PortfolioListRequest
	2, // 2: Telega.GetPortfolioSummaryMessage:input_type -> PortfolioRequest
	3, // 3: Telega.GetPortfolioList:output_type -> PortfolioListResponse
	4, // 4: Telega.GetPortfolioSummaryMessage:output_type -> PortfolioSummaryResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_telegram_proto_init() }
func file_proto_telegram_proto_init() {
	if File_proto_telegram_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_telegram_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Portfolio); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_telegram_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortfolioListRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_telegram_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortfolioRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_telegram_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortfolioListResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_telegram_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortfolioSummaryResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_telegram_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_telegram_proto_goTypes,
		DependencyIndexes: file_proto_telegram_proto_depIdxs,
		MessageInfos:      file_proto_telegram_proto_msgTypes,
	}.Build()
	File_proto_telegram_proto = out.File
	file_proto_telegram_proto_rawDesc = nil
	file_proto_telegram_proto_goTypes = nil
	file_proto_telegram_proto_depIdxs = nil
}