package grpcRouter

import (
	"challenge/internal/cache"
	"challenge/internal/converter"
	"challenge/internal/entity"
	pb "challenge/internal/grpc"
	"challenge/internal/logger"
	"challenge/internal/repository"
	"context"
	"encoding/json"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
)

type TaskService struct {
	pb.UnimplementedTaskServiceServer
	repository repository.Repository
	cache      *cache.Cache
}

func NewService(repo repository.Repository, cache *cache.Cache) *TaskService {
	return &TaskService{repository: repo,
		cache: cache}
}

func (ts *TaskService) ListTask(_ context.Context, _ *empty.Empty) (*pb.TaskList, error) {

	t := new(pb.TaskList)
	t.Tasks = make([]*pb.Task, 0)
	task, err := ts.repository.ListTask()
	if err != nil {
		logger.Log().Warn("%v", err)
		return nil, err
	}

	for i := range task {
		aux := converter.EntitytoPb(task[i])
		t.Tasks = append(t.Tasks, aux)
	}

	return t, nil

}

func (ts *TaskService) GetTask(_ context.Context, taskPb *pb.Task) (*pb.Task, error) {

	id := taskPb.GetId()

	tk, err := ts.cache.Get(strconv.Itoa(int(id)))
	if err == nil {
		aux := converter.EntitytoPb(*tk)
		return aux, nil
	}
	task, err := ts.repository.GetTask(int(id))
	if err != nil {
		logger.Log().Warn("%v", err)
		return nil, err
	}

	jt, err := json.Marshal(task)
	if err != nil {
		logger.Log().Warn("%v", err)
	}

	err = ts.cache.Set(strconv.Itoa(int(id)), string(jt))
	if err != nil {
		logger.Log().Warn("coudn't set task: %v", err)
	}
	aux := converter.EntitytoPb(*task)

	return aux, nil

}

func (ts *TaskService) NewTask(_ context.Context, taskPb *pb.Task) (*pb.IntId, error) {

	i := new(pb.IntId)
	aux := converter.PbtoEntity(taskPb)
	id, err := ts.repository.NewTask(aux)
	if err != nil {
		logger.Log().Warn("%v", err)
		return nil, err
	}

	i.Id = int32(id)
	jt, err := json.Marshal(i)
	if err != nil {
		logger.Log().Warn("%v", err)
	}

	err = ts.cache.Set(strconv.Itoa(int(i.Id)), string(jt))
	if err != nil {
		logger.Log().Warn("coudn't set task: %v", err)
	}
	return i, nil

}

func (ts *TaskService) UpdateTask(_ context.Context, taskPb *pb.Task) (*empty.Empty, error) {

	i := new(pb.IntId)
	id := taskPb.GetId()
	err := ts.repository.UpdateTask(int(id))
	if err != nil {
		logger.Log().Warn("%v", err)
		return nil, err
	}

	tk, err := ts.cache.Get(strconv.Itoa(int(id)))
	if err == nil {
		jt, err := json.Marshal(tk)
		if err != nil {
			logger.Log().Warn("%v", err)
		}

		err = ts.cache.Set(strconv.Itoa(int(id)), string(jt))
		if err != nil {
			err = ts.cache.Del(strconv.Itoa(int(id)))
			if err != nil {
				logger.Log().Warn("%v", err)
			}
		}
	}

	i.Id = int32(0)
	return &empty.Empty{}, nil

}

func (ts *TaskService) ListComp(_ context.Context, taskPb *pb.Task) (*pb.TaskList, error) {
	t := new(pb.TaskList)
	t.Tasks = make([]*pb.Task, 0)

	var task []entity.Task
	var err error
	comp := taskPb.GetCompleted()
	if comp == "true" {
		task, err = ts.repository.ListComp("yes")
	} else {
		task, err = ts.repository.ListComp("no")
	}

	if err != nil {
		logger.Log().Warn("%v", err)
		return nil, err
	}

	for i := range task {
		aux := converter.EntitytoPb(task[i])
		t.Tasks = append(t.Tasks, aux)
	}

	return t, nil

}

func (ts *TaskService) mustEmbedUnimplementedTaskServiceServer() {
	panic("not implemented") // TODO: Implement
}
