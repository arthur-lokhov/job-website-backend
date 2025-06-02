-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Location (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    priority INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS Level (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    priority INT,
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updatedAt TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS Vacancy (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    department_url VARCHAR(255),
    level_id UUID NOT NULL,
    location_id UUID NOT NULL,
    info TEXT NOT NULL,
    application_form JSONB NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    important BOOLEAN NOT NULL DEFAULT false,
    priority INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (level_id) REFERENCES Level (id),
    FOREIGN KEY (location_id) REFERENCES Location (id)
);

CREATE TABLE IF NOT EXISTS Application (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vacancy_id UUID NOT NULL,
    answer JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (vacancy_id) REFERENCES Vacancy (id)
);

CREATE TABLE IF NOT EXISTS Status (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO Status (name) VALUES 
('PENDING'),
('VIEWED'),
('IN_PROGRESS'),
('APPROVED'),
('REJECTED');

CREATE TABLE IF NOT EXISTS Application_Status (
    application_id UUID NOT NULL,
    status_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

    PRIMARY KEY (application_id, status_id),
    FOREIGN KEY (application_id) REFERENCES Application (id),
    FOREIGN KEY (status_id) REFERENCES Status (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Application_Status;
DROP TABLE IF EXISTS Status;
DROP TABLE IF EXISTS Application;
DROP TABLE IF EXISTS Vacancy;
DROP TABLE IF EXISTS Level;
DROP TABLE IF EXISTS Location;
-- +goose StatementEnd

