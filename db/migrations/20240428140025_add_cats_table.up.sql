CREATE TYPE race AS ENUM (
    'Persian',
    'Maine Coon',
    'Siamese',
    'Ragdoll',
    'Bengal',
    'Sphynx',
    'British Shorthair',
    'Abyssinian',
    'Scottish Fold',
    'Birman'
);

CREATE TYPE sex AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS "cats" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL,
    name VARCHAR(30) NOT NULL,
    race race NOT NULL,
    sex sex NOT NULL,
    age_in_month INT NOT NULL,
    image_urls TEXT[] NOT NULL,
    description VARCHAR(200) NOT NULL,
    hasMatched BOOLEAN NOT NULL default FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

CREATE INDEX IF NOT EXISTS id_cat_email ON cats(user_id);

CREATE INDEX IF NOT EXISTS idx_cat_name ON cats USING GIN(name);

CREATE INDEX IF NOT EXISTS idx_cat_created_at ON cats(created_at);

CREATE INDEX IF NOT EXISTS idx_cat_sex ON cats USING HASH(sex);

CREATE INDEX IF NOT EXISTS idx_cat_age ON cats(age_in_month);