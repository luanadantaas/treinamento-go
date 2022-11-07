package converter

import (
	"challenge/internal/entity"
	pb "challenge/internal/grpc"
)

//transforma informacoes no formato protobuf para ser lida em go
func PbtoEntity(task *pb.Task) (entity.Task){
	return entity.Task{
		ID: int(task.Id), Name: task.Name, Completed: task.Completed,
	}
}

//transforma informações em go para ser lida pelo protobuf
func EntitytoPb(task entity.Task) (*pb.Task){
	return &pb.Task{
		Id: int32(task.ID), Name: task.Name, Completed: task.Completed,
	}
}

