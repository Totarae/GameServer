package hero

type Warrior struct {
	BaseHero
}

func NewWarrior() Hero {
	return &Warrior{BaseHero{health: 100, speed: 10, class: "Warrior"}}
}
