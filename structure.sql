CREATE DATABASE bluemine;
USE bluemine;

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    username STRING NOT NULL,
    user_fio STRING NOT NULL,
    isAdmin BOOLEAN DEFAULT false,
    rating INT DEFAULT 0
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task_name STRING NOT NULL,
    task_creator INT NOT NULL,
    executor_id INT NOT NULL,
    executor_type STRING NOT NULL,
    stat STRING NOT NULL,
    priority STRING NOT NULL,
    date_added STRING NOT NULL,
    date_last_update STRING NOT NULL,
    date_start STRING NOT NULL,
    date_end STRING NOT NULL,
    rating INT DEFAULT 0
);

CREATE TABLE task_template (
    id SERIAL PRIMARY KEY,
    tmpl_name STRING NOT NULL,
    stat STRING NOT NULL,
    priority STRING NOT NULL,
    rating INT DEFAULT 0
);

CREATE TABLE checkboxs (
    id INT,
    task_id INT,
    checked BOOLEAN DEFAULT false,
    desk STRING NOT NULL,
    UNIQUE (task_id)
);

CREATE TABLE groups_profiles (
    group_id INT NOT NULL,
    profile_id INT NOT NULL
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    group_name STRING NOT NULL
);

CREATE TABLE wiki (
    id SERIAL PRIMARY KEY,
    author_id INT NOT NULL,
    father_id INT DEFAULT 0,
    title STRING NOT NULL
);
