package kube

import (
	"GomokuGame/app/conf"
	"encoding/json"
	"github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
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

func NewKubePodsClient() corev1.PodInterface {
	return clientset.CoreV1().Pods(apiv1.NamespaceDefault)
}

func CreateMatchPodResourceFile(gameId string, player1FirstHand string, maxThinkingTime string, player1Name string, player2Name string, player1ImageAddr string, player2ImageAddr string) *apiv1.Pod {
	var r apiv1.ResourceRequirements
	resourceLimitStr := "{\"limits\": {\"cpu\":\"" + conf.CPU_LIMIT + "\", \"memory\":\"" + conf.MEMEORY_LIMIT + "\"}}"
	json.Unmarshal([]byte(resourceLimitStr), &r)

	matchPod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "metav1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "match-game-" + gameId,
			Labels: map[string]string{
				"app": "match-game-" + gameId,
			},
		},
		Spec: apiv1.PodSpec{
			Volumes: []apiv1.Volume{apiv1.Volume{
				Name:         "match-logs",
				VolumeSource: apiv1.VolumeSource{},
			}},
			//InitContainers:                nil,
			Containers: []apiv1.Container{
				// referee container
				apiv1.Container{
					Name:            "gomoku-game-referee",
					Image:           conf.GAME_REFEREE_IMAGE_ADDR,
					ImagePullPolicy: "Never",
					Ports: []apiv1.ContainerPort{
						apiv1.ContainerPort{
							Name:          "referee-port",
							ContainerPort: 10003,
						},
					},
					VolumeMounts: []apiv1.VolumeMount{
						apiv1.VolumeMount{
							Name:      "match-logs",
							MountPath: "/match_logs",
						},
					},
					Env: []apiv1.EnvVar{
						apiv1.EnvVar{
							Name:  "PLAYER1_FIRST_HAND",
							Value: player1FirstHand,
						},
						apiv1.EnvVar{
							Name:  "MAX_THINKING_TIME",
							Value: maxThinkingTime,
						},
						apiv1.EnvVar{
							Name:  "GAME_ID",
							Value: gameId,
						},
						apiv1.EnvVar{
							Name:  "PLAYER1_ID",
							Value: player1Name,
						},
						apiv1.EnvVar{
							Name:  "PLAYER2_ID",
							Value: player2Name,
						},
						apiv1.EnvVar{
							Name:  "NSQ_PUBLISH_ADDR",
							Value: conf.NSQ_PUB_ADDR,
						},
						apiv1.EnvVar{
							Name:  "LOG_VOLUME_ADDR_PLAYER1",
							Value: "/match_logs/player1_log",
						},
						apiv1.EnvVar{
							Name:  "LOG_VOLUME_ADDR_PLAYER2",
							Value: "/match_logs/player2_log",
						},
					},
					Command: []string{
						"./main",
					},
				},
				// player1 container
				apiv1.Container{
					Name:            "gomoku-game-player1",
					Image:           player1ImageAddr,
					ImagePullPolicy: "Never",
					Ports: []apiv1.ContainerPort{
						apiv1.ContainerPort{
							Name:          "player1-port",
							ContainerPort: 10001,
						},
					},
					VolumeMounts: []apiv1.VolumeMount{
						apiv1.VolumeMount{
							Name:      "match-logs",
							MountPath: "/match_logs",
						},
					},
					Resources: r,
					Env: []apiv1.EnvVar{
						apiv1.EnvVar{
							Name:  "LOG_VOLUME_ADDR",
							Value: "/match_logs/player1_log",
						},
						apiv1.EnvVar{
							Name:  "PORT",
							Value: "10001",
						},
					},
					Command: []string{
						"./main",
					},
				},
				// player2 container
				apiv1.Container{
					Name:            "gomoku-game-player2",
					Image:           player2ImageAddr,
					ImagePullPolicy: "Never",
					Ports: []apiv1.ContainerPort{
						apiv1.ContainerPort{
							Name:          "player2-port",
							ContainerPort: 10002,
						},
					},
					VolumeMounts: []apiv1.VolumeMount{
						apiv1.VolumeMount{
							Name:      "match-logs",
							MountPath: "/match_logs",
						},
					},
					Resources: r,
					Env: []apiv1.EnvVar{
						apiv1.EnvVar{
							Name:  "LOG_VOLUME_ADDR",
							Value: "/match_logs/player2_log",
						},
						apiv1.EnvVar{
							Name:  "PORT",
							Value: "10002",
						},
					},
					Command: []string{
						"./main",
					},
				},
			},
		},
	}

	return matchPod
}
