package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type gameSnapshot struct {
	id          int    `json:"id"`
	raw_game_id int    `json:"raw_game_id"`
	move_number int    `json:"move_number"`
	turn        string `json:"turn"`
	move        int    `json:"move"`
	fen         string `json:"fen"`
}

var snapshots = []gameSnapshot{
	{id: 0, raw_game_id: 2, move_number: 3, turn: "k32a", move: 0, fen: "dafsd"},
	{id: 0, raw_game_id: 21, move_number: 13, turn: "k32a", move: 0, fen: "dafsd"},
	{id: 1, raw_game_id: 23, move_number: 30, turn: "k32a", move: 0, fen: "dafsd"},
	{id: 2, raw_game_id: 20, move_number: 31, turn: "k32a", move: 0, fen: "dafsd"},
	{id: 3, raw_game_id: 24, move_number: 33, turn: "k32a", move: 0, fen: "dafsd"},
}

func getSnapshot(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, snapshots)
}

func main() {
	fmt.Println("Hello, World!")
	router := gin.Default()
	router.GET("/snapshots", getSnapshot)

	router.Run("localhost:8080")
}
