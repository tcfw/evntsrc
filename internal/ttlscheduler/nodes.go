package ttlscheduler

import (
	"fmt"

	"github.com/spf13/viper"
	pb "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type basicNodeFetcher struct {
	k8s *kubernetes.Clientset
}

func (bnf *basicNodeFetcher) GetNodes() ([]*pb.Node, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	bnf.k8s, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	workerSelector := viper.GetString("workerSelector")

	pods, err := bnf.k8s.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: workerSelector})
	if err != nil {
		return nil, fmt.Errorf("k8s: %s", err.Error())
	}

	nodes := []*pb.Node{}

	for _, pod := range pods.Items {
		nodes = append(nodes, &pb.Node{Id: pod.ObjectMeta.Name})
	}

	return nodes, nil
}
