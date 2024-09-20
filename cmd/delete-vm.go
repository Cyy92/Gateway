package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//v1beta1 "k8s.io/api/extensions/v1beta1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"kubevirt.io/client-go/kubecli"
)

func DeleteVM(VMName, VMNamespace string, clientset *kubernetes.Clientset) error {

	if err := ValidateServiceName(VMName); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		log.Fatalf("cannot obtain KubeVirt client: %v\n", err)
	}

	getOpts := metav1.GetOptions{}

	vm, findVMErr := virtClient.VirtualMachine(VMNamespace).Get(VMName, &getOpts)

	if findVMErr != nil {
		if errors.IsNotFound(findVMErr) {
			return status.Error(codes.NotFound, findVMErr.Error())
		}
		return status.Error(codes.Internal, findVMErr.Error())
	}

	if vm != nil {
		log.Println("Deleted VM ...")
		if _, found := vm.Spec.Template.ObjectMeta.Labels["faas_vm"]; found {
			if err := deleteVM(VMName, VMNamespace, clientset); err != nil {
				return err
			}
		} else {
			return status.Error(codes.Internal, "Not a VM: "+VMName)
		}
	}

	return nil
}

func deleteVM(VMName, VMNamespace string, clientset *kubernetes.Clientset) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}
	getOpts := metav1.GetOptions{}
	ctx := context.Background()

	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		log.Fatalf("cannot obtain KubeVirt client: %v\n", err)
	}

	if vmErr := virtClient.VirtualMachine(VMNamespace).Delete(VMName, &opts); vmErr != nil {
		if errors.IsNotFound(vmErr) {
			return status.Error(codes.NotFound, vmErr.Error())

		}
		return status.Error(codes.Internal, vmErr.Error())
	}

	for {
		result, err := virtClient.VirtualMachineInstance(VMNamespace).Get(ctx, VMName, &getOpts)
		if err != nil {
			log.Println(status.Error(codes.Internal, err.Error()))
			break
		}

		if result == nil {
			break
		}

		time.Sleep(time.Second)
	}

	if svcErr := clientset.CoreV1().
		Services(VMNamespace).
		Delete(ctx, VMName, opts); svcErr != nil {

		if errors.IsNotFound(svcErr) {
			return status.Error(codes.NotFound, svcErr.Error())
		}
		return status.Error(codes.Internal, svcErr.Error())
	}

	return nil
}
