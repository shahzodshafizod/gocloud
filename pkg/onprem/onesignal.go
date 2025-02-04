package onprem

import (
	"context"
	"net/http"
	"os"

	"github.com/OneSignal/onesignal-go-api"
	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type notification struct {
	appID      string
	restApiKey string
	client     *onesignal.APIClient
}

func NewNotification() pkg.Notification {
	return &notification{
		appID:      os.Getenv("ONESIGNAL_APP_ID"),
		restApiKey: os.Getenv("ONESIGNAL_REST_API_KEY"),
		client:     onesignal.NewAPIClient(onesignal.NewConfiguration()),
	}
}

func (n *notification) SendPush(ctx context.Context, to string, body string, priority int) (string, error) {
	ctx = context.WithValue(ctx, onesignal.AppAuth, n.restApiKey)

	contents := onesignal.NewStringMap()
	contents.SetEn(body)

	headings := onesignal.NewStringMap()
	headings.SetEn("Notification Header")

	var name = "Notification Name"

	var prty int32 = int32(priority)
	notification := onesignal.Notification{
		IncludePlayerIds: []string{to}, // equals to "include_subscription_ids"
		AppId:            n.appID,
		Contents:         *onesignal.NewNullableStringMap(contents),
		Headings:         *onesignal.NewNullableStringMap(headings),
		Data:             map[string]any{},
		Priority:         *onesignal.NewNullableInt32(&prty),
		Name:             &name,
	}

	successResp, resp, err := n.client.DefaultApi.
		CreateNotification(ctx).
		Notification(notification).
		Execute()
	if err != nil {
		return "", errors.Wrap(err, "CreateNotification.Execute")
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	return successResp.Id, nil
}

func (n *notification) SendToMany(ctx context.Context, tos []string, body string, priority int) (string, error) {
	contents := onesignal.NewStringMap()
	contents.SetEn(body)

	headings := onesignal.NewStringMap()
	headings.SetEn("Notification Header")

	var prty int32 = int32(priority)
	var name = "Notification Name"

	notification := onesignal.Notification{
		IncludePlayerIds: tos, // equals to "include_subscription_ids"
		AppId:            n.appID,
		Contents:         *onesignal.NewNullableStringMap(contents),
		Headings:         *onesignal.NewNullableStringMap(headings),
		Data:             map[string]any{},
		Priority:         *onesignal.NewNullableInt32(&prty),
		Name:             &name,
	}

	successResp, resp, err := n.client.DefaultApi.
		CreateNotification(ctx).
		Notification(notification).
		Execute()
	if err != nil {
		return "", errors.Wrap(err, "CreateNotification.Execute")
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	return successResp.Id, nil
}

func (n *notification) SendToTopic(ctx context.Context, topic string, data map[string]string) (string, error) {

	headings := onesignal.NewStringMap()
	headings.SetEn("Notification Header")

	var priority int32 = 10
	var name = "Notification Name"

	notification := onesignal.Notification{
		WebPushTopic: &topic,
		AppId:        n.appID,
		Headings:     *onesignal.NewNullableStringMap(headings),
		Priority:     *onesignal.NewNullableInt32(&priority),
		Name:         &name,
	}
	for key, value := range data {
		notification.Data[key] = value
	}

	successResp, resp, err := n.client.DefaultApi.
		CreateNotification(ctx).
		Notification(notification).
		Execute()
	if err != nil {
		return "", errors.Wrap(err, "CreateNotification.Execute")
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	return successResp.Id, nil
}
