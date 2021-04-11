package pod

import (
	cachev1alpha1 "github.com/Sher-Chowdhury/mysql-operator2/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *cachev1alpha1.MySQL) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	mysqlEnvVars := cr.Spec.Environment
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
					Env: []corev1.EnvVar{
						{
							Name:  "MYSQL_ROOT_PASSWORD",
							Value: mysqlEnvVars.MysqlRootPassword,
						},
						{
							Name:  "MYSQL_DATABASE",
							Value: mysqlEnvVars.MysqlDatabase,
						},
						{
							Name:  "MYSQL_USER",
							Value: mysqlEnvVars.MysqlUser,
						},
						{
							Name:  "MYSQL_PASSWORD",
							Value: mysqlEnvVars.MysqlPassword,
						},
					},
				},
			},
		},
	}
}