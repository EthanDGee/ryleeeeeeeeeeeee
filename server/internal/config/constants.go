package config

const (
	LICHESS_BASE_URL    = "https://database.lichess.org/standard/"
	GAME_COUNTS_URL     = LICHESS_BASE_URL + "counts.txt"
	CACHE_FOLDER        = ".cache/"
	LOCAL_DATABASE_PATH = CACHE_FOLDER + "app.db"
	FILES_FOLDER        = CACHE_FOLDER + "files"
	GAME_COUNT          = 5_000_000
	BATCH_SIZE          = 500
)
