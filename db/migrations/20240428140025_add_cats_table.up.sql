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
    id BIGSERIAL PRIMARY KEY NOT  NULL,
    user_email VARCHAR(50) NOT NULL,
    name VARCHAR(30) NOT NULL,
    race race NOT NULL,
    sex sex NOT NULL,
    age_in_month INT NOT NULL,
    image_urls TEXT[] NOT NULL,
    description VARCHAR(200) NOT NULL,
    hasMatched BOOLEAN NOT NULL default FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE "cats" ADD FOREIGN KEY ("user_email") REFERENCES "users" ("email");