// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package k8s

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"

	"k8s.io/apimachinery/pkg/version"
	discoveryfake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlClientFake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type testOverrides struct {
	getEnv                          func(envName string, defaultValue string) string
	getClientSetWrapper             func(c Client, data []byte) (kubernetes.Interface, error)
	serverPreferredResourcesWrapper func(k8sClientSet kubernetes.Interface) ([]*metav1.APIResourceList, error)
	ctrlClientNewWrapper            func(clientConfig *rest.Config, scheme *runtime.Scheme) (ctrlClient.Client, error)
	getConfigWrapper                func() (*rest.Config, error)
	retryableGetPodsTimeOut         time.Duration
	minPollChTimeOut                time.Duration
}

var (
	namespace        = "csm"
	csmDataCollector = "csm-data-collector"
	secretName       = "testing-secret"
	configMapName    = "testing-configMap"
)

func Test_CreateNameSpaceFromName(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, *ControllerRuntimeClient){
		"success- namespace already exist": func(*testing.T) (bool, testOverrides, *ControllerRuntimeClient) {
			ns := &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        namespace,
					Annotations: map[string]string{},
				},
			}
			client := ctrlClientFake.NewClientBuilder().WithObjects(ns).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, testOverrides{}, controllerRuntimeClient
		},
		"success-  namespace does not exist": func(*testing.T) (bool, testOverrides, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, testOverrides{}, controllerRuntimeClient
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, controllerRuntimeClient := tc(t)
			if patch.getEnv != nil {
				oldGetEnv := utils.GetEnv
				defer func() { utils.GetEnv = oldGetEnv }()
				utils.GetEnv = patch.getEnv
			}

			err := controllerRuntimeClient.CreateNameSpaceFromName(context.TODO(), namespace)
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				namespaces := &corev1.NamespaceList{}
				err := controllerRuntimeClient.client.List(context.TODO(), namespaces)
				assert.NoError(t, err)
				found := false
				for _, ns := range namespaces.Items {
					if ns.Name == namespace {
						found = true
					}
				}
				assert.True(t, found)
			}

		})
	}
}

func Test_DeleteNameSpaceByName(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, *ControllerRuntimeClient){
		"success- namespace already exist": func(*testing.T) (bool, testOverrides, *ControllerRuntimeClient) {
			ns := &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        namespace,
					Annotations: map[string]string{},
				},
			}
			client := ctrlClientFake.NewClientBuilder().WithObjects(ns).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, testOverrides{}, controllerRuntimeClient
		},
		"success-  namespace does not exist": func(*testing.T) (bool, testOverrides, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, testOverrides{}, controllerRuntimeClient
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, controllerRuntimeClient := tc(t)
			if patch.getEnv != nil {
				oldGetEnv := utils.GetEnv
				defer func() { utils.GetEnv = oldGetEnv }()
				utils.GetEnv = patch.getEnv
			}

			err := controllerRuntimeClient.DeleteNameSpaceByName(context.TODO(), namespace)
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				namespaces := &corev1.NamespaceList{}
				err := controllerRuntimeClient.client.List(context.TODO(), namespaces)
				assert.NoError(t, err)
				found := false
				for _, ns := range namespaces.Items {
					if ns.Name == namespace {
						found = true
					}
				}
				assert.False(t, found)
			}

		})
	}
}

func Test_CreateSecretFromName(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, *ControllerRuntimeClient){
		"success- secret already exist": func(*testing.T) (bool, testOverrides, *ControllerRuntimeClient) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        secretName,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			}
			client := ctrlClientFake.NewClientBuilder().WithObjects(secret).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, testOverrides{}, controllerRuntimeClient
		},
		"success- secret exist but not in another namespace": func(*testing.T) (bool, testOverrides, *ControllerRuntimeClient) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        secretName,
					Namespace:   "new" + namespace,
					Annotations: map[string]string{},
				},
			}
			client := ctrlClientFake.NewClientBuilder().WithObjects(secret).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, testOverrides{}, controllerRuntimeClient
		},
		"success-  secret does not exist": func(*testing.T) (bool, testOverrides, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}
			return true, testOverrides{}, controllerRuntimeClient
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, controllerRuntimeClient := tc(t)
			if patch.getEnv != nil {
				oldGetEnv := utils.GetEnv
				defer func() { utils.GetEnv = oldGetEnv }()
				utils.GetEnv = patch.getEnv
			}

			err := controllerRuntimeClient.CreateSecretFromName(context.TODO(), secretName, namespace, []byte(""))
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				secret := &corev1.Secret{}
				err := controllerRuntimeClient.client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: namespace}, secret)
				assert.NoError(t, err)
			}

		})
	}

}

func Test_CreateTLSSecretFromName(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, *v1.Namespace, *ControllerRuntimeClient){
		"success-  secret already exist in namspace": func(*testing.T) (bool, *v1.Namespace, *ControllerRuntimeClient) {
			ns := &v1.Namespace{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Namespace",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        namespace,
					Annotations: map[string]string{},
				},
			}
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        secretName,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			}
			client := ctrlClientFake.NewClientBuilder().WithObjects(secret).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, ns, controllerRuntimeClient
		},
		"success- secret exist in another namespace": func(*testing.T) (bool, *v1.Namespace, *ControllerRuntimeClient) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        secretName,
					Namespace:   "new" + namespace,
					Annotations: map[string]string{},
				},
			}

			ns := &v1.Namespace{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Namespace",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        namespace,
					Annotations: map[string]string{},
				},
			}

			client := ctrlClientFake.NewClientBuilder().WithObjects(secret).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, ns, controllerRuntimeClient
		},
		"fail-  deserializing": func(*testing.T) (bool, *v1.Namespace, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}
			return false, &v1.Namespace{}, controllerRuntimeClient
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, ns, controllerRuntimeClient := tc(t)

			namespaceData, _ := json.Marshal(ns)

			err := controllerRuntimeClient.CreateTLSSecretFromName(context.TODO(), secretName, namespaceData, []byte(""), []byte(""))
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				secret := &corev1.Secret{}
				err := controllerRuntimeClient.client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: namespace}, secret)
				assert.NoError(t, err)
			}

		})
	}
}

func Test_CreateNameSpace(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, *v1.Namespace, *ControllerRuntimeClient){
		"success- namespace already exist": func(*testing.T) (bool, *v1.Namespace, *ControllerRuntimeClient) {
			ns := &v1.Namespace{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Namespace",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        namespace,
					Annotations: map[string]string{},
				},
			}
			client := ctrlClientFake.NewClientBuilder().WithObjects(ns).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return true, ns, controllerRuntimeClient
		},
		"success-  namespace does not exist": func(*testing.T) (bool, *v1.Namespace, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			ns := &v1.Namespace{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Namespace",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        namespace,
					Annotations: map[string]string{},
				},
			}

			return true, ns, controllerRuntimeClient
		},
		"fail-  deserializing": func(*testing.T) (bool, *v1.Namespace, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}
			return false, &v1.Namespace{}, controllerRuntimeClient
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, ns, controllerRuntimeClient := tc(t)

			namespaceData, _ := json.Marshal(ns)

			err := controllerRuntimeClient.CreateNameSpace(context.TODO(), namespaceData)
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				namespaces := &corev1.NamespaceList{}
				err := controllerRuntimeClient.client.List(context.TODO(), namespaces)
				assert.NoError(t, err)
				found := false
				for _, ns := range namespaces.Items {
					if ns.Name == namespace {
						found = true
					}
				}
				assert.True(t, found)
			}

		})
	}
}

func Test_CreateSecret(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, *corev1.Secret, *ControllerRuntimeClient){
		"success-  secret does not already exist": func(*testing.T) (bool, *corev1.Secret, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}
			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        secretName,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			}

			return true, secret, controllerRuntimeClient
		},
		"fail- secret already exist": func(*testing.T) (bool, *corev1.Secret, *ControllerRuntimeClient) {
			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        secretName,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			}
			client := ctrlClientFake.NewClientBuilder().WithObjects(secret).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return false, secret, controllerRuntimeClient
		},
		"fail-  deserializing": func(*testing.T) (bool, *corev1.Secret, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}
			return false, &corev1.Secret{}, controllerRuntimeClient
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, ns, controllerRuntimeClient := tc(t)

			secretData, _ := json.Marshal(ns)

			err := controllerRuntimeClient.CreateSecret(context.TODO(), secretData)
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				secret := &corev1.Secret{}
				err := controllerRuntimeClient.client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: namespace}, secret)
				assert.NoError(t, err)
			}

		})
	}

}

func Test_CreateConfigMap(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, *corev1.ConfigMap, *ControllerRuntimeClient){
		"success-  configMap does not already exist": func(*testing.T) (bool, *corev1.ConfigMap, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}
			secret := &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        configMapName,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			}
			return true, secret, controllerRuntimeClient
		},
		"fail- configMap already exist": func(*testing.T) (bool, *corev1.ConfigMap, *ControllerRuntimeClient) {
			secret := &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        configMapName,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			}
			client := ctrlClientFake.NewClientBuilder().WithObjects(secret).Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}

			return false, secret, controllerRuntimeClient
		},
		"fail-  deserializing": func(*testing.T) (bool, *corev1.ConfigMap, *ControllerRuntimeClient) {
			client := ctrlClientFake.NewClientBuilder().WithObjects().Build()
			controllerRuntimeClient := &ControllerRuntimeClient{
				Logger: echo.New().Logger,
				client: client,
			}
			return false, &corev1.ConfigMap{}, controllerRuntimeClient
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, ns, controllerRuntimeClient := tc(t)

			configMapData, _ := json.Marshal(ns)

			err := controllerRuntimeClient.CreateConfigMap(context.TODO(), configMapData)
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				secret := &corev1.ConfigMap{}
				err := controllerRuntimeClient.client.Get(context.TODO(), types.NamespacedName{Name: configMapName, Namespace: namespace}, secret)
				assert.NoError(t, err)
			}

		})
	}

}

func Test_GetControllerClient(t *testing.T) {
	createBadTestConfig := func() *clientcmdapi.Config {
		const (
			server = "https://127.0.0.1:6443"
			token  = "the-token"
		)

		config := clientcmdapi.NewConfig()
		config.Clusters["clean"] = &clientcmdapi.Cluster{
			Server: server,
		}
		config.AuthInfos["clean"] = &clientcmdapi.AuthInfo{
			Token: token,
		}
		config.Contexts["clean"] = &clientcmdapi.Context{
			Cluster:  "clean",
			AuthInfo: "clean",
		}
		config.CurrentContext = "clean"

		return config
	}

	tests := map[string]func(t *testing.T) (bool, testOverrides, *rest.Config){
		"success - no clientConfig": func(*testing.T) (bool, testOverrides, *rest.Config) {
			return true, testOverrides{
				ctrlClientNewWrapper: func(clientConfig *rest.Config, scheme *runtime.Scheme) (ctrlClient.Client, error) {
					return ctrlClientFake.NewClientBuilder().WithObjects().Build(), nil
				},
				getConfigWrapper: func() (*rest.Config, error) {
					config := createBadTestConfig()
					clientBuilder := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{})

					clientConfig, _ := clientBuilder.ClientConfig()
					return clientConfig, nil
				},
			}, nil
		},
		"fail - getting clientConfig": func(*testing.T) (bool, testOverrides, *rest.Config) {
			return false, testOverrides{
				ctrlClientNewWrapper: func(clientConfig *rest.Config, scheme *runtime.Scheme) (ctrlClient.Client, error) {
					return ctrlClientFake.NewClientBuilder().WithObjects().Build(), nil
				},
				getConfigWrapper: func() (*rest.Config, error) {
					config := createBadTestConfig()
					clientBuilder := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{})

					clientConfig, _ := clientBuilder.ClientConfig()
					return clientConfig, errors.New("getting clientConfig")
				},
			}, nil
		},
		"fail - test config ": func(*testing.T) (bool, testOverrides, *rest.Config) {
			config := createBadTestConfig()
			clientBuilder := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{})

			clientConfig, _ := clientBuilder.ClientConfig()

			return false, testOverrides{}, clientConfig
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, clientConfig := tc(t)

			if patch.ctrlClientNewWrapper != nil {
				oldCtrlClientNewWrapper := CtrlClientNewWrapper
				defer func() { CtrlClientNewWrapper = oldCtrlClientNewWrapper }()
				CtrlClientNewWrapper = patch.ctrlClientNewWrapper
			}
			if patch.getConfigWrapper != nil {
				oldGetConfigWrapper := GetConfigWrapper
				defer func() { GetConfigWrapper = oldGetConfigWrapper }()
				GetConfigWrapper = patch.getConfigWrapper
			}

			scheme := runtime.NewScheme()
			utilruntime.Must(clientgoscheme.AddToScheme(scheme))

			_, err := GetControllerClient(clientConfig, scheme)

			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func Test_ControllerRuntimeClient(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, string){
		"success": func(*testing.T) (bool, testOverrides, string) {
			content := `
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8080
    extensions:
    - name: client.authentication.k8s.io/exec
      extension:
        audience: foo
        other: bar
  name: foo-cluster
contexts:
- context:
    cluster: foo-cluster
    user: foo-user
    namespace: bar
  name: foo-context
current-context: foo-context
kind: Config
users:
- name: foo-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - arg-1
      - arg-2
      command: foo-command
      provideClusterInfo: true
`
			return true, testOverrides{
				ctrlClientNewWrapper: func(clientConfig *rest.Config, scheme *runtime.Scheme) (ctrlClient.Client, error) {
					return ctrlClientFake.NewClientBuilder().WithObjects().Build(), nil
				},
			}, content
		},
		"fail - getting client": func(*testing.T) (bool, testOverrides, string) {
			content := `
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8080
    extensions:
    - name: client.authentication.k8s.io/exec
      extension:
        audience: foo
        other: bar
  name: foo-cluster
contexts:
- context:
    cluster: foo-cluster
    user: foo-user
    namespace: bar
  name: foo-context
current-context: foo-context
kind: Config
users:
- name: foo-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - arg-1
      - arg-2
      command: foo-command
      provideClusterInfo: true
`
			return false, testOverrides{
				ctrlClientNewWrapper: func(clientConfig *rest.Config, scheme *runtime.Scheme) (ctrlClient.Client, error) {
					return ctrlClientFake.NewClientBuilder().WithObjects().Build(), errors.New(" error getting client")
				},
			}, content
		},

		"fail - RESTConfigFromKubeConfig": func(*testing.T) (bool, testOverrides, string) {
			content := `
apiVersion: v1
clusters:
- cluster:
	certificate-authority-data: ZGF0YS1oZXJl
	server: https://127.0.0.1:6443
	name: kubernetes
contexts:
- context:
	cluster: kubernetes
	user: kubernetes-admin
	name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
	user:
	client-certificate-data: ZGF0YS1oZXJl
	client-key-data: ZGF0YS1oZXJl			
`
			return false, testOverrides{}, content
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, data := tc(t)

			if patch.ctrlClientNewWrapper != nil {
				oldCtrlClientNewWrapper := CtrlClientNewWrapper
				defer func() { CtrlClientNewWrapper = oldCtrlClientNewWrapper }()
				CtrlClientNewWrapper = patch.ctrlClientNewWrapper
			}

			_, err := NewControllerRuntimeClient([]byte(data), echo.New().Logger)

			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func Test_IsOpenShift(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides){
		"success ": func(*testing.T) (bool, testOverrides) {
			return true, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					fakeDiscovery, ok := fakeClientSet.Discovery().(*discoveryfake.FakeDiscovery)
					if !ok {
						t.Fatalf("couldn't convert Discovery() to *FakeDiscovery")
					}
					fakeDiscovery.Resources = []*metav1.APIResourceList{
						{
							APIResources: []metav1.APIResource{
								{Name: "security.openshift.io"},
							},
							GroupVersion: "security.openshift.io/v1",
						},
					}
					return fakeClientSet, nil
				},
			}
		},
		"fail - not found ": func(*testing.T) (bool, testOverrides) {
			return false, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					fakeDiscovery, ok := fakeClientSet.Discovery().(*discoveryfake.FakeDiscovery)
					if !ok {
						t.Fatalf("couldn't convert Discovery() to *FakeDiscovery")
					}
					fakeDiscovery.Resources = []*metav1.APIResourceList{
						{
							APIResources: []metav1.APIResource{
								{Name: "security.k8s.io"},
							},
							GroupVersion: "security.k8s.io/v1",
						},
					}
					return fakeClientSet, nil
				},
			}
		},
		"fail- bad version ": func(*testing.T) (bool, testOverrides) {
			return false, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					fakeDiscovery, ok := fakeClientSet.Discovery().(*discoveryfake.FakeDiscovery)
					if !ok {
						t.Fatalf("couldn't convert Discovery() to *FakeDiscovery")
					}
					fakeDiscovery.Resources = []*metav1.APIResourceList{
						{
							APIResources: []metav1.APIResource{
								{Name: "security.openshift.io"},
							},
							GroupVersion: "security.openshift.io////v1",
						},
					}
					return fakeClientSet, nil
				},
			}
		},
		"bad config data ": func(*testing.T) (bool, testOverrides) {
			return false, testOverrides{}
		},
		"fail - to get client set": func(*testing.T) (bool, testOverrides) {
			return false, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					return fake.NewSimpleClientset(), errors.New(" error listing pods")
				},
			}

		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch := tc(t)

			if patch.getClientSetWrapper != nil {
				oldGetClientSetWrapper := GetClientSetWrapper
				defer func() { GetClientSetWrapper = oldGetClientSetWrapper }()
				GetClientSetWrapper = patch.getClientSetWrapper
			}

			c := Client{}
			isOpenshift, err := c.IsOpenShift([]byte(""))
			if !success {
				assert.False(t, isOpenshift)
			} else {
				assert.NoError(t, err)
				assert.True(t, isOpenshift)
			}

		})
	}

}
func Test_GetAPIResource(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, string, testOverrides){
		"success ": func(*testing.T) (bool, string, testOverrides) {
			resourceName := "hello-world"
			return true, resourceName, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					return fakeClientSet, nil
				},
				serverPreferredResourcesWrapper: func(_ kubernetes.Interface) ([]*metav1.APIResourceList, error) {
					return []*metav1.APIResourceList{
						{
							APIResources: []metav1.APIResource{
								{Name: "resource-name", Kind: "hello-world"},
							},
							GroupVersion: "v1",
						},
					}, nil
				},
			}

		},
		"fail-resources ": func(*testing.T) (bool, string, testOverrides) {
			resourceName := "not-found"
			return false, resourceName, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					return fakeClientSet, nil
				},
			}
		},
		"fail - to get client set": func(*testing.T) (bool, string, testOverrides) {
			return false, "", testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					return fake.NewSimpleClientset(), errors.New(" error listing pods")
				},
			}

		},
		"fail - to get client server": func(*testing.T) (bool, string, testOverrides) {
			return false, "", testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					return fakeClientSet, nil
				},
				serverPreferredResourcesWrapper: func(_ kubernetes.Interface) ([]*metav1.APIResourceList, error) {
					return []*metav1.APIResourceList{}, errors.New(" error listing pods")
				},
			}

		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, resourceName, patch := tc(t)

			if patch.getClientSetWrapper != nil {
				oldGetClientSetWrapper := GetClientSetWrapper
				defer func() { GetClientSetWrapper = oldGetClientSetWrapper }()
				GetClientSetWrapper = patch.getClientSetWrapper
			}

			if patch.serverPreferredResourcesWrapper != nil {
				oldServerPreferredResourcesWrapper := ServerPreferredResourcesWrapper
				defer func() { ServerPreferredResourcesWrapper = oldServerPreferredResourcesWrapper }()
				ServerPreferredResourcesWrapper = patch.serverPreferredResourcesWrapper
			}

			c := Client{}
			_, _, err := c.GetAPIResource([]byte(""), resourceName)
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func Test_GetVersion(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, string, string, testOverrides){
		"success ": func(*testing.T) (bool, string, string, testOverrides) {
			major := "2"
			minor := "9"
			return true, major, minor, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					fakeClientSet.Discovery().(*discoveryfake.FakeDiscovery).FakedServerVersion = &version.Info{
						Major: major,
						Minor: minor,
					}
					return fakeClientSet, nil
				},
			}

		},
		"fail - to get client set": func(*testing.T) (bool, string, string, testOverrides) {
			return false, "", "", testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					return fake.NewSimpleClientset(), errors.New(" error listing pods")
				},
			}

		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, major, minor, patch := tc(t)

			if patch.getClientSetWrapper != nil {
				oldGetClientSetWrapper := GetClientSetWrapper
				defer func() { GetClientSetWrapper = oldGetClientSetWrapper }()
				GetClientSetWrapper = patch.getClientSetWrapper
			}

			c := Client{}
			out, err := c.GetVersion([]byte(""))
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, out, fmt.Sprintf("%s.%s", major, minor))
			}

		})
	}
}

func Test_GetClientSet(t *testing.T) {
	c := Client{}
	content := `
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8080
    extensions:
    - name: client.authentication.k8s.io/exec
      extension:
        audience: foo
        other: bar
  name: foo-cluster
contexts:
- context:
    cluster: foo-cluster
    user: foo-user
    namespace: bar
  name: foo-context
current-context: foo-context
kind: Config
users:
- name: foo-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - arg-1
      - arg-2
      command: foo-command
      provideClusterInfo: true
`
	_, err := c.GetClientSet([]byte(content))
	assert.NoError(t, err)

	content = `
apiVersion: v1
clusters:
- cluster:
	certificate-authority-data: data-here
	server: https://127.0.0.1:6443
	name: kubernetes
contexts:
- context:
	cluster: kubernetes
	user: kubernetes-admin
	name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
	user:
	client-certificate-data: data-here
	client-key-data: data-here
	`
	_, err = c.GetClientSet([]byte(content))
	assert.Error(t, err)
}

func Test_DiscoverK8sDetails(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides){
		"success ": func(*testing.T) (bool, testOverrides) {
			return true, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					fakeDiscovery, ok := fakeClientSet.Discovery().(*discoveryfake.FakeDiscovery)
					if !ok {
						t.Fatalf("couldn't convert Discovery() to *FakeDiscovery")
					}
					fakeDiscovery.Resources = []*metav1.APIResourceList{
						{
							APIResources: []metav1.APIResource{
								{Name: "security.openshift.io"},
							},
							GroupVersion: "security.openshift.io/v1",
						},
					}
					fakeClientSet.Discovery().(*discoveryfake.FakeDiscovery).FakedServerVersion = &version.Info{
						Major: "1",
						Minor: "2",
					}
					return fakeClientSet, nil
				},
			}
		},
		"fail- getting version ": func(*testing.T) (bool, testOverrides) {
			return false, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					return fake.NewSimpleClientset(), errors.New(" error listing pods")
				},
			}
		},
		"fail- getting isopenshift ": func(*testing.T) (bool, testOverrides) {
			return false, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset()
					fakeDiscovery, ok := fakeClientSet.Discovery().(*discoveryfake.FakeDiscovery)
					if !ok {
						t.Fatalf("couldn't convert Discovery() to *FakeDiscovery")
					}
					fakeDiscovery.Resources = []*metav1.APIResourceList{
						{
							APIResources: []metav1.APIResource{
								{Name: "security.openshift.io"},
							},
							GroupVersion: "security.openshift.io////v1",
						},
					}
					return fakeClientSet, nil
				},
			}
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch := tc(t)

			if patch.getClientSetWrapper != nil {
				oldGetClientSetWrapper := GetClientSetWrapper
				defer func() { GetClientSetWrapper = oldGetClientSetWrapper }()
				GetClientSetWrapper = patch.getClientSetWrapper
			}

			c := Client{}
			_, _, _, err := c.DiscoverK8sDetails([]byte(""))
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func Test_GetCertManagerPods(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, string, testOverrides){
		"success ": func(*testing.T) (bool, string, testOverrides) {
			label := "hello-world"
			return true, label, testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					fakeClientSet := fake.NewSimpleClientset(&v1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name:        csmDataCollector + "-testing1",
							Namespace:   namespace,
							Annotations: map[string]string{},
							Labels:      map[string]string{"app": label},
						},
						Status: v1.PodStatus{
							InitContainerStatuses: []v1.ContainerStatus{
								{
									Name:  "csm-init",
									State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 0}},
								},
							},
							HostIP: "1.2.3",
						},
					},
					)
					return fakeClientSet, nil
				},
				retryableGetPodsTimeOut: 1 * time.Second,
			}

		},
		"fail - to get client set": func(*testing.T) (bool, string, testOverrides) {
			return false, "", testOverrides{
				getClientSetWrapper: func(c Client, data []byte) (kubernetes.Interface, error) {
					return nil, errors.New(" error listing pods")
				},
			}

		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, label, patch := tc(t)

			if patch.getClientSetWrapper != nil {
				oldGetClientSetWrapper := GetClientSetWrapper
				defer func() { GetClientSetWrapper = oldGetClientSetWrapper }()
				GetClientSetWrapper = patch.getClientSetWrapper
			}

			if patch.retryableGetPodsTimeOut != 0 {
				oldRetryableGetPodsTimeOut := RetryableGetPodsTimeOut
				defer func() { RetryableGetPodsTimeOut = oldRetryableGetPodsTimeOut }()
				RetryableGetPodsTimeOut = patch.retryableGetPodsTimeOut
			}

			c := Client{}
			_, err := c.GetCertManagerPods([]byte(""), namespace, label)
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}

func Test_NodeDataCollectorInstall(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, *NodeDataCollector){
		"success- with out csm-data-collector and namespace(greenfield)": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset()

			collector := &NodeDataCollector{
				Logger:          echo.New().Logger,
				ClientSet:       fakeClientSet,
				InstallWaitTime: 1 * time.Second,
			}
			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
			}
			return true, patch, collector
		},
		"success - with existing csm-data-collector and namespace(brownfield)": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset(&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        namespace,
					Annotations: map[string]string{},
				},
			}, &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			})

			collector := &NodeDataCollector{
				Logger:          echo.New().Logger,
				ClientSet:       fakeClientSet,
				InstallWaitTime: 1 * time.Second,
			}
			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
			}
			return true, patch, collector
		},
		"fail - invalid configuration. Data collector image not set": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			return false, testOverrides{}, &NodeDataCollector{}
		},
		"fail - create namespace": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset()
			fakeClientSet.PrependReactor("create", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &v1.Namespace{}, errors.New(" error creating namespace")
			})

			collector := &NodeDataCollector{
				Logger:          echo.New().Logger,
				ClientSet:       fakeClientSet,
				InstallWaitTime: time.Duration(1),
			}

			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
			}

			return false, patch, collector
		},
		"fail - create deamonsets": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset(&appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			})
			fakeClientSet.PrependReactor("create", "daemonsets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &appsv1.DaemonSet{}, errors.New(" error creating daemonsets")
			})

			collector := &NodeDataCollector{
				Logger:          echo.New().Logger,
				ClientSet:       fakeClientSet,
				InstallWaitTime: time.Duration(1),
			}

			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
			}

			return false, patch, collector
		},
		"fail - create not found deamonsets": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset()
			fakeClientSet.PrependReactor("create", "daemonsets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &appsv1.DaemonSet{}, errors.New(" error creating daemonsets")
			})

			collector := &NodeDataCollector{
				Logger:          echo.New().Logger,
				ClientSet:       fakeClientSet,
				InstallWaitTime: time.Duration(1),
			}

			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
			}

			return false, patch, collector
		},
		"fail - delete deamonsets": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset(&appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			})
			fakeClientSet.PrependReactor("delete", "daemonsets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &appsv1.DaemonSet{}, errors.New(" error deleting daemonsets")
			})

			collector := &NodeDataCollector{
				Logger:          echo.New().Logger,
				ClientSet:       fakeClientSet,
				InstallWaitTime: time.Duration(1),
			}

			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
			}

			return false, patch, collector
		},
		"fail - get deamonsets(other kind of error)": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset(&appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector,
					Namespace:   namespace,
					Annotations: map[string]string{},
				},
			})
			fakeClientSet.PrependReactor("get", "daemonsets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &appsv1.DaemonSet{}, errors.New(" error getting daemonsets")
			})

			collector := &NodeDataCollector{
				Logger:          echo.New().Logger,
				ClientSet:       fakeClientSet,
				InstallWaitTime: time.Duration(1),
			}

			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
			}

			return false, patch, collector
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, collector := tc(t)
			if patch.getEnv != nil {
				oldGetEnv := utils.GetEnv
				defer func() { utils.GetEnv = oldGetEnv }()
				utils.GetEnv = patch.getEnv
			}

			err := collector.Install()
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// check that namespace and deamonser exist
				_, err = collector.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
				assert.NoError(t, err)
				_, err = collector.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), csmDataCollector, metav1.GetOptions{})
				assert.NoError(t, err)
			}

		})
	}

}

// tests is too slow
func Test_NodeDataCollectorCollect(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, *NodeDataCollector){
		"success": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector + "-testing1",
					Namespace:   namespace,
					Annotations: map[string]string{},
					Labels:      map[string]string{"app": "csm-data-collector"},
				},
				Status: v1.PodStatus{
					InitContainerStatuses: []v1.ContainerStatus{
						{
							Name:  "csm-init",
							State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 0}},
						},
					},
					HostIP: "1.2.3",
				},
			}, &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector + "-testing2",
					Namespace:   namespace,
					Annotations: map[string]string{},
					Labels:      map[string]string{"app": "csm-data-collector"},
				},
				Status: v1.PodStatus{
					InitContainerStatuses: []v1.ContainerStatus{
						{
							Name:  "csm-init",
							State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 2}},
						},
					},
					HostIP: "1.2.4",
				},
			}, &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector + "-testing3",
					Namespace:   namespace,
					Annotations: map[string]string{},
					Labels:      map[string]string{"app": "csm-data-collector"},
				},
				Status: v1.PodStatus{
					InitContainerStatuses: []v1.ContainerStatus{
						{
							Name:  "csm-init",
							State: v1.ContainerState{Running: &v1.ContainerStateRunning{}},
						},
					},
					HostIP: "1.2.5",
				},
			})

			collector := &NodeDataCollector{
				Logger:                     echo.New().Logger,
				ClientSet:                  fakeClientSet,
				InstallWaitTime:            1 * time.Second,
				HandleTerminatedPodTimeout: 1 * time.Second,
				HandlePendingPodTimeout:    1 * time.Second,
				HandlePendingPodsWaitTime:  500 * time.Millisecond,
			}
			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
				minPollChTimeOut:        1 * time.Second,
				retryableGetPodsTimeOut: 1 * time.Second,
			}

			return true, patch, collector
		},
		"success -- all terminated": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector + "-testing1",
					Namespace:   namespace,
					Annotations: map[string]string{},
					Labels:      map[string]string{"app": "csm-data-collector"},
				},
				Status: v1.PodStatus{
					InitContainerStatuses: []v1.ContainerStatus{
						{
							Name:  "csm-init",
							State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 0}},
						},
					},
					HostIP: "1.2.3",
				},
			},
			)

			collector := &NodeDataCollector{
				Logger:                     echo.New().Logger,
				ClientSet:                  fakeClientSet,
				InstallWaitTime:            1 * time.Second,
				HandleTerminatedPodTimeout: 1 * time.Second,
				HandlePendingPodTimeout:    1 * time.Second,
				HandlePendingPodsWaitTime:  500 * time.Millisecond,
			}
			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
				minPollChTimeOut:        1 * time.Second,
				retryableGetPodsTimeOut: 1 * time.Second,
			}

			return true, patch, collector
		},
		"success -- all pending": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset(&v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        csmDataCollector + "-testing3",
					Namespace:   namespace,
					Annotations: map[string]string{},
					Labels:      map[string]string{"app": "csm-data-collector"},
				},
				Status: v1.PodStatus{
					InitContainerStatuses: []v1.ContainerStatus{
						{
							Name:  "csm-init",
							State: v1.ContainerState{Running: &v1.ContainerStateRunning{}},
						},
					},
					HostIP: "1.2.5",
				},
			})

			collector := &NodeDataCollector{
				Logger:                     echo.New().Logger,
				ClientSet:                  fakeClientSet,
				InstallWaitTime:            1 * time.Second,
				HandleTerminatedPodTimeout: 1 * time.Second,
				HandlePendingPodTimeout:    1 * time.Second,
				HandlePendingPodsWaitTime:  500 * time.Millisecond,
			}
			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
				minPollChTimeOut:        1 * time.Second,
				retryableGetPodsTimeOut: 1 * time.Second,
			}

			return true, patch, collector
		},
		"fail - failed to install the data collector daemonset": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			return false, testOverrides{}, &NodeDataCollector{}
		},
		"fail - failed to query for the daemonset pods": func(*testing.T) (bool, testOverrides, *NodeDataCollector) {
			fakeClientSet := fake.NewSimpleClientset()
			fakeClientSet.PrependReactor("list", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &v1.PodList{}, errors.New(" error listing pods")
			})

			collector := &NodeDataCollector{
				Logger:                     echo.New().Logger,
				ClientSet:                  fakeClientSet,
				InstallWaitTime:            1 * time.Second,
				HandleTerminatedPodTimeout: 1 * time.Second,
				HandlePendingPodTimeout:    1 * time.Second,
				HandlePendingPodsWaitTime:  500 * time.Millisecond,
			}

			patch := testOverrides{
				getEnv: func(envName string, defaultValue string) string {
					return envName
				},
				retryableGetPodsTimeOut: 1 * time.Second,
				minPollChTimeOut:        1 * time.Second,
			}

			return false, patch, collector
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, collector := tc(t)
			if patch.getEnv != nil {
				oldGetEnv := utils.GetEnv
				defer func() { utils.GetEnv = oldGetEnv }()
				utils.GetEnv = patch.getEnv
			}
			if patch.retryableGetPodsTimeOut != 0 {
				oldRetryableGetPodsTimeOut := RetryableGetPodsTimeOut
				defer func() { RetryableGetPodsTimeOut = oldRetryableGetPodsTimeOut }()
				RetryableGetPodsTimeOut = patch.retryableGetPodsTimeOut
			}
			if patch.minPollChTimeOut != 0 {
				oldMinPollChTimeOut := MinPollChTimeOut
				defer func() { MinPollChTimeOut = oldMinPollChTimeOut }()
				MinPollChTimeOut = patch.retryableGetPodsTimeOut
			}

			_, err := collector.Collect()
			if !success {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ConsolidateNodeDiscoveryDetails(t *testing.T) {
	tests := map[string]func(t *testing.T) ([]Node, []Node){
		"success": func(*testing.T) ([]Node, []Node) {
			inputNodes := []Node{
				{
					HostName: "host-1", InstalledSoftware: map[string]string{"sdc": "enabled"},
				},
				{
					HostName: "host-1", InstalledSoftware: map[string]string{"iscsi": "enabled"},
				},
				{
					HostName: "host-2", InstalledSoftware: map[string]string{"iscsi": "enabled"},
				},
				{
					HostName: "host-3", InstalledSoftware: nil,
				},
				{},
			}
			expectedOutputNodes := []Node{
				{
					HostName: "host-1", InstalledSoftware: map[string]string{"sdc": "enabled", "iscsi": "enabled"},
				},
				{
					HostName: "host-2", InstalledSoftware: map[string]string{"iscsi": "enabled"},
				},
				{
					HostName: "host-3", InstalledSoftware: map[string]string{},
				},
			}
			return inputNodes, expectedOutputNodes
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inputNodes, expectedOutputNodes := tc(t)
			outputNodes := consolidateNodeDiscoveryDetails(inputNodes)
			assert.Equal(t, len(expectedOutputNodes), len(outputNodes))
			for _, node := range outputNodes {
				foundNode := findNode(t, node.HostName, expectedOutputNodes)
				assert.Equal(t, foundNode, node)
			}
		})
	}
}

func findNode(t *testing.T, hostname string, nodes []Node) Node {
	for _, node := range nodes {
		if node.HostName == hostname {
			return node
		}
	}
	t.Fatalf("unable to find host with HostName=%s", hostname)
	return Node{}
}

func Test_GetInstalledSoftware(t *testing.T) {
	tests := map[string]func(t *testing.T) (string, map[string]string){
		"success": func(*testing.T) (string, map[string]string) {
			logs := "sdc=123\niscsi=abc"
			expectedResult := map[string]string{
				"sdc":   "enabled",
				"iscsi": "enabled",
			}
			return logs, expectedResult
		},
		"skip empty software values": func(*testing.T) (string, map[string]string) {
			logs := "sdc=123\niscsi="
			expectedResult := map[string]string{
				"sdc": "enabled",
			}
			return logs, expectedResult
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			logs, expectedResult := tc(t)
			outputResult := getInstalledSoftware(logs)
			assert.True(t, reflect.DeepEqual(outputResult, expectedResult))
		})
	}
}
