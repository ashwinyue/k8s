package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// LogEntry 结构化日志条目
type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Service    string                 `json:"service"`
	TraceID    string                 `json:"trace_id,omitempty"`
	SpanID     string                 `json:"span_id,omitempty"`
	UserID     string                 `json:"user_id,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Duration   string                 `json:"duration,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// Logger 自定义日志器
type Logger struct {
	logger  *logrus.Logger
	service string
	logFile *os.File
}

// NewLogger 创建新的日志器
func NewLogger(service string) (*Logger, error) {
	// 创建日志目录
	logDir := "/app/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建日志文件
	logFile, err := os.OpenFile(
		filepath.Join(logDir, "app.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return nil, fmt.Errorf("创建日志文件失败: %v", err)
	}

	// 配置 logrus
	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	// 设置日志级别
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return &Logger{
		logger:  logger,
		service: service,
		logFile: logFile,
	}, nil
}

// Close 关闭日志器
func (l *Logger) Close() error {
	return l.logFile.Close()
}

// LogWithContext 记录带上下文的日志
func (l *Logger) LogWithContext(level string, message string, ctx map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Service:   l.service,
		Extra:     ctx,
	}

	// 从上下文中提取常用字段
	if traceID, ok := ctx["trace_id"]; ok {
		entry.TraceID = fmt.Sprintf("%v", traceID)
	}
	if spanID, ok := ctx["span_id"]; ok {
		entry.SpanID = fmt.Sprintf("%v", spanID)
	}
	if userID, ok := ctx["user_id"]; ok {
		entry.UserID = fmt.Sprintf("%v", userID)
	}
	if requestID, ok := ctx["request_id"]; ok {
		entry.RequestID = fmt.Sprintf("%v", requestID)
	}
	if method, ok := ctx["method"]; ok {
		entry.Method = fmt.Sprintf("%v", method)
	}
	if path, ok := ctx["path"]; ok {
		entry.Path = fmt.Sprintf("%v", path)
	}
	if statusCode, ok := ctx["status_code"]; ok {
		if code, err := strconv.Atoi(fmt.Sprintf("%v", statusCode)); err == nil {
			entry.StatusCode = code
		}
	}
	if duration, ok := ctx["duration"]; ok {
		entry.Duration = fmt.Sprintf("%v", duration)
	}
	if errorMsg, ok := ctx["error"]; ok {
		entry.Error = fmt.Sprintf("%v", errorMsg)
	}

	// 序列化并写入日志
	logData, _ := json.Marshal(entry)
	fmt.Fprintln(l.logFile, string(logData))
	l.logFile.Sync()
}

// Info 记录信息日志
func (l *Logger) Info(message string, ctx map[string]interface{}) {
	l.LogWithContext("info", message, ctx)
	logEntriesTotal.WithLabelValues("info").Inc()
}

// Warn 记录警告日志
func (l *Logger) Warn(message string, ctx map[string]interface{}) {
	l.LogWithContext("warn", message, ctx)
	logEntriesTotal.WithLabelValues("warn").Inc()
}

// Error 记录错误日志
func (l *Logger) Error(message string, ctx map[string]interface{}) {
	l.LogWithContext("error", message, ctx)
	logEntriesTotal.WithLabelValues("error").Inc()
}

// Debug 记录调试日志
func (l *Logger) Debug(message string, ctx map[string]interface{}) {
	l.LogWithContext("debug", message, ctx)
	logEntriesTotal.WithLabelValues("debug").Inc()
}

// 全局变量
var (
	appLogger *Logger
	userIDs   = []string{"user1", "user2", "user3", "user4", "user5"}
	paths     = []string{"/api/users", "/api/orders", "/api/products", "/api/payments", "/health"}

	// Prometheus 指标
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	logEntriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_entries_total",
			Help: "Total number of log entries",
		},
		[]string{"level"},
	)
)

// 中间件：请求日志记录
// initTracer 初始化 Jaeger 追踪
func initTracer() (*trace.TracerProvider, error) {
	// 创建 Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger:14268/api/traces")))
	if err != nil {
		return nil, err
	}

	// 创建资源
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-demo"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	// 创建 TracerProvider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	return tp, nil
}

// initMetrics 初始化 Prometheus 指标
func initMetrics() {
	// 注册指标
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(logEntriesTotal)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		traceID := uuid.New().String()
		spanID := uuid.New().String()[:8]

		// 设置上下文
		c.Set("request_id", requestID)
		c.Set("trace_id", traceID)
		c.Set("span_id", spanID)

		// 增加活跃连接数
		activeConnections.Inc()
		defer activeConnections.Dec()

		// 记录请求开始
		appLogger.Info("请求开始", map[string]interface{}{
			"request_id": requestID,
			"trace_id":   traceID,
			"span_id":    spanID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"user_agent": c.Request.UserAgent(),
			"remote_ip":  c.ClientIP(),
		})

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(start)
		status := c.Writer.Status()

		// 记录 Prometheus 指标
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, fmt.Sprintf("%d", status)).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(duration.Seconds())

		// 记录请求完成
		logLevel := "info"
		if status >= 400 {
			logLevel = "error"
		} else if status >= 300 {
			logLevel = "warn"
		}

		ctx := map[string]interface{}{
			"request_id":    requestID,
			"trace_id":      traceID,
			"span_id":       spanID,
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"status_code":   status,
			"duration":      duration.String(),
			"response_size": c.Writer.Size(),
		}

		if status >= 400 {
			ctx["error"] = fmt.Sprintf("HTTP %d", status)
		}

		appLogger.LogWithContext(logLevel, "请求完成", ctx)
	}
}

// API 处理函数
func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "go-demo",
		"version":   os.Getenv("APP_VERSION"),
	})
}

func usersHandler(c *gin.Context) {
	// 模拟随机延迟
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	// 模拟随机错误
	if rand.Float32() < 0.1 {
		appLogger.Error("数据库连接失败", map[string]interface{}{
			"request_id": c.GetString("request_id"),
			"trace_id":   c.GetString("trace_id"),
			"error":      "connection timeout",
			"database":   "users_db",
		})
		c.JSON(500, gin.H{"error": "内部服务器错误"})
		return
	}

	users := []map[string]interface{}{
		{"id": 1, "name": "张三", "email": "zhangsan@example.com"},
		{"id": 2, "name": "李四", "email": "lisi@example.com"},
		{"id": 3, "name": "王五", "email": "wangwu@example.com"},
	}

	appLogger.Info("获取用户列表成功", map[string]interface{}{
		"request_id": c.GetString("request_id"),
		"trace_id":   c.GetString("trace_id"),
		"user_count": len(users),
	})

	c.JSON(200, gin.H{"users": users})
}

func ordersHandler(c *gin.Context) {
	// 模拟业务逻辑
	userID := userIDs[rand.Intn(len(userIDs))]
	orderID := uuid.New().String()

	appLogger.Info("处理订单请求", map[string]interface{}{
		"request_id": c.GetString("request_id"),
		"trace_id":   c.GetString("trace_id"),
		"user_id":    userID,
		"order_id":   orderID,
		"action":     "create_order",
	})

	// 模拟支付处理
	if rand.Float32() < 0.05 {
		appLogger.Error("支付处理失败", map[string]interface{}{
			"request_id": c.GetString("request_id"),
			"trace_id":   c.GetString("trace_id"),
			"user_id":    userID,
			"order_id":   orderID,
			"error":      "payment gateway timeout",
			"gateway":    "stripe",
		})
		c.JSON(500, gin.H{"error": "支付处理失败"})
		return
	}

	appLogger.Info("订单创建成功", map[string]interface{}{
		"request_id": c.GetString("request_id"),
		"trace_id":   c.GetString("trace_id"),
		"user_id":    userID,
		"order_id":   orderID,
		"amount":     rand.Float64() * 1000,
		"currency":   "CNY",
	})

	c.JSON(200, gin.H{
		"order_id": orderID,
		"user_id":  userID,
		"status":   "created",
	})
}

// 后台任务：生成定期日志
func backgroundTasks(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 生成系统指标日志
			appLogger.Info("系统指标报告", map[string]interface{}{
				"metric_type":        "system",
				"cpu_usage":          rand.Float64() * 100,
				"memory_usage":       rand.Float64() * 100,
				"disk_usage":         rand.Float64() * 100,
				"active_connections": rand.Intn(1000),
				"goroutines":         rand.Intn(100),
			})

			// 随机生成一些业务事件
			if rand.Float32() < 0.3 {
				userID := userIDs[rand.Intn(len(userIDs))]
				appLogger.Info("用户活动", map[string]interface{}{
					"event_type": "user_login",
					"user_id":    userID,
					"ip_address": fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
					"user_agent": "Mozilla/5.0 (compatible)",
				})
			}

			// 随机生成一些警告
			if rand.Float32() < 0.1 {
				appLogger.Warn("性能警告", map[string]interface{}{
					"warning_type":  "high_response_time",
					"endpoint":      paths[rand.Intn(len(paths))],
					"response_time": fmt.Sprintf("%dms", 1000+rand.Intn(2000)),
					"threshold":     "1000ms",
				})
			}
		}
	}
}

func main() {
	// 初始化日志器
	var err error
	appLogger, err = NewLogger("go-demo")
	if err != nil {
		log.Fatalf("初始化日志器失败: %v", err)
	}
	defer appLogger.Close()

	// 记录应用启动
	appLogger.Info("应用启动", map[string]interface{}{
		"service": "go-demo",
		"version": os.Getenv("APP_VERSION"),
		"port":    "8080",
		"env":     os.Getenv("LOG_LEVEL"),
	})

	// 初始化 Prometheus 指标
	initMetrics()

	// 初始化 Jaeger 追踪
	tp, err := initTracer()
	if err != nil {
		appLogger.Error("Failed to initialize tracer", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				appLogger.Error("Error shutting down tracer provider", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
		appLogger.Info("Jaeger tracer initialized", nil)
	}

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由器
	r := gin.New()
	r.Use(otelgin.Middleware("go-demo"))
	r.Use(LoggingMiddleware())
	r.Use(gin.Recovery())

	// 注册路由
	r.GET("/health", healthHandler)
	r.GET("/api/users", usersHandler)
	r.POST("/api/orders", ordersHandler)
	r.GET("/api/products", func(c *gin.Context) {
		c.JSON(200, gin.H{"products": []string{"产品A", "产品B", "产品C"}})
	})
	r.POST("/api/payments", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "success", "transaction_id": uuid.New().String()})
	})

	// Prometheus 指标端点
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 创建 HTTP 服务器，明确监听所有接口
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	// 启动后台任务
	ctx, cancel := context.WithCancel(context.Background())
	go backgroundTasks(ctx)

	// 启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("服务器启动失败", map[string]interface{}{
				"error": err.Error(),
			})
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	appLogger.Info("服务器启动成功", map[string]interface{}{
		"port": "8080",
		"pid":  os.Getpid(),
	})

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("服务器关闭中...", map[string]interface{}{})

	// 优雅关闭
	cancel() // 停止后台任务
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("服务器强制关闭", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		appLogger.Info("服务器已优雅关闭", map[string]interface{}{})
	}
}
