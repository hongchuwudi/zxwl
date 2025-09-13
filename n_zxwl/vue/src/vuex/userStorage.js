// userStorage.js - Vue 3用户信息存储工具
import { createStore } from 'vuex';
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

// 创建Vuex store
const store = createStore({
    state() {
        return {
            userName: null,
            userEmail: null,
            userId: null
        };
    },
    mutations: {
        // 保存用户信息到state和localStorage
        SAVE_USER(state, userData) {
            state.userName = userData.name;
            state.userEmail = userData.email;
            state.userId = userData.id;

            // 加密并保存到localStorage
            const encryptedData = encryptData(userData);
            localStorage.setItem('userData', encryptedData);
        },
        // 从localStorage加载用户信息到state
        LOAD_USER(state) {
            const encryptedData = localStorage.getItem('userData');
            if (encryptedData) {
                const decryptedData = decryptData(encryptedData);
                if (decryptedData) {
                    state.userName = decryptedData.name;
                    state.userEmail = decryptedData.email;
                    state.userId = decryptedData.id;
                }
            }
        },
        // 删除用户信息
        DELETE_USER(state) {
            state.userName = null;
            state.userEmail = null;
            state.userId = null;

            localStorage.removeItem('userData');
        }
    },
    actions: {
        // 保存用户信息
        saveUser({ commit }, userData) {
            commit('SAVE_USER', userData);
            return true;
        },
        // 读取用户信息
        getUser({ commit }) {
            commit('LOAD_USER');
            return {
                id: this.state.userId,
                name: this.state.userName,
                email: this.state.userEmail
            };
        },
        // 删除用户信息
        deleteUser({ commit }) {
            commit('DELETE_USER');
            return true;
        }
    }
});

// 导出store实例
export default store;

// 导出单独的方法，方便直接使用
export const saveUser = (userData) => store.dispatch('saveUser', userData)

export const getUser = () => store.dispatch('getUser')

export const deleteUser = () => store.dispatch('deleteUser')