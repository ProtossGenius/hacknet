package main

import (
	"fmt"
	"testing"

	"github.com/ProtossGenius/hacknet/cltimpl/cltric/clt_rpc_cltitf"
	"github.com/ProtossGenius/hacknet/cltimpl/svrric/svr_rpc_cltitf"
)

// TestClient test clientf.
type TestClient struct {
}

func (t *TestClient) SendForward(email string, data []byte) {
	fmt.Println("recv email .. ", email)
}

func (t *TestClient) OnForward(data []byte) {
	panic("not implemented") // TODO: Implement
}

func TestHelloWorld(t *testing.T) {
	svr := svr_rpc_cltitf.NewSvrRpcCltOperItf(&TestClient{})
	svr.OnMessage(nil, nil)
	clt := clt_rpc_cltitf.NewCltRpcCltOperItf(nil, 0)
	clt.OnErr(nil)
}
