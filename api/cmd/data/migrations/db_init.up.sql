CREATE DATABASE friends_db;

-- Create Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create Relationships Table
CREATE TABLE relationships (
    id UUID PRIMARY KEY,
    requestor_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    target_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    relationship_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT no_self_friendship CHECK (requestor_id != target_id)
);


