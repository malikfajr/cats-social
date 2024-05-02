CREATE TYPE ISSUER AS (
    name VARCHAR(50),
    email VARCHAR(50),
    created_at TIMESTAMP
);

CREATE TYPE CAT_DETAIL AS (
    id VARCHAR(20) ,
    name VARCHAR(30),
    race race,
    sex sex,
    description VARCHAR(200),
    age_in_month SMALLINT,
    image_urls TEXT[],
    hasMatched BOOLEAN,
    created_at TIMESTAMP
);

-- CREATE TABLE IF NOT EXISTS matches (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     issued_by ISSUER NOT NULL,
--     match_cat_detail CAT_DETAIL NOT NULL,
--     user_cat_detail CAT_DETAIL NOT NULL,
--     message VARCHAR(150) NOT NULL,
--     created_at TIMESTAMP DEFAULT NOW()
-- );

CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issued_by JSONB NOT NULL,
    match_cat_detail JSONB NOT NULL,
    match_user_email VARCHAR(50) NOT NULL,
    user_cat_detail JSONB NOT NULL,
    message VARCHAR(150) NOT NULL,
    status ENUM("pending", "approved", "rejected") DEFAULT "pending",
    created_at TIMESTAMP DEFAULT NOW()
);