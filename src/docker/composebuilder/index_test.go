package composebuilder_test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/util"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("composebuilder", func() {
	var _ = Describe("GetApplicationDockerCompose", func() {
		It("should return the proper docker configs for development", func() {
			appDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			err = helpers.CheckoutApp(appDir, "complex-setup-app")
			Expect(err).NotTo(HaveOccurred())
			internalServices := []string{"html-server", "todo-service", "users-service"}
			externalServices := []string{"external-service"}
			internalDependencies := []string{"exocom"}
			externalDependencies := []string{"mongo"}
			appContext, err := context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())

			dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
				AppContext: appContext,
				BuildMode: types.BuildMode{
					Type:        types.BuildModeTypeLocal,
					Environment: types.BuildModeEnvironmentDevelopment,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			By("generate an image name for each dependency and external service")
			for _, serviceRole := range helpers.JoinStringSlices(internalDependencies, externalDependencies, externalServices) {
				Expect(len(dockerCompose.Services[serviceRole].Image)).ToNot(Equal(0))
			}

			By("should have the correct build command for each internal service and dependency")
			for _, serviceRole := range internalServices {
				Expect(dockerCompose.Services[serviceRole].Command).To(Equal(`echo "does not run"`))
			}
			Expect(dockerCompose.Services["exocom"].Command).To(Equal(""))

			By("should include 'exocom' in the dependencies of every service")
			for _, serviceRole := range append(internalServices, externalServices...) {
				exists := util.DoesStringArrayContain(dockerCompose.Services[serviceRole].DependsOn, "exocom")
				Expect(exists).To(Equal(true))
			}

			By("should include external dependencies as dependencies")
			exists := util.DoesStringArrayContain(dockerCompose.Services["todo-service"].DependsOn, "mongo")
			Expect(exists).To(Equal(true))

			By("should properly reserve ports for services")
			actualApiPort := dockerCompose.Services["api-service"].Ports
			expectedApiPort := []string{"3000:80"}
			Expect(actualApiPort).To(Equal(expectedApiPort))

			actualExternalServicePort := dockerCompose.Services["external-service"].Ports
			expectedExternalServicePort := []string{"3010:5000"}
			Expect(actualExternalServicePort).To(Equal(expectedExternalServicePort))

			actualHtmlPort := dockerCompose.Services["html-server"].Ports
			expectedHtmlPort := []string{"3020:80"}
			Expect(actualHtmlPort).To(Equal(expectedHtmlPort))

			By("should inject the proper service endpoint environment variables")
			expectedApiEndpointKey := "API_SERVICE_EXTERNAL_ORIGIN"
			expectedApiEndpointValue := "http://localhost:3000"
			expectedHtmlEndpointKey := "HTML_SERVER_EXTERNAL_ORIGIN"
			expectedHtmlEndpointValue := "http://localhost:3020"

			skipServices := []string{"api-service", "exocom", "mongo"}
			for serviceRole, dockerConfig := range dockerCompose.Services {
				if util.DoesStringArrayContain(skipServices, serviceRole) {
					continue
				}
				Expect(dockerConfig.Environment[expectedApiEndpointKey]).To(Equal(expectedApiEndpointValue))
			}
			skipServices = []string{"html-server", "exocom", "mongo"}
			for serviceRole, dockerConfig := range dockerCompose.Services {
				if util.DoesStringArrayContain(skipServices, serviceRole) {
					continue
				}
				Expect(dockerConfig.Environment[expectedHtmlEndpointKey]).To(Equal(expectedHtmlEndpointValue))
			}
			nonPublicServiceKeys := []string{
				"USERS_SERVICE_EXTERNAL_ORIGIN",
				"TODO_SERVICE_EXTERNAL_ORIGIN",
			}
			for _, dockerConfig := range dockerCompose.Services {
				for _, nonPublicKey := range nonPublicServiceKeys {
					Expect(dockerConfig.Environment[nonPublicKey]).To(Equal(""))
				}
			}

			By("should include the correct exocom environment variables")
			environment := dockerCompose.Services["exocom"].Environment
			serviceData, err := json.Marshal(map[string]map[string]interface{}{
				"api-service":      {},
				"external-service": {},
				"html-server": {
					"receives": []interface{}{"todo.created"},
					"sends":    []interface{}{"todo.create"},
				},
				"todo-service": {
					"receives": []interface{}{"todo.create"},
					"sends":    []interface{}{"todo.created"},
				},
				"users-service": {
					"receives": []interface{}{"mongo.list", "mongo.create"},
					"sends":    []interface{}{"mongo.listed", "mongo.created"},
					"translations": []interface{}{
						map[string]interface{}{
							"internal": "mongo create",
							"public":   "users create",
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(environment).To(Equal(map[string]string{
				"SERVICE_DATA": string(serviceData),
			}))

			By("should include exocom environment variables in every services' environment")
			for _, serviceRole := range append(internalServices, externalServices...) {
				environment := dockerCompose.Services[serviceRole].Environment
				Expect(environment["EXOCOM_HOST"]).To(Equal("exocom"))
			}

			By("should generate a volume path for an external dependency that mounts a volume")
			Expect(len(dockerCompose.Services["mongo"].Volumes)).NotTo(Equal(0))

			By("should have the specified image for the external service")
			serviceRole := "external-service"
			imageName := "originate/test-web-server:0.0.1"
			Expect(dockerCompose.Services[serviceRole].Image).To(Equal(imageName))

			By("should have the volumes for the external dependency defined in application.yml")
			serviceRole = "mongo"
			Expect(dockerCompose.Services[serviceRole].Volumes).To(Equal([]string{"mongo__data_db:/data/db"}))
		})
	})

	var _ = Describe("compiles the docker compose project name properly", func() {
		expected := "spacetweet123"

		It("converts all characters to lowercase", func() {
			actual := composebuilder.GetDockerComposeProjectName("SpaceTweet123")
			Expect(actual).To(Equal(expected))
		})

		It("strips non-alphanumeric characters", func() {
			actual := composebuilder.GetDockerComposeProjectName("$Space-Tweet_123")
			Expect(actual).To(Equal(expected))
		})

		It("strips whitespace characters", func() {
			actual := composebuilder.GetDockerComposeProjectName("Space   Tweet  123")
			Expect(actual).To(Equal(expected))
		})
	})
})
