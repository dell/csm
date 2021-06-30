package main

import (
	"fmt"
	"net/url"
	"path"

	"github.com/dell/csm-deployment/db"
	"github.com/dell/csm-deployment/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/dell/csm-deployment/handler"
	"github.com/dell/csm-deployment/kapp"
	"github.com/dell/csm-deployment/router"
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
)

// @title CSM Deployment API
// @version 1.0
// @description CSM Deployment API
// @title CSM Deployment API

// @BasePath /api

// @produce	application/json
// @consumes application/json

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	scheme := utils.GetEnv("SCHEME", "https")
	hostName := utils.GetEnv("HOST", "127.0.0.1")
	port := utils.GetEnv("PORT", "8080")
	certDir := utils.GetEnv("CERT_DIR", "samplecerts")
	certFileName := utils.GetEnv("CERT_FILE", "samplecert.crt")
	keyFileName := utils.GetEnv("KEY_FILE", "samplecert.key")
	dbDir := utils.GetEnv("DB_DIR", "")
	hostNameWithPort := fmt.Sprintf("%s:%s", hostName, port)

	// Update docs
	docs.SwaggerInfo.Schemes = append(docs.SwaggerInfo.Schemes, scheme)
	docs.SwaggerInfo.Host = hostNameWithPort

	swaggerURL := url.URL{
		Scheme: scheme,
		Host:   hostNameWithPort,
		Path:   "swagger/doc.json",
	}

	rt := router.New()

	rt.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.URL(swaggerURL.String())))

	api := rt.Group("/api")

	d, err := db.New(dbDir)
	if err != nil {
		rt.Logger.Fatal("Error in initializing db", err.Error())
	}

	db.AutoMigrate(d)
	db.PopulateInventory(d)

	us := store.NewUserStore(d)
	h := handler.New(us)
	h.Register(api)

	applicationStateChanges := store.NewApplicationStateChangeStore(d)

	clusters := store.NewClusterStore(d)
	hc := handler.NewClusterHandler(clusters)
	hc.Register(api)

	tasks := store.NewTaskStore(d)

	applications := store.NewApplicationStore(d)
	arrays := store.NewStorageArrayStore(d)
	modules := store.NewModuleStore(d)

	as := handler.NewApplicationHandler(applications, tasks, clusters, applicationStateChanges, arrays, modules)
	as.Register(api)

	th := handler.NewTaskHandler(tasks, applications, applicationStateChanges, clusters, kapp.NewClient(""))
	th.Register(api)

	storageArrays := store.NewStorageArrayStore(d)
	sah := handler.NewStorageArrayHandler(storageArrays)
	sah.Register(api)

	if scheme == "http" {
		rt.Logger.Fatal(rt.Start(hostNameWithPort))
	} else if scheme == "https" {
		rt.Logger.Fatal(rt.StartTLS(hostNameWithPort, path.Join(certDir, certFileName), path.Join(certDir, keyFileName)))
	} else {
		rt.Logger.Fatal("unknown scheme specified")
	}

}
