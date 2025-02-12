package handlers

import (
	"awesomeProject5/game"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все соединения (для упрощения)
	},
}

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var gameState = game.NewGameState()

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при обновлении соединения:", err)
		return
	}
	defer conn.Close()

	player := gameState.AddPlayer()
	player.Conn = conn // Привязываем соединение к игроку!

	// Отправляем начальное состояние игроку
	// скрываем соединение
	conn.WriteJSON(Event{
		Type: "init",
		Payload: map[string]interface{}{
			"id":    player.ID,
			"class": player.Hero.Description(),
			"x":     player.Position.X,
			"y":     player.Position.Y,
		},
	})

	for {
		var event Event
		err := conn.ReadJSON(&event)
		if err != nil {
			log.Println("Ошибка при чтении сообщения:", err)
			break
		}

		switch event.Type {
		case "move":
			// Обработка команды движения
			payload := event.Payload.(map[string]interface{})
			xFloat, xOk := payload["x"].(float64)
			yFloat, yOk := payload["y"].(float64)
			if !xOk || !yOk {
				log.Println("Ошибка: неверные координаты для move")
				continue
			}
			// Приведение float64 → int
			x, y := int(xFloat), int(yFloat)
			err := gameState.MovePlayer(player.ID, x, y)
			if err != nil {
				log.Printf("Ошибка при перемещении игрока %d: %v\n", player.ID, err)
				// TODO: Выдать игроку ошибку адресно
				continue
			}

			// Отправляем обновленное состояние всем игрокам
			broadcastGameState()
		case "attack":
			// Обработка атаки
		default:
			log.Println("Неизвестный тип события:", event.Type)
		}
	}
	//Удаляем игрока при отключении
	gameState.RemovePlayer(player.ID)
	log.Printf("Игрок %d отключился", player.ID)

}

func HeartbeatLoop() {
	ticker := time.NewTicker(10 * time.Second) // Каждые 10 секунд
	defer ticker.Stop()
	for range ticker.C {
		for id, player := range gameState.Players {
			if player.Conn == nil {
				continue
			}
			if err := player.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Отключение клиента %s из-за отсутствия ответа на пинг", id)
				player.Conn.Close()
				gameState.RemovePlayer(id)
				/*delete(clients, id)
				idInt, err := strconv.Atoi(id) // не распарсили id - падаем
				// TODO: протестить этот кусок, могут быть паники или зависания
				if err != nil {
					log.Printf("Ошибка при преобразовании ID игрока %s в int: %v", id, err)
					continue
				}
				gameState.RemovePlayer(idInt)*/
			}
		}
	}
}

// broadcastGameState отправляет текущее состояние игры всем подключенным клиентам.
// TODO: тут можно fanout попробовать, подрезать Connection
// TODO: можно исполдьзовать замыкания для тикающего урона
func broadcastGameState() {
	state := gameState.GetState()
	// Получаем текущее состояние игры

	// Преобразуем состояние в JSON
	stateJSON, err := json.Marshal(state)
	if err != nil {
		log.Println("Ошибка при сериализации состояния игры:", err)
		return
	}

	// Отправляем состояние всем клиентам
	for _, player := range gameState.Players {
		if player.Conn != nil {
			player.Conn.WriteJSON(Event{
				Type:    "update",
				Payload: state,
			})
		}
	}
	log.Println("Состояние игры:", string(stateJSON))
}

func GameLoop() {
	ticker := time.NewTicker(50 * time.Millisecond)
	for range ticker.C {
		// Сам игровой движок
	}
}
