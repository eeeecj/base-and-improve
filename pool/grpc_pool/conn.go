package grpc_pool

import "google.golang.org/grpc"

type ConnImp interface {
	Value() *grpc.ClientConn
	Close() error
}
type conn struct {
	conn *grpc.ClientConn
	pool *pool
	once bool
}

func (c *conn) Value() *grpc.ClientConn {
	return c.conn
}

func (c *conn) Close() error {
	c.pool.decrRef()
	if c.once {
		return c.reset()
	}
	return nil
}

func (c *conn) reset() error {
	con := c.conn
	c.conn = nil
	c.once = false
	if con != nil {
		return con.Close()
	}
	return nil
}

func (p *pool) wrapConn(con *grpc.ClientConn, once bool) *conn {
	return &conn{
		conn: con,
		once: once,
		pool: p,
	}
}
