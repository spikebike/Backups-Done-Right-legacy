// Code generated by protoc-gen-go from "bdr_proto.proto"
// DO NOT EDIT!

package bdrservice

import proto "code.google.com/p/goprotobuf/proto"
import "math"

import "net"
import "net/rpc"
import "github.com/kylelemons/go-rpcgen/codec"
import "net/url"
import "net/http"
import "github.com/kylelemons/go-rpcgen/webrpc"

// Reference proto and math imports to suppress error if they are not otherwise used.
var _ = proto.GetString
var _ = math.Inf

type RequestMessage struct {
	Blobarray        []*RequestMessageBlob `protobuf:"bytes,1,rep,name=blobarray" json:"blobarray,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (this *RequestMessage) Reset()         { *this = RequestMessage{} }
func (this *RequestMessage) String() string { return proto.CompactTextString(this) }

type RequestMessageBlob struct {
	Sha256           *string `protobuf:"bytes,1,req,name=sha256" json:"sha256,omitempty"`
	Bsize            *int32  `protobuf:"varint,2,req,name=bsize" json:"bsize,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *RequestMessageBlob) Reset()         { *this = RequestMessageBlob{} }
func (this *RequestMessageBlob) String() string { return proto.CompactTextString(this) }

type RequestACKMessage struct {
	Received         *int32 `protobuf:"varint,1,req,name=received" json:"received,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *RequestACKMessage) Reset()         { *this = RequestACKMessage{} }
func (this *RequestACKMessage) String() string { return proto.CompactTextString(this) }

func init() {
}

// RequestService is an interface satisfied by the generated client and
// which must be implemented by the object wrapped by the server.
type RequestService interface {
	Request(in *RequestMessage, out *RequestACKMessage) error
}

// internal wrapper for type-safe RPC calling
type rpcRequestServiceClient struct {
	*rpc.Client
}

func (this rpcRequestServiceClient) Request(in *RequestMessage, out *RequestACKMessage) error {
	return this.Call("RequestService.Request", in, out)
}

// NewRequestServiceClient returns an *rpc.Client wrapper for calling the methods of
// RequestService remotely.
func NewRequestServiceClient(conn net.Conn) RequestService {
	return rpcRequestServiceClient{rpc.NewClientWithCodec(codec.NewClientCodec(conn))}
}

// ServeRequestService serves the given RequestService backend implementation on conn.
func ServeRequestService(conn net.Conn, backend RequestService) error {
	srv := rpc.NewServer()
	if err := srv.RegisterName("RequestService", backend); err != nil {
		return err
	}
	srv.ServeCodec(codec.NewServerCodec(conn))
	return nil
}

// DialRequestService returns a RequestService for calling the RequestService servince at addr (TCP).
func DialRequestService(addr string) (RequestService, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewRequestServiceClient(conn), nil
}

// ListenAndServeRequestService serves the given RequestService backend implementation
// on all connections accepted as a result of listening on addr (TCP).
func ListenAndServeRequestService(addr string, backend RequestService) error {
	clients, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := rpc.NewServer()
	if err := srv.RegisterName("RequestService", backend); err != nil {
		return err
	}
	for {
		conn, err := clients.Accept()
		if err != nil {
			return err
		}
		go srv.ServeCodec(codec.NewServerCodec(conn))
	}
	panic("unreachable")
}

// RequestServiceWeb is the web-based RPC version of the interface which
// must be implemented by the object wrapped by the webrpc server.
type RequestServiceWeb interface {
	Request(r *http.Request, in *RequestMessage, out *RequestACKMessage) error
}

// internal wrapper for type-safe webrpc calling
type rpcRequestServiceWebClient struct {
	remote   *url.URL
	protocol webrpc.Protocol
}

func (this rpcRequestServiceWebClient) Request(in *RequestMessage, out *RequestACKMessage) error {
	return webrpc.Post(this.protocol, this.remote, "/RequestService/Request", in, out)
}

// Register a RequestServiceWeb implementation with the given webrpc ServeMux.
// If mux is nil, the default webrpc.ServeMux is used.
func RegisterRequestServiceWeb(this RequestServiceWeb, mux webrpc.ServeMux) error {
	if mux == nil {
		mux = webrpc.DefaultServeMux
	}
	if err := mux.Handle("/RequestService/Request", func(c *webrpc.Call) error {
		in, out := new(RequestMessage), new(RequestACKMessage)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Request(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	return nil
}

// NewRequestServiceWebClient returns a webrpc wrapper for calling the methods of RequestService
// remotely via the web.  The remote URL is the base URL of the webrpc server.
func NewRequestServiceWebClient(pro webrpc.Protocol, remote *url.URL) RequestService {
	return rpcRequestServiceWebClient{remote, pro}
}