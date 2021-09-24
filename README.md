# GifCreator

Easily create GIFs from a series of images in a Telegram bot.
Link of the official instance: https://t.me/CreateGIFsBot

## Features

- Multiple photo types supported
- Matomo Analytics Support
- Customizable via environment variables

## Deploy

### Replit/Heroku

1. Set env variables as explained below

   Note: You might need to rely on a MySQL As A Service for the database.

2. Deploy as you would normally

### Docker

1. Copy `.example_env` to `.env`
2. Edit it according to the env variables explanation below. Leave the defaults for the database credentials!
3. Execute `docker-compose up`

### Environment Variables Explaination

| Variable Name  | Explanation                                                  | ⚠️*   | Sample Value                                      |
| -------------- | ------------------------------------------------------------ | ---- | ------------------------------------------------- |
| TOKEN          | The bot token @BotFather gave you.                           |      | 8babd9a8dd:HHOf03z2Mi6sYqOaJsKZrhQ33o5S7Y2go_a    |
| DB_USERNAME    | Username of the user accessing the MySQL database            | ✖️    | user                                              |
| DB_PASSWORD    | Password of the user accessing the MySQL database            | ✖️    | password                                          |
| DB_ADDRESS     | Address (either IP or host, with port if not default) of the database | ✖️    | db.radommysqlhost.com:1234                        |
| DB_NAME        | Name of the database where the user data will be stored in the table 'users' | ✖️    | db                                                |
| DEV_LINK       | Link to the developer's website or channel                   |      | https://mywebsite.com                             |
| DEV_CHANNEL    | Tag of the developer's channel                               |      | @mychannel                                        |
| SOURCE         | Link to the repository where the source code is stored       |      | https://mygitprovider.com/myusername/myrepository |
| PRIVACY        | Link to the privacy policy detailing how the user data is used in the bot |      | https://mywebsite.com/privacy                     |
| TERMS          | Link to the terms of service that regulate how the user shall interact with the bot |      | https://mywebsite.com/terms                       |
| MATOMO_ENABLED | Enable Matomo Analytics? 1 = yes, 0 = no                     |      | 1                                                 |
| MATOMO_SITEID  | Site ID used by Matomo                                       |      | 1                                                 |
| MATOMO_HOST    | Host where Matomo is located                                 |      | matomo.mywebsite.com                              |

*If you're deploying with Docker, don't change these values from the sample configuration file!

## License

![GPLv3 Logo.svg](https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/GPLv3_Logo.svg/220px-GPLv3_Logo.svg.png)

```
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
```