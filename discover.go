package main

import (
	"fmt"
	"log"
	"net"

	"github.com/onesimus-systems/dhcp4"
)

func doDiscover(conn *conn, mac net.HardwareAddr, serverIP, relayIP net.IP, fullDiscover bool) {
	log.Println("Sending Discover:")
	opts := []dhcp4.Option{
		dhcp4.Option{
			Code:  dhcp4.OptionParameterRequestList,
			Value: []byte{0x1, 0x3, 0x6, 0xf, 0x23},
		},
	}

	packet := dhcp4.RequestPacket(dhcp4.Discover, mac, nil, randomID(), false, opts)
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

	if !fullDiscover {
		return
	}

	fmt.Println("")
	doRequestWithXID(conn, mac, serverIP, relayIP, resp.YIAddr(), resp.XId())
}
