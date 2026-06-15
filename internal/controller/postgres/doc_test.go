package postgres

import (
	"testing"

	opsmatev1alpha1 "github.com/JungyuOO/Cywell_OpsMate/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDeploymentBuildsPostgresShape(t *testing.T) {
	config := sampleConfig()

	deployment := Deployment(config)

	if deployment.Name != "sample-postgres" {
		t.Fatalf("name = %q, want sample-postgres", deployment.Name)
	}
	container := deployment.Spec.Template.Spec.Containers[0]
	if container.Image != "registry.example.com/pgvector:pg16" {
		t.Fatalf("image = %q, want custom pgvector image", container.Image)
	}
	assertEnv(t, container.Env, "POSTGRES_DB", DefaultDBName)
	assertEnv(t, container.Env, "POSTGRES_SHARED_BUFFERS", "128MB")
	assertEnv(t, container.Env, "POSTGRES_MAX_CONNECTIONS", "100")
	assertSecretEnv(t, container.Env, "POSTGRES_PASSWORD", "sample-postgres-credentials", "password")
	if container.Ports[0].ContainerPort != Port {
		t.Fatalf("container port = %d, want %d", container.Ports[0].ContainerPort, Port)
	}
}

func TestServiceTargetsPostgresPort(t *testing.T) {
	config := sampleConfig()

	service := Service(config)

	if service.Name != "sample-postgres" {
		t.Fatalf("name = %q, want sample-postgres", service.Name)
	}
	if service.Spec.Ports[0].Port != Port {
		t.Fatalf("service port = %d, want %d", service.Spec.Ports[0].Port, Port)
	}
	if service.Spec.Selector["app.kubernetes.io/component"] != "postgres" {
		t.Fatalf("component selector = %q, want postgres", service.Spec.Selector["app.kubernetes.io/component"])
	}
}

func sampleConfig() *opsmatev1alpha1.OpsMateConfig {
	return &opsmatev1alpha1.OpsMateConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample",
			Namespace: "opsmate",
		},
		Spec: opsmatev1alpha1.OpsMateConfigSpec{
			Database: opsmatev1alpha1.DatabaseSpec{
				Type:           "postgres",
				Image:          "registry.example.com/pgvector:pg16",
				SharedBuffers:  "128MB",
				MaxConnections: 100,
			},
		},
	}
}

func TestDeploymentDefaultsToPGVectorImageWhenRequired(t *testing.T) {
	config := sampleConfig()
	config.Spec.Database.Image = ""
	config.Spec.Embedding.RequirePGVector = true

	deployment := Deployment(config)
	container := deployment.Spec.Template.Spec.Containers[0]

	if container.Image != PGVectorImage {
		t.Fatalf("image = %q, want %q", container.Image, PGVectorImage)
	}
}

func TestPGVectorMigrationJobRequiresApproval(t *testing.T) {
	config := sampleConfig()
	config.Spec.Database.PGVectorMigrationApproved = false
	config.Spec.Database.DSNSecretRef = "postgres-dsn"
	config.Spec.Embedding.Dimensions = 768

	if _, err := PGVectorMigrationJob(config); err == nil {
		t.Fatal("expected approval error")
	}
}

func TestPGVectorMigrationJobBuildsSecretBackedTemplate(t *testing.T) {
	config := sampleConfig()
	config.Spec.Database.PGVectorMigrationApproved = true
	config.Spec.Database.DSNSecretRef = "postgres-dsn"
	config.Spec.Database.DSNSecretKey = "url"
	config.Spec.Embedding.Dimensions = 768

	job, err := PGVectorMigrationJob(config)
	if err != nil {
		t.Fatal(err)
	}

	if job.Name != "sample-pgvector-migration" {
		t.Fatalf("name = %q", job.Name)
	}
	container := job.Spec.Template.Spec.Containers[0]
	if container.Command[0] != "cyops-pgvector-migrate" {
		t.Fatalf("command = %v", container.Command)
	}
	assertSecretEnv(t, container.Env, "CYOPS_POSTGRES_DSN", "postgres-dsn", "url")
	assertEnv(t, container.Env, "CYOPS_EMBEDDING_DIMENSIONS", "768")
	if container.Env[0].Value != "" {
		t.Fatal("dsn value is set directly, want SecretKeyRef only")
	}
}

func TestApplyPGVectorMigrationJobStatusMarksReadyOnCompletion(t *testing.T) {
	config := sampleConfig()
	job := &batchv1.Job{
		Status: batchv1.JobStatus{
			Conditions: []batchv1.JobCondition{
				{Type: batchv1.JobComplete, Status: corev1.ConditionTrue},
			},
		},
	}

	ApplyPGVectorMigrationJobStatus(config, job)

	if !config.Status.PGVectorReady {
		t.Fatal("pgVectorReady = false, want true")
	}
	if config.Status.PGVectorLastError != "" {
		t.Fatalf("last error = %q, want empty", config.Status.PGVectorLastError)
	}
}

func TestApplyPGVectorMigrationJobStatusRecordsFailure(t *testing.T) {
	config := sampleConfig()
	config.Status.PGVectorReady = true
	job := &batchv1.Job{
		Status: batchv1.JobStatus{
			Conditions: []batchv1.JobCondition{
				{Type: batchv1.JobFailed, Status: corev1.ConditionTrue},
			},
		},
	}

	ApplyPGVectorMigrationJobStatus(config, job)

	if config.Status.PGVectorReady {
		t.Fatal("pgVectorReady = true, want false")
	}
	if config.Status.PGVectorLastError != "pgvector migration job failed" {
		t.Fatalf("last error = %q, want migration failure", config.Status.PGVectorLastError)
	}
}

func assertEnv(t *testing.T, env []corev1.EnvVar, name string, want string) {
	t.Helper()
	for _, item := range env {
		if item.Name == name {
			if item.Value != want {
				t.Fatalf("%s = %q, want %q", name, item.Value, want)
			}
			return
		}
	}
	t.Fatalf("missing env %s", name)
}

func assertSecretEnv(t *testing.T, env []corev1.EnvVar, name string, secretName string, key string) {
	t.Helper()
	for _, item := range env {
		if item.Name != name {
			continue
		}
		if item.ValueFrom == nil || item.ValueFrom.SecretKeyRef == nil {
			t.Fatalf("%s does not use a SecretKeyRef", name)
		}
		if item.ValueFrom.SecretKeyRef.Name != secretName {
			t.Fatalf("%s secret = %q, want %q", name, item.ValueFrom.SecretKeyRef.Name, secretName)
		}
		if item.ValueFrom.SecretKeyRef.Key != key {
			t.Fatalf("%s key = %q, want %q", name, item.ValueFrom.SecretKeyRef.Key, key)
		}
		return
	}
	t.Fatalf("missing env %s", name)
}
