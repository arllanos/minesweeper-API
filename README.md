# Minesweeper-API
Minesweeper game written as a simple REST API using Golang.

It uses Redis as database although it is designed to easily swap out to other database vendors (e.g., postgresql, sqllite)

It provides the ability to easily change underlying http routing framework (e.g., switch from Chi to Mux or viceversa)

The game engine has been written by adapting the classical [Flood Fill algorithm](https://en.wikipedia.org/wiki/Flood_fill).

**Technical highlights**
- Follows clean architecture principles
- Independent of http framework and data repository technology
- Use Docker
- Deployable in Kubernetes (`deployment/k8s-deployment.yaml`)

## Running the application locally
There are two ways to run the API server locally. 
1. Running the app using docker-compose (recommended). 
2. Build and running locally.

### Option 1. Running using docker-compose
Starting the server
```
docker compose up -d
```
Stopping the server
```
docker compose down --remove-orphans
```
### Option 2. Building and running locally
Pre-requisites: 
- If Golang is not installed, click [here](https://golang.org/doc/install ) for installation instructions
- Provision and run a local Redis instance.

#### Build
```
make build
```
#### Run
```
make run
```
## API Endpoints
### Create User

Creates a user for playing. The user should be created before starting a new game.

**POST**  `http://localhost:8080/users`

| Code | Description  |
| ---- | ------------ |
| 201  | User created |
| 400  | Bad request  |
| 500  | Server error |

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

Starts a new game or restart a game

**PUT** `http://localhost:8080/games`

| Code | Description  |
| ---- | ------------ |
| 201  | Game created/restarted |
| 400  | Bad request  |
| 500  | Server error |

**Body**
```json
{
	"name": "game1",
	"username": "player1",
	"rows": 4,
	"cols": 4,
	"mines": 5
}
```
**Example Request**
```bash
curl --location --request PUT 'http://localhost:8080/games' \
--data-raw '{
	"name": "game1",
	"username": "player1",
	"rows": 4,
	"cols": 4,
	"mines": 5
}'
```
**Example Response**
```json
{
    "name": "game1",
    "username": "player1",
    "rows": 4,
    "cols": 4,
    "mines": 5,
    "status": "ready",
    "board": [
        "RU1FRUVNRQ==",
        "RUVFRUVFRQ==",
        "RUVFRUVFRQ==",
        "RUVFRUVFTQ=="
    ],
    "clicks": 0,
    "created_at": "2020-06-11T13:05:54.943472481-03:00",
    "started_at": "0001-01-01T00:00:00Z",
    "time_spent": 0
}
```
### Click

Click or flag a cell in the game board. Use the `kind` field to indicate either `click` or `flag`

**POST** `http://localhost:8080/games/game1/player1/click`

| Code | Description  |
| ---- | ------------ |
| 200  | Successfully applied click on cell |
| 400  | Bad request (wrong click kind, already won / lost) |
| 404  | Cell not found |
| 500  | Server error |


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
        "RUVFRUVFTQ=="
    ],
    "clicks": 1,
    "created_at": "2020-06-11T13:05:54.943472481-03:00",
    "started_at": "2020-06-11T13:06:30.513938447-03:00",
    "time_spent": 700
}
```
### Get the Game Board

Get the board in JSON format.

**GET** `http://localhost:8080/games/game1/player1/board`

| Code | Description  |
| ---- | ------------ |
| 200  | OK |
| 404  | User / Game not found |
| 500  | Server error |

**Example Request**
```
curl --location --request GET 'http://localhost:8080/games/game1/player1/board'
```

**Example Response**
```json
[
    ["E", "M", "M", "E"],
    ["E", "2", "E", "M"],
    ["E", "E", "E", "E"],
    ["E", "M", "M", "E"]
]
```
**Visualizing board in Postman**

To render the GET board response using the Postman Visualization feature, follow these steps:
1. In Postman, add following code to the `Scripts > Post-response` tab.
    ```js
    let jsonResponse = pm.response.json();

    // Generate an HTML table from the JSON array
    let table = '<table border="1" style="border-collapse: collapse; text-align: center;">';
    for (let row of jsonResponse) {
        table += '<tr>';
        for (let cell of row) {
            table += `<td style="padding: 10px; width: 30px;">${cell}</td>`;
        }
        table += '</tr>';
    }
    table += '</table>';

    // Set the visualization
    pm.visualizer.set(table);
    ```
2. Click `Send` to run the request.
3. Click the `Visualize` tab to render the game board.

## Game engine logic and how to interpret the board
The game **board** is part of the Game structure `internal/domain/game.go`.
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
- When a cell transition to 'B' it triggers a **recursive** reveal of adjacent cells.
- Flagging a cell will change any unrevealed value to its corresponding lowercase letter value.
