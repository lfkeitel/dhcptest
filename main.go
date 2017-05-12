package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/onesimus-systems/dhcp4"
)

var (
	macAddress   string
	dhcpType     string
	dhcpServerIP string
	fullDiscover bool
	gatewayIP    string
	requestIP    string
	readTimeout  int
	clientPort   int
)

func init() {
	flag.StringVar(&macAddress, "mac", "", "MAC address")
	flag.StringVar(&dhcpType, "type", "discover", "DHCP message type")
	flag.StringVar(&dhcpServerIP, "server", "", "DHCP server IP")
	flag.StringVar(&gatewayIP, "gw", "0.0.0.0", "DHCP relay IP")
	flag.StringVar(&requestIP, "ci", "0.0.0.0", "IP to request")
	flag.IntVar(&readTimeout, "timeout", 5, "Read timeout in seconds")
	flag.IntVar(&clientPort, "port", 68, "DHCP client port")
	flag.BoolVar(&fullDiscover, "full", true, "Perform a full discover, request sequence")

	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()

	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		log.Fatal(err.Error())
	}

	serverIP := net.ParseIP(dhcpServerIP)
	if serverIP == nil {
		log.Fatal("Bad server IP")
	}

	serverIP = serverIP.To4()
	if serverIP == nil {
		log.Fatal("Only DHCPv4 is supported")
	}

	relayIP := net.ParseIP(gatewayIP)
	if relayIP == nil {
		log.Fatal("Bad relay IP")
	}

	relayIP = relayIP.To4()
	if relayIP == nil {
		log.Fatal("Only DHCPv4 is supported")
	}

	requestedIP := net.ParseIP(requestIP)
	if requestedIP == nil {
		log.Fatal("Bad request IP")
	}

	requestedIP = requestedIP.To4()
	if requestedIP == nil {
		log.Fatal("Only DHCPv4 is supported")
	}

	conn, err := newConn(time.Duration(readTimeout)*time.Second, serverIP, clientPort)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	switch dhcpType {
	case "discover":
		doDiscover(conn, mac, serverIP, relayIP, fullDiscover)
	case "request":
		doRequest(conn, mac, serverIP, relayIP, requestedIP)
	case "inform":
		doInform(conn, mac, serverIP, requestedIP)
	default:
		log.Fatal("Unsupported DHCP message type")
	}
}

type conn struct {
	readTimeout time.Duration
	server      net.IP
	c           net.PacketConn
}

func newConn(readTimeout time.Duration, server net.IP, port int) (*conn, error) {
	c, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return &conn{
		readTimeout: readTimeout,
		server:      server,
		c:           c,
	}, nil
}

func (c *conn) Close() error {
	return c.c.Close()
}

func (c *conn) sendPacket(packet dhcp4.Packet) (int, error) {
	return c.c.WriteTo(packet, &net.UDPAddr{IP: c.server, Port: 67})
}

func (c *conn) nextPacket() (dhcp4.Packet, error) {
	buffer := make([]byte, 1500)
	c.c.SetReadDeadline(time.Now().Add(c.readTimeout))

	n, addr, err := c.c.ReadFrom(buffer)
	if err != nil {
		return nil, err
	}

	ipStr, _, _ := net.SplitHostPort(addr.String())
	if ipStr != c.server.String() {
		return nil, errors.New("received from wrong server")
	}

	if n < 240 { // Packet too small to be DHCP
		return nil, errors.New("invalid DHCP packet")
	}

	req := dhcp4.Packet(buffer[:n])
	if req.HLen() > 16 { // Invalid size
		return nil, errors.New("invalid DHCP packet")
	}

	options := req.ParseOptions()

	t := options[dhcp4.OptionDHCPMessageType]
	if len(t) != 1 {
		return nil, errors.New("invalid DHCP packet")
	}

	reqType := dhcp4.MessageType(t[0])
	if reqType < dhcp4.Discover || reqType > dhcp4.Inform {
		return nil, errors.New("invalid DHCP packet")
	}

	return req, nil
}

const idAlphaNum = `ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`

func randomID() []byte {
	xid := make([]byte, 4)
	for i := range xid {
		xid[i] = idAlphaNum[rand.Intn(len(idAlphaNum))]
	}
	return xid
}
