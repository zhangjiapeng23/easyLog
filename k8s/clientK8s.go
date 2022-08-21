package k8s

import (
	"bufio"
	"context"
	"easyLog/filters"
	"fmt"
	"io"
	"os"
	"os/signal"

	"syscall"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// single instance
var (
	clientInstance = make(map[string]*Client)
)

type Client struct {
	Clientset *kubernetes.Clientset
	DataCache *DataCache
}

// global data cache
type DataCache struct {
	Namespaces *v1.NamespaceList
	Apps       map[string]*appsv1.DeploymentList
}

func NewClient(env string) *Client {
	// check exits env client instance, if exit use cache
	if _, ok := clientInstance[env]; !ok {
		//use the current  context  in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", env)
		if err != nil {
			panic(err.Error())
		}

		//create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		// cache k8s env connect
		clientInstance[env] = &Client{Clientset: clientset, DataCache: &DataCache{Apps: make(map[string]*appsv1.DeploymentList)}}
	}
	return clientInstance[env]
}

func PathExits(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (c *Client) ListNamespaces() *v1.NamespaceList {
	// check namespace whether exit cache
	if c.DataCache.Namespaces == nil {
		namespaceList, err := c.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		c.DataCache.Namespaces = namespaceList
	}
	return c.DataCache.Namespaces
}

// get all deployments info under a namespace
func (c *Client) ListAppsForNamespace(namespace string) *appsv1.DeploymentList {
	if _, ok := c.DataCache.Apps[namespace]; !ok {
		deployments, err := c.Clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		c.DataCache.Apps[namespace] = deployments
	}
	return c.DataCache.Apps[namespace]
}

// get pod info by namespace and deployment
func (c *Client) ListPodsForApp(ns string, app string) *v1.PodList {
	deployments := c.ListAppsForNamespace(ns)
	var deployment appsv1.Deployment
	for _, dp := range deployments.Items {
		if app == dp.Name {
			deployment = dp
		}
	}

	listOpt := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels).String(),
	}

	pods, err := c.Clientset.CoreV1().Pods(ns).List(context.Background(), listOpt)
	if err != nil {
		panic(err.Error())
	}

	return pods
}

// follo log output
func (c *Client) FollowLogForPods(ns string, podList *v1.PodList,
	filter func(log chan []byte, filterLog chan *filters.Log, extra ...string), extra ...string) {
	log := make(chan []byte)
	filterLog := make(chan *filters.Log)
	quit := make(chan int, 1)
	for _, pod := range podList.Items {
		go c.followLogForPods(log, quit, ns, pod.Name)
	}
	go filter(log, filterLog, extra...)
	SetupCloseHandler(quit)
	for {
		select {
		case signal := <-quit:
			if signal == 0 {
				fmt.Println("Closing log output...")
			}
			return
		case log := <-filterLog:
			log.String()
		}
	}
}

// Only print log to current time, doesn't keep output
func (c *Client) PrintLogForPods(ns string, PodList *v1.PodList,
	filter func(log chan []byte, filterLog chan *filters.Log, extra ...string), extra ...string) {
	log := make(chan []byte)
	filterLog := make(chan *filters.Log)
	quit := make(chan int, 1)
	for _, pod := range PodList.Items {
		go c.printLogForPod(log, quit, ns, pod.Name)
	}
	go filter(log, filterLog, extra...)
	SetupCloseHandler(quit)
	for {
		select {
		case signal := <-quit:
			if signal == 0 {
				fmt.Println("Closing log output...")
			}
			return
		case log := <-filterLog:
			log.String()
		}
	}
}

func (c *Client) followLogForPods(log chan []byte, quit chan int, ns string, podName string) {
	var sinceTime int64 = 60 * 60 * 2
	opts := &v1.PodLogOptions{
		Follow:       true,
		SinceSeconds: &sinceTime,
	}
	resp := c.Clientset.CoreV1().Pods(ns).GetLogs(podName, opts)
	readCloser, err := resp.Stream(context.TODO())
	if err != nil {
		// panic(err.Error())
		fmt.Fprintln(os.Stdout, err.Error())
		return
	}
	defer readCloser.Close()
	r := bufio.NewReader(readCloser)
	for {
		bytes, err := r.ReadBytes('\n')
		if err != nil {
			// panic(err.Error())
			fmt.Fprintln(os.Stderr, err.Error())
			quit <- 1
			return
		}
		log <- bytes
	}
}

func (c *Client) printLogForPod(log chan []byte, quit chan int, ns string, podName string) {
	var sinceTime int64 = 60 * 60 * 2
	opts := &v1.PodLogOptions{
		Follow:       false,
		SinceSeconds: &sinceTime,
	}

	resp := c.Clientset.CoreV1().Pods(ns).GetLogs(podName, opts)
	readCloser, err := resp.Stream(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	defer readCloser.Close()
	r := bufio.NewReader(readCloser)
	for {
		bytes, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err.Error())
				quit <- 1
			}
			return
		}
		log <- bytes
	}

}

func SetupCloseHandler(quit chan int) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		quit <- 0
	}()
}
