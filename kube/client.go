package kube

import (
	"GomokuGame/app/conf"
	"github.com/sirupsen/logrus"
	apiV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

var clientset *kubernetes.Clientset
var clientsetOnce sync.Once

func InitClientset() {
	clientsetOnce.Do(func() {
		config, err := clientcmd.BuildConfigFromFlags("", conf.KUBE_CONFIG_PATH)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"config_path": conf.KUBE_CONFIG_PATH,
				"err":         err.Error(),
			}).Fatal("BuildConfigFromFlags failed.")
		}

		cs, err := kubernetes.NewForConfig(config)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Fatal("Init clientset failed.")
		}
		clientset = cs
	})
}

func NewKubePodsClient() coreV1.PodInterface {
	return clientset.CoreV1().Pods(apiV1.NamespaceDefault)
}
