package cap

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	device     = "eth0"                       // Change this to your network interface
	forwardURL = "http://example.com/forward" // Change this to your endpoint
	bufferSize = 1600
	timeout    = 5 * time.Second
)

type HTTPStream struct {
	sync.Mutex
	packets map[uint32][]byte
}

func (hs *HTTPStream) AddPacket(seq uint32, payload []byte) {
	hs.Lock()
	defer hs.Unlock()
	hs.packets[seq] = payload
}

func (hs *HTTPStream) GetAndRemovePackets() [][]byte {
	hs.Lock()
	defer hs.Unlock()
	var packets [][]byte
	for seq, payload := range hs.packets {
		packets = append(packets, payload)
		delete(hs.packets, seq)
	}
	return packets
}

func main() {
	// Open the device for capturing
	handle, err := pcap.OpenLive(device, bufferSize, true, timeout)
	if err != nil {
		log.Fatalf("Error opening device: %v", err)
	}
	defer handle.Close()

	// Set filter to capture only TCP packets
	err = handle.SetBPFFilter("tcp")
	if err != nil {
		log.Fatalf("Error setting BPF filter: %v", err)
	}

	// Create a packet source
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	httpStreams := make(map[uint32]*HTTPStream)
	var mu sync.Mutex

	fmt.Println("Starting TCP traffic capture on", device)

	// Loop through packets
	for packet := range packetSource.Packets() {
		// Get the TCP layer from this packet
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)

			// Get the IP layer from this packet
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer != nil {
				ip, _ := ipLayer.(*layers.IPv4)

				// Check if this is an HTTP request
				if tcp.SrcPort == layers.TCPPort(80) || tcp.DstPort == layers.TCPPort(80) {
					mu.Lock()
					streamID := getStreamID(ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort)
					if _, exists := httpStreams[streamID]; !exists {
						httpStreams[streamID] = &HTTPStream{packets: make(map[uint32][]byte)}
					}
					httpStreams[streamID].AddPacket(tcp.Seq, tcp.Payload)
					mu.Unlock()

					// Check if we have a complete HTTP request
					if tcp.FIN {
						mu.Lock()
						payloads := httpStreams[streamID].GetAndRemovePackets()
						mu.Unlock()

						var fullPayload bytes.Buffer
						for _, payload := range payloads {
							fullPayload.Write(payload)
						}

						// Forward the HTTP payload to the specified endpoint
						go forwardHTTPPayload(fullPayload.Bytes())
					}
				}
			}
		}
	}
}

func getStreamID(srcIP net.IP, srcPort layers.TCPPort, dstIP net.IP, dstPort layers.TCPPort) uint32 {
	// Simple hash function to generate a stream ID
	hash := hex.EncodeToString([]byte(fmt.Sprintf("%s:%d->%s:%d", srcIP, srcPort, dstIP, dstPort)))
	return uint32(hash[0]) + uint32(hash[1])<<8 + uint32(hash[2])<<16 + uint32(hash[3])<<24
}

func forwardHTTPPayload(payload []byte) {
	resp, err := http.Post(forwardURL, "application/octet-stream", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error forwarding HTTP payload: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing HTTP response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to forward HTTP payload: %s", resp.Status)
	}
}
