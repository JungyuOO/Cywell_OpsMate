package postgres

import (
	"fmt"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultImage   = "postgres:16-alpine"
	PGVectorImage  = "pgvector/pgvector:pg16"
	MigrationImage = "ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.38"
	PortName       = "postgres"
	Port           = int32(5432)
	DefaultDBName  = "opsmate"
	DataMountPath  = "/var/lib/postgresql/data"
	RunMountPath   = "/var/run/postgresql"
	PGDataPath     = "/var/lib/postgresql/data/pgdata"
)

func Deployment(config *opsmatev1alpha1.OpsMateConfig) *appsv1.Deployment {
	labels := labelsFor(config)
	replicas := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ResourceName(config),
			Namespace: config.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "postgres",
							Image: imageFor(config),
							Ports: []corev1.ContainerPort{
								{Name: PortName, ContainerPort: Port},
							},
							Env: []corev1.EnvVar{
								{Name: "POSTGRES_DB", Value: DefaultDBName},
								{Name: "POSTGRES_USER", Value: DefaultDBName},
								{Name: "PGDATA", Value: PGDataPath},
								{
									Name: "POSTGRES_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: SecretName(config)},
											Key:                  "password",
										},
									},
								},
								{Name: "POSTGRES_SHARED_BUFFERS", Value: config.Spec.Database.SharedBuffers},
								{Name: "POSTGRES_MAX_CONNECTIONS", Value: fmt.Sprintf("%d", config.Spec.Database.MaxConnections)},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "postgres-data", MountPath: DataMountPath},
								{Name: "postgres-run", MountPath: RunMountPath},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: ptr(false),
								RunAsNonRoot:             ptr(true),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: ptr(true),
					},
					Volumes: []corev1.Volume{
						{Name: "postgres-data", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
						{Name: "postgres-run", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
					},
				},
			},
		},
	}
}

func imageFor(config *opsmatev1alpha1.OpsMateConfig) string {
	if config.Spec.Database.Image != "" {
		return config.Spec.Database.Image
	}
	if config.Spec.Embedding.RequirePGVector {
		return PGVectorImage
	}
	return DefaultImage
}

func Service(config *opsmatev1alpha1.OpsMateConfig) *corev1.Service {
	labels := labelsFor(config)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ResourceName(config),
			Namespace: config.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       PortName,
					Port:       Port,
					TargetPort: intstr.FromString(PortName),
				},
			},
		},
	}
}

func PGVectorMigrationJob(config *opsmatev1alpha1.OpsMateConfig) (*batchv1.Job, error) {
	if !config.Spec.Database.PGVectorMigrationApproved {
		return nil, fmt.Errorf("pgvector migration approval is required")
	}
	if config.Spec.Database.DSNSecretRef == "" {
		return nil, fmt.Errorf("postgres dsn secret reference is required")
	}
	if config.Spec.Embedding.Dimensions <= 0 {
		return nil, fmt.Errorf("embedding dimensions are required")
	}

	labels := labelsFor(config)
	labels["app.kubernetes.io/component"] = "pgvector-migration"
	backoffLimit := int32(0)
	dsnKey := config.Spec.Database.DSNSecretKey
	if dsnKey == "" {
		dsnKey = "dsn"
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PGVectorMigrationJobName(config),
			Namespace: config.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "pgvector-migration",
							Image:   MigrationImage,
							Command: []string{"cyops-pgvector-migrate"},
							Env: []corev1.EnvVar{
								{
									Name: "CYOPS_POSTGRES_DSN",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: config.Spec.Database.DSNSecretRef},
											Key:                  dsnKey,
										},
									},
								},
								{Name: "CYOPS_EMBEDDING_DIMENSIONS", Value: fmt.Sprintf("%d", config.Spec.Embedding.Dimensions)},
							},
						},
					},
				},
			},
		},
	}, nil
}

func ResourceName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-postgres", config.Name)
}

func SecretName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-postgres-credentials", config.Name)
}

func PGVectorMigrationJobName(config *opsmatev1alpha1.OpsMateConfig) string {
	return fmt.Sprintf("%s-pgvector-migration", config.Name)
}

func ApplyPGVectorMigrationJobStatus(config *opsmatev1alpha1.OpsMateConfig, job *batchv1.Job) {
	if job == nil {
		return
	}
	for _, condition := range job.Status.Conditions {
		switch condition.Type {
		case batchv1.JobComplete:
			if condition.Status == corev1.ConditionTrue {
				config.Status.PGVectorReady = true
				config.Status.PGVectorLastError = ""
				return
			}
		case batchv1.JobFailed:
			if condition.Status == corev1.ConditionTrue {
				config.Status.PGVectorReady = false
				config.Status.PGVectorLastError = "pgvector migration job failed"
			}
		}
	}
}

func labelsFor(config *opsmatev1alpha1.OpsMateConfig) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "cywell-opsmate",
		"app.kubernetes.io/component":  "postgres",
		"app.kubernetes.io/instance":   config.Name,
		"app.kubernetes.io/managed-by": "cywell-opsmate-operator",
	}
}

func ptr[T any](value T) *T {
	return &value
}

// Package postgres will manage the PostgreSQL database used by Lightspeed-style
// conversation cache and later OpsMate operational history storage.
