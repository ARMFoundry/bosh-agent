package infrastructure_test

import (
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-agent/infrastructure"
	boshdpresolv "github.com/cloudfoundry/bosh-agent/infrastructure/devicepathresolver"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	fakeplatform "github.com/cloudfoundry/bosh-agent/platform/fakes"
)

var _ = Describe("Provider", func() {
	var (
		logger   boshlog.Logger
		platform *fakeplatform.FakePlatform
		provider Provider
	)

	BeforeEach(func() {
		platform = fakeplatform.NewFakePlatform()
		logger = boshlog.NewLogger(boshlog.LevelNone)

		providerOptions := ProviderOptions{
			MetadataService: MetadataServiceOptions{
				UseConfigDrive: true,
			},
		}

		provider = NewProvider(logger, platform, providerOptions)
	})

	Describe("Get", func() {
		It("returns aws infrastructure", func() {
			resolver := NewRegistryEndpointResolver(
				NewDigDNSResolver(logger),
			)

			metadataService := NewAwsMetadataServiceProvider(resolver).Get()
			registry := NewAwsRegistry(metadataService)

			expectedDevicePathResolver := boshdpresolv.NewMappedDevicePathResolver(
				500*time.Millisecond,
				platform.GetFs(),
			)

			expectedInf := NewAwsInfrastructure(
				metadataService,
				registry,
				platform,
				expectedDevicePathResolver,
				logger,
			)

			inf, err := provider.Get("aws")
			Expect(err).ToNot(HaveOccurred())
			Expect(inf).To(Equal(expectedInf))
		})

		It("returns openstack infrastructure", func() {
			resolver := NewRegistryEndpointResolver(
				NewDigDNSResolver(logger),
			)

			metadataServiceOptions := MetadataServiceOptions{
				UseConfigDrive: true,
			}

			metadataService := NewOpenstackMetadataServiceProvider(resolver, platform, metadataServiceOptions, logger).Get()
			registry := NewOpenstackRegistry(metadataService)

			expectedDevicePathResolver := boshdpresolv.NewMappedDevicePathResolver(
				500*time.Millisecond,
				platform.GetFs(),
			)

			expectedInf := NewOpenstackInfrastructure(
				metadataService,
				registry,
				platform,
				expectedDevicePathResolver,
				logger,
			)

			inf, err := provider.Get("openstack")
			Expect(err).ToNot(HaveOccurred())
			Expect(inf).To(Equal(expectedInf))
		})

		It("returns vsphere infrastructure", func() {
			expectedDevicePathResolver := boshdpresolv.NewVsphereDevicePathResolver(
				500*time.Millisecond,
				platform.GetFs(),
			)

			expectedInf := NewVsphereInfrastructure(platform, expectedDevicePathResolver, logger)

			inf, err := provider.Get("vsphere")
			Expect(err).ToNot(HaveOccurred())
			Expect(inf).To(Equal(expectedInf))
		})

		It("returns dummy infrastructure", func() {
			expectedDevicePathResolver := boshdpresolv.NewDummyDevicePathResolver()

			expectedInf := NewDummyInfrastructure(
				platform.GetFs(),
				platform.GetDirProvider(),
				platform,
				expectedDevicePathResolver,
			)

			inf, err := provider.Get("dummy")
			Expect(err).ToNot(HaveOccurred())
			Expect(inf).To(Equal(expectedInf))
		})

		It("returns warden infrastructure", func() {
			expectedDevicePathResolver := boshdpresolv.NewDummyDevicePathResolver()
			fs := platform.GetFs()
			boshDir := platform.GetDirProvider().BoshDir()

			wardenMetadataService := NewFileMetadataService(
				filepath.Join(boshDir, "warden-cpi-user-data.json"),
				fs,
				logger,
			)
			expectedRegistryProvider := NewRegistryProvider(
				wardenMetadataService,
				filepath.Join(boshDir, "warden-cpi-agent-env.json"),
				fs,
			)

			expectedInf := NewWardenInfrastructure(
				platform,
				expectedDevicePathResolver,
				expectedRegistryProvider,
			)

			inf, err := provider.Get("warden")
			Expect(err).ToNot(HaveOccurred())
			Expect(inf).To(Equal(expectedInf))
		})

		It("returns an error on unknown infrastructure", func() {
			_, err := provider.Get("some unknown infrastructure name")
			Expect(err).To(HaveOccurred())
		})
	})
})
