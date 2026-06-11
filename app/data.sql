PRAGMA foreign_keys = ON;

CREATE TABLE masters (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id         INTEGER UNIQUE NOT NULL,
    name                TEXT,
    phone               TEXT,
    description         TEXT,
    subscription_until  TEXT,
    tariff_id           INTEGER,
    is_active           INTEGER DEFAULT 1,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at          TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_masters_telegram_id ON masters(telegram_id);


CREATE TABLE clients (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id         INTEGER UNIQUE,
    name                TEXT,
    phone               TEXT,
    no_show_count       INTEGER DEFAULT 0,
    is_blocked          INTEGER DEFAULT 0,
    blocked_until       TEXT,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at          TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clients_telegram_id ON clients(telegram_id);


CREATE TABLE master_ref_codes (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    master_id           INTEGER,
    short_id            TEXT UNIQUE NOT NULL,
    qr_code_url         TEXT,
    is_active           INTEGER DEFAULT 1,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE
);

CREATE INDEX idx_ref_codes_short_id ON master_ref_codes(short_id);
CREATE INDEX idx_ref_codes_master ON master_ref_codes(master_id);


CREATE TABLE master_client_bindings (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    master_id           INTEGER,
    client_id           INTEGER,
    binded_at           TEXT DEFAULT CURRENT_TIMESTAMP,
    bind_type           TEXT,
    is_active           INTEGER DEFAULT 1,

    UNIQUE(master_id, client_id),

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

CREATE INDEX idx_bindings_master ON master_client_bindings(master_id);
CREATE INDEX idx_bindings_client ON master_client_bindings(client_id);


CREATE TABLE services (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    master_id           INTEGER,
    name                TEXT NOT NULL,
    duration_minutes    INTEGER NOT NULL,
    price               REAL,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE
);

CREATE INDEX idx_services_master ON services(master_id);


CREATE TABLE schedule_slots (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    master_id           INTEGER,
    client_id           INTEGER,
    service_id          INTEGER,
    start_time          TEXT NOT NULL,
    end_time            TEXT NOT NULL,
    status              TEXT DEFAULT 'pending_confirmation',
    details             TEXT,
    confirm_code        TEXT,
    confirm_deadline    TEXT,
    checkin_time        TEXT,
    cancelled_by        TEXT,
    cancelled_at        TEXT,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at          TEXT DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE SET NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE SET NULL
);

CREATE INDEX idx_slots_master ON schedule_slots(master_id);
CREATE INDEX idx_slots_client ON schedule_slots(client_id);
CREATE INDEX idx_slots_start_time ON schedule_slots(start_time);
CREATE INDEX idx_slots_status ON schedule_slots(status);


CREATE TABLE working_hours (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    master_id           INTEGER,
    day_of_week         INTEGER NOT NULL,
    time_start          TEXT NOT NULL,
    time_end            TEXT NOT NULL,
    break_start         TEXT,
    break_end           TEXT,
    is_day_off          INTEGER DEFAULT 0,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at          TEXT DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(master_id, day_of_week),

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE
);

CREATE INDEX idx_working_hours_master ON working_hours(master_id);


CREATE TABLE tariffs (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    name                TEXT NOT NULL,
    price               REAL NOT NULL,
    meeting_limit       INTEGER,
    client_limit        INTEGER,
    duration_days       INTEGER DEFAULT 30,
    is_active           INTEGER DEFAULT 1,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE subscription_payments (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    master_id           INTEGER,
    tariff_id           INTEGER,
    amount              REAL NOT NULL,
    payment_method      TEXT,
    payment_id          TEXT,
    status              TEXT DEFAULT 'pending',
    paid_at             TEXT,
    valid_from          TEXT,
    valid_until         TEXT,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE,
    FOREIGN KEY (tariff_id) REFERENCES tariffs(id)
);

CREATE INDEX idx_payments_master ON subscription_payments(master_id);
CREATE INDEX idx_payments_status ON subscription_payments(status);


CREATE TABLE no_show_settings (
    master_id           INTEGER PRIMARY KEY,
    enable_penalty      INTEGER DEFAULT 1,
    penalty_percent     INTEGER DEFAULT 50,
    no_show_limit       INTEGER DEFAULT 3,
    block_days          INTEGER DEFAULT 30,
    confirm_minutes     INTEGER DEFAULT 15,
    checkin_method      TEXT DEFAULT 'manual',
    remind_minutes      INTEGER DEFAULT 60,
    updated_at          TEXT DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE
);


CREATE TABLE master_blacklist (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    master_id           INTEGER,
    client_id           INTEGER,
    reason              TEXT,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(master_id, client_id),

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

CREATE INDEX idx_blacklist_master ON master_blacklist(master_id);


CREATE TABLE blocked_slots (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    master_id           INTEGER,
    start_time          TEXT NOT NULL,
    end_time            TEXT NOT NULL,
    reason              TEXT,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE
);

CREATE INDEX idx_blocked_slots_master ON blocked_slots(master_id);
CREATE INDEX idx_blocked_slots_time ON blocked_slots(start_time, end_time);


CREATE TABLE notifications_log (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    user_type           TEXT NOT NULL,
    user_id             INTEGER NOT NULL,
    telegram_id         INTEGER NOT NULL,
    notification_type   TEXT,
    message             TEXT,
    status              TEXT DEFAULT 'sent',
    error_message       TEXT,
    created_at          TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications_log(user_type, user_id);
CREATE INDEX idx_notifications_created ON notifications_log(created_at);


INSERT INTO tariffs (
    name,
    price,
    meeting_limit,
    client_limit,
    duration_days,
    is_active
) VALUES
('Базовый', 300.00, 30, 10, 30, 1),
('Про', 700.00, 100, 50, 30, 1),
('Безлимит', 1500.00, NULL, NULL, 30, 1);