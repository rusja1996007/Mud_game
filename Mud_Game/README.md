```
my-mud-game/
├── cmd/
│   └── server/
│       └── main.go                    # Только инициализация и запуск
│
├── internal/
│   │
│   ├── domain/                        # ЯДРО БИЗНЕС-ЛОГИКИ
│   │   ├── player/                    # Домен "Игрок" (РАСШИРЕНО)
│   │   │   ├── entity.go              # Структура Player
│   │   │   ├── value_object.go        # Value objects
│   │   │   ├── service.go             # Бизнес-логика игрока
│   │   │   ├── repository.go           # Интерфейс хранилища
│   │   │   │
│   │   │   ├── attributes/             # ← НОВОЕ: Система характеристик
│   │   │   │   ├── core/
│   │   │   │   │   ├── strength.go    # Сила - ХП, тяжелое оружие
│   │   │   │   │   ├── agility.go     # Ловкость - уклонение, меткость
│   │   │   │   │   ├── energy.go      # Энергия - мана, реген
│   │   │   │   │   ├── tracking.go    # Следопытство - поиск ресурсов
│   │   │   │   │   ├── luck.go        # Удача - сохранение вещей, крит
│   │   │   │   │   ├── intelligence.go # Интеллект - изучение, крафт
│   │   │   │   │   └── stamina.go     # Выносливость - реген, сопротивления
│   │   │   │   ├── calculator.go      # Калькулятор бонусов
│   │   │   │   ├── requirements.go    # Проверка требований
│   │   │   │   └── growth.go          # Рост характеристик
│   │   │   │
│   │   │   ├── death_system/           # ← НОВОЕ: Система смерти
│   │   │   │   ├── death_handler.go   # Обработка смерти
│   │   │   │   ├── respawn.go         # Воскрешение
│   │   │   │   ├── item_loss.go       # Потеря предметов
│   │   │   │   ├── item_salvage.go    # Шанс сохранения вещей
│   │   │   │   └── party_resurrection.go # Воскрешение при победе над боссом
│   │   │   │
│   │   │   └── infirmary.go            # ← НОВОЕ: Лазарет
│   │   │       ├── bed.go             # Койка в лазарете
│   │   │       ├── timer.go           # 2-часовой таймер
│   │   │       └── recovery.go        # Восстановление
│   │   │
│   │   ├── location/                   # Домен "Локация"
│   │   │   ├── entity.go
│   │   │   ├── service.go
│   │   │   └── repository.go
│   │   │
│   │   ├── item/                       # Домен "Предмет" (РАСШИРЕНО)
│   │   │   ├── entity.go
│   │   │   ├── weapon.go
│   │   │   ├── armor.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── binding.go              # ← НОВОЕ: Система привязки (BoP/BoE)
│   │   │
│   │   ├── npc/                        # Домен "NPC"
│   │   │   ├── entity.go
│   │   │   ├── merchant.go
│   │   │   ├── quest_giver.go
│   │   │   └── repository.go
│   │   │
│   │   ├── quest/                       # Домен "Квест" (РАСШИРЕНО)
│   │   │   ├── entity.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── condition.go            # ← НОВОЕ: Условия выполнения
│   │   │   ├── reward.go                # ← НОВОЕ: Награды
│   │   │   └── chain.go                 # ← НОВОЕ: Цепочки квестов
│   │   │
│   │   ├── dialog/                      # Домен "Диалог"
│   │   │   ├── entity.go
│   │   │   ├── tree.go
│   │   │   └── service.go
│   │   │
│   │   ├── combat/                       # Домен "Бой" (РАСШИРЕНО)
│   │   │   ├── entity.go
│   │   │   ├── service.go
│   │   │   ├── calculator.go
│   │   │   ├── role_system.go           # ← НОВОЕ: Роли (танк, химер, дд)
│   │   │   │   ├── tank_role.go
│   │   │   │   ├── healer_role.go
│   │   │   │   └── damage_role.go
│   │   │   ├── threat.go                 # ← НОВОЕ: Система агр
│   │   │   ├── party_combat.go           # ← НОВОЕ: Групповой бой
│   │   │   └── boss_mechanics.go         # ← НОВОЕ: Механики боссов
│   │   │       ├── phase_system.go
│   │   │       ├── special_abilities.go
│   │   │       └── enrage_timer.go
│   │   │
│   │   ├── magic/                        # ← НОВОЕ: Магическая система
│   │   │   ├── spell.go
│   │   │   ├── school.go
│   │   │   ├── mana.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   │
│   │   ├── crafting/                     # ← НОВОЕ: Крафтинг (РАСШИРЕНО)
│   │   │   ├── alchemy/                  # Алхимия
│   │   │   │   ├── potion.go
│   │   │   │   ├── recipe.go
│   │   │   │   ├── cauldron.go
│   │   │   │   └── discovery.go
│   │   │   │
│   │   │   ├── herbalism/                 # ← НОВОЕ: Травничество
│   │   │   │   ├── plant.go
│   │   │   │   ├── gathering.go
│   │   │   │   └── zone_herbs.go
│   │   │   │
│   │   │   ├── recipe.go
│   │   │   ├── profession.go
│   │   │   ├── materials.go
│   │   │   ├── workbench.go
│   │   │   └── service.go
│   │   │
│   │   ├── dungeon/                       # ← НОВОЕ: Подземелья
│   │   │   ├── instance.go
│   │   │   ├── lockout.go                 # Блокировка на 30 мин
│   │   │   │   ├── instance_lock.go
│   │   │   │   ├── player_tracking.go
│   │   │   │   ├── queue_system.go
│   │   │   │   └── reset_timer.go
│   │   │   ├── boss.go
│   │   │   │   ├── mechanics.go
│   │   │   │   ├── phases.go
│   │   │   │   └── party_requirements.go
│   │   │   ├── room.go
│   │   │   └── generator.go
│   │   │
│   │   ├── loot/                          # ← НОВОЕ: Система лута
│   │   │   ├── distribution.go            # Need/Greed система
│   │   │   │   ├── need_before_greed.go
│   │   │   │   ├── timer.go               # 1 минута на выбор
│   │   │   │   ├── auto_pass.go
│   │   │   │   └── random_selector.go
│   │   │   ├── boss_loot.go
│   │   │   │   ├── loot_table.go
│   │   │   │   ├── personal_loot.go
│   │   │   │   └── group_loot.go
│   │   │   └── binding.go
│   │   │
│   │   ├── housing/                        # ← НОВОЕ: Жилье и участки
│   │   │   ├── house.go
│   │   │   │   ├── entity.go
│   │   │   │   ├── storage.go
│   │   │   │   ├── upgrade.go
│   │   │   │   └── location.go
│   │   │   ├── land.go
│   │   │   │   ├── plot.go
│   │   │   │   ├── purchase.go
│   │   │   │   ├── farming.go
│   │   │   │   │   ├── crop.go
│   │   │   │   │   ├── growth_timer.go    # 4 часа реального времени
│   │   │   │   │   └── harvest.go
│   │   │   │   └── expansion.go
│   │   │   └── furniture.go
│   │   │
│   │   ├── economy/                        # ← НОВОЕ: Экономика
│   │   │   ├── auction_house.go            # Аукцион
│   │   │   │   ├── listing.go
│   │   │   │   ├── bid.go
│   │   │   │   ├── buyout.go
│   │   │   │   ├── duration.go
│   │   │   │   └── commission.go
│   │   │   ├── currency.go
│   │   │   └── market.go
│   │   │
│   │   ├── social/                         # ← НОВОЕ: Социальная система
│   │   │   ├── party.go                    # Группы игроков
│   │   │   ├── chat.go
│   │   │   └── friendship.go
│   │   │
│   │   └── world/                          # ← НОВОЕ: Мир и путешествия
│   │       ├── travel/
│   │       │   ├── path.go
│   │       │   ├── random_encounter.go     # Случайные встречи (цветы)
│   │       │   ├── travel_time.go
│   │       │   └── discovery.go
│   │       └── weather.go
│   │
│   ├── usecase/                            # СЦЕНАРИИ ИСПОЛЬЗОВАНИЯ (РАСШИРЕНО)
│   │   ├── movement/
│   │   │   ├── move_player.go
│   │   │   └── look_around.go
│   │   │
│   │   ├── combat/
│   │   │   ├── start_combat.go
│   │   │   ├── attack.go
│   │   │   ├── use_ability.go
│   │   │   ├── heal_party.go
│   │   │   ├── manage_threat.go
│   │   │   └── boss_phase_change.go
│   │   │
│   │   ├── death/                           # ← НОВОЕ: Сценарии смерти
│   │   │   ├── handle_death.go
│   │   │   ├── respawn_in_infirmary.go
│   │   │   ├── party_resurrect.go
│   │   │   └── recover_items.go
│   │   │
│   │   ├── interaction/
│   │   │   ├── take_item.go
│   │   │   ├── talk_to_npc.go
│   │   │   └── use_item.go
│   │   │
│   │   ├── trade/
│   │   │   ├── buy_item.go
│   │   │   ├── sell_item.go
│   │   │   └── haggle.go
│   │   │
│   │   ├── quest/
│   │   │   ├── accept_quest.go
│   │   │   ├── complete_quest.go
│   │   │   └── check_progress.go
│   │   │
│   │   ├── housing/                          # ← НОВОЕ: Сценарии жилья
│   │   │   ├── purchase_house.go
│   │   │   ├── upgrade_house.go
│   │   │   ├── store_item.go
│   │   │   ├── withdraw_item.go
│   │   │   ├── buy_land.go
│   │   │   └── plant_crop.go
│   │   │
│   │   ├── auction/                          # ← НОВОЕ: Сценарии аукциона
│   │   │   ├── list_item.go
│   │   │   ├── place_bid.go
│   │   │   ├── buyout_item.go
│   │   │   ├── cancel_auction.go
│   │   │   └── claim_item.go
│   │   │
│   │   ├── dungeon/                          # ← НОВОЕ: Сценарии подземелий
│   │   │   ├── enter_dungeon.go
│   │   │   ├── check_instance_lock.go
│   │   │   ├── join_queue.go
│   │   │   └── reset_instance.go
│   │   │
│   │   └── crafting/                         # ← НОВОЕ: Сценарии крафта
│   │       ├── gather_herb.go
│   │       ├── craft_potion.go
│   │       ├── discover_recipe.go
│   │       └── use_potion.go
│   │
│   ├── delivery/                          # ДОСТАВКА (без изменений)
│   │   ├── tcp/
│   │   │   ├── handler.go
│   │   │   ├── command_parser.go
│   │   │   └── protocol.go
│   │   │
│   │   ├── websocket/
│   │   │   ├── handler.go
│   │   │   └── message.go
│   │   │
│   │   └── api/
│   │       ├── handler.go
│   │       ├── middleware.go
│   │       └── response.go
│   │
│   ├── repository/                         # ХРАНИЛИЩА (РАСШИРЕНО)
│   │   ├── player/
│   │   │   ├── memory_repo.go
│   │   │   ├── postgres_repo.go
│   │   │   └── sqlite_repo.go
│   │   │
│   │   ├── location/
│   │   │   ├── memory_repo.go
│   │   │   ├── json_repo.go
│   │   │   └── yaml_repo.go
│   │   │
│   │   ├── item/
│   │   │   ├── memory_repo.go
│   │   │   └── postgres_repo.go
│   │   │
│   │   ├── npc/
│   │   │   └── memory_repo.go
│   │   │
│   │   ├── quest/
│   │   │   └── memory_repo.go
│   │   │
│   │   ├── spell/                           # ← НОВОЕ
│   │   │   └── memory_repo.go
│   │   │
│   │   ├── recipe/                          # ← НОВОЕ
│   │   │   └── yaml_repo.go
│   │   │
│   │   ├── herb/                            # ← НОВОЕ
│   │   │   └── yaml_repo.go
│   │   │
│   │   ├── dungeon/                         # ← НОВОЕ
│   │   │   └── template_repo.go
│   │   │
│   │   ├── auction/                         # ← НОВОЕ
│   │   │   └── postgres_repo.go
│   │   │
│   │   └── shared/
│   │       ├── database.go
│   │       └── cache.go
│   │
│   └── pkg/                                 # ОБЩИЕ УТИЛИТЫ (РАСШИРЕНО)
│       ├── logger/
│       │   ├── logger.go
│       │   └── zap_logger.go
│       │
│       ├── config/
│       │   ├── loader.go
│       │   └── validator.go
│       │
│       ├── uuid/
│       │   └── generator.go
│       │
│       ├── time/
│       │   ├── game_clock.go
│       │   └── timer_manager.go             # ← НОВОЕ: Управление таймерами
│       │
│       ├── random/                           # ← НОВОЕ: Генерация случайных чисел
│       │   ├── weighted.go
│       │   └── dice.go
│       │
│       └── formula/                          # ← НОВОЕ: Математические формулы
│           ├── calculator.go
│           └── balancer.go
│
├── web/                                    # Веб-клиент (без изменений)
│   ├── static/
│   │   ├── js/
│   │   ├── css/
│   │   └── images/
│   └── index.html
│
├── scripts/                                # Скрипты (РАСШИРЕНО)
│   ├── migrate_db.sh
│   ├── import_world.py
│   ├── backup_world.sh
│   ├── generate_herbs.py                    # ← НОВОЕ: Генерация растений
│   └── balance_checker.py                   # ← НОВОЕ: Проверка баланса
│
├── configs/                                 # КОНФИГИ (РАСШИРЕНО)
│   ├── server.yaml
│   ├── database.yaml
│   ├── game.yaml
│   ├── attributes.yaml                       # ← НОВОЕ: Настройки характеристик
│   ├── combat.yaml                           # ← НОВОЕ: Баланс боя
│   ├── death.yaml                            # ← НОВОЕ: Штрафы при смерти
│   ├── housing.yaml                          # ← НОВОЕ: Цены на дома
│   ├── auction.yaml                          # ← НОВОЕ: Комиссии аукциона
│   ├── dungeons.yaml                         # ← НОВОЕ: Настройки подземелий
│   ├── alchemy.yaml                          # ← НОВОЕ: Рецепты зелий
│   ├── herbalism.yaml                        # ← НОВОЕ: Растения и шансы
│   └── travel.yaml                           # ← НОВОЕ: Время перемещения
│
├── data/                                    # ДАННЫЕ МИРА (РАСШИРЕНО)
│   ├── locations/
│   │   ├── starting_area.yaml
│   │   └── dark_forest.yaml
│   │
│   ├── items/
│   │   ├── weapons.yaml
│   │   ├── armor.yaml
│   │   └── consumables.yaml
│   │
│   ├── npcs/
│   │   ├── merchants.yaml
│   │   ├── quest_givers.yaml
│   │   └── enemies.yaml
│   │
│   ├── quests/
│   │   ├── main_questline.yaml
│   │   └── side_quests.yaml
│   │
│   ├── dialogs/
│   │   ├── village_dialogs.yaml
│   │   └── story_dialogs.yaml
│   │
│   ├── herbs/                                # ← НОВОЕ: Растения по зонам
│   │   ├── starting_area_herbs.yaml
│   │   ├── forest_herbs.yaml
│   │   └── mountain_herbs.yaml
│   │
│   ├── potions/                              # ← НОВОЕ: Рецепты зелий
│   │   ├── healing_potions.yaml
│   │   ├── buff_potions.yaml
│   │   └── utility_potions.yaml
│   │
│   ├── bosses/                               # ← НОВОЕ: Данные боссов
│   │   ├── dungeon1_bosses.yaml
│   │   └── dungeon2_bosses.yaml
│   │
│   ├── loot_tables/                          # ← НОВОЕ: Таблицы лута
│   │   ├── boss_loot.yaml
│   │   └── dungeon_loot.yaml
│   │
│   ├── houses/                               # ← НОВОЕ: Шаблоны домов
│   │   ├── starter_houses.yaml
│   │   ├── medium_houses.yaml
│   │   └── large_houses.yaml
│   │
│   └── spells/                               # ← НОВОЕ: Заклинания
│       ├── fire_magic.yaml
│       ├── healing_magic.yaml
│       └── utility_magic.yaml
│
└── go.mod
```

## 🎮 Особенности игры

- 7 характеристик персонажа (Сила, Ловкость, Энергия, Следопытство, Интеллект, еда, вода)
- Подземелья с боссами и системой блокировки
- Система крафта и алхимии
- Сбор растений и фермерство
- Собственное жилье с возможностью улучшения
- Аукцион для торговли между игроками
- Система смерти с лазаретом и шансом сохранения вещей

 sudo systemctl start postgresql - Запуск сервера PostgreSQL
 sudo systemctl status postgresql - Проверка статуса
 sudo ss -nltp | grep 5432 - Проверяет, открыт ли (порт 5432) для подключений. PostgreSQL ждёт гостей именно на этой двери.
 sudo -u postgres psql - Заходим в PostgreSQL как главный администратор
 CREATE USER muduser WITH PASSWORD 'mudpassword'; - Создаём нового жителя (пользователя) с именем muduser и паролем. Теперь у него есть свой пропуск.
 CREATE DATABASE mudgame OWNER muduser; - Создаём квартиру (базу данных) с названием mudgame и отдаём ключи muduser. Теперь это его личное пространство.
 GRANT ALL PRIVILEGES ON DATABASE mudgame TO muduser; - Говорим "muduser, ты здесь главный - можешь делать всё что хочешь".

  показываем пропуск (пароль) и попадаем внутрь.
 Дом (компьютер)
├── Электричество (сервер PostgreSQL) - включили (start)
├── Дверь (порт 5432) - открыта (ss -nltp)
├── Жильцы:
│   ├── postgres (администратор дома)
│   └── muduser (наш игровой пользователь)
└── Квартиры (базы данных):
    └── mudgame (наша игровая БД) - ключи у muduser

   //psql -U muduser -d mudgame -h localhost - подключение к БД PgAdmin
   //go run cmd/server/main.go - запуск самого сервера
   SELECT * FROM player_models; - вытащить всю таблицу
   DELETE FROM player_models; - удалить всю таблицу, при повторном запуске опять создасться сама



/////////////////////КОМАНДЫ///////////////////
Команда	                       Описание
inventory	               Показать инвентарь и экипировку
wear <предмет>	           Надеть предмет
remove <слот>	           Снять предмет со слота
take / drop	               Взять/бросить предмет
garden / plant / harvest   Работа с огородом
look / move	               Осмотр и перемещение
destroy <предмет>          Уничтожить
stats                      Характеристики
remove <экипировка>        Снять + заменить если уже имеется.
fill empty bottle          Набрать воды у колодца
eat <еда>                  сьесть ..
drink <напиток>            выпить...
yes                        команда для подтверждения (для охоты)
hunt                       отправка на охоту

1. Базовая система (без классов)

    В игре нет классов. Игрок создаёт персонажа с нуля, распределяя четыре характеристики: Сила, Ловкость, Следопытство, Интеллект.

    Нет маны как ресурса.

    Магия существует только в виде одноразовых свитков (найти или создать).

    Перманентная смерть — персонаж удаляется при гибели. Есть малый шанс выжить (привязан к Следопытству).

2. Интеллект (Intellect)

Отвечает за:

    Активацию свитков (шанс успеха, усиление эффекта, мгновенное использование в бою)

    Идентификацию предметов, свитков, магических замков, свойств растений

    Сложные навыки (алхимия, руны, инженерия, создание свитков, древние языки)

    Взлом магических замков (рунные, словесные, временные, энергетические печати)

Доступ к навыкам:

    Простые — без требований

    Средние — Инт 9+

    Сложные — Инт 13+

    Элитные — Инт 16+

Скорость обучения навыков зависит от Интеллекта (выше Инт — быстрее прокачка).
3. Следопытство (Wits / Survival)

Отвечает за:

    Сбор ресурсов (растения, грибы, минералы) — количество, качество, сезонность

    Огород — урожайность, качество плодов, устойчивость к вредителям, семена

    Поиск в пещерах — скрытые тайники, следы монстров, залежи руды и кристаллов

    Навигацию — не теряется в лесу/горах, ведёт группу, находит безопасный лагерь

    Охоту и рыбалку — гарантированная добыча, редкие виды, разделка туш

    Чутьё на опасность — предчувствие засад и ловушек

    Полевую медицину — лечение без алхимической лаборатории



Шанс выжить при смертельном уроне — привязан к Следопытству (инстинктивное чутьё). При успехе персонаж сбегает с 1 HP с поле битвы, и разрушеном или сломаном снаряжении(будет какой то шанс)
4. Сила (Strength)

Отвечает за:

    Здоровье (HP) — чем выше Сила, тем больше HP и регенерация

    Щиты — лёгкие (Сил 6+), средние (10+), тяжёлые (14+), бастионные (18+)

    Тяжёлое оружие — среднее (8+), тяжёлое (12+), осадное (16+), титаническое (18+)

    Бонус к урону в ближнем бою

    

    Физическое взаимодействие — открыть дверь, сдвинуть камень, вырвать решётку, поднять союзника

    Грузоподъёмность — сколько веса может нести без штрафов


    Кузнечное дело и добычу руды

5. Ловкость (Dexterity)

Отвечает за:

    Уклонение (Класс защиты) — базовая защита без брони или в лёгкой броне

    Инициативу в бою — кто ходит первым, возможность двойного действия в первом раунде

    Точность атак — бонус к попаданию для всех видов оружия

    

 

    Взлом механических замков 

   

    Метательное оружие и луки — урон и точность для дистанционного боя

    Точные ремёсла (ювелирка, механика, портняжное дело, ,

