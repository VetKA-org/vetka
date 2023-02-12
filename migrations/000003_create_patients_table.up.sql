DROP TYPE IF EXISTS gender;
CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS animal_species (
    species_id uuid DEFAULT gen_random_uuid () primary key,
    title      varchar(32) not null unique
);

INSERT INTO animal_species (title) VALUES (
    unnest(array['Amphibian', 'Bird', 'Cat', 'Dog', 'Exotic', 'Reptile', 'Rodent'])
) ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS patients (
    patient_id    uuid DEFAULT gen_random_uuid () primary key,
    name          varchar(32) not null,
    species_id    uuid REFERENCES animal_species(species_id),
    gender        gender not null,
    breed         varchar(64),
    birth         date not null,
    aggressive    boolean not null DEFAULT false,
    vaccinated_at date,
    sterilized_at date
);
