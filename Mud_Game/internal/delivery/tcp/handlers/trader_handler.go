package handlers

import "Mud_game/Mud_Game/internal/repository/trader"

// чтобы работать с одним репозиторием?
type TraderHandler struct {
	traderRepo *trader.PostgresTraderRepository
}

// нужен ли?
func NewTraderHandler(repo *trader.PostgresTraderRepository) *TraderHandler {
	return &TraderHandler{traderRepo: repo}
}
