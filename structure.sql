CREATE DATABASE bluemine;
USE bluemine;

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    username STRING NOT NULL,
    user_fio STRING NOT NULL,
    rating INT DEFAULT 0,
    department INT,
    group_id INT,
    isAdmin BOOLEAN DEFAULT false
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task_name STRING NOT NULL,
    rating INT DEFAULT 0,
    path_to_desc STRING,
    stat SMALLINT DEFAULT 0,
    date_start DATE NOT NULL,
    date_end DATE
);

CREATE TABLE checkboxs (
    id INT,
    task_id INT,
    checked BOOLEAN DEFAULT false,
    desk STRING NOT NULL,
    UNIQUE (task_id)
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    group_name STRING NOT NULL
);

CREATE TABLE wiki (
    id SERIAL PRIMARY KEY,
    father_id INT,
    title STRING NOT NULL,
    path_to_article STRING NOT NULL
);