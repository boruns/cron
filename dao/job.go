package dao

import (
	"crontab/global"
	"crontab/models"
	"crontab/utils"

	"gorm.io/gorm"
)

type JobDaoInterface interface {
	GetJobByJobId(jobId string) (*models.JobModel, error)
	GetJobList(page, pageSize int) (int64, []models.JobModel)
	SaveJob(db *gorm.DB, job *models.Job) (string, error)
	DestoryJobByJobId(db *gorm.DB, jobId string) error
	UpdateJobByJobId(db *gorm.DB, jobId string, jobData map[string]interface{}) error
}

func NewJobHandler() JobDaoInterface {
	return &jobDao{}
}

type jobDao struct {
}

func (j *jobDao) UpdateJobByJobId(db *gorm.DB, jobId string, jobData map[string]interface{}) error {
	return db.Model(&models.JobModel{}).Where("job_id = ?", jobId).Updates(jobData).Error
}

func (j *jobDao) DestoryJobByJobId(db *gorm.DB, jobId string) error {
	return db.Where(&models.JobModel{JobId: jobId}).Delete(&models.JobModel{}).Error
}

func (j *jobDao) GetJobList(page, pageSize int) (int64, []models.JobModel) {
	var jobs []models.JobModel
	var count int64
	global.DB.Model(&models.JobModel{}).Count(&count)
	global.DB.Model(&models.JobModel{}).Scopes(models.Paginate(page, pageSize)).Find(&jobs)
	return count, jobs
}

//通过jobId获取模型
func (j *jobDao) GetJobByJobId(jobId string) (*models.JobModel, error) {
	job := &models.JobModel{}
	if err := global.DB.Where(&models.JobModel{JobId: jobId}).First(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

//保存任务
func (j *jobDao) SaveJob(db *gorm.DB, job *models.Job) (string, error) {
	uuid, _ := utils.GetUUID()
	jobModel := &models.JobModel{
		JobId:      uuid,
		Name:       job.Name,
		Command:    job.Command,
		CronExpire: job.CronExpire,
		UpdatedAt:  utils.GetNow(),
		CreatedAt:  utils.GetNow(),
	}
	if err := db.Create(jobModel).Error; err != nil {
		return "", err
	}
	return uuid, nil
}
