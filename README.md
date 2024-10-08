# turnBasedPVPWebServer

## About

This project is a web server for a turn-based PvP game that uses a card-based mechanic. Players can attack other players (currently limited to 2-player rooms), increase their defense, and heal during gameplay.

## How to Use

1. **Clone the repository**
   To get started, clone the repository using the following command:

   ```bash
   git clone https://github.com/Panzer-prime/turnBasedPVPWebServer.git
   ```

2. **Change directory to `/cmd`**
   
   Navigate to the `/cmd` directory where the main file is located:

   ```bash
   cd turnBasedPVPWebServer/cmd
   ```

3. **Run the server**
   
   Run the main Go file to start the server:

   ```bash
   go run main.go
   ```

4. **How to Play**

   - **Create a Room**

     To create a room, send a `POST` request to the following endpoint:

     ```
     localhost:{yourPort}/create-room
     ```

     Example JSON body for creating a room:

     ```json
     {
       "roomName": "exampleRoom",
       "roomPassword": "examplePass",
       "userName": "exampleUser"
     }
     ```

   - **Receive Room and Player IDs**

     The server will respond with a JSON object containing both the `playerID` and `roomID`, which are necessary for connecting via WebSocket.

     Example WebSocket connection string:

     ```
     ws://localhost:3000/connect?playerID={yourPlayerID}&roomID={yourRoomID}
     ```

   - **Game Initialization**

     At the start of the game, each player will receive an initial game state, including a pre-selected deck of cards. You can refer to the `jsonFile.json` for an example of the initial game state.

   - **Taking a Turn**

     During the game, the client must send a turn in JSON format to play. An example of a turn JSON object looks like this:

     ```json
     {
       "skip": "",      // a string indicating whether the player wants to skip their turn
       "cardsID": []    // an array of card IDs representing the cards the player plays
     }
     ```

   The `skip` field should contain a value if the player wishes to skip their turn, and `cardsID` should be an array of card IDs for the turn.

