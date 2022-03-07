package main

import (
	flag "flag"
	log "log"
	filepath "path/filepath"
	time "time"

	client "github.com/alias-rahil/cluster-configs/pkg/client/clientset/versioned"
	infFac "github.com/alias-rahil/cluster-configs/pkg/client/informers/externalversions"
	clusterconfigmapcontroller "github.com/alias-rahil/cluster-configs/pkg/clusterconfigmapcontroller"
	clustersecretcontroller "github.com/alias-rahil/cluster-configs/pkg/clustersecretcontroller"
	rest "k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	homedir "k8s.io/client-go/util/homedir"
)

func main() {
	home := homedir.HomeDir()
	defaultKubeconfig := filepath.Join(home, ".kube", "config")
	kubeconfig := flag.String("kubeconfig", defaultKubeconfig, "(optional) absolute path to the kubeconfig file")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		log.Printf("error getting config from flags: %s\n", err.Error())
		log.Println("using in-cluster config instead")

		config, err = rest.InClusterConfig()

		if err != nil {
			log.Printf("error getting in-cluster config: %s\n", err.Error())

			return
		}
	}

	clientset, err := client.NewForConfig(config)

	if err != nil {
		log.Printf("error building clientset from config: %s\n", err.Error())

		return
	}

	sharedInformer := infFac.NewSharedInformerFactory(clientset, 20*time.Minute)
	clusterConfigMapInformer := sharedInformer.Clusterconfigs().V1alpha1().ClusterConfigMaps()
	clusterSecretInformer := sharedInformer.Clusterconfigs().V1alpha1().ClusterSecrets()
	csc := clustersecretcontroller.NewController(clientset, clusterSecretInformer)
	ccmc := clusterconfigmapcontroller.NewController(clientset, clusterConfigMapInformer)
	ch := make(chan struct{})

	sharedInformer.Start(ch)

	if err := csc.Run(ch); err != nil {
		log.Printf("error running cluster-config-map controller: %s\n", err.Error())
	}

	if err := ccmc.Run(ch); err != nil {
		log.Printf("error running cluster-secret controller: %s\n", err.Error())
	}
}
