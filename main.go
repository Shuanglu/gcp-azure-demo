package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azarm "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/sirupsen/logrus"
)

var (
	serviceAccount           = os.Getenv("GCP_SERVICE_ACCOUNT")
	endpoint                 = "http://metadata/computeMetadata/v1/instance/service-accounts/" + serviceAccount + "/identity?audience=api://AzureADTokenExchange"
	gcpMetadataRequestHeader = map[string]string{
		"Metadata-Flavor": "Google",
	}
	gcpOpenIDToken             = ""
	err                        error
	log                        = logrus.New()
	azTenantID                 = os.Getenv("AZURE_TENANT_ID")
	azClientID                 = os.Getenv("AZURE_CLIENT_ID")
	azSubID                    = os.Getenv("AZURE_SUBSCRIPTION_ID")
	errMsg_unauthorized_client = regexp.MustCompile("\"error\": \"unauthorized_client\\")
)

type User struct {
}

func main() {
	if serviceAccount == "" {
		log.Fatal("Please add environment variable \"GCP_SERVICE_ACCOUNT\"")
	}
	if azTenantID == "" {
		log.Fatal("Please add environment variable \"AZURE_TENANT_ID\"")
	}
	if azClientID == "" {
		log.Fatal("Please add environment variable \"AZURE_CLIENT_ID\"")
	}
	if azSubID == "" {
		log.Fatal("Please add environment variable \"AZURE_SUBSCRIPTION_ID\"")
	}
	user := User{}
	client := http.Client{}
	gcpOpenIDToken, err = user.GetOpenIdTokenGSA(context.TODO(), endpoint, client)
	if err != nil {
		log.Errorf("failed to get the OpenID token from GCP", err)
		errorHandling()
	}
	log.Infof("gcpToken: %v", gcpOpenIDToken)
	azClientAssertCredential, err := user.GetAzureCredential(context.TODO(), gcpOpenIDToken)
	if err != nil {
		log.Errorf("failed to get the credential from Azure", err)
		errorHandling()
	}
	azToken, err := azClientAssertCredential.GetToken(context.TODO(), policy.TokenRequestOptions{
		Scopes: []string{
			"https://vault.azure.net/.default",
		},
	})
	if err != nil {
		for {
			azToken, err = azClientAssertCredential.GetToken(context.TODO(), policy.TokenRequestOptions{
				Scopes: []string{
					"https://vault.azure.net/.default",
				},
			})
			if err == nil {
				break
			} else {
				if res := errMsg_unauthorized_client.FindString(err.Error()); res != "" {
					log.Errorf("Will retry in 10s because it failed to get the access token from Azure", err)
					time.Sleep(10 * time.Second)
				} else {
					log.Errorf("Won't retry due to error: %v", err)
					errorHandling()
				}
			}
		}
	}
	log.Infof("Azure Token: %v", azToken)
	azRGClient, err := azarm.NewResourceGroupsClient(azSubID, azClientAssertCredential, &arm.ClientOptions{})
	if err != nil {
		log.Errorf("failed to construct Azure ARM ResourceGroup client", err)
		errorHandling()
	}

	pager := azRGClient.NewListPager(nil)

	for pager.More() {
		nextResult, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Errorf("failed to read resource group page", err)
			errorHandling()
		}
		if nextResult.ResourceGroupListResult.Value != nil {
			for _, rg := range nextResult.ResourceGroupListResult.Value {
				log.Infof("Found Resource group: %v", *rg.Name)
			}
		}
	}
	for {
		select {}
	}
}

func errorHandling() {
	for {
		select {}
	}
}

func (u *User) GetOpenIdTokenGSA(ctx context.Context, endpoint string, client http.Client) (string, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Errorf("failed to construct request to retrive token from GCP endpoint", err)
		return "", err
	}
	for k, v := range gcpMetadataRequestHeader {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("failed to send request to retrieve token from GCP endpoint", err)
		return "", err
	}
	if resp.StatusCode != 200 {
		log.Errorf("response code of the request to retrieve token from GCP endpoint is non-200", resp.StatusCode)
		err = fmt.Errorf("response code of the request to retrieve token from GCP endpoint is non-200: %v", resp.StatusCode)
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to extract the response body of the request to retrieve token from GCP endpoint", err)
		return "", err
	}
	token := string(body)
	return token, nil
}

func (u *User) getAssertion(context.Context) (string, error) {
	if gcpOpenIDToken == "" {
		return "", fmt.Errorf("empty gcpOpenIDToken")
	}
	return gcpOpenIDToken, nil
}

func (u *User) GetAzureCredential(ctx context.Context, gcpOpenIDToken string) (*azidentity.ClientAssertionCredential, error) {
	clientAssertCredential, err := azidentity.NewClientAssertionCredential(azTenantID, azClientID, u.getAssertion, &azidentity.ClientAssertionCredentialOptions{})
	if err != nil {
		log.Errorf("failed to create clientAssertCredential", err)
		return nil, err
	}
	return clientAssertCredential, nil
	/*
		azAccessToken, err := clientAssertCredential.GetToken(ctx, policy.TokenRequestOptions{})
		if err != nil {
			log.Errorf("failed to get azure Access Token: %w", err)
			return "", err
		}
		return azAccessToken.Token, nil
	*/
}
