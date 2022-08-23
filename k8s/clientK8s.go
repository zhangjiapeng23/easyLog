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

	"golang.org/x/term"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// single instance
var (
	clientInstance = make(map[string]*Client)
)

type Client struct {
	Clientset *kubernetes.Clientset
	DataCache *DataCache
	config    *rest.Config
}

// global data cache
type DataCache struct {
	Namespaces *corev1.NamespaceList
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
		clientInstance[env] = &Client{Clientset: clientset, DataCache: &DataCache{Apps: make(map[string]*appsv1.DeploymentList)}, config: config}
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

func (c *Client) ListNamespaces() *corev1.NamespaceList {
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
func (c *Client) ListPodsForApp(ns string, app string) *corev1.PodList {
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
func (c *Client) FollowLogForPods(ns string, podList *corev1.PodList,
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
func (c *Client) PrintLogForPods(ns string, PodList *corev1.PodList,
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
	opts := &corev1.PodLogOptions{
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
	opts := &corev1.PodLogOptions{
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

func (c *Client) ExecPod(ns string, podName string) {
	rep := c.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(ns).SubResource("exec").VersionedParams(
		&corev1.PodExecOptions{
			Command: []string{"sh", "-c", "/bin/sh"},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, _ := remotecommand.NewSPDYExecutor(c.config, "POST", rep.URL())
	fmt.Println(os.Stdin.Fd())
	fmt.Println(term.IsTerminal(88))
	fmt.Println(term.IsTerminal(1))

	// if !term.IsTerminal(0) || !term.IsTerminal(1) {
	// 	panic("stdin/stdout should be terminal")
	// 	// fmt.Errorf("stdin/stdout should be terminal")
	// 	// fmt.Println("stdin/stdout should be terminal")
	// }

	// oldState, err := term.MakeRaw(0)
	// if err != nil {
	// 	panic(err.Error())
	// }

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer term.Restore(fd, oldState)

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	if err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  screen,
		Stdout: screen,
		Stderr: screen,
		Tty:    true,
	}); err != nil {
		fmt.Print(err)
	}

}
