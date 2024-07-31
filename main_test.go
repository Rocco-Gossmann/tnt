package main_test

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

func panicCatcher(test *testing.T, msg string, args ...any) func() {
	return func() {
		if p := recover(); p != nil {
			test.Fatal(msg, p, args)
		}
	}
}

func TestEverything(test *testing.T) {

	const DB_FILE = "./testdata.db"
	const TEST_TASK_1 = "test case 1"
	const TEST_TASK_2 = "another case"

	var taskId1 uint

	log.SetOutput(io.Discard)

	var handleunexpectedError = func(err error, wherewhy string) {
		if err != nil {
			test.Fatal("encountered unexpected error", wherewhy, ":", err)
		}
	}

	test.Run("Creating db", func(test *testing.T) {
		test.Helper()
		os.Remove(DB_FILE)
		database.InitDB(DB_FILE)
	})

	test.Run("Creating a new Task", func(test *testing.T) {
		err := database.AddTask(TEST_TASK_1)
		handleunexpectedError(err, "while creating task "+TEST_TASK_1)

		test.Run("Get taskID", func(t *testing.T) {
			defer panicCatcher(t, "failed to get id for '"+TEST_TASK_1+"'")()
			taskId1 = database.GetTaskIDByName(TEST_TASK_1)
		})

	})

	test.Run("prevent doubles", func(test *testing.T) {

		checkErr := func(err error) {

			if err == nil {
				test.Fatal("failed to stop '", TEST_TASK_1, "' from being added twice")
			}

			if err.Error() != "constraint failed: UNIQUE constraint failed: tasks.textkey (2067)" {
				test.Fatal("Wrong Error messages: Expected: 'UNIQUE constraint failed: tasks.textkey' => Got '", err.Error(), "'")
			}
		}

		test.Run("exact same taskKey", func(test *testing.T) {
			checkErr(database.AddTask(TEST_TASK_1))
		})

		test.Run("same task different casing", func(test *testing.T) {
			checkErr(database.AddTask(strings.ToUpper(TEST_TASK_1)))
		})

		test.Run("same task different spaces/tabs around", func(test *testing.T) {
			checkErr(database.AddTask(" " + TEST_TASK_1 + "\t"))
		})

		test.Run("same task different spaces/tabs + casing", func(test *testing.T) {
			checkErr(database.AddTask(" " + strings.ToUpper(TEST_TASK_1) + "\t"))
		})
	})

	test.Run("Make sure task is listed by ls", func(test *testing.T) {
		tasks, err := database.GetTaskList("")
		handleunexpectedError(err, "while listing tasks")

		if len(tasks) != 1 {
			test.Fatal("Tasklist returned nothing. Should have returned at least on entry")
		}

		if TEST_TASK_1 != tasks[0].Name {
			test.Fatal("returned enty is messed up\nShould be: ", TEST_TASK_1, "\nIs: ", tasks[0])
		}
	})

	test.Run("make sure times are empty to start", func(test *testing.T) {

		times, err := database.GetTimes(0)
		handleunexpectedError(err, "while listing Times")

		if len(times) > 0 {
			test.Fatal("times table already contains data, when it should not")
		}

	})

	test.Run("start the timer for a new task", func(t *testing.T) {
		defer panicCatcher(t, "failed to start timer because DB")()
		database.StartNewTime(taskId1)
	})

	test.Run("make sure task timer is running", func(t *testing.T) {
		defer panicCatcher(t, "encountered error while checking for running tasks")()
		if !database.TimedTaskIsRunning(taskId1) {
			test.Fatal("task is not running yet")
		}
	})

	test.Run("make sure timer is not listed in sums, while it is running", func(t *testing.T) {
		defer panicCatcher(t, "encountered error while checking for time sums")()

		lst := database.GetTimeSums(0)
		if len(lst) > 0 {
			test.Fatal("task yelds a sum", lst, len(lst))
		}
	})

	test.Run("stop all timers", func(t *testing.T) {
		defer panicCatcher(t, "db issue on ending all tasks")
		database.FinishCurrentlyRunningTimes()
	})

	test.Run("make sure task timer is not running anymore", func(t *testing.T) {
		defer panicCatcher(t, "encountered error while checking for running tasks")()
		if database.TimedTaskIsRunning(taskId1) {
			test.Fatal("task is not supposed to be running anymore")
		}
	})

	test.Run("timer should now be included in sums", func(t *testing.T) {
		defer panicCatcher(t, "encountered error while checking for time sums")()

		lst := database.GetTimeSums(0)
		if len(lst) == 0 {
			test.Fatal("task did not yeld a sum")
		}
	})

	test.Cleanup(func() {
		database.DeInitDB()
		os.Remove(DB_FILE)
	})
}
