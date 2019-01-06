CREATE TABLE `cab_location` (
	`id` INT(255) NOT NULL AUTO_INCREMENT,
	`name` varchar(255) NOT NULL,
	`cabtype` varchar(255) NOT NULL,
	`lat` FLOAT(3),
	`lng` FLOAT(3),
	`ontrip` BOOLEAN DEFAULT '0',
	`last_updated` DATETIME(3),
	`time_start` FLOAT(3),
	PRIMARY KEY (`id`)
);

CREATE TABLE `user_location` (
	`id` INT(255) NOT NULL AUTO_INCREMENT,
	`name` varchar(255),
	`lat` FLOAT(3),
	`lng` FLOAT(3),
	`time_start` varchar(100),
	`last_updated` DATETIME(3),
	PRIMARY KEY (`id`)
);

CREATE TABLE `cost_detail` (
	`id` INT(255) NOT NULL AUTO_INCREMENT,
	`user_id` INT(255) NOT NULL,
	`cab_id` INT(255) NOT NULL,
	`distance` FLOAT(3),
	`minute_travel` FLOAT(3),
	`final_cost` FLOAT(3),
	`last_updated` DATETIME(3),
	PRIMARY KEY (`id`)
);

