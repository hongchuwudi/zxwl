// file: main.go
package main

import (
	"log"
	"mymod/api/aiApi"
	"mymod/api/baseApi"
	"mymod/api/chatApi"
	"mymod/api/mixedApi"
	"mymod/api/newsApi"
	"mymod/api/oldChatApi"
	"mymod/api/oldManyMgrApi"
	"mymod/api/recomApi"
	"mymod/api/schoolApi"
	"mymod/api/scoreApi"
	"mymod/api/specialApi"
	"mymod/api/userApi"
	"mymod/api/userChooseApi"
	userFriendsApi "mymod/api/userFriendApi"
	"mymod/api/websocketApi"
	"mymod/config"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	uploadHandler      = baseApi.NewUploadHandler()
	userHandler        = userApi.NewUserHandler()
	userChooseHandler  = userChooseApi.NewUserChooseHandler()
	userFriendsHandler = userFriendsApi.NewFriendHandler()
	chatHandler        = chatApi.NewChatHandler()
	wsHandler          = websocketApi.NewWebSocketHandler()
	scoreHandler       = scoreApi.NewScoreHandler()
)

func main() {
	config.InitOSS()
	oldManyMgrApi.InitDB()
	defer oldManyMgrApi.Db.Close()
	oldChatApi.InitDB()
	defer oldChatApi.Dbd.Close()
	router := mux.NewRouter() // gorilla/mux 路由器

	if err := oldChatApi.CreateTable(oldManyMgrApi.Db); err != nil {
		log.Fatal("创建聊天表失败:", err)
	}
	go oldChatApi.HandleMessages()

	// WebSocket路由
	router.HandleFunc("/ws/user", wsHandler.HandleWebSocket).Methods("GET")        // 建立 WebSocket 连接
	router.HandleFunc("/ws/send/{userID}", wsHandler.SendMessage).Methods("POST")  // 发送消息（HTTP API）
	router.HandleFunc("/ws/online-users", wsHandler.GetOnlineUsers).Methods("GET") // 获取在线用户列表

	// 聊天接口
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("新的WebSocket连接: %s", r.RemoteAddr)
		oldChatApi.ChatHandler(w, r)
		log.Printf("WebSocket连接关闭: %s", r.RemoteAddr)
	})
	router.HandleFunc("/chat/unread", oldManyMgrApi.UnreadHandler).Methods("GET")
	router.HandleFunc("/chat/batch-unread", oldManyMgrApi.BatchUnreadHandler).Methods("GET")
	router.HandleFunc("/chat/mark-read", oldManyMgrApi.MarkReadHandler).Methods("POST")

	// 政策接口
	router.HandleFunc("/policy", oldManyMgrApi.PolicyHandler).Methods("GET")
	router.HandleFunc("/policy/instert", oldManyMgrApi.PolicyCreateHandler).Methods("POST") //创建政策
	router.HandleFunc("/policy/search", oldManyMgrApi.PolicySearchHandler).Methods("POST")
	router.HandleFunc("/policy/{id}", oldManyMgrApi.PolicyDeleteHandler).Methods("DELETE")

	// 日志接口
	router.HandleFunc("/log", oldManyMgrApi.LogInsertHandler).Methods("POST")
	router.HandleFunc("/logs", oldManyMgrApi.LogRetrieveHandler).Methods("GET")

	// 新加的 学校/专业/推荐
	router.HandleFunc("/recommends", recomApi.RecommendationHandlers).Methods("POST")                    // 推荐服务
	router.HandleFunc("/schools/profiles/{id}", schoolApi.GetSchoolProfileByIDHandler).Methods("GET")    // 学校查询
	router.HandleFunc("/specials/profiles/{id}", specialApi.GetSpecialProfileByIDHandler).Methods("GET") // 特殊查询
	router.HandleFunc("/professional/querys", specialApi.ProfessionalQueryHandlers).Methods("POST")      // 专业查询
	router.HandleFunc("/gAllSchools/", schoolApi.CollegesHandlers).Methods("GET")                        // 获取所有学校
	router.HandleFunc("/searchSchAndSpe/", mixedApi.SearchSchoolAndSpecial).Methods("GET")               // 搜索学校和专业
	router.HandleFunc("/special/name", specialApi.GetSpecialIDByNameHandler).Methods("GET")              // 根据专业名称获取专业ID
	router.HandleFunc("/school/name", schoolApi.GetSchoolIDByNameHandler).Methods("GET")                 // 根据学校名称获取学校ID

	// 新加的 一分一段
	router.HandleFunc("/scores/data", scoreHandler.GetScoreDataHandler).Methods("GET")
	router.HandleFunc("/scores/years", scoreHandler.GetAvailableYearsHandler).Methods("GET")

	// 新加的 资讯论坛相关路由
	router.HandleFunc("/news/list", newsApi.NewsQueryHandlers).Methods("POST")                 // 分页条件查询资讯
	router.HandleFunc("/news/detail", newsApi.GetNewsByIDHandlers).Methods("GET")              // 按ID查询
	router.HandleFunc("/news/insert", newsApi.CreateNewsHandlers).Methods("POST")              // 创建资讯
	router.HandleFunc("/news/updateAll", newsApi.UpdateNewsHandlers).Methods("PUT")            // 更新资讯
	router.HandleFunc("/news/delete", newsApi.DeleteNewsHandlers).Methods("DELETE")            // 删除资讯
	router.HandleFunc("/news/updateContent", newsApi.UpdateNewsContentHandlers).Methods("PUT") // 更新内容
	router.HandleFunc("/news/count/update", newsApi.UpdateNewsCountHandler).Methods("POST")
	router.HandleFunc("/news/count/get", newsApi.GetNewsCountHandler).Methods("POST") // 获取文章赞评数据

	// 新加的 添加资讯论坛评论相关路由
	router.HandleFunc("/news/comments/list", newsApi.GetCommentsByNewsIDHandler).Methods("POST") // 根据文章ID获取评论
	router.HandleFunc("/news/comments/add", newsApi.AddCommentHandler).Methods("POST")           // 添加评论
	router.HandleFunc("/news/comments/like", newsApi.ToggleCommentLikeHandler).Methods("POST")   // 点赞/取消点赞

	// 新加的 给后端返回百度地图ak
	router.HandleFunc("/getBaiDuKey", baseApi.GetAPIKeyHandler).Methods("GET")

	// 新加的 后端的AI服务
	router.HandleFunc("/aiChat", aiApi.AIHandlers).Methods("POST")

	// 新加的 用户相关路由
	router.HandleFunc("/get_varifycode", oldManyMgrApi.GetVerifyCodeHandler).Methods("POST") // 获取验证码
	router.HandleFunc("/user/register", userHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/user/login", userHandler.LoginHandler).Methods("POST")
	router.HandleFunc("/user/{id}", userHandler.GetUserHandler).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.DeleteUserHandler).Methods("DELETE")
	router.HandleFunc("/users", userHandler.GetUsersHandler).Methods("GET")                       // 分页查询用户
	router.HandleFunc("/user/update/{id}", userHandler.UpdateUserHandler).Methods("PUT")          // 更新用户信息
	router.HandleFunc("/user/change-password", userHandler.ChangePasswordHandler).Methods("POST") // 更改密码
	router.HandleFunc("/user/profile/{email}", userHandler.GetUserHandlerByEmail).Methods("GET")  // 根据邮箱获取用户信息

	// 添加文件上传路由
	router.HandleFunc("/upload/avatar", uploadHandler.UploadAvatarHandler).Methods("POST")
	router.HandleFunc("/upload/{folder}", uploadHandler.UploadFileHandler).Methods("POST")
	router.HandleFunc("/upload/delete", uploadHandler.DeleteFileHandler).Methods("DELETE")

	// 新加的 志愿服务
	router.HandleFunc("/user/{userID}/choices", userChooseHandler.GetUserChoices).Methods("GET")
	router.HandleFunc("/user/{userID}/choices", userChooseHandler.CreateUserChoice).Methods("POST")
	//router.HandleFunc("/user/{userID}/choices/{choiceID}", userChooseHandler.DeleteUserChoice).Methods("DELETE")
	router.HandleFunc("/user/{userID}/choices/all", userChooseHandler.DeleteAllUserChoices).Methods("DELETE")

	// 好友相关路由
	router.HandleFunc("/user/{userID}/friend/request/{friendID}", userFriendsHandler.SendFriendRequest).Methods("POST")   //发送好友请求
	router.HandleFunc("/user/{userID}/friend/accept/{requestID}", userFriendsHandler.AcceptFriendRequest).Methods("POST") //接受好友请求
	router.HandleFunc("/user/{userID}/friend/reject/{requestID}", userFriendsHandler.RejectFriendRequest).Methods("POST") // 拒绝好友请求
	router.HandleFunc("/user/{userID}/friends", userFriendsHandler.GetUserFriends).Methods("GET")                         //获取我的好友列表
	router.HandleFunc("/user/{userID}/friend-requests", userFriendsHandler.GetFriendRequestsToMe).Methods("GET")          //获取发给我的好友请求
	router.HandleFunc("/user/{userID}/friend/{friendID}/nickname", userFriendsHandler.SetFriendNickname).Methods("PUT")   //设置好友昵称
	router.HandleFunc("/user/{userID}/friend/{friendID}", userFriendsHandler.DeleteFriend).Methods("DELETE")              //删除好友

	// 聊天相关路由
	router.HandleFunc("/user/{userID}/chat/friend/{friendID}/history", chatHandler.GetFriendChatHistory).Methods("GET") // 获得好友历史消息
	router.HandleFunc("/user/{userID}/chat/mark-read/{friendID}", chatHandler.MarkMessagesAsRead).Methods("POST")
	router.HandleFunc("/user/{userID}/chat/friend/send", chatHandler.SendFriendMessage).Methods("POST")
	router.HandleFunc("/user/{userID}/chat/unread-count", chatHandler.GetUnreadCount).Methods("GET")
	router.HandleFunc("/user/{userID}/chat/unread-count/{friendID}", chatHandler.GetUnreadCountByFriend).Methods("GET")
	router.HandleFunc("/user/{userID}/chat/recent", chatHandler.GetRecentChats).Methods("GET")
	router.HandleFunc("/user/{userID}/chat/message/{messageID}", chatHandler.DeleteFriendMessage).Methods("DELETE")
	router.HandleFunc("/user/{userID}/chat/clear/{friendID}", chatHandler.ClearChatHistory).Methods("DELETE")
	// 群聊相关路由
	router.HandleFunc("/chat/room/{meetingID}/history", chatHandler.GetRoomChatHistory).Methods("GET")
	router.HandleFunc("/user/{userID}/chat/room/send", chatHandler.SendRoomMessage).Methods("POST")

	// 添加CORS中间件
	handler := corsMiddleware(router)
	log.Println("zxwl-服务器启动，监听 :8792 端口...")
	log.Fatal(http.ListenAndServe(":8792", handler))
}

// 跨域中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
