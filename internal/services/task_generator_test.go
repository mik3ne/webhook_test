package services_test

import (
	"testing"
	"webhook/internal/services"
)

func Test_GenerateNoTasks(t *testing.T) {

	generator := services.NewTaskGenerator()
	genChan := generator.GenerateTasks(0)

	if len(genChan) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(genChan))
	}
}

// Using table tests here
func Test_GenerateTasks(t *testing.T) {

	tCases := []struct {
		Name          string
		ExpectedTasks int
	}{
		{
			Name:          "one task",
			ExpectedTasks: 1,
		},
		{
			Name:          "two tasks",
			ExpectedTasks: 2,
		},
		{
			Name:          "ten tasks",
			ExpectedTasks: 10,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.Name, func(t *testing.T) {
			generator := services.NewTaskGenerator()
			genChan := generator.GenerateTasks(tc.ExpectedTasks)

			if len(genChan) != tc.ExpectedTasks {
				t.Errorf("[%s] expected %d tasks, got %d", tc.Name, tc.ExpectedTasks, len(genChan))
			}
		})
	}

}
