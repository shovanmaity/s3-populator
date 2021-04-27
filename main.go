package main

import (
	"flag"
	"os"

	populator_machinery "github.com/kubernetes-csi/lib-volume-populator/populator-machinery"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	prefix     = "example.io"
	mountPath  = "/mnt"
	devicePath = "/dev/block"
)

func main() {
	var (
		imageName string
	)
	flag.StringVar(&imageName, "image-name", "", "Image to use for populating")
	flag.Parse()

	namespace := os.Getenv("POD_NAMESPACE")
	const (
		groupName  = "example.io"
		apiVersion = "v1"
		kind       = "S3Populator"
		resource   = "s3populators"
	)
	var (
		gk  = schema.GroupKind{Group: groupName, Kind: kind}
		gvr = schema.GroupVersionResource{Group: groupName, Version: apiVersion, Resource: resource}
	)
	populator_machinery.RunController("", "", imageName,
		namespace, prefix, gk, gvr, mountPath, devicePath, getPopulatorPodArgs)
}

func getPopulatorPodArgs(rawBlock bool, u *unstructured.Unstructured) ([]string, error) {
	if rawBlock {
		return nil, errors.New("Block volume is not imaplementated in populator")
	}
	s3populator := S3Populator{}
	err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(u.UnstructuredContent(), &s3populator)
	if nil != err {
		return nil, err
	}
	args := []string{}
	if s3populator.Spec.URL != "" {
		args = append(args, "-u "+s3populator.Spec.URL)
	}
	if s3populator.Spec.ID != "" {
		args = append(args, "-i "+s3populator.Spec.ID)
	}
	if s3populator.Spec.Secret != "" {
		args = append(args, "-s "+s3populator.Spec.Secret)
	}
	if s3populator.Spec.Region != "" {
		args = append(args, "-r "+s3populator.Spec.Region)
	}
	if s3populator.Spec.Bucket != "" {
		args = append(args, "-b "+s3populator.Spec.Bucket)
	}
	if s3populator.Spec.Key != "" {
		args = append(args, "-k "+s3populator.Spec.Key)
	}
	args = append(args, "-p "+mountPath)
	return args, nil
}
