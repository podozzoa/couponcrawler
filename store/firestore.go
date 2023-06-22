package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/podozzoa/couponcrawler/model"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var firestoreClient *firestore.Client
var m sync.Mutex
var latestPost model.PostData

func InitFirestoreClient(ctx context.Context) {
	// sa := option.WithCredentialsFile("firebase-adminsdk.json")
	firebaseKeyJSON := os.Getenv("FIREBASE_ADMINSDK_JSON")
	if len(firebaseKeyJSON) == 0 {
		log.Fatal("FIREBASE_ADMINSDK_JSON 환경 변수가 설정되지 않았습니다.")
	}
	credsJSON := []byte(firebaseKeyJSON)
	sa := option.WithCredentialsJSON(credsJSON)

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func CloseFirestoreClient() {
	firestoreClient.Close()
}

func SavePosts(ctx context.Context, postList []model.PostData) {
	if len(postList) >= 1 {
		m.Lock()

		newPost := postList[0]

		if latestPost.Num != newPost.Num {
			latestPost = newPost
			numDiff := newPost.Num - latestPost.Num

			for i := 0; i < numDiff; i++ {
				post := postList[i]
				postMap, err := postDataToMap([]model.PostData{post})
				if err != nil {
					log.Printf("Failed to convert PostData to map: %v", err)
					m.Unlock()
					return
				}
				postId := createPostID(post.Num)
				postRef := firestoreClient.Collection("coupon_post").Doc(postId)
				// 포스트가 이미 존재하는지 확인합니다.
				_, err = postRef.Get(context.Background())
				if status.Code(err) == codes.NotFound {
					_, err = postRef.Set(context.Background(), postMap[0])
					if err != nil {
						log.Fatalf("Failed to save data for post %d: %v", post.Num, err)
					}
				} else if err != nil {
					log.Printf("Failed to get post %d: %v", post.Num, err)
				} // 존재하면 이미 저장되어 있으므로 아무 것도 하지 않습니다.
			}
		}

		m.Unlock()
	}
}

func postDataToMap(postList []model.PostData) ([]map[string]interface{}, error) {
	postListMap := make([]map[string]interface{}, len(postList))
	for i, post := range postList {
		jsonPost, err := json.Marshal(post)
		if err != nil {
			log.Printf("Failed to marshal post data: %v", err)
			return nil, err
		}

		var postMap map[string]interface{}
		err = json.Unmarshal(jsonPost, &postMap)
		if err != nil {
			log.Printf("Failed to unmarshal post data: %v", err)
			return nil, err
		}

		postListMap[i] = postMap
	}

	return postListMap, nil
}

func createPostID(postNum int) string {
	return fmt.Sprintf("Post%d", postNum)
}
