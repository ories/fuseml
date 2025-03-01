package acceptance_test

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MLflow Model", func() {
	var org = "mlflow-org"
	BeforeEach(func() {
		out, err := Fuseml("create-org "+org, "")
		Expect(err).ToNot(HaveOccurred(), out)
		out, err = Fuseml("target "+org, "")
		Expect(err).ToNot(HaveOccurred(), out)
	})
	Describe("push and delete", func() {
		var appName string
		BeforeEach(func() {
			appName = "mlflow-" + strconv.Itoa(int(time.Now().Nanosecond()))
		})

		It("pushes predicts and deletes an mlflow model", func() {
			By("pushing the mlflow model")
			currentDir, err := os.Getwd()
			time.Sleep(10 * time.Second)
			Expect(err).ToNot(HaveOccurred())
			appDir := path.Join(currentDir, "../examples/mlflow-model")

			pushCmd := fmt.Sprintf("push --serve %s ", serve)
			out, err := Fuseml(pushCmd+appName, appDir)
			Expect(err).ToNot(HaveOccurred(), out)
			out, err = Fuseml("apps", "")
			Expect(err).ToNot(HaveOccurred(), out)
			routeRegex := `.*\|.*1\/1.*\|.*`
			if serve == "knative" || serve == "kfserving" {
				routeRegex = `.*\|.*[0,1]\/[0,1].*\|.*`
			}
			Expect(out).To(MatchRegexp(appName + routeRegex))

			By("predict")
			re := regexp.MustCompile("http.*[a-z]")
			url := re.FindString(out)
			predict := predictMLFlowApp(url, serve)
			Expect(predict).To(MatchRegexp(`\[6.48634480\d+\]`))

			By("deleting the mlflow model")
			out, err = Fuseml("delete "+appName, "")
			Expect(err).ToNot(HaveOccurred(), out)
			// TODO: Fix `fuseml delete` from returning before the app is deleted #131
			Eventually(func() string {
				out, err := Fuseml("apps", "")
				Expect(err).ToNot(HaveOccurred(), out)
				return out
			}, "1m").ShouldNot(MatchRegexp(`.*%s.*`, appName))
		})
	})
})
