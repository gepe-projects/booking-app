-- Drop index dulu sebelum drop tabel
DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_user_identities_email;
DROP INDEX IF EXISTS idx_user_identities_phone;

-- Drop tabel child dulu (karena ada FK ke users)
DROP TABLE IF EXISTS user_identities;

-- Baru drop tabel parent
DROP TABLE IF EXISTS users;
