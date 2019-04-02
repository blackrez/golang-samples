// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Contains examples calls to Cloud Security Center ListAssets API method.
package assets

// [START list_project_assets_at_time]
import (
	"context"
	"fmt"
	"io"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/iterator"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
)

// listAllProjectAssets lists all GCP Projects in orgID at asOf time and prints
// out results to w. listAllProjectAssets returns the number of project assets
// found.  orgID is the numeric organization ID of interest.
func listAllProjectAssetsAtTime(w io.Writer, orgID string, asOf time.Time) (int, error) {
	// orgID := "12321311"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error instantiating client %v\n", err)
		return -1, err
	}
	defer client.Close() // Closing the client safely cleans up background resources.
	onlyProjects := "security_center_properties.resource_type=" +
		"\"google.cloud.resourcemanager.Project\""

	// Convert the time to a Timestamp protobuf
	readTime, err := ptypes.TimestampProto(asOf)
	if err != nil {
		fmt.Printf("Error converting %v: %v", asOf, err)
		return 0, err
	}

	req := &securitycenterpb.ListAssetsRequest{
		Parent:   fmt.Sprintf("organizations/%s", orgID),
		Filter:   onlyProjects,
		ReadTime: readTime,
	}

	assetsFound := 0
	it := client.ListAssets(ctx, req)
	for {
		result, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, fmt.Errorf("Error listing assets: %v", err)
		}
		asset := &result.Asset
		properties := &(*asset).SecurityCenterProperties
		fmt.Fprintf(w, "Asset Name: %s, Resource Name %s, Resource Type %s\n",
			(*asset).Name,
			(*properties).ResourceName,
			(*properties).ResourceType)
		assetsFound++
	}
	return assetsFound, nil
}

// [END list_project_assets_at_time]
