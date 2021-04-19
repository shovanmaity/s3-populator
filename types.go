package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type S3Populator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   S3PopulatorSpec   `json:"spec"`
	Status S3PopulatorStatus `json:"status"`
}

type S3PopulatorSpec struct {
	URL    string `json:"url"`
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Region string `json:"region"`
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type S3PopulatorStatus struct {
}
