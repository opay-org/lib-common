package clients

import (
	"math/rand"

	"google.golang.org/grpc"
)

const get_conn_retry = 3

type GrpcPool interface {
	Get() (conn *grpc.ClientConn, err error)
	Put(conn *grpc.ClientConn) error
	Close()
}

// TODO: add health check optimize for short conn
func randAddr(addrs []string) (string, int) {
	size := len(addrs)
	if size == 0 {
		return "", 0
	}
	idx := rand.Intn(size)
	return addrs[idx], idx
}

/**
####################################################################################
USE SHORT CONNECTION
*/

func newShortGrpcClientPool(conf GrpcClientConfig, opt ...grpc.DialOption) (pool *ShortGrpcPool, err error) {
	pool = &ShortGrpcPool{
		conf:     conf,
		dialOpts: opt,
	}
	if len(pool.dialOpts) == 0 {
		pool.dialOpts = []grpc.DialOption{grpc.WithInsecure()}
	}
	return
}

type ShortGrpcPool struct {
	conf     GrpcClientConfig
	dialOpts []grpc.DialOption
}

func (pool *ShortGrpcPool) randAddr() string {
	addr, _ := randAddr(pool.conf.Addrs)
	return addr
}

func (pool *ShortGrpcPool) Get() (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(pool.randAddr(), pool.dialOpts...)
	if err != nil {
		conn, err = grpc.Dial(pool.randAddr(), pool.dialOpts...)
	}
	return
}

func (pool *ShortGrpcPool) Put(conn *grpc.ClientConn) error {
	return conn.Close()
}

func (pool *ShortGrpcPool) Close() {

}
