import { ref, computed } from 'vue';

const userName = ref(null);
const userEmail = ref(null);

// 从localStorage初始化数据
const initFromStorage = () => {
    userName.value = localStorage.getItem('userName');
    userEmail.value = localStorage.getItem('userEmail');
};

// 初始化
initFromStorage();

export function useUserStore() {
    // 保存用户信息
    const saveUser = (userData) => {
        // 先删除已有数据
        localStorage.removeItem('userName');
        localStorage.removeItem('userEmail');

        // 保存新数据
        userName.value = userData.name;
        userEmail.value = userData.email;

        localStorage.setItem('userName', userData.name);
        localStorage.setItem('userEmail', userData.email);

        return true;
    };

    // 读取用户信息
    const getUser = () => {
        // 确保数据最新
        initFromStorage();
        return {
            name: userName.value,
            email: userEmail.value
        };
    };

    // 删除用户信息
    const deleteUser = () => {
        userName.value = null;
        userEmail.value = null;

        localStorage.removeItem('userName');
        localStorage.removeItem('userEmail');

        return true;
    };

    // 判断是否登录
    const isLoggedIn = computed(() => {
        return !!userName.value && !!userEmail.value;
    });

    // 检查登录状态（同步方式）
    const checkLoginStatus = () => {
        initFromStorage(); // 确保数据最新
        return !!userName.value && !!userEmail.value;
    };

    return {
        userName,
        userEmail,
        saveUser,
        getUser,
        deleteUser,
        isLoggedIn, // 响应式计算属性
        checkLoginStatus // 同步检查函数
    };
}