package main

import (
	"context"
	"fmt"
	"io"
	"cloud.google.com/go/storage"
)

const (
	BUCKET_NAME = "my-around-bucket"
)

func saveToGCS(r io.Reader, objectName string) (string ,error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	object := client.Bucket(BUCKET_NAME).Object(objectName)//objectName是上传之后在云端的名字
	wc := object.NewWriter(ctx)

	if _, err := io.Copy(wc, r); err != nil {//r是local file，copy到云端的wc
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {//access control
		//在object level 设置 public read
		//RoleReader,读权限
		//RoleWriter，写权限
		//RoleOwner，读写权限
		return "", err
	}

	attrs, err := object.Attrs(ctx)//MediaLink is the url
	if err != nil {
		return "", err
	}
	fmt.Printf("Image is saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}

