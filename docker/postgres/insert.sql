CREATE TABLE users (ID SERIAL PRIMARY KEY NOT NULL,
		NAME text,
		PASSWORD text,
		EMAIL text
	);

INSERT INTO users(name, password, email) VALUES 
	('Peter Parker', 'spiderman', 'spiderman@marvel.com'),
	('Matt Murdoch', 'daredevil', 'daredevil@marvel.com'),
	('Bruce Banner', 'hulk', 'hulk@marvel.com'),
	('Stephen Strange', 'drstrange', 'drstrange@marvel.com'),
	('Reed Richards', 'mrfantastic', 'mrfantastic@marvel.com'),
	('Bruce Wayne', 'batman', 'batman@dc.com'), 
	('Clark Kent', 'superman', 'superman@dc.com'),
	('Barry Allen', 'flash', 'flash@dc.com'),
	('Allan Scott', 'greenlantern', 'greenlantern@dc.com');