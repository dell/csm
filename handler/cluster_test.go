package handler

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dell/csm-deployment/router"
	"github.com/dell/csm-deployment/store/mocks"
	"github.com/golang/mock/gomock"

	handlerMocks "github.com/dell/csm-deployment/handler/mocks"
	"github.com/dell/csm-deployment/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCluster(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemResponseJSON := `{"cluster_id":0,"cluster_name":"cluster-name","nodes":""}`

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().Create(gomock.Any()).Times(1)
			isOpenShift := false

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("1.18", &isOpenShift, nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusCreated, handler, body, writer, createStorageSystemResponseJSON, ctrl
		},
		"success with openshift": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemResponseJSON := `{"cluster_id":0,"cluster_name":"cluster-name","nodes":""}`

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().Create(gomock.Any()).Times(1)
			isOpenShift := true

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("1.18", &isOpenShift, nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusCreated, handler, body, writer, createStorageSystemResponseJSON, ctrl
		},
		"error saving to database": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().Create(gomock.Any()).Times(1).Return(errors.New("error saving to database"))
			isOpenShift := true

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("1.18", &isOpenShift, nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusInternalServerError, handler, body, writer, "", ctrl
		},
		"error communicating with cluster": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("", nil, nil, errors.New("error communicating with cluster"))
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusBadRequest, handler, body, writer, "", ctrl
		},
		"error empty cluster name": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusBadRequest, handler, body, writer, "", ctrl
		},
		"error empty file upload": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "cluster-name")
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusBadRequest, handler, body, writer, "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, body, writer, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPost, "/", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, handler.CreateCluster(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_UpdateCluster(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemResponseJSON := `{"cluster_id":0,"cluster_name":"new-cluster-name","nodes":""}`

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "new-cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Cluster{ClusterName: "old-cluster-name"}, nil)
			clusterStore.EXPECT().Update(gomock.Any()).Times(1)
			isOpenShift := false

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("1.18", &isOpenShift, nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusOK, handler, body, writer, "1", createStorageSystemResponseJSON, ctrl
		},
		"success updating name of cluster": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemResponseJSON := `{"cluster_id":0,"cluster_name":"new-cluster-name","nodes":""}`

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "new-cluster-name")
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Cluster{ClusterName: "old-cluster-name"}, nil)
			clusterStore.EXPECT().Update(gomock.Any()).Times(1)
			isOpenShift := false

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("1.18", &isOpenShift, nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusOK, handler, body, writer, "1", createStorageSystemResponseJSON, ctrl
		},
		"success with openshift": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemResponseJSON := `{"cluster_id":0,"cluster_name":"cluster-name","nodes":""}`

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Cluster{ClusterName: "old-cluster-name"}, nil)
			clusterStore.EXPECT().Update(gomock.Any()).Times(1)
			isOpenShift := true

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("1.18", &isOpenShift, nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusOK, handler, body, writer, "1", createStorageSystemResponseJSON, ctrl
		},
		"error looking up cluster in database": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "new-cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("cluster not found"))
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusInternalServerError, handler, body, writer, "1", "", ctrl
		},
		"error cluster not found in database": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "new-cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusNotFound, handler, body, writer, "1", "", ctrl
		},
		"error saving to database": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Cluster{ClusterName: "old-cluster-name"}, nil)
			clusterStore.EXPECT().Update(gomock.Any()).Times(1).Return(errors.New("error saving to database"))
			isOpenShift := true

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("1.18", &isOpenShift, nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusInternalServerError, handler, body, writer, "1", "", ctrl
		},
		"error communicating with cluster": func(*testing.T) (int, *ClusterHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "cluster-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample kubeconfig file contents`))
			writer.Close()

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			k8sClient := handlerMocks.NewMockK8sClientInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Cluster{ClusterName: "old-cluster-name"}, nil)

			k8sClient.EXPECT().DiscoverK8sDetails(gomock.Any()).Times(1).Return("", nil, nil, errors.New("error communicating with cluster"))
			handler := &ClusterHandler{clusterStore: clusterStore, k8sClient: k8sClient}
			return http.StatusBadRequest, handler, body, writer, "1", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, body, writer, clusterID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPatch, "/", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/clusters/:id")
			c.SetParamNames("id")
			c.SetParamValues(clusterID)

			assert.NoError(t, handler.UpdateCluster(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_GetCluster(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ClusterHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ClusterHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getStorageSystemResponseJSON := `{"cluster_id":1,"cluster_name":"cluster-1","nodes":""}`

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			cluster := model.Cluster{
				ClusterName: "cluster-1",
			}
			cluster.ID = 1
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&cluster, nil)
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusOK, handler, "1", getStorageSystemResponseJSON, ctrl
		},
		"nil result from db": func(*testing.T) (int, *ClusterHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusNotFound, handler, "1", "", ctrl
		},
		"error querying db": func(*testing.T) (int, *ClusterHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusInternalServerError, handler, "1", "", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *ClusterHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusUnprocessableEntity, handler, "abc", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, clusterID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/clusters/:id")
			c.SetParamNames("id")
			c.SetParamValues(clusterID)

			assert.NoError(t, handler.GetCluster(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_ListClusters(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ClusterHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ClusterHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			listStorageSystemResponseJSON := "[{\"cluster_id\":0,\"cluster_name\":\"cluster-1\",\"nodes\":\"\"},{\"cluster_id\":0,\"cluster_name\":\"cluster-2\",\"nodes\":\"\"}]"

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)

			clusters := make([]model.Cluster, 0)
			clusters = append(clusters, model.Cluster{
				ClusterName: "cluster-1",
			})
			clusters = append(clusters, model.Cluster{
				ClusterName: "cluster-2",
			})
			clusterStore.EXPECT().GetAll().Times(1).Return(clusters, nil)
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusOK, handler, listStorageSystemResponseJSON, ctrl
		},
		"error querying database": func(*testing.T) (int, *ClusterHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetAll().Times(1).Return(nil, errors.New("error"))
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusInternalServerError, handler, "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, handler.ListClusters(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_DeleteCluster(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ClusterHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ClusterHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Cluster{}, nil)
			clusterStore.EXPECT().Delete(gomock.Any()).Times(1)
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusNoContent, handler, "1", ctrl
		},
		"nil result from db": func(*testing.T) (int, *ClusterHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusNotFound, handler, "1", ctrl
		},
		"error getting from db": func(*testing.T) (int, *ClusterHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
		"error deleting from db": func(*testing.T) (int, *ClusterHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Cluster{}, nil)
			clusterStore.EXPECT().Delete(gomock.Any()).Times(1).Return(errors.New("error"))
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *ClusterHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			handler := &ClusterHandler{clusterStore: clusterStore}
			return http.StatusUnprocessableEntity, handler, "abc", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, storageSystemID, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/clusters/:id")
			c.SetParamNames("id")
			c.SetParamValues(storageSystemID)

			assert.NoError(t, handler.DeleteCluster(c))
			assert.Equal(t, expectedStatus, rec.Code)
			ctrl.Finish()
		})
	}
}
