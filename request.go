package main

import (
	"fmt"
	"log"
	"net"

	"github.com/onesimus-systems/dhcp4"
)

func doRequest(conn *conn, mac net.HardwareAddr, serverIP, relayIP, requestIP net.IP) {
	doRequestWithXID(conn, mac, serverIP, relayIP, requestIP, nil)
}

func doRequestWithXID(conn *conn, mac net.HardwareAddr, serverIP, relayIP, requestIP net.IP, xid []byte) {
	log.Println("Sending Request:")
	if xid == nil {
		xid = randomID()
	}

	opts := []dhcp4.Option{
		dhcp4.Option{
			Code:  dhcp4.OptionParameterRequestList,
			Value: []byte{0x1, 0x3, 0x6, 0xf, 0x23},
		},
		dhcp4.Option{
			Code:  dhcp4.OptionServerIdentifier,
			Value: []byte(serverIP),
		},
		dhcp4.Option{
			Code:  dhcp4.OptionRequestedIPAddress,
			Value: []byte(requestIP),
		},
	}

	packet := dhcp4.RequestPacket(dhcp4.Request, mac, nil, xid, false, opts)
	packet.SetGIAddr(relayIP)
	printPacket(packet)

	if _, err := conn.sendPacket(packet); err != nil {
		log.Println(err.Error())
		return
	}

	resp, err := conn.nextPacket()
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("")
	log.Println("Server reply:")
	printPacket(resp)
}
