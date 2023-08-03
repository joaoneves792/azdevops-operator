package controllers

//Get id from pool name:
//curl -u "joao.neves@vortal.biz:$AZDEVOPSTOKEN" -H "Content-Type: application/json" "https://dev.azure.com/vortal-projects/_apis/distributedtask/pools?api-version=7.0" | jq '.value[] | select(.name == "vision-ci-dotnet6") | .id'

//Get queued and running jobs
//jobRequests=$(curl -u peterjgrainger:${{ YOUR_DEVOPS_TOKEN }} https://dev.azure.com/{your_org}/_apis/distributedtask/pools/{your_pool}/jobrequests?api-version=6.0)
//queuedJobs=$(echo $jobRequests | jq '.value | map(select(has("assignTime") | not)) | length')
//runningJobs=$(echo $jobRequests | jq '.value | map(select(.result == null)) | length')

//List agents
//https://dev.azure.com/vortal-projects/_apis/distributedtask/pools/90/agents?api-version=7.0

//Disable agent
//curl -X PATCH -u "joao.neves@vortal.biz:$AZDEVOPSTOKEN" -H "Content-Type: application/json" "https://dev.azure.com/vortal-projects/_apis/distributedtask/pools/90/agents/269?api-version=7.0" -d '{"id": 269, "enabled":true}'

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"strings"
	"time"
	vortalbizv1 "vortal.biz/joaoneves/azdevops-operator/api/v1"
)

func (r *AzDevopsAgentPoolReconciler) performDevopsRESTRequest(method string, url string, body string) ([]byte, error) {
	log := r.log
	client := http.Client{}

	bodyReader := strings.NewReader(body)
	req, err := http.NewRequest(method, url, bodyReader) //this body needs to be io.Reader not string
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", r.PAT)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Failed to execute request", url)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Info("Unexpected response code", strconv.Itoa(resp.StatusCode), http.StatusText(resp.StatusCode))
		return nil, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, "Failed to read response body")
		return nil, err
	}

	return []byte(bodyBytes), nil
}

type TaskAgentPoolResponse struct {
	Value []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"value"`
}

type TaskAgentsResponse struct {
	Value []struct {
		Name    string `json:"name"`
		Id      int    `json:"id"`
		Enabled bool   `json:"enabled"`
		Status  string `json:"status"`
	} `json:"value"`
}

func (r *AzDevopsAgentPoolReconciler) getPoolID(instance *vortalbizv1.AzDevopsAgentPool) (int, error) {
	log := r.log
	baseUrl := instance.Spec.Project.Url
	poolName := instance.Spec.Project.PoolName
	apiVersion := "api-version=7.0"

	//Get all pools
	url := fmt.Sprintf("%s/_apis/distributedtask/pools?%s", baseUrl, apiVersion)

	response, err := r.performDevopsRESTRequest("GET", url, "")
	if err != nil {
		log.Error(err, "Failed to execute request")
		return 0, err
	}

	var robj TaskAgentPoolResponse

	err = json.Unmarshal(response, &robj)
	if err != nil {
		log.Error(err, "Failed to unmarshal Pool json")
		return 0, err
	}

	for _, pool := range robj.Value {
		if pool.Name == poolName {
			log.Info("PoolID", pool.Name, pool.Id)
			return pool.Id, nil
		}
	}

	return 0, nil
}

func (r *AzDevopsAgentPoolReconciler) getAgentsStatus(instance *vortalbizv1.AzDevopsAgentPool, poolId int) (TaskAgentsResponse, error) {
	log := r.log

	baseUrl := instance.Spec.Project.Url
	apiVersion := "api-version=7.0"

	var robj TaskAgentsResponse

	url := fmt.Sprintf("%s/_apis/distributedtask/pools/%s/agents?%s", baseUrl, strconv.Itoa(poolId), apiVersion)
	response, err := r.performDevopsRESTRequest("GET", url, "")
	if err != nil {
		log.Error(err, "Failed to execute request")
		return robj, err
	}

	err = json.Unmarshal(response, &robj)
	if err != nil {
		log.Error(err, "Failed to unmarshal Agents json")
		return robj, err
	}

	return robj, nil

}

func (r *AzDevopsAgentPoolReconciler) disableAgent(instance *vortalbizv1.AzDevopsAgentPool, poolId int, agentId int) error {
	log := r.log

	baseUrl := instance.Spec.Project.Url
	apiVersion := "api-version=7.0"

	url := fmt.Sprintf("%s/_apis/distributedtask/pools/%s/agents/%s?%s", baseUrl, strconv.Itoa(poolId), strconv.Itoa(agentId), apiVersion)
	body := fmt.Sprintf("{\"id\": %s, \"enabled\":false}", strconv.Itoa(agentId))
	_, err := r.performDevopsRESTRequest("PATCH", url, body)
	if err != nil {
		log.Error(err, "Failed to execute request")
		return err
	}
	return nil

}

func (r *AzDevopsAgentPoolReconciler) enableAgent(instance *vortalbizv1.AzDevopsAgentPool, poolId int, agentId int) error {
	log := r.log

	baseUrl := instance.Spec.Project.Url
	apiVersion := "api-version=7.0"

	url := fmt.Sprintf("%s/_apis/distributedtask/pools/%s/agents/%s?%s", baseUrl, strconv.Itoa(poolId), strconv.Itoa(agentId), apiVersion)
	body := fmt.Sprintf("{\"id\": %s, \"enabled\":true}", strconv.Itoa(agentId))
	_, err := r.performDevopsRESTRequest("PATCH", url, body)
	if err != nil {
		log.Error(err, "Failed to execute request")
		return err
	}
	return nil

}

func (r *AzDevopsAgentPoolReconciler) isAgentIdle(instance *vortalbizv1.AzDevopsAgentPool, poolId int, agentId int) (bool, error) {
	log := r.log

	baseUrl := instance.Spec.Project.Url
	apiVersion := "api-version=7.0"
	url := fmt.Sprintf("%s/_apis/distributedtask/pools/%s/agents/%s?includeAssignedRequest=true&%s", baseUrl, strconv.Itoa(poolId), strconv.Itoa(agentId), apiVersion)

	response, err := r.performDevopsRESTRequest("GET", url, "")
	if err != nil {
		log.Error(err, "Failed to execute request")
		return false, err
	}

	var obj map[string]json.RawMessage
	err = json.Unmarshal(response, &obj)
	if err != nil {
		log.Error(err, "Failed to unmarshal Agents json")
		return false, err
	}
	_, exists := obj["assignedRequest"]
	if exists {
		return false, nil //busy
	} else {
		return true, nil //idle
	}

}

func (r *AzDevopsAgentPoolReconciler) calculateScheduleReplicas(instance *vortalbizv1.AzDevopsAgentPool) (int32, error) {
	log := r.log

	tz := instance.Spec.Autoscaling.Schedule.TZ
	layout := "15:04"
	scaleUp := instance.Spec.Autoscaling.Schedule.ScaleUp
	scaleDown := instance.Spec.Autoscaling.Schedule.ScaleDown

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Error(err, "Failed to load timeZone")
		return 0, err
	}

	n := time.Now().In(loc)

	t, err := time.Parse(layout, scaleUp)
	if err != nil {
		log.Error(err, "Failed to parse scaleUp time")
		return 0, err
	}
	scaleUpTime := time.Date(n.Year(), n.Month(), n.Day(), t.Hour(), t.Minute(), 0, n.Nanosecond(), n.Location())

	t, _ = time.Parse(layout, scaleDown)
	if err != nil {
		log.Error(err, "Failed to parse scaleDown time")
		return 0, err
	}
	scaleDownTime := time.Date(n.Year(), n.Month(), n.Day(), t.Hour(), t.Minute(), 0, n.Nanosecond(), n.Location())

	skipWeekends := instance.Spec.Autoscaling.Schedule.SkipWeekends
	if skipWeekends && (n.Weekday() == time.Saturday || n.Weekday() == time.Sunday) {
		return instance.Spec.Autoscaling.Min, nil
	}

	if n.After(scaleUpTime) && n.Before(scaleDownTime) {
		return instance.Spec.Autoscaling.Max, nil
	} else {
		return instance.Spec.Autoscaling.Min, nil
	}
}

func (r *AzDevopsAgentPoolReconciler) autoscale(request reconcile.Request,
	instance *vortalbizv1.AzDevopsAgentPool,
	sts *appsv1.StatefulSet,
	ctx context.Context) (int32, error) {

	log := r.log

	currentReplicas := *sts.Spec.Replicas
	desiredReplicas, err := r.calculateScheduleReplicas(instance)
	if err != nil {
		log.Error(err, "Schedule error, defaulting to current replicas")
		desiredReplicas = currentReplicas
	}

	//Retrieve the PAT from the secret
	secret := new(corev1.Secret)
	if err := r.Get(ctx, types.NamespacedName{Namespace: instance.Namespace, Name: instance.Spec.Project.PatSecretRef}, secret); err != nil {
		log.Error(err, "Failed to find PAT secret.",
			"namespace", instance.Namespace,
			"name", instance.Spec.Project.PatSecretRef)
		return currentReplicas, err
	}
	PAT := string(secret.Data["token"])
	r.PAT = PAT

	//Retrieve Pool ID from Pool name
	poolID, err := r.getPoolID(instance)
	if err != nil {
		log.Error(err, "Failed to resolve pool name to a pool ID")
		return currentReplicas, err
	}
	if poolID == 0 {
		log.Info("Failed to find a matching pool, check the pool name")
		return currentReplicas, nil
	}

	//Fetch agents
	agentStatus, err := r.getAgentsStatus(instance, poolID)
	if err != nil {
		log.Error(err, "Failed to fetch agents status")
		return currentReplicas, err
	}

	for _, agent := range agentStatus.Value {
		log.Info("Status", agent.Name, agent.Enabled)

		//Retrieve agent number from name
		splitName := strings.Split(agent.Name, instance.Name)
		if splitName[0] == agent.Name {
			continue //Not an agent controlled by us
		}
		i, err := strconv.Atoi(strings.Trim(splitName[1], "-"))
		if err != nil {
			log.Error(err, "Failed to retrieve agent number from agent name")
			return currentReplicas, err
		}

		//Ensure only the desired agents are enabled
		if int32(i) > (desiredReplicas - 1) {
			if agent.Enabled {
				log.Info("Disable agent", agent.Name, agent.Enabled)
				err = r.disableAgent(instance, poolID, agent.Id)
				if err != nil {
					log.Error(err, "Failed to disable agent")
					continue
				}
			}
		} else if int32(i) <= (desiredReplicas - 1) {
			if !agent.Enabled {
				log.Info("Enable agent", agent.Name, agent.Enabled)
				err = r.enableAgent(instance, poolID, agent.Id)
				if err != nil {
					log.Error(err, "Failed to enable agent")
					continue
				}
				agent.Enabled = true
			}
		}

		//Try to shut down idle disabled agents, but only the highest numbered one
		if int32(i) == (currentReplicas-1) && !agent.Enabled {
			idle, err := r.isAgentIdle(instance, poolID, agent.Id)
			if err != nil {
				log.Error(err, "Failed to check if agent was idle", agent.Name)
			}
			if idle {
				log.Info("Shutting down idle disabled agent", agent.Name, agent.Enabled)
				return currentReplicas - 1, nil
			}
		}
	}

	return desiredReplicas, nil
}
