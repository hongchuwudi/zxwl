package websocketApi

import (
	"log"
	"mymod/config"
	"mymod/repositories/userFriendsRepo"
	"time"
)

/**
 * 通知好友在线状态变化
 */
// notifyFriendsOnlineStatus 通知好友在线状态变化
func (h *WebSocketHandler) notifyFriendsOnlineStatus(userID int64, isOnline bool) {
	// 获取用户的好友列表
	friendIDs, err := h.getUserFriends(userID)
	if err != nil {
		log.Printf("获取用户 %s 的好友列表失败: %v", userID, err)
		return
	}

	// 创建状态通知消息
	statusMsg := map[string]interface{}{
		"type": "friend_status",
		"data": map[string]interface{}{
			"user_id":   userID,
			"is_online": isOnline,
			"timestamp": time.Now().Unix(),
		},
	}

	// 向每个在线好友发送状态通知
	for _, friendID := range friendIDs {
		if h.isUserOnline(friendID) {
			//log.Printf("向好友 %s 发送状态通知", friendID)
			if err := h.sendToUser(friendID, statusMsg); err != nil {
				log.Printf("向好友 %s 发送状态通知失败: %v", friendID, err)
			}
		}
	}

	//log.Printf("用户 %s 的好友在线状态已更新为: %t", userID, isOnline)
}

// getUserFriends 获取用户的好友列表
func (h *WebSocketHandler) getUserFriends(userID int64) ([]int64, error) {
	// 根据用户id获取好友id列表
	friendRepo := userFriendsRepo.NewFriendRepository(config.GetDB())

	friends, err := friendRepo.GetUserFriends(userID)
	if err != nil {
		return nil, err
	}

	friendIDs := []int64{}
	for _, friend := range friends {
		//log.Println("获取用户好友列表:", friend.UserAID)
		friendIDs = append(friendIDs, friend.UserBID)
	}
	return friendIDs, nil
}
