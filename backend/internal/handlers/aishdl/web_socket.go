package aishdl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

const (
	url                       = "wss://stream.aisstream.io/v0/stream"
	messageTypePositionReport = "PositionReport"
)

type SubscriptionMessage struct {
	APIKey        string        `json:"APIKey"`
	BoundingBoxes [][][]float64 `json:"BoundingBoxes"`
}

type AISPacket struct {
	MetaData    map[string]interface{} `json:"MetaData"`
	MessageType string                 `json:"MessageType"`
	Message     AISPacketMessage       `json:"Message"`
}

type AISPacketMessage struct {
	PositionReport *PositionReport `json:"PositionReport,omitempty"`
}

type PositionReport struct {
	MessageID                 int32   `json:"MessageID"`
	RepeatIndicator           int32   `json:"RepeatIndicator"`
	UserID                    int32   `json:"UserID"`
	Valid                     bool    `json:"Valid"`
	NavigationalStatus        int32   `json:"NavigationalStatus"`
	RateOfTurn                int32   `json:"RateOfTurn"`
	Sog                       float64 `json:"Sog"`
	PositionAccuracy          bool    `json:"PositionAccuracy"`
	Longitude                 float64 `json:"Longitude"`
	Latitude                  float64 `json:"Latitude"`
	Cog                       float64 `json:"Cog"`
	TrueHeading               int32   `json:"TrueHeading"`
	Timestamp                 int32   `json:"Timestamp"`
	SpecialManoeuvreIndicator int32   `json:"SpecialManoeuvreIndicator"`
	Spare                     int32   `json:"Spare"`
	Raim                      bool    `json:"Raim"`
	CommunicationState        int32   `json:"CommunicationState"`
}

type WebSocketListener struct {
	apiKey           string
	collectorService ports.CollectorService
	conn             *websocket.Conn
}

func NewWebSocketListener(cfg config.Config, collectorService ports.CollectorService) (*WebSocketListener, error) {
	if strings.TrimSpace(cfg.WebSocketAPIKey) == "" {
		return nil, errors.New("web socket API key must be set")
	}

	return &WebSocketListener{
		apiKey:           cfg.WebSocketAPIKey,
		collectorService: collectorService,
	}, nil
}

func (wsl *WebSocketListener) Listen(ctx context.Context) error {
	var dialErr error
	wsl.conn, _, dialErr = websocket.DefaultDialer.Dial(url, nil)
	if dialErr != nil {
		return fmt.Errorf("error on dialing web socket URL: %w", dialErr)
	}

	subMsgBytes, _ := json.Marshal(SubscriptionMessage{
		APIKey:        wsl.apiKey,
		BoundingBoxes: [][][]float64{{{-90.0, -180.0}, {90.0, 180.0}}}, // bounding box for the entire world
	})
	if err := wsl.conn.WriteMessage(websocket.TextMessage, subMsgBytes); err != nil {
		return fmt.Errorf("error on subscribing to web socket: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := wsl.readAndProcessMessage(); err != nil {
				return err
			}
		}
	}
}

func (wsl *WebSocketListener) Shutdown() {
	err := wsl.conn.Close()
	if err != nil {
		clog.Errorf("error on shutting down webhook listener: %s", err.Error())
	}
}

func (wsl *WebSocketListener) readAndProcessMessage() error {
	_, p, err := wsl.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("error on reading message: %w", err)
	}

	var packet AISPacket
	err = json.Unmarshal(p, &packet)
	if err != nil {
		return fmt.Errorf("error on unmarshalling packet: %w", err)
	}

	var shipName string
	if packetShipName, ok := packet.MetaData["ShipName"]; ok {
		shipName = packetShipName.(string)
	}

	if packet.MessageType == messageTypePositionReport {
		var positionReport PositionReport
		positionReport = *packet.Message.PositionReport

		err := wsl.collectorService.Process(positionReport.UserID, shipName, positionReport.Latitude, positionReport.Longitude)
		if err != nil {
			return fmt.Errorf("error on processing webhook message: %w", err)
		}
	}

	return nil
}
