
# Minesweeper-API
Minesweeper game written as a simple REST API using Golang.

It uses Redis as database although it is designed to easily add another database vendor.

It provides the ability to easily change underlying http framework (e.g., switch from Chi to Mux http router or viceversa)

The game engine has been written by adapting the classical [Flood Fill algorithm](https://en.wikipedia.org/wiki/Flood_fill).

## API Endpoints

### POST - Create User
```
http://localhost:8080/users
```
Create a user for playing. A user should be created in order to be able to start a new game.

**Body**
```json
{
	"username": "myuser"
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
### PUT - Start/Restart Game
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
### POST - Click
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
### GET - Obtain the Game Board
```
http://localhost:8080/games/mygame1/myuser/board
```
Get the board in JSON format.

**Example Request**
```
curl --location --request GET 'http://localhost:8080/games/mygame1/myuser/board'
```

**Example Response**
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
	... omitted some data to shorten
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
Using Postman? You can pretty render the reponse using the Postman Visualize feature. 
Follow this steps to render the game board in Postman:
1. Add [this](https://gist.github.com/arllanos/6a57c6b293c0c7280562aef3d97eb248) code to the `Tests` script for the request.
2. Click `Send` to run the request.
3. Click the `Visualize` tab to render the game board.

## Game engine logic and how to interpret the board
The game **board** is part of the Game structure `types/game.go`.
This board is a 2-Dimensional array of bytes an its data is coded as follows:

**Unrevealed values**
```
* M -----> Veiled mine 
* E -----> Veiled Empty 
```
By clicking on a cell (using the click endpoint) it will reveal the corresponding cell.

**Revealed values**
```
* X -----> Exploded mine
* B -----> Revealed Blank
* Digit -> Revealed with adjacent mine count
```
**Logic**
- When a mine ('M') is clicked, it changes to ('X') and the game is over.
- When a Veiled Empty ('E') cell is clicked, it can either, transition to Revealed blank ('B') or to a digit (1 to 8) indicating the number of adjacent mines.
- Flagging a cell will change any unrevealed value to its corresponding lowercase letter value.