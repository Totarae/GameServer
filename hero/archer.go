package hero

type Archer struct {
	BaseHero
}

func NewArcher() Hero {
	return &Archer{BaseHero{health: 90, speed: 12, class: "Archer"}}
}
