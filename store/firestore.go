package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/podozzoa/couponcrawler/model"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client
var m sync.Mutex
var LatestPost model.PostData

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

	err = GetLatestPost(ctx)
	if err != nil {
		fmt.Println("문서 데이터를 가져오는 데 실패했습니다.", err)
	}
}

func CloseFirestoreClient() {
	firestoreClient.Close()
}

func GetLatestPost(ctx context.Context) error {
	m.Lock()
	postQuery := firestoreClient.Collection("coupon_post").OrderBy("num", firestore.Desc).Limit(1)
	docs, err := postQuery.Documents(context.Background()).GetAll()
	if err != nil {
		fmt.Println("쿼리 실행에 실패했습니다.", err)
	}

	for _, doc := range docs {

		if err := doc.DataTo(&LatestPost); err != nil {

			return err
		}
		fmt.Printf("DB 내 가장 최근 포스트: %s (번호: %d)\n", LatestPost.Title, LatestPost.Num)
	}

	defer m.Unlock()
	return nil
}

func SavePosts(ctx context.Context, postList []model.PostData) {
	if len(postList) >= 1 {
		m.Lock()

		writer := firestoreClient.BulkWriter(ctx)
		collectionRef := firestoreClient.Collection("coupon_post")
		for _, post := range postList {
			if LatestPost.Num >= post.Num {
				break
			}
			postRef := collectionRef.Doc(createPostID(post.Num))
			writer.Set(postRef, post)
		}

		LatestPost = postList[0]

		writer.Flush()
		m.Unlock()
	}
}

//

func createPostID(postNum int) string {
	return fmt.Sprintf("Post%d", postNum)
}
