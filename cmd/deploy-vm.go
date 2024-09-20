package cmd

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/keti-openfx/openfx/pb"

	//"HV/faas/pb"

	"github.com/spf13/pflag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "kubevirt.io/api/core/v1"

	//v1beta1 "k8s.io/api/extensions/v1beta1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"kubevirt.io/client-go/kubecli"
)

/* initialReplicasCount how many replicas to start of creating for a function */
//const initialReplicasCount = 1

/* initialCpuUtilization limit of CPU Utilization per pod */
//const initialCpuUtilization = 80

type DeployVMConfig struct {
	VMNamespace string
	//EnableHttpProbe bool
	//ImagePullPolicy string
	FxWatcherPort int
	FxMeshPort    int
	//SecretMountPath string
}

/* ValidateDeployRequest validates that the service name is valid for Kubernetes */
func ValidateVMName(vm string) error {
	/* Regex for RFC-1123 validation:
	 *	k8s.io/kubernetes/pkg/util/validation/validation.go */
	var validDNS = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	matched := validDNS.MatchString(vm)
	if matched {
		return nil
	}

	return fmt.Errorf("(%s) must be a valid DNS entry for service name", vm)
}

func VMDeploy(req *pb.CreateVMRequest, clientset *kubernetes.Clientset, config *DeployVMConfig) (string, error) {
	if err := ValidateVMName(req.Instance); err != nil {
		return "", status.Error(codes.InvalidArgument, err.Error())
	}
	log.Printf("Deploying... \ndeploy VM config:%+v\n create VM request:%+v\n", config, req)

	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		log.Fatalf("cannot obtain KubeVirt client: %v\n", err)
	}

	ctx := context.Background()
	createOpts := metav1.CreateOptions{}
	getOpts := metav1.GetOptions{}

	/*

	   persistentVolumeClaim := clientset.CoreV1().PersistentVolumeClaims(config.VMNamespace)
	   persistentVolumeClaimSpec := makePersistentVolumeClaimSpec(req)
	   _, err = persistentVolumeClaim.Create(ctx, persistentVolumeClaimSpec, createOpts)

	   if err != nil {
	           if k8sErrors.IsAlreadyExists(err) {
	                   return "", status.Error(codes.AlreadyExists, err.Error())
	           }
	           return "", status.Error(codes.Internal, err.Error())
	   }
	*/

	vmSpec, err := makeVMSpec(req, config)
	if err != nil {
		return "", status.Error(codes.InvalidArgument, err.Error())
	}

	vmdeploy := virtClient.VirtualMachine(config.VMNamespace)
	_, err = vmdeploy.Create(vmSpec)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return "", status.Error(codes.AlreadyExists, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	log.Println("Created vm - " + req.Instance)

	vmServiceSpec := makeVMServiceSpec(req, config.FxWatcherPort, config.FxMeshPort)
	vmService := clientset.CoreV1().Services(config.VMNamespace)
	_, err = vmService.Create(ctx, vmServiceSpec, createOpts)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return "", status.Error(codes.AlreadyExists, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	log.Println("Created service - " + req.Instance)

	var ipaddr string
	for {
		time.Sleep(time.Second)
		result, err := virtClient.VirtualMachineInstance(config.VMNamespace).Get(ctx, req.Instance, &getOpts)
		if err != nil {
			log.Println(status.Error(codes.Internal, err.Error()))
		}

		if result.Status.Phase == "Running" {
			ipaddr = result.Status.Interfaces[0].IP
			break
		}
	}

	log.Println(ipaddr)
	return ipaddr, nil
}

func makeVMSpec(req *pb.CreateVMRequest, config *DeployVMConfig) (*v1.VirtualMachine, error) {
	labels := map[string]string{
		"faas_vm":            req.Instance,
		"kubevirt.io/size":   "small",
		"kubevirt.io/domain": req.Domain,
	}

	var img string
	if req.Domain == "ubuntu" {
		img = "tedezed/ubuntu-container-disk:20.0"
	}

	resources, resourceErr := createVMResources(req)

	if resourceErr != nil {
		return nil, resourceErr
	}

	vmSpec := &v1.VirtualMachine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualMachine",
			APIVersion: "kubevirt.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Instance,
		},
		Spec: v1.VirtualMachineSpec{
			Running: newTrue(),
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.VirtualMachineInstanceSpec{
					Domain: v1.DomainSpec{
						Devices: v1.Devices{
							Disks: []v1.Disk{
								{
									Name: "containerdisk",
									DiskDevice: v1.DiskDevice{
										Disk: &v1.DiskTarget{
											Bus: v1.DiskBusVirtio,
										},
									},
								},
								{
									Name: "cloudinitdisk",
									DiskDevice: v1.DiskDevice{
										Disk: &v1.DiskTarget{
											Bus: v1.DiskBusVirtio,
										},
									},
								},
							},
							Interfaces: []v1.Interface{
								{
									Name: "default",
									InterfaceBindingMethod: v1.InterfaceBindingMethod{
										Masquerade: &v1.InterfaceMasquerade{},
									},
								},
							},
						},
						Resources: *resources,
					},
					Networks: []v1.Network{
						{
							Name: "default",
							NetworkSource: v1.NetworkSource{
								Pod: &v1.PodNetwork{
									VMNetworkCIDR: "10.244.0.0/16",
								},
							},
						},
					},
					DNSPolicy: apiv1.DNSClusterFirstWithHostNet,
					Volumes: []v1.Volume{
						{
							Name: "containerdisk",
							VolumeSource: v1.VolumeSource{
								ContainerDisk: &v1.ContainerDiskSource{
									Image: img,
								},
							},
						},
						{
							Name: "cloudinitdisk",
							VolumeSource: v1.VolumeSource{
								CloudInitNoCloud: &v1.CloudInitNoCloudSource{
									UserDataBase64: req.UserData,
								},
							},
						},
					},
				},
			},
		},
	}

	return vmSpec, nil
}

func makeVMServiceSpec(req *pb.CreateVMRequest, fxWatcherPort int, fxMeshPort int) *apiv1.Service {
	serviceSpec := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Instance,
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"faas_vm":            req.Instance,
				"kubevirt.io/size":   "small",
				"kubevirt.io/domain": req.Domain,
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Name:     "fxwatcher",
					Port:     int32(fxWatcherPort),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(fxWatcherPort),
					},
				},
				{
					Protocol: apiv1.ProtocolTCP,
					Name:     "fxmesh",
					Port:     int32(fxMeshPort),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(fxMeshPort),
					},
				},
			},
		},
	}
	return serviceSpec
}

/*
   func makePersistentVolumeClaimSpec(req *pb.CreateFunctionRequest) *apiv1.PersistentVolumeClaim {
           resources := apiv1.ResourceRequirements{
                   Requests: apiv1.ResourceList{
			Storage:
		   },
           }

           resources.Requests[apiv1.ResourceStorage] = resource.MustParse("2Gi")

           var storageClassNamePointer *string
           var storageClassName string
           storageClassName = "nfs-client"
           storageClassNamePointer = &storageClassName

           pvClaimSpec := &apiv1.PersistentVolumeClaim{
                   TypeMeta: metav1.TypeMeta{
                           Kind: "PersistentVolumeClaim",
                           APIVersion: "v1",
                   },

                   ObjectMeta: metav1.ObjectMeta{
                           Name: req.Service,
                   },

                   Spec: apiv1.PersistentVolumeClaimSpec{
                           VolumeName: req.Service,
                           AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteMany},
                           Resources: resources,
                           StorageClassName: storageClassNamePointer,
                   },
           }
           return pvClaimSpec
   }
*/

func createVMResources(req *pb.CreateVMRequest) (*v1.ResourceRequirements, error) {
	resources := &v1.ResourceRequirements{
		Requests: apiv1.ResourceList{},
	}

	if req.Requests != nil {
		if len(req.Requests.Memory) > 0 {
			qty, err := resource.ParseQuantity(req.Requests.Memory)
			if err != nil {
				return resources, err
			}
			resources.Requests[apiv1.ResourceMemory] = qty
		}

	}

	return resources, nil
}

func newTrue() *bool {
	b := true
	return &b
}
