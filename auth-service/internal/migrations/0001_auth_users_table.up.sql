CREATE TABLE IF NOT EXISTS users(
		id SERIAL primary key,
		login varchar(64) not null,
		email varchar(255) not null unique,
		password varchar(255) not null unique,
		img_path varchar(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)