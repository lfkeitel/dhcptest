package main

import (
	"fmt"
	"log"
	"net"

	"github.com/onesimus-systems/dhcp4"
)

func doInform(conn *conn, mac net.HardwareAddr, serverIP, ciaddr net.IP) {
	log.Println("Sending Inform:")

	opts := []dhcp4.Option{
		dhcp4.Option{
			Code:  dhcp4.OptionParameterRequestList,
			Value: []byte{0x1, 0x3, 0x6, 0xf},
		},
		dhcp4.Option{
			Code:  dhcp4.OptionServerIdentifier,
			Value: []byte(serverIP),
		},
	}

	packet := dhcp4.RequestPacket(dhcp4.Inform, mac, ciaddr, randomID(), false, opts)
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
