package awsssmenv

import "github.com/joho/godotenv"

func Load(filename string) error {
	return godotenv.Load(filename)
}
