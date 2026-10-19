CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number TEXT NOT NULL UNIQUE,
    room_type TEXT NOT NULL,
    floor INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'available' CHECK (status IN ('available','occupied','needs_cleaning','out_of_order')),
    rate NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (rate >= 0),
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);

CREATE TABLE guests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL, last_name TEXT NOT NULL, email TEXT, phone TEXT,
    nationality TEXT, document_number TEXT, notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX guests_email_unique ON guests (lower(email)) WHERE email IS NOT NULL AND archived_at IS NULL;
CREATE INDEX guests_name_idx ON guests (lower(last_name), lower(first_name)) WHERE archived_at IS NULL;

CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL, last_name TEXT NOT NULL, email TEXT UNIQUE, phone TEXT,
    role TEXT NOT NULL, department TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','on_leave','inactive')),
    hired_on DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);

CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id), guest_id UUID NOT NULL REFERENCES guests(id),
    check_in DATE NOT NULL, check_out DATE NOT NULL, status TEXT NOT NULL DEFAULT 'confirmed' CHECK (status IN ('confirmed','checked_in','checked_out','cancelled','no_show')),
    adults INTEGER NOT NULL DEFAULT 1 CHECK (adults > 0), children INTEGER NOT NULL DEFAULT 0 CHECK (children >= 0),
    rate NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (rate >= 0), source TEXT NOT NULL DEFAULT 'direct', special_requests TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ,
    CHECK (check_out > check_in),
    EXCLUDE USING gist (room_id WITH =, daterange(check_in, check_out, '[)') WITH &&)
        WHERE (archived_at IS NULL AND status IN ('confirmed','checked_in'))
);
CREATE INDEX reservations_guest_idx ON reservations (guest_id, check_in DESC) WHERE archived_at IS NULL;
CREATE INDEX reservations_arrivals_idx ON reservations (check_in, status) WHERE archived_at IS NULL;

CREATE TABLE housekeeping_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), room_id UUID NOT NULL REFERENCES rooms(id),
    assigned_employee_id UUID REFERENCES employees(id), task_type TEXT NOT NULL,
    priority TEXT NOT NULL DEFAULT 'normal' CHECK (priority IN ('low','normal','high','urgent')),
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','in_progress','done','cancelled')),
    due_at TIMESTAMPTZ, notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);

CREATE TABLE folios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), guest_id UUID NOT NULL REFERENCES guests(id), reservation_id UUID REFERENCES reservations(id),
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','closed','void')),
    currency CHAR(3) NOT NULL DEFAULT 'SAR', total_amount NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), folio_id UUID NOT NULL REFERENCES folios(id),
    amount NUMERIC(12,2) NOT NULL CHECK (amount > 0), currency CHAR(3) NOT NULL DEFAULT 'SAR',
    method TEXT NOT NULL CHECK (method IN ('cash','card','bank_transfer','online','other')),
    reference TEXT, status TEXT NOT NULL DEFAULT 'completed' CHECK (status IN ('pending','completed','refunded','failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);
CREATE INDEX payments_folio_idx ON payments (folio_id, created_at DESC) WHERE archived_at IS NULL;

CREATE TABLE maintenance_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), room_id UUID REFERENCES rooms(id), assigned_employee_id UUID REFERENCES employees(id),
    title TEXT NOT NULL, description TEXT NOT NULL DEFAULT '', priority TEXT NOT NULL DEFAULT 'normal' CHECK (priority IN ('low','normal','high','urgent')),
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','in_progress','done','cancelled')),
    reported_at TIMESTAMPTZ NOT NULL DEFAULT now(), resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);

CREATE TABLE inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), name TEXT NOT NULL, sku TEXT UNIQUE, category TEXT NOT NULL DEFAULT '', unit TEXT NOT NULL DEFAULT 'each',
    quantity NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (quantity >= 0), reorder_level NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (reorder_level >= 0),
    unit_cost NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (unit_cost >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);
CREATE TABLE inventory_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), item_id UUID NOT NULL REFERENCES inventory_items(id), employee_id UUID REFERENCES employees(id),
    quantity NUMERIC(12,3) NOT NULL CHECK (quantity > 0), movement_type TEXT NOT NULL CHECK (movement_type IN ('in','out','restock','use')),
    note TEXT NOT NULL DEFAULT '', created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), archived_at TIMESTAMPTZ
);
CREATE INDEX inventory_movement_item_idx ON inventory_movements (item_id, created_at DESC);

CREATE TABLE audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), record_kind TEXT NOT NULL, record_id UUID NOT NULL,
    action TEXT NOT NULL, details JSONB NOT NULL DEFAULT '{}'::jsonb, occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX audit_events_record_idx ON audit_events (record_kind, record_id, occurred_at DESC);
