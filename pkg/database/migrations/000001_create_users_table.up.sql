CREATE TABLE IF NOT EXISTS users (
  id              UUID PRIMARY KEY,
  name            VARCHAR(100) NOT NULL,
  image_url       VARCHAR(255),
  role            VARCHAR(50) NOT NULL DEFAULT 'user',
  created_at      TIMESTAMP NOT NULL DEFAULT now(),
  updated_at      TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_identities (
  id              UUID PRIMARY KEY,
  user_id         UUID NOT NULL,
  provider        VARCHAR(50) NOT NULL,   -- 'local', 'google', 'facebook', 'github', 'phone'
  provider_id     VARCHAR(255) NOT NULL,  -- unique ID dari provider (sub Google, ID FB, nomor HP, atau email utk local)
  email           VARCHAR(255),           -- optional, tergantung provider
  phone           VARCHAR(20),            -- optional, tergantung provider
  password_hash   VARCHAR(255),           -- hanya dipakai untuk provider = 'local'
  verified        BOOLEAN NOT NULL DEFAULT FALSE,
  created_at      TIMESTAMP NOT NULL DEFAULT now(),
  updated_at      TIMESTAMP NOT NULL DEFAULT now(),

  -- constraint: kombinasi provider + provider_id harus unik global
  CONSTRAINT uq_provider_providerid UNIQUE (provider, provider_id),

  -- constraint: 1 user tidak boleh punya 2 login method yang sama
  CONSTRAINT uq_user_provider UNIQUE (user_id, provider),

  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index tambahan untuk query
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_user_identities_email ON user_identities(email);
CREATE INDEX idx_user_identities_phone ON user_identities(phone);
