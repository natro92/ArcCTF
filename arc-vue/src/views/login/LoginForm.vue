<script setup>
import {ref, reactive, inject, computed} from 'vue';
import {useRouter} from "vue-router";
import {useToast} from "vue-toastification";
import {createNamespacedHelpers, useStore, mapActions} from "vuex";

// * 响应式对象存储用户信息
const userInfo = reactive({
    name: '',
    telephone: '',
    password: ''
});
const phoneNumberError = ref(false);
const userNameError = ref(false);

const toast = useToast();
const axios = inject('axios');
const router = useRouter();
const store = useStore();

function login() {
    // * 验证数据
    if (phoneNumberError.value) {
        return;
    }
    // * 请求
    store.dispatch('userModule/login', userInfo).then(() => {
        router.replace({name: 'home'});
    }).catch(err => {
        const message = err.response && err.response.data && err.response.data.msg ? err.response.data.msg : 'An unknown error occurred';
        toast.error(message)
    });
}

// * 检测手机号是否满足11位
function validatePhoneNumber() {
    phoneNumberError.value = userInfo.telephone.length !== 11;
}

</script>

<template>
    <div class="login">
        <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center w-full <!-- dark:bg-gray-950-->">
            <div class="bg-white <!--dark:bg-gray-900--> shadow-md rounded-lg px-8 py-6 max-w-md w-96">
                <h1 class="text-2xl font-bold text-center mb-4 <!--dark:text-gray-200-->">登 录</h1>
                <form action="#">
                    <div class="mb-2">
                        <label for="telephone"
                               class="text-left block text-sm font-medium text-gray-700 <!--dark:text-gray-300--> mb-2">Tel
                            number</label>
                        <input type="tel" id="telephone" v-model="userInfo.telephone"
                               class="shadow-sm rounded-md w-full px-3 py-2 border border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                               autocomplete="current-password" placeholder="+86" @input="validatePhoneNumber" required>
                        <div class="text-sm mt-2 text-red-500" v-if="phoneNumberError">手机号必须为11位</div>
                    </div>
                    <div class="mb-4">
                        <label for="password"
                               class="text-left block text-sm font-medium text-gray-700 <!--dark:text-gray-300--> mb-2">Password</label>
                        <input type="password" id="password" v-model="userInfo.password"
                               class="shadow-sm rounded-md w-full px-3 py-2 border border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                               autocomplete="current-password" placeholder="Enter your password" required>
                        <a href="#"
                           class="text-xs text-gray-600 hover:text-indigo-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">Forgot
                            Password?</a>
                    </div>
                    <div class="flex items-center justify-between mb-4">
                        <div class="flex items-center">
                            <input type="checkbox" id="remember"
                                   class="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 focus:outline-none"
                                   autocomplete="current-password" checked>
                            <label for="remember" class="ml-2 block text-sm text-gray-700 <!--dark:text-gray-300-->">Remember
                                me</label>
                        </div>
                        <a href="/register"
                           class="text-xs text-indigo-500 hover:text-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">Create
                            Account</a>
                    </div>
                    <button @click.prevent="login"
                            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                        Login
                    </button>
                </form>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">

</style>
