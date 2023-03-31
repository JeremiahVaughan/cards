package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"strings"
)

func main() {
	//assignment_two.AssignmentTwo(os.Args[1])
	//channels.Channels()
	//channel_experiments.ChannelExperiments()

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		log.Fatalf("error, when creating AWS session: %v", err)
	}

	// Get ECR authentication token
	ecrSvc := ecr.New(sess)
	authInput := &ecr.GetAuthorizationTokenInput{}
	authOutput, err := ecrSvc.GetAuthorizationToken(authInput)

	if err != nil {
		log.Fatalf("error, when getting ECR authorization token: %v", err)

	}

	// Extract username, password, and registry URL from the ECR token
	authData := authOutput.AuthorizationData[0]
	authToken, err := base64.StdEncoding.DecodeString(aws.StringValue(authData.AuthorizationToken))

	if err != nil {
		log.Fatalf("error, when decoding ECR authorization token: %v", err)
	}

	tokenParts := strings.SplitN(string(authToken), ":", 2)
	username := tokenParts[0]
	password := tokenParts[1]
	registryURL := aws.StringValue(authData.ProxyEndpoint)

	// Login to Docker registry using the ECR token
	ctx := context.Background()
	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("error, when creating Docker client: %v", err)
	}

	authConfig := dockerTypes.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: registryURL,
	}
	encodedAuth, err := dockerCli.RegistryLogin(ctx, authConfig)

	if err != nil {
		log.Fatalf("error, when logging in to Docker registry: %v", err)
	}

	ecrImageURI := fmt.Sprintf("%s/%s", strings.TrimPrefix(registryURL, "https://"), "lambda-house")
	ecrImageName := fmt.Sprintf("%s:%s", ecrImageURI, "busybox")

	if err = dockerCli.ImageTag(ctx, "busybox", ecrImageName); err != nil {
		fmt.Println("Error tagging Docker image:", err)
		return
	}

	authConfigBytes, err := json.Marshal(authConfig)
	if err != nil {
		fmt.Println("Error marshaling auth config:", err)
		return
	}

	// Push the Docker image to ECR
	pushOptions := dockerTypes.ImagePushOptions{
		RegistryAuth: base64.URLEncoding.EncodeToString(authConfigBytes),
	}
	pushResponse, err := dockerCli.ImagePush(ctx, ecrImageName, pushOptions)

	if err != nil {
		fmt.Println("Error pushing Docker image to ECR:", err)
		return
	}
	defer pushResponse.Close()

	fmt.Println("Login to ECR successful. Status:", encodedAuth.Status)

}
