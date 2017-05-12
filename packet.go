package main

import (
	"fmt"

	"github.com/onesimus-systems/dhcp4"
)

func printPacket(p dhcp4.Packet) {
	// Standard BOOTP fields
	fmt.Printf("  OpCode: %d\n", p.OpCode())
	fmt.Printf("  HType: %d\n", p.HType())
	fmt.Printf("  HLen: %d\n", p.HLen())
	fmt.Printf("  Hops: %d\n", p.Hops())
	fmt.Printf("  XID: %s\n", string(p.XId()))
	fmt.Printf("  Secs: %d\n", p.Secs())
	fmt.Printf("  CIAddr: %s\n", p.CIAddr().String())
	fmt.Printf("  YIAddr: %s\n", p.YIAddr().String())
	fmt.Printf("  SIAddr: %s\n", p.SIAddr().String())
	fmt.Printf("  GIAddr: %s\n", p.GIAddr().String())
	fmt.Printf("  CHAddr: %s\n", p.CHAddr().String())

	// Flags
	fmt.Println("  Flags:")
	fmt.Printf("    Broadcast: %t\n", p.Broadcast())

	// Options
	fmt.Println("  Options:")
	opts := p.ParseOptions()
	for code, val := range opts {
		if code == dhcp4.OptionDHCPMessageType {
			val = []byte(dhcp4.MessageType(val[0]).String())
			fmt.Printf("    %s: %s\n", code, val)
			continue
		}
		fmt.Printf("    %s: %v\n", code, val)
	}
}
