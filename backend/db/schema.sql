CREATE TABLE IF NOT EXISTS shelters (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    name           TEXT    NOT NULL,
    address        TEXT    NOT NULL,
    capacity_total INTEGER NOT NULL CHECK (capacity_total >= 0),
    capacity_used  INTEGER NOT NULL DEFAULT 0 CHECK (capacity_used >= 0)
);

CREATE TABLE IF NOT EXISTS responders (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    name   TEXT NOT NULL,
    skill  TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'available'
           CHECK (status IN ('available', 'assigned', 'off_duty'))
);

CREATE TABLE IF NOT EXISTS incidents (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    type            TEXT    NOT NULL
                    CHECK (type IN ('flood', 'fire', 'storm', 'medical', 'other')),
    location        TEXT    NOT NULL,
    description     TEXT    NOT NULL DEFAULT '',
    people_affected INTEGER NOT NULL CHECK (people_affected BETWEEN 1 AND 500),
    vulnerable      INTEGER NOT NULL DEFAULT 0 CHECK (vulnerable IN (0, 1)),
    severity        INTEGER NOT NULL CHECK (severity BETWEEN 1 AND 5),
    priority_score  INTEGER NOT NULL,
    priority_band   TEXT    NOT NULL
                    CHECK (priority_band IN ('Critical', 'High', 'Medium', 'Low')),
    status          TEXT    NOT NULL DEFAULT 'registered'
                    CHECK (status IN ('registered', 'assessed', 'assigned', 'in_progress', 'resolved')),
    shelter_id      INTEGER REFERENCES shelters (id),
    responder_id    INTEGER REFERENCES responders (id),
    reported_at     TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS audit (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    at          TEXT NOT NULL,
    actor       TEXT NOT NULL,
    action      TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id   INTEGER NOT NULL,
    detail      TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_incidents_priority ON incidents (priority_score DESC, id);
CREATE INDEX IF NOT EXISTS idx_incidents_status ON incidents (status);
CREATE INDEX IF NOT EXISTS idx_audit_at ON audit (at DESC, id DESC);
