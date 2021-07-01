package main

import (
	"context"
	"fmt"
	"os"

	"github.com/whywaita/ingress-external-caddy/cmd"
	"github.com/whywaita/ingress-external-caddy/driver"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

var (
	o = &cmd.Options{}
)

func main() {
	if err := run(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

func run() error {
	c := cmd.New(o)
	if err := c.Execute(); err != nil {
		return fmt.Errorf("failed to parse flag")
	}
	if _, err := driver.GetProvider(o.Provider); err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	clientset, err := driver.NewClient(o.KubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to new client: %w", err)
	}

	ingressWatcher, err := driver.NewWatcher(clientset)
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	_, informer := cache.NewInformer(
		ingressWatcher,
		&networkingv1.Ingress{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				addFunc(obj, o, clientset)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				updateFunc(oldObj, newObj, o, clientset)
			},
			DeleteFunc: func(obj interface{}) {
				deleteFunc(obj, o, clientset)
			},
		},
	)

	stop := make(chan struct{})
	defer close(stop)
	klog.Info("Starting Ingress informer...")
	go informer.Run(stop)

	select {}
}

func updateCaddyConfig(o *cmd.Options, clientset kubernetes.Interface) error {
	ingresses, err := clientset.NetworkingV1().Ingresses(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list ingresses: %w", err)
	}

	conf, err := driver.GenerateCaddy(ingresses.Items, *o)
	if err != nil {
		return fmt.Errorf("failed to generate config of caddy: %w", err)
	}

	if err := driver.UpdateConfig(*conf, o.CaddyHost); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

func addFunc(obj interface{}, o *cmd.Options, clientset kubernetes.Interface) {
	if _, ok := obj.(*networkingv1.Ingress); ok {
		if err := updateCaddyConfig(o, clientset); err != nil {
			klog.Error(err)
			return
		}
	} else {
		if _, err := cache.MetaNamespaceKeyFunc(obj); err != nil {
			klog.Error(err)
			return
		}
	}
}

func updateFunc(old, new interface{}, o *cmd.Options, clientset kubernetes.Interface) {
	_, okNew := new.(*networkingv1.Ingress)
	_, okOld := old.(*networkingv1.Ingress)
	if okNew && okOld {
		if err := updateCaddyConfig(o, clientset); err != nil {
			klog.Error(err)
			return
		}
	} else {
		if _, err := cache.MetaNamespaceKeyFunc(new); err != nil {
			klog.Error(err)
			return
		}
	}
}

func deleteFunc(obj interface{}, o *cmd.Options, clientset kubernetes.Interface) {
	if _, ok := obj.(*networkingv1.Ingress); ok {
		if err := updateCaddyConfig(o, clientset); err != nil {
			klog.Error(err)
			return
		}
	} else {
		if _, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err != nil {
			klog.Error(err)
			return
		}
	}
}
