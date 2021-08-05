package controllers

import (
	"crontab/forms"
	"crontab/models"
	"crontab/response"
	"crontab/services"
	"crontab/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type JobController struct {
}

//保存任务
func (j *JobController) Store(c *gin.Context) {
	request := forms.JobStoreRequest{}
	if err := c.ShouldBind(&request); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	jobService := services.NewJobServiceHandler()
	job := &models.Job{
		Name:       request.Name,
		Command:    request.Command,
		CronExpire: request.CronExpire,
	}
	jobId, err := jobService.SaveJob(c, job)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), "")
		return
	}
	response.Success(c, 200, "保存成功", gin.H{
		"job_id": jobId,
	})
}

//删除任务
func (j *JobController) Destory(c *gin.Context) {
	req := forms.JobDestoryRequest{}
	if err := c.ShouldBindUri(&req); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	jobService := services.NewJobServiceHandler()
	if err := jobService.DestoryJob(c, req.JobId); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), "")
		return
	}
	response.Success(c, 200, "删除成功", "")
}

//任务列表
func (j *JobController) Index(c *gin.Context) {
	req := &forms.JobListRequest{}
	if err := c.ShouldBind(req); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	jobService := services.NewJobServiceHandler()
	jobList, total := jobService.JobList(c, req.Page, req.PageSize)
	data := map[string]interface{}{
		"total":     total,
		"list":      jobList,
		"page":      req.Page,
		"page_size": req.PageSize,
	}
	response.Success(c, 200, "获取列表成功", data)
}

//杀死任务
func (j *JobController) KillJob(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", "")
		return
	}
	jobService := services.NewJobServiceHandler()
	if ok := jobService.KillJob(c, jobId); !ok {
		response.Error(c, http.StatusBadRequest, 400, "强杀任务失败", "")
		return
	}
	response.Success(c, 200, "强杀成功", "")
}

//更新任务
func (j *JobController) UpdateJob(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", "")
		return
	}
	req := &forms.JobUpdateRequest{}
	if err := c.ShouldBind(req); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	job := &models.Job{
		Name:       req.Name,
		Command:    req.Command,
		CronExpire: req.CronExpire,
	}
	jobService := services.NewJobServiceHandler()
	if err := jobService.UpdateJob(c, jobId, job); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), "")
		return
	}
	response.Success(c, 200, "修改成功", gin.H{
		"job_id": jobId,
	})
}

//任务详情
func (j *JobController) Show(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", "")
		return
	}
	jobService := services.NewJobServiceHandler()
	job, err := jobService.Show(c, jobId)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "资源不存在", "")
		return
	}
	response.Success(c, 200, "查询成功", job)
}

//任务日志列表
func (j *JobController) Logs(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", "")
		return
	}
	req := &forms.JobListRequest{}
	if err := c.ShouldBind(req); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	logService := services.NewJobLogService()
	total, result, err := logService.GetLog(c, jobId, req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "获取日志列表错误", err.Error())
		return
	}
	response.Success(c, 200, "获取日志列表成功", gin.H{
		"total":     total,
		"list":      result,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

//获取work节点
func (j *JobController) WorkNode(c *gin.Context) {
	jobService := services.NewJobServiceHandler()
	data, err := jobService.WorkNode(c)
	if err != nil {
		utils.LoggerFor(c).Info("get work node failed", zap.Error(err))
		response.Error(c, http.StatusBadRequest, 400, "获取节点列表错误", err.Error())
		return
	}
	response.Success(c, 200, "获取成功", data)
}
