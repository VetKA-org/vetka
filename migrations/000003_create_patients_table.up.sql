DROP TYPE IF EXISTS gender;
CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS animal_species (
    id   uuid DEFAULT gen_random_uuid () primary key,
    name varchar(32) not null unique
);

INSERT INTO animal_species (name) VALUES (
    unnest(array['amphibian', 'bird', 'cat', 'dog', 'exotic', 'reptile', 'rodent'])
) ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS patients (
    id            uuid DEFAULT gen_random_uuid () primary key,
    name          varchar(32) not null,
    species       uuid REFERENCES animal_species(id),
    gender        gender not null,
    breed         varchar(64),
    birth         date not null,
    aggressive    boolean not null,
    vaccinated_at date,
    sterilized_at date
);
