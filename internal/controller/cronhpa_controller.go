/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/robfig/cron"
	jumpserverv1 "github.com/wuyuetianjian/jumpcronsynchpa/api/v1"
)

// CronHPAReconciler reconciles a CronHPA object
type CronHPAReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=jumpserver.sunny.io,resources=cronhpas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jumpserver.sunny.io,resources=cronhpas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=jumpserver.sunny.io,resources=cronhpas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CronHPA object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *CronHPAReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Reconciling CronHPA")
	var cronhpa jumpserverv1.CronHPA
	if err := r.Get(ctx, req.NamespacedName, &cronhpa); err != nil {
		if errors.IsNotFound(err) {
			log.Info("CronHPA resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	now := time.Now()
	var earliestNextRunTime *time.Time

	// 遍历用户jobs，检查调度时间并更新内容
	for _, job := range cronhpa.Spec.SyncUserGroup {
		lastRunTime := cronhpa.Status.LastRunTimes[job.Name]
		// 计算上次运行时间之后的下一个调度时间
		nextScheduledTime, err := r.getNextScheduledTime(job.Schedule, lastRunTime.Time)
		if err != nil {
			log.Error(err, "Failed to calculate next scheduled time")
			return reconcile.Result{}, err
		}

		log.Info("Job info", "name", job.Name, "lastRunTime", lastRunTime, "nextScheduledTime", nextScheduledTime, "now", now)

		// 检查当前时间是否已经到达或超过了计划的运行时间
		if now.After(nextScheduledTime) || now.Equal(nextScheduledTime) {
			// 更新内容

		} else {
			// 如果当前时间未到达计划时间，将这个时间作为下一次运行时间
			if earliestNextRunTime == nil || nextScheduledTime.Before(*earliestNextRunTime) {
				earliestNextRunTime = &nextScheduledTime
			}
		}
	}

	// 遍历主机jobs，检查调度时间并更新内容
	for _, job := range cronhpa.Spec.SyncAssetsGroup {
		lastRunTime := cronhpa.Status.LastRunTimes[job.Name]
		// 计算上次运行时间之后的下一个调度时间
		nextScheduledTime, err := r.getNextScheduledTime(job.Schedule, lastRunTime.Time)
		if err != nil {
			log.Error(err, "Failed to calculate next scheduled time")
			return reconcile.Result{}, err
		}

		log.Info("Job info", "name", job.Name, "lastRunTime", lastRunTime, "nextScheduledTime", nextScheduledTime, "now", now)

		// 检查当前时间是否已经到达或超过了计划的运行时间
		if now.After(nextScheduledTime) || now.Equal(nextScheduledTime) {
			// 更新内容

		} else {
			// 如果当前时间未到达计划时间，将这个时间作为下一次运行时间
			if earliestNextRunTime == nil || nextScheduledTime.Before(*earliestNextRunTime) {
				earliestNextRunTime = &nextScheduledTime
			}
		}
	}

	// 更新 CronHPA 实例状态
	if err := r.Status().Update(ctx, &cronhpa); err != nil {
		return reconcile.Result{}, err
	}

	// 如果有下一次运行时间，设置重新入队
	if earliestNextRunTime != nil {
		requeueAfter := earliestNextRunTime.Sub(now)
		if requeueAfter < 0 {
			requeueAfter = time.Second // 如果计算出的时间已经过去，则在1秒后重新入队
		}
		log.Info("Requeue after", "time", requeueAfter)
		return reconcile.Result{RequeueAfter: requeueAfter}, nil
	}

	return ctrl.Result{}, nil
}

func (r *CronHPAReconciler) getNextScheduledTime(schedule string, after time.Time) (time.Time, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cronSchedule, err := parser.Parse(schedule)
	if err != nil {
		return time.Time{}, err
	}

	return cronSchedule.Next(after), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CronHPAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jumpserverv1.CronHPA{}).
		Named("cronhpa").
		Complete(r)
}
