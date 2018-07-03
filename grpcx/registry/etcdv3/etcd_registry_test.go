package etcdv3_test

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/grpcx/registry/etcdv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"testing"
	"time"
)

const testSvrName = "test"

func init() {
	go newMockServer(":12345", "default")
}

func newMockServer(listen string, nodeId string) error {

	l, e := net.Listen("tcp", listen)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	server := grpc.NewServer()
	RegisterHelloServiceServer(server, &hello{listen})

	r, err := etcdv3.NewRegistry(registry.Timeout(5 * time.Minute))
	if err != nil {
		panic(err)
	}
	err = r.Register(testSvrName, &registry.Node{Id: nodeId, Address: listen})
	if err != nil {
		panic(err)
	}
	server.Serve(l)
	return nil
}

func TestEtcdv3Registry_Build(t *testing.T) {
	r, _ := etcdv3.NewRegistry()
	b := r.(resolver.Builder)
	resolver.Register(b)

	_, err := grpc.Dial(testSvrName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

}

func TestEtcdv3Registry_Register(t *testing.T) {
	r, err := etcdv3.NewRegistry(registry.Timeout(5 * time.Minute))
	if err != nil {
		t.Error(err)
	}
	err = r.Register(testSvrName, &registry.Node{Id: testSvrName, Address: ":12345"})
	if err != nil {
		t.Error(err)
	}
}

func TestMulti(t *testing.T) {
	//
	go newMockServer(":12346", "n2")
	go newMockServer(":12347", "n3")
	r, _ := etcdv3.NewRegistry(registry.Timeout(30 * time.Second))
	b := r.(resolver.Builder)
	resolver.Register(b)
	//d1
	d1, err := grpc.Dial(b.Scheme()+"://author/"+testSvrName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	for {
		out := new(Payload)
		err = d1.Invoke(context.Background(), "/pb.HelloService/Echo", &Payload{Data: "hello"}, out)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(out)
		}
		<-time.After(time.Second)
	}
}

//grpc server
type hello struct {
	listen string
}

func (t *hello) Echo(ctx context.Context, req *Payload) (*Payload, error) {
	req.Data = req.Data + ", from:" + t.listen
	return req, nil
}

type Payload struct {
	Data string `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *Payload) Reset()                    { *m = Payload{} }
func (m *Payload) String() string            { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()               {}
func (*Payload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Payload) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*Payload)(nil), "pb.Payload")
	proto.RegisterFile("pb/hello.proto", fileDescriptor0)
}

type HelloServiceServer interface {
	Echo(context.Context, *Payload) (*Payload, error)
}

func RegisterHelloServiceServer(s *grpc.Server, srv HelloServiceServer) {
	s.RegisterService(&_HelloService_serviceDesc, srv)
}

func _HelloService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Payload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.HelloService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServiceServer).Echo(ctx, req.(*Payload))
	}
	return interceptor(ctx, in, info, handler)
}

var _HelloService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.HelloService",
	HandlerType: (*HelloServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _HelloService_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/hello.proto",
}

var fileDescriptor0 = []byte{
	// 112 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x48, 0xd2, 0xcf,
	0x48, 0xcd, 0xc9, 0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x92,
	0xe5, 0x62, 0x0f, 0x48, 0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0x11, 0x12, 0xe2, 0x62, 0x49, 0x49, 0x2c,
	0x49, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x8d, 0x8c, 0xb8, 0x78, 0x3c, 0x40,
	0x3a, 0x82, 0x53, 0x8b, 0xca, 0x32, 0x93, 0x53, 0x85, 0x94, 0xb8, 0x58, 0x5c, 0x93, 0x33, 0xf2,
	0x85, 0xb8, 0xf5, 0x0a, 0x92, 0xf4, 0xa0, 0x1a, 0xa5, 0x90, 0x39, 0x4a, 0x0c, 0x49, 0x6c, 0x60,
	0xd3, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x7d, 0x47, 0x10, 0xcf, 0x6f, 0x00, 0x00, 0x00,
}
