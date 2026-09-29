-- パスワードは保持せず、OIDC(Google等)のプロバイダー情報でユーザーを識別する。
-- provider + provider_user_id の組がプロバイダー側のアカウントと1対1に対応する。
CREATE TABLE identity.users (
    id               uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    email            citext      NOT NULL UNIQUE,
    name             text        NOT NULL,
    provider         text        NOT NULL,  -- 'google' など
    provider_user_id text        NOT NULL,  -- IDトークンの sub
    role             text        NOT NULL DEFAULT 'member'
                                 CHECK (role IN ('member', 'admin')),
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_user_id)
);
