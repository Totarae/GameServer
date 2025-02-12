package game

import (
	"awesomeProject5/hero"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"sync"
	"time"
)

// TODO: добавить какую-то авториазцию token, xApiKey
type Player struct {
	ID       int
	Hero     hero.Hero
	Position struct {
		X, Y int
	}
	Conn *websocket.Conn
}

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// состояние мира
type GameState struct {
	Players map[int]*Player
	mu      sync.Mutex
}

func NewGameState() *GameState {
	return &GameState{
		Players: make(map[int]*Player),
		// mu инициализируется автоматически
	}
}

// Проверяет, заняты ли координаты (x, y) другими игроками
func (gs *GameState) isPositionOccupied(x, y int) bool {
	for _, p := range gs.Players {
		if p.Position.X == x && p.Position.Y == y {
			return true
		}
	}
	return false
}

// добавить нового игрока
func (gs *GameState) AddPlayer() *Player {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	rand.Seed(time.Now().UnixNano())

	// Случайный выбор класса
	var h hero.Hero
	switch rand.Intn(4) {
	case 0:
		h = hero.NewWarrior()
	case 1:
		h = hero.NewMage()
	case 2:
		h = hero.NewArcher()
	case 3:
		h = hero.NewThief()
	}

	player := &Player{
		ID:   rand.Intn(1000),
		Hero: h,
	}
	// TODO: добавить ретраи
	for {

		// TODO: добавить безопасную дистанцию спауна
		player.Position.X = rand.Intn(100)
		player.Position.Y = rand.Intn(100)

		// Проверяем, заняты ли координаты
		if !gs.isPositionOccupied(player.Position.X, player.Position.Y) {
			break // Выбрали свободное место
		}
	}

	gs.Players[player.ID] = player
	return player
}

// удалит игрока
func (gs *GameState) RemovePlayer(playerID int) {
	// TODO: lock-uinlock дублирование, нужен защищенный режим
	gs.mu.Lock()
	defer gs.mu.Unlock()

	delete(gs.Players, playerID)
}

// обновляет позицию игрока.
func (gs *GameState) MovePlayer(playerID int, x, y int) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	player, exists := gs.Players[playerID]
	if !exists {
		return fmt.Errorf("игрок с ID %d не найден", playerID)
	}

	// Проверяем, заняты ли новые координаты
	if gs.isPositionOccupied(x, y) {
		return fmt.Errorf("координаты (%d, %d) уже заняты", x, y)
	}

	// Обновляем позицию игрока
	player.Position.X = x
	player.Position.Y = y

	return nil
}

// возвращает текущее состояние игры.
func (gs *GameState) GetState() map[int]*Player {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	state := make(map[int]*Player)
	// Conn и Hero надо скрывать наверно
	for id, player := range gs.Players {
		state[id] = &Player{
			ID:   player.ID,
			Hero: player.Hero,
			Position: struct {
				X, Y int
			}{
				X: player.Position.X,
				Y: player.Position.Y,
			},
		}
	}
	return state
}
