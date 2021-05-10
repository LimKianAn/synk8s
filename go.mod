module github.com/LimKianAn/syncrd

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
	repo-url v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.8.3
)

replace repo-url => github.com/fi-ts/cloud-gateway-controller v0.0.0-20210510085353-3ef1d442031d
