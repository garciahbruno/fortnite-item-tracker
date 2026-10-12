package web

import "embed"

//go:embed index.html styles.css app.js
var assets embed.FS

func Read(name string) []byte {
	data, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return data
}
