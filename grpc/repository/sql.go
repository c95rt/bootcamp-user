package repository

const (
	createUserTableQuery = `
	CREATE TABLE IF NOT EXISTS user (
		id int(11) NOT NULL AUTO_INCREMENT,
		email varchar(255) NOT NULL,
		firstname varchar(255) NOT NULL,
		lastname varchar(255) NOT NULL,
		active tinyint(1) NOT NULL DEFAULT '1',
		password varchar(255) NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4
`
)
