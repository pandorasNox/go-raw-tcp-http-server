// taken / reduced from https://github.com/jpillora/go-tcp-proxy

package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/mgutz/ansi"
)

var (
	connid = uint64(0)
	logger ColorLogger

	localAddr  = flag.String("l", ":9999", "local address")
	targetAddr = flag.String("t", "localhost:80", "remote target address")
	verbose    = flag.Bool("v", false, "display server actions")
	hex        = flag.Bool("h", false, "output hex")
)

func main() {
	flag.Parse()

	logger := ColorLogger{
		Verbose:     *verbose,
		VeryVerbose: *verbose,
	}

	logger.Info("go-tcp-proxy - proxing from %v to %v ", *localAddr, *targetAddr)

	laddr, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		logger.Warn("Failed to resolve local address: %s", err)
		os.Exit(1)
	}
	targetAddr, err := net.ResolveTCPAddr("tcp", *targetAddr)
	if err != nil {
		logger.Warn("Failed to resolve remote address: %s", err)
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		logger.Warn("Failed to open local port to listen: %s", err)
		os.Exit(1)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Warn("Failed to accept connection '%s'", err)
			continue
		}
		connid++

		// var p *Proxy
		p := NewProxy(conn, laddr, targetAddr)

		p.OutputHex = *hex
		p.Log = ColorLogger{
			Verbose:     *verbose,
			VeryVerbose: *verbose,
			Prefix:      fmt.Sprintf("Connection #%03d ", connid),
			// Color:   false,
		}

		go p.Start()
	}
}

// ---------------

// Proxy - Manages a Proxy connection, piping data between local and remote.
type Proxy struct {
	sentBytes                   uint64
	receivedBytes               uint64
	incommingAddr, targetAddr   *net.TCPAddr
	incommingConn, outgoingConn io.ReadWriteCloser
	erred                       bool
	errsig                      chan bool

	// Settings
	Log       Logger
	OutputHex bool
}

// New - Create a new Proxy instance. Takes over local connection passed in,
// and closes it when finished.
func NewProxy(incommingConn *net.TCPConn, incommingAddr, targetAddr *net.TCPAddr) *Proxy {
	return &Proxy{
		incommingConn: incommingConn,
		incommingAddr: incommingAddr,
		targetAddr:    targetAddr,
		erred:         false,
		errsig:        make(chan bool),
		Log:           NullLogger{},
	}
}

// Start - open connection to remote and start proxying data.
func (p *Proxy) Start() {
	defer p.incommingConn.Close()

	var err error
	//connect to remote
	p.outgoingConn, err = net.DialTCP("tcp", nil, p.targetAddr)
	if err != nil {
		p.Log.Warn("Remote connection failed: %s", err)
		return
	}
	defer p.outgoingConn.Close()

	//display both ends
	p.Log.Info("Opened %s >>> %s", p.incommingAddr.String(), p.targetAddr.String())

	//bidirectional copy
	go p.pipe(p.incommingConn, p.outgoingConn)
	go p.pipe(p.outgoingConn, p.incommingConn)

	//wait for close...
	<-p.errsig
	p.Log.Info("Closed (%d bytes sent, %d bytes recieved)", p.sentBytes, p.receivedBytes)
}

func (p *Proxy) err(s string, err error) {
	if p.erred {
		return
	}
	if err != io.EOF {
		p.Log.Warn(s, err)
	}
	p.errsig <- true
	p.erred = true
}

func (p *Proxy) pipe(src, dst io.ReadWriter) {
	islocal := src == p.incommingConn

	var dataDirection string
	if islocal {
		dataDirection = ">>> %d bytes sent%s"
	} else {
		dataDirection = "<<< %d bytes recieved%s"
	}

	var byteFormat string
	if p.OutputHex {
		byteFormat = "%x"
	} else {
		byteFormat = "%s"
	}

	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			p.err("Read failed '%s'\n", err)
			return
		}
		b := buff[:n]

		//show output
		p.Log.Debug(dataDirection, n, "")
		p.Log.Trace(byteFormat, b)

		//write out result
		n, err = dst.Write(b)
		if err != nil {
			p.err("Write failed '%s'\n", err)
			return
		}
		if islocal {
			p.sentBytes += uint64(n)
		} else {
			p.receivedBytes += uint64(n)
		}
	}
}

// -----------

// Logger - Interface to pass into Proxy for it to log messages
type Logger interface {
	Trace(f string, args ...interface{})
	Debug(f string, args ...interface{})
	Info(f string, args ...interface{})
	Warn(f string, args ...interface{})
}

// NullLogger - An empty logger that ignores everything
type NullLogger struct{}

// Trace - no-op
func (l NullLogger) Trace(f string, args ...interface{}) {}

// Debug - no-op
func (l NullLogger) Debug(f string, args ...interface{}) {}

// Info - no-op
func (l NullLogger) Info(f string, args ...interface{}) {}

// Warn - no-op
func (l NullLogger) Warn(f string, args ...interface{}) {}

// ColorLogger - A Logger that logs to stdout in color
type ColorLogger struct {
	VeryVerbose bool
	Verbose     bool
	Prefix      string
	Color       bool
}

// Trace - Log a very verbose trace message
func (l ColorLogger) Trace(f string, args ...interface{}) {
	if !l.VeryVerbose {
		return
	}
	l.output("blue", f, args...)
}

// Debug - Log a debug message
func (l ColorLogger) Debug(f string, args ...interface{}) {
	if !l.Verbose {
		return
	}
	l.output("green", f, args...)
}

// Info - Log a general message
func (l ColorLogger) Info(f string, args ...interface{}) {
	l.output("green", f, args...)
}

// Warn - Log a warning
func (l ColorLogger) Warn(f string, args ...interface{}) {
	l.output("red", f, args...)
}

func (l ColorLogger) output(color, f string, args ...interface{}) {
	if l.Color && color != "" {
		f = ansi.Color(f, color)
	}
	fmt.Printf(fmt.Sprintf("%s%s\n", l.Prefix, f), args...)
}
