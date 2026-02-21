package room

type Room struct {
	ID          string
	Name        string
	Description string
	Exits       map[string]string //выходы: направление → ID комнаты
}

func (r *Room) GetID() string {
	return r.ID
}

func (r *Room) GetName() string {
	return r.Name
}

func (r *Room) GetDescription() string {
	return r.Description
}

func (r *Room) GetExits() map[string]string {
	return r.Exits
}
func (r *Room) Look(playerID string) string {
	return r.Description
}
func (r *Room) OnEnter() string {
	return "Ты вошел в " + r.Name
}
func (r *Room) OnExit() string {
	return "Ты покинул " + r.Name
}
