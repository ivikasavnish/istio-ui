package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivikasavnish/istio-ui/internal/scheduler"
)

// ListSchedules returns all scheduled actions
func ListSchedules(sched *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actions := sched.ListActions()
		c.JSON(http.StatusOK, actions)
	}
}

// CreateSchedule creates a new scheduled action
func CreateSchedule(sched *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var action scheduler.ScheduledAction
		if err := c.ShouldBindJSON(&action); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate cron expression
		if err := scheduler.ValidateCronExpression(action.Schedule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cron expression: " + err.Error()})
			return
		}

		if err := sched.AddAction(&action); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, action)
	}
}

// DeleteSchedule deletes a scheduled action
func DeleteSchedule(sched *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := sched.RemoveAction(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
	}
}
