// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package k8s

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// ControllerRuntimeInterface is an interface to support operations for a runtime client
//go:generate mockgen -destination=mocks/controller_runtime_interface.go -package=mocks github.com/dell/csm-deployment/k8s ControllerRuntimeInterface
type ControllerRuntimeInterface interface {
	CreateSecret(ctx context.Context, bytes []byte) error
	CreateNameSpace(ctx context.Context, data []byte) error
	CreateNameSpaceFromName(ctx context.Context, name string) error
	DeleteNameSpaceByName(ctx context.Context, name string) error
	CreateSecretFromName(ctx context.Context, name string, namespace string, secretData []byte) error
	CreateTLSSecretFromName(ctx context.Context, name string, namespace, key, cert []byte) error
	CreateConfigMap(ctx context.Context, bytes []byte) error
}

// ControllerRuntimeClient provides functionality for the controller runtime
type ControllerRuntimeClient struct {
	client ctrlClient.Client
	Logger echo.Logger
}

// Client provides functionality for accessing the client go library
type Client struct{}

// CreateNameSpaceFromName will create the Namespace in a k8s cluster using the name
func (c ControllerRuntimeClient) CreateNameSpaceFromName(ctx context.Context, name string) error {
	// check if it is a supported type
	namespaceObj := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}

	namespaces := &corev1.NamespaceList{}
	err := c.client.List(ctx, namespaces)
	if err != nil {
		return err
	}
	for _, namespace := range namespaces.Items {
		if namespace.Name == namespaceObj.Name {
			c.Logger.Info("Namespace already exists. Returning with success")
			return nil
		}
	}

	err = c.client.Create(ctx, &namespaceObj)
	if err != nil {
		return err
	}
	fmt.Println("Successfully created namespace: ", namespaceObj.Name)

	return nil
}

// DeleteNameSpaceByName will delete the Namespace in a k8s cluster using the name
func (c ControllerRuntimeClient) DeleteNameSpaceByName(ctx context.Context, name string) error {
	// check if it is a supported type
	namespaceObj := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}

	namespaces := &corev1.NamespaceList{}
	err := c.client.List(ctx, namespaces)
	if err != nil {
		return err
	}
	for _, namespace := range namespaces.Items {
		if namespace.Name == namespaceObj.Name {
			c.Logger.Info("Namespace exists, Proceeding to delete")
			err = c.client.Delete(ctx, &namespaceObj)
			if err != nil {
				return err
			}
			fmt.Println("Successfully deleted namespace: ", namespaceObj.Name)
		}
	}
	return nil
}

// CreateSecretFromName - Creates secret using the name and payload data
func (c ControllerRuntimeClient) CreateSecretFromName(ctx context.Context, name string, namespace string, secretData []byte) error {
	// check if it is a supported type
	m := make(map[string]string)
	m["data"] = string(secretData)

	secretObj := corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		StringData: m,
	}

	secrets := &corev1.SecretList{}
	err := c.client.List(ctx, secrets, &ctrlClient.ListOptions{Namespace: namespace})
	if err != nil {
		return err
	}

	for _, secret := range secrets.Items {
		if secret.Name == name {
			c.Logger.Info("Secret already exists. Returning with success")
			return nil
		}
	}

	err = c.client.Create(ctx, &secretObj)
	if err != nil {
		return err
	}
	fmt.Println("Successfully created secret: ", secretObj.Name)

	return nil
}

// CreateTLSSecretFromName - Creates tls secret using the name, key and cert
func (c ControllerRuntimeClient) CreateTLSSecretFromName(ctx context.Context, name string, namespaceData []byte, key, cert []byte) error {

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	runtimeObj, _, err := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode(namespaceData, nil, nil)
	if err != nil {
		return err
	}
	// check if it is a supported type
	namespaceObj, ok := runtimeObj.(*corev1.Namespace)
	if !ok {
		return fmt.Errorf("unsupported object type")
	}

	m := make(map[string][]byte)
	m["tls.key"] = key
	m["tls.crt"] = cert

	secretObj := corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespaceObj.Name,
			Name:      name,
		},
		Type: corev1.SecretTypeTLS,
		Data: m,
	}

	secrets := &corev1.SecretList{}
	err = c.client.List(ctx, secrets, &ctrlClient.ListOptions{Namespace: namespaceObj.Name})
	if err != nil {
		return err
	}

	for _, secret := range secrets.Items {
		if secret.Name == name {
			c.Logger.Info("Secret already exists. Returning with success")
			return nil
		}
	}

	err = c.client.Create(ctx, &secretObj)
	if err != nil {
		return err
	}
	fmt.Println("Successfully created secret: ", secretObj.Name)

	return nil
}

// CreateNameSpace will create the given Namespace in a k8s cluster
func (c ControllerRuntimeClient) CreateNameSpace(ctx context.Context, data []byte) error {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	runtimeObj, _, err := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode(data, nil, nil)
	if err != nil {
		return err
	}
	// check if it is a supported type
	namespaceObj, ok := runtimeObj.(*corev1.Namespace)
	if !ok {
		return fmt.Errorf("unsupported object type")
	}

	namespaces := &corev1.NamespaceList{}
	err = c.client.List(ctx, namespaces)
	if err != nil {
		// ignore this and continue to create namespace
	} else {
		for _, namespace := range namespaces.Items {
			if namespace.Name == namespaceObj.Name {
				c.Logger.Info("Namespace already exists. Returning with success")
				return nil
			}
		}
	}
	err = c.client.Create(ctx, namespaceObj)
	if err != nil {
		return err
	}
	fmt.Println("Successfully created namespace: ", namespaceObj.Name)

	return nil
}

// CreateSecret will create the given Secret resource in a k8s cluster
func (c ControllerRuntimeClient) CreateSecret(ctx context.Context, data []byte) error {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	runtimeObj, _, err := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode(data, nil, nil)
	if err != nil {
		return err
	}
	// check if it is a supported type
	secretObj, ok := runtimeObj.(*corev1.Secret)
	if !ok {
		return fmt.Errorf("unsupported object type")
	}
	if ok {
		err := c.client.Create(ctx, secretObj)
		if err != nil {
			return err
		}
		fmt.Println("Successfully created secret: ", secretObj.Name)
	}
	return nil
}

// CreateConfigMap will create the given config map resource in a k8s cluster
func (c ControllerRuntimeClient) CreateConfigMap(ctx context.Context, data []byte) error {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	runtimeObj, _, err := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode(data, nil, nil)
	if err != nil {
		return err
	}
	// check if it is a supported type
	configMapObj, ok := runtimeObj.(*corev1.ConfigMap)
	if !ok {
		return fmt.Errorf("unsupported object type")
	}
	if ok {
		err := c.client.Create(ctx, configMapObj)
		if err != nil {
			return err
		}
		fmt.Println("Successfully created configmap: ", configMapObj.Name)
	}
	return nil
}

// NewControllerRuntimeClient will return a new controller runtime client
func NewControllerRuntimeClient(data []byte, logger echo.Logger) (ControllerRuntimeInterface, error) {
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(data)
	if err != nil {
		return nil, err
	}
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	client, err := GetControllerClient(restConfig, scheme)
	if err != nil {
		return nil, err
	}

	return ControllerRuntimeClient{
		client: client,
		Logger: logger,
	}, nil
}

// CtrlClientNewWrapper -
var CtrlClientNewWrapper = func(clientConfig *rest.Config, scheme *runtime.Scheme) (ctrlClient.Client, error) {
	return ctrlClient.New(clientConfig, ctrlClient.Options{Scheme: scheme})
}

// GetConfigWrapper -
var GetConfigWrapper = func() (*rest.Config, error) {
	return config.GetConfig()
}

// GetControllerClient - Returns a controller client which reads and writes directly to API server
func GetControllerClient(restConfig *rest.Config, scheme *runtime.Scheme) (ctrlClient.Client, error) {
	// Create a temp client and use it
	var clientConfig *rest.Config
	var err error
	if restConfig == nil {
		clientConfig, err = GetConfigWrapper()
		if err != nil {
			return nil, err
		}
	} else {
		clientConfig = restConfig
	}
	return CtrlClientNewWrapper(clientConfig, scheme)
}

// IsOpenShift returns true if the cluster is OpenShift, otherwise returns false
func (c Client) IsOpenShift(data []byte) (bool, error) {
	k8sClientSet, err := GetClientSetWrapper(c, data)
	if err != nil {
		return false, err
	}

	serverGroups, _, err := k8sClientSet.Discovery().ServerGroupsAndResources()
	if err != nil {
		return false, err
	}
	openshiftAPIGroup := "security.openshift.io"
	for i := 0; i < len(serverGroups); i++ {
		if serverGroups[i].Name == openshiftAPIGroup {
			return true, nil
		}
	}
	return false, nil
}

// ServerPreferredResourcesWrapper -
var ServerPreferredResourcesWrapper = func(k8sClientSet kubernetes.Interface) ([]*metav1.APIResourceList, error) {
	return k8sClientSet.Discovery().ServerPreferredResources()
}

// GetAPIResource will return details about a specific CRD resource
func (c Client) GetAPIResource(data []byte, resourceName string) (*metav1.APIResource, string, error) {
	k8sClientSet, err := GetClientSetWrapper(c, data)
	if err != nil {
		return nil, "", err
	}
	list, err := ServerPreferredResourcesWrapper(k8sClientSet)
	if err != nil {
		return nil, "", err
	}
	for _, resourceList := range list {
		for _, resource := range resourceList.APIResources {
			if resource.Kind == resourceName {
				return &resource, resourceList.GroupVersion, nil
			}
		}
	}
	return nil, "", fmt.Errorf("resource %s not found", resourceName)
}

// GetVersion returns version of the k8s cluster
func (c Client) GetVersion(data []byte) (string, error) {
	k8sClientSet, err := GetClientSetWrapper(c, data)
	if err != nil {
		return "", err
	}
	sv, err := k8sClientSet.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", sv.Major, sv.Minor), nil
}

// DiscoverK8sDetails - discover k8s details
func (c Client) DiscoverK8sDetails(data []byte) (string, *bool, *kubernetes.Clientset, error) {

	version, err := c.GetVersion(data)
	if err != nil {
		return "", nil, nil, err
	}

	isOpenShift, err := c.IsOpenShift(data)
	if err != nil {
		return "", nil, nil, err
	}

	k8sClientSet, err := GetClientSetWrapper(c, data)
	if err != nil {
		return "", nil, nil, err
	}

	if _, ok := k8sClientSet.(*fake.Clientset); ok {
		k8sClientSet = &kubernetes.Clientset{}
	}

	return version, &isOpenShift, k8sClientSet.(*kubernetes.Clientset), nil
}

// GetClientSet returns a reference to the Clientset for the given cluster
func (c Client) GetClientSet(data []byte) (*kubernetes.Clientset, error) {
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(data)
	if err != nil {
		return nil, err
	}
	return newClientSet(restConfig)
}

// GetClientSetWrapper -
var GetClientSetWrapper = func(c Client, data []byte) (kubernetes.Interface, error) {
	return c.GetClientSet(data)
}

// GetCertManagerPods returns the pods in cert-manager namespace
func (c Client) GetCertManagerPods(data []byte, namespace, label string) (*corev1.PodList, error) {
	k8sClientSet, err := GetClientSetWrapper(c, data)
	if err != nil {
		return nil, err
	}
	return retryableGetPods(k8sClientSet, namespace, metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", label)})
}

func newClientSet(restConfig *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// Node provides details about a node in a k8s cluster
type Node struct {
	HostName          string            `json:"host_name"`
	InstalledSoftware map[string]string `json:"installed_software"`
}

// NodeDataCollector can gather logs from an installed daemon set on the cluster
type NodeDataCollector struct {
	clusterID                  uint
	nPods                      int
	readyPods                  chan corev1.Pod
	pendingPods                chan corev1.Pod
	failedPods                 chan corev1.Pod
	nodes                      chan Node
	wg                         sync.WaitGroup
	ClientSet                  kubernetes.Interface
	Logger                     echo.Logger
	InstallWaitTime            time.Duration
	HandleTerminatedPodTimeout time.Duration
	HandlePendingPodTimeout    time.Duration
	HandlePendingPodsWaitTime  time.Duration
}

func (collector *NodeDataCollector) init() {
	collector.readyPods = make(chan corev1.Pod, 100)
	collector.pendingPods = make(chan corev1.Pod, 100)
	collector.failedPods = make(chan corev1.Pod, 100)
	collector.nodes = make(chan Node, 100)
}

// Install will install the data collector daemonset on a k8s cluster
func (collector *NodeDataCollector) Install() error {
	// Read the image name from env
	img := utils.GetEnv("DATA_COLLECTOR_IMAGE", "")
	if img == "" {
		return fmt.Errorf("invalid configuration. Data collector image not set")
	}
	collector.Logger.Info("Querying if the data collector daemonset is already installed in the cluster")
	// First create the namespace "csm" if required
	createNS := true
	_, err := collector.ClientSet.CoreV1().Namespaces().Get(context.TODO(), "csm", metav1.GetOptions{})
	if err != nil {
		// Lets ignore this (any error apart from NotFound) and still try to create the namespace
	} else {
		createNS = false
	}
	if createNS {
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "csm",
			},
		}
		_, err := collector.ClientSet.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
		if err != nil {
			collector.Logger.Error("Failed to create namespace: ", err.Error())
			return err
		}
	}
	_, err = collector.ClientSet.AppsV1().DaemonSets("csm").Get(
		context.TODO(), "csm-data-collector", metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		collector.Logger.Info("Installing the data collector daemonset")

		_, err := collector.ClientSet.AppsV1().DaemonSets("csm").Create(
			context.TODO(), getDaemonset(img), metav1.CreateOptions{})
		if err != nil {
			collector.Logger.Error("failed to install the data collector daemonset", err.Error())
			return err
		}
	} else if err != nil {
		collector.Logger.Error("failed to query if the data collector daemonset is already installed", err.Error())
		return err
	} else {
		collector.Logger.Info("Deleting the existing data collector daemonset")
		err = collector.ClientSet.AppsV1().DaemonSets("csm").Delete(
			context.TODO(), "csm-data-collector", metav1.DeleteOptions{})
		if err != nil {
			collector.Logger.Error("failed to delete the data collector daemonset", err.Error())
			return err
		}
		collector.Logger.Info("Installing the data collector daemonset")
		_, err := collector.ClientSet.AppsV1().DaemonSets("csm").Create(
			context.TODO(), getDaemonset(img), metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	collector.Logger.Info("Waiting for 30 seconds before querying for the pod status")
	time.Sleep(collector.InstallWaitTime)
	return nil
}

// Collect will gather output from the data-collector daemonset that is installed on each node in the cluster
func (collector *NodeDataCollector) Collect() ([]Node, error) {
	collector.init()
	// First install the daemonset
	err := collector.Install()
	if err != nil {
		fmt.Printf("Failed to install the data collector daemonset. Error: %s\n", err.Error())
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: "app=csm-data-collector",
	}

	pods, err := retryableGetPods(collector.ClientSet, "csm", opts)
	if err != nil {
		fmt.Printf("Failed to query for the daemonset pods. Error: %s\n", err.Error())
		return nil, err
	}
	collector.nPods = len(pods.Items)

	for _, pod := range pods.Items {
		collector.calculatePodStatus(pod)
	}
	collector.wg.Add(1)
	go collector.handleTerminatedPod(collector.HandleTerminatedPodTimeout)
	collector.wg.Add(1)
	go collector.handleFailedPods()
	collector.wg.Add(1)
	go collector.handlePendingPods(collector.HandlePendingPodTimeout)
	collector.wg.Wait()
	close(collector.nodes)

	//consolidate node data because collector may have returned more than 1 entry per node
	nodeSlice := make([]Node, 0)
	for node := range collector.nodes {
		nodeSlice = append(nodeSlice, node)
	}
	nodes := consolidateNodeDiscoveryDetails(nodeSlice)
	return nodes, nil
}

// consolidateNodeDiscoveryDetails will consolidate the slice of nodes to have one entry per
// hostname in case input contains multiple entries per node
func consolidateNodeDiscoveryDetails(nodes []Node) []Node {
	nodeInfo := make(map[string]Node)
	nodeResults := make([]Node, 0)
	for _, node := range nodes {
		if node.HostName == "" {
			continue
		}
		if _, ok := nodeInfo[node.HostName]; !ok {
			nodeInfo[node.HostName] = Node{
				HostName:          node.HostName,
				InstalledSoftware: make(map[string]string),
			}
		}
		for k, v := range node.InstalledSoftware {
			nodeInfo[node.HostName].InstalledSoftware[k] = v
		}
	}
	for _, node := range nodeInfo {
		nodeResults = append(nodeResults, node)
	}
	return nodeResults
}

// getInstalledSoftware will scan the logs and get list of installed software
func getInstalledSoftware(logs string) map[string]string {
	installedSoftware := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(logs))
	for scanner.Scan() {
		// expect enabled software to be of the form "<software-name>=<value>"
		software := strings.Split(scanner.Text(), "=")
		if len(software) == 2 {
			softwareName := strings.TrimSpace(software[0])
			softwareEnabled := strings.TrimSpace(software[1])
			if len(softwareName) != 0 && len(softwareEnabled) != 0 {
				installedSoftware[softwareName] = "enabled"
			}
		}
	}
	return installedSoftware
}

func (collector *NodeDataCollector) handleTerminatedPod(duration time.Duration) {
	for afterCh := time.After(duration); ; {
		select {
		case pod := <-collector.readyPods:
			logs := getPodLogs(collector.ClientSet, pod.Name, pod.Namespace)
			installedSoftware := getInstalledSoftware(logs)
			node := Node{
				HostName:          pod.Status.HostIP,
				InstalledSoftware: installedSoftware,
			}
			collector.nodes <- node
			if len(collector.nodes) == collector.nPods {
				collector.wg.Done()
				close(collector.readyPods)
				collector.Logger.Info("Exiting handleTerminatedPod as all node status has been updated")
				return
			}
		case <-afterCh:
			collector.wg.Done()
			close(collector.readyPods)
			return
		}
	}
}

// MinPollChTimeOut -
var MinPollChTimeOut = 5 * time.Second

func (collector *NodeDataCollector) handlePendingPods(duration time.Duration) {
	for afterCh := time.After(duration); ; {
		minPollCh := time.After(MinPollChTimeOut)
		stopWaiting := false
		select {
		case pod := <-collector.pendingPods:
			updatedPod, err := refreshPodInfo(collector.ClientSet, pod)
			if err != nil {
				// Put the pod back in the pending queue
				collector.pendingPods <- pod
			} else {
				collector.calculatePodStatus(updatedPod)
			}
			if len(collector.nodes) == collector.nPods {
				stopWaiting = true
				collector.Logger.Info("Exiting handlePendingPods as all node status has been updated")
				break
			}
			collector.Logger.Info("Sleeping for 3 seconds before querying for next pod status")
			time.Sleep(collector.HandlePendingPodsWaitTime)
		case <-minPollCh:
			if len(collector.nodes) == collector.nPods {
				stopWaiting = true
				collector.Logger.Info("Exiting handlePendingPods as all node status has been updated")
				break
			}
		case <-afterCh:
			stopWaiting = true
			break
		}
		if stopWaiting {
			break
		}
	}
	close(collector.pendingPods)
	for pod := range collector.pendingPods {
		// Update all these as unknown
		node := Node{
			HostName:          pod.Status.HostIP,
			InstalledSoftware: map[string]string{},
		}
		collector.nodes <- node
	}
	collector.wg.Done()
}

func (collector *NodeDataCollector) handleFailedPods() {
	close(collector.failedPods)
	for pod := range collector.failedPods {
		// Update all these as unknown
		node := Node{
			HostName:          pod.Status.HostIP,
			InstalledSoftware: map[string]string{},
		}
		collector.nodes <- node
	}
	collector.wg.Done()
}

// RetryableGetPodsTimeOut -
var RetryableGetPodsTimeOut = time.Duration(30 * time.Second)

func retryableGetPods(clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
	deadline := time.Now().Add(RetryableGetPodsTimeOut)
	fmt.Println(RetryableGetPodsTimeOut.Seconds())
	var pods *corev1.PodList
	var err error
	for tries := 0; time.Now().Before(deadline); tries++ {
		pods, err = clientset.CoreV1().Pods(namespace).List(context.TODO(), opts)
		if err != nil {
			fmt.Println("failed to get the list of pods")
			time.Sleep(time.Second << uint(tries)) // exponential backoff
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func refreshPodInfo(clientset kubernetes.Interface, pod corev1.Pod) (corev1.Pod, error) {
	updatedPod, err := clientset.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
	if err != nil {
		return corev1.Pod{}, err
	}
	return *updatedPod, nil
}

func (collector *NodeDataCollector) calculatePodStatus(pod corev1.Pod) {
	for _, initContainerStatus := range pod.Status.InitContainerStatuses {
		if initContainerStatus.Name == "csm-init" {
			if initContainerStatus.State.Terminated != nil {
				if initContainerStatus.State.Terminated.ExitCode == 0 {
					collector.readyPods <- pod
				} else {
					collector.failedPods <- pod
				}
			} else {
				if initContainerStatus.State.Waiting != nil || initContainerStatus.State.Running != nil {
					collector.pendingPods <- pod
				}
			}
			break
		}
	}
}

func getPodLogs(clientset kubernetes.Interface, podName, podNamespace string) string {
	podLogOpts := corev1.PodLogOptions{
		Container: "csm-init",
	}
	req := clientset.CoreV1().Pods(podNamespace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()
	return str
}

func getDaemonset(imageName string) *appsv1.DaemonSet {
	// TODO: This shouldn't be hardcoded
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csm-data-collector",
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "csm-data-collector",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "csm-data-collector",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "pause",
							Image: "gcr.io/google_containers/pause",
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:            "csm-init",
							Image:           imageName,
							ImagePullPolicy: corev1.PullAlways,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "etc",
									MountPath: "/hostetc",
									ReadOnly:  true,
								},
								{
									Name:      "proc",
									MountPath: "/hostproc",
									ReadOnly:  true,
								},
								{
									Name:      "usr",
									MountPath: "/hostusr",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "etc",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc",
								},
							},
						},
						{
							Name: "proc",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/proc",
								},
							},
						},
						{
							Name: "usr",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/usr",
								},
							},
						},
					},
				},
			},
		},
	}
	return daemonSet
}
