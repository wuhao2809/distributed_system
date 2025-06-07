package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Coordinator struct {
	MapTasks []MapTask // list of map tasks
	NReduce  int       // number of reduce tasks
}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) RequestTask(args *TaskRequest, reply *TaskReply) error {
	// Check if there are any idle map tasks
	for i, mapTask := range c.MapTasks {
		if mapTask.FileStatus == "idle" {
			c.MapTasks[i].FileStatus = "in-process"

			reply.FileID = mapTask.FileID
			reply.IsTask = true
			reply.TaskType = "map"
			reply.Message = mapTask.Filename
			reply.NReduce = c.NReduce
			return nil
		}
	}
	reply.IsTask = false
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, NReduce int) *Coordinator {
	c := Coordinator{}
	for i, file := range files {
		task := MapTask{
			Filename:   file,
			FileID:     i,
			FileStatus: "idle",
		}
		c.MapTasks = append(c.MapTasks, task)
		log.Printf("Added map task for file %s with ID %d", file, i)
		c.NReduce = NReduce
	}
	// Your code here.

	c.server()
	return &c
}
