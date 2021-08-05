package forms

type JobStoreRequest struct {
	Name       string `json:"name" binding:"required"`
	Command    string `json:"command" binding:"required"`
	CronExpire string `json:"cron_expire" binding:"required"`
}

type JobDestoryRequest struct {
	JobId string `json:"jobId" uri:"jobId" binding:"required"`
}

type JobUpdateRequest struct {
	JobId      string `json:"jobId"`
	Name       string `json:"name" binding:"required"`
	Command    string `json:"command" binding:"required"`
	CronExpire string `json:"cron_expire" binding:"required"`
}

type JobListRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type JobLogRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}
