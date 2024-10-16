package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/CamusSisyphus/blog_aggregator/internal/config"
	"github.com/CamusSisyphus/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *(config.Config)
}

func main() {

	c, err := config.Read()

	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", c.DBURL)
	if err != nil {
		fmt.Println(err)
	}

	dbQueries := database.New(db)

	s := &state{db: dbQueries, cfg: &c}

	commands := commands{make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not Enough Arguments Providedbootdev run dca1352a-7600-4d1d-bfdf-f9d741282e55")
	}
	commandName := args[1]
	commandArgs := args[2:]

	err = commands.run(s, command{Name: commandName, Args: commandArgs})
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
	// handlerLogin(&s, command{name: "handlerLogin", args: []string{"Mason"}})

}
