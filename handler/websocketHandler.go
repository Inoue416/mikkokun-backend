package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketアップグレードの設定
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketコネクションを格納するマップ
var connections = make(map[string]*websocket.Conn)
var clientInformations = make(map[string]string)

const Broadcast = "broadcast"
const Alert = "alert"

// WebSocketメッセージの構造体
type WebSocketRequest struct {
	ActionType string `json:"ActionType"`
	SeatNumber string `json:"SeatNumber"`
}

type AlertMessageType struct {
	TargetSeatNumber string
	Message          string
	TimeLimitSec     int
}

type ResponseMessageOnly struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

func broadcastMessage(myId string, message string) {
	// メッセージをブロードキャスト
	for clientId, c := range connections {
		if clientId != myId {
			if err := c.WriteJSON(ResponseMessageOnly{Message: message, IsSuccess: true}); err != nil {
				break
			}
		}
	}
}

// アラートメッセージを送信
func AlertMessage(myId string, targetSeatNumber string, timeLimitSec int) {
	alert := AlertMessageType{
		TargetSeatNumber: targetSeatNumber,
		Message:          "密告されました。\nタイマーを止めると密告を防ぐことができます。",
		TimeLimitSec:     timeLimitSec,
	}
	errorMessage := "そのユーザーは現在接続していません。"
	successMessave := "ターゲットにアラートを送信しました。"

	for clientId, c := range connections {
		if clientId == targetSeatNumber {
			if err := c.WriteJSON(alert); err != nil {
				connections[myId].WriteJSON(ResponseMessageOnly{
					IsSuccess: false,
					Message:   errorMessage})
				return
			}
			connections[myId].WriteJSON(ResponseMessageOnly{
				IsSuccess: true,
				Message:   successMessave})
			return
		}
	}
	connections[myId].WriteJSON(ResponseMessageOnly{
		IsSuccess: false,
		Message:   errorMessage})
}

// TODO: コネクションが切れた際のブロードキャスト

// TODO: タイムアップ時にブロードキャスト　(これはフロント側なのでこちらでは実装せず、実装ずみのブロードキャスト機能を使う)

// WebSocketハンドラー
func WebsocketHandler(c *gin.Context) {
	// WebSocketのアップグレード
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	fmt.Println("WebsocketHandler...")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer conn.Close()

	// クライアントの識別子を生成
	clientId := c.Query("uuid")
	clientSeatNumber := c.Query("seatnumber")

	// WebSocketコネクションをマップに格納
	clientInformations[clientSeatNumber] = clientId
	connections[clientId] = conn
	conn.WriteJSON(ResponseMessageOnly{
		IsSuccess: true,
		Message: fmt.Sprintf(
			"Success Connection!!\nYour id is %s\nYour seat number is %s\n",
			clientId,
			clientSeatNumber,
		),
	})

	for {
		// メッセージの読み取り
		var request WebSocketRequest
		if err := conn.ReadJSON(&request); err != nil {
			break
		}
		switch request.ActionType {
		case Broadcast:
			broadcastMessage(clientId, "sample")
		case Alert:
			break
		default:
			break
		}
	}

	// コネクションをマップから削除
	delete(connections, clientId)
	delete(clientInformations, clientId)
}
