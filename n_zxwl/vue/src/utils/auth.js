import { ref, computed } from 'vue';
import CryptoJS from 'crypto-js';

// 加密密钥，建议存储在环境变量中
const SECRET_KEY = 'your-secret-key-here';

// 加密函数
const encryptData = (data) => {
    return CryptoJS.AES.encrypt(JSON.stringify(data), SECRET_KEY).toString();
};

// 解密函数
const decryptData = (cipherText) => {
    if (!cipherText) return null;
    try {
        const bytes = CryptoJS.AES.decrypt(cipherText, SECRET_KEY);
        return JSON.parse(bytes.toString(CryptoJS.enc.Utf8));
    } catch (error) {
        console.error('Decryption error:', error);
        return null;
    }
};

const userName = ref(null);
const userEmail = ref(null);
const userId = ref(null);

// 从localStorage初始化数据
const initFromStorage = () => {
    const encryptedData = localStorage.getItem('userData');
    if (encryptedData) {
        const decryptedData = decryptData(encryptedData);
        if (decryptedData) {
            userName.value = decryptedData.name;
            userEmail.value = decryptedData.email;
            userId.value = decryptedData.id;
        }
    }
};

// 初始化
initFromStorage();

export function useUserStore() {
    // 保存用户信息
    const saveUser = (userData) => {
        // 加密数据
        const encryptedData = encryptData(userData);

        // 保存到localStorage
        localStorage.setItem('userData', encryptedData);

        // 更新响应式数据
        userName.value = userData.name;
        userEmail.value = userData.email;
        userId.value = userData.id;

        return true;
    };

    // 读取用户信息
    const getUser = () => {
        // 确保数据最新
        initFromStorage();
        return {
            id: userId.value,
            name: userName.value,
            email: userEmail.value
        };
    };

    // 删除用户信息
    const deleteUser = () => {
        userName.value = null;
        userEmail.value = null;
        userId.value = null;

        localStorage.removeItem('userData');

        return true;
    };

    // 判断是否登录
    const isLoggedIn = computed(() => {
        return !!userId.value && !!userName.value;
    });

    // 检查登录状态（同步方式）
    const checkLoginStatus = () => {
        initFromStorage(); // 确保数据最新
        return !!userId.value && !!userName.value
    }

    return {
        userId,
        userName,
        userEmail,
        saveUser,
        getUser,
        deleteUser,
        isLoggedIn, // 响应式计算属性
        checkLoginStatus // 同步检查函数
    }
}