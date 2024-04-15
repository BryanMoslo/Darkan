CREATE TABLE IF NOT EXISTS "keywords" (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    value varchar(255) NOT NULL,
    callback_url text NOT NULL,
    found bool NOT NULL DEFAULT false 
);