package k8s

import (
	"bufio"
	"context"
	"easyLog/filters"

	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/fatih/color"
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

var (
	// single instance
	clientInstance       = make(map[string]*Client)
	sinceTime      int64 = 2 * 60 * 60
	printYellow          = color.New(color.FgHiYellow)
	printRed             = color.New(color.FgHiRed)
)

type Client struct {
	clientset *kubernetes.Clientset
	dataCache *DataCache
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
		clientInstance[env] = &Client{clientset: clientset, dataCache: &DataCache{Apps: make(map[string]*appsv1.DeploymentList)}, config: config}
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
	if c.dataCache.Namespaces == nil {
		namespaceList, err := c.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		c.dataCache.Namespaces = namespaceList
	}
	return c.dataCache.Namespaces
}

// get all deployments info under a namespace
func (c *Client) ListAppsForNamespace(namespace string) *appsv1.DeploymentList {
	if _, ok := c.dataCache.Apps[namespace]; !ok {
		deployments, err := c.clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		c.dataCache.Apps[namespace] = deployments
	}
	return c.dataCache.Apps[namespace]
}

// get pod info by namespace and deployment
func (c *Client) ListPodsForApp(ns string, app string) *corev1.PodList {
	deployments := c.ListAppsForNamespace(ns)
	var deployment appsv1.Deployment
	for _, dp := range deployments.Items {
		if app == dp.Name {
			deployment = dp
			break
		}
	}

	listOpt := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels).String(),
	}

	pods, err := c.clientset.CoreV1().Pods(ns).List(context.Background(), listOpt)
	if err != nil {
		panic(err.Error())
	}

	return pods
}

// follo log output
func (c *Client) FollowLogForPods(ns string, podList *corev1.PodList,
	filter func(log chan []byte, filterLog chan *filters.Log, extra ...string), extra ...string) {
	// 初始化关闭广播通道
	filters.Done = make(chan struct{})
	log := make(chan []byte, 10)
	filterLog := make(chan *filters.Log, 10)
	for _, pod := range podList.Items {
		go c.followLogForPods(log, ns, pod.Name)
	}
	go filter(log, filterLog, extra...)
	// listent ctrl+c quit
	SetupCloseHandler()

	for {
		select {
		case <-filters.Done:
			printRed.Println("Closing log output...")
			return
		// from filter channel get new log and print
		case log := <-filterLog:
			log.String()
		}
	}
}

// Only print log to current time, doesn't keep output
func (c *Client) PrintLogForPods(ns string, PodList *corev1.PodList,
	filter func(log chan []byte, filterLog chan *filters.Log, extra ...string), extra ...string) {
	// 初始化关闭广播通道
	filters.Done = make(chan struct{})
	log := make(chan []byte, 10)
	filterLog := make(chan *filters.Log, 10)
	for _, pod := range PodList.Items {
		go c.printLogForPod(log, ns, pod.Name)
	}
	go filter(log, filterLog, extra...)
	SetupCloseHandler()

	for {
		select {
		case <-filters.Done:
			printRed.Println("Closing log output...")
			return
		case log := <-filterLog:
			log.String()
		}
	}
}

func (c *Client) followLogForPods(log chan []byte, ns string, podName string) {
	// 如果提前终止则直接返回
	if cancelled() {
		return
	}

	opts := &corev1.PodLogOptions{
		Follow:       true,
		SinceSeconds: &sinceTime,
	}
	resp := c.clientset.CoreV1().Pods(ns).GetLogs(podName, opts)
	readCloser, err := resp.Stream(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	defer readCloser.Close()
	r := bufio.NewReader(readCloser)
	for {
		select {
		case <-filters.Done:
			return
		default:
			bytes, err := r.ReadBytes('\n')
			if err != nil {
				if err != io.EOF && !cancelled() {
					close(filters.Done)
				}
				return
			}
			log <- bytes
		}
	}
}

func (c *Client) printLogForPod(log chan []byte, ns string, podName string) {
	// 如果提前终止则直接返回
	if cancelled() {
		return
	}

	opts := &corev1.PodLogOptions{
		Follow:       false,
		SinceSeconds: &sinceTime,
	}

	resp := c.clientset.CoreV1().Pods(ns).GetLogs(podName, opts)
	readCloser, err := resp.Stream(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	defer readCloser.Close()
	r := bufio.NewReader(readCloser)
	for {
		select {
		case <-filters.Done:
			return
		default:
			bytes, err := r.ReadBytes('\n')
			if err != nil {
				if err != io.EOF && !cancelled() {
					close(filters.Done)
				}
				return
			}
			log <- bytes
		}
	}
}

// listen Ctrl+C to termination log output
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
			close(filters.Done)
		case <-filters.Done:
			return
		}
	}()
}

func cancelled() bool {
	select {
	case <-filters.Done:
		return true
	default:
		return false
	}
}

func (c *Client) ExecPod(ns string, podName string) {
	// enter pod defualt container shell by exec 'sh -c /bin/sh'
	rep := c.clientset.CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(ns).SubResource("exec").VersionedParams(
		&corev1.PodExecOptions{
			Command: []string{"sh", "-c", "/bin/sh"},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.config, "POST", rep.URL())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	defer term.Restore(fd, oldState)

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	if err := exec.Stream(remotecommand.StreamOptions{
		Stdin:  screen,
		Stdout: screen,
		Stderr: screen,
		Tty:    true,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	printYellow.Print("\r输入回车继续...")
	screen.Read(make([]byte, 0))
	fmt.Print("\r")
	time.Sleep(time.Second * 1)
	fmt.Print("\r")
}
