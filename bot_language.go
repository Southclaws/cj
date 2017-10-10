package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"go.uber.org/zap"
)

// Locale stores map from language ID to Language object
type Locale struct {
	Languages map[string]Language
}

// Language stores map of key to text body
type Language struct {
	Strings map[string]string
}

// GetLangString returns a text body based on language ID and key.
func (l Locale) GetLangString(lang string, key string, vargs ...interface{}) string {
	str := fmt.Sprintf(l.Languages[lang].Strings[key], vargs...)
	if str == "" {
		logger.Warn("undefined lang key", zap.String("key", key))
	}
	return str
}

// LoadLanguages reads language definition files form ./lang
func (app *App) LoadLanguages() {
	files, err := ioutil.ReadDir("lang")
	if err != nil {
		panic(err)
	}

	var key string
	app.locale.Languages = make(map[string]Language)

	for _, f := range files {
		if f.IsDir() {
			key = f.Name()
			la := loadLanguageFromDir(filepath.Join("lang", key))
			app.locale.Languages[key] = la
		}
	}
}

func loadLanguageFromDir(dir string) Language {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var language Language
	var key string

	language = Language{
		Strings: make(map[string]string),
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		key = f.Name()
		str := loadLanguageStringFromFile(filepath.Join(dir, key))
		language.Strings[key] = str
	}

	return language
}

func loadLanguageStringFromFile(file string) string {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return string(contents)
}
