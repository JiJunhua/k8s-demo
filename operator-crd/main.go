package main

import (
	"context"

	crdv1 "github.com/jijunhua/k8s-demo/operator-crd/pkg/apis/crd.example.com/v1"
	clientset "github.com/jijunhua/k8s-demo/operator-crd/pkg/generated/clientset/versioned"
	"github.com/jijunhua/k8s-demo/operator-crd/pkg/generated/informers/externalversions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		klog.Fatalln(err)
		return
	}
	client, err := clientset.NewForConfig(config)
	if err != nil {
		klog.Fatalln(err)
		return
	}
	list, err := client.CrdV1().Foos("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Fatalln(err)
		return
	}

	for _, foo := range list.Items {
		klog.Infof("%+v", foo.Name)
	}

	// informer, 使用生成的externalversions
	factory := externalversions.NewSharedInformerFactory(client, 0)
	// 注册对应的实践，完成自己的定制功能
	factory.Crd().V1().Foos().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				foo := obj.(*crdv1.Foo)
				klog.Infof("add %+v", foo)
			},
			UpdateFunc: func(old, new interface{}) {
				foo := new.(*crdv1.Foo)
				klog.Infof("update %+v", foo)
			},
			DeleteFunc: func(obj interface{}) {
				foo := obj.(*crdv1.Foo)
				klog.Infof("delete %+v", foo)
			},
		},
	)

	//controller := pkg.NewController(factory.Crd().V1().Foos().Informer())
	//
	//stopCh := make(chan struct{})
	//// 启动informer
	//factory.Start(stopCh)
	//// 等待同步完成
	//factory.WaitForCacheSync(stopCh)
	//
	//// controller 处理事件消息
	//controller.Run(stopCh)
}
