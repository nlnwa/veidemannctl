// Code generated by protoc-gen-go. DO NOT EDIT.
// source: report.proto

package veidemann_api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Filter_Operator int32

const (
	Filter_EQ    Filter_Operator = 0
	Filter_NE    Filter_Operator = 1
	Filter_MATCH Filter_Operator = 2
	Filter_LT    Filter_Operator = 3
	Filter_GT    Filter_Operator = 4
)

var Filter_Operator_name = map[int32]string{
	0: "EQ",
	1: "NE",
	2: "MATCH",
	3: "LT",
	4: "GT",
}
var Filter_Operator_value = map[string]int32{
	"EQ":    0,
	"NE":    1,
	"MATCH": 2,
	"LT":    3,
	"GT":    4,
}

func (x Filter_Operator) String() string {
	return proto.EnumName(Filter_Operator_name, int32(x))
}
func (Filter_Operator) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{0, 0}
}

type Filter struct {
	FieldName            string          `protobuf:"bytes,1,opt,name=field_name,json=fieldName,proto3" json:"field_name,omitempty"`
	Op                   Filter_Operator `protobuf:"varint,2,opt,name=op,proto3,enum=veidemann.api.Filter_Operator" json:"op,omitempty"`
	Value                string          `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Filter) Reset()         { *m = Filter{} }
func (m *Filter) String() string { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()    {}
func (*Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{0}
}
func (m *Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Filter.Unmarshal(m, b)
}
func (m *Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Filter.Marshal(b, m, deterministic)
}
func (dst *Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Filter.Merge(dst, src)
}
func (m *Filter) XXX_Size() int {
	return xxx_messageInfo_Filter.Size(m)
}
func (m *Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Filter proto.InternalMessageInfo

func (m *Filter) GetFieldName() string {
	if m != nil {
		return m.FieldName
	}
	return ""
}

func (m *Filter) GetOp() Filter_Operator {
	if m != nil {
		return m.Op
	}
	return Filter_EQ
}

func (m *Filter) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type CrawlLogListRequest struct {
	WarcId               []string  `protobuf:"bytes,1,rep,name=warc_id,json=warcId,proto3" json:"warc_id,omitempty"`
	ExecutionId          string    `protobuf:"bytes,2,opt,name=execution_id,json=executionId,proto3" json:"execution_id,omitempty"`
	Filter               []*Filter `protobuf:"bytes,3,rep,name=filter,proto3" json:"filter,omitempty"`
	PageSize             int32     `protobuf:"varint,14,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Page                 int32     `protobuf:"varint,15,opt,name=page,proto3" json:"page,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CrawlLogListRequest) Reset()         { *m = CrawlLogListRequest{} }
func (m *CrawlLogListRequest) String() string { return proto.CompactTextString(m) }
func (*CrawlLogListRequest) ProtoMessage()    {}
func (*CrawlLogListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{1}
}
func (m *CrawlLogListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrawlLogListRequest.Unmarshal(m, b)
}
func (m *CrawlLogListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrawlLogListRequest.Marshal(b, m, deterministic)
}
func (dst *CrawlLogListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrawlLogListRequest.Merge(dst, src)
}
func (m *CrawlLogListRequest) XXX_Size() int {
	return xxx_messageInfo_CrawlLogListRequest.Size(m)
}
func (m *CrawlLogListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CrawlLogListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CrawlLogListRequest proto.InternalMessageInfo

func (m *CrawlLogListRequest) GetWarcId() []string {
	if m != nil {
		return m.WarcId
	}
	return nil
}

func (m *CrawlLogListRequest) GetExecutionId() string {
	if m != nil {
		return m.ExecutionId
	}
	return ""
}

func (m *CrawlLogListRequest) GetFilter() []*Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *CrawlLogListRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *CrawlLogListRequest) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

type CrawlLogListReply struct {
	Value                []*CrawlLog `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
	Count                int64       `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	PageSize             int32       `protobuf:"varint,14,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Page                 int32       `protobuf:"varint,15,opt,name=page,proto3" json:"page,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *CrawlLogListReply) Reset()         { *m = CrawlLogListReply{} }
func (m *CrawlLogListReply) String() string { return proto.CompactTextString(m) }
func (*CrawlLogListReply) ProtoMessage()    {}
func (*CrawlLogListReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{2}
}
func (m *CrawlLogListReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrawlLogListReply.Unmarshal(m, b)
}
func (m *CrawlLogListReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrawlLogListReply.Marshal(b, m, deterministic)
}
func (dst *CrawlLogListReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrawlLogListReply.Merge(dst, src)
}
func (m *CrawlLogListReply) XXX_Size() int {
	return xxx_messageInfo_CrawlLogListReply.Size(m)
}
func (m *CrawlLogListReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CrawlLogListReply.DiscardUnknown(m)
}

var xxx_messageInfo_CrawlLogListReply proto.InternalMessageInfo

func (m *CrawlLogListReply) GetValue() []*CrawlLog {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *CrawlLogListReply) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *CrawlLogListReply) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *CrawlLogListReply) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

type PageLogListRequest struct {
	WarcId               []string  `protobuf:"bytes,1,rep,name=warc_id,json=warcId,proto3" json:"warc_id,omitempty"`
	ExecutionId          string    `protobuf:"bytes,2,opt,name=execution_id,json=executionId,proto3" json:"execution_id,omitempty"`
	Filter               []*Filter `protobuf:"bytes,3,rep,name=filter,proto3" json:"filter,omitempty"`
	PageSize             int32     `protobuf:"varint,14,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Page                 int32     `protobuf:"varint,15,opt,name=page,proto3" json:"page,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *PageLogListRequest) Reset()         { *m = PageLogListRequest{} }
func (m *PageLogListRequest) String() string { return proto.CompactTextString(m) }
func (*PageLogListRequest) ProtoMessage()    {}
func (*PageLogListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{3}
}
func (m *PageLogListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PageLogListRequest.Unmarshal(m, b)
}
func (m *PageLogListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PageLogListRequest.Marshal(b, m, deterministic)
}
func (dst *PageLogListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PageLogListRequest.Merge(dst, src)
}
func (m *PageLogListRequest) XXX_Size() int {
	return xxx_messageInfo_PageLogListRequest.Size(m)
}
func (m *PageLogListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PageLogListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PageLogListRequest proto.InternalMessageInfo

func (m *PageLogListRequest) GetWarcId() []string {
	if m != nil {
		return m.WarcId
	}
	return nil
}

func (m *PageLogListRequest) GetExecutionId() string {
	if m != nil {
		return m.ExecutionId
	}
	return ""
}

func (m *PageLogListRequest) GetFilter() []*Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *PageLogListRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *PageLogListRequest) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

type PageLogListReply struct {
	Value                []*PageLog `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
	Count                int64      `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	PageSize             int32      `protobuf:"varint,14,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Page                 int32      `protobuf:"varint,15,opt,name=page,proto3" json:"page,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *PageLogListReply) Reset()         { *m = PageLogListReply{} }
func (m *PageLogListReply) String() string { return proto.CompactTextString(m) }
func (*PageLogListReply) ProtoMessage()    {}
func (*PageLogListReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{4}
}
func (m *PageLogListReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PageLogListReply.Unmarshal(m, b)
}
func (m *PageLogListReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PageLogListReply.Marshal(b, m, deterministic)
}
func (dst *PageLogListReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PageLogListReply.Merge(dst, src)
}
func (m *PageLogListReply) XXX_Size() int {
	return xxx_messageInfo_PageLogListReply.Size(m)
}
func (m *PageLogListReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PageLogListReply.DiscardUnknown(m)
}

var xxx_messageInfo_PageLogListReply proto.InternalMessageInfo

func (m *PageLogListReply) GetValue() []*PageLog {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *PageLogListReply) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *PageLogListReply) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *PageLogListReply) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

type ScreenshotListRequest struct {
	Id                   []string  `protobuf:"bytes,1,rep,name=id,proto3" json:"id,omitempty"`
	ExecutionId          string    `protobuf:"bytes,2,opt,name=execution_id,json=executionId,proto3" json:"execution_id,omitempty"`
	Filter               []*Filter `protobuf:"bytes,3,rep,name=filter,proto3" json:"filter,omitempty"`
	ImgData              bool      `protobuf:"varint,4,opt,name=img_data,json=imgData,proto3" json:"img_data,omitempty"`
	PageSize             int32     `protobuf:"varint,14,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Page                 int32     `protobuf:"varint,15,opt,name=page,proto3" json:"page,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ScreenshotListRequest) Reset()         { *m = ScreenshotListRequest{} }
func (m *ScreenshotListRequest) String() string { return proto.CompactTextString(m) }
func (*ScreenshotListRequest) ProtoMessage()    {}
func (*ScreenshotListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{5}
}
func (m *ScreenshotListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScreenshotListRequest.Unmarshal(m, b)
}
func (m *ScreenshotListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScreenshotListRequest.Marshal(b, m, deterministic)
}
func (dst *ScreenshotListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScreenshotListRequest.Merge(dst, src)
}
func (m *ScreenshotListRequest) XXX_Size() int {
	return xxx_messageInfo_ScreenshotListRequest.Size(m)
}
func (m *ScreenshotListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ScreenshotListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ScreenshotListRequest proto.InternalMessageInfo

func (m *ScreenshotListRequest) GetId() []string {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *ScreenshotListRequest) GetExecutionId() string {
	if m != nil {
		return m.ExecutionId
	}
	return ""
}

func (m *ScreenshotListRequest) GetFilter() []*Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *ScreenshotListRequest) GetImgData() bool {
	if m != nil {
		return m.ImgData
	}
	return false
}

func (m *ScreenshotListRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ScreenshotListRequest) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

type ScreenshotListReply struct {
	Value                []*Screenshot `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
	Count                int64         `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	PageSize             int32         `protobuf:"varint,14,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Page                 int32         `protobuf:"varint,15,opt,name=page,proto3" json:"page,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ScreenshotListReply) Reset()         { *m = ScreenshotListReply{} }
func (m *ScreenshotListReply) String() string { return proto.CompactTextString(m) }
func (*ScreenshotListReply) ProtoMessage()    {}
func (*ScreenshotListReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{6}
}
func (m *ScreenshotListReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScreenshotListReply.Unmarshal(m, b)
}
func (m *ScreenshotListReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScreenshotListReply.Marshal(b, m, deterministic)
}
func (dst *ScreenshotListReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScreenshotListReply.Merge(dst, src)
}
func (m *ScreenshotListReply) XXX_Size() int {
	return xxx_messageInfo_ScreenshotListReply.Size(m)
}
func (m *ScreenshotListReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ScreenshotListReply.DiscardUnknown(m)
}

var xxx_messageInfo_ScreenshotListReply proto.InternalMessageInfo

func (m *ScreenshotListReply) GetValue() []*Screenshot {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *ScreenshotListReply) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *ScreenshotListReply) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ScreenshotListReply) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

type ExecuteDbQueryRequest struct {
	// The query to execute
	Query string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	// Maximum number of rows to return. A limit of -1 indicates no limit. If unset or zero, use default limit.
	Limit                int32    `protobuf:"varint,14,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExecuteDbQueryRequest) Reset()         { *m = ExecuteDbQueryRequest{} }
func (m *ExecuteDbQueryRequest) String() string { return proto.CompactTextString(m) }
func (*ExecuteDbQueryRequest) ProtoMessage()    {}
func (*ExecuteDbQueryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{7}
}
func (m *ExecuteDbQueryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecuteDbQueryRequest.Unmarshal(m, b)
}
func (m *ExecuteDbQueryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecuteDbQueryRequest.Marshal(b, m, deterministic)
}
func (dst *ExecuteDbQueryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteDbQueryRequest.Merge(dst, src)
}
func (m *ExecuteDbQueryRequest) XXX_Size() int {
	return xxx_messageInfo_ExecuteDbQueryRequest.Size(m)
}
func (m *ExecuteDbQueryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteDbQueryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteDbQueryRequest proto.InternalMessageInfo

func (m *ExecuteDbQueryRequest) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

func (m *ExecuteDbQueryRequest) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type ExecuteDbQueryReply struct {
	Record               string   `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExecuteDbQueryReply) Reset()         { *m = ExecuteDbQueryReply{} }
func (m *ExecuteDbQueryReply) String() string { return proto.CompactTextString(m) }
func (*ExecuteDbQueryReply) ProtoMessage()    {}
func (*ExecuteDbQueryReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_report_7cb662b0d011aad5, []int{8}
}
func (m *ExecuteDbQueryReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecuteDbQueryReply.Unmarshal(m, b)
}
func (m *ExecuteDbQueryReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecuteDbQueryReply.Marshal(b, m, deterministic)
}
func (dst *ExecuteDbQueryReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteDbQueryReply.Merge(dst, src)
}
func (m *ExecuteDbQueryReply) XXX_Size() int {
	return xxx_messageInfo_ExecuteDbQueryReply.Size(m)
}
func (m *ExecuteDbQueryReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteDbQueryReply.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteDbQueryReply proto.InternalMessageInfo

func (m *ExecuteDbQueryReply) GetRecord() string {
	if m != nil {
		return m.Record
	}
	return ""
}

func init() {
	proto.RegisterType((*Filter)(nil), "veidemann.api.Filter")
	proto.RegisterType((*CrawlLogListRequest)(nil), "veidemann.api.CrawlLogListRequest")
	proto.RegisterType((*CrawlLogListReply)(nil), "veidemann.api.CrawlLogListReply")
	proto.RegisterType((*PageLogListRequest)(nil), "veidemann.api.PageLogListRequest")
	proto.RegisterType((*PageLogListReply)(nil), "veidemann.api.PageLogListReply")
	proto.RegisterType((*ScreenshotListRequest)(nil), "veidemann.api.ScreenshotListRequest")
	proto.RegisterType((*ScreenshotListReply)(nil), "veidemann.api.ScreenshotListReply")
	proto.RegisterType((*ExecuteDbQueryRequest)(nil), "veidemann.api.ExecuteDbQueryRequest")
	proto.RegisterType((*ExecuteDbQueryReply)(nil), "veidemann.api.ExecuteDbQueryReply")
	proto.RegisterEnum("veidemann.api.Filter_Operator", Filter_Operator_name, Filter_Operator_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ReportClient is the client API for Report service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReportClient interface {
	// List crawl logs
	ListCrawlLogs(ctx context.Context, in *CrawlLogListRequest, opts ...grpc.CallOption) (*CrawlLogListReply, error)
	// List page logs
	ListPageLogs(ctx context.Context, in *PageLogListRequest, opts ...grpc.CallOption) (*PageLogListReply, error)
	// List screenshots
	ListScreenshots(ctx context.Context, in *ScreenshotListRequest, opts ...grpc.CallOption) (*ScreenshotListReply, error)
	// Execute a query against the database
	ExecuteDbQuery(ctx context.Context, in *ExecuteDbQueryRequest, opts ...grpc.CallOption) (Report_ExecuteDbQueryClient, error)
}

type reportClient struct {
	cc *grpc.ClientConn
}

func NewReportClient(cc *grpc.ClientConn) ReportClient {
	return &reportClient{cc}
}

func (c *reportClient) ListCrawlLogs(ctx context.Context, in *CrawlLogListRequest, opts ...grpc.CallOption) (*CrawlLogListReply, error) {
	out := new(CrawlLogListReply)
	err := c.cc.Invoke(ctx, "/veidemann.api.Report/ListCrawlLogs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportClient) ListPageLogs(ctx context.Context, in *PageLogListRequest, opts ...grpc.CallOption) (*PageLogListReply, error) {
	out := new(PageLogListReply)
	err := c.cc.Invoke(ctx, "/veidemann.api.Report/ListPageLogs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportClient) ListScreenshots(ctx context.Context, in *ScreenshotListRequest, opts ...grpc.CallOption) (*ScreenshotListReply, error) {
	out := new(ScreenshotListReply)
	err := c.cc.Invoke(ctx, "/veidemann.api.Report/ListScreenshots", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportClient) ExecuteDbQuery(ctx context.Context, in *ExecuteDbQueryRequest, opts ...grpc.CallOption) (Report_ExecuteDbQueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Report_serviceDesc.Streams[0], "/veidemann.api.Report/ExecuteDbQuery", opts...)
	if err != nil {
		return nil, err
	}
	x := &reportExecuteDbQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Report_ExecuteDbQueryClient interface {
	Recv() (*ExecuteDbQueryReply, error)
	grpc.ClientStream
}

type reportExecuteDbQueryClient struct {
	grpc.ClientStream
}

func (x *reportExecuteDbQueryClient) Recv() (*ExecuteDbQueryReply, error) {
	m := new(ExecuteDbQueryReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReportServer is the server API for Report service.
type ReportServer interface {
	// List crawl logs
	ListCrawlLogs(context.Context, *CrawlLogListRequest) (*CrawlLogListReply, error)
	// List page logs
	ListPageLogs(context.Context, *PageLogListRequest) (*PageLogListReply, error)
	// List screenshots
	ListScreenshots(context.Context, *ScreenshotListRequest) (*ScreenshotListReply, error)
	// Execute a query against the database
	ExecuteDbQuery(*ExecuteDbQueryRequest, Report_ExecuteDbQueryServer) error
}

func RegisterReportServer(s *grpc.Server, srv ReportServer) {
	s.RegisterService(&_Report_serviceDesc, srv)
}

func _Report_ListCrawlLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CrawlLogListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServer).ListCrawlLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/veidemann.api.Report/ListCrawlLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServer).ListCrawlLogs(ctx, req.(*CrawlLogListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Report_ListPageLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PageLogListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServer).ListPageLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/veidemann.api.Report/ListPageLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServer).ListPageLogs(ctx, req.(*PageLogListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Report_ListScreenshots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScreenshotListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServer).ListScreenshots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/veidemann.api.Report/ListScreenshots",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServer).ListScreenshots(ctx, req.(*ScreenshotListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Report_ExecuteDbQuery_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExecuteDbQueryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReportServer).ExecuteDbQuery(m, &reportExecuteDbQueryServer{stream})
}

type Report_ExecuteDbQueryServer interface {
	Send(*ExecuteDbQueryReply) error
	grpc.ServerStream
}

type reportExecuteDbQueryServer struct {
	grpc.ServerStream
}

func (x *reportExecuteDbQueryServer) Send(m *ExecuteDbQueryReply) error {
	return x.ServerStream.SendMsg(m)
}

var _Report_serviceDesc = grpc.ServiceDesc{
	ServiceName: "veidemann.api.Report",
	HandlerType: (*ReportServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListCrawlLogs",
			Handler:    _Report_ListCrawlLogs_Handler,
		},
		{
			MethodName: "ListPageLogs",
			Handler:    _Report_ListPageLogs_Handler,
		},
		{
			MethodName: "ListScreenshots",
			Handler:    _Report_ListScreenshots_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ExecuteDbQuery",
			Handler:       _Report_ExecuteDbQuery_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "report.proto",
}

func init() { proto.RegisterFile("report.proto", fileDescriptor_report_7cb662b0d011aad5) }

var fileDescriptor_report_7cb662b0d011aad5 = []byte{
	// 783 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x55, 0xc1, 0x6e, 0xdb, 0x46,
	0x10, 0xed, 0x92, 0x92, 0x62, 0x8d, 0x1d, 0x45, 0x5d, 0x59, 0x16, 0xed, 0x24, 0x0d, 0x43, 0xb4,
	0x80, 0x50, 0x44, 0xa4, 0xa1, 0xa2, 0x97, 0xde, 0x1c, 0xc5, 0x6d, 0x0c, 0xb8, 0xa9, 0xc2, 0x18,
	0x3d, 0xf4, 0x22, 0xac, 0xc8, 0x0d, 0xbd, 0x00, 0xb9, 0xcb, 0x90, 0xab, 0xa8, 0x0a, 0x7a, 0xea,
	0xa5, 0xed, 0xa1, 0x97, 0xf6, 0x17, 0xfa, 0x01, 0xed, 0x47, 0xf4, 0x0b, 0xfa, 0x0b, 0xbd, 0xf4,
	0x2f, 0x8a, 0x5d, 0x52, 0x36, 0x25, 0x2b, 0x29, 0x02, 0xb8, 0x87, 0x9c, 0x96, 0x33, 0xfb, 0x38,
	0x6f, 0xe6, 0xcd, 0xec, 0x2e, 0xec, 0x64, 0x34, 0x15, 0x99, 0x74, 0xd3, 0x4c, 0x48, 0x81, 0x6f,
	0xbe, 0xa4, 0x2c, 0xa4, 0x09, 0xe1, 0xdc, 0x25, 0x29, 0x3b, 0x68, 0x25, 0x34, 0xcf, 0x49, 0x44,
	0xf3, 0x62, 0xfb, 0xe0, 0x4e, 0x24, 0x44, 0x14, 0x53, 0x8f, 0xa4, 0xcc, 0x23, 0x9c, 0x0b, 0x49,
	0x24, 0x13, 0x7c, 0xb9, 0xfb, 0x40, 0x2f, 0xc1, 0x20, 0xa2, 0x7c, 0x90, 0xcf, 0x49, 0x14, 0xd1,
	0xcc, 0x13, 0xa9, 0x46, 0x5c, 0x45, 0x3b, 0xbf, 0x21, 0x68, 0x7c, 0xce, 0x62, 0x49, 0x33, 0x7c,
	0x17, 0xe0, 0x39, 0xa3, 0x71, 0x38, 0xe1, 0x24, 0xa1, 0x16, 0xb2, 0x51, 0xbf, 0xe9, 0x37, 0xb5,
	0xe7, 0x09, 0x49, 0x28, 0x76, 0xc1, 0x10, 0xa9, 0x65, 0xd8, 0xa8, 0xdf, 0x1a, 0x7e, 0xe0, 0xae,
	0x64, 0xe8, 0x16, 0x11, 0xdc, 0xaf, 0x52, 0x9a, 0x11, 0x29, 0x32, 0xdf, 0x10, 0x29, 0xde, 0x85,
	0xfa, 0x4b, 0x12, 0xcf, 0xa8, 0x65, 0xea, 0x48, 0x85, 0xe1, 0x7c, 0x0a, 0x5b, 0x4b, 0x14, 0x6e,
	0x80, 0x71, 0xfc, 0xb4, 0xfd, 0x9e, 0x5a, 0x9f, 0x1c, 0xb7, 0x11, 0x6e, 0x42, 0xfd, 0xcb, 0xa3,
	0xb3, 0xd1, 0xe3, 0xb6, 0xa1, 0x5c, 0xa7, 0x67, 0x6d, 0x53, 0xad, 0x5f, 0x9c, 0xb5, 0x6b, 0xce,
	0x1f, 0x08, 0x3a, 0xa3, 0x8c, 0xcc, 0xe3, 0x53, 0x11, 0x9d, 0xb2, 0x5c, 0xfa, 0xf4, 0xc5, 0x8c,
	0xe6, 0x12, 0xf7, 0xe0, 0xc6, 0x9c, 0x64, 0xc1, 0x84, 0x85, 0x16, 0xb2, 0xcd, 0x7e, 0xd3, 0x6f,
	0x28, 0xf3, 0x24, 0xc4, 0xf7, 0x61, 0x87, 0x7e, 0x4b, 0x83, 0x99, 0xaa, 0x55, 0xed, 0x1a, 0x3a,
	0x89, 0xed, 0x0b, 0xdf, 0x49, 0x88, 0x07, 0xd0, 0x78, 0xae, 0xf3, 0xb6, 0x4c, 0xdb, 0xec, 0x6f,
	0x0f, 0xbb, 0x1b, 0x8b, 0xf2, 0x4b, 0x10, 0xbe, 0x0d, 0xcd, 0x94, 0x44, 0x74, 0x92, 0xb3, 0x57,
	0xd4, 0x6a, 0xd9, 0xa8, 0x5f, 0xf7, 0xb7, 0x94, 0xe3, 0x19, 0x7b, 0x45, 0x31, 0x86, 0x9a, 0xfa,
	0xb6, 0x6e, 0x69, 0xbf, 0xfe, 0x76, 0x7e, 0x42, 0xf0, 0xfe, 0x6a, 0xce, 0x69, 0xbc, 0xc0, 0x83,
	0xa5, 0x2c, 0x48, 0x93, 0xf6, 0xd6, 0x48, 0x97, 0x3f, 0x94, 0x7a, 0x29, 0x15, 0x03, 0x31, 0xe3,
	0x52, 0x17, 0x60, 0xfa, 0x85, 0xf1, 0xf6, 0xb9, 0xfc, 0x8e, 0x00, 0x8f, 0x49, 0x44, 0xdf, 0x21,
	0xf9, 0x7e, 0x40, 0xd0, 0x5e, 0x49, 0x59, 0xa9, 0xf7, 0x60, 0x55, 0xbd, 0xbd, 0x35, 0xce, 0x12,
	0x7f, 0xcd, 0xe2, 0xfd, 0x89, 0xa0, 0xfb, 0x2c, 0xc8, 0x28, 0xe5, 0xf9, 0xb9, 0x90, 0x55, 0xfd,
	0x5a, 0x60, 0x5c, 0x48, 0x67, 0xb0, 0xff, 0x43, 0xb6, 0x7d, 0xd8, 0x62, 0x49, 0x34, 0x09, 0x89,
	0x24, 0x56, 0xcd, 0x46, 0xfd, 0x2d, 0xff, 0x06, 0x4b, 0xa2, 0x47, 0x44, 0x92, 0xb7, 0xaf, 0xe3,
	0x67, 0x04, 0x9d, 0xf5, 0x3a, 0x94, 0xa8, 0xde, 0xaa, 0xa8, 0xfb, 0x6b, 0x19, 0x5d, 0xfe, 0x72,
	0xcd, 0xba, 0x8e, 0xa0, 0x7b, 0xac, 0x95, 0xa1, 0x8f, 0xa6, 0x4f, 0x67, 0x34, 0x5b, 0x2c, 0x65,
	0xdd, 0x85, 0xfa, 0x0b, 0x65, 0x97, 0x97, 0x50, 0x61, 0x28, 0x6f, 0xcc, 0x12, 0x26, 0xcb, 0xd8,
	0x85, 0xe1, 0x0c, 0xa0, 0xb3, 0x1e, 0x44, 0xd5, 0xb4, 0x07, 0x8d, 0x8c, 0x06, 0x22, 0x0b, 0xcb,
	0x18, 0xa5, 0x35, 0xfc, 0xc7, 0x84, 0x86, 0xaf, 0xef, 0x5a, 0x2c, 0xe1, 0xa6, 0xd2, 0x60, 0x79,
	0xe2, 0x72, 0xec, 0xbc, 0xe6, 0x2c, 0x56, 0x3a, 0x7e, 0x60, 0xbf, 0x11, 0x93, 0xc6, 0x0b, 0xe7,
	0xee, 0xf7, 0x7f, 0xfd, 0xfd, 0xab, 0xd1, 0xc3, 0x5d, 0x7d, 0x3f, 0x17, 0xf7, 0xba, 0x17, 0x28,
	0x58, 0xac, 0x48, 0x52, 0xd8, 0x51, 0xd8, 0x72, 0x52, 0x73, 0x7c, 0x7f, 0xf3, 0x08, 0x57, 0x39,
	0xef, 0xbd, 0x09, 0xa2, 0x28, 0xef, 0x68, 0xca, 0x3d, 0xbc, 0x5b, 0xa5, 0x54, 0x12, 0x6b, 0xc6,
	0xef, 0xe0, 0x96, 0x82, 0x5e, 0xb6, 0x31, 0xc7, 0x1f, 0xbe, 0xb6, 0xc5, 0x55, 0x5e, 0xe7, 0x3f,
	0x50, 0x8a, 0xfa, 0x9e, 0xa6, 0xde, 0xc7, 0xbd, 0x2a, 0x75, 0x5e, 0xa1, 0x5a, 0x40, 0x6b, 0xb5,
	0x3f, 0x57, 0xc8, 0x37, 0xce, 0xc0, 0x15, 0xf2, 0x0d, 0x4d, 0x76, 0x6e, 0x6b, 0xf2, 0x2e, 0xee,
	0x54, 0xc9, 0xc3, 0xa9, 0x1e, 0x97, 0x43, 0xf4, 0xf0, 0x47, 0xf4, 0xcb, 0xd1, 0x18, 0x3f, 0x06,
	0xeb, 0xeb, 0x65, 0x24, 0x7b, 0x24, 0xb8, 0xcc, 0x44, 0x1c, 0xd3, 0xcc, 0x3e, 0x1a, 0x9f, 0x38,
	0x1f, 0x41, 0xf3, 0x62, 0x0f, 0x5b, 0xe7, 0x52, 0xa6, 0xf9, 0x67, 0x9e, 0x17, 0x31, 0x79, 0x3e,
	0x9b, 0xba, 0x81, 0x48, 0x3c, 0x1e, 0xf3, 0x39, 0x19, 0xd6, 0x0f, 0xdd, 0xa1, 0x7b, 0xf8, 0x71,
	0x0d, 0x19, 0x66, 0x0d, 0x7a, 0x5c, 0xb8, 0x7c, 0xea, 0x72, 0x4e, 0x56, 0x33, 0x7c, 0xb8, 0x5d,
	0xcc, 0xd5, 0x58, 0xbd, 0xab, 0xdf, 0x5c, 0x3e, 0xe1, 0x13, 0x92, 0xb2, 0x69, 0x43, 0xbf, 0xb6,
	0x9f, 0xfc, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x82, 0x3a, 0xa6, 0xc3, 0xe8, 0x07, 0x00, 0x00,
}
