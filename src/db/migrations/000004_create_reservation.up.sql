CREATE TABLE reservation.reservations (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id         uuid        NOT NULL,  -- facility.rooms への論理参照（FKなし）
    user_id         uuid        NOT NULL,  -- identity.users への論理参照（FKなし）
    title           text        NOT NULL,
    attendees       int         NOT NULL DEFAULT 1 CHECK (attendees > 0),
    period          tstzrange   NOT NULL,
    status          text        NOT NULL
                                CHECK (status IN ('HELD', 'CONFIRMED', 'CANCELLED', 'EXPIRED')),
    hold_expires_at timestamptz,
    cancelled_at    timestamptz,
    idempotency_key text        UNIQUE,     -- 同じリクエストの二重送信対策
    version         int         NOT NULL DEFAULT 1,  -- 楽観ロック（更新時に version を比較）
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),

    -- 期間は [開始, 終了) の半開区間で、空でないこと
    CONSTRAINT period_valid CHECK (
        NOT isempty(period)
        AND lower_inc(period) AND NOT upper_inc(period)
        AND NOT lower_inf(period) AND NOT upper_inf(period)
    ),
    -- 仮押さえ中は期限が必須
    CONSTRAINT hold_has_expiry CHECK (
        status <> 'HELD' OR hold_expires_at IS NOT NULL
    ),
    -- キャンセル時は日時を記録
    CONSTRAINT cancelled_has_time CHECK (
        status <> 'CANCELLED' OR cancelled_at IS NOT NULL
    ),
    -- 二重予約の防止: 同じ部屋で有効な予約の期間が重なることを禁止
    CONSTRAINT no_double_booking EXCLUDE USING gist (
        room_id WITH =,
        period  WITH &&
    ) WHERE (status IN ('HELD', 'CONFIRMED'))
);

-- 自分の予約一覧（新しい順）
CREATE INDEX reservations_user_idx
    ON reservation.reservations (user_id, lower(period) DESC);

-- 期限切れジョブが HELD を拾うため
CREATE INDEX reservations_hold_expiry_idx
    ON reservation.reservations (hold_expires_at)
    WHERE status = 'HELD';

-- 部屋ごとの空き状況検索は no_double_booking の GiST インデックスが使われる

-- Transactional Outbox（フェーズ3で本格利用。モノリス段階から書いておくと移行が楽）
CREATE TABLE reservation.outbox (
    id             bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    aggregate_type text        NOT NULL,   -- 'reservation'
    aggregate_id   uuid        NOT NULL,
    event_type     text        NOT NULL,   -- 'ReservationConfirmed' など
    payload        jsonb       NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT now(),
    published_at   timestamptz            -- NULL = 未送信
);

CREATE INDEX outbox_unpublished_idx
    ON reservation.outbox (id)
    WHERE published_at IS NULL;
