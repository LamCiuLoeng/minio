// +build ignore

package main

/*
 * Minio Cloud Storage, (C) 2017 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"fmt"
	"log"

	"github.com/minio/minio/pkg/madmin"
)

func main() {

	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY are
	// dummy values, please replace them with original values.

	// API requests are secure (HTTPS) if secure=true and insecure (HTTPS) otherwise.
	// New returns an Minio Admin client object.
	madmClnt, err := madmin.New("your-minio.example.com:9000", "YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", true)
	if err != nil {
		log.Fatalln(err)
	}

	bucket := "mybucket"
	prefix := "myprefix"

	// Create a done channel to control 'ListUploadsHeal' go routine.
	doneCh := make(chan struct{})
	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// Set true if recursive listing is needed.
	isRecursive := true
	// List objects that need healing for a given bucket and
	// prefix.
	healUploadsCh, err := madmClnt.ListUploadsHeal(bucket, prefix, isRecursive, doneCh)
	if err != nil {
		log.Fatalln("Failed to get list of uploads to be healed: ", err)
	}

	for upload := range healUploadsCh {
		if upload.Err != nil {
			log.Println("upload listing error: ", upload.Err)
		}

		if upload.HealUploadInfo != nil {
			switch healInfo := *upload.HealUploadInfo; healInfo.Status {
			case madmin.CanHeal:
				fmt.Println(upload.Key, " can be healed.")
			case madmin.QuorumUnavailable:
				fmt.Println(upload.Key, " can't be healed until quorum is available.")
			case madmin.Corrupted:
				fmt.Println(upload.Key, " can't be healed, not enough information.")
			}
		}
	}
}
