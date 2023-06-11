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
	"bytes"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"text/template"
	vortalbizv1 "vortal.biz/joaoneves/azdevops-operator/api/v1"
)

func (r *AzDevopsAgentPoolReconciler) autoscale(request reconcile.Request,
	instance *vortalbizv1.AzDevopsAgentPool,
	sts *appsv1.StatefulSet,
	ctx context.Context) (int32, error) {

	log := log.FromContext(ctx).WithValues("AzDevopsControllerAutoScaler", sts.Name)

	currentReplicas := *sts.Spec.Replicas
	baseUrl := instance.Spec.Project.Url
	//projectName := instance.Spec.Project.ProjectName
	//poolName := instance.Spec.Project.PoolName
	//PATSecretRef := instance.Spec.Project.PatSecretRef
	apiVersion := "api-version=7.0"

	//Retrieve the PAT from the secret
	secret := new(corev1.Secret)
	if err := r.Get(ctx, types.NamespacedName{Namespace: instance.Namespace, Name: instance.Spec.Project.PatSecretRef}, secret); err != nil {
		log.Error(err, "Failed to find PAT secret.",
			"namespace", instance.Namespace,
			"name", instance.Spec.Project.PatSecretRef)
		return currentReplicas, err
	}
	PAT := string(secret.Data["token"])
	log.Info(PAT)

	var buf bytes.Buffer
	url := "{{ .baseUrl }}/_apis/distributedtask/pools?{{ .apiVersion }}"
	templ := template.Must(template.New("getPools").Parse(url))
	templ.Execute(&buf, map[string]interface{}{
		"baseUrl":    baseUrl,
		"apiVersion": apiVersion,
	})
	resp, err := http.Get(string(buf.Bytes()))
	if err != nil {
		log.Error(err, "Failed to execute request", string(buf.Bytes()))
		return currentReplicas, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Info("Unexpected response code", strconv.Itoa(resp.StatusCode), http.StatusText(resp.StatusCode))
		return currentReplicas, nil
	}

	return currentReplicas, nil
}
