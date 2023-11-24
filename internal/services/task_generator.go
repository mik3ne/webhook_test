package services

type TaskGenerator struct {
}

func NewTaskGenerator() *TaskGenerator {
	return &TaskGenerator{}
}

func (t *TaskGenerator) GenerateTasks(requestAmount int) chan int {

	tasksChan := make(chan int, requestAmount)
	for i := 0; i < requestAmount; i++ {
		tasksChan <- i
	}
	return tasksChan
}
