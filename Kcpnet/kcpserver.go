package Kcpnet

// by udp

import (
	"fmt"

	"github.com/xtaci/kcp-go"
	//cli "gopkg.in/urfave/cli.v2"
	"context"
	"os"
	"sync"

	"github.com/Peakchen/xgameCommon/akLog"
	"github.com/Peakchen/xgameCommon/pprof"
	cli "github.com/urfave/cli"
)

type KcpServer struct {
	sw      sync.WaitGroup
	svrName string
	pack    IMessagePack
	addr    string
	ppAddr  string
	cancel  context.CancelFunc
	offCh   chan *KcpServerSession
}

func NewKcpServer(Name string, addr string, pprofAddr string) *KcpServer {
	return &KcpServer{
		svrName: Name,
		addr:    addr,
		ppAddr:  pprofAddr,
		offCh:   make(chan *KcpServerSession, 1000),
	}
}

func (this *KcpServer) Run() {
	os.Setenv("GOTRACEBACK", "crash")

	ctx, _ := context.WithCancel(context.Background())
	pprof.Run(ctx)

	app := &cli.App{
		Name:    this.svrName,
		Usage:   "a server...",
		Version: "v1.0",
		Flags:   []cli.Flag{},
		Action: func(c *cli.Context) error {
			akLog.FmtPrintln("action begin...")

			//setup net param
			config := &KcpSvrConfig{
				listen:       c.String("listen"),
				readDeadline: c.Duration("read-deadline"),
				sockbuf:      c.Int("sockbuf"),
				udp_sockbuf:  c.Int("udp-sockbuf"),
				txqueuelen:   c.Int("txqueuelen"),
				dscp:         c.Int("dscp"),
				sndwnd:       c.Int("udp-sndwnd"),
				rcvwnd:       c.Int("udp-rcvwnd"),
				mtu:          c.Int("udp-mtu"),
				nodelay:      c.Int("nodelay"),
				interval:     c.Int("interval"),
				resend:       c.Int("resend"),
				nc:           c.Int("nc"),
			}
			// init services
			//startup(c)
			// init timer
			//initTimer(c.Int("rpm-limit"))

			// start udp server...
			this.sw.Add(1)
			go this.kcpAccept(config)
			go this.loopOffline()
			this.sw.Wait()
			return nil
		},
	}

	app.Run(os.Args)
}

func (this *KcpServer) kcpAccept(c *KcpSvrConfig) {
	l, err := kcp.Listen(this.addr)
	if err != nil {
		panic(err)
	}

	akLog.FmtPrintln("kcp listening on:", l.Addr())
	kcplis := l.(*kcp.Listener)
	if err := kcplis.SetReadBuffer(c.sockbuf); err != nil {
		panic(fmt.Errorf("SetReadBuffer, err: %v.", err))
	}

	if err := kcplis.SetWriteBuffer(c.sockbuf); err != nil {
		panic(fmt.Errorf("SetWriteBuffer, err: %v.", err))
	}

	if err := kcplis.SetDSCP(c.dscp); err != nil {
		panic(fmt.Errorf("SetDSCP, err: %v.", err))
	}

	// loop accepting
	for {
		conn, err := kcplis.AcceptKCP()
		if err != nil {
			akLog.FmtPrintln("accept failed:", err)
			continue
		}

		// set kcp parameters
		conn.SetWindowSize(c.sndwnd, c.rcvwnd)
		conn.SetNoDelay(c.nodelay, c.interval, c.resend, c.nc)
		conn.SetStreamMode(true)
		conn.SetMtu(c.mtu)

		// start a goroutine for every incoming connection for read and write
		//go handleClient(conn, config)
		sess := NewKcpSvrSession(conn, this.offCh, c)
		sess.Handler()
	}
}

func (this *KcpServer) loopOffline() {
	for {
		offsession := <-this.offCh
		offsession.Offline()
	}
}
