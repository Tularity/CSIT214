# CSIT214 — Emergency Support Coordination System

Group project by **Tularity** for CoastLink Council, a fictional coastal council.

During a local emergency such as a flood or a fire, the council currently
coordinates by phone, email and spreadsheets. This system links resident
requests, shelters and responders in one place, so a coordinator can see which
incident to deal with first, place affected people somewhere that still has
room, dispatch someone who is actually free, and leave a record of every step.

The prototype concentrates on one path, end to end:

> register an incident → the system scores and bands it → place the affected
> people in a shelter with remaining capacity → dispatch an available responder
> → resolve, releasing the capacity and the responder

## Scope

The council asked for nine capability areas. The interface design covers all
nine. The runnable prototype implements only the subset one working chain
needs, rather than nine that are half built.

## Layout

| Path | Contents |
| --- | --- |
| `backend/` | Go API server, standard library HTTP, SQLite via `modernc.org/sqlite` |
| `frontend/` | React and TypeScript client, built with Vite |
| `seed/` | Sample data loaded into an empty database on first start |

## Screens

| Capability area | Screens | In the prototype |
| --- | --- | --- |
| Incident registration | New incident form | Yes |
| Assessment and prioritisation | Incident list, incident detail | Yes |
| Shelter capacity | Shelter list, shelter detail | List only |
| Responder dispatch | Responder list, roster | List only |
| Resource allocation | Equipment inventory, allocation request | No |
| Progress tracking | Incident timeline | Yes |
| Operations dashboard | Dashboard | Yes |
| Reporting | Report builder | No |
| Audit trail | Activity log | Yes, read only |

The shell is a header plus left navigation holding Dashboard, Incidents,
Shelters, Responders and Activity log. Every list needs three states, not just
the populated one: loading, empty, and failed to load. Priority bands are shown
as colour **and** text, never colour alone.

## Running in development

Go 1.23 or newer and Node 20 or newer. There is no database server to install:
the backend keeps its data in a local SQLite file.

**Backend**

```sh
cd backend
go run .
curl http://localhost:8080/api/health     # {"status":"ok"}
```

On first start it creates `backend/data.db`, applies the schema and loads
`seed/seed.sql`. Delete that file to start again from the sample data.

**Frontend**

```sh
cd frontend
npm install
npm run dev                               # http://localhost:5173
```

The dev server proxies `/api` to the backend, so start the backend first.

## Configuration

Every setting has a working default; the prototype starts with no environment
variables set at all.

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | Port the API listens on |
| `CSIT214_DB_PATH` | `data.db` | SQLite file location |
| `CSIT214_SEED_PATH` | `seed/seed.sql`, then `../seed/seed.sql` | Sample data used when the database is empty |
| `CSIT214_API_URL` | `http://localhost:8080` | Proxy target for the frontend dev server |

## API

Agreed between backend and frontend on 2026-09-19. Both sides build against
this table; if something here turns out to be wrong, change it here first.

Times are stored and returned as UTC in ISO 8601 (`2026-09-19T03:42:00Z`); the
frontend converts to local time for display.

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/api/health` | Liveness probe, returns `{"status":"ok"}` |
| `GET` | `/api/incidents` | List incidents ordered by `priority_score` descending |
| `POST` | `/api/incidents` | Create an incident. The server derives `priority_score` and `priority_band`, and sets `status` to `registered` |
| `GET` | `/api/incidents/{id}` | Single incident |
| `POST` | `/api/incidents/{id}/assess` | Move `status` to `assessed` |
| `POST` | `/api/incidents/{id}/assign-shelter` | Body `{"shelter_id": n}`. Adds `people_affected` to the shelter's `capacity_used`, moves `status` to `assigned` |
| `POST` | `/api/incidents/{id}/assign-responder` | Body `{"responder_id": n}`. Sets the responder to `assigned`, moves `status` to `in_progress` |
| `POST` | `/api/incidents/{id}/resolve` | Moves `status` to `resolved`, releases the shelter capacity and returns the responder to `available` |
| `GET` | `/api/shelters` | Shelters including remaining capacity |
| `GET` | `/api/responders` | Responders and their current status |
| `GET` | `/api/audit` | Audit entries, newest first |
| `GET` | `/api/dashboard` | Counts per status, counts per priority band, shelter occupancy, number of available responders |

Every write operation appends a row to `audit`. The client asked for this
explicitly, so it is not optional.

### Incident fields

| Field | Type | Notes |
| --- | --- | --- |
| `id` | int | |
| `type` | string | `flood`, `fire`, `storm`, `medical`, `other` |
| `location` | string | Street address or landmark |
| `description` | string | |
| `people_affected` | int | 1–500 |
| `vulnerable` | bool | Elderly, children or people with reduced mobility are involved |
| `severity` | int | 1–5 |
| `priority_score` | int | Derived, see below |
| `priority_band` | string | `Critical`, `High`, `Medium`, `Low` |
| `status` | string | `registered` → `assessed` → `assigned` → `in_progress` → `resolved` |
| `shelter_id` | int or null | |
| `responder_id` | int or null | |
| `reported_at` | string | ISO 8601, UTC |

`shelters`: `id`, `name`, `address`, `capacity_total`, `capacity_used`

`responders`: `id`, `name`, `skill`, `status` (`available`, `assigned`, `off_duty`)

`audit`: `id`, `at`, `actor`, `action`, `entity_type`, `entity_id`, `detail`

### Priority formula

```
priority_score = severity * 10
               + min(people_affected, 30)
               + (vulnerable ? 20 : 0)
```

| Score | Band |
| --- | --- |
| 60 and above | Critical |
| 40 to 59 | High |
| 20 to 39 | Medium |
| below 20 | Low |

Worked example: `severity = 5`, `people_affected = 30`, `vulnerable = true`
gives `50 + 30 + 20 = 100`, which is Critical. The formula is deliberately
simple so that during the demonstration an operator can explain out loud why
one incident outranks another.

### Errors

Failures return the matching status and a body of `{"error": "<message>"}`. The
frontend shows the message as it arrives and does not rewrite it.

| Condition | Status | Message |
| --- | --- | --- |
| `location` empty or longer than 120 characters | 400 | `Location is required (max 120 characters)` |
| `people_affected` outside 1–500 | 400 | `People affected must be between 1 and 500` |
| `severity` outside 1–5 | 400 | `Severity must be between 1 and 5` |
| Shelter does not have enough remaining capacity | 409 | `Shelter has only N places remaining` |
| Responder is not available | 409 | `Responder is not available` |
| Status skips a step, for example resolving a registered incident | 409 | `Incident must be assessed first` |
