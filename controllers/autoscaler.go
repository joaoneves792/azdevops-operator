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
	"github.com/go-logr/logr"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"strings"
	vortalbizv1 "vortal.biz/joaoneves/azdevops-operator/api/v1"
)

func (r *AzDevopsAgentPoolReconciler) performDevopsRESTRequest(method string, url string, PAT string, body string, log *logr.Logger) ([]byte, error) {
	client := http.Client{}

	bodyReader := strings.NewReader(body)
	req, err := http.NewRequest(method, url, bodyReader) //this body needs to be io.Reader not string
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", PAT)

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

func (r *AzDevopsAgentPoolReconciler) getPoolID(instance *vortalbizv1.AzDevopsAgentPool, PAT string, log *logr.Logger) (int, error) {
	baseUrl := instance.Spec.Project.Url
	poolName := instance.Spec.Project.PoolName
	apiVersion := "api-version=7.0"

	//Get all pools
	url := fmt.Sprintf("%s/_apis/distributedtask/pools?%s", baseUrl, apiVersion)

	response, err := r.performDevopsRESTRequest("GET", url, PAT, "", log)
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

func (r *AzDevopsAgentPoolReconciler) autoscale(request reconcile.Request,
	instance *vortalbizv1.AzDevopsAgentPool,
	sts *appsv1.StatefulSet,
	ctx context.Context) (int32, error) {

	log := log.FromContext(ctx).WithValues("AzDevopsControllerAutoScaler", sts.Name)

	currentReplicas := *sts.Spec.Replicas
	//projectName := instance.Spec.Project.ProjectName

	//Retrieve the PAT from the secret
	secret := new(corev1.Secret)
	if err := r.Get(ctx, types.NamespacedName{Namespace: instance.Namespace, Name: instance.Spec.Project.PatSecretRef}, secret); err != nil {
		log.Error(err, "Failed to find PAT secret.",
			"namespace", instance.Namespace,
			"name", instance.Spec.Project.PatSecretRef)
		return currentReplicas, err
	}
	PAT := string(secret.Data["token"])

	poolID, err := r.getPoolID(instance, PAT, &log)
	if err != nil {
		log.Error(err, "Failed to resolve pool name to a pool ID")
		return currentReplicas, err
	}

	//log.Info("Got ", "pools", strconv.Itoa(int(obj["count"].(float64))))
	log.Info("Success ", "PoolID", strconv.Itoa(poolID))

	return currentReplicas, nil
}
