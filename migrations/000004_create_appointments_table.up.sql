DROP TYPE IF EXISTS appointment_status;
CREATE TYPE appointment_status AS ENUM ('scheduled', 'opened', 'closed', 'canceled');

CREATE TABLE IF NOT EXISTS appointments (
    appointment_id uuid DEFAULT gen_random_uuid (),
    patient_id     uuid REFERENCES patients(patient_id),
    assignee_id    uuid REFERENCES users(user_id),
    scheduled_for  timestamp not null,
    primary key    (patient_id, scheduled_for),
    status         appointment_status,
    reason         varchar(255) not null,
    details        text
);
