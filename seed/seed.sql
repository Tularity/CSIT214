-- Sample data for the CoastLink Council prototype.
-- Loaded automatically the first time the server starts against an empty database.
-- Place names are fictional. Times are UTC, ISO 8601.

INSERT INTO shelters (id, name, address, capacity_total, capacity_used) VALUES
    (1, 'Harbour View Community Hall', '12 Marine Parade, Coastvale',   120, 24),
    (2, 'Northbay Primary School Gym', '8 Kingfisher Road, Northbay',    40, 38),
    (3, 'Pelican Point Surf Club',     '3 Esplanade, Pelican Point',     60, 32),
    (4, 'Riverside Civic Centre',      '45 Willow Street, Riverside',   200,  0);

INSERT INTO responders (id, name, skill, status) VALUES
    (1, 'Aaron Whitfield', 'Paramedic',              'assigned'),
    (2, 'Bianca Torres',   'Swift water rescue',     'available'),
    (3, 'Callum Reid',     'Logistics',              'available'),
    (4, 'Dana Nakamura',   'Evacuation coordinator', 'available'),
    (5, 'Elias Brandt',    'Structural assessment',  'off_duty'),
    (6, 'Farah Haddad',    'Field medic',            'available');

-- priority_score = severity * 10 + min(people_affected, 30) + (vulnerable ? 20 : 0)
INSERT INTO incidents
    (id, type, location, description, people_affected, vulnerable, severity,
     priority_score, priority_band, status, shelter_id, responder_id, reported_at)
VALUES
    (1, 'flood', 'Lowtide Crescent, Riverside',
     'Stormwater drain overflow has cut off vehicle access to six houses.',
     18, 1, 4, 78, 'Critical', 'registered', NULL, NULL, '2026-09-18T21:10:00Z'),

    (2, 'fire', 'Unit 7, 210 Beacon Street, Coastvale',
     'Kitchen fire in an apartment block; the top two floors have been evacuated.',
     24, 0, 5, 74, 'Critical', 'in_progress', 1, 1, '2026-09-19T03:42:00Z'),

    (3, 'storm', 'Pelican Point Esplanade',
     'Fallen power line blocking the foreshore road in both directions.',
     5, 0, 4, 45, 'High', 'assessed', NULL, NULL, '2026-09-19T08:15:00Z'),

    (4, 'medical', 'Northbay Aged Care, 8 Kingfisher Road',
     'Extended power outage affecting oxygen concentrators and refrigerated medication.',
     32, 1, 4, 90, 'Critical', 'assigned', 3, NULL, '2026-09-19T11:05:00Z'),

    (5, 'flood', 'Marine Parade, Coastvale',
     'Sandbagging requested for shopfronts near the boat ramp before the evening tide.',
     12, 0, 2, 32, 'Medium', 'registered', NULL, NULL, '2026-09-19T14:30:00Z'),

    (6, 'other', 'Willow Street Bridge, Riverside',
     'Debris build-up under the bridge is restricting water flow.',
     3, 0, 1, 13, 'Low', 'registered', NULL, NULL, '2026-09-20T01:20:00Z'),

    (7, 'storm', 'Gull Rise Estate, Northbay',
     'Roof damage across eleven properties after overnight winds.',
     28, 1, 3, 78, 'Critical', 'resolved', 4, 5, '2026-09-17T19:55:00Z'),

    (8, 'medical', 'Coastvale Foreshore Markets',
     'Heat exhaustion cases reported at a community event.',
     7, 0, 3, 37, 'Medium', 'resolved', NULL, 6, '2026-09-16T05:40:00Z');

INSERT INTO audit (at, actor, action, entity_type, entity_id, detail) VALUES
    ('2026-09-16T05:40:00Z', 'duty.officer',  'register',         'incident', 8, 'Registered with priority score 37 (Medium)'),
    ('2026-09-16T05:52:00Z', 'duty.officer',  'assess',           'incident', 8, 'Assessed on site by the event first aid post'),
    ('2026-09-16T06:05:00Z', 'm.alvarez',     'assign-responder', 'incident', 8, 'Assigned responder 6 (Farah Haddad)'),
    ('2026-09-16T08:30:00Z', 'm.alvarez',     'resolve',          'incident', 8, 'Resolved; responder 6 returned to available'),
    ('2026-09-17T19:55:00Z', 'duty.officer',  'register',         'incident', 7, 'Registered with priority score 78 (Critical)'),
    ('2026-09-17T20:14:00Z', 'duty.officer',  'assess',           'incident', 7, 'Assessed; vulnerable residents confirmed on site'),
    ('2026-09-17T20:40:00Z', 'm.alvarez',     'assign-shelter',   'incident', 7, 'Assigned shelter 4 (Riverside Civic Centre) for 28 people'),
    ('2026-09-17T21:02:00Z', 'm.alvarez',     'assign-responder', 'incident', 7, 'Assigned responder 5 (Elias Brandt)'),
    ('2026-09-18T11:18:00Z', 'm.alvarez',     'resolve',          'incident', 7, 'Resolved; 28 places released at shelter 4'),
    ('2026-09-18T21:10:00Z', 'duty.officer',  'register',         'incident', 1, 'Registered with priority score 78 (Critical)'),
    ('2026-09-19T03:42:00Z', 'duty.officer',  'register',         'incident', 2, 'Registered with priority score 74 (Critical)'),
    ('2026-09-19T03:58:00Z', 'duty.officer',  'assess',           'incident', 2, 'Assessed; fire service already on scene'),
    ('2026-09-19T04:15:00Z', 's.okafor',      'assign-shelter',   'incident', 2, 'Assigned shelter 1 (Harbour View Community Hall) for 24 people'),
    ('2026-09-19T04:31:00Z', 's.okafor',      'assign-responder', 'incident', 2, 'Assigned responder 1 (Aaron Whitfield)'),
    ('2026-09-19T08:15:00Z', 'duty.officer',  'register',         'incident', 3, 'Registered with priority score 45 (High)'),
    ('2026-09-19T08:44:00Z', 'duty.officer',  'assess',           'incident', 3, 'Assessed; awaiting network operator isolation'),
    ('2026-09-19T11:05:00Z', 'duty.officer',  'register',         'incident', 4, 'Registered with priority score 90 (Critical)'),
    ('2026-09-19T11:20:00Z', 'duty.officer',  'assess',           'incident', 4, 'Assessed; residential aged care facility'),
    ('2026-09-19T11:48:00Z', 's.okafor',      'assign-shelter',   'incident', 4, 'Assigned shelter 3 (Pelican Point Surf Club) for 32 people'),
    ('2026-09-19T14:30:00Z', 'duty.officer',  'register',         'incident', 5, 'Registered with priority score 32 (Medium)'),
    ('2026-09-20T01:20:00Z', 'duty.officer',  'register',         'incident', 6, 'Registered with priority score 13 (Low)');
