import storageService from "@/service/storageService";
import userService from "@/service/userService";
import {useToast} from "vue-toastification";
import {useStore} from "vuex";
import {useRouter} from "vue-router";

const toast = useToast();
const store = useStore();
const router = useRouter();



const userModule = {
    namespaced: true,
    state: {
        token: storageService.get(storageService.USER_TOKEN),
        userinfo: storageService.get(storageService.USER_INFO) ? JSON.parse(storageService.get(storageService.USER_INFO)) : null,
    },
    mutations: {
        SET_TOKEN(state, token) {
            // * 更新本地缓存
            storageService.set(storageService.USER_TOKEN, token);
            // * 更新state
            state.token = token;
        },
        SET_USERINFO(state, userinfo) {
            // * 更新本地缓存
            storageService.set(storageService.USER_INFO, JSON.stringify(userinfo));
            // * 更新state
            state.userinfo = userinfo;
        },
    },

    actions: {
        register(context, {name, telephone, password}) {
            // * 注册逻辑
            return new Promise((resolve, reject) => {
                userService.register({name, telephone, password}).then((res) => {
                    toast.success('注册成功');
                    // * 保存token
                    context.commit('SET_TOKEN', res.data.data.token);
                    return userService.info();
                }).then(res => {
                    // * 保存用户信息
                    context.commit('SET_USERINFO', res.data.data.user);
                    resolve(res);
                }).catch(err => {
                    reject(err);
                });
            });
        },
        login(context, {telephone, password}) {
            // * 注册逻辑
            return new Promise((resolve, reject) => {
                userService.login({telephone, password}).then((res) => {
                    toast.success('登录成功');
                    // * 保存token
                    context.commit('SET_TOKEN', res.data.data.token);
                    return userService.info();
                }).then(res => {
                    // * 保存用户信息
                    context.commit('SET_USERINFO', res.data.data.user);
                    resolve(res);
                }).catch(err => {
                    reject(err);
                });
            });
        },
        logout({commit}) {
            commit('SET_TOKEN', '');
            storageService.set(storageService.USER_TOKEN, '');
            commit('SET_USERINFO', '');
            storageService.set(storageService.USER_INFO, '');

            window.location.reload();
        }
    }

}

export default userModule;
