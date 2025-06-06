package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/rpc"
	"os"
	"time"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
	for {
		args := TaskRequest{} // example worker ID
		reply := TaskReply{}
		ok := call("Coordinator.RequestTask", &args, &reply)
		if ok {
			log.Printf("Task assigned: %s\n", reply.Message)
			if reply.IsTask && reply.TaskType == "map" {
				MapWork(mapf, reply.Message, reply.NReduce, reply.FileID)
			}
		} else {
			log.Printf("call failed!\n")
		}
		time.Sleep(time.Second * 1)
	}
}

func MapWork(mapf func(string, string) []KeyValue, filename string, NReduce int, mapID int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	kva := mapf(filename, string(content))

	// create files and encoder
	files := make([]*os.File, NReduce)
	encoders := make([]*json.Encoder, NReduce)
	for r := 0; r < NReduce; r++ {
		name := fmt.Sprintf("mr-%d-%d", mapID, r)
		f, err := os.Create(name)
		if err != nil {
			log.Fatalf("cannot create %v", name)
		}
		files[r] = f
		encoders[r] = json.NewEncoder(f)
	}

	// write each kv into the correct file
	for _, kv := range kva {
		r := ihash(kv.Key) % NReduce
		err := encoders[r].Encode(&kv)
		if err != nil {
			log.Fatalf("Cannot encode %v", kv)
		}
	}

	// close file
	for _, f := range files {
		f.Close()
		log.Printf("Successfully close %v", f.Name())
	}

}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.AssignTask", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
