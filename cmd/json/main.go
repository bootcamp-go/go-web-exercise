package main

import (
	"app/internal/application"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// env
	// ...
	
	// application
	// - config
	dbFilePathLastId, err := strconv.Atoi(os.Getenv("DB_FILE_PATH_LAST_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg := &application.ConfigAppJSON{
		Addr:  os.Getenv("SERVER_ADDR"),
		Token: os.Getenv("API_TOKEN"),
		FilePath: os.Getenv("DB_FILE_PATH"),
		FilePathLastId: dbFilePathLastId,
		LayoutDate: os.Getenv("LAYOUT_DATE"),
	}
	app := application.NewApplicationJSON(cfg)
	// - set-up
	if err := app.SetUp(); err != nil {
		fmt.Println(err)
		return
	}

	// run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}