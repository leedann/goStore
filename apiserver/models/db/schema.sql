CREATE USER pgstest WITH SUPERUSER;
create table users (
	ID serial primary key,
	Email varchar(255),
	PassHash varchar(255),
	UserName varchar(100),
    FirstName varchar(50),
    LastName varchar(50),
    PhotoURL varchar(100),
    MobilePhone varchar(12)
);