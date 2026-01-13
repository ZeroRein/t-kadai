-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    language VARCHAR(5) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Events Table
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID, -- Foreign key ommitted for simplicity or add REFERENCES users(id)
    title VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    description TEXT,
    color VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Memos Table
CREATE TABLE IF NOT EXISTS memos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    title VARCHAR(255),
    content TEXT,
    linked_date TIMESTAMP, -- Nullable for free memos
    theme_color VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed Data
INSERT INTO users (username, password_hash, language) VALUES 
('demo_user', 'hash_placeholder', 'ja');

INSERT INTO events (title, start_time, end_time, description, color) VALUES
('Project Kickoff', CURRENT_DATE + TIME '10:00:00', CURRENT_DATE + TIME '11:00:00', 'Discuss project requirements', '#4a90e2'),
('Lunch with Team', CURRENT_DATE + TIME '12:00:00', CURRENT_DATE + TIME '13:00:00', 'Italian place', '#2ecc71');

INSERT INTO memos (title, content, linked_date, theme_color) VALUES
('Meeting Notes', 'Kickoff went well. Need to set up Docker.', CURRENT_DATE, '#f1c40f'),
('Idea: Bookshelf', 'Use a grid layout for memos.', NULL, '#9b59b6');
