package hero

type Thief struct {
	BaseHero
}

func NewThief() Hero {
	return &Thief{BaseHero{health: 80, speed: 16, class: "Thief"}}
}
