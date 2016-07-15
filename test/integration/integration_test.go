package integration

import (
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	k8sapi "k8s.io/kubernetes/pkg/api"
	k8serrors "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/client/unversioned/testclient"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	clientpkg "github.com/wikiwi/kube-volume-freezer/pkg/client"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/master"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/fs/fstest"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/volumes"
)

var _ = Describe("Server", func() {
	// Setup test environment.
	var k8sFixture = []runtime.Object{
		&k8sapi.Pod{
			ObjectMeta: k8sapi.ObjectMeta{Namespace: "default", Name: "pod1", UID: "11111111-1111-1111-1111-111111111111", Labels: map[string]string{"abc": "123"}},
			Status:     k8sapi.PodStatus{PodIP: "127.0.0.1"},
			Spec:       k8sapi.PodSpec{NodeName: "node1"},
		},
		&k8sapi.Pod{
			ObjectMeta: k8sapi.ObjectMeta{Namespace: "default", Name: "pod2", UID: "22222222-2222-2222-2222-222222222222", Labels: map[string]string{"abc": "123"}},
			Status:     k8sapi.PodStatus{PodIP: "127.0.0.1"},
			Spec:       k8sapi.PodSpec{NodeName: "node1"},
		},
		&k8sapi.Pod{
			ObjectMeta: k8sapi.ObjectMeta{Namespace: "default", Name: "pod3", UID: "33333333-3333-3333-3333-333333333333", Labels: map[string]string{"def": "123"}},
			Status:     k8sapi.PodStatus{PodIP: "127.0.0.1"},
			Spec:       k8sapi.PodSpec{NodeName: "node1"},
		},
	}

	var fsFixture = []string{
		volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~empty-dir/empty1",
		volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~empty-dir/empty2",
		volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~gce-pd/pd1",
		volumes.PodsBasePath + "/22222222-2222-2222-2222-222222222222/volumes/kubernetes.io~gce-pd/pd1",
	}

	var masterToken string
	var master2MinionToken string
	var minionToken string
	var minionSelector string
	var minionNamespace string
	var clientToken string

	var resetConfig = func() {
		masterToken = ""
		master2MinionToken = ""
		minionToken = ""
		minionSelector = "abc=123"
		minionNamespace = "default"
		clientToken = ""
	}
	resetConfig()

	var minionServer *httptest.Server
	var masterServer *httptest.Server
	var client clientpkg.Interface
	var fsFake *fstest.FakeFS

	JustBeforeEach(func() {
		// Setup minion.
		fsFake = fstest.NewFake(fsFixture)
		minionRESTServer, err := minion.NewRESTServer(&minion.Options{FS: fsFake, Token: minionToken})
		if err != nil {
			panic(err)
		}
		minionServer = httptest.NewServer(minionRESTServer.Handler())

		u, err := url.Parse(minionServer.URL)
		if err != nil {
			panic(err)
		}
		port, err := strconv.Atoi(strings.Split(u.Host, ":")[1])
		if err != nil {
			panic(err)
		}

		// Setup master.
		fakeClient := testclient.NewSimpleFake(k8sFixture...)
		fakeClient.PrependReactor("get", "*", func(action testclient.Action) (handled bool, ret runtime.Object, err error) {
			a := action.(testclient.GetAction)
			if a.GetName() == "not-existing" {
				return true, nil, k8serrors.NewNotFound(k8sapi.Resource("pods"), a.GetName())
			}
			return false, nil, nil
		})
		masterRESTServer, err := master.NewRESTServer(&master.Options{
			Token:           masterToken,
			MinionNamespace: minionNamespace,
			MinionSelector:  minionSelector,
			MinionPort:      port,
			MinionToken:     master2MinionToken,
			KubeClient:      fakeClient,
		})
		if err != nil {
			panic(err)
		}
		masterServer = httptest.NewServer(masterRESTServer.Handler())

		client, err = clientpkg.New(masterServer.URL, &clientpkg.Options{Token: clientToken})
		if err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		masterServer.Close()
		minionServer.Close()
		resetConfig()
	})

	Describe("GET /volumes/{namespace}/{podName}", func() {
		It("should return list of volumes when found", func() {
			volumesList, err := client.Volumes().List("default", "pod1")
			Expect(err).To(BeNil())
			Expect(volumesList).To(Equal(&api.VolumeList{
				PodUID: "11111111-1111-1111-1111-111111111111",
				Items:  []string{"empty1", "empty2", "pd1"},
			}))
		})
		It("should return error when not found", func() {
			_, err := client.Volumes().List("default", "not-existing")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(404))
		})
		It("should reject invalid parameters", func() {
			By("validating namespace")
			_, err := client.Volumes().List("$nvalid", "pod1")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(422))

			By("validating pod name")
			_, err = client.Volumes().List("default", "$nvalid")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(422))
		})
	})

	Describe("GET /volumes/{namespace}/{podName}/{volName}", func() {
		It("should return volume when found", func() {
			volumesList, err := client.Volumes().Get("default", "pod1", "empty1")
			Expect(err).To(BeNil())
			Expect(volumesList).To(Equal(&api.Volume{
				PodUID: "11111111-1111-1111-1111-111111111111",
				Name:   "empty1",
			}))
		})
		It("should return error when not found", func() {
			_, err := client.Volumes().Get("default", "pod1", "not-existing")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(404))
		})
		It("should reject invalid parameters", func() {
			By("validating namespace")
			_, err := client.Volumes().Get("$nvalid", "pod1", "empty1")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(422))

			By("validating pod name")
			_, err = client.Volumes().Get("default", "$nvalid", "empty1")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(422))

			By("validating volume name")
			_, err = client.Volumes().Get("default", "pod1", "$invalid")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(422))
		})
	})

	Describe("POST /volumes/{namespace}/{podName}/{volName}", func() {
		It("should freeze and return volume when found", func() {
			volumesList, err := client.Volumes().Freeze("default", "pod1", "empty1")
			Expect(err).To(BeNil())
			Expect(volumesList).To(Equal(&api.Volume{
				PodUID: "11111111-1111-1111-1111-111111111111",
				Name:   "empty1",
			}))
			Expect(fsFake.Frozen).To(HaveLen(1))
		})
		It("should thaw and return volume when found", func() {
			volumesList, err := client.Volumes().Thaw("default", "pod1", "empty1")
			Expect(err).To(BeNil())
			Expect(volumesList).To(Equal(&api.Volume{
				PodUID: "11111111-1111-1111-1111-111111111111",
				Name:   "empty1",
			}))
			Expect(fsFake.Thawed).To(HaveLen(1))
		})
		It("should return error when not found", func() {
			_, err := client.Volumes().Freeze("default", "pod1", "not-existing")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(404))
		})
		It("should reject invalid parameters", func() {
			By("validating namespace")
			_, err := client.Volumes().Freeze("$nvalid", "pod1", "empty1")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(422))

			By("validating pod name")
			_, err = client.Volumes().Freeze("default", "$nvalid", "empty1")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(422))

			By("validating volume name")
			_, err = client.Volumes().Freeze("default", "pod1", "$invalid")
			Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
			Expect(err.(*api.Error).Code).To(Equal(422))
		})
	})

	Describe("swagger", func() {
		Context("master", func() {
			It("should return swagger spec at /apidocs.json", func() {
				var json interface{}
				client := generic.NewOrDie(masterServer.URL, nil)
				req := client.NewRequestOrDie("GET", "/apidocs.json", nil)
				_, err := client.Do(req, &json)
				Expect(err).To(BeNil())
				Expect(json).To(HaveKey("swaggerVersion"))
			})
		})
		Context("minion", func() {
			It("should return swagger spec at /apidocs.json", func() {
				var json interface{}
				client := generic.NewOrDie(minionServer.URL, nil)
				req := client.NewRequestOrDie("GET", "/apidocs.json", nil)
				_, err := client.Do(req, &json)
				Expect(err).To(BeNil())
				Expect(json).To(HaveKey("swaggerVersion"))
			})
		})
	})

	Describe("health", func() {
		Context("master", func() {
			It("should report health status at /healthz", func() {
				var health api.Health
				client := generic.NewOrDie(masterServer.URL, nil)
				req := client.NewRequestOrDie("GET", "/healthz", nil)
				_, err := client.Do(req, &health)
				Expect(err).To(BeNil())
				Expect(health.Status).To(Equal("healthy"))
			})
		})
		Context("minion", func() {
			It("should report health status at /healthz", func() {
				var health api.Health
				client := generic.NewOrDie(minionServer.URL, nil)
				req := client.NewRequestOrDie("GET", "/healthz", nil)
				_, err := client.Do(req, &health)
				Expect(err).To(BeNil())
				Expect(health.Status).To(Equal("healthy"))
			})
		})
	})

	Describe("auth", func() {
		Context("with auth turned on", func() {
			BeforeEach(func() {
				masterToken = "protectedMaster"
				minionToken = "protectedMinion"
				master2MinionToken = minionToken
			})
			Context("without credentials", func() {
				It("should block unauthenticated requests to resources", func() {
					_, err := client.Volumes().List("default", "pod1")
					Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
					Expect(err.(*api.Error).Code).To(Equal(403))
				})
				It("should allow health requests to master", func() {
					client := generic.NewOrDie(masterServer.URL, nil)
					req := client.NewRequestOrDie("GET", "/healthz", nil)
					_, err := client.Do(req, nil)
					Expect(err).To(BeNil())
				})
				It("should allow health requests to minion", func() {
					client := generic.NewOrDie(minionServer.URL, nil)
					req := client.NewRequestOrDie("GET", "/healthz", nil)
					_, err := client.Do(req, nil)
					Expect(err).To(BeNil())
				})
				It("should allow swagger requests to master", func() {
					client := generic.NewOrDie(masterServer.URL, nil)
					req := client.NewRequestOrDie("GET", "/apidocs.json", nil)
					_, err := client.Do(req, nil)
					Expect(err).To(BeNil())
				})
				It("should allow swagger requests to minion", func() {
					client := generic.NewOrDie(minionServer.URL, nil)
					req := client.NewRequestOrDie("GET", "/apidocs.json", nil)
					_, err := client.Do(req, nil)
					Expect(err).To(BeNil())
				})
			})
			Context("with credentials", func() {
				BeforeEach(func() {
					clientToken = "protectedMaster"
				})
				It("should allow authenticated requests", func() {
					_, err := client.Volumes().List("default", "pod1")
					Expect(err).To(BeNil())
				})
			})
		})
		Context("with only minion auth turned on", func() {
			BeforeEach(func() {
				minionToken = "protectedMinion"
				master2MinionToken = minionToken
			})
			It("should allow all requests", func() {
				_, err := client.Volumes().List("default", "pod1")
				Expect(err).To(BeNil())
			})
		})
		Context("with auth turned on but wrong master2MinionToken ", func() {
			BeforeEach(func() {
				minionToken = "protectedMinion"
				master2MinionToken = "wrongToken"
			})
			It("should return 403", func() {
				_, err := client.Volumes().List("default", "pod1")
				Expect(err).To(BeAssignableToTypeOf(&api.Error{}))
				Expect(err.(*api.Error).Code).To(Equal(403))
			})
		})
	})
})
