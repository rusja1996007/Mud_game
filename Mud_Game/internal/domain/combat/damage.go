package combat

// тип урона damage.go нужен для атак монстров!
type DamageType string

const (
	DamagePhysical DamageType = "physical"
	DamageFire     DamageType = "fire"
	DamagePoison   DamageType = "poison"
	DamageMagic    DamageType = "magic"
)
