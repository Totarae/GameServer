package game

import "testing"

func TestAddPlayer(t *testing.T) {
	gs := NewGameState()
	player := gs.AddPlayer()

	if player.ID == 0 {
		t.Error("ID игрока не должен быть 0")
	}

	if len(gs.Players) != 1 {
		t.Error("Игрок не был добавлен в состояние игры")
	}
}

func TestMovePlayer(t *testing.T) {
	gs := NewGameState()
	player := gs.AddPlayer()

	gs.MovePlayer(player.ID, 10, 20)

	if player.Position.X != 10 || player.Position.Y != 20 {
		t.Error("Позиция игрока не была обновлена")
	}
}

func TestRemovePlayer(t *testing.T) {
	gs := NewGameState()
	player := gs.AddPlayer()

	gs.RemovePlayer(player.ID)

	if len(gs.Players) != 0 {
		t.Error("Игрок не был удален из состояния игры")
	}
}
