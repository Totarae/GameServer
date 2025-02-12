package hero

type Mage struct {
	BaseHero
}

func NewMage() Hero {
	return &Mage{BaseHero{health: 80, speed: 8, class: "Mage"}}
}
