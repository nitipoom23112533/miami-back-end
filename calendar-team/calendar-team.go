package calendarteam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CalendarService struct {
}

func NewCalendarService() *CalendarService {
	return &CalendarService{}
}


type CalendarInfo struct {
	ClientId          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	TtenantId          string `json:"tenantId"`
	UserId            string `json:"userId"`
	AccessToken       string `json:"accessToken"`
	OnlineMeetingUrl  string `json:"onlineMeetingUrl"`
}

// func (c *CalendarService)CreatCalendarAndTeam(ct *CalendarInfo) error {
func (c *CalendarService)CreatCalendarAndTeam() error {


	// newCalendar := CalendarInfo{
	// 	ClientId:          c.ClientId,
	// 	ClientSecret:      c.ClientSecret,
	// 	TtenantId:         c.TtenantId,
	// 	UserId:            c.UserId, // อีเมลที่มี calendar และ license Teams
	// 	AccessToken:       "",
	// 	OnlineMeetingUrl:  "",
	// }

	clientId := "YOUR_CLIENT_ID"
	clientSecret := "YOUR_CLIENT_SECRET"
	tenantId := "YOUR_TENANT_ID"
	userId := "user_email@domain.com" // อีเมลที่มี calendar และ license Teams

	// ขอ Access Token
	tokenUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantId)
	data := "client_id=" + clientId +
		"&scope=https%3A%2F%2Fgraph.microsoft.com%2F.default" +
		"&client_secret=" + clientSecret +
		"&grant_type=client_credentials"

	req, _ := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp map[string]interface{}
	json.Unmarshal(body, &tokenResp)

	accessToken := tokenResp["access_token"].(string)

	// เตรียมข้อมูล Event
	event := map[string]interface{}{
		"subject": "ประชุมทีม Go API",
		"body": map[string]string{
			"contentType": "HTML",
			"content":     "ประชุมเพื่อวางแผนโปรเจคผ่าน Teams",
		},
		"start": map[string]string{
			"dateTime": "2025-07-18T10:00:00",
			"timeZone": "Asia/Bangkok",
		},
		"end": map[string]string{
			"dateTime": "2025-07-18T11:00:00",
			"timeZone": "Asia/Bangkok",
		},
		"location": map[string]string{
			"displayName": "Online Meeting",
		},
		"isOnlineMeeting":      true,
		"onlineMeetingProvider": "teamsForBusiness",
	}

	eventBody, _ := json.Marshal(event)

	// เรียก Graph API
	eventUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/events", userId)
	eventReq, _ := http.NewRequest("POST", eventUrl, bytes.NewBuffer(eventBody))
	eventReq.Header.Set("Authorization", "Bearer "+accessToken)
	eventReq.Header.Set("Content-Type", "application/json")

	eventResp, err := http.DefaultClient.Do(eventReq)
	if err != nil {
		panic(err)
	}
	defer eventResp.Body.Close()

	eventRespBody, _ := io.ReadAll(eventResp.Body)

	// แสดงผล
	var eventRespMap map[string]interface{}
	json.Unmarshal(eventRespBody, &eventRespMap)

	if onlineMeeting, ok := eventRespMap["onlineMeeting"].(map[string]interface{}); ok {
		fmt.Println("Teams Meeting URL:", onlineMeeting["joinUrl"])
	} else {
		fmt.Println("ไม่พบลิงก์ online meeting")
	}

	// Optional: แสดง raw response
	// fmt.Println(string(eventRespBody))

	return nil

}