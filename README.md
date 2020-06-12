
# Minesweeper-API
Minesweeper game written as a simple REST API using Golang.

It uses Redis as database although it is designed to easily swap out to other database vendors (e.g, postgresql, sqllite)

It provides the ability to easily change underlying http routing framework (e.g., switch from Chi to Mux or viceversa)

The [game engine]([https://github.com/arllanos/minesweeper-API#game-engine-logic-and-how-to-interpret-the-board](https://github.com/arllanos/minesweeper-API#game-engine-logic-and-how-to-interpret-the-board)) has been written by adapting the classical [Flood Fill algorithm](https://en.wikipedia.org/wiki/Flood_fill).



## Running the application locally
There are two ways to run the API server locally. 
1. Running the app using docker-compose (recommended). 
2. Build and running locally.

### Option 1. Running using docker-compose
Starting the server
```
docker-compose up
```
Stopping the server
```
docker-compose down --remove-orphans
```
### Option 2. Building and running locally
Pre-requisites: 
- If Golang is not installed click [here](https://golang.org/doc/install ) for installation instructions
- Provision a local Redis instance.

#### Build
```
go build -o minesweeper-api
```
#### Run
```
go run .
```
## REST API Endpoints

### Create User

**POST**  `http://localhost:8080/users`

Creates a user for playing. The user should be created before starting a new game.

**Body**
```json
{
	"username": "player1"
}
```
**Example Request**
```bash
curl --location --request POST 'http://localhost:8080/users' \
--data-raw '{
	"username": "player1"
}'
```
**Example Response**
```json
{
    "username": "player1",
    "createdAt": "2020-06-11T13:03:30.917771715-03:00"
}
```
### Start/Restart Game

**PUT** `http://localhost:8080/## Headinggames`

**Body**
```json
{
	"name": "game1",
	"username": "player1",
	"rows": 7,
	"cols": 7,
	"mines": 5
}
```
**Example Request**
```bash
curl --location --request PUT 'http://localhost:8080/games' \
--data-raw '{
	"name": "game1",
	"username": "player1",
	"rows": 7,
	"cols": 7,
	"mines": 5
}'
```
**Example Response**
```json
{
    "name": "game1",
    "username": "player1",
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
### Click

**POST** `http://localhost:8080/games/game1/player1/click`

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
curl --location --request POST 'http://localhost:8080/games/game1/player1/click' \ --data-raw '{ "row": 1, "col": 0, "kind": "click" }'
```
**Example Response**
```json
{
    "name": "game1",
    "username": "player1",
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
### Obtain the Game Board

**GET** `http://localhost:8080/games/game1/player1/board`

Get the board in JSON format.

**Example Request**
```
curl --location --request GET 'http://localhost:8080/games/game1/player1/board'
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
**Note!**
Using Postman? You can pretty render the response using the Postman Visualize feature. 
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
