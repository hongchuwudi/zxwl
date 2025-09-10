package main

import (
	"log"
	chat "mymod/ChatMgr"
	"mymod/MysqlMgr"
	"mymod/new_zxwl/api"
	"net/http"
)

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

func main() {
	MysqlMgr.InitDB()
	defer MysqlMgr.Db.Close()
	chat.InitDB()
	defer chat.Dbd.Close()
	mux := http.NewServeMux()
	if err := chat.CreateTable(MysqlMgr.Db); err != nil {
		log.Fatal("创建聊天表失败:", err)
	}
	go chat.HandleMessages()

	// 用户接口
	mux.HandleFunc("/auth", MysqlMgr.AuthHandler)                    // 登录
	mux.HandleFunc("/get_varifycode", MysqlMgr.GetVerifyCodeHandler) // 获取验证码
	mux.HandleFunc("/reset_pwd", MysqlMgr.ChangePasswordHandler)     // 修改密码
	mux.HandleFunc("/user_register", MysqlMgr.UserRegisterHandler)   // 注册
	mux.HandleFunc("/profile", MysqlMgr.ProfileHandler)              // 获取用户信息
	mux.HandleFunc("/profile/update", MysqlMgr.UpdateProfileHandler) // 修改用户信息
	mux.HandleFunc("/profile/delete", MysqlMgr.DeleteProfileHandler) // 删除用户
	mux.HandleFunc("/profile/list", MysqlMgr.ProfileListHandler)     // 管理员查看所有用户列表

	// 家庭接口
	mux.HandleFunc("/family/remove", MysqlMgr.FamilyRemoveHandler)
	mux.HandleFunc("/family/add", MysqlMgr.FamilyAddHandler)
	mux.HandleFunc("/family/find", MysqlMgr.FamilyFindHandler)

	// 聊天接口
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("新的WebSocket连接: %s", r.RemoteAddr)
		chat.ChatHandler(w, r)
		log.Printf("WebSocket连接关闭: %s", r.RemoteAddr)
	})
	mux.HandleFunc("/chat/unread", MysqlMgr.UnreadHandler)
	mux.HandleFunc("/chat/batch-unread", MysqlMgr.BatchUnreadHandler)
	mux.HandleFunc("/chat/mark-read", MysqlMgr.MarkReadHandler)
	mux.HandleFunc("/chat/send", MysqlMgr.SendHandler)

	// 政策接口
	mux.HandleFunc("/policy", MysqlMgr.PolicyHandler)
	mux.HandleFunc("/policy/instert", MysqlMgr.PolicyCreateHandler)
	mux.HandleFunc("/policy/search", MysqlMgr.PolicySearchHandler)

	// 志愿接口
	mux.HandleFunc("/volunteer/fetch", MysqlMgr.FetchVolunteerHandler)

	// 日志接口
	mux.HandleFunc("/log", MysqlMgr.LogInsertHandler)
	mux.HandleFunc("/logs", MysqlMgr.LogRetrieveHandler)

	// 旧接口/未使用的接口(下面已经重新实现)
	mux.HandleFunc("/professional/upsert", MysqlMgr.ProfessionalUpsertHandler)
	mux.HandleFunc("/professional/insert", MysqlMgr.ProfessionalInsertHandler)
	mux.HandleFunc("/professional/query", MysqlMgr.ProfessionalQueryHandler)
	mux.HandleFunc("/schools/profile", MysqlMgr.SchoolProfileHandler)

	// 新加的 学校/专业/推荐
	mux.HandleFunc("/recommends", api.RecommendationHandlers)               // 推荐服务
	mux.HandleFunc("/schools/profiles/", api.GetSchoolProfileByIDHandler)   // 学校查询
	mux.HandleFunc("/specials/profiles/", api.GetSpecialProfileByIDHandler) // 特殊查询
	mux.HandleFunc("/professional/querys", api.ProfessionalQueryHandlers)   // 专业查询
	mux.HandleFunc("/gAllSchools/", api.CollegesHandlers)                   // 获取所有学校
	mux.HandleFunc("/searchSchAndSpe/", api.SearchSchoolAndSpecial)         // 搜索学校和专业
	mux.HandleFunc("/special/name", api.GetSpecialIDByNameHandler)          // 根据专业名称获取专业ID
	mux.HandleFunc("/school/name", api.GetSchoolIDByNameHandler)            // 根据学校名称获取学校ID

	// 新加的 资讯论坛相关路由
	mux.HandleFunc("/news/list", api.NewsQueryHandlers)                  // 分页条件查询资讯
	mux.HandleFunc("/news/detail", api.GetNewsByIDHandlers)              // 按ID查询
	mux.HandleFunc("/news/insert", api.CreateNewsHandlers)               // 创建资讯
	mux.HandleFunc("/news/updateAll", api.UpdateNewsHandlers)            // 更新资讯
	mux.HandleFunc("/news/delete", api.DeleteNewsHandlers)               // 删除资讯
	mux.HandleFunc("/news/updateContent", api.UpdateNewsContentHandlers) // 更新内容
	mux.HandleFunc("/news/count/update", api.UpdateNewsCountHandler)
	mux.HandleFunc("/news/count/get", api.GetNewsCountHandler)

	// 新加的 添加资讯论坛评论相关路由
	mux.HandleFunc("/news/comments/list", api.GetCommentsByNewsIDHandler) // 根据文章ID获取评论
	mux.HandleFunc("/news/comments/add", api.AddCommentHandler)           // 添加评论
	mux.HandleFunc("/news/comments/like", api.ToggleCommentLikeHandler)   // 点赞/取消点赞

	// 新加的 给后端返回AiKey
	mux.HandleFunc("/getBaiDuKey", api.GetAPIKeyHandler)
	mux.HandleFunc("/aiChat", api.AIHandlers)
	handler := corsMiddleware(mux)

	log.Println("服务器启动，监听 :8792 端口...")
	log.Fatal(http.ListenAndServe(":8792", handler))
}
