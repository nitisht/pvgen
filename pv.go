/*
 * Copyright (C) 2020, MinIO, Inc.
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3,
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License, version 3,
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 *
 */

package main

import (
	"context"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	accessMode    = "ReadWriteOnce"
	reclaimPolicy = corev1.PersistentVolumeReclaimRetain
)

var fs = "xfs"

func createPVs(kubeClient kubernetes.Interface) error {
	ip, err := parseInput(inputPath)
	if err != nil {
		return err
	}
	for i, h := range ip.Hosts {
		for j, p := range ip.Paths {
			if err = createPV(kubeClient, ip.Capacity, ip.StorageClass, h, p, ip.Namespace, "minio-pv-"+strconv.Itoa(i)+strconv.Itoa(j)); err != nil {
				return err
			}
		}
	}

	return nil
}

func createPV(kubeClient kubernetes.Interface, capacity, sc, host, path, ns, name string) error {
	c, err := resource.ParseQuantity(capacity)
	if err != nil {
		return err
	}
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceName("storage"): c,
			},
			PersistentVolumeReclaimPolicy: reclaimPolicy,
			AccessModes:                   []corev1.PersistentVolumeAccessMode{accessMode},
			StorageClassName:              sc,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				Local: &corev1.LocalVolumeSource{
					Path:   path,
					FSType: &fs,
				},
			},
			NodeAffinity: &corev1.VolumeNodeAffinity{
				Required: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      "kubernetes.io/hostname",
									Operator: corev1.NodeSelectorOpIn,
									Values:   []string{host},
								},
							},
						},
					},
				},
			},
		},
	}

	ctx := context.Background()
	cOpts := metav1.CreateOptions{}
	_, err = kubeClient.CoreV1().PersistentVolumes().Create(ctx, pv, cOpts)
	return nil
}
