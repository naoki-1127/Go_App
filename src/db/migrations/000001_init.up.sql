-- =====================================================================
-- 会議室予約システム（モジュラーモノリス版）
--   モジュールごとに PostgreSQL のスキーマを分ける:
--     identity    : ユーザー
--     facility    : 拠点・会議室・設備
--     reservation : 予約・outbox
--   ルール: 外部キーは同一スキーマ内だけ。スキーマをまたぐ参照は
--           uuid を持つだけにする（将来のサービス分割に備える）。
-- =====================================================================

CREATE EXTENSION IF NOT EXISTS btree_gist;  -- EXCLUDE 制約で uuid の = を使うため
CREATE EXTENSION IF NOT EXISTS citext;      -- 大文字小文字を区別しないメールアドレス

CREATE SCHEMA IF NOT EXISTS identity;
CREATE SCHEMA IF NOT EXISTS facility;
CREATE SCHEMA IF NOT EXISTS reservation;
