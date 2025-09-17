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

// 创建响应式变量来存储完整的用户信息
const userData = ref(null);

// 从localStorage初始化数据
const initFromStorage = () => {
    const encryptedData = localStorage.getItem('userData');
    if (encryptedData) {
        const decryptedData = decryptData(encryptedData);
        if (decryptedData) {
            userData.value = decryptedData;
        }
    }
};

// 初始化
initFromStorage();

export function useUserStore() {
    // 保存用户信息
    const saveUser = (userInfo) => {
        // 加密数据
        const encryptedData = encryptData(userInfo);

        // 保存到localStorage
        localStorage.setItem('userData', encryptedData);

        // 更新响应式数据
        userData.value = userInfo;

        return true;
    };

    // 读取用户信息
    const getUser = () => {
        // 确保数据最新
        initFromStorage();
        return userData.value;
    };

    // 获取特定用户字段
    const getUserField = (field) => {
        initFromStorage();
        return userData.value ? userData.value[field] : null;
    };

    // 删除用户信息
    const deleteUser = () => {
        userData.value = null;
        localStorage.removeItem('userData');
        return true;
    };

    // 判断是否登录
    const isLoggedIn = computed(() => {
        return !!userData.value && !!userData.value.id;
    });

    // 检查登录状态（同步方式）
    const checkLoginStatus = () => {
        initFromStorage(); // 确保数据最新
        return !!userData.value && !!userData.value.id;
    };

    return {
        userData,
        saveUser,
        getUser,
        getUserField,
        deleteUser,
        isLoggedIn, // 响应式计算属性
        checkLoginStatus // 同步检查函数
    };
}