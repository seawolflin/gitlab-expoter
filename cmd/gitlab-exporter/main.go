package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"log"
	"net/http"
)

// 匿名导入，为了执行collector的init方法，用于注册prometheus的方法注册collector
import (
	_ "github.com/seawolflin/gitlab-exporter/internal/collector"
	_ "github.com/seawolflin/gitlab-exporter/internal/db"
)

func main() {
	context.GetInstance().Parse()

	initializer.InitAll()

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
