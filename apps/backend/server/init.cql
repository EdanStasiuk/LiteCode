-- Create keyspace
CREATE KEYSPACE IF NOT EXISTS litecode
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

-- Switch to keyspace
USE litecode;

-- submission_code
CREATE TABLE IF NOT EXISTS submission_code (
    submission_id TEXT PRIMARY KEY,
    code TEXT,
    language TEXT
);


-- submissions_by_user table
CREATE TABLE submissions_by_user (
    user_id TEXT,
    submission_id TEXT,
    problem_id TEXT,
    status TEXT,
    runtime DOUBLE,
    memory BIGINT,
    result TEXT,
    language TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id, submission_id)
);

-- submissions_by_problem table
CREATE TABLE submissions_by_problem (
    problem_id TEXT,
    submission_id TEXT,
    user_id TEXT,
    status TEXT,
    runtime DOUBLE,
    memory BIGINT,
    result TEXT,
    language TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (problem_id, submission_id)
);

-- submissions_by_problem_and_user table
CREATE TABLE IF NOT EXISTS submissions_by_problem_and_user (
    problem_id TEXT,
    user_id TEXT,
    submission_id TEXT,
    status TEXT,
    runtime DOUBLE,
    memory BIGINT,
    result TEXT,
    language TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY ((problem_id, user_id), submission_id)
);
