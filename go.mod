module github.com/LimKianAn/sync-crd

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
	repo-url v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.8.3
)

replace repo-url => REPO_URL REPO_VERSION
