CREATE TABLE users
(
    id            serial       primary key,
    name          varchar(255) not null default '',
    username      varchar(255) not null unique,
    password      varchar(255) not null,
    email         varchar(255) not null unique,
    created_at    timestamp    not null,
    activated_at  timestamp,
    CONSTRAINT proper_users_name      CHECK (name = '' OR name ~* '^[a-zA-Z+ ]*[a-zA-Z]*$'),
    CONSTRAINT proper_users_username  CHECK (username ~* '^[A-Za-z]\w{2,}$'),
    CONSTRAINT proper_users_email     CHECK (email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    CONSTRAINT proper_users_time      CHECK (activated_at IS NULL OR activated_at < created_at)
);

CREATE TABLE tests
(
    id                      serial        primary key,
    title                   varchar(255)  not null unique,
    description             varchar(2048) not null default '',
    random_questions_order  boolean       not null default false,
    questions_visibility    varchar(12)   not null,
    start_time              timestamp     not null,
    end_time                timestamp     not null,
    duration_sec            integer       not null,
    CONSTRAINT proper_tests_title      CHECK (length(title) > 0),
    CONSTRAINT proper_tests_type       CHECK (questions_visibility in ('ShowOneByOne', 'ShowAll')),
    CONSTRAINT proper_tests_time       CHECK (start_time < end_time),
    CONSTRAINT proper_tests_duration   CHECK (duration_sec < EXTRACT(EPOCH FROM end_time - start_time))
);

CREATE TABLE questions
(
    id             serial                                               primary key,
    test_id        int          references tests (id) on delete cascade not null,
    text           varchar(2048)                                        not null,
    answer_type    varchar(10)                                          not null,
    show_answers   varchar[]                                            not null,
    true_answers   varchar[]                                            not null,
    CONSTRAINT proper_questions_text CHECK (length(text) > 0),
    CONSTRAINT proper_questions_type CHECK (answer_type in ('freeField', 'oneSelect', 'manySelect')),

    CHECK (array_length(true_answers, 1) <= array_length(show_answers, 1)),
    CHECK (array_length(show_answers, 1) < 100),

    CHECK (
        answer_type = 'freeField' AND array_length(show_answers, 1) = 0 OR
        answer_type in ('oneSelect', 'manySelect') AND array_length(show_answers, 1) > 1),
    CHECK (
        answer_type in ('freeField', 'oneSelect') AND array_length(true_answers, 1) = 1 OR
        answer_type = 'manySelect' AND array_length(true_answers, 1) > 0)
);

CREATE TABLE question_answers
(
    user_id        int          references users (id)     on delete cascade not null,
    question_id    int          references questions (id) on delete cascade not null,
    answer         varchar[]    not null,
    time           timestamp    not null,
    PRIMARY KEY (user_id, question_id),
    CONSTRAINT proper_answer CHECK (array_length(answer, 1) > 0)
);

CREATE TABLE test_answers
(
    user_id        int        references users (id) on delete cascade not null,
    test_id        int        references tests (id) on delete cascade not null,
    complete_time  timestamp  not null,
    PRIMARY KEY (user_id, test_id)
);

CREATE TABLE refresh_sessions (
    --id            serial       primary key,
    refresh_token uuid         primary key,
    user_id       int          references users (id) on delete cascade not null,
    user_agent    varchar(200) not null,
    fingerprint   varchar(200) not null,
    ip            inet         not null,
    expires_at    timestamp    not null
    --createdAt timestamp with time zone not null DEFAULT now(),
    --CREATE INDEX ON ip_ranges USING GIST (ip_range inet_ops)
);