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

const BroadcastType = "broadcast"
const AlertType = "alert"
const TimeupType = "timeup"
const LogoutType = "logout"

// TODO: Debug
const TIMELITMISEC = 300

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

func TimeupBroadcast(targetSeatNumber string) {
	message := "座席番号 " + targetSeatNumber + " は居眠りをしています！！！"
	for clientId, c := range connections {
		if clientId == targetSeatNumber {
			continue
		}
		if err := c.WriteJSON(ResponseMessageOnly{
			ActionType: "broadcast",
			Message:    message,
		}); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}

// リクエストがあったシート番号がすでにないかを確認
// ある場合はエラーメッセージを返す
func isExistSeatNumber(seatNumber string) bool {
	for seatNum, _ := range connections {
		if seatNum == seatNumber {
			fmt.Printf("Exist seat number: %s\n", seatNumber)
			fmt.Printf("Compare seat number: %s\n", seatNum)
			return true
		}
	}
	fmt.Printf("Is not exist seat number : %s\n", seatNumber)
	return false
}

type CheckSameSeatNumberRepspose struct {
	IsExists bool   `json:"isExists"`
	Message  string `json:"message"`
}

// CheckSameSeatNumber godoc
// @Summary      受け取った座席番号がすでにないかを確認
// @Description  使用されていればtrue、使用されていなければfalseを返す (およびメッセージ)
// @Accept       json
// @Produce      json
// @Param        seatnumber  query  string  true  "Seat Number"
// @Success      200 {object} CheckSameSeatNumberRepspose
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /checkSameSeatNumber [get]
func CheckSameSeatNumber(c *gin.Context) {
	seatNumber := c.Query("seatnumber")
	if seatNumber == "" {
		c.JSON(http.StatusOK, CheckSameSeatNumberRepspose{
			IsExists: true,
			Message:  "座席を指定して下さい。",
		})
	}
	if isExistSeatNumber(seatNumber) {
		c.JSON(http.StatusOK, CheckSameSeatNumberRepspose{
			IsExists: true,
			Message:  "すでにこの座席番号は使用されています。",
		})
		return
	}
	c.JSON(http.StatusOK, CheckSameSeatNumberRepspose{
		IsExists: false,
		Message:  "この座席番号は使用可能です。",
	})
}

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
		Message:    "Success Connection!!",
	})

	for {
		// メッセージの読み取り
		var request WebSocketRequest
		var logoutFlag bool
		if err := conn.ReadJSON(&request); err != nil {
			break
		}
		logoutFlag = false
		switch request.ActionType {
		case BroadcastType:
			broadcastMessage("Broadcast: sample")
		case AlertType:
			AlertMessage(clientSeatNumber, request.SeatNumber, TIMELITMISEC)
		case TimeupType:
			TimeupBroadcast(request.SeatNumber)
		case LogoutType:
			fmt.Println("Action: Logout")
			logoutFlag = true
		default:
			break
		}
		if logoutFlag {
			break
		}
	}
	conn.Close()
	delete(connections, clientSeatNumber)
	fmt.Println("--- Close Connection ---")
}
