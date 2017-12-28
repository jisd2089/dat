package distribute

// 任务仓库
type TaskBase struct {
	Tasks chan *Task
}

func NewTaskBase() *TaskBase {
	return &TaskBase{
		Tasks: make(chan *Task, 1024),
	}
}

// 服务器向仓库添加一个任务
func (self *TaskBase) Push(task *Task) {
	id := len(self.Tasks)
	task.Id = id
	self.Tasks <- task
}

// 客户端从本地仓库获取一个任务
func (self *TaskBase) Pull() *Task {
	return <-self.Tasks
}

// 仓库任务总数
func (self *TaskBase) Len() int {
	return len(self.Tasks)
}

// 主节点从仓库发送一个任务
func (self *TaskBase) Send(clientNum int) Task {
	return *<-self.Tasks
}

// 从节点接收一个任务到仓库
func (self *TaskBase) Receive(task *Task) {
	self.Tasks <- task
}
