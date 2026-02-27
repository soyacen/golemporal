package constant

import (
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	ErrorsPackage = protogen.GoImportPath("errors")
	NewErrorIdent = ErrorsPackage.Ident("New")
)

var (
	BytesPackage = protogen.GoImportPath("bytes")
	Buffer       = BytesPackage.Ident("Buffer")
)

var (
	ProtoJsonPackage               = protogen.GoImportPath("google.golang.org/protobuf/encoding/protojson")
	ProtoJsonMarshalOptionsIdent   = ProtoJsonPackage.Ident("MarshalOptions")
	ProtoJsonUnmarshalOptionsIdent = ProtoJsonPackage.Ident("UnmarshalOptions")
)

var (
	ContextPackage = protogen.GoImportPath("context")
	ContextIdent   = ContextPackage.Ident("Context")
)

var (
	WorkflowPackage      = protogen.GoImportPath("go.temporal.io/sdk/workflow")
	WorkflowContextIdent = WorkflowPackage.Ident("Context")
)

var (
	HttpPackage                 = protogen.GoImportPath("net/http")
	NewServeMuxIndent           = HttpPackage.Ident("NewServeMux")
	RouterIdent                 = HttpPackage.Ident("ServeMux")
	ClientIdent                 = HttpPackage.Ident("Client")
	HttpHandlerIdent            = HttpPackage.Ident("Handler")
	HttpHandlerFuncIdent        = HttpPackage.Ident("HandlerFunc")
	ResponseWriterIdent         = HttpPackage.Ident("ResponseWriter")
	RequestIdent                = HttpPackage.Ident("Request")
	ResponseIdent               = HttpPackage.Ident("Response")
	Handler                     = HttpPackage.Ident("Handler")
	Header                      = HttpPackage.Ident("Header")
	NewRequestWithContextIndent = HttpPackage.Ident("NewRequestWithContext")
)

var (
	FmtPackage   = protogen.GoImportPath("fmt")
	SprintfIdent = FmtPackage.Ident("Sprintf")
)

var (
	ProtoPackage     = protogen.GoImportPath("google.golang.org/protobuf/proto")
	ProtoStringIdent = ProtoPackage.Ident("String")
)

var (
	WrapperspbPackage     = protogen.GoImportPath("google.golang.org/protobuf/types/known/wrapperspb")
	WrapperspbStringIdent = WrapperspbPackage.Ident("String")
)
