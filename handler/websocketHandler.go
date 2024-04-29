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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketコネクションを格納するマップ
var connections = make(map[string]*websocket.Conn)

const Broadcast = "broadcast"
const Alert = "alert"

// WebSocketメッセージの構造体
type WebSocketRequest struct {
	ActionType string `json:"ActionType"`
	SeatNumber string `json:"TargetSeatNumber"`
}

type AlertMessageType struct {
	ActionType       string
	TargetSeatNumber string
	Message          string
	TimeLimitSec     int
}

type ResponseMessageOnly struct {
	ActionType string `json:"ActionType"`
	Message    string `json:"Message"`
}

func broadcastMessage(message string) {
	// メッセージをブロードキャスト
	for idx, c := range connections {
		println("Index (SeatNumber): %s", idx)
		if err := c.WriteJSON(ResponseMessageOnly{Message: message, ActionType: "broadcast"}); err != nil {
			println("Error: %v", err)
			c.Close()
			delete(connections, idx)
			break
		}
	}
}

// アラートメッセージを送信
func AlertMessage(myId string, targetSeatNumber string, timeLimitSec int) {
	alert := AlertMessageType{
		ActionType:       "alert",
		TargetSeatNumber: targetSeatNumber,
		Message:          "密告されました。\nタイマーを止めると密告を防ぐことができます。",
		TimeLimitSec:     timeLimitSec,
	}
	successMessage := "ターゲットにアラートを送信しました。"

	for clientId, c := range connections {
		if clientId == targetSeatNumber {
			if err := c.WriteJSON(alert); err != nil {
				fmt.Println("Error: ", err)
				return
			}
			connections[myId].WriteJSON(ResponseMessageOnly{
				ActionType: "broadcast",
				Message:    successMessage})
			return
		}
	}
}

// TODO: コネクションが切れた際のブロードキャスト

// TODO: タイムアップ時にブロードキャスト　(これはフロント側なのでこちらでは実装せず、実装ずみのブロードキャスト機能を使う)

// リクエストがあったシート番号がすでにないかを確認
// ある場合はエラーメッセージを返す
// func isExistSeatNumber(seatNumber string) bool {
// 	for seatNum, _ := range connections {
// 		if seatNum == seatNumber {
// 			fmt.Printf("Exist seat number: %s\n", seatNumber)
// 			fmt.Printf("Compare seat number: %s\n", seatNum)
// 			return true
// 		}
// 	}
// 	fmt.Printf("Is not exist seat number : %s\n", seatNumber)
// 	return false
// }

// WebSocketハンドラー
func WebsocketHandler(c *gin.Context) {
	// WebSocketのアップグレード
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	fmt.Println("*** WebsocketHandler ***")
	if err != nil {
		fmt.Printf("Can not upgrade: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed"})
		return
	}

	// クライアントの識別子を生成
	clientSeatNumber := c.Query("seatnumber")

	// WebSocketコネクションをマップに格納
	connections[clientSeatNumber] = conn
	conn.WriteJSON(ResponseMessageOnly{
		ActionType: "broadcast",
		Message: fmt.Sprintf(
			"Success Connection!!\nYour seat number is %s\n",
			clientSeatNumber,
		),
	})

	for {
		// メッセージの読み取り
		var request WebSocketRequest
		if err := conn.ReadJSON(&request); err != nil {
			break
		}
		fmt.Println("*** Request Data ***")
		fmt.Println(request.ActionType)
		fmt.Println(request.SeatNumber)
		fmt.Println("*********")
		switch request.ActionType {
		case Broadcast:
			fmt.Println("--- Broadcast case ---")
			broadcastMessage("sample")
		case Alert:
			AlertMessage(clientSeatNumber, request.SeatNumber, 0)
		default:
			break
		}
	}
	conn.Close()
	delete(connections, clientSeatNumber)
}
