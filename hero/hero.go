package hero

import "fmt"

// Hero - интерфейс героя
// Напилил базовую структуру, без аттрибутов (сила, ловкость, инт) и доп. пропов (дальность атаки, тип оружия)
// TODO: можно обернуть BaseHero в тип атаки DistantCombat, CloseCombat
type Hero interface {
	GetHealth() int
	GetSpeed() int
	Description() string
}

// BaseHero - базовая структура героя
type BaseHero struct {
	health int
	speed  int
	class  string
}

func (h *BaseHero) GetHealth() int {
	return h.health
}

func (h *BaseHero) GetSpeed() int {
	return h.speed
}

func (h *BaseHero) Description() string {
	return fmt.Sprintf("%s (Health: %d, Speed: %d)", h.class, h.health, h.speed)
}
