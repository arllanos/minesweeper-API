
# Minesweeper-API

## POST - Create User
```
http://localhost:8080/users
```
Create a user for playing. A user should be created in order to be able to start a new game.

**Body**
```json
{
	"username": "carouser"
}
```
**Example Request**
```bash
curl --location --request POST 'http://localhost:8080/users' \
--data-raw '{
	"username": "myuser"
}'
```
**Example Response**
```json
{
    "username": "myuser",
    "createdAt": "2020-06-11T13:03:30.917771715-03:00"
}
```
## PUT - Start/Restart Game
```
http://localhost:8080/games
```
**Body**
```json
{
	"name": "mygame1",
	"username": "myuser",
	"rows": 7,
	"cols": 7,
	"mines": 5
}
```
**Example Request**
```bash
curl --location --request PUT 'http://localhost:8080/games' \
--data-raw '{
	"name": "mygame1",
	"username": "myuser",
	"rows": 7,
	"cols": 7,
	"mines": 5
}'
```
**Example Response**
```json
{
    "name": "mygame1",
    "username": "myuser",
    "rows": 7,
    "cols": 7,
    "mines": 5,
    "status": "ready",
    "board": [
        "RU1FRUVNRQ==",
        "RUVFRUVFRQ==",
        "RUVFRUVFRQ==",
        "RUVFRUVFTQ==",
        "TUVNRUVFRQ==",
        "RUVFRUVFRQ==",
        "RUVFRUVFRQ=="
    ],
    "clicks": 0,
    "created_at": "2020-06-11T13:05:54.943472481-03:00",
    "started_at": "0001-01-01T00:00:00Z",
    "time_spent": 0
}
```
## POST - Click
```
http://localhost:8080/games/mygame1/myuser/click
```
Click or flag a cell in the game board. Use the `kind` field to indicate either `click` or `flag`
**Body**
```json
{
	"row": 1,
	"col": 0,
	"kind": "click"
}
```
**Example Request**
```bash
curl --location --request GET 'http://localhost:8080/games/mygame1/myuser/board'
```
**Example Response**
```json
{
    "name": "mygame1",
    "username": "myuser",
    "rows": 7,
    "cols": 7,
    "mines": 5,
    "status": "in_progress",
    "board": [
        "RU1FRUVNRQ==",
        "MUVFRUVFRQ==",
        "RUVFRUVFRQ==",
        "RUVFRUVFTQ==",
        "TUVNRUVFRQ==",
        "RUVFRUVFRQ==",
        "RUVFRUVFRQ=="
    ],
    "clicks": 1,
    "created_at": "2020-06-11T13:05:54.943472481-03:00",
    "started_at": "2020-06-11T13:06:30.513938447-03:00",
    "time_spent": 700
}
```
## GET - Obtain the Game Board
```
http://localhost:8080/games/mygame1/myuser/board
```
Get the board in JSON format.

**Example Request**
curl --location --request GET 'http://localhost:8080/games/mygame1/myuser/board'

**Example Response**

```
```json
[
    [
        "E",
        "M",
        "E",
        "E",
        "E",
        "M",
        "E"
    ],
	... ommited some data to shorten
    [
        "E",
        "E",
        "E",
        "E",
        "E",
        "E",
        "E"
    ]
]
```
Using Postman? You can pretty render using the Postman Visualize feature. For rendering in Postman [this](https://gist.github.com/arllanos/6a57c6b293c0c7280562aef3d97eb248) code in the Test tab
