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
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// ControllerRuntimeInterface is an interface to support operations for a runtime client
type ControllerRuntimeInterface interface {
	CreateSecret(ctx context.Context, bytes []byte) error
	CreateNameSpace(ctx context.Context, data []byte) error
}

// ControllerRuntimeClient provides functionality for the controller runtime
type ControllerRuntimeClient struct {
	client ctrlClient.Client
	Logger echo.Logger
}

// Client provides functionality for accessing the client go library
type Client struct{}

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
	if ok {
		namespaces := &corev1.NamespaceList{}
		err := c.client.List(ctx, namespaces)
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
	}
	return nil
}

// CreateSecret will create the given Secret resource in a k8s cluster
func (c ControllerRuntimeClient) CreateSecret(ctx context.Context, data []byte) error {
	fmt.Println(data)
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

// NewControllerRuntimeClient will return a new controller runtime client for the given kubeconfig
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

// GetControllerClient - Returns a controller client which reads and writes directly to API server
func GetControllerClient(restConfig *rest.Config, scheme *runtime.Scheme) (ctrlClient.Client, error) {
	// Create a temp client and use it
	var clientConfig *rest.Config
	var err error
	if restConfig == nil {
		clientConfig, err = config.GetConfig()
		if err != nil {
			return nil, err
		}
	} else {
		clientConfig = restConfig
	}
	client, err := ctrlClient.New(clientConfig, ctrlClient.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// IsOpenShift returns true if the cluster is OpenShift, otherwise returns false
func (c Client) IsOpenShift(data []byte) (bool, error) {
	k8sClientSet, err := c.GetClientSet(data)
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

// GetVersion returns version of the k8s cluster
func (c Client) GetVersion(data []byte) (string, error) {
	k8sClientSet, err := c.GetClientSet(data)
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

	k8sClientSet, err := c.GetClientSet(data)
	if err != nil {
		return "", nil, nil, err
	}

	return version, &isOpenShift, k8sClientSet, nil
}

// GetClientSet returns a reference to the Clientset for the given cluster
func (c Client) GetClientSet(data []byte) (*kubernetes.Clientset, error) {
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(data)
	if err != nil {
		return nil, err
	}
	return newClientSet(restConfig)
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
	clusterID   uint
	nPods       int
	readyPods   chan corev1.Pod
	pendingPods chan corev1.Pod
	failedPods  chan corev1.Pod
	nodes       chan Node
	wg          sync.WaitGroup
	ClientSet   *kubernetes.Clientset
	Logger      echo.Logger
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
		_, err = collector.ClientSet.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
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
	time.Sleep(30 * time.Second)
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
	pods, err := retryableGetPods(collector.ClientSet)
	if err != nil {
		fmt.Printf("Failed to query for the daemonset pods. Error: %s\n", err.Error())
		return nil, err
	}
	collector.nPods = len(pods.Items)

	for _, pod := range pods.Items {
		collector.calculatePodStatus(pod)
	}
	collector.wg.Add(1)
	go collector.handleTerminatedPod(5 * time.Minute)
	collector.wg.Add(1)
	go collector.handleFailedPods()
	collector.wg.Add(1)
	go collector.handlePendingPods(3 * time.Minute)
	collector.wg.Wait()
	close(collector.nodes)
	nodes := make([]Node, 0)
	for node := range collector.nodes {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (collector *NodeDataCollector) handleTerminatedPod(duration time.Duration) {
	for afterCh := time.After(duration); ; {
		select {
		case pod := <-collector.readyPods:
			logs := getPodLogs(collector.ClientSet, pod.Name, pod.Namespace)

			installedSoftware := make(map[string]string)
			scanner := bufio.NewScanner(strings.NewReader(logs))
			for scanner.Scan() {
				// expect enabled software to be of the form "<software-name>=<value>"
				software := strings.Split(scanner.Text(), "=")
				if len(software) == 2 {
					installedSoftware[software[0]] = "enabled"
				}
			}
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

func (collector *NodeDataCollector) handlePendingPods(duration time.Duration) {
	for afterCh := time.After(duration); ; {
		minPollCh := time.After(5 * time.Second)
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
			time.Sleep(3 * time.Second)
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

func retryableGetPods(clientset *kubernetes.Clientset) (*corev1.PodList, error) {
	opts := metav1.ListOptions{
		LabelSelector: "app=csm-data-collector",
	}
	deadline := time.Now().Add(time.Duration(30 * time.Second))
	var pods *corev1.PodList
	var err error
	for tries := 0; time.Now().Before(deadline); tries++ {
		pods, err = clientset.CoreV1().Pods("csm").List(context.TODO(), opts)
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

func refreshPodInfo(clientset *kubernetes.Clientset, pod corev1.Pod) (corev1.Pod, error) {
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

func getPodLogs(clientset *kubernetes.Clientset, podName, podNamespace string) string {
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
					},
				},
			},
		},
	}
	return daemonSet
}
