package combat

// тип урона
type DamageType string

const (
	DamagePhysical DamageType = "physical"
	DamageFire     DamageType = "fire"
	DamagePoison   DamageType = "poison"
	DamageMagic    DamageType = "magic"
)
