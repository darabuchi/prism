package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/prism/core/internal/service"
	"github.com/prism/core/internal/storage"
	"github.com/prism/core/internal/testutil"
)

func setupSubscriptionHandler(t *testing.T) (*SubscriptionHandler, *storage.Database) {
	db := testutil.SetupTestDB(t)
	subscriptionSvc := service.NewSubscriptionService(db)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // 减少测试日志输出

	handler := NewSubscriptionHandler(subscriptionSvc, logger)
	return handler, db
}

func TestSubscriptionHandler_Create(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)

	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	req := service.CreateSubscriptionRequest{
		Name:           "Test Subscription",
		URL:            "https://example.com/subscription",
		UserAgent:      "Test-Agent/1.0",
		AutoUpdate:     true,
		UpdateInterval: 3600,
		NodePoolIDs:    []uint{},
	}

	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/subscriptions", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// 验证响应
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	data := response["data"].(map[string]interface{})
	if data["name"] != req.Name {
		t.Errorf("Expected name %q, got %v", req.Name, data["name"])
	}
}

func TestSubscriptionHandler_Create_InvalidRequest(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)

	// 发送无效请求（缺少必需字段）
	req := map[string]interface{}{
		"name": "", // 空名称应该失败
		"url":  "invalid-url",
	}

	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/subscriptions", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSubscriptionHandler_List(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/subscriptions?page=1&size=10", nil)

	handler.List(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	data := response["data"].(map[string]interface{})
	if data["total"].(float64) == 0 {
		t.Error("Expected at least one subscription")
	}
}

func TestSubscriptionHandler_GetByID(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	gin.SetMode(gin.TestMode)

	// 获取测试订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/subscriptions/"+strconv.Itoa(int(subscription.ID)), nil)
	c.Params = []gin.Param{
		{Key: "id", Value: strconv.Itoa(int(subscription.ID))},
	}

	handler.GetByID(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	data := response["data"].(map[string]interface{})
	if data["id"].(float64) != float64(subscription.ID) {
		t.Errorf("Expected ID %d, got %v", subscription.ID, data["id"])
	}
}

func TestSubscriptionHandler_GetByID_NotFound(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/subscriptions/9999", nil)
	c.Params = []gin.Param{
		{Key: "id", Value: "9999"},
	}

	handler.GetByID(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestSubscriptionHandler_Update(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	gin.SetMode(gin.TestMode)

	// 获取测试订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	newName := "Updated Subscription"
	req := service.UpdateSubscriptionRequest{
		Name: &newName,
	}

	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/subscriptions/"+strconv.Itoa(int(subscription.ID)), bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{
		{Key: "id", Value: strconv.Itoa(int(subscription.ID))},
	}

	handler.Update(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	data := response["data"].(map[string]interface{})
	if data["name"] != newName {
		t.Errorf("Expected name %q, got %v", newName, data["name"])
	}
}

func TestSubscriptionHandler_Delete(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	gin.SetMode(gin.TestMode)

	// 获取测试订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/subscriptions/"+strconv.Itoa(int(subscription.ID)), nil)
	c.Params = []gin.Param{
		{Key: "id", Value: strconv.Itoa(int(subscription.ID))},
	}

	handler.Delete(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// 验证订阅已被删除
	var count int64
	db.DB.Model(&storage.Subscription{}).Where("id = ?", subscription.ID).Count(&count)
	if count != 0 {
		t.Error("Expected subscription to be deleted")
	}
}

func TestSubscriptionHandler_UpdateSubscription(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)

	gin.SetMode(gin.TestMode)

	// 创建测试服务器
	mockServer := testutil.MockHTTPServer(map[string]string{
		"/subscription": testutil.MockSubscriptionData,
	})
	defer mockServer.Close()

	// 创建订阅
	subscription := &storage.Subscription{
		Name:   "Test Subscription",
		URL:    mockServer.URL + "/subscription",
		Status: "active",
	}
	db.DB.Create(subscription)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/subscriptions/"+strconv.Itoa(int(subscription.ID))+"/update", nil)
	c.Params = []gin.Param{
		{Key: "id", Value: strconv.Itoa(int(subscription.ID))},
	}

	handler.UpdateSubscription(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}
}

func TestSubscriptionHandler_GetStats(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	gin.SetMode(gin.TestMode)

	// 获取测试订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/subscriptions/"+strconv.Itoa(int(subscription.ID))+"/stats", nil)
	c.Params = []gin.Param{
		{Key: "id", Value: strconv.Itoa(int(subscription.ID))},
	}

	handler.GetStats(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	data := response["data"].(map[string]interface{})
	if data["subscription_id"].(float64) != float64(subscription.ID) {
		t.Errorf("Expected subscription ID %d, got %v", subscription.ID, data["subscription_id"])
	}
}

func TestSubscriptionHandler_GetLogs(t *testing.T) {
	handler, db := setupSubscriptionHandler(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	gin.SetMode(gin.TestMode)

	// 获取测试订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/subscriptions/"+strconv.Itoa(int(subscription.ID))+"/logs?page=1&size=10", nil)
	c.Params = []gin.Param{
		{Key: "id", Value: strconv.Itoa(int(subscription.ID))},
	}

	handler.GetLogs(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}
}
