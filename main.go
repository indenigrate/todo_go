package main

func main() {
	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	storage.load(&todos)
	cmdflags := NewCmdFlags()
	cmdflags.Execute(&todos)
	storage.Save(&todos)
}
