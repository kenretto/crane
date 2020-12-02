package timer

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/kenretto/crane/util/stack"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	// ErrJobExist 同名字的任务已经被添加过了
	ErrJobExist = errors.New("job has been exist")
)

// JobName 任务名称
type JobName string

// JobInterface 定时任务接口
type JobInterface interface {
	// Name 定时任务名称, 保持唯一
	Name() JobName
	// Description 对定时任务的描述
	Description() string
	// Spec 执行触发时间表达式, 下边是一些例子
	//  0 * * * * ?                      每1分钟触发一次
	//  0 0 * * * ?                      每天每1小时触发一次
	//  0 0 10 * * ?                     每天10点触发一次
	//  0 * 14 * * ?                     在每天下午2点到下午2:59期间的每1分钟触发
	//  0 30 9 1 * ?                     每月1号上午9点半
	//  0 15 10 15 * ?                   每月15日上午10:15触发
	//  */5 * * * * ?                    每隔5秒执行一次
	//  0 */1 * * * ?                    每隔1分钟执行一次
	//  0 0 5-15 * * ?                   每天5-15点整点触发
	//  0 0/3 * * * ?                    每三分钟触发一次
	//  0 0-5 14 * * ?                   在每天下午2点到下午2:05期间的每1分钟触发
	//  0 0/5 14 * * ?                   在每天下午2点到下午2:55期间的每5分钟触发
	//  0 0/5 14,18 * * ?                在每天下午2点到2:55期间和下午6点到6:55期间的每5分钟触发
	//  0 0/30 9-17 * * ?                朝九晚五工作时间内每半小时
	//  0 0 10,14,16 * * ?               每天上午10点，下午2点，4点
	//  0 0 12 ? * WED                   表示每个星期三中午12点
	//  0 0 17 ? * TUES,THUR,SAT         每周二、四、六下午五点
	//  0 10,44 14 ? 3 WED               每年三月的星期三的下午2:10和2:44触发
	//  0 15 10 ? * MON-FRI              周一至周五的上午10:15触发
	//  0 0 23 L * ?                     每月最后一天23点执行一次
	//  0 15 10 L * ?                    每月最后一日的上午10:15触发
	//  0 15 10 ? * 6L                   每月的最后一个星期五上午10:15触发
	//  0 15 10 * * ? 2005               2005年的每天上午10:15触发
	//  0 15 10 ? * 6L 2002-2005         2002年至2005年的每月的最后一个星期五上午10:15触发
	//  0 15 10 ? * 6#3                  每月的第三个星期五上午10:15触发
	Spec() string
	Runnable() bool
	Pause()
	// Start 再次启动
	Start()
	Run()
	Status() Status
}

// Status 定时任务状态
type Status int

const (
	// Run 定时任务处于正常可以执行状态
	Run Status = iota
	// Pause 定时任务处于暂停状态
	Pause
)

// MarshalJSON json
func (status Status) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + status.String() + `"`), nil
}

// String ...
func (status Status) String() string {
	switch status {
	case Run:
		return "正常"
	case Pause:
		return "暂停"
	}

	return "UnKnow"
}

// Task 单个定时任务的一些信息
type Task struct {
	EntryID cron.EntryID // 定时任务标识
	Name    JobName
	Job     JobInterface // 定时任务本身
}

// Run 运行
func (info *Task) Run() {
	defer func() {
		if err := recover(); err != nil {
			err := map[string]interface{}{
				"msg":   err,
				"stack": stack.Stack(3),
			}
			var s, _ = jsoniter.MarshalToString(err)
			driver.Set(info.Name, "error", s)
		}
	}()
	var runnable = info.Job.Runnable()
	if info.EntryID == 0 || info.Job.Status() == Pause || !runnable {
		return
	}
	driver.Set(info.Name, "prev", time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"))
	info.Job.Run()
}

// Status 状态
func (info *Task) Status() Status {
	return info.Job.Status()
}

// Tasks 定时任务集合
type Tasks struct {
	lock  sync.Mutex
	tasks map[JobName]*Task
}

// Pause 暂停某个任务
func (t *Tasks) Pause(name JobName) {
	t.tasks[name].Job.Pause()
}

// Start 再次启动
func (t *Tasks) Start(name JobName) {
	t.tasks[name].Job.Start()
}

// Exec 直接执行
func (t *Tasks) Exec(name JobName) {
	t.tasks[name].Run()
}

// All 所有任务列表
func (t *Tasks) All() map[JobName]*Task {
	return t.tasks
}

// GetJob 获取任务
func (t *Tasks) GetJob(name JobName) *Task {
	return t.tasks[name]
}

// AddJob 添加定时任务
func (t *Tasks) AddJob(jobInterface JobInterface) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.tasks == nil {
		t.tasks = make(map[JobName]*Task)
	}

	if _, ok := t.tasks[jobInterface.Name()]; ok {
		return ErrJobExist
	}

	t.tasks[jobInterface.Name()] = &Task{
		Name: jobInterface.Name(),
		Job:  jobInterface,
	}
	return nil
}

// Timer 计划任务
type Timer struct {
	Cron  *cron.Cron
	tasks *Tasks
}

// NewTimer 开启新的定时器
// 使用 quartz 规则, 这是一种 Java 定时任务同用的规则, 他可以精确到秒
//  see http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html
func NewTimer(log *logrus.Entry) *Timer {
	return &Timer{Cron: cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local),
		cron.WithLogger(&ILogger{log}),
		cron.WithChain(cron.Recover(&ILogger{log})),
	), tasks: new(Tasks)}
}

// Run 启动定时器
func (timer *Timer) Run() {
	timer.Cron.Start()
}

// Tasks 任务列表
func (timer *Timer) Tasks() *Tasks {
	return timer.tasks
}

// Prev 上次执行的时间
func (timer *Timer) Prev(name JobName) time.Time {
	var t, _ = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", driver.Get(name, "prev"))
	return t
}

// LastError 上次执行的时间
func (timer *Timer) LastError(name JobName) map[string]interface{} {
	var err map[string]interface{}
	_ = jsoniter.UnmarshalFromString(driver.Get(name, "error"), &err)
	return err
}

// Next 下一次执行的时间
func (timer *Timer) Next(name JobName) time.Time {
	return timer.Cron.Entry(timer.Tasks().GetJob(name).EntryID).Next
}

// AddJob 添加任务
func (timer *Timer) AddJob(job JobInterface) error {
	err := timer.tasks.AddJob(job)
	if err != nil {
		return err
	}

	var jobInfo = timer.tasks.GetJob(job.Name())
	entryID, err := timer.Cron.AddJob(jobInfo.Job.Spec(), jobInfo)
	if err != nil {
		return err
	}
	jobInfo.EntryID = entryID
	return err
}
