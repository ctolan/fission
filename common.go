/*
Copyright 2016 The Fission Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fission

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/gorilla/handlers"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	rbac "k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

func UrlForFunction(name, namespace string) string {
	prefix := "/fission-function"
	if namespace != metav1.NamespaceDefault {
		prefix += "/" + namespace
	}
	return fmt.Sprintf("%v/%v", prefix, name)
}

func SetupStackTraceHandler() {
	// register signal handler for dumping stack trace.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Received SIGTERM : Dumping stack trace")
		debug.PrintStack()
		os.Exit(1)
	}()
}

// IsNetworkError returns true if an error is a network error, and false otherwise.
func IsNetworkError(err error) bool {
	_, ok := err.(net.Error)
	return ok
}

// GetFunctionIstioServiceName return service name of function for istio feature
func GetFunctionIstioServiceName(fnName, fnNamespace string) string {
	return fmt.Sprintf("istio-%v-%v", fnName, fnNamespace)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI := r.RequestURI
		if !strings.Contains(requestURI, "healthz") {
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			handlers.LoggingHandler(os.Stdout, next).ServeHTTP(w, r)
		}
	})
}

// following set of utilities are for setting up cluster role bindings for fission env deployments.
func makeClusterRoleBindingObj(ns, sa, clusterRoleBinding, clusterRole string) *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBinding,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      sa,
				Namespace: ns,
			},
		},
		RoleRef: rbac.RoleRef{
			Kind: "ClusterRole",
			Name: clusterRole,
		},
	}
}

func addSAToClusterRoleBinding(k8sClient *kubernetes.Clientset, crbObj *rbac.ClusterRoleBinding, sa, ns string) error {
	subjects := crbObj.Subjects
	subjects = append(subjects, rbac.Subject{
		Kind:      "ServiceAccount",
		Name:      sa,
		Namespace: ns,
	})
	crbObj.Subjects = subjects

	_, err := k8sClient.RbacV1beta1().ClusterRoleBindings().Update(crbObj)
	return err
}

func isSAInClusterRoleBinding(crbObj *rbac.ClusterRoleBinding, sa, ns string) bool {
	for _, subject := range crbObj.Subjects {
		if subject.Name == sa && subject.Namespace == ns {
			return true
		}
	}

	return false
}

func setupClusterRoleBinding(k8sClient *kubernetes.Clientset, sa, ns, clusterRoleBinding, clusterRole string) error {
	// get the cluster role binding object
	crbObj, err := k8sClient.RbacV1beta1().ClusterRoleBindings().Get(
		clusterRoleBinding, metav1.GetOptions{})

	if err == nil {
		// if cluster role binding exists, check if this sa is part of the binding. if not, add it
		if !isSAInClusterRoleBinding(crbObj, sa, ns) {
			return addSAToClusterRoleBinding(k8sClient, crbObj, sa, ns)
		}
		return nil
	}

	// if cluster role binding is missing, create it. also add this sa to the binding.
	if k8serrors.IsNotFound(err) {
		crbObj = makeClusterRoleBindingObj(ns, sa, clusterRoleBinding, clusterRole)
		crbObj, err = k8sClient.RbacV1beta1().ClusterRoleBindings().Create(crbObj)
	}

	return err
}

func makeSAObj(sa, ns string) *apiv1.ServiceAccount {
	return &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      sa,
		},
	}
}

func setupSA(k8sClient *kubernetes.Clientset, sa, ns string) (*apiv1.ServiceAccount, error) {
	saObj, err := k8sClient.CoreV1Client.ServiceAccounts(ns).Get(sa, metav1.GetOptions{})
	if err == nil {
		return saObj, nil
	}

	if k8serrors.IsNotFound(err) {
		saObj = makeSAObj(sa, ns)
		saObj, err = k8sClient.CoreV1Client.ServiceAccounts(ns).Create(saObj)
	}

	return saObj, err
}

func SetupRBAC(k8sClient *kubernetes.Clientset, sa, ns, clusterRoleBinding, clusterRole string) error {
	_, err := setupSA(k8sClient, sa, ns)
	if err != nil {
		return err
	}

	return setupClusterRoleBinding(k8sClient, sa, ns, clusterRoleBinding, clusterRole)
}
