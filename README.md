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

The prototype is a demonstration, not a system to run an emergency on:

- It runs on one machine for one operator at a time, with no sign-in.
- Its data is the fictional sample set; nothing is imported from, or sent to,
  any council system, and no message reaches a resident or a responder.
- Times are stored in UTC and shown in the browser's own time zone.
- A shelter's capacity is a single number of places; it does not model rooms,
  accessibility or special needs.
- Priority uses the agreed formula only. Real triage would add the judgement of
  the person assessing the incident.

## Layout

| Path | Contents |
| --- | --- |
| `backend/` | Go API server, standard library HTTP, SQLite via `modernc.org/sqlite` |
| `frontend/` | React and TypeScript client, built with Vite |
| `seed/` | Sample data loaded into an empty database on first start |
| `bin/` | Ready-built programs for Windows, macOS and Linux |
| `scripts/` | Release build script for `bin/` |

## Screens

| Capability area | Screens in the design |
| --- | --- |
| Incident registration | New incident form |
| Assessment and prioritisation | Incident list, incident detail |
| Shelter capacity | Shelter list, shelter detail |
| Responder dispatch | Responder list, roster |
| Resource allocation | Equipment inventory, allocation request |
| Progress tracking | Incident timeline |
| Operations dashboard | Dashboard |
| Reporting | Report builder |
| Audit trail | Activity log |

The shell is a header plus left navigation holding Dashboard, Incidents,
Shelters, Responders and Activity log. Every list needs three states, not just
the populated one: loading, empty, and failed to load. Priority bands are shown
as colour **and** text, never colour alone.

## Demonstration walkthrough

About five minutes, starting from fresh sample data (`docker compose down -v`,
or delete `bin/data.db`). Each step names what to point out.

1. **Dashboard.** Six open incidents, four responders free, Northbay Primary
   School Gym nearly full at 38 of 40 places. Keep these numbers in mind.
2. **Incidents.** The list is already in priority order. Explain one score out
   loud: Lowtide Crescent is severity 4 (40) + 18 people (18) + vulnerable
   people (20) = 78, Critical.
3. **Register incident.** Submit with the location empty to show the rule
   coming back from the server. Then enter a storm at *Bayside Caravan Park*,
   10 people, severity 2, vulnerable people ticked. It scores 20 + 10 + 20 = 50,
   High, and appears highlighted between the 74 and the 45.
4. **Open it and assess it.** Before assessing, only *Assess* is offered.
5. **Assign shelter.** Pick Northbay Primary School Gym first: it is refused
   with "Shelter has only 2 places remaining", nothing changes, and the page
   lists the shelters that can take all 10 people, roomiest first. Choose
   Harbour View Community Hall from that list with *Place here*.
6. **Dispatch responder.** Pick someone marked as on another job to show
   "Responder is not available", then pick an available responder.
7. **Resolve.** The history at the bottom of the page now shows all five steps
   with who and when.
8. **Dashboard again.** The places and the responder have come back; resolved
   is up by one. Finish on the **Activity log** to show nothing can be edited.


What the full design contains against what the prototype does. Each feature is
in exactly one column.

| Feature | Implemented | Partly | Designed only | Out of scope |
| --- | :---: | :---: | :---: | :---: |
| Register an incident, with the client's validation rules | ✓ | | | |
| Priority score and band worked out by the system | ✓ | | | |
| Incident list in priority order | ✓ | | | |
| Incident detail with only the steps its status allows | ✓ | | | |
| Assessment step | | ✓ | | |
| Place people in a shelter, refusing one without room | ✓ | | | |
| Offer the shelters with room when one is full (change CR-01) | ✓ | | | |
| Shelter list with remaining places | ✓ | | | |
| Shelter detail: adjust capacity, close a shelter | | | ✓ | |
| Dispatch an available responder | ✓ | | | |
| Responder list with availability | ✓ | | | |
| Roster and shifts | | | ✓ | |
| Resolve, releasing shelter places and the responder | ✓ | | | |
| Equipment inventory and allocation requests | | | ✓ | |
| Incident timeline | ✓ | | | |
| Operations dashboard | ✓ | | | |
| Report builder and export | | | ✓ | |
| Activity log of every change, read only | ✓ | | | |
| Sign-in, roles and permissions | | | | ✓ |
| Messages to residents and responders | | | | ✓ |

**Partly:** assessment moves the incident on and is logged, but records no
assessment notes and does not re-score the incident.

**Out of scope:** without sign-in, the activity log names the operator from the
`X-Actor` header or records `duty.officer`.

Every implemented row on the end-to-end path, and each validation rule, is
exercised by `backend/handlers/acceptance_test.go`:

```sh
cd backend
go test -v -run Acceptance ./handlers/
```

## Running with Docker

Needs only Docker with the Compose plugin; Go and Node are not required.

```sh
docker compose up --build
```

Open http://localhost:8080. The image builds the client, embeds it in the Go
server and runs both from one container. On first start the database is created
in a named volume and filled with the sample data.

| Task | Command |
| --- | --- |
| Use another port | `CSIT214_PORT=9090 docker compose up --build` |
| Stop, keeping the data | `docker compose down` |
| Start again from the sample data | `docker compose down -v` then `docker compose up` |

## Running without Docker

`bin/` holds ready-built single-file programs with the interface and the sample
data handling built in. Nothing needs to be installed and no network access is
used.

| System | File |
| --- | --- |
| Windows | `bin/csit214-windows.exe` |
| macOS on Apple silicon (M1 and later) | `bin/csit214-macos-apple-silicon` |
| macOS on Intel | `bin/csit214-macos-intel` |
| Linux | `bin/csit214-linux` |

Double-click the file for your system, or run it from a terminal. A browser
opens at http://localhost:8080; if it does not, open that address yourself. The
window that appears is the server: close it, or press Ctrl+C, to stop.

The database is created next to the program as `bin/data.db` and filled from
`seed/seed.sql`, so keep the repository's folder layout. Delete `bin/data.db` to
start again from the sample data.

The programs are not signed, so the operating system may ask first:

- **Windows** shows "Windows protected your PC": choose *More info*, then *Run anyway*.
- **macOS** may refuse a file downloaded as a ZIP: right-click it, choose *Open*,
  then *Open* again. Alternatively run
  `xattr -d com.apple.quarantine bin/csit214-macos-*` once. A `git clone` is not affected.
- **Linux** may need `chmod +x bin/csit214-linux` if the execute bit was lost.

If port 8080 is taken, start it from a terminal with another port, for example
`PORT=9090 bin/csit214-linux` (on Windows: `set PORT=9090` then the `.exe`).

To rebuild the four files after a change (needs Go and Node):

```sh
./scripts/build-binaries.sh
```

## Running in development

Go 1.23 or newer and Node 22 or newer. There is no database server to install:
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
Opening the backend's own port in a browser shows a short notice instead of the
interface until the client has been built into it; the Docker image always has.

## Configuration

Every setting has a working default; the prototype starts with no environment
variables set at all.

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | Port the server listens on |
| `CSIT214_DB_PATH` | `data.db` in the working directory; next to the program for the `bin/` builds | SQLite file location |
| `CSIT214_SEED_PATH` | `seed/seed.sql`, then `../seed/seed.sql`; `../seed/seed.sql` beside the program for the `bin/` builds | Sample data used when the database is empty |
| `CSIT214_OPEN_BROWSER` | on for the `bin/` builds | Set to `0` to stop a `bin/` build opening a browser |
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
| `GET` | `/api/audit` | Audit entries, newest first. Optional `entity_type`, `entity_id` and `limit` (1–1000, default 200) |
| `GET` | `/api/dashboard` | Counts per status, counts per priority band, shelter occupancy, number of available responders |

Every write operation appends a row to `audit` in the same transaction as the
change itself. The client asked for this explicitly, so it is not optional.
There is no sign-in in the prototype: a request may name the operator in an
`X-Actor` header, otherwise the entry is recorded against `duty.officer`.

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
| Shelter does not have enough remaining capacity | 409 | `Shelter has only N places remaining`, plus `alternatives` (below) |
| Responder is not available | 409 | `Responder is not available` |
| Status skips a step, for example resolving a registered incident | 409 | `Incident must be assessed first` |

Other replies follow the same shape: `400` for an unknown incident `type`, a
malformed body or id; `404` for an incident, shelter or responder that does not
exist; `409` for a step that has already happened, such as assessing twice or
resolving a resolved incident.

When a shelter is refused for lack of room, the 409 also lists the shelters
that could take the whole group, roomiest first, so the operator can choose
again without guessing. The list is empty when no shelter is large enough.

```json
{
  "error": "Shelter has only 2 places remaining",
  "alternatives": [
    { "id": 4, "name": "Riverside Civic Centre", "capacity_remaining": 200 },
    { "id": 1, "name": "Harbour View Community Hall", "capacity_remaining": 96 },
    { "id": 3, "name": "Pelican Point Surf Club", "capacity_remaining": 28 }
  ]
}
```

A shelter is optional: an incident may be dispatched or resolved straight after
assessment, as a medical call often needs no shelter. Resolving releases the
incident's shelter places and returns its responder to `available`.
