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
    executor_type STRING NOT NULL,
    executor_id INT NOT NULL,
    stat STRING DEFAULT 'В процессе',
    date_start STRING NOT NULL,
    date_end STRING NOT NULL,
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
    father_id INT DEFAULT 0,
    title STRING NOT NULL,
);