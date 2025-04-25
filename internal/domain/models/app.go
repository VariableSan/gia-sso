package models

type App struct {
	ID   int
	Name string
	// It is advisable not to pass the secret through the model
	Secret string
}
