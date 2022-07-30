package constanta

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load("config.env")

var (
	PORT   = os.Getenv("PORT")
	BROKER = os.Getenv("BROKER")
)
