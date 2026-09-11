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
| `backend/` | Go API server |
| `frontend/` | React and TypeScript client |
| `seed/` | Sample data for demonstrations |

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
