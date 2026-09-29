CREATE TABLE facility.locations (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       text        NOT NULL UNIQUE,
    address    text,
    timezone   text        NOT NULL DEFAULT 'Asia/Tokyo',  -- 表示・営業時間判定用
    open_time  time        NOT NULL DEFAULT '08:00',
    close_time time        NOT NULL DEFAULT '22:00',
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK (open_time < close_time)
);

CREATE TABLE facility.rooms (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id  uuid        NOT NULL REFERENCES facility.locations (id),
    name         text        NOT NULL,
    capacity     int         NOT NULL CHECK (capacity > 0),
    hourly_price int         NOT NULL DEFAULT 0 CHECK (hourly_price >= 0),  -- 円。決済フェーズで使用
    is_active    boolean     NOT NULL DEFAULT true,  -- 削除せず無効化する
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    UNIQUE (location_id, name)
);

CREATE TABLE facility.equipment (
    id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL UNIQUE  -- 例: プロジェクター、ホワイトボード
);

CREATE TABLE facility.room_equipment (
    room_id      uuid NOT NULL REFERENCES facility.rooms (id) ON DELETE CASCADE,
    equipment_id uuid NOT NULL REFERENCES facility.equipment (id),
    PRIMARY KEY (room_id, equipment_id)
);
