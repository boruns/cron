package models

type Job struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	CronExpire string `json:"corn_expire"`
}

type JobModel struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	JobId      string `json:"job_id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	CronExpire string `json:"cron_expire"`
	UpdatedAt  string `json:"updated_at"`
	CreatedAt  string `json:"created_at"`
}

type JobLog struct {
	JobId        string `bson:"job_id" json:"job_id"`
	JobName      string `bson:"job_name" json:"job_name"`
	Command      string `bson:"command" json:"command"`
	Err          string `bson:"err" json:"err"`
	Output       string `bson:"output" json:"output"`
	PlanTime     int64  `bson:"plan_time" json:"plan_time"`
	ScheduleTime int64  `bson:"schedule_time" json:"schedule_time"`
	StartTime    int64  `bson:"start_time" json:"start_time"`
	EndTime      int64  `bson:"end_time" json:"end_time"`
}

func (j *JobModel) TableName() string {
	return "jobs"
}
