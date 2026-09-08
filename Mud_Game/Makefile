#Эта команда просто показывает что можно сделать:
help:
	@echo "📋 Доступные команды:"
	@echo "make clean-db        -очистка БД(с подтверждением)"
	@echo "make clean-db-force  -очистка БД(без подтверждения)"
	@echo "make run             -запустить сервер"
	@echo "make git             -отправить изменения в GitHub"
clean-db:
	@echo "🧹 Очистка БД..."
	@read -p "Точно удалить все данные?(y/n) " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		psql -U muduser -d mudgame -h localhost -c "DELETE FROM room_models;"; \
		psql -U muduser -d mudgame -h localhost -c "DELETE FROM npc_models;"; \
		psql -U muduser -d mudgame -h localhost -c "DELETE FROM player_models;"; \
		echo "✅ БД очищена!"; \
	else \
		echo "❌ Отменено"; \
		fi
clean-db-force:
	@echo "🧹 принудительная очистка БД"
	psql -U muduser -d mudgame -h localhost -c "DELETE FROM room_models;"
	psql -U muduser -d mudgame -h localhost -c "DELETE FROM npc_models;"
	psql -U muduser -d mudgame -h localhost -c "DELETE FROM player_models;"
	@echo "✅ БД очищена!"

run:
	-go run cmd/server/main.go || true

git:
	@echo "📦 Отправка на GitHub..."
	git add .
	@read -p "Введите сообщение коммита: " msg; \
	git commit -m "msg"; \
	git push
	@echo "✅ Готово!"

