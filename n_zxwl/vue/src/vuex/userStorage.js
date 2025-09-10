// userStorage.js - Vue 3用户信息存储工具

import { createStore } from 'vuex';

// 创建Vuex store
const store = createStore({
    state() {
        return {
            userName: null,
            userEmail: null
        };
    },
    mutations: {
        // 保存用户信息到state和localStorage
        SAVE_USER(state, { name, email }) {
            state.userName = name;
            state.userEmail = email;

            // 同时保存到localStorage
            localStorage.setItem('userName', name);
            localStorage.setItem('userEmail', email);
        },
        // 从localStorage加载用户信息到state
        LOAD_USER(state) {
            state.userName = localStorage.getItem('userName');
            state.userEmail = localStorage.getItem('userEmail');
        },
        // 删除用户信息
        DELETE_USER(state) {
            state.userName = null;
            state.userEmail = null;

            // 同时从localStorage删除
            localStorage.removeItem('userName');
            localStorage.removeItem('userEmail');
        }
    },
    actions: {
        // 保存用户信息
        saveUser({ commit }, userData) {
            // 如果有就删掉再保存
            if (localStorage.getItem('userName') || localStorage.getItem('userEmail')) {
                commit('DELETE_USER');
            }

            commit('SAVE_USER', userData);
            return true;
        },
        // 读取用户信息
        getUser({ commit }) {
            commit('LOAD_USER');
            return {
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
export const saveUser = (userData) => {
    return store.dispatch('saveUser', userData);
};

export const getUser = () => {
    return store.dispatch('getUser');
};

export const deleteUser = () => {
    return store.dispatch('deleteUser');
};