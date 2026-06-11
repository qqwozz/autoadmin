.PHONY: env-up env-down env-rebuild env-logs env-shell db-shell db-test db-status db-logs db-schema db-tables db-counts db-reset db-clean build build-api build-bot run-api run-bot test

env-up:
	docker-compose up -d --build
	$(MAKE) db-test

env-down:
	docker-compose down

env-rebuild:
	docker-compose down
	docker-compose up -d --build
	$(MAKE) db-test

env-logs:
	docker-compose logs -f

env-shell:
	docker exec -it autoadmin-api sh

db-shell:
	docker exec -it autoadmin-api sqlite3 /data/database.sqlite

db-status:
	docker ps -a | grep autoadmin || true

db-logs:
	docker logs autoadmin-api

db-schema:
	docker exec -it autoadmin-api sqlite3 /data/database.sqlite ".schema"

db-tables:
	docker exec -it autoadmin-api sqlite3 /data/database.sqlite ".tables"

db-counts:
	docker exec -it autoadmin-api sh -c 'sqlite3 /data/database.sqlite "SELECT '\''masters'\'' AS table_name, COUNT(*) FROM masters UNION ALL SELECT '\''clients'\'', COUNT(*) FROM clients UNION ALL SELECT '\''tariffs'\'', COUNT(*) FROM tariffs UNION ALL SELECT '\''services'\'', COUNT(*) FROM services UNION ALL SELECT '\''schedule_slots'\'', COUNT(*) FROM schedule_slots;"'

db-test:
	docker exec -it autoadmin-api sh -c 'sqlite3 /data/database.sqlite "PRAGMA foreign_keys = ON; \
	SELECT '\''CLEAN TEST DATA'\''; \
	DELETE FROM schedule_slots WHERE start_time = '\''2026-01-01 10:00:00'\''; \
	DELETE FROM services WHERE name = '\''Test Service'\''; \
	DELETE FROM master_client_bindings WHERE master_id IN (SELECT id FROM masters WHERE telegram_id = 100001); \
	DELETE FROM clients WHERE telegram_id = 200001; \
	DELETE FROM masters WHERE telegram_id = 100001; \
	SELECT '\''TEST 1: tables exist'\''; \
	SELECT name FROM sqlite_master WHERE type='\''table'\'' ORDER BY name; \
	SELECT '\''TEST 2: tariffs count'\''; \
	SELECT COUNT(*) FROM tariffs; \
	SELECT '\''TEST 3: insert master'\''; \
	INSERT INTO masters (telegram_id, name, phone) VALUES (100001, '\''Test Master'\'', '\''+100000000'\''); \
	SELECT id, telegram_id, name FROM masters WHERE telegram_id = 100001; \
	SELECT '\''TEST 4: insert client'\''; \
	INSERT INTO clients (telegram_id, name, phone) VALUES (200001, '\''Test Client'\'', '\''+200000000'\''); \
	SELECT id, telegram_id, name FROM clients WHERE telegram_id = 200001; \
	SELECT '\''TEST 5: insert binding'\''; \
	INSERT INTO master_client_bindings (master_id, client_id, bind_type) \
	SELECT m.id, c.id, '\''test'\'' FROM masters m, clients c WHERE m.telegram_id = 100001 AND c.telegram_id = 200001; \
	SELECT * FROM master_client_bindings WHERE master_id = (SELECT id FROM masters WHERE telegram_id = 100001); \
	SELECT '\''TEST 6: insert service'\''; \
	INSERT INTO services (master_id, name, duration_minutes, price) \
	SELECT id, '\''Test Service'\'', 60, 500.00 FROM masters WHERE telegram_id = 100001; \
	SELECT * FROM services WHERE name = '\''Test Service'\''; \
	SELECT '\''TEST 7: insert schedule slot'\''; \
	INSERT INTO schedule_slots (master_id, client_id, service_id, start_time, end_time, status) \
	SELECT m.id, c.id, s.id, '\''2026-01-01 10:00:00'\'', '\''2026-01-01 11:00:00'\'', '\''confirmed'\'' \
	FROM masters m, clients c, services s \
	WHERE m.telegram_id = 100001 AND c.telegram_id = 200001 AND s.name = '\''Test Service'\'' LIMIT 1; \
	SELECT * FROM schedule_slots WHERE start_time = '\''2026-01-01 10:00:00'\''; \
	SELECT '\''ALL TESTS DONE'\'';\""'

db-reset:
	docker-compose down -v
	docker-compose up -d --build
	$(MAKE) db-test

db-clean:
	docker-compose down -v

build:
	cd app && CGO_ENABLED=1 go build -o ../autoadmin-api ./cmd/api && CGO_ENABLED=1 go build -o ../autoadmin-bot ./cmd/bot

build-api:
	cd app && CGO_ENABLED=1 go build -o ../autoadmin-api ./cmd/api

build-bot:
	cd app && CGO_ENABLED=1 go build -o ../autoadmin-bot ./cmd/bot

run-api: build-api
	./autoadmin-api

run-bot: build-bot
	./autoadmin-bot

test:
	cd app && go test ./internal/... -v
