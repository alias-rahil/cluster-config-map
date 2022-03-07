package clusterconfigmapcontroller

import (
	errors "errors"
	log "log"
	time "time"

	clientset "github.com/alias-rahil/cluster-configs/pkg/client/clientset/versioned"
	inf "github.com/alias-rahil/cluster-configs/pkg/client/informers/externalversions/clusterconfigs.io/v1alpha1"
	lister "github.com/alias-rahil/cluster-configs/pkg/client/listers/clusterconfigs.io/v1alpha1"
	wait "k8s.io/apimachinery/pkg/util/wait"
	cache "k8s.io/client-go/tools/cache"
	workqueue "k8s.io/client-go/util/workqueue"
)

type Controller struct {
	clientset clientset.Interface
	synced    cache.InformerSynced
	lister    lister.ClusterConfigMapLister
	wq        workqueue.RateLimitingInterface
}

func NewController(clientset clientset.Interface, clusterConfigMapInformer inf.ClusterConfigMapInformer) *Controller {
	controller := &Controller{
		clientset: clientset,
		synced:    clusterConfigMapInformer.Informer().HasSynced,
		lister:    clusterConfigMapInformer.Lister(),
		wq:        workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	clusterConfigMapInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.handleAdd,
		DeleteFunc: controller.handleDelete,
		UpdateFunc: controller.handleUpdate,
	})

	return controller
}

func (controller *Controller) handleAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)

	if err != nil {
		log.Printf("cluster-config-map controller: error getting key for %s: %s\n", obj, err.Error())

		return
	}

	controller.wq.Add(key)
}

func (controller *Controller) handleDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)

	if err != nil {
		log.Printf("cluster-config-map controller: error getting key for %s: %s\n", obj, err.Error())

		return
	}

	controller.wq.Add(key)
}

func (controller *Controller) handleUpdate(old, new interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(new)

	if err != nil {
		log.Printf("cluster-config-map controller: error getting key for %s: %s\n", new, err.Error())

		return
	}

	controller.wq.Add(key)
}

func (controller *Controller) Run(ch chan struct{}) error {
	if ok := cache.WaitForCacheSync(ch, controller.synced); !ok {
		return errors.New("cluster-config-map controller: failed to sync informer cache for ClusterConfigMap controller")
	}

	go wait.Until(controller.worker, 1*time.Second, ch)

	<-ch

	return nil
}

func (controller *Controller) worker() {
	for controller.processNextItem() {
	}
}

func (controller *Controller) processNextItem() bool {
	key, quit := controller.wq.Get()

	if quit {
		return false
	}

	defer controller.wq.Done(key)

	if err := controller.processItem(key.(string)); err != nil {
		log.Printf("cluster-config-map controller: error processing %s: %s\n", key, err.Error())
		log.Printf("cluster-config-map controller: %s will be requeued\n", key)

		controller.wq.AddRateLimited(key)

		return true
	}

	controller.wq.Forget(key)

	return true
}

func (controller *Controller) processItem(key string) error {
	log.Printf("cluster-config-map controller: key: %s\n", key)

	return nil
}
