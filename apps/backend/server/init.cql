-- Create keyspace
CREATE KEYSPACE IF NOT EXISTS litecode
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

-- Switch to keyspace
USE litecode;

-- SubmissionCode table
CREATE TABLE IF NOT EXISTS submission_code (
    submission_id TEXT PRIMARY KEY,
    code TEXT
);

-- UserSubmission table
CREATE TABLE IF NOT EXISTS user_submission (
    user_id TEXT,
    submission_id TEXT,
    problem_id TEXT,
    status TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id, submission_id)
);

-- ProblemSubmission table
CREATE TABLE IF NOT EXISTS problem_submission (
    problem_id TEXT,
    submission_id TEXT,
    user_id TEXT,
    status TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (problem_id, submission_id)
);
