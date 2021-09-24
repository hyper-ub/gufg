package main

/*

	GifCreator - Telegram Bot to create GIFs from a series of images.
	Copyright (C) 2021  MassiveBox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.

 */

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func databaseConnect() (*sql.DB, error) {
	return sql.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@("+os.Getenv("DB_ADDRESS")+")/"+os.Getenv("DB_NAME"))
}

func createTable(db *sql.DB) {
	db.Exec("CREATE TABLE IF NOT EXISTS `"+os.Getenv("DB_NAME")+"`.`users` ( `id` INT NOT NULL AUTO_INCREMENT , `user_id` BIGINT NOT NULL DEFAULT '0' , `username` VARCHAR(64) NOT NULL DEFAULT '' , `page` VARCHAR(32) NOT NULL DEFAULT '' , `frames` MEDIUMTEXT NOT NULL , `has_voted` TINYINT NOT NULL DEFAULT '0' , `registration_time` BIGINT NOT NULL DEFAULT '0' , PRIMARY KEY (`id`)) ENGINE = InnoDB;")
}

type userData struct {
	ID int
	UserID int64
	Username string
	Page string
	Frames string
	HasVoted bool
	RegistrationTime int64
}

func databaseRegisterUser(db *sql.DB, data userData) {
	db.Exec("INSERT INTO `users` (`id`, `user_id`, `username`, `page`, `frames`, `has_voted`, `registration_time`) VALUES (NULL, '"+fmt.Sprintf("%d", data.UserID)+"', '"+data.Username+"', '', '', '0', '0');")
}

func databaseGetUserData(db *sql.DB, userID int64) (userData, error) {

	rows, err := db.Query("SELECT * FROM `users` WHERE `user_id` = ?", userID)
	if err != nil {
		return userData{}, err
	}
	defer rows.Close()
	rows.Next()
	data := userData{}
	rows.Scan(&data.ID, &data.UserID, &data.Username, &data.Page, &data.Frames, &data.HasVoted, &data.RegistrationTime)

	return data, nil
}

func databaseUpdateUserData(db *sql.DB, data userData) {

	var strHasVoted string
	if data.HasVoted {
		strHasVoted = "1"
	}else{
		strHasVoted = "0"
	}

	db.Exec("UPDATE `users` SET `username` = '"+data.Username+"', `page` = '"+data.Page+"', `frames` = '"+data.Frames+"', `has_voted` = '" + strHasVoted + "' WHERE `users`.`user_id` = " + fmt.Sprintf("%d", data.UserID))

}